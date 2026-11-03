# Service Unit Test Guideline

## Quick Start: Dependency Detection

Detect these key dependencies before generating tests:

{% raw %}```bash
# Core detection
SERVICE_NAME=$(grep "^module" go.mod | cut -d'/' -f3)
HAS_REDIS=$(grep -c "redis.UniversalClient\|dcache.DCache" internal/svc/*.go || echo 0)
HAS_GRPC_CLIENTS=$(grep -E "Client\s+\*" internal/svc/svc.go | grep -v redis | wc -l)
USES_CACHE=$(grep -l "-- -- cache:" internal/repos/*/query.sql | wc -l)
```{% endraw %}

| Component | Detection | Impact |
| Service Name | `go.mod` module path | Use in test DB name, env vars |
| Redis/Cache | `redis.UniversalClient` or `dcache.DCache` | Include cache setup, Redis cleanup |
| gRPC Clients | Parameters ending with `Client` | Create mocks with CallOption |
| Cache Queries | `-- -- cache:` in query.sql | Requires cache infrastructure |

## Critical Setup Points

### 1. Database Connection
{% raw %}```go
// ✅ ALWAYS use Pool.WConn() for repos
suite.dbClient = repos.NewClient(suite.Pool.WConn(), cache, redis, logger)
```{% endraw %}

### 2. gRPC Mock Pattern
{% raw %}```go
type Mock{SERVICE}Client struct {
    mock.Mock
}

func (m *Mock{SERVICE}Client) Get{RESOURCE}(
    ctx context.Context,
    req *pb.Get{RESOURCE}Request,
    opts ...grpc.CallOption,  // REQUIRED for gRPC!
) (*pb.Get{RESOURCE}Response, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*pb.Get{RESOURCE}Response), args.Error(1)
}
```{% endraw %}

### 3. Environment Setup
{% raw %}```go
func init() {
    if env.IsUnitTest() {
        _ = os.Setenv("POSTGRES_PASSWORD", "changeme")
        _ = os.Setenv("POSTGRES_APPNAME", "{SERVICE_NAME}")
        _ = os.Setenv("POSTGRES_SSLMODE", "disable")
        _ = os.Setenv("{SERVICE_NAME_UPPER}_DEBUG", "true")

        // Add Redis if detected
        if HAS_REDIS {
            _ = os.Setenv("REDIS_HOST", "localhost")
            _ = os.Setenv("REDIS_PORT", "6379")
        }
    }
}
```{% endraw %}

## Test File Structure

{% raw %}```go
package svc

import (
    "context"
    "os"
    "testing"
    "time"

    // Core dependencies
    "github.com/kelseyhightower/envconfig"
    "github.com/rs/zerolog"
    "github.com/stretchr/testify/suite"
    "github.com/stumble/wpgx/testsuite"
    "github.com/ggsrc/gglib/env"
    "github.com/ggsrc/gglib/rpcutil"

    // Project imports
    "{MODULE_PATH}/internal/config"
    "{MODULE_PATH}/internal/repos"

    // Conditional imports
    // "github.com/redis/go-redis/v9"      // If Redis detected
    // "github.com/stumble/dcache"         // If cache detected
    // "github.com/stretchr/testify/mock"  // If mocks needed
)

type SvcTestSuite struct {
    *testsuite.WPgxTestSuite
    svc      *Service
    dbClient *repos.Client

    // Optional based on detection
    resource *rpcutil.Resource        // If Redis/cache
    redisCli redis.UniversalClient    // If Redis
    cache    *dcache.DCache           // If cache
    // mock{Service}Client *Mock{Service}Client  // If gRPC
}
```{% endraw %}

## Initialization Patterns

### Pattern A: With Redis/Cache

{% raw %}```go
func newSvcTestSuite() *SvcTestSuite {
    var wpgxConfig wpgx.Config
    envconfig.MustProcess("postgres", &wpgxConfig)
    wpgxConfig.DBName = "{SERVICE_NAME}_test_db"
    wpgxConfig.MaxConns = 1000

    resource := rpcutil.MustInitResource(context.Background(),
        rpcutil.WithCacheInit(),  // Redis + cache
        rpcutil.WithAppName("{SERVICE_NAME}"),
        rpcutil.WithDebug(true),
    )

    return &SvcTestSuite{
        WPgxTestSuite: testsuite.NewWPgxTestSuiteFromConfig(
            &wpgxConfig,
            "{SERVICE_NAME}_test_db",
            []string{{TABLE}.Schema},  // Add all table schemas
        ),
        resource: resource,
        redisCli: resource.RedisClient,
        cache:    resource.DCache,
    }
}
```{% endraw %}

### Pattern B: Without Cache (Database Only)

{% raw %}```go
func newSvcTestSuite() *SvcTestSuite {
    var wpgxConfig wpgx.Config
    envconfig.MustProcess("postgres", &wpgxConfig)
    wpgxConfig.DBName = "{SERVICE_NAME}_test_db"
    wpgxConfig.MaxConns = 1000

    return &SvcTestSuite{
        WPgxTestSuite: testsuite.NewWPgxTestSuiteFromConfig(
            &wpgxConfig,
            "{SERVICE_NAME}_test_db",
            []string{{TABLE}.Schema},
        ),
        // No Redis/cache components
    }
}
```{% endraw %}

## Setup and Teardown

### With Cache
{% raw %}```go
func (suite *SvcTestSuite) SetupTest() {
    suite.WPgxTestSuite.SetupTest()  // Reset DB tables

    var c config.Config
    envconfig.MustProcess("{SERVICE_NAME}", &c)

    logger := zerolog.New(os.Stderr)
    suite.dbClient = repos.NewClient(
        suite.Pool.WConn(),  // IMPORTANT: Use WConn()
        suite.cache,
        suite.redisCli,
        logger,
    )

    // Create service (match actual NewService signature)
    suite.svc = NewService(
        suite.dbClient,
        suite.cache,
        suite.redisCli,
        &c,
    )

    if suite.resource != nil {
        suite.resource.OK(context.Background())
    }
}

func (suite *SvcTestSuite) AfterTest(suiteName, testName string) {
    if suite.redisCli != nil {
        suite.redisCli.FlushDB(context.Background())
    }
}
```{% endraw %}

### Without Cache
{% raw %}```go
func (suite *SvcTestSuite) SetupTest() {
    suite.WPgxTestSuite.SetupTest()

    var c config.Config
    envconfig.MustProcess("{SERVICE_NAME}", &c)

    logger := zerolog.New(os.Stderr)
    suite.dbClient = repos.NewClient(
        suite.Pool.WConn(),
        nil,  // No cache
        nil,  // No Redis
        logger,
    )

    suite.svc = NewService(
        suite.dbClient,
        &c,
    )
}
```{% endraw %}

## Test Entry Point

{% raw %}```go
func TestSvcTestSuite(t *testing.T) {
    suite.Run(t, newSvcTestSuite())
}
```{% endraw %}

## Core Test Patterns

### 1. Golden Test Pattern
**When**: Testing database state changes
**Requires**: `serde.go` in table directory

{% raw %}```go
// internal/repos/{TABLE}/serde.go
type {table}TableSerde struct {
    {table} *Queries
}

func (t {table}TableSerde) Load(data []byte) error {
    err := t.{table}.Load(context.Background(), data)
    if err != nil {
        return err
    }
    // Reset auto-increment sequence if needed
    return t.{table}.RefreshID(context.Background())
}

func (t {table}TableSerde) Dump() ([]byte, error) {
    return t.{table}.Dump(context.Background(), func(m *{Model}) {
        // Normalize timestamps for stable golden files
        if m.CreatedAt != nil {
            *m.CreatedAt = time.Unix(0, 0).UTC()
        }
        if m.UpdatedAt != nil {
            *m.UpdatedAt = time.Unix(0, 0).UTC()
        }
    })
}
```{% endraw %}

{% raw %}```go
func (suite *SvcTestSuite) Test{METHOD}() {
    ctx := context.Background()

    // Load initial state
    {table}Serde := {table}TableSerde{
        {table}: suite.svc.{table}Repo,
    }
    suite.LoadState("TestSvcTestSuite/Test{METHOD}.{table}.input.json", {table}Serde)

    // Execute operation
    resp, err := suite.svc.{METHOD}(ctx, &pb.{REQUEST}{/* fields */})
    suite.NoError(err)

    // Verify DB state
    suite.Golden("{table}_table", {table}Serde)

    // Verify response
    suite.GoldenVarJSON("response", resp)
}

// Update golden files: go test ./internal/svc -update
```{% endraw %}

**File structure:**
{% raw %}```
testdata/TestSvcTestSuite/
├── Test{Method}.{table}.input.json     # Initial state
└── Test{Method}/
    ├── {table}_table.golden            # Expected DB state
    └── response.var.golden             # Expected response
```{% endraw %}

### 2. Table-Driven Test
**When**: Testing multiple scenarios/edge cases

{% raw %}```go
func (suite *SvcTestSuite) Test{METHOD}() {
    testCases := []struct {
        name        string
        input       *pb.{REQUEST}
        expected    *pb.{RESPONSE}
        wantErr     bool
        errContains string
    }{
        {
            name:     "valid_input",
            input:    &pb.{REQUEST}{/* valid */},
            expected: &pb.{RESPONSE}{/* expected */},
        },
        {
            name:        "missing_field",
            input:       &pb.{REQUEST}{/* invalid */},
            wantErr:     true,
            errContains: "required",
        },
    }

    for _, tc := range testCases {
        suite.Run(tc.name, func() {
            result, err := suite.svc.{METHOD}(context.Background(), tc.input)
            if tc.wantErr {
                suite.Error(err)
                if tc.errContains != "" {
                    suite.Contains(err.Error(), tc.errContains)
                }
            } else {
                suite.NoError(err)
                suite.Equal(tc.expected, result)
            }
        })
    }
}
```{% endraw %}

### 3. CRUD Test
**When**: Testing Create, Read, Update, Delete lifecycle

{% raw %}```go
func (suite *SvcTestSuite) TestCRUD{MODEL}() {
    ctx := context.Background()

    // CREATE
    createResp, err := suite.svc.Create{MODEL}(ctx, &pb.Create{MODEL}Request{
        Name: "Test Item",
    })
    suite.NoError(err)
    id := createResp.{MODEL}.Id

    // READ
    getResp, err := suite.svc.Get{MODEL}(ctx, &pb.Get{MODEL}Request{Id: id})
    suite.NoError(err)
    suite.Equal("Test Item", getResp.{MODEL}.Name)

    // UPDATE
    _, err = suite.svc.Update{MODEL}(ctx, &pb.Update{MODEL}Request{
        Id:   id,
        Name: "Updated",
    })
    suite.NoError(err)

    getResp, err = suite.svc.Get{MODEL}(ctx, &pb.Get{MODEL}Request{Id: id})
    suite.NoError(err)
    suite.Equal("Updated", getResp.{MODEL}.Name)

    // DELETE
    _, err = suite.svc.Delete{MODEL}(ctx, &pb.Delete{MODEL}Request{Id: id})
    suite.NoError(err)

    _, err = suite.svc.Get{MODEL}(ctx, &pb.Get{MODEL}Request{Id: id})
    suite.Error(err)  // Should not exist
}
```{% endraw %}

### 4. Error Handling Test
**When**: Always - test error paths for every method

{% raw %}```go
func (suite *SvcTestSuite) TestErrors{METHOD}() {
    ctx := context.Background()

    // Not found
    _, err := suite.svc.Get{MODEL}(ctx, &pb.Get{MODEL}Request{
        Id: 99999,
    })
    suite.Error(err)
    suite.Contains(err.Error(), "not found")

    // Validation error
    _, err = suite.svc.Create{MODEL}(ctx, &pb.Create{MODEL}Request{
        // Missing required field
    })
    suite.Error(err)
    suite.Contains(err.Error(), "required")

    // Duplicate constraint
    _, err = suite.svc.Create{MODEL}(ctx, &pb.Create{MODEL}Request{
        Email: "duplicate@example.com",
    })
    suite.Error(err)
    // Check for PostgreSQL error code 23505 (unique_violation)

    // Context timeout
    ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
    defer cancel()
    _, err = suite.svc.{METHOD}(ctxTimeout, &pb.{REQUEST}{})
    suite.Error(err)
    suite.True(errors.Is(err, context.DeadlineExceeded))
}
```{% endraw %}

### 5. Mock External Service Test
**When**: Method calls external gRPC services

{% raw %}```go
func (suite *SvcTestSuite) TestWithMock{METHOD}() {
    ctx := context.Background()

    // Setup mock expectation
    suite.mock{SERVICE}Client.On("Get{RESOURCE}", ctx, mock.Anything).Return(
        &pb.Get{RESOURCE}Response{
            {RESOURCE}: &pb.{RESOURCE}{
                Id:   "123",
                Name: "Mocked",
            },
        }, nil,
    ).Once()

    // Execute and verify
    result, err := suite.svc.{METHOD}(ctx, "123")
    suite.NoError(err)
    suite.Equal("Mocked", result.Name)

    suite.mock{SERVICE}Client.AssertExpectations(suite.T())
}
```{% endraw %}

## Pattern Selection Guide

| Scenario | Test Pattern |
|----------|-------------|
| Database state changes | Golden Test |
| Multiple scenarios/edge cases | Table-Driven Test |
| CRUD operations | CRUD Test |
| Error conditions | Error Handling Test |
| External service calls | Mock Test |

## Test Implementation Checklist

### Step 1: Detect Dependencies
{% raw %}```bash
# Check for Redis/Cache
grep -c "redis.UniversalClient\|dcache.DCache" internal/svc/svc.go

# Check for cache in queries
grep -l "-- -- cache:" internal/repos/*/query.sql

# Check NewService signature
grep "func NewService" internal/svc/svc.go
```{% endraw %}

### Step 2: Choose Pattern
- **Has Redis/Cache** → Use Pattern A with rpcutil.Resource
- **No Redis/Cache** → Use Pattern B (database only)
- **Has gRPC Clients** → Add mock implementations

### Step 3: Verify Setup
1. Check `repos.NewClient` signature
2. Verify Pool has `WConn()` method
3. Match test data with actual schema
4. Start with simple test first

## Running Tests

{% raw %}```bash
# Set required environment
export POSTGRES_PASSWORD=changeme
export POSTGRES_APPNAME=service-name
export POSTGRES_SSLMODE=disable

# Run all tests
go test ./internal/svc -v

# Run specific test
go test ./internal/svc -run TestSvcTestSuite/Test{METHOD} -v

# Update golden files
go test ./internal/svc -update

# Check coverage
go test -coverprofile=coverage.out ./internal/svc
go tool cover -html=coverage.out
```{% endraw %}

## Common Pitfalls

### Database Connection
{% raw %}```go
// ❌ WRONG
suite.dbClient = repos.NewClient(suite.Pool, ...)

// ✅ CORRECT
suite.dbClient = repos.NewClient(suite.Pool.WConn(), ...)
```{% endraw %}

### gRPC Mocks
{% raw %}```go
// ❌ WRONG - Missing CallOption
func (m *MockClient) GetResource(ctx context.Context, req *pb.Request) (*pb.Response, error)

// ✅ CORRECT
func (m *MockClient) GetResource(ctx context.Context, req *pb.Request, opts ...grpc.CallOption) (*pb.Response, error)
```{% endraw %}

### Cache Verification
{% raw %}```go
// ❌ WRONG - Direct DB check
record, err := suite.dbClient.GetByID(ctx, id)

// ✅ CORRECT - Through API
result, err := suite.svc.GetModel(ctx, req)
```{% endraw %}

## Troubleshooting

### Redis Connection Failed
{% raw %}```go
// If service doesn't need Redis:
// Use Pattern B (no cache)

// If Redis optional for tests:
if suite.cache == nil {
    suite.T().Skip("Redis not available")
}
```{% endraw %}

### Repository Requires Cache
{% raw %}```go
// Create no-op cache
type NoOpCache struct{}

func (n *NoOpCache) Get(ctx context.Context, key string) ([]byte, error) {
    return nil, errors.New("cache disabled")
}

func (n *NoOpCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    return nil
}
```{% endraw %}

## Golden Files

### Workflow
{% raw %}```bash
# Generate/update golden files
go test ./internal/svc -update

# Review changes
git diff internal/svc/testdata/

# Commit with code changes
git add internal/svc/svc.go internal/svc/testdata/
git commit -m "feat: update with golden files"
```{% endraw %}

### Debugging
{% raw %}```go
func (suite *SvcTestSuite) TestWithDebug() {
    // Log intermediate state
    suite.T().Logf("Data: %+v", testData)

    // Use JSONEq for better diff
    suite.JSONEq(expected, actual)
}
```{% endraw %}

## Best Practices

### DO:
- Detect dependencies before generating tests
- Use `Pool.WConn()` for database connections
- Include `grpc.CallOption` in all mock methods
- Verify through API responses, not direct DB queries
- Match test data with schema requirements
- Use golden tests for state verification

### DON'T:
- Assume all services need Redis
- Force unnecessary dependencies
- Use Pool directly as WGConn
- Forget to clean Redis after tests
- Mix inconsistent field values
- Hardcode error messages

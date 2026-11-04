package svc

import (
	"context"
	"os"
	"testing"

	"github.com/ggsrc/gglib/env"
	ggcache "github.com/ggsrc/gglib/resource/cache"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/suite"
	"github.com/stumble/wpgx/testsuite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"

	"{{ cookiecutter.go_module_path.strip('/') }}/internal/config"
	"{{ cookiecutter.go_module_path.strip('/') }}/internal/repos/users"
)

// Test container images. Keep them in step with docker-compose.yml and with
// the versions the production database and Redis run.
const (
	postgresImage = "postgres:18.6"
	redisImage    = "redis:8.10.3-alpine"
)

func init() {
	if env.IsUnitTest() {
		// Set environment variables for service configuration
		// NewWPgxTestSuiteTcDefault automatically handles PostgreSQL container setup
		_ = os.Setenv("{{ cookiecutter.project_name.strip().upper() }}_DEBUG", "true")
	}
}

// SvcTestSuite embeds WPgxTestSuite which provides:
// - suite.Suite (testify) with assertion methods (NoError, Error, Equal, NotNil, etc.)
// - Database test utilities (GetPool, SetupTest, etc.)
type SvcTestSuite struct {
	*testsuite.WPgxTestSuite

	cache          *ggcache.Cache
	svc            *Service
	redisContainer testcontainers.Container
}

func newSvcTestSuite() *SvcTestSuite {
	// Create WPgxTestSuite using NewWPgxTestSuiteTcDefault (recommended since v0.4)
	// This automatically creates a dedicated PostgreSQL container for this test suite,
	// allowing packages to run in parallel without conflicts
	// Note: Redis container will be started in SetupSuite() for proper initialization order
	return &SvcTestSuite{
		WPgxTestSuite: testsuite.NewWPgxTestSuiteTcDefault(
			[]string{
				users.Schema, // wpgx will auto-apply schema
			},
			testsuite.WithPostgresImage(postgresImage),
		),
	}
}

// newCache reads REDIS_* and DCACHE_*, so call it after the Redis container
// has set REDIS_HOST and REDIS_PORT.
func newCache() *ggcache.Cache {
	var redisConfig ggcache.RedisConfig
	envconfig.MustProcess("redis", &redisConfig)
	var dcacheConfig ggcache.DCacheConfig
	envconfig.MustProcess("dcache", &dcacheConfig)
	return ggcache.NewCacheWithOptions(
		ggcache.WithAppName("{{ cookiecutter.project_name.strip() }}"),
		ggcache.WithRedisConfig(&redisConfig),
		ggcache.WithDCacheConfig(&dcacheConfig),
	)
}

func (suite *SvcTestSuite) setupRedisContainer() error {
	// A Redis started outside the tests (REDIS_HOST / REDIS_PORT) is used as is.
	if os.Getenv("REDIS_HOST") != "" {
		return nil
	}

	// Otherwise start one with testcontainers (wpgx doesn't manage Redis), in
	// CI as on a laptop.
	ctx := context.Background()
	redisContainer, err := redis.Run(ctx, redisImage)
	if err != nil {
		return err
	}

	// Get Redis connection details
	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		return err
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		return err
	}

	// Set Redis environment variables for ggcache
	_ = os.Setenv("REDIS_HOST", redisHost)
	_ = os.Setenv("REDIS_PORT", redisPort.Port())

	// Store container for cleanup
	suite.redisContainer = redisContainer
	return nil
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSvcTestSuite(t *testing.T) {
	suite.Run(t, newSvcTestSuite())
}

// SetupSuite runs once before all tests
func (suite *SvcTestSuite) SetupSuite() {
	// Step 1: Start PostgreSQL container (managed by wpgx)
	// This will update the config with container connection details
	suite.WPgxTestSuite.SetupSuite()

	// Step 2: Get updated config with container connection details
	// In container mode, GetConfig() returns the updated config with actual container host/port
	pgConfig := suite.GetConfig()
	suite.T().Logf("PostgreSQL container: %s:%d/%s", pgConfig.Host, pgConfig.Port, pgConfig.DBName)

	// Step 3: Start Redis container (manually managed)
	// This must be done before initializing cache, as cache needs Redis connection
	err := suite.setupRedisContainer()
	suite.Require().NoError(err, "Failed to setup Redis container")

	// Step 4: Create and initialize the cache once Redis is up: its
	// configuration is read from the variables set in step 3.
	suite.cache = newCache()
	err = suite.cache.Init(context.Background())
	suite.Require().NoError(err, "Failed to initialize cache")
}

// SetupTest runs before each test
func (suite *SvcTestSuite) SetupTest() {
	// wpgx automatically resets database for each test
	// This ensures test isolation
	suite.WPgxTestSuite.SetupTest()

	var c config.Config
	envconfig.MustProcess("{{ cookiecutter.project_name.strip().replace("-", "_") }}", &c)

	// Create new service instance for each test
	// suite.Pool is available after SetupTest() resets the database
	// GetPool() internally uses GetConfig() to ensure correct container connection details
	suite.svc = NewService(
		suite.GetPool(),              // GetPool() uses GetConfig() internally for container connection
		suite.cache.GetDCache(),      // DCache instance
		suite.cache.GetRedisClient(), // Redis client (manually managed)
		&c,
	)
}

func (suite *SvcTestSuite) AfterTest(suiteName, testName string) {
	// Clear redis cache after each test to ensure test isolation
	if suite.cache != nil && suite.cache.GetRedisClient() != nil {
		ctx := context.Background()
		suite.cache.GetRedisClient().FlushDB(ctx)
	}
}

// TearDownSuite runs once after all tests
func (suite *SvcTestSuite) TearDownSuite() {
	// Step 1: Clean up Redis container if it exists
	// This should be done before cleaning up PostgreSQL to ensure proper shutdown order
	if suite.redisContainer != nil {
		ctx := context.Background()
		if err := suite.redisContainer.Terminate(ctx); err != nil {
			suite.T().Logf("Failed to terminate Redis container: %v", err)
		}
		suite.redisContainer = nil
	}

	if suite.cache != nil {
		_ = suite.cache.Stop(context.Background())
	}

	// Step 2: Clean up PostgreSQL container (managed by wpgx)
	suite.WPgxTestSuite.TearDownSuite()
}

// Note: TestAddUser has been replaced with comprehensive table-driven test: TestAddUserRPC

// Note: TestGetUser, TestUpdateUser, and TestDeleteUser have been replaced
// with comprehensive table-driven tests: TestGetUserRPC, TestUpdateUserRPC, and TestDeleteUserRPC

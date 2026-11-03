<p align="center">
    <h1 align="center">{{ cookiecutter.project_name.strip() }}</h1>
</p>
<p align="center">
    <em>{{ cookiecutter.project_description }}</em>
</p>
<hr>

[![goSec](https://img.shields.io/badge/Security-GoSec-lightgrey?style=for-the-badge&logo=springsecurity)](https://github.com/securego/gosec)

## 🚀 Quick Start (5 minutes)

```bash
# 1. Clone and setup
git clone https://{{ cookiecutter.go_module_path.strip('/') }}
cd {{ cookiecutter.project_name.strip() }}

# 2. Start dependencies (PostgreSQL + Redis)
make docker-up

# 3. Generate code
make proto     # Generate protobuf code
make repos     # Generate database code from SQL

# 4. Run the service
make local-run

# 5. Test it works
grpcurl -plaintext localhost:9090 {{ cookiecutter.project_name.strip().replace('-', '_') }}.v1.GService/Ping
```

## 📚 Documentation

**New to the project?** Start here → [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md)

| Task | Guide | Quick Command |
|------|-------|---------------|
| 🗄️ **Add database table** | [REPOSITORY_GUIDE.md](./internal/repos/REPOSITORY_GUIDE.md) | `make repos` |
| ✅ **Write service tests** | [SERVICE_TEST_GUIDE.md](./SERVICE_TEST_GUIDE.md) | `make test` |
| 🐛 **Debug issues** | [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md#getting-help) | `make lint` |

## 💻 Development Commands

### Essential Commands
```bash
make build        # Build all binaries → bin/server
make test         # Run all tests with coverage
make lint         # Run linters (golangci-lint, buf)
make format       # Format code (go, sql, proto)
```

### Code Generation
```bash
make proto        # Generate from *.proto files → api/go/*/v1/*.pb.go
make repos        # Generate from *.sql files → internal/repos/*/*.go
make mock         # Generate mocks for testing
```

### Running Locally
```bash
make docker-up    # Start PostgreSQL + Redis
make local-run    # Run gRPC server on :9090

# Clean up
make docker-down
```

## 📂 Project Structure

```
.
├── api/                          # API definitions
│   └── proto/*/v1/*.proto       # Protobuf service definitions
├── cmd/
│   └── server/main.go           # gRPC server entry point
├── internal/
│   ├── svc/                     # Service layer (business logic)
│   │   ├── *.go                # Service methods
│   │   └── *_test.go           # Tests with golden files
│   ├── repos/                   # Database layer (sqlc generated)
│   │   ├── README.md            # Quick start guide
│   │   ├── REPOSITORY_GUIDE.md  # Detailed database guide
│   │   └── {table}/
│   │       ├── schema.sql      # Table schema
│   │       ├── query.sql       # SQL queries (sqlc)
│   │       └── *.go            # Generated code
│   └── config/                  # Configuration
├── DEVELOPMENT_GUIDE.md         # 👈 Development index (start here)
├── PROJECT_SETUP_GUIDE.md       # Project initialization guide
├── SERVICE_TEST_GUIDE.md        # Testing guide
└── Makefile                     # All commands
```

## 🎯 Common Tasks

### 1. Add a New API Endpoint

```bash
# Step 1: Define in proto
vim api/proto/{{ cookiecutter.project_name.strip() }}/v1/*.proto

# Step 2: Generate code
make proto

# Step 3: Implement in service
vim internal/svc/my_feature.go

# Step 4: Write tests
vim internal/svc/my_feature_test.go
make test
```

See [SERVICE_TEST_GUIDE.md](./SERVICE_TEST_GUIDE.md) for testing patterns.

### 2. Add a Database Table

```bash
# Step 1: Create schema
vim internal/repos/mytable/schema.sql

# Step 2: Write queries
vim internal/repos/mytable/query.sql

# Step 3: Generate Go code
make repos

# Step 4: Use in service
# Import: "{{ cookiecutter.go_module_path.strip('/') }}/internal/repos/mytable"
```

See [REPOSITORY_GUIDE.md](./internal/repos/REPOSITORY_GUIDE.md) for SQL patterns.

### 3. Debug Issues

```bash
# Check linter errors
make lint

# Check specific file
golangci-lint run internal/svc/myfile.go

# Run specific test
go test -v ./internal/svc -run TestMyFeature

# Check database queries
make repos  # Will show sqlc errors
```

## 🔧 Environment Variables

Required for local development:

```bash
# Database
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=my-secret
export POSTGRES_DBNAME={{ cookiecutter.project_name.strip() }}_db
export POSTGRES_APPNAME={{ cookiecutter.project_name.strip() }}

# Redis
export REDIS_HOST=localhost
export REDIS_PORT=6379

# Service
export {{ cookiecutter.project_name.strip().upper().replace('-', '_') }}_DEBUG=true

# gglib: environment, logs, ports (see .env.example)
export ENV=local
export GRPC_PORT=9090
```

💡 **Tip**: Use `direnv` to auto-load these from `.envrc`

## 🧪 Testing

```bash
# Run all tests
make test

# Run specific package
go test ./internal/svc/...

# Run with coverage
go test -cover ./...

# Update golden files
go test ./internal/svc -update

# Run integration tests
make test-integration
```

See [SERVICE_TEST_GUIDE.md](./SERVICE_TEST_GUIDE.md) for:
- How to write tests with golden files
- Database testing patterns
- Mock usage
- Best practices

## 🐳 Docker

```bash
# Start all services
make docker-up

# View logs
docker compose logs -f

# Stop all services
make docker-down
```

Services available:
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`
- gRPC Server: `localhost:9090`

## 📊 Monitoring & Observability

- **Metrics**: Prometheus metrics at `:4014/metrics`
- **Health**: `/health/ready` and `/health/alive` on `:8080` (gglib/health): not ready while Postgres or Redis is down, restarted only if stuck
- **Logs**: JSON with `trace_id` and `span_id` in clusters, readable output locally (`LOG_LEVEL`, `LOG_FORMAT`)
- **Tracing**: OpenTelemetry over OTLP when `OTEL_EXPORTER_OTLP_ENDPOINT` is set

## 🚀 CI/CD & GitHub Actions

This project includes pre-configured GitHub Actions workflows:

### Available Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `check.yml` | Push/PR | Run tests and linting |
| `lint.yml` | Push/PR | Run golangci-lint |
| `deploy-beta.yml` | Manual | Deploy to beta environment |
| `deploy-stg.yml` | Manual | Deploy to staging |
| `deploy-prd.yml` | Manual | Deploy to production |
| `deploy-hotfix.yml` | Manual | Emergency hotfix deployment |
| `cruft-update.yml` | Schedule | Check for template updates |

### Setup

1. **Enable Workflows**: Workflows are ready to use
2. **Add Secrets**: Configure repository secrets
   - Add deployment credentials as needed
3. **Customize**: Edit `.github/workflows/*.yml` for your environment

### Dependabot

Automatic dependency updates are enabled:
- Weekly Go module updates
- GitHub Actions updates  
- Docker base image updates

To disable: Remove `.github/dependabot.yml`

## 🚨 Troubleshooting

### Proto generation fails
```bash
# Reinstall buf
ONLY=buf scripts/install-buf-deps.sh

# Or manually
brew upgrade buf
```

### sqlc generation fails
```bash
# Check SQL syntax
cat internal/repos/{table}/query.sql

# Regenerate
make repos
```

### Tests fail
```bash
# Make sure Docker is running (testcontainers requires it)
docker info >/dev/null

# Run tests
make test
```

### Lint errors
```bash
# Auto-fix formatting
make format

# Check what's wrong
make lint

# Fix imports
goimports -w .
```

## 🔗 Related Resources

- **[DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md)** - Complete development workflow
- **service-demo** - Reference implementation (see examples in this project)
- **gglib** - Shared Go libraries for common functionality

## 📝 Code Examples

### Call the service (gRPC)
```go
conn, _ := grpc.Dial("localhost:9090", grpc.WithInsecure())
client := pb.NewGServiceClient(conn)
resp, _ := client.Ping(context.Background(), &emptypb.Empty{})
```

### Query database
```go
// In internal/svc/myservice.go
records, err := s.projectsRepo.GetAllProjects(ctx)
if err != nil {
    return nil, err
}
```

### Add cache
```go
// Cache for 5 minutes
err := s.cache.Set(ctx, "key", value, 5*time.Minute)

// Get from cache
var value MyStruct
err := s.cache.Get(ctx, "key", &value)
```

## 🤝 Contributing

1. Read [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md)
2. Create feature branch: `git checkout -b feature/my-feature`
3. Write tests: Follow [SERVICE_TEST_GUIDE.md](./SERVICE_TEST_GUIDE.md)
4. Run checks: `make lint && make test`
5. Submit PR

## 📄 License

[Your License Here]

---

**Questions?** Check [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md) or ask the team!

<!-- Reference Links -->
[makefile-file]: ./Makefile
[github-actions]: ../../actions
[github-releases]: ../../releases
[precommit-config]: ./.pre-commit-config.yaml
[gomod-file]: ./go.mod
[github-actions-tests]: ../../actions/workflows/check.yml
[dependabot-pulls]: ../../pulls?utf8=%E2%9C%93&q=is%3Apr%20author%3Aapp%2Fdependabot

[semver-link]: https://semver.org
[pre-commit]: https://pre-commit.com
[github-repo]: https://github.com/new
[gitlab-repo]: https://gitlab.com/new
[dependabot-link]: https://dependabot.com
[githooks]: https://git-scm.com/docs/githooks
[python-install]: https://www.python.org/downloads
[release-drafter-config]: ./.github/release-drafter.yml
[makefile-official]: https://www.gnu.org/software/make
[pip-install]: https://pip.pypa.io/en/stable/installation
[golangci-install]: https://golangci-lint.run/usage/install
[cookiecutter-link]: https://github.com/cookiecutter/cookiecutter
[release-drafter]: https://github.com/marketplace/actions/release-drafter
[creating-secrets]: https://docs.github.com/en/actions/security-guides/encrypted-secrets#creating-encrypted-secrets-for-a-repository

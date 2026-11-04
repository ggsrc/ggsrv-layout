# ggsrv-layout

<div align="center">

**Production-ready Go gRPC service template with comprehensive guidelines**

A cookiecutter template for building scalable Go microservices with gRPC, PostgreSQL and Redis, built on [gglib](https://github.com/ggsrc/gglib).

[Quick Start](#-quick-start) • [Features](#-features) • [Documentation](#-documentation) • [Examples](#-examples)

</div>

---

## 🚀 Quick Start

### Prerequisites

```bash
# Recommended: Install cruft (supports template updates)
brew install cruft

# Or install cookiecutter (basic usage)
brew install cookiecutter
# pip install cookiecutter
```

### Generate a New Project

**Option 1: Using cruft (Recommended)**

Cruft allows you to keep your project in sync with template updates.

```bash
# Interactive mode
cruft create https://github.com/ggsrc/ggsrv-layout

# Non-interactive mode
cruft create https://github.com/ggsrc/ggsrv-layout \
  --extra-context '{"project_name": "my-service", "org": "ggsrc"}' \
  --no-input
```

**Option 2: Using cookiecutter (Basic)**

```bash
# Interactive mode (recommended for first-time users)
cookiecutter https://github.com/ggsrc/ggsrv-layout

# Non-interactive mode with defaults
cookiecutter https://github.com/ggsrc/ggsrv-layout --no-input \
  project_name=my-service \
  org=ggsrc
```

**Why use cruft?**
- ✅ Track which template version you're using
- ✅ Update your project when template improves
- ✅ Check if your project is out of sync with the template

### What You Get

```
my-service/
├── 📚 Complete documentation suite (6 guides)
├── 🔧 Ready-to-use Makefile commands
├── 🗄️ PostgreSQL integration with sqlc
├── 🚀 gRPC API with Protocol Buffers
├── ✅ Comprehensive test suite
└── 🐳 Docker & Docker Compose configs
```

---

## ✨ Features

### Development Experience

- **🤖 AI-Friendly Documentation** - 6 comprehensive guides designed for both humans and AI assistants
  - PROJECT_SETUP_GUIDE.md - 17 step-by-step TODOs for project initialization
  - DEVELOPMENT_GUIDE.md - Central hub linking all documentation
  - SERVICE_TEST_GUIDE.md - Testing patterns and best practices
  - REPOSITORY_GUIDE.md - Database patterns with sqlc

- **⚡ Instant Productivity** - Everything pre-configured and working out of the box
  - Proto generation with `buf`
  - SQL-to-Go code generation with `sqlc`
  - Fast iteration with `go run`
  - Docker Compose for dependencies

- **🎯 Type Safety** - End-to-end type safety from database to API
  - SQL queries → Go structs (sqlc)
  - Proto definitions → Go code (protoc)
  - Compile-time safety throughout

### Technology Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **API** | gRPC + Protocol Buffers | High-performance RPC with strong typing |
| **Database** | PostgreSQL + sqlc | Type-safe SQL queries |
| **Cache** | Redis + dcache | Distributed caching with automatic invalidation |
| **Testing** | testify + wpgx | Unit & integration tests with golden files |
| **Linting** | golangci-lint | 20+ linters pre-configured |

### Code Quality

- **Automated Formatting** - gofmt, gofumpt, goimports, gci
- **Comprehensive Linting** - golangci-lint with 20+ linters
- **Pre-commit Hooks** - Catch issues before they reach CI
- **Test Coverage** - Coverage reports with golden file testing

### Deployment Ready

- **Docker Support** - Multi-stage builds for minimal image sizes
- **GitHub Actions** - Pre-configured CI/CD workflows
  - Automated testing and linting
  - Multi-environment deployment (beta, staging, production)
  - Cruft template updates
- **Health Checks** - Built-in health check endpoints
- **Observability** - Structured logging with zerolog
- **Resource Management** - Graceful shutdown and connection pooling
- **Dependabot** - Automatic dependency updates

---

## 📚 Documentation

Each generated project includes **7 comprehensive guides**:

### 1. **PROJECT_SETUP_GUIDE.md** (⭐ Start here for new projects)
Complete initialization workflow with 17 actionable TODOs across 7 phases:
- Phase 1: Environment Setup
- Phase 2: Database Setup (schema, queries, serde)
- Phase 3: API Definition (protobuf)
- Phase 4: Service Implementation
- Phase 5: Build and Test
- Phase 6: Run the Service

### 2. **DEVELOPMENT_GUIDE.md** (📍 Central documentation hub)
Index of all guides with quick navigation and best practices summary.

### 3. **SERVICE_TEST_GUIDE.md** (Testing patterns)
- Test suite setup with wpgx
- Golden file testing
- Mock usage
- Database test patterns

### 4. **internal/repos/REPOSITORY_GUIDE.md** (Database layer)
- Schema design patterns
- sqlc query patterns
- Cache integration
- Transaction handling

### 5. **internal/repos/README.md** (Quick reference)
Fast lookup for common database operations.

### 6. **README.md** (Project overview)
Quick start, common tasks, and troubleshooting for the generated project.

---

## 🎯 Template Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `project_name` | `my-service` | Project name (becomes directory name) |
| `org` | `ggsrc` | Organization name (e.g., your GitHub org) |
| `go_module_path` | Auto-generated | Go module path (e.g., github.com/ggsrc/my-service) |
| `go_version` | `1.27.3` | Go version to use (gglib requires 1.27 or later) |
| `base_branch` | `main` | Default git branch |
| `contact_email` | `""` | Contact email (optional) |
| `use_precommit` | `y` | Include pre-commit hooks |

### Example Configuration

```yaml
# cookiecutter-config-file.yml (auto-generated in your project)
project_name: my-awesome-service
org: ggsrc
go_module_path: github.com/ggsrc/my-awesome-service
go_version: 1.27.3
```

---

## 💻 Generated Project Structure

```
my-service/
├── README.md                       # Project overview & quick start
├── DEVELOPMENT_GUIDE.md            # Documentation index
├── PROJECT_SETUP_GUIDE.md          # Initialization guide (AI-friendly)
├── SERVICE_TEST_GUIDE.md           # Testing patterns
│
├── api/
│   └── proto/*/v1/*.proto          # API definitions
│
├── cmd/
│   └── server/main.go              # gRPC server
│
├── pkg/
│   ├── svc/                        # Service layer (business logic)
│   │   ├── *.go                    # Service methods
│   │   └── *_test.go               # Tests with golden files
│   │
│   ├── repos/                      # Database layer
│   │   ├── README.md               # Quick reference
│   │   ├── REPOSITORY_GUIDE.md     # Detailed guide
│   │   └── {table}/
│   │       ├── schema.sql          # Table definitions
│   │       ├── query.sql           # SQL queries (sqlc)
│   │       └── *.go                # Generated code
│   │
│   ├── config/                     # Configuration
│   └── convert/                    # Model converters
│
├── Makefile                        # All commands
├── docker-compose.yml              # Local dependencies
└── Dockerfile                      # Production image
```

---

## 📖 Examples

### Example 1: Generate a New Service

```bash
# Generate project
cookiecutter https://github.com/ggsrc/ggsrv-layout

# Follow the PROJECT_SETUP_GUIDE.md
cd my-service
cat PROJECT_SETUP_GUIDE.md

# Quick start (automated)
make docker-up       # Start PostgreSQL and Redis
make proto           # Generate API code
make repos           # Generate database code
make build           # Build binaries
make test            # Run tests
make test-local-run  # Test local run (start, verify, stop)
make local-run       # Start server
```

### Example 2: AI-Assisted Project Initialization

The PROJECT_SETUP_GUIDE.md is designed for AI assistants to execute:

1. AI reads PROJECT_SETUP_GUIDE.md
2. AI creates 17 TODOs from the guide
3. AI executes each TODO sequentially:
   - Writes schema.sql
   - Writes query.sql
   - Defines protobuf
   - Implements service methods
   - Writes tests
   - Verifies everything works

Result: Fully functional service in ~60 minutes

### Example 3: Add a New Feature

```bash
# 1. Define database table
vim pkg/repos/users/schema.sql
vim pkg/repos/users/query.sql
make repos

# 2. Define API
vim api/proto/my-service/v1/users.proto
make proto

# 3. Implement service
vim pkg/svc/users.go
vim pkg/svc/users_test.go
make test

# 4. Run
make local-run
```

---

## 🛠️ Available Commands

Generated projects include a comprehensive Makefile:

```bash
# Essential
make build          # Build all binaries
make test           # Run tests with coverage
make test-local-run # Test local run (start, verify, stop)
make lint           # Run all linters
make format         # Format all code

# Code Generation
make proto          # Generate from .proto files
make repos          # Generate from .sql files
make mock           # Generate test mocks

# Development
make local-run      # Run gRPC server
make docker-up      # Start dependencies
make docker-down    # Stop dependencies

# Cleaning
make clean          # Remove generated files
```

---

## 🎓 Best Practices

The template enforces industry best practices:

### Database Layer
- ✅ Use sqlc for type-safe SQL
- ✅ Add cache annotations to queries
- ✅ Add timeout annotations
- ✅ Use transactions for multi-table operations
- ❌ Never use string concatenation in SQL

### Service Layer
- ✅ Write tests with golden files
- ✅ Use table-driven tests
- ✅ Mock external dependencies
- ✅ Test error cases
- ❌ Don't test external services directly

## 🚀 CI/CD & GitHub Actions

Generated projects include pre-configured GitHub Actions workflows:

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

### Setup GitHub Actions

1. **Enable Workflows**: Workflows are included by default
2. **Add Secrets**: Configure repository secrets for deployment
   - Add deployment credentials as needed
3. **Customize**: Edit `.github/workflows/*.yml` for your needs

### Dependabot

Automatic dependency updates are enabled by default:
- Weekly checks for Go module updates
- GitHub Actions updates
- Docker base image updates

To disable: Remove `.github/dependabot.yml`

---

## 🔄 Keeping Your Project Updated (cruft)

If you created your project with `cruft`, you can easily sync template improvements:

### Check for Updates

```bash
cd my-service
cruft check

# Output:
# ✅ SUCCESS: Good to go, project's cruft is up to date!
# or
# ⚠️  WARNING: Project's cruft is out of date!
```

### Update Your Project

```bash
cd my-service
cruft update

# cruft will:
# 1. Fetch the latest template
# 2. Generate updated files
# 3. Show you what changed
# 4. Merge changes while preserving your custom code
# 5. Ask you to resolve conflicts (if any)
```

### View Template Info

```bash
cruft info

# Shows:
# - Template URL
# - Current template version (commit hash)
# - Configuration used during generation
```

**Important**: Always commit `.cruft.json` to your repository to track the template version.

---

## 🤝 Contributing

We welcome contributions! The template is designed to be:

- **Opinionated** - Best practices baked in
- **Flexible** - Easy to customize for your needs
- **Well-documented** - Clear guidelines for everything

To contribute:
1. Fork the repository
2. Create a feature branch
3. Test with: `cookiecutter .` or `cruft create .` (local testing)
4. Submit a pull request

---

## 🔗 Related Projects

- **service-demo** - Reference implementation (see examples in generated projects)
- **gglib** - Shared Go libraries for common functionality

---

## 📄 License

This template is licensed under the MIT License. Generated projects can use any license you choose.

---

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/ggsrc/ggsrv-layout/issues)
- **Documentation**: Check the guides in generated projects
- **Examples**: See the comprehensive guides in each generated project

---

<div align="center">

**Built with ❤️ for the Go community**

[⬆ Back to Top](#ggsrv-layout)

</div>

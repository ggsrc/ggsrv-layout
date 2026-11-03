# Development Guide

This document provides an overview of all development guidelines in this project.

## 🚀 New to This Project?

**Start here**: [PROJECT_SETUP_GUIDE.md](./PROJECT_SETUP_GUIDE.md)

This step-by-step guide will help you (or an AI assistant) initialize a fully functional project from scratch in ~60 minutes. It breaks down the setup into clear TODOs that can be followed sequentially.

## 📚 Available Guidelines

### 0. Project Setup Guide (START HERE!)
**File**: `PROJECT_SETUP_GUIDE.md`  
**Purpose**: Step-by-step initialization guide for new projects

**Key Topics**:
- Environment setup
- Database schema and queries
- API definition (protobuf)
- Service implementation
- Testing and verification

**When to use**: When initializing a new project from this template, or when an AI needs to scaffold a complete project

**Quick Start**:
```bash
# Generate project from template
cookiecutter ggsrv-layout/

# Follow the setup guide
cat PROJECT_SETUP_GUIDE.md

# Or let AI follow it step-by-step
# AI will create TODOs and execute them sequentially
```

### 1. Service Testing Guide
**File**: `SERVICE_TEST_GUIDE.md`  
**Purpose**: Complete guide for writing service tests using wpgx test suite

**Key Topics**:
- Test suite setup and initialization
- Database test patterns (CRUD operations)
- Golden file testing
- Mock usage
- Transaction handling
- Best practices

**When to use**: When writing or maintaining service layer tests (`internal/svc/*_test.go`)

**Quick Start**:
```bash
# Read the guide
cat SERVICE_TEST_GUIDE.md

# Run tests
make test
```

### 2. Repository Layer Guide
**Files**: 
- `internal/repos/REPOSITORY_GUIDE.md` - Detailed guidelines
- `internal/repos/README.md` - Quick reference

**Purpose**: Guide for working with database repositories using sqlc

**Key Topics**:
- Schema design and migrations
- SQL query writing with sqlc
- Cache integration
- Transaction patterns
- Testing database operations
- Best practices for queries

**When to use**: When creating new database tables, writing SQL queries, or working with `internal/repos/*`

**Quick Start**:
```bash
# Read the guidelines
cat internal/repos/REPOSITORY_GUIDE.md

# Generate repository code
make repos
```

## 🎯 Development Workflow

### For Brand New Projects

**Follow**: [PROJECT_SETUP_GUIDE.md](./PROJECT_SETUP_GUIDE.md)

This guide provides a complete initialization workflow with clear TODOs:
- Phase 1: Environment Setup
- Phase 2: Database Setup
- Phase 3: API Definition
- Phase 4: Service Implementation
- Phase 5: Build and Test
- Phase 6: Run the Service

### For New Features

1. **Define Data Model** → Follow `internal/repos/REPOSITORY_GUIDE.md`
   - Create schema in `internal/repos/{table}/schema.sql`
   - Write queries in `internal/repos/{table}/query.sql`
   - Run `make repos` to generate code

2. **Implement Service Logic** → Follow `SERVICE_TEST_GUIDE.md`
   - Add methods to `internal/svc/*.go`
   - Write tests in `internal/svc/*_test.go`
   - Run `make test`

### For Bug Fixes

1. Write a failing test first (follow `SERVICE_TEST_GUIDE.md`)
2. Fix the bug
3. Verify test passes
4. Check related queries (follow `internal/repos/REPOSITORY_GUIDE.md`)

### For Code Review

Check that the PR follows:
- ✅ Database changes use proper sqlc patterns
- ✅ Tests use golden files for complex responses
- ✅ Code passes `make lint`

## 🚀 Quick Reference Commands

```bash
# Generate all code
make proto          # Generate protobuf code
make repos          # Generate sqlc repository code

# Development
make build          # Build all binaries
make test           # Run all tests
make lint           # Run linters
make format         # Format code

# Running
make local-run      # Run server locally

# Docker
make docker-up      # Start all services
make docker-down    # Stop all services
```

## 📖 Reading Order for New Developers

1. **Start Here**: README.md (project overview)
2. **Initialize Project**: `PROJECT_SETUP_GUIDE.md` (step-by-step setup)
3. **Database Work**: `internal/repos/README.md` (quick start)
4. **Service Layer**: `SERVICE_TEST_GUIDE.md` (testing patterns)
5. **Advanced Topics**: 
   - `internal/repos/REPOSITORY_GUIDE.md` (detailed DB patterns)

## 🔍 Finding Information

| I want to... | Read this guide | Section |
|-------------|----------------|---------|
| **Initialize a new project** | `PROJECT_SETUP_GUIDE.md` | All phases |
| Add a new database table | `internal/repos/REPOSITORY_GUIDE.md` | Schema Design |
| Write a SQL query | `internal/repos/REPOSITORY_GUIDE.md` | Query Patterns |
| Add a service method | `SERVICE_TEST_GUIDE.md` | Service Structure |
| Write a test | `SERVICE_TEST_GUIDE.md` | Testing Patterns |
| Add caching | `internal/repos/REPOSITORY_GUIDE.md` | Cache Integration |
| Handle transactions | `internal/repos/REPOSITORY_GUIDE.md` | Transaction Patterns |

## 🎓 Best Practices Summary

### Database (repos)
- ✅ Use named parameters in SQL (`@param_name`)
- ✅ Add cache annotations for read queries
- ✅ Add timeout annotations for all queries
- ✅ Use transactions for multi-table operations
- ❌ Don't use string concatenation in SQL
- ❌ Don't forget to add invalidation annotations

### Service Layer (svc)
- ✅ Write tests with golden files
- ✅ Use table-driven tests
- ✅ Mock external dependencies
- ✅ Test error cases
- ❌ Don't test external services directly
- ❌ Don't commit with failing tests

## 🆘 Getting Help

1. **Search the guides** - Use Cmd+F in the relevant guide
2. **Check examples** - Look at existing code in the same area
3. **Run tests** - Often the tests show the expected patterns
4. **Ask the team** - If still stuck, reach out

## 📝 Contributing to Guidelines

Found an issue or want to improve the guidelines?

1. Update the relevant guide file
2. Test the changes by following your own instructions
3. Submit a PR with clear description
4. Tag for review

## 🔗 Related Resources

- **Project README**: `README.md`
- **API Documentation**: `api/proto/*/v1/*.proto`
- **Makefile**: All available commands
- **Service Demo**: Reference implementation examples

# 🔄 Git Workflow

## Branch Strategy

- **main**: Production-ready releases
- **develop**: Development branch (default)
- **feature/***: Feature branches
- **bugfix/***: Bug fix branches
- **hotfix/***: Hotfix branches

## Development Setup

```bash
# Clone the repository
git clone https://github.com/drjyhpaz/new_falcon.git
cd new_falcon

# Switch to develop branch
git checkout develop

# Create feature branch
git checkout -b feature/your-feature-name
```

## Commit Guidelines

- Use clear, descriptive commit messages
- Start with emoji for quick identification:
  - 🎨 Style changes
  - 🐛 Bug fixes
  - ✨ New features
  - 📚 Documentation
  - ⚡ Performance
  - 🔒 Security
  - 🧪 Tests
  - 🔨 Build/CI

## Merge Process

1. Push feature branch to remote
2. Create Pull Request
3. Wait for code review
4. Merge to develop
5. Test in develop environment
6. Create release PR to main

## Code Style

- Follow Go conventions
- Run `gofmt` before committing
- Use `golint` for linting
- Add comments for exported functions

```bash
# Format code
gofmt -w .

# Check for issues
golint ./...
```

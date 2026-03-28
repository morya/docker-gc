# GitHub Actions Specifications

## Node.js Version Requirements

GitHub Actions has deprecated Node.js 20 runtime. All actions should be updated to use Node.js 22 or later.

## Current Actions Configuration

### Workflow: Docker Build and Push
- File: `.github/workflows/docker-build.yml`
- Runs on: `ubuntu-latest`
- Uses deprecated Node.js 20 runtime

## Required Updates

### 1. Update Actions to Latest Versions
Ensure all GitHub Actions are using the latest versions that support Node.js 22+:

```yaml
# Current versions (may use Node.js 20)
- uses: actions/checkout@v4
- uses: docker/setup-buildx-action@v3
- uses: docker/login-action@v3
- uses: docker/metadata-action@v5
- uses: docker/build-push-action@v5

# Recommended updates (check for latest versions)
- uses: actions/checkout@v4  # Already latest
- uses: docker/setup-buildx-action@v3  # Already latest
- uses: docker/login-action@v3  # Already latest
- uses: docker/metadata-action@v5  # Already latest
- uses: docker/build-push-action@v5  # Already latest
```

### 2. Explicitly Set Node.js Version
Add Node.js version specification to the workflow:

```yaml
jobs:
  build-and-push:
    runs-on: ubuntu-latest
    # Add this to specify Node.js version
    env:
      NODE_VERSION: '22.x'
    
    steps:
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
      
      # ... rest of steps
```

### 3. Alternative: Update Ubuntu Runner
The `ubuntu-latest` runner may be using an older version. Consider specifying a newer runner version:

```yaml
jobs:
  build-and-push:
    # Use Ubuntu 24.04 or newer
    runs-on: ubuntu-24.04
```

## Implementation Priority

1. **High Priority**: Update the workflow to explicitly set Node.js 22+
2. **Medium Priority**: Update actions to their latest versions
3. **Low Priority**: Update runner to a newer Ubuntu version

## Testing

After making changes:
1. Push changes to a feature branch
2. Create a pull request to trigger the workflow
3. Verify no "Node.js 20 actions are deprecated" warnings appear
4. Ensure the Docker build and push still works correctly

## References

- [GitHub Actions: Node.js version support](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idruns-on)
- [GitHub Actions: Setting up Node.js](https://github.com/actions/setup-node)
- [GitHub Actions: Deprecation schedule](https://github.com/github/roadmap/issues/727)
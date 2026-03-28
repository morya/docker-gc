# GitHub Actions 规范

## Node.js 版本要求

GitHub Actions 已弃用 Node.js 20 运行时。所有 actions 都应更新为使用 Node.js 22 或更高版本。

## 当前 Actions 配置

### 工作流：Docker 构建和推送
- 文件：`.github/workflows/docker-build.yml`
- 运行环境：`ubuntu-latest`
- 使用已弃用的 Node.js 20 运行时

## 需要进行的更新

### 1. 更新 Actions 到最新版本
确保所有 GitHub Actions 都使用支持 Node.js 22+ 的最新版本：

```yaml
# 当前版本（可能使用 Node.js 20）
- uses: actions/checkout@v4
- uses: docker/setup-buildx-action@v3
- uses: docker/login-action@v3
- uses: docker/metadata-action@v5
- uses: docker/build-push-action@v5

# 推荐更新（检查最新版本）
- uses: actions/checkout@v4  # 已是最新
- uses: docker/setup-buildx-action@v3  # 已是最新
- uses: docker/login-action@v3  # 已是最新
- uses: docker/metadata-action@v5  # 已是最新
- uses: docker/build-push-action@v5  # 已是最新
```

### 2. 显式设置 Node.js 版本
在工作流中添加 Node.js 版本规范：

```yaml
jobs:
  build-and-push:
    runs-on: ubuntu-latest
    # 添加此设置以指定 Node.js 版本
    env:
      NODE_VERSION: '22.x'
    
    steps:
      - name: 设置 Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
      
      # ... 其余步骤
```

### 3. 替代方案：更新 Ubuntu 运行器
`ubuntu-latest` 运行器可能使用较旧版本。考虑指定更新的运行器版本：

```yaml
jobs:
  build-and-push:
    # 使用 Ubuntu 24.04 或更新版本
    runs-on: ubuntu-24.04
```

## 实施优先级

1. **高优先级**：更新工作流以显式设置 Node.js 22+
2. **中优先级**：更新 actions 到最新版本
3. **低优先级**：更新运行器到更新的 Ubuntu 版本

## 测试

进行更改后：
1. 将更改推送到功能分支
2. 创建拉取请求以触发工作流
3. 验证没有出现 "Node.js 20 actions are deprecated" 警告
4. 确保 Docker 构建和推送仍然正常工作

## 参考链接

- [GitHub Actions: Node.js 版本支持](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idruns-on)
- [GitHub Actions: 设置 Node.js](https://github.com/actions/setup-node)
- [GitHub Actions: 弃用时间表](https://github.com/github/roadmap/issues/727)
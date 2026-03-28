# AGENTS.md - Docker GC 项目指南

本文档为在 Docker GC 项目上工作的 AI 代理提供指南，这是一个使用 Golang 实现的带有 cron 调度的 Docker 垃圾回收工具。

## 项目概述

- **语言**: Go 1.21+
- **目的**: 带有 cron 调度的 Docker 垃圾回收
- **架构**: 单二进制应用程序，集成 Docker CLI
- **部署**: 使用多阶段构建的 Docker 容器

## 构建命令

### 开发构建
```bash
# 构建 Go 二进制文件
go build -o docker-gc-cron ./cmd/docker-gc-cron

# 本地运行
./docker-gc-cron
```

### Docker 构建
```bash
# 构建 Docker 镜像
docker build -t ghcr.io/morya/docker-gc:latest .

# 使用 Docker socket 运行
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/morya/docker-gc:latest
```

### 依赖管理
```bash
# 下载依赖
go mod download

# 整理依赖
go mod tidy

# 验证依赖
go mod verify
```

## 测试命令

**注意**: 当前项目没有测试套件。添加测试时：

```bash
# 运行所有测试
go test ./...

# 运行带覆盖率的测试
go test -cover ./...

# 运行特定测试
go test -run TestFunctionName ./path/to/package

# 运行详细输出的测试
go test -v ./...
```

## 代码检查和代码质量

### Go 格式化
```bash
# 格式化所有 Go 文件
go fmt ./...

# 检查格式化而不应用
gofmt -d .
```

### Go Vet（静态分析）
```bash
# 运行 go vet 进行静态分析
go vet ./...
```

### 建议的代码检查工具（当前未配置）
- `golangci-lint run` - 全面的代码检查器
- `staticcheck ./...` - 高级静态分析
- `revive ./...` - 快速、可配置的代码检查器

## 代码风格指南

### 导入组织
```go
import (
    // 标准库
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"

    // 第三方包
    "github.com/robfig/cron/v3"
)
```

### 命名约定
- **包**: 小写，单个单词（例如 `main`）
- **变量**: camelCase（例如 `containerIDs`, `gracePeriodSeconds`）
- **常量**: camelCase 或 UPPER_SNAKE_CASE 用于环境默认值
- **类型**: PascalCase（例如 `Config`, `ContainerInfo`）
- **接口**: PascalCase，如果合适则以 "er" 结尾（例如 `Cleaner`）

### 错误处理
- 使用 `log.Printf` 处理非致命错误并提供上下文
- 仅对不可恢复的启动错误使用 `log.Fatalf`
- 出错时尽早从函数返回
- 包含有意义的错误消息和上下文

示例：
```go
output, err := cmd.Output()
if err != nil {
    log.Printf("Failed to list exited containers: %v", err)
    return
}
```

### 函数结构
- 保持函数专注且短小（< 50 行）
- 使用描述性函数名（例如 `cleanContainers`, `loadConfig`）
- 使用注释记录公共函数
- 将相关函数分组在一起

### 配置模式
- 使用 `Config` 结构体处理环境变量
- 提供合理的默认值
- 在 `loadConfig()` 函数中验证配置
- 使用布尔标志进行功能切换（例如 `FORCE_CONTAINER_REMOVAL=1`）

### 日志记录指南
- 使用 `log.Println` 记录信息性消息
- 重要操作前加上 `[Docker GC]` 前缀
- 包含相关标识符（容器 ID、镜像 ID）
- 在重要操作前后记录日志

### Docker CLI 集成
- 使用 `exec.Command` 进行 Docker 操作
- 始终检查命令输出和错误
- 优雅处理空结果
- 支持干运行模式进行测试

## 项目结构

```
docker-gc/
├── cmd/
│   └── docker-gc-cron/
│       └── main.go          # 主应用程序入口点
├── Dockerfile               # 多阶段 Docker 构建
├── entrypoint.sh           # 容器入口点脚本
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖校验和
├── .github/
│   └── workflows/
│       └── docker-build.yml # CI/CD 流水线
└── tasks/
    └── task.md             # 项目需求
```

## Docker 构建指南

### 多阶段构建
- 构建阶段使用 `golang:alpine`
- 运行时阶段使用 `alpine:latest`
- 设置 `CGO_ENABLED=0` 用于静态二进制文件
- 仅在阶段之间复制必要的文件

### 运行时依赖
- 运行时镜像中包含 `docker-cli`
- 安装 `bash` 和 `tzdata` 以支持 cron
- 创建必要的目录（`/var/log`）
- 设置适当的文件权限

### 环境变量
- 在 README 中记录所有环境变量
- 提供合理的默认值
- 使用布尔标志（`0`/`1`）进行切换
- 支持通过 Docker run 命令进行配置

## GitHub Actions 工作流

项目使用 GitHub Actions 进行 CI/CD：
- 推送到 main/master 分支时构建
- 推送标签（`v*`）时构建
- 拉取请求时构建（不推送）
- 推送到 GitHub Container Registry

**重要**: 工作流使用 Node.js 22.x，符合 GitHub Actions 规范要求。

## 开发工作流

1. **进行更改**: 编辑 `cmd/docker-gc-cron/main.go` 或其他文件
2. **本地测试**: `go build -o docker-gc-cron ./cmd/docker-gc-cron && ./docker-gc-cron`
3. **构建 Docker**: `docker build -t ghcr.io/morya/docker-gc:latest .`
4. **测试容器**: `docker run -v /var/run/docker.sock:/var/run/docker.sock ghcr.io/morya/docker-gc:latest`
5. **提交更改**: 遵循约定式提交消息
6. **创建 PR**: 更改将自动构建和测试

## 添加新功能

添加新功能时：
1. 在 `main.go` 的 `Config` 结构体中添加配置
2. 在 `loadConfig()` 中添加环境变量处理
3. 在 `README.md` 中更新新环境变量
4. 如果需要，在 `Dockerfile` 中更新默认值
5. 考虑向后兼容性
6. 使用各种环境配置进行测试

## AI 代理的常见任务

### 添加新的清理功能
1. 创建新函数（例如 `cleanNetworks()`）
2. 添加到 `runGarbageCollection()` 调用序列中
3. 如果需要，添加配置选项
4. 更新文档

### 改进错误处理
1. 添加更具体的错误消息
2. 考虑为暂时性故障添加重试逻辑
3. 添加错误率的指标/日志记录
4. 测试错误场景

### 性能优化
1. 考虑并行执行清理任务
2. 尽可能批量处理 Docker 操作
3. 添加资源使用监控
4. 分析内存和 CPU 使用情况

## AI 代理注意事项

- 这是一个相对简单的 Go 应用程序，依赖项最少
- 代码遵循标准的 Go 惯用模式和模式
- 专注于清晰性和可维护性，而不是过早优化
- 如有疑问，请遵循现有的代码模式
- 进行更改后始终测试 Docker 构建
- 考虑容器化部署环境
# Docker GC（Golang 实现）

一个使用 Golang 实现的 Docker 镜像，允许定时清理未使用的 Docker 镜像、容器和卷。

## 功能特性

- 使用 cron 进行定时垃圾回收
- 清理超过宽限期的已退出容器
- 清理未使用（悬空）的镜像
- 可选清理悬空卷
- 通过环境变量配置
- 用于测试的干运行模式

## 环境变量

| 变量 | 默认值 | 描述 |
|------|--------|------|
| `CRON` | `0 0 * * *`（每天午夜） | 垃圾回收的 cron 时间表 |
| `GRACE_PERIOD_SECONDS` | `3600`（1小时） | 容器/镜像被移除的最小年龄 |
| `FORCE_CONTAINER_REMOVAL` | `0` | 强制移除容器（设置为 `1` 启用） |
| `FORCE_IMAGE_REMOVAL` | `0` | 强制移除镜像（设置为 `1` 启用） |
| `CLEAN_UP_VOLUMES` | `0` | 清理悬空卷（设置为 `1` 启用） |
| `DRY_RUN` | `0` | 干运行模式 - 不实际移除（设置为 `1` 启用） |
| `EXCLUDE_VOLUMES_IDS_FILE` | （无） | 包含要排除的卷 ID 的文件路径 |
| `VOLUME_DELETE_ONLY_DRIVER` | （无） | 仅删除具有特定驱动程序的卷 |

## 快速开始

```bash
# 使用默认设置运行（每天午夜）
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/morya/docker-gc:latest

# 使用自定义时间表运行（每6小时）
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e CRON="0 */6 * * *" \
  ghcr.io/morya/docker-gc:latest

# 使用强制移除和卷清理运行
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e FORCE_CONTAINER_REMOVAL=1 \
  -e FORCE_IMAGE_REMOVAL=1 \
  -e CLEAN_UP_VOLUMES=1 \
  ghcr.io/morya/docker-gc:latest

# 干运行模式（测试而不实际移除）
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e DRY_RUN=1 \
  ghcr.io/morya/docker-gc:latest
```

## 从源代码构建

```bash
# 克隆仓库
git clone https://github.com/morya/docker-gc.git
cd docker-gc

# 构建 Docker 镜像
docker build -t ghcr.io/morya/docker-gc:latest .

# 本地运行
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/morya/docker-gc:latest
```

## 开发

```bash
# 构建 Go 二进制文件
go build -o docker-gc-cron ./cmd/docker-gc-cron

# 本地运行
./docker-gc-cron
```

## Dockerfile 特性

- 两阶段构建：golang:alpine 用于构建，alpine 用于运行时
- CGO_ENABLED=0 用于静态二进制文件
- 最终镜像体积小
- 包含 docker-cli 用于 Docker 操作

## GitHub Actions

通过 GitHub Actions 配置了自动化构建。在以下情况下构建镜像并推送到 GitHub Container Registry：
- 推送到 main/master 分支
- 推送标签（v*）
- 拉取请求（仅构建，不推送）

## 许可证

Apache License 2.0
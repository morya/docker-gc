# docker-gc 项目

帮我用 golang 重新实现一遍  `clockworksoul/docker-gc-cron` 功能。

## 实现要求

- 提供 Dockerfile 构建 当前项目
    - golang 的构建基础镜像，使用 golang:alpine 并且，构建中，实现为 2 步 构建，构建出来的 /app 执行时，基础镜像使用 alpine 。
    - 构建golang app时，为方便构建能够执行，需要提前设置 CGO_ENABLED=0 环境变量
- 为当前项目 编写 github actions 配置文件
- 不需要实现 tests 相关测试单元


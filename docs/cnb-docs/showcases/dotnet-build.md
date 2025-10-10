# CNB 配置 .NET 项目构建并推送 docker 镜像到 CNB 制品库

本文档地址：https://cnb.cool/examples/showcase

文档摘要：本文档指导构建.NET项目的Docker镜像并通过CNB制品库进行推送。关键步骤包括：①创建Dockerfile实现多阶段构建，基于SDK镜像编译后在运行时镜像部署，通过参数化处理多平台构建；②配置.cnb.yml定义流水线任务，包含Docker注册表登录、镜像构建（使用制品库命名规范）和推送操作，并集成环境变量实现安全认证。技术点涵盖镜像优化策略、构建参数传递及制品库发布流程。

.NET 项目如何构建镜像并发布到 CNB Docker 仓库。

## 1. 准备 Dockerfile 文件

在项目根目录下创建一个 `Dockerfile` 文件，用于构建 docker 镜像。

```dockerfile
# Learn about building .NET container images:
# https://github.com/dotnet/dotnet-docker/blob/main/samples/README.md
FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/sdk:9.0 AS build
ARG TARGETARCH
WORKDIR /source

# Copy project file and restore as distinct layers
COPY --link *.csproj .
RUN dotnet restore -a $TARGETARCH

# Copy source code and publish app
COPY --link . .
RUN dotnet publish -a $TARGETARCH --no-restore -o /app


# Runtime stage
FROM mcr.microsoft.com/dotnet/runtime:9.0
WORKDIR /app
COPY --link --from=build /app .
USER $APP_UID
ENTRYPOINT ["./dotnetapp"]
```

## 2. 配置 .cnb.yml 流水线配置文件

在项目根目录下创建一个 `.cnb.yml` 文件，用于配置构建任务：安装依赖 -> 编译 -> 上传到目标服务器。

```yaml
main:
  push:
    # 上传 dotnet-build 的 dotnetapp 应用示例镜像到 CNB 制品库
    - services:
        - docker
      stages:
        - name: docker login
          script: docker login -u ${CNB_TOKEN_USER_NAME} -p "${CNB_TOKEN}" ${CNB_DOCKER_REGISTRY}
        - name: docker build
          script: docker build -t ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}/dotnetapp:latest .
        - name: docker push
          script: docker push ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}/dotnetapp:latest
```

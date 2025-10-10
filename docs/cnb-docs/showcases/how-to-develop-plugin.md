# 如何开发一个插件示例

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该示例指导开发了一个CNB插件，支持构建多架构Docker镜像。通过创建包含`hello world`功能的Dockerfile和配置`.cnb.yml`文件，实现AMD64与ARM64架构的并行构建。文件定义了分平台构建流水线任务，每个任务包含Docker登录、镜像构建推送步骤，并通过`cnb:await`组件同步完成后，使用`manifest`插件合并生成多架构镜像。关键配置涉及环境变量设置、平台任务标记、目标镜像标签模板及多平台兼容性声明。

这个示例用于展示，如何同时构建 `amd64` 和 `arm64` 平台的 `hello world` 镜像，并在用于 CNB 插件。

## 1、准备一个 Dockerfile

准备一个 Dockerfile，用于实现具体的插件逻辑。例如下面这个 Dockerfile 实现了一个简单的 `hello world` 插件：

```docker
FROM alpine

ENTRYPOINT ["echo", "Hello, World!"]
```

## 2、配置 .cnb.yml 配置文件

在项目根目录下创建一个 `.cnb.yml` 文件，用于配置插件的构建任务。
下面的配置文件定义了一个多架构构建流水线，包含了 AMD64 和 ARM64 两个架构的镜像构建任务，以及最后将它们合并为一个多架构镜像的任务。

```yaml
main:
  push:
    # AMD64架构的构建任务
    - runner:
        tags: cnb:arch:amd64  # 指定运行在AMD64架构的runner上
      services:
        - docker  # 使用Docker服务
      env:
        IMAGE_TAG: ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest-linux-amd64  # 设置Docker镜像标签
      stages:
        - name: docker login
          script: docker login -u ${CNB_TOKEN_USER_NAME} -p "${CNB_TOKEN}" ${CNB_DOCKER_REGISTRY}  # 登录Docker registry
        - name: docker build
          script: docker build -t ${IMAGE_TAG} .  # 构建Docker镜像
        - name: docker push
          script: docker push ${IMAGE_TAG}  # 推送Docker镜像到registry
        - name: resolve
          type: cnb:resolve
          options:
            key: build-amd64  # 标记AMD64构建完成
    
    # ARM64架构的构建任务
    - runner:
        tags: cnb:arch:arm64:v8  # 指定运行在ARM64架构的runner上
      services:
        - docker
      env:
        IMAGE_TAG: ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest-linux-arm64
      stages:
        - name: docker login
          script: docker login -u ${CNB_TOKEN_USER_NAME} -p "${CNB_TOKEN}" ${CNB_DOCKER_REGISTRY}
        - name: docker build
          script: docker build -t ${IMAGE_TAG} .
        - name: docker push
          script: docker push ${IMAGE_TAG}
        - name: resolve
          type: cnb:resolve
          options:
            key: build-arm64  # 标记ARM64构建完成

    # 创建多架构镜像
    - services:
        - docker
      env:
        IMAGE_TAG: ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest
      stages:
        - name: await the amd64
          type: cnb:await
          options:
            key: build-amd64  # 等待AMD64构建完成
        - name: await the arm64
          type: cnb:await
          options:
            key: build-arm64  # 等待ARM64构建完成
        - name: manifest
          image: plugins/manifest  # 使用manifest插件创建多架构镜像
          settings:
            username: $CNB_TOKEN_USER_NAME
            password: $CNB_TOKEN
            target: ${IMAGE_TAG}  # 目标镜像标签
            template: ${IMAGE_TAG}-OS-ARCH  # 模板格式
            platforms:  # 支持的平台
              - linux/amd64
              - linux/arm64
```

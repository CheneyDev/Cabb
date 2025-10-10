# CNB 配置 SpringBoot + Gradle 项目， 并且构建 docker 镜像

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该文档指导完成SpringBoot+Gradle项目的云原生构建流程：通过.cnb.yml定义使用gradle:6.8-jdk8构建环境，挂载缓存目录并启用Docker服务，执行构建、登录制品库、构建Docker镜像及推送的操作。指定基于OpenJDK 8的Dockerfile完成JAR部署与环境配置，同时提供云原生开发Dockerfile以安装Code Server、扩展插件及SSH服务，实现开发环境集成。所有操作均通过流水线完成制品化交付。

通过云原生构建实现，打包 springboot+gradle 项目, 构建 Docker 镜像并发布到制品库

gradle:6.8-jdk8 构建环境 -> 挂载 gradle 缓存目录 -> 启用 docker 服务 -> 打包项目 -> 登录制品库 -> 构建镜像 -> 推送镜像

## 1. 配置 .cnb.yml 文件
```yaml
main:
  push:
    - docker:
        # 声明式的构建环境 https://docs.cnb.cool/
        # 可以去dockerhub https://hub.docker.com/_/gradle 找到您需要的gradle和jdk版本
        image: gradle:6.8-jdk8
        volumes:
          # 声明式的构建缓存 https://docs.cnb.cool/zh/grammar/pipeline.html#volumes
          - /root/.gradle:copy-on-write
      services:
        # 流水线中启用 docker 服务
        - docker
      stages:
        - name: gradle build
          script:
            - ./gradlew build
        # 云原生构建自动构建Docker镜像并将它发布到制品库参考【上传Docker制品】https://docs.cnb.cool/zh/artifact/docker.html
        - name: docker login
          script:
            - docker login -u ${CNB_TOKEN_USER_NAME} -p "${CNB_TOKEN}" ${CNB_DOCKER_REGISTRY}
        - name: docker build
          script:
            - docker build -t ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:${CNB_COMMIT} .
        - name: docker push
          script:
            - docker push ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:${CNB_COMMIT}

```

## 2. 配置 Dockerfile 文件

```dockerfile
# 使用官方的OpenJDK镜像作为基础镜像
FROM openjdk:8

# 设置工作目录
WORKDIR /app

# 将JAR包复制到镜像中
COPY ./build/libs/gradle-deploy-0.1-SNAPSHOT.jar /app/gradle-deploy-0.1-SNAPSHOT.jar

# 暴露应用程序的端口（如果需要）
EXPOSE 8081

# 运行JAR包
CMD ["java", "-jar", "gradle-deploy-0.1-SNAPSHOT.jar"]
```


## 3. 配置 springboot + gradle 项目云原生开发

在 `.ide/Dockerfile` 文件中配置云原生开发环境

```Dockerfile
# 帮助文档地址: https://docs.cnb.cool/zh/vscode/quick-start.html
# 可以去 dockerhub: https://hub.docker.com/_/gradle 找到您需要的 gradle 和 jdk 版本
FROM gradle:6.8-jdk8

# 腾讯云软件源使用示例: https://cnb.cool/examples/mirrors/mirrors.cloud.tencent.com
RUN sed -Ei "s@(security|ports|archive).ubuntu.com@mirrors.cloud.tencent.com@g" /etc/apt/sources.list

# 以及按需安装其他软件
# RUN apt-get update && apt-get install -y git

# 安装 code-server 和 vscode 常用插件
RUN curl -fsSL https://code-server.dev/install.sh | sh \
  && code-server --install-extension redhat.vscode-yaml \
  && code-server --install-extension eamodio.gitlens \
  && code-server --install-extension tencent-cloud.coding-copilot \
  && code-server --install-extension vscjava.vscode-java-pack \
  && echo done

# 安装 ssh 服务，用于支持 VSCode 客户端通过 Remote-SSH 访问开发环境
RUN apt-get update && apt-get install -y wget unzip openssh-server

# 指定字符集支持命令行输入中文（根据需要选择字符集）
ENV LANG C.UTF-8
ENV LANGUAGE C.UTF-8
```
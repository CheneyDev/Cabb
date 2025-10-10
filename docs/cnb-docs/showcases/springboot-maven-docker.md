# CNB 配置 SpringBoot + Maven 项目，并且构建 docker 镜像

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该配置实现SpringBoot+Maven项目的云原生构建流程：通过定义.cnb.yml文件指定Maven 3.8.6+OpenJDK 8构建环境，使用挂载的Maven缓存加速构建，自动执行"mvn clean package"构建JAR包，经制品库认证后构建Docker镜像（基于openjdk:8基础镜像，复制JAR至/app目录），最后将包含版本号的镜像推送到制品库。Dockerfile明确设置了运行时端口和JAR启动命令。

将通过云原生构建实现，打包 springboot+maven 项目，构建 Docker 镜像并发布到制品库。

## 1. 配置 .cnb.yml 文件

在项目根目录下创建一个 `.cnb.yml` 文件，用于配置构建 Maven 项目任务。步骤如下：

maven:3.8.6-openjdk-8 构建环境 -> 挂载 Maven 缓存目录 -> 启用 docker 服务 -> 打包项目 -> 登录制品库 -> 构建镜像 -> 推送镜像


```yaml
main:
  push:
    - docker:
        # 声明式的构建环境 https://docs.cnb.cool/
        # 可以去dockerhub上 https://hub.docker.com/_/maven 找到您需要maven和jdk版本
        image: maven:3.8.6-openjdk-8
        volumes:
          # 声明式的构建缓存 https://docs.cnb.cool/zh/grammar/pipeline.html#volumes
          - /root/.m2:copy-on-write
      services:
        # 流水线中启用 docker 服务
        - docker
      stages:
        - name: mvn package
          script:
            # 合并./settings.xml和/root/.m2/settings.xml
            - mvn clean package -s ./settings.xml
        # 云原生构建自动构建Docker镜像并将它发布到制品库，【上传Docker制品】https://docs.cnb.cool/zh/artifact/docker.html
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
COPY ./target/maven-deploy-0.1-SNAPSHOT.jar /app/maven-deploy-0.1-SNAPSHOT.jar

# 暴露应用程序的端口（如果需要）
EXPOSE 8081

# 运行JAR包
CMD ["java", "-jar", "maven-deploy-0.1-SNAPSHOT.jar"]
```
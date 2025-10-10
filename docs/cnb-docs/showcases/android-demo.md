# CNB 配置打包构建 Android 安卓应用

本文档地址：https://cnb.cool/examples/showcase

文档摘要：CNB构建配置通过项目根目录的.cnb.yml文件实现Android应用构建任务，核心配置包括指定基于34.0.1版本的移动开发安卓SDK镜像进行构建，并挂载根目录的.gradle缓存目录以优化构建效率。任务执行包含两个阶段：执行gradlew build构建APK，然后通过ls命令列出应用build/outputs/apk/release目录下的输出文件。支持使用Hub.docker.com的现有镜像或自定义SDK版本Docker镜像满足不同环境需求。

## 1. 配置 .cnb.yml 文件

在项目根目录下创建一个 `.cnb.yml` 文件，用于配置构建 Android 安卓应用任务。 

```yaml
master:
  push:
    - docker:
        # 可以在 hub.docker.com 上找需要的 android sdk 版本的 docker 镜像
        # https://hub.docker.com/r/mobiledevops/android-sdk-image
        # https://github.com/docker-android-sdk
        # 当这些都不满足您的需求时，您可以制作自己的 docker 镜像安装您需要的 sdk 版本和工具
        image: mobiledevops/android-sdk-image:34.0.1
        # 挂载 gradle 缓存目录，加快构建速度
        volumes:
          - /root/.gradle:cow
      stages:
        - name: android-build
          script: ./gradlew build
        - name: "ls"
          script: ls ./app/build/outputs/apk/release
```
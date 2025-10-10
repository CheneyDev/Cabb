---
title: 构建环境
permalink: https://docs.cnb.cool/zh/build/build-env.html
summary: 构建环境由 Docker 容器提供，支持两种配置方式：指定已有镜像或使用 Dockerfile 动态构建。未指定时，默认使用 cnbcool/default-build-env:latest。镜像中的 VOLUME 可通过参数共享给插件任务，以 NodeJS 构建为例，可使用官方 NodeJS 镜像进行依赖安装和测试。
redirectFrom: /zh/build-env.html
---
构建环境确定了在运行构建任务的时候，环境中拥有哪些软件；

`云原生构建` 使用 Docker 容器作为构建环境。

相比于传统的虚拟机容器，Docker 容器有非常巨大的优势，
行业中各个 CI 厂商也在向这方面发展。

请确保使用 `云原生构建` 时，拥有使用 Docker 容器的相关经验，
或至少了解 Docker 入门知识。

## 配置方式

Docker 容器的配置方式有两种：

1. **指定一个 image（已有的镜像）**。这些镜像是之前由其他人制作出来并且推送到镜像仓库中的。
2. **指定一个 Dockerfile**。在构建开始时，`云原生构建` 会根据 Dockerfile 来即时制作镜像（或某种条件下使用缓存）来使用。

### 指定 Image

```yaml
main:
  push:
    - docker:
        # 通过此参数控制使用的 image
        image: node:22
      stages:
        - stage1
        - stage2
        - stage3
```

用 [pipeline.docker.image](./grammar.md#image) 参数指定 Docker Image。
其值可以是 公网中的官方镜像源或其他可以被访问到的镜像源 中的公开镜像。

### 指定 Dockerfile

```yaml
main:
  push:
    - docker:
        # 通过此参数控制使用的 Dockerfile
        build: ./image/Dockerfile
      stages:
        - stage1
        - stage2
        - stage3
```

用 [pipeline.docker.build](./grammar.md#build) 参数指定 Dockerfile。

### 缺省镜像

当未指定 image 和 Dockerfile 时，image 会被设置为缺省镜像：[cnbcool/default-build-env:latest](cnbcool/default-build-env:latest_LINK)

```yaml
main:
  push:
    - stages:
        - stage1
        - stage2
        - stage3
      # 未指定 image 和 Dockerfile 相当于如下声明
      # docker:
      # image: cnbcool/default-build-env:latest
```

### 镜像中的 VOLUME

镜像中可能包含`VOLUME`命令。
用**插件任务**的容器启动时会通过 `--volumes-from` 参数，将这些 volume 共享给 **插件任务**。

例如，在 Dockerfile 中准备了文件：

```dockerfile
RUN mkdir /cache && echo 'hello world' > /cache/data.txt

VOLUME /cache
```

在后续的 image-commands 中也可以访问到它。

```yaml
- name: image-commands中访问pipeline volume
  image: alpine
  commands:
    - cat /cache/data.txt
```

## 使用案例

### NodeJS

当需要使用 NodeJS 的构建环境时，可以直接使用 Docker Hub 之中的官方 NodeJS 镜像。

```yaml
main:
  push:
    - docker:
        image: node:20
      stages:
        - name: 依赖安装
          script: npm install
        - name: 测试用例检查
          script: npm test
```

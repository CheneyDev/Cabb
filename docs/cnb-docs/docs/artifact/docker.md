---
title: Docker 制品库
permalink: https://docs.cnb.cool/zh/artifact/docker.html
summary: 该文档介绍了Docker制品库的使用方法，包括使用CNB访问令牌登录、制品路径规则（同名制品和非同名制品）以及推送和拉取制品的命令。此外，还介绍了在云原生构建和开发环境中如何使用Docker制品，以及在命令行中使用制品。更多详细用法可查阅Docker官方文档。
---
## 登录 CNB Docker 制品库

您可以使用 CNB 的访问令牌作为登录凭据，登录命令：

```bash
docker login docker.cnb.cool -u cnb -p {token-value}
```

## Docker 制品路径规则

制品在发布到某一仓库时，支持两种命名规则

1. 同名制品 - 制品路径与仓库路径一致，如：`docker.cnb.cool/{repository-path}`
2. 非同名制品 - 仓库路径作为制品的命名空间，制品路径=仓库路径/制品名称，如：`docker.cnb.cool/{repository-path}/{artifact-name}`

## 推送制品

### 本地命令行推送

同名制品

```bash
docker build -t docker.cnb.cool/{repository-path}:latest .
docker push docker.cnb.cool/{repository-path}:latest
```

非同名制品

```bash
docker build -t docker.cnb.cool/{repository-path}/{image-name}:latest .
docker push docker.cnb.cool/{repository-path}/{image-name}:latest
```

### 云原生构建中推送

```yaml
main:
  push:
    - services:
        - docker
      stages:
        - name: docker build
          script: docker build -t ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest .
        - name: docker push
          script: docker push ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest
```

### 云原生开发中推送

同名制品

```bash
docker build -t docker.cnb.cool/{repository-path}:latest .
docker push docker.cnb.cool/{repository-path}:latest
```

非同名制品

```bash
docker build -t docker.cnb.cool/{repository-path}/{image-name}:latest .
docker push docker.cnb.cool/{repository-path}/{image-name}:latest
```

## 使用制品

### 在命令行使用

```bash
docker pull docker.cnb.cool/{artifact-path}:latest

# ...
```

### 定制云原生构建环境

```yaml{4}
main:
  push:
    - docker:
        image: docker.cnb.cool/{artifact-path}:latest
      stages:
        - name: hello world
          script: echo "Hello World"
```

### 定制云原生开发环境

```yaml{4}
$:
  vscode:
    - docker:
        image: docker.cnb.cool/{artifact-path}:latest
      services:
        - vscode
        - docker
```

## 更多用法

更多 Docker 用法，请查阅官方文档

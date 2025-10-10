# CNB 配置 Python Flask 项目，并且构建 docker 镜像

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该文档指导基于CNB构建Python-Flask项目并生成Docker镜像的具体流程。通过云原生架构，使用Python 3.11-alpine作为构建环境，经Docker服务整合，执行制品库登录、镜像构建与推送操作，最终将构建产物上传至Docker仓库。需配置.cnb.yml定义构建阶段包含镜像登录(Docker login)、镜像构建(docker build)和镜像推送(docker push)三个脚本步骤，同时明确应用依赖通过腾讯云镜像源安装。Dockerfile采用Python 3.8基础镜像，设置应用工作目录并运行Flask入口文件。

将通过云原生构建 CNB 实现，打包 python-flask-docker 项目，构建并将构建产物上传 Docker 制品库中

## 1. 配置 .cnb.yml 文件

python:3.11-alpine 构建环境 -> 启用 docker 服务 -> 打包项目 -> 登录制品库 -> 构建镜像 -> 推送镜像

```yaml
# 帮助文档地址: https://docs.cnb.cool/zh/artifact/docker.html
main:
  push:
    - services:
        - docker
      stages:
        - name: docker login
          script: docker login -u ${CNB_TOKEN_USER_NAME} -p "${CNB_TOKEN}" ${CNB_DOCKER_REGISTRY}
        - name: docker build
          script: docker build -t ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:${CNB_COMMIT} .
        - name: docker push
          script: docker push ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:${CNB_COMMIT}

```


## 2. 配置 Dockerfile 文件

```Dockerfile
FROM python:3.8

WORKDIR /app

COPY . /app

RUN pip install -r requirements.txt -i https://mirrors.cloud.tencent.com/pypi/simple

# 假设 app.py 是 Flask 应用的入口文件
CMD python3 app.py
```
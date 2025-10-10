# 在远程开发环境中启动 MySQL 和 Redis 等数据库中间件

本文档地址：https://cnb.cool/examples/showcase

文档摘要：在远程开发环境中，通过配置.cnb.yml文件定义基于Docker Compose的任务来启动MySQL和Redis数据库。该文件触发构建流程时执行docker compose up -d命令，而独立的docker-compose.yml则定义了服务配置：MySQL使用5.7镜像，设置根密码和初始化脚本，映射3306端口；Redis采用alpine镜像并对外开放6379端口。开发环境通过VSCode的Remote-SSH协议连接，任务执行步骤覆盖构建容器化环境到启动中间件的全过程。

基本思路：
1. 配置 .cnb.yml 文件，用于配置远程开发启动 MySQL 和 Redis 等数据库中间件的构建任务。
2. 配置 Docker Compose 文件，用于配置 MySQL 和 Redis 等数据库中间件。
3. 配置 vscode 远程开发环境，用于支持 vscode 客户端通过 Remote-SSH 访问开发环境。参考[《云原生开发快速入门》](https://docs.cnb.cool/zh/vscode/quick-start.html)。

仓库地址：[https://cnb.cool/examples/ecosystem/nodejs-docker-compose](https://cnb.cool/examples/ecosystem/nodejs-docker-compose)

## 1、配置 .cnb.yml 文件

在项目根目录下创建一个 `.cnb.yml` 文件，用于配置远程开发启动 MySQL 和 Redis 等数据库中间件的构建任务。

```yaml
$:
  # vscode 事件：专供页面中启动远程开发用
  vscode:
    - docker:
        build: .ide/Dockerfile
      services:
        - vscode
        - docker
      stages:
        # 启动 MySQL 和 Redis 等数据库中间件
        - name: start docker compose
          script: docker compose up -d
```

## 2、配置 Docker Compose 文件

在项目根目录下创建一个 `docker-compose.yml` 文件，用于配置 MySQL 和 Redis 等数据库中间件。

```yaml
services:
  # MySQL 数据库
  mysql:
    image: mysql:5.7
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: demo
    # 挂载 MySQL 初始化脚本，用于在容器启动时初始化数据库
    volumes:
      # 安装实际的 MySQL 初始化脚本路径填写
      - ./mysql-init:/docker-entrypoint-initdb.d
    ports:
      - "3306:3306"

  # Redis 数据库
  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
```


## 参考链接

- [在远程开发环境中启动 MySQL 和 Redis 等数据库中间件](https://cnb.cool/examples/ecosystem/nodejs-docker-compose)
- [Docker Compose 官方文档](https://docs.docker.com/compose/)
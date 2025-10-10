# 腾讯云软件源使用示例

本文档地址：https://cnb.cool/examples/showcase

文档摘要：文本展示了利用腾讯云软件源加速不同操作系统及开发工具包下载的具体方法。针对Alpine（修改apk仓库）、Debian/Ubuntu（替换apt/deb源）、PHP（调整apt源并配置Composer镜像）系统通过Dockerfile实现镜像加速；Python使用pip3切换国内PyPI镜像；npm/yarn通过命令行设置注册表指向腾讯云节点。所有配置均采用替换官方源地址为核心策略，并通过清理缓存目录优化镜像体积，适配不同包管理器的差异化操作路径。

使用腾讯云软件源，可以加速软件包的下载。

仓库地址：[腾讯云软件源使用示例](https://cnb.cool/examples/mirrors/mirrors.cloud.tencent.com)

## Alpine 使用腾讯云软件源加速
配置 `Dockerfile` 文件，使用腾讯云软件源加速 Alpine 镜像的下载。

```dockerfile
FROM alpine:3.20.3

RUN sed -i 's@dl-cdn.alpinelinux.org@mirrors.cloud.tencent.com@g' /etc/apk/repositories && \
    apk --no-cache add git
```

## Debian 使用腾讯云软件源加速
配置 `Dockerfile` 文件，使用腾讯云软件源加速 Debian 镜像的下载。

```dockerfile
FROM debian:12.7

RUN sed -i "s@http://deb.debian.org/debian@http://mirrors.cloud.tencent.com/debian@g" /etc/apt/sources.list.d/debian.sources

RUN apt-get update && \
    apt-get install -y git && \
    rm -rf /var/lib/apt/lists/*
```

## Ubuntu 使用腾讯云软件源加速
配置 `Dockerfile` 文件，使用腾讯云软件源加速 Ubuntu 软件包的下载。

```dockerfile
FROM ubuntu:24.04

RUN sed -Ei "s@(security|ports|archive).ubuntu.com@mirrors.cloud.tencent.com@g" /etc/apt/sources.list.d/ubuntu.sources && \
    apt-get update && \
    apt-get install -y git && \
    rm -rf /var/lib/apt/lists/*
```

## PHP 使用腾讯云软件源加速
配置 `Dockerfile` 文件，使用腾讯云软件源加速 PHP 软件包的下载。

```dockerfile
FROM composer:2.2.24 as composer

FROM php:8.1-fpm-bookworm

COPY --from=composer /usr/bin/composer /usr/bin

# 安装 zip 扩展，composer 下载依赖需要
# 不同版本的 php 的 docker 容器安装扩展的方式不同， 具体请查看 https://hub.docker.com/_/php 和查阅相关资料
RUN sed -i "s@http://deb.debian.org/debian@http://mirrors.cloud.tencent.com/debian@g" /etc/apt/sources.list.d/debian.sources && \
    apt-get update && \
    apt-get install -y \
            libzip-dev \
            zip && \
    rm -rf /var/lib/apt/lists/* && \
    docker-php-ext-install zip

# composer 安装依赖
RUN composer config -g repo.packagist composer https://mirrors.cloud.tencent.com/composer/ && \
    composer require guzzlehttp/guzzle
```


## Python 使用腾讯云软件源加速
配置 `Dockerfile` 文件，使用腾讯云软件源加速 Python 软件包的下载。

```dockerfile
FROM python:3.13.0-bullseye

WORKDIR /app

COPY requirements.txt .

# 临时使用
RUN pip3 install -i https://mirrors.cloud.tencent.com/pypi/simple psutil && \
    pip3 install -r requirements.txt -i https://mirrors.cloud.tencent.com/pypi/simple

# 设置全局
RUN pip3 config set global.index-url https://mirrors.cloud.tencent.com/pypi/simple && \
    pip3 install requests

```

## npm 使用腾讯云软件源加速

运行命令配置 npm 为腾讯云软件源
```bash
npm config set registry https://mirrors.cloud.tencent.com/npm/
```

## yarn 使用腾讯云软件源加速

运行命令配置 yarn 为腾讯云软件源
```bash
yarn config set registry  https://mirrors.cloud.tencent.com/npm/
```

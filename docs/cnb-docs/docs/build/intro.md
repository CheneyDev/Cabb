---
title: 云原生构建介绍
permalink: https://docs.cnb.cool/zh/build/intro.html
summary: 云原生构建基于 Docker 生态，通过声明式语法、与环境代码同源管理、资源池化等特点助力软件构建。其有声明式构建环境、缓存、插件等特性，支持按需获取计算资源，还具备 CPU 自由、读秒克隆、缓存并发等高性能表现 。
---

基于 Docker 生态，对环境、缓存、插件进行抽象，通过声明式的语法，帮助开发者以更酷的方式构建软件。

- 声明式：声明式语法，可编程、易分享。
- 易管理：与代码一起，同源管理。
- 云原生：资源池化，屏蔽基础设施复杂性。

### 声明式的构建环境

```yaml{4}
main:
  push:
    - docker:
        image: node:20
      stages:
        - node -v
        - npm install
        - npm test
```

### 声明式的构建缓存

```yaml{5,6}
main:
  push:
    - docker:
        image: node:20
        volumes:
          - /root/.npm:copy-on-write
      stages:
        - node -v
        - npm install
        - npm test
```

### Docker 作为任务的运行环境

```yaml{5,8}
main:
  push:
    - stages:
        - name: run with node 20
          image: node:20
          script: node -v
        - name: run with node 21
          image: node:21
          script: node -v
```

### 基于 Docker 生态的插件

```yaml{5}
main:
  push:
    - stages:
        - name: hello world
          image: cnbcool/hello-world
```

### 按需获取计算资源

```yaml{4}
main:
  push:
    - runner:
        cpus: 64
      docker:
        image: node:20
      stages:
        - node -v
        - npm install
        - npm test
```

### 云原生开发

```yaml{5,6}
$:
  vscode:
    - runner:
        cpus: 64
      services:
        - vscode
      docker:
        image: node:20
        volumes:
          - node_modules:copy-on-write
      stages:
        - npm install
```

## 高性能

### CPU自由

- [runner.cpus](./grammar.md#cpus)

通过 `runner.cpus` 可按需声明需要的 CPU资源，最高可达 `64核`。

### 读秒克隆

基于 `OverlayFS` 的 [git-clone-yyds](https://cnb.cool/cnb/cool/git-clone-yyds) 可以在数秒内完成代码准备，轻松支持 `100GB+` 超大仓库。

### 缓存并发

- [copy-on-write](./grammar.md#volumes)

`copy-on-write` 可以实现缓存的写时复制，在并发场景下，无需再担心缓存读写冲突问题。

```yaml{7,8}
main:
  push:
    - runner:
        cpus: 64
      docker:
        image: node:20
        volumes:
          - /root/.npm:copy-on-write
      stages:
        - node -v
        - npm install
        - npm test
```

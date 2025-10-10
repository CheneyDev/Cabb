---
title: 构建集群
permalink: https://docs.cnb.cool/zh/build/saas/build-node.html
summary: 云原生构建通过将任务下发到指定 Docker 镜像的构建集群执行，用户可在流水线配置中通过`pipeline.runner.tags`指定不同集群，并通过`pipeline.runner.cpus`配置最大 CPU 核数，如`cnb:arch:amd64`支持 1 - 64 核，默认 8 核等不同配置选项 。
redirectFrom: /zh/build/build-node.html
---
当您在使用云原生构建时，本质上是将构建任务下发到各构建集群中执行。集群以指定的 Docker 镜像作为构建环境执行构建任务。

## 配置方式

在流水线配置里指定 `pipeline.runner.tags` 属性，即可选择不同构建集群。指定 `pipeline.runner.cpus` 属性，即可配置需使用的最大 `CPU` 核数。
云原生构建会以实际分配的 `核数` 乘以 流水线的耗时，作为流水线使用的核时。

官方可用构建集群 `tags` 及可声明的 `cpus` 如下：

1. `cnb:arch:amd64` 代表 `amd64` 架构的 `CPU` 集群
   - `cpus` 可配置范围为 1 ~ 64，默认为 8
2. `cnb:arch:arm64:v8` 代表 `arm64/v8` 架构的 `CPU` 集群
   - `cpus` 可配置范围为 1 ~ 16，默认为 8
3. `cnb:arch:amd64:gpu` 代表 `amd64` 架构的 `GPU` 集群
   - `cpus` 固定为 32
   - `GPU` 显存最大为 96GB，共享模式
4. `cnb:arch:amd64:gpu:L20` 代表 `amd64` 架构的 `GPU` 集群
   - `cpus` 固定为 16
   - `GPU` 显存最大为 48GB，共享模式

示例：

```yaml
main:
  push:
    # 指定在 amd64 架构构建集群上执行
    - runner:
        tags: cnb:arch:amd64
      stages:
        - name: uname
          script: uname -a
    # 指定在 arm64/v8 架构构建集群上执行
    - runner:
        tags: cnb:arch:arm64:v8
      stages:
        - name: uname
          script: uname -a

# 启动一个能使用 gpu 的远程开发环境
$:
  vscode:
    - runner:
        tags: cnb:arch:amd64:gpu
      services:
        - vscode
```

---
title: 简化配置文件
permalink: https://docs.cnb.cool/zh/build/simplify-configuration.html
summary: **简化配置文件**部分介绍了在 `.cnb.yml` 配置文件中如何利用 YAML 的高级语法，如锚点 `&` 、别名 `*` 和对象合并 `<<` 符号，以及变量引用来简化配置，同时还介绍了云原生构建为了跨文件复用流水线配置实现的特性。 

YAML 高级语法可用于减少单个文件内的重复，而云原生构建则通过文件引用和变量引用来解决跨文件的配置复用问题 。
---

## YAML 高级语法

由于`.cnb.yml`配置文件是 `YAML` 格式，所以可以利用更多 `YAML` 特性（如锚点 `&` 、别名 `*` 和对象合并 `<<` 符号）来简化配置文件。

一个简单的运用锚点和别名简化的例子如下：

```yaml
# pull_request 和 push 事件的流水线完全一致，这种方式可以减少重复
.pipeline: &pipeline
  docker:
    image: node:22
  stages:
    - name: install
      script: npm install
    - name: test
      script: npm test

main:
  pull_request:
    - <<: *pipeline
  push:
    - <<: *pipeline
```

支持多级嵌套：

```yaml
.jobs: &jobs
  - name: install
    script: npm install
  - name: test
    script: npm test

.pipeline: &pipeline
  docker:
    image: node:22
  stages: *jobs

main:
  pull_request:
    - <<: *pipeline
  push:
    - <<: *pipeline
```

:::tip

以上是 `YAML` 自带特性，仅在解析单个YAML文件时有效，不能跨文件使用。

:::

## 文件引用

为了方便复用流水线配置，`云原生构建` 实现了以下特性：

- [include](./grammar.md#include)：可以跨文件引用流水线模板。
- [imports](./grammar.md#Pipeline-imports)：可以跨文件引用变量。
- [optionsFrom](./grammar.md#optionsfrom)：可以跨文件引用内置任务参数。
- [settingsFrom](./grammar.md#settingsfrom)：可以跨文件引用插件任务参数。

详细文件引用说明见 [文件引用](./file-reference.md)

## 变量引用

为了解决YAML锚点和别名不能跨文件使用的问题。`云原生构建` 实现了以下特性：

[reference](./grammar.md#reference): 可以跨文件按属性路径引用变量。

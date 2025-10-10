---
title: 配置文件
permalink: https://docs.cnb.cool/zh/build/configuration.html
summary: 云原生构建配置文件（.cnb.yml）用于定义仓库事件触发构建任务的条件和步骤，采用YAML格式，存放在仓库根目录下，支持通过PR进行变更。配置文件包括触发分支和事件、执行环境和任务脚本，支持锚点复用和跨文件配置导入，提供语法检查和自动补全功能。
redirectFrom: /zh/configuration.html
---
## 简介

云原生构建配置文件(`.cnb.yml`)描述了当仓库发生一些事件时（有新的 Commit 被推送、有新的 PR 请求等），
`云原生构建` 是否应该启动构建任务，如果启动构建的话，构建任务的每一步分别做什么。

`云原生构建` 的配置文件格式是 `YAML`，这一点与业界主流 CI 服务相同。

这是一个简单的、可工作的 `云原生构建` 配置文件：

```yaml
main: # 触发分支
  push: # 触发事件
    - docker:
        image: node:22 # 流水线执行环境，可以指定任意docker镜像
      stages:
        - name: install
          script: npm install
        - name: test
          script: npm test
```

这个案例描述的流程如下：

1. 声明了在 `main` 分支在收到 `push` 事件时（即有新的 Commit 推送到 main 分支）
1. 会选择 Docker 镜像 `node:22` 作为执行环境
1. 依次执行任务 `npm install` 和 `npm test`

## 存放位置

`云原生构建` 约定的配置文件命名为 `.cnb.yml`，存放于仓库根目录下，配置文件即代码。

这意味着，配置文件可以通过 PR 进行变更，开源协作场景下，这十分重要。

构建流程纳入版本管理，与源代码保持相同的透明度和变更流程，修改历史很容易追溯。

## 基本语法结构

配置文件的基本语法结构如下所示：

```yaml
# 流水线结构：数组形式
main:
  push:
    # main 分支 - push 事件包含两条流水线：push-pipeline1 和 push-pipeline2
    - name: push-pipeline1 # 流水线名称，可省略
      stages:
        - name: job1
          script: echo 1
    - name: push-pipeline2 # 流水线名称，可省略
      stages:
        - name: job2
          script: echo 2

  pull_request:
    # main 分支 - pull_request 事件包含两条流水线：pr-pipeline1 和 pr-pipeline2
    - name: pr-pipeline1 # 流水线名称，可省略
      stages:
        - name: job1
          script: echo 1
    - name: pr-pipeline2 # 流水线名称，可省略
      stages:
        - name: job2
          script: echo 2
```

等价与以下写法：

```yaml
# 流水线结构：对象形式
main:
  push:
    # main 分支 - push 事件包含两条流水线：push-pipeline1 和 push-pipeline2
    push-pipeline1: # 流水线名称，必须唯一
      stages:
        - name: job1
          script: echo 1
    push-pipeline2: # 流水线名称，必须唯一
      stages:
        - name: job2
          script: echo 2

  pull_request:
    # main 分支 - pull_request 事件包含两条流水线：pr-pipeline1 和 pr-pipeline2
    pr-pipeline1: # 流水线名称，必须唯一
      stages:
        - name: job1
          script: echo 1
    pr-pipeline2: # 流水线名称，必须唯一
      stages:
        - name: job2
          script: echo 3
```

其中 `main` 表示分支名称， `push` 和 `pull_request` 表示[触发事件](./grammar.md#trigger-event)。

一个事件包含多个 `pipeline`，支持数组和对象两种形式，并发执行。

一个 `pipeline` 包含一组顺序执行的任务，在同一个构建环境（物理机、虚拟机或 Docker 容器）中执行。

详细语法说明可参考： [流水线语法](./grammar.md)

## 配置文件版本选择

同[代码版本选择](./trigger-rule.md#代码版本选择)

## 语法检查和自动补全

### VSCode

推荐使用 [云原生开发](../workspaces/intro.md) 书写配置文件，因为原生支持语法检查和自动补全，效果如下：

![yaml-auto](https://docs.cnb.cool/images/yaml-auto.gif)

本地开发配置方法，以 VSCode 为例：

先安装 `redhat.vscode-yaml` 插件，然后在 `settings.json` 中加入以下配置：

```text
"yaml.schemas": {
  "https://docs.cnb.cool/conf-schema-zh.json": ".cnb.yml"
},
```

### Jetbrains

![yaml-auto](https://docs.cnb.cool/images/jetbrains.png)

1. 「Settings」「Languages & Frameworks」「Schemas and DTDs」「JSON Schema Mappings」

2. 点击新增按钮，设置 「Name」（名称随意填写）

3. 设置「Schema file or URL」

4. 「Add mapping for a file」

```text
https://docs.cnb.cool/conf-schema-zh.json
```

## 配置复用

### anchor

在 YAML 中，锚点（Anchor）和引用（Alias）允许同一个文件中复用配置，从而避免重复并保持文件的简洁。

示例：

```yaml
# .cnb.yml
# 通用的流水线配置
.pipeline-config: &pipeline-config
  stages:
    - echo "do something"
main:
  push:
    # 引用 pipeline-config
    - *pipeline-config
dev:
  push:
    # 引用 pipeline-config
    - *pipeline-config
```

### include

利用 `include` 参数，可以在当前文件导入当前仓库或其他仓库上的文件。依此可以对配置文件进行拆分，方便复用和维护。

#### 使用示例

template.yml

```yaml
# template.yml
main:
  push:
    pipeline_2:
      env:
        ENV_KEY1: xxx
        ENV_KEY3: inner
      services:
        - docker
      stages:
        - name: echo
          script: echo 222
```

.cnb.yml

```yaml
# .cnb.yml
include:
  - https://cnb.cool/<your-repo-slug>/-/blob/main/xxx/template.yml

main:
  push:
    pipeline_1:
      stages:
        - name: echo
          script: echo 111
    pipeline_2:
      env:
        ENV_KEY2: xxx
        ENV_KEY3: outer
      stages:
        - name: echo
          script: echo 333
```

合并后的配置

```yaml{3,4,5,6,12,13,15,16}
main:
  push:
    pipeline_1: # key不存在，合并时新增
      stages:
        - name: echo
          script: echo 111
    pipeline_2:
      env:
        ENV_KEY1: xxx
        ENV_KEY2: xxx # key不存在，合并时新增
        ENV_KEY3: outer # 同名 key， 合并时覆盖
      services:
        - docker
      stages: # 数组在合并时，追加
        - name: echo
          script: echo 222
        - name: echo
          script: echo 333
```

#### 语法说明

```yaml
include:
  # 1、可直接传入配置文件路径
  - "https://cnb.cool/<your-repo-slug>/-/blob/main/xxx/template.yml"
  - "template.yml"
  # 2、可传入一个对象。
  # path: 配置文件路径
  # ignoreError：未读取到配置文件是否抛出错误。true 不抛出错误；false 抛出错误。默认为 false
  - path: "template1.yml"
    ignoreError: true
  # 3、可传入一个对象
  # config: 传入 yaml 配置
  - config:
      main:
        push:
          - stages:
              - name: echo
                script: echo "hello world"
```

#### 合并规则

不同文件的流水线配置合并规则：

- 数组(Array)和数组(Array)合并：子元素追加
- 对象(Map)和对象(Map)合并：同名 key 覆盖
- 数组(Array)和对象(Map)合并：仅保留数组
- 对象(Map)和数组(Array)合并：仅保留数组

引用配置文件权限控制参考 [配置文件引用鉴权](./file-reference.md)
  
:::tip
合并后的流水线配置会展示在构建详情页，与密钥仓库内容保护的理念不符，include 无法引用密钥仓库文件。
:::

:::tip

1. 本地的 .cnb.yml 会覆盖 include 中的配置，include 数组中后面的配置会覆盖前面的配置。
2. 支持嵌套 include，include 的本地文件路径相对于项目根目录。
3. 最多支持 include 50个配置文件。
4. 不支持引用 submodule 中的文件。
5. 不支持跨文件使用 YAML 锚点功能。

:::

### reference

YAML 不支持跨文件引用，`云原生构建` 通过扩展 YAML 自定义标签 `reference` 实现按属性路径引用变量值，可结合 [include](#include) 跨文件复用配置。

:::tip

1. 第一层同名变量会被覆盖，不会合并。本地的 `.cnb.yml` 会覆盖 `include` 中的变量，`include` 数组中后面的变量会覆盖前面的变量。
2. `reference` 支持嵌套引用，最多 10 层。

:::

#### 示例 {#reference-example}

a.yml

```yaml
.val1:
  echo1: echo hello
.val2:
  friends:
    - one:
      name: tom
      say: !reference [.val1, echo1]
```

.cnb.yml

```yaml
include:
  - ./a.yml
.val3:
  size: 100
main:
  push:
    - stages:
        - name: echo hello
          script: !reference [.val2, friends, "0", say]
        - name: echo size
          env:
            SIZE: !reference [".val3", "size"]
          script: echo my size ${SIZE}
```

解析后相当于：

```yaml
main:
  push:
    - stages:
        - name: echo hello
          script: echo hello
        - name: echo size
          env:
            SIZE: 100
          script: echo my size ${SIZE}
```

#### 进阶示例

可以将流水线作为整体配置引用：

```yaml
.common-pipeline:
  - stages:
      - name: echo
        script: echo hello

main:
  push: !reference [.common-pipeline]
test:
  push: !reference [.common-pipeline]
```

解析后相当于：

```yaml
main:
  push:
    - stages:
        - name: echo
          script: echo hello
test:
  push:
    - stages:
        - name: echo
          script: echo hello
```

#### VSCode 配置

安装 `VSCode YAML` 插件后，
为了在 `VSCode` 编写带自定义标签 `reference` 的 `YAML` 文件时不报错，需要如下配置：

setting.json

```json
{
  "yaml.customTags": ["!reference sequence"]
}
```

:::tip
  为避免编写时 `YAML` 插件根据 `Schema` 把第一层变量名当做分支名，
  有错误提示，`reference` 所在的第一层变量名可用 `.` 开头，如：`.var`。
:::

---
title: 自定义部署流程
permalink: https://docs.cnb.cool/zh/build/deploy.html
summary: 云原生构建支持自定义部署流程，用户可通过在仓库根目录下添加 `.cnb/tag_deploy.yml` 文件来配置部署环境、审批流程和部署前置条件。部署环境包括 `development`、`staging` 和 `production`，每个环境可定义按钮、审批流程、部署前置条件等，流水线基于当前 tag 对应的代码进行部署操作，且部署需有仓库写权限和推送 Tag 权限的用户才能进行 。
redirectFrom: /zh/deploy.html
---

云原生构建支持自定义部署流程，通过自定义部署环境、审批流程、部署流水线，实现自动化的部署流程。

操作示例：

![](https://docs.cnb.cool/images/build/deploy-intro.png)

## 自定义部署环境

在仓库根目录下添加 `.cnb/tag_deploy.yml` 文件用于配置部署环境。
如下示例中定义了 `development`、`staging`、`production` 三种环境，用户可在页面上选择需要部署的环境类型。

```yaml
# .cnb/tag_deploy.yml
environments:
  # name: 环境名，点击该环境对应的部署按钮将触发 .cnb.yml 中的 tag_deploy.development 事件流水线
  - name: development
    description: Development environment
    # 环境变量（触发流水线时，会将环境变量传入流水线，包括部署流水线、web_trigger 流水线）
    env:
      name: development
      # CNB_BRANCH: 环境变量，部署事件中，为 tag 名
      tag_name: $CNB_BRANCH

  - name: staging
    description: Staging environment
    env:
      name: staging
      # CNB_BRANCH: 环境变量，部署事件中，为 tag 名
      tag_name: $CNB_BRANCH

  - name: production
    description: Production environment
    # 环境变量（触发流水线时，会将环境变量传入流水线，包括部署流水线、web_trigger 流水线）
    env:
      name: production
      # CNB_BRANCH: 环境变量，部署事件中，为 tag 名
      tag_name: $CNB_BRANCH
    button:
      - name: 创建审批单
        # 如存在，则将作为流水线 title，否则流水线使用默认 title
        description: 自动创建审批单流程
        # 需要在 .cnb.yml 中自定义 web_trigger_approval 事件流水线
        event: web_trigger_approval
        # 权限控制，不配置则有仓库写权限的用户可触发构建
        # 如果配置，则需要有仓库写权限，并且满足 roles 或 users 其中之一才有权限触发构建
        permissions:
          # roles 和 users 配置其中之一或都配置均可，二者满足其一即可
          roles:
            - owner
            - developer
          users:
            - name1
            - name2
        # 传给 web_trigger_approval 事件流水线的环境变量
        # 可继承上一级别环境变量，优先级高于上一级别环境变量
        env:
          name1: value1
          name2: value2

    # 部署前置条件检查（支持对环境、元数据、审批流程的检查），满足所有前置条件才可进行部署操作
    require:
      # 1 对部署环境是否满足要求的检查

      # 1.1 要求 development 环境部署成功
      - environmentName: development

      # 1.2 要求 staging 环境部署成功 30 分钟后
      - environmentName: staging
        after: 1800

      # 2 对元数据是否满足要求的检查

      # 2.1 键值 key1 对应的 value 不为空，即有值
      - annotation: key1

      # 2.2 键值 key1 对应的 value 值需等于 value1
      - annotation: key1
        expect:
          eq: value1

      # 2.3 键值 key2 对应的 value 值需大于 1 且小于 10
      - annotation: key2
        expect:
          and:
            gt: 1
            lt: 10
        # 自定义按钮，点击可触发执行 web_trigger_annotation 事件。
        # 可定义与 require 信息有关的按钮事件，当 require 满足条件后隐藏按钮
        button:
          - name: 生成元数据
            event: web_trigger_annotation
            # 如存在，则将作为流水线 title，否则流水线使用默认 title
            description: 生成元数据流程
            # 权限控制，不配置则有仓库写权限的用户可触发构建
            # 如果配置，则需要有仓库写权限，并且满足 roles 或 users 其中之一才有权限触发构建
            permissions:
              # roles 和 users 配置其中之一或都配置均可，二者满足其一即可
              roles:
                - owner
                - developer
              users:
                - name1
                - name2
            # 传给 web_trigger_annotation 事件流水线的环境变量
            # 可继承上一级别环境变量，优先级高于上一级别环境变量
            env:
              name1: value1
              name2: value2

      # 3 对审批流程是否满足要求的检查（可按以下方式自定义审批流程）
      # - 审批顺序：如下 1、2、3 审批流程需按顺序进行，即 1 审批通过，2 才能进行审批。1、2、3 审批流程全部通过才算通过审批
      # - 审批操作：包括 同意、拒绝。一人同意即算通过。如果拒绝，其他审批人无法再操作，直到拒绝的审批人再修改审批结果为同意

      # 3.1 按用户名审批，其中一人审批通过即可
      - approver:
          users:
            - user1
            - user2
            - user3
        title: 测试审批

      # 3.2 按角色审批，其中一人审批通过即可
      - approver:
          roles:
            - developer
            - master
        title: 开发审批

      # 3.3 按用户名或角色审批（审批人满足 users 或 roles 其一即可），其中一人审批通过才行
      - approver:
          users:
            - user4
            - user5
          roles:
            - master
            - owner
        title: 运维审批

    # 自定义部署按钮（缺省值：默认展示一个部署按钮）
    # 使用场景：有多个不同模块（例如仓库、CI、制品库等），需要分开独立部署时，可以配置多个不同的按钮
    # 注意：部署流水线中要区分是哪个模块，可以通过传入流水线的环境变量来区分
    deploy:
      - name: 部署按钮名1
        description: 部署按钮描述
        # 环境变量（触发部署流水线时，会将环境变量传入流水线），优先级高于上一级 env
        env:
          name1: value1
          name2: value2
      - name: 部署按钮名2
        description: 部署按钮描述
        # 环境变量（触发部署流水线时，会将环境变量传入流水线），优先级高于上一级 env
        env:
          name1: value1
          name2: value2
```

- `name`: 必填，环境名，需唯一。例如 `name: development`，点击该环境对应的部署按钮，将触发 `.cnb.yml` 中的 `tag_deploy.development` 事件流水线
- `description`: 选填，环境描述
- `env`: 选填，传给部署流水线的环境变量。用户可根据需要传入需要的环境变量。
- `button`: 选填，对象数组格式，自定义按钮。点击按钮可触发云原生构建流水线，执行参数 event 对应的事件。

  - `name`: 必填，按钮名。
  - `description`: 选填，按钮描述。如存在，则将作为流水线 title，否则流水线使用默认 title。
  - `event`: 必填，自定义事件，仅支持 web_trigger 事件。
  - `env`: 选填，传给 web_trigger 流水线的环境变量，可继承上一级别环境变量，优先级高于上一级别环境变量。
  - `permissions`: 选填，权限控制，满足 `users` 或 `roles` 其中之一即有权限触发构建（还需要有仓库写权限）。如果未配置 `permissions`，则有仓库写权限即可出发构建
    - `users`: 选填，`Array<String>`，用户名数组。可定义多个。
    - `roles`: 选填，`Array<String>`，仓库角色数组。可定义多种仓库角色。
    `owner`(负责人)、`master`(管理员Administrator)、`developer`(开发者)、`reporter`(助手)、`guest`(访客)
  
- `deploy`: 选填，对象数组格式，自定义部署按钮。点击按钮可触发云原生构建流水线，执行部署事件（`tag_deploy.*`）。

  - `name`: 必填，按钮名。
  - `description`: 选填，按钮描述。
  - `env`: 选填，传给部署流水线的环境变量，优先级高于上一级 env。
  
- `require`: 选填，对象数组格式。部署的前置条件，需满足了前置条件（部署环境要求、元数据要求、审批流程）才可进行部署操作。

  1、部署环境要求的参数包括

  - `environmentName`: 必填。环境名。
  - `after`: 选填。时间，单位 s(秒)。表示 `environmentName` 的环境部署成功后 after 时间后才算满足前置条件。
  - `description`: 选填。`require` 的描述信息，附注用户理解 `require` 要求的内容。
  - `button`: 选填。自定义按钮，点击可触发执行 event 传入的事件。可定义与 require 信息有关的按钮事件，注意：当 require 满足条件后隐藏按钮。

    - `name`: 必填，按钮名。
    - `event`: 必填，自定义事件，仅支持 `web_trigger_*` 事件。
    - `description`: 选填，按钮描述。如存在，则将作为流水线 title，否则流水线使用默认 title。
    - `env`: 选填，传给 web_trigger 流水线的环境变量，可继承上一级别环境变量，优先级高于上一级别环境变量。
    - `permissions`: 选填，权限控制，满足 `users` 或 `roles` 其中之一即有权限触发构建（还需要有仓库写权限）。如果未配置 `permissions`，则有仓库写权限即可出发构建
      - `users`: 选填，`Array<String>`，用户名数组。可定义多个。
      - `roles`: 选填，`Array<String>`，仓库角色数组。可定义多种仓库角色。
      `owner`(负责人)、`master`(管理员Administrator)、`developer`(开发者)、`reporter`(助手)、`guest`(访客)

  2、元数据要求的参数包括

  - `annotation`: 必填。元数据的 `key` 值。
  - `expect`: 选填。对元数据的 `value` 值的要求。对象格式，支持 `eq`、`ne`、`gt`、`lt`、`gte`、`lte`、`and`、`or`、`reg` 操作符。
    - `eq`: 等于
    - `ne`: 不等于
    - `gt`: 大于
    - `lt`: 小于
    - `gte`: 大于等于
    - `lte`: 小于等于
    - `and`: 与
    - `or`: 或
    - `reg`: 能和正则表达式匹配
  - `description`: 选填。`require` 的描述信息，附注用户理解 `require` 要求的内容。
  - `button`: 选填。自定义按钮，点击可触发执行 event 传入的事件。可定义与 require 信息有关的按钮事件，注意：当 require 满足条件后隐藏按钮。

    - `name`: 必填，按钮名。
    - `description`: 选填，按钮描述。如存在，则将作为流水线 title，否则流水线使用默认 title。
    - `event`: 必填，自定义事件，仅支持 web_trigger 事件。
    - `env`: 选填，传给 web_trigger 流水线的环境变量，可继承上一级别环境变量，优先级高于上一级别环境变量。
    - `permissions`: 选填，权限控制，满足 `users` 或 `roles` 其中之一即有权限触发构建（还需要有仓库写权限）。如果未配置 `permissions`，则有仓库写权限即可出发构建
      - `users`: 选填，`Array<String>`，用户名数组。可定义多个。
      - `roles`: 选填，`Array<String>`，仓库角色数组。可定义多种仓库角色。
      `owner`(负责人)、`master`(管理员Administrator)、`developer`(开发者)、`reporter`(助手)、`guest`(访客)

  3、审批流程要求的参数包括

  - `approver`: 必填，审批人定义，满足 `users` 或 `role` 的审批人中，一人审批通过即可。
    - `users`: 用户名数组。可定义多个审批人。
    - `roles`: 仓库角色数组。可定义多种仓库角色。
    `owner`(负责人)、`master`(管理员Administrator)、`developer`(开发者)、`reporter`(助手)、`guest`(访客)
  - `title`: 选填，审批标题，如 `测试审批`。

## 自定义部署前置条件

对于每个环境可定义部署前置条件，只有满足所有前置条件才可进行部署操作。可定义如下三种前置条件：

- 环境部署要求：要求指定环境已经部署成功，且满足 `after` 部署成功时间要求
- 元数据值要求：要求指定元数据对应的值是否满足要求
- 审批流程要求：可自定义审批流程指定审批人，并进行审批操作，当全部审批流程都审批通过后，才算满足要求

### 环境部署前置条件示例

```yaml
# .cnb/tag_deploy.yml
environments:
  - name: development
    description: Development environment
    env:
      name: development
      tag_name: $CNB_BRANCH

  - name: staging
    description: Staging environment
    env:
      name: staging
      tag_name: $CNB_BRANCH
    require:
      # 要求 development 环境部署成功
      - environmentName: development

  - name: production
    description: Production environment
    require:
      # 要求 staging 环境部署成功 30 分钟后
      - environmentName: staging
        after: 1800
```

### 元数据前置条件示例

```yaml
# .cnb/tag_deploy.yml
environments:
  - name: production
    description: Production environment
    require:
      # 对元数据是否满足要求的检查

      # 键值 key1 对应的 value 不为空，即有值
      - annotation: key1

      # 键值 key2 对应的 value 值需等于 value1
      - annotation: key2
        expect:
          eq: value2

      # 键值 key3 对应的 value 值需大于 1 且小于 10
      - annotation: key3
        expect:
          and:
            gt: 1
            lt: 10
        # 自定义按钮，点击可触发执行 web_trigger_annotation 事件。
        # 可定义与 require 信息有关的按钮事件，当 require 满足条件后隐藏按钮
        button:
          - name: 生成元数据
            event: web_trigger_annotation
            # 如存在，则将作为流水线 title，否则流水线使用默认 title
            description: 生成元数据流程
            # 传给 web_trigger_annotation 事件流水线的环境变量
            # 可继承上一级别环境变量，优先级高于上一级别环境变量
            env:
              name1: value1
              name2: value2
```

### 审批流程前置条件示例

可自定义审批流程和指定审批人。有权限的审批人可进行审批操作（同意、拒绝）。全部流程审批通过后，即算满足要求

```yaml
# .cnb/tag_deploy.yml
environments:
  - name: production
    description: Production environment
    require:
      # 对审批流程是否满足要求的检查（可按以下方式自定义审批流程）
      # - 审批顺序：如下 1、2、3 审批流程需按顺序进行，即 1 审批通过，2 才能进行审批。1、2、3 审批流程全部通过才算通过审批
      # - 审批操作：包括 同意、拒绝。一人同意即算通过。如果拒绝，其他审批人无法再操作，直到拒绝的审批人再修改审批结果为同意

      # 按用户名审批，其中一人审批通过即可
      - approver:
          users:
            - user1
            - user2
            - user3
        title: 测试审批

      # 按角色审批，其中一人审批通过即可
      - approver:
          roles:
            - developer
            - master
        title: 开发审批

      # 按用户名或角色审批（审批人满足 users 或 roles 其一即可），其中一人审批通过才行
      - approver:
          users:
            - user4
            - user5
          roles:
            - master
            - owner
        title: 运维审批
```

## 自定义部署流水线

如下示例定义了三种环境的部署事件流水线，
当在页面中选择部署 `development` 环境时，则触发 `tag_deploy.development` 事件。
流水线基于当前 tag 对应的代码进行部署操作。

```yaml
# .cnb.yml
$:
  tag_deploy.development:
    - name: dev
      stages:
        - name: 部署环境名
          script: echo $name
        - name: tag 名
          script: echo $tag_name
  tag_deploy.staging:
    - name: staging
      stages:
        - name: 部署环境名
          script: echo $name
        - name: tag 名
          script: echo $tag_name
  tag_deploy.production:
    - name: production
      stages:
        - name: 部署环境名
          script: echo $name
        - name: tag 名
          script: echo $tag_name
```

示例中的流水线事件名和部署环境类型对应关系如下：

- `tag_deploy.development`：`development`
- `tag_deploy.staging`：`staging`
- `tag_deploy.production`：`production`

### 自定义按钮触发的 web_trigger 事件

`tag_deploy.yml` 中的自定义按钮，仅支持触发 [web_trigger事件](./grammar.md#web-trigger-自定义事件) 事件。
如下流水线配置中，`web_trigger_annotation` 事件执行时，会进行上传[元数据](../repo/annotations.md)操作。

```yaml
# .cnb.yml
$:
  # 自定义按钮可触发的事件
  web_trigger_annotation:
    - stages:
        - name: 上传元数据
          image: cnbcool/annotations:latest
          settings:
            data: |
              key1=value1
              key2=value2
```

## 部署权限说明

需有 `仓库写权限` 且有 `推送 Tag` 权限的用户才能进行部署操作。

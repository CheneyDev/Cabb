---
title: 手动触发流水线
permalink: https://docs.cnb.cool/zh/build/web-trigger.html
summary: 云原生构建支持在`代码-分支详情页`配置自定义按钮，通过点击按钮触发`web_trigger`事件流水线，并支持输入环境变量。需在仓库根目录下创建`.cnb/web_trigger.yml`文件进行按钮配置，包括按钮显示规则、事件名、环境变量及权限控制等；同时，在`.cnb.yml`中定义`web_trigger`事件流水线。只有具备仓库写权限的用户可触发执行。
---

云原生构建支持在页面上配置自定义按钮，点击按钮触发执行流水线（仅支持触发 [页面自定义事件](./grammar.md#web_trigger)）。同时支持在页面上输入环境变量。

效果如下：

![](https://docs.cnb.cool/images/build/web-trigger-btn.png)

:::tip
目前仅支持在 `代码-分支详情页` 上配置的自定义按钮。
:::

## 自定义按钮配置

在仓库根目录下增加 `.cnb/web_trigger.yml` 文件用于配置自定义按钮。
注意 `.cnb/web_trigger.yml` 配置文件中不支持 `include` 和 `imports` 语法。

```yaml
# .cnb/web_trigger.yml
branch:
  # 如下按钮在分支名以 release 开头的分支详情页面显示
  - reg: "^release"
    buttons:
      - name: 按钮名1
        # 如存在，则将作为流水线 title，否则流水线使用默认 title
        description: 按钮描述
        event: web_trigger_one # 触发的 CI 事件名
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
        # 环境变量
        env:
          # 默认传入的环境变量，其中 key 值（a,b,c）为环境变量名，支持如下两种格式
          a: 1
          b: 2
          c:
            # 环境变量别名
            name: 变量c
            # 环境变量值
            value: 3
        # 可输入环境变量，可覆盖上述 env 的变量配置
        inputs:
          # 目前支持以下三种格式：输入框（input）、多行文本输入框（textarea）、下拉选择框（select 支持单选和多选），switch 开关
          # 其中 key 值（var1、var2、var3、var4、var5）为环境变量名
          var1:
            # 输入框
            name: 变量1
            description: 输入变量1
            required: true # 是否必填
            type: input
            default: 默认值1
          var2:
            # 输入框
            name: 变量2
            description: 输入变量2
            required: true
            type: textarea
            default: 默认值2
          var3:
            # 单选下拉选择框
            name: 变量3
            description: 输入变量3
            required: false
            type: select
            default: value1
            options:
              - name: 选项1
                value: value1
                description: 选项1描述
              - name: 选项2
                value: value2
                description: 选项2描述
          var4:
            # 多选下拉选择框
            name: 变量4
            description: 输入变量4
            required: false
            type: select
            # 是否支持多选，多选结果用分号分隔
            multiple: true
            default: value1,value2
            options:
              - name: 选项1
                value: value1
                description: 选项1描述
              - name: 选项2
                value: value2
                description: 选项2描述
              - name: 选项3
                value: value3
                description: 选项3描述
          var5:
            # switch 开关
            name: 变量5
            description: 选择变量5
            required: false
            type: switch
            default: value1
            options:
              - name: 选项1
                value: value1
                description: 选项1描述
              - name: 选项2
                value: value2
                description: 选项2描述

  # 如下按钮在分支名以 dev 开头的分支详情页面显示
  - reg: "^dev"
    buttons:
      - name: 按钮名2
        description: 按钮描述
        event: web_trigger_two
      - name: 按钮名3
        description: 按钮描述
        event: web_trigger_three

  # 如下自定义按钮在所有分支详情页面显示
  - buttons:
      - name: 按钮名4
        description: 按钮描述
        event: web_trigger_four
```

目前仅支持配置分支详情页面的自定义按钮，yaml 中键值为 `branch`，值为数组格式，数组元素定义如下：

- `reg`: 选填，`String`，正则表达式，用于匹配分支名（仅匹配到的分支显示 `buttons` 配置的自定义按钮），
未填则匹配全部分支
- `buttons`: 必填，`Array<Button>`，自定义按钮定义。`Button` 类型定义如下
  - `name`: 必填，`String`，自定义按钮名
  - `description`: 选填，`String`，按钮描述。如存在，则将作为流水线 title，否则流水线使用默认 title
  - `event`: 必填，`String`，事件名
  - `env`: 选填，`Object<String, String|EnvType>`，传给 web_trigger 自定义事件流水线的默认环境变量，不支持编辑。
  对象键值为环境变量名；
  对象值支持两种格式，`String` 和 `EnvType`。`EnvType` 类型定义如下：
    - `name`: 必填，`String`，环境变量别名。非传给流水线的环境变量名
    - `value`: 必填，`String`，环境变量值
  - `permissions`: 选填，权限控制，满足 `users` 或 `roles` 其中之一即有权限触发构建（还需要有仓库写权限）。如果未配置 `permissions`，则有仓库写权限即可出发构建
    - `users`: 选填，`Array<String>`，用户名数组。可定义多个。
    - `roles`: 选填，`Array<String>`，仓库角色数组。可定义多种仓库角色。
    `owner`(负责人)、`master`(管理员Administrator)、`developer`(开发者)、`reporter`(助手)、`guest`(访客)
  - `inputs`: 选填，`Object<String, Input>`，可手动输入的环境变量，对象的键值为变量名，`Input` 类型定义如下：
    - `description`：选填，`String`，描述信息
    - `required`：选填，`Boolean`，是否必填
    - `type`：选填，`String`，输入框类型，可选 `input`、`textarea`、`select`、`switch`，默认为 `input`
    - `default`: 选填，`String`，默认值
    - `multiple`: 选填，`Boolean`，是否支持多选，仅当 type=`select` 时有效，多选结果用分号分隔。单个选项的 `value` 中需避免出现逗号
    - `options`: 选填，`Array<Option>`，当 `type: select` 或 `type: switch` 时的选项。`Option` 类型定义如下：
      - `name`: 必填，`String`，选项名
      - `value`: 必填，`String`，选项值，作为环境变量值
      - `description`: 选填，`String`，选项描述信息

## 自定义 web_trigger 流水线

`.cnb/web_trigger.yml` 中的自定义按钮，仅支持触发 [web_trigger](./grammar.md#web_trigger) 事件。

`web_trigger` 事件流水线在 `.cnb.yml` 中配置

```yaml
# .cnb.yml

# 匹配以 release 开头的分支名
release*:
  # 自定义按钮可触发的事件
  web_trigger_one:
    - stages:
        - name: 输出环境变量
          script:
            - echo $a
            - echo $b
            - echo $var1
            - echo $var2
            - echo $var3

# 匹配以 dev 开头的分支名
dev*:
  web_trigger_two:
    - stages:
        - name: 执行任务
          script: echo "job"

  web_trigger_three:
    - stages:
        - name: 执行任务
          script: echo "job"

# 匹配所有分支名
"**":
  web_trigger_four:
    - stages:
        - name: 执行任务
          script: echo "job"
```

## 权限说明

仅有 `仓库写权限` 的用户可点击自定义按钮执行 `web_trigger` 流水线。

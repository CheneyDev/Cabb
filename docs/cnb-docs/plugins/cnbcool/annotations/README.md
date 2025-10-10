# GIT 元数据

上传/获取/删除元数据，元数据是对 tag 或 commit 的注解，用于存储关联数据。

## 镜像

`cnbcool/annotations:latest`

## 支持的事件

以下事件支持 `上传`、`获取`、`删除` 元数据操作：

- push
- branch.create
- branch.delete
- pull_request.target
- pull_request.approved
- pull_request.changes_requested
- pull_request.mergeable
- pull_request.merged
- tag_push
- vscode
- auto_tag
- tag_deploy.*
- api_trigger*
- web_trigger*

以下事件仅支持 `获取` 元数据，不支持 `上传` 和 `删除` 元数据`：

- pull_request

[事件含义说明](https://docs.cnb.cool/event.html#pull_request)

## 元数据

### 格式说明

`key`:`value` 形式

### 元数据存储载体

目前可以针对 `tag` 或 `commit` 存储元数据：

- `tag`: 元数据是对 `tag` 的注解
- `commit`: 元数据是对 `commitId` 的注解

## 参数说明

- `type`: 操作类型，默认为 `ADD`。
  - `ADD`：上传元数据，其他可用参数：`data`、`fromJsonFile`。
  - `DELETE`：删除元数据，其他可用参数：`key`
  - `GET`：获取元数据，其他可用参数：`toFile`
- `data`: 需要上传的元数据。`type: ADD` 时有效。
支持以下两种格式（二选一，优先用 `\n`），当需要使用 `;` 分隔时，不能含有 `\n`。
元数据 `key` 不能为空，为空会被过滤掉，且只能包含`0-9`、`a-z`、`A-Z`、`_`、`-`。
元数据 `value` 可能包含 `\n` 或 `;` 时，建议将 value 编码，或使用 `fromJsonFile` 传入数据，避免得到非预期结果。
  - 用分号分隔：`key1=value1;key2=value2;key3=value3`
  - 用换行分隔：`key1=value1\nkey2=value2\nkey3=value3`
- `key`：支持多个，用英文分号、英文逗号、换行分隔。`type: DELETE` 时，表示需要删除的元数据 `key`。
- `toFile`: `type: GET` 时有效。将查询出的数据以 JSON 字符串格式存入指定文件，传入相对路径（如 `text.json`）
- `fromJsonFile`：`type: ADD` 时有效。
JSON 数据（如：`{"key1": "value1","key2":"value2"}`）存储的文件相对路径（如 `text.json`）。
当 `data` 和 `fromJsonFile` 同时传时，都有效，但两者同 `key` 数据，优先使用 `data` 中的数据。
- `tag`: tag 名称，非必传。操作指定 tag 的元数据。
- `commit`: `commitID`，长 `hash`，非必传。操作指定 `commit` 的元数据。

其中 `tag` 或 `commit` 都不传时，元数据载体默认为：

- `tag_push`和`tag_deploy.*` 事件: 默认对 `tag` 的元数据进行操作。`tag` 取环境变量 `CNB_BRANCH`
- 其他事件：默认对 `commit` 的元数据进行操作。`commit` 取环境变量 `CNB_COMMIT`

## 在 云原生构建 中使用

### 上传元数据

```yaml
main:
  # push 事件中上传元数据，数据是对 commitId 的注解
  push:
    - stages:
      - name: 上传元数据
        image: cnbcool/annotations:latest
        settings:
          # 传入数据的四种方式
          # 方式一：多行文本格式示例，实际是\n分隔
          data: |
            key1=value1
            key2=value2
          # 方式二：分号分隔
          # data: key1=value1;key2=value2
          # 方式三：\n 分隔
          # data: key1=value1\nkey2=value2
          # 方式四：从文件中读取数据
          # fromFile: text.json
          type: ADD

$:
  # tag_push 事件中上传元数据，数据是对 tag 的注解
  tag_push:
    - stages:
      - name: 上传元数据
        image: cnbcool/annotations:latest
        settings:
          # 传入数据的四种方式
          # 方式一：多行文本格式示例，实际是\n分隔
          data: |
            key1=value1
            key2=value2
          # 方式二：分号分隔
          # data: key1=value1;key2=value2
          # 方式三：\n 分隔
          # data: key1=value1\nkey2=value2
           # 方式四：从文件中读取数据
          # fromFile: text.json
          type: ADD

main:
  # 非 tag_push/tag_deploy.* 事件中上传指定 tag 的元数据
  api_trigger_test:
    - stages:
      - name: 上传元数据
        image: cnbcool/annotations:latest
        settings:
          # 传入数据的四种方式
          # 方式一：多行文本格式示例，实际是\n分隔
          data: |
            key1=value1
            key2=value2
          # 方式二：分号分隔
          # data: key1=value1;key2=value2
          # 方式三：\n 分隔
          # data: key1=value1\nkey2=value2
           # 方式四：从文件中读取数据
          # fromFile: text.json
          type: ADD
          tag: v1.0.0

```

### 删除元数据

```yaml
main:
  # push 事件中删除元数据，删除的是 commitId 的元数据
  push:
    - stages:
      - name: 删除元数据
        image: cnbcool/annotations:latest
        settings:
          # 支持分号、\n 分隔，也可写成多行文本
          key: key1;key2;key3
          type: DELETE

$:
  # tag_push 事件中删除元数据，删除的是 tag 的元数据
  tag_push:
    - stages:
      - name: 删除元数据
        image: cnbcool/annotations:latest
        settings:
          # 支持分号、\n 分隔，也可写成多行文本
          key: key1;key2;key3
          type: DELETE

main:
  # 非 tag_push/tag_deploy.* 事件中删除指定 tag 的元数据
  api_trigger_test:
    - stages:
      - name: 删除元数据
        image: cnbcool/annotations:latest
        settings:
          # 支持分号、\n 分隔，也可写成多行文本
          key: key1;key2;key3
          type: DELETE
          tag: v1.0.0
```

### 获取元数据

```yaml
main:
  # push 事件中获取元数据，获取的是 commitId 的元数据
  push:
    - stages:
      - name: 获取元数据
        image: cnbcool/annotations:latest
        settings:
          type: GET
        exports:
          # annotations 为所有元数据的 json 字符串格式
          # 例如: "{"key1":"value1","key2":"value2","key3":"value3"}"
          annotations: ANNOTATIONS
          key1: KEY1
          key2: KEY2
          key3: KEY3
      - name: 输出元数据
        script: 
          - echo $ANNOTATIONS
          - echo $KEY1
          - echo $KEY2
          - echo $KEY3

$:
  # tag_push 事件中获取元数据，获取的是 tag 的元数据
  tag_push:
    - stages:
      - name: 获取所有元数据
        image: cnbcool/annotations:latest
        settings:
          type: GET
        exports:
          # annotations 为所有元数据的 json 字符串格式
          # 例如: "{"key1":"value1","key2":"value2","key3":"value3"}"
          annotations: ANNOTATIONS
          key1: KEY1
          key2: KEY2
          key3: KEY3
      - name: 输出元数据
        script: 
          - echo $ANNOTATIONS
          - echo $KEY1
          - echo $KEY2
          - echo $KEY3

main:
  # 非 tag_push/tag_deploy.* 事件中获取指定 tag 的元数据
  api_trigger_test:
    - stages:
      - name: 获取所有元数据
        image: cnbcool/annotations:latest
        settings:
          type: GET
          tag: v1.0.0
        exports:
          # annotations 为所有元数据的 json 字符串格式
          # 例如: "{"key1":"value1","key2":"value2","key3":"value3"}"
          annotations: ANNOTATIONS
          key1: KEY1
          key2: KEY2
          key3: KEY3
      - name: 输出元数据
        script: 
          - echo $ANNOTATIONS
          - echo $KEY1
          - echo $KEY2
          - echo $KEY3
```

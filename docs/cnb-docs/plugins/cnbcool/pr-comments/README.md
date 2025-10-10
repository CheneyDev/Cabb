# PR 评论插件

上传 PR 评论，支持普通评论和评审评论。

普通评论指的是不进入评审，对整个 PR 进行评论；评审评论指的是进入评审，对整个 PR，或某个文件，或文件的具体行，进行评论，并且完成一次评审。

## 镜像

`cnbcool/pr-comments:latest`

## 支持的事件

- pull_request
- pull_request.target
- pull_request.approved
- pullrequest.changesrequested
- pull_request.mergeable
- pull_request.merged

## 参数说明

- `type`：评论类型，非必填，值为`COMMENT`或者`REVIEW`。当不传入该值时，如果同时没有传入`comments`，则默认为`COMMENT`，否则为`REVIEW`。
  - `COMMENT`：普通评论
  - `REVIEW`：评审评论
- `content`: 评论内容，类型为`string`。`content`和`fromFile`必填其一。
- `fromFile`: 评论内容文件相对路径，类型为`string`。`content`和`fromFile`必填其一。当`content`和`fromFile`同时存在时，以`fromFile`为准。
- `comments`: 评审评论信息，类型为`string`，JSON 字符串（数组的 JSON 字符串格式），非必填。目前只支持一条评论，因此数组长度最大只能为1。

其中，`comments`中的参数说明如下：

- `path`：文件绝对路径，类型为`string`，必填。
- `content`: 评论内容，类型为`string`。`content`和`fromFile`必填其一。
- `fromFile`: 评论内容文件相对路径，类型为`string`。`content`和`fromFile`必填其一。当`content`和`fromFile`同时存在时，以`fromFile`为准。
- `start_line`：评论起始行，类型为`number`，非必填。当该值不填，但`end_line`存在时，默认值为`end_line`。
- `end_line`：评论终止行，类型为`number`，非必填。当该值不填，但`start_line`存在时，默认值为`start_line`。
- `start_side`: 评论起始方向，类型为`string`，非必填，可选值为`LEFT`或`RIGHT`。不填默认传`LEFT`。
- `end_side`: 评论终止方向，类型为`string`，非必填，可选值为`LEFT`或`RIGHT`。不填默认传`LEFT`。

关于评审评论几种类型的具体说明：

- 当`comments`不存在时，表示对整个 PR 进行评审评论。
- 当`comments`只填`path`和`body`时，表示对整个文件进行评审评论。
- 当`start_line`或者`end_line`存在时，表示对文件的具体行进行评审评论。

## 在 云原生构建 中使用

### 上传普通评论

```yaml
main:
  pull_request:
    - stages:
      - name: 上传普通评论
        image: cnbcool/pr-comments:latest
        settings:
          content: "test"
          fromFile: "test.txt"
```

### 上传评审评论

```yaml
main:
  pull_request:
    - stages:
      - name: 上传评审评论——对整个 pr 进行评审评论
        image: cnbcool/pr-comments:latest
        settings:
          content: "test"
          fromFile: "test.txt"
          type: 'REVIEW'

      - name: 上传评审评论——对文件 main.ts 进行评审评论
        image: cnbcool/pr-comments:latest
        settings:
          content: "test"
          fromFile: "test.txt"
          comments:
            - path: "main.ts"
              content: "test"
      
      - name: 上传评审评论——对文件 main.ts 的某些行进行评审评论
        image: cnbcool/pr-comments:latest
        settings:
          content: "test"
          fromFile: "test.txt"
          comments:
            - path: "main.ts"
              content: "test"
              start_line: 1
              end_line: 10
```

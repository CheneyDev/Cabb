# changelog

根据两个版本之间的历史提交信息生成 changelog，默认将生成到工作目录下的 `CHANGELOG.md` 文件

注意：不支持在 `branch.delete` 事件中使用

## 参数

### tag

- type: String
- required: 否
- default: 触发流水线的 tag 名

用于生成 changelog 的 tag（获取 `previousTag` | `previousBranch` 和 `tag` 之间的变更记录来生成 changelog）

### previousTag

- type: String
- required: 否
- default: 按创建时间排序的上一个版本（tag）

用于生成 changelog 的上一个版本 tag，优先级大于 `previousBranch`。
如果有传入，则获取 `previousTag` 和 `tag` 之间的变更记录生成 changelog；
如果没传入，则获取 `tag` 和上一个 tag 之间的变更记录生成 changelog

### previousBranch

- type: String
- required: 否

用于生成 changelog 的上一个版本（分支名）
如果有传入，则获取 `previousBranch` 和 `tag` 之间的变更记录生成 changelog；
如果没传入，则获取 `tag` 和上一个 tag 之间的变更记录生成 changelog

### from

- type: String
- required: 否

用于生成写入 `target` 文件的 changelog，如果传入，则从 `from` 版本开始生成日志，否则从第一个版本开始生成日志

### target

- type: String
- required: 否
- default: CHANGELOG.md

修改生成日志的文件路径。当设置 `dryRun: true` 或 `day: xxx` 时，该字段无效。
写入文件的日志起始版本为 `from`，如果 `from` 不存在，则为第一个版本。
不受 `tag` 和 `previousTag` 影响。

### dryRun

- type: Boolean
- required: 否
- default: false

仅生成并直接返回上一个版本至今的变更日志
(如果传入了 `previousTag` 和 `tag`，则生成这两者之间的变更日志), 内容见输出结果;

如果设置该字段，则不会将日志写入文件，即 `target` 字段无效。且 `day` 无效

### day

- type: String
- required: 否

获取最近一段时间内的提交汇总, 内容见输出结果

- 1: 最近1天内，从当天 00:00 开始
- 2: 最近2天内，从昨天 00:00 开始
- 7: 最近7天
- 2022-12-02: 从具体某一天开始
- 2022-12-02 11:00:00: 从具体某一个时间点开始

如果设置该字段，则不会将日志写入文件，即 `target` 字段无效。且 `tag` 和 `previousTag` 无效

### discard

- type: Boolean
- required: 否
- default: false

是否丢弃部分提交，只保留语义化版本中有有版本变更影响的提交

### linkReferences

- type: Boolean
- required: 否
- default: true

是否显示关联的链接，包括提交和issue等。

### showIssue

- type: Boolean
- required: 否
- default: true

是否识别并显示 issue id，支持识别 `#123` 格式

### showPr

- type: Boolean
- required: 否
- default: true

是否识别并显示 pr id。支持识别 `PR-URL: #123` 格式

### showCommitTime

- type: Boolean
- required: 否
- default: true

是否显示代码提交时间

## 输出结果

```js
{
  latestChangeLog, // 最新的变更日志
}
```

## 在 Docker 上使用

```shell
docker run --rm -t -v $(pwd):$(pwd) -w $(pwd) cnbcool/changelog .
```

## 在 云原生构建 上使用

生成 `CHANGELOG.md` 文件：

```yaml
# .cnb.yml
$:
  tag_push:
    - stages:
      - name: changelog
        image: cnbcool/changelog
      - name: showfile
        script: ls -al CHANGELOG.md
```

导出为变量，并发送到群：

```yaml
# .cnb.yml
main:
  tag_push:
    - stages:
      - name: 生成 changelog
        image: cnbcool/changelog
        settings:
          dryRun: true
          linkReferences: false
        exports:
          latestChangeLog: RELEASE_NOTES
      - name: 回显 changelog
        script: echo "$RELEASE_NOTES"
      - name: 发送到企业群
        image: tencentcom/wecom-message
        settings:
          robot: your-bot-xxxxxxx
          content: |
            ## ChangeLog
            
            $RELEASE_NOTES      
```

每周五下午三点，自动生成工作周报：

```yaml
# .cnb.yml
main:
  # 每周五下午三点
  "crontab: 0 15 * * 5":
    - stages:
        - name: week changelog
          image: cnbcool/changelog
          settings:
            day: 7
            dryRun: true
            linkReferences: false
          exports:
            latestChangeLog: WEEK_REPORT
        - name: wework notice
          image: tencentcom/wecom-message
          settings:
            robot: your-bot-xxxxxxx
            content: |
              ## ChangeLog
              
              $WEEK_REPORT
```

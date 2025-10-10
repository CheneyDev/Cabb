# AI 代码评审

支持使用 AI 进行四类操作，代码评审、PR 标题和描述的可读性检测、提交注释可读性检测、变更总结

## 镜像名

cnbcool/ai-review

## 参数说明

- `type`: 操作类型，默认为 `code-review`
  - `code-review`: 代码评审
  - `pr-info-readability-check`: PR 标题和描述的可读性检测
  - `commit-message-readability-check`: 提交注释可读性检测
  - `diff-summary`: 变更总结
- `message`: `type=commit-message-readability-check` 时，可指定 `message` 参数，表示提交注释，
不指定默认取 `CNB_COMMIT_MESSAGE`
- `pr_comment`: `true` 或 `false`，默认为 `true`。是否将 pr 评论发表为 pull_request 行内评论（目前仅支持发送一条评论）。
`type=code-review` 时有效
- `max_comments`: 最大评论数，默认为 10，最多 10 条。
`type=code-review` 时有效

## 输出结果

- `pr-info-readability-check`: 输出环境变量名 `status`，值为 yes 或 no
- `commit-message-readability-check`: 输出环境变量名 `status`，值为 yes 或 no
- `diff-summary`: 输出环境变量名 `summary`，表示变更总结内容

## 在 CNB 中使用

### 代码评审

默认会将 ai 评审意见发送评论到 pr，目前仅支持发送一条评论

```yaml
main:
  pull_request:
    - stages:
      - name: 代码评审
        image: cnbcool/ai-review:latest
        settings:
          type: code-review
```

### PR 标题和描述的可读性检测

```yaml
main:
  pull_request:
    - stages:
      - name: 标题和描述的可读性检测
        image: cnbcool/ai-review:latest
        settings:
          type: pr-info-readability-check
        exports:
          status: STATUS
      - name: 标题和描述的可读性检测结果
        script: echo $STATUS
```

### 提交注释的可读性检测

```yaml
main:
  push:
    - stages:
      - name: 提交注释的可读性检测
        image: cnbcool/ai-review:latest
        settings:
          type: commit-message-readability-check
        exports:
          status: STATUS
      - name: 提交注释的可读性检测结果
        script: echo $STATUS
```

### 变更总结

```yaml
main:
  pull_request:
    - stages:
      - name: 变更总结
        image: cnbcool/ai-review:latest
        settings:
          type: diff-summary
        exports:
          summary: SUMMARY
      - name: 变更总结
        script: echo $SUMMARY
```

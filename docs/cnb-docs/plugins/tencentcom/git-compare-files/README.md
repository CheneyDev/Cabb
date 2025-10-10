# git-compare-files

获取差异文件列表，输出到文件，一行一个文件名，相对构建目录。

## 输入

### from

Type: `String`

变更前的 hash、branch name 或 tag。

### to

Type: `String`

变更后的 hash、branch name 或 tag。
注意 from 和 to 不要写反了。

### changed

Type: `String`

变更列表存放文件名，包含修改过的、新增加的文件。

### deleted

Type: `String`

删除列表存放文件名。

### added

Type: `String`

新增列表存放文件名。

## 输出

通过 `##[set-output key=value]` 输出。

### `changed`

Type: `String`

变更文件列表，`\n` 分割。

### `deleted`

Type: `String`

删除文件列表，`\n` 分割。

### `added`

Type: `String`

新增文件列表，`\n` 分割。

## 在 Docker 上使用

```shell
docker run --rm -it \
  -v $(pwd):$(pwd) -w $(pwd) \
  -e PLUGIN_FROM="HEAD~1" \
  -e PLUGIN_TO="HEAD" \
  tencentcom/git-compare-files
```

## 在 云原生构建 上使用

输出合并请求相对于目标分支的变更。

```yaml
main:
  pull_request:
  - stages:
    - name: 对比变更文件
      image: tencentcom/git-compare-files
      settings:
        from: $CNB_BRANCH_SHA
        to: $CNB_COMMIT
```

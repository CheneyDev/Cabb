# git-commit-list

获取提交记录列表，支持输出提交记录列表 json 数据到文件

## 在 云原生构建 中使用

```yaml
main:
  pull_request:
  - stages:
    - name: git-commit-list
      image: tencentcom/git-commit-list:latest
      settings:
        # 可选，文件列表输出到文件中
        toFile: commit-list.json
```

commit-list.json 数据格式

```json
[
  {
    "id": "ee0738cc39660baa47e08b6ef3b9f942165d1a0a",
    "short_id": "ee0738cc",
    "author_name": "name",
    "author_email": "name@xxx.com",
    "create_at": "2022-11-24T16:26:48+08:00",
    "message": "fix: 获取提交记录列表\n\n#123",
    "title": "fix: 获取提交记录列表"
  }
]
```

## docker

push 事件：

```shell
docker run --rm \
    -e TZ=Asia/Shanghai \
    # 事件名
    -e CNB_EVENT="push" \
    # 本次提交之前的一个 commitid
    -e CNB_BEFORE_SHA="xxx" \
    # 本次的最新一个 commitid
    -e ORANGE_COMMIT="xxx" \
    # 存放提交记录列表 json 数据的文件相对路径
    -e PLUGIN_TOFILE="commit-list.json" \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    tencentcom/git-commit-list:latest
```

mr 事件：

```shell
docker run --rm \
    -e TZ=Asia/Shanghai \
    # 事件名
    -e CNB_EVENT="merge_request" \
    # mr 源分支名
    -e CNB_PULL_REQUEST_BRANCH="xxx" \
    # mr 目标分支名
    -e CNB_BRANCH="master" \
    # 本次的最新一个 commitid
    -e CNB_COMMIT="xxx" \
    # 存放提交记录列表 json 数据的文件相对路径
    -e PLUGIN_TOFILE="commit-list.json" \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    tencentcom/git-commit-list:latest
```

# git-mr-limit

检查 MR 中 commit 次数 和 代码变更行数 。

## 输入

### maxCommitCount

Type: `Number`

单次 MR 最大 commit 次数

### maxCommitLine

Type: `Number`

单次 MR 最大代码变更行数，包含增加、删除和修改

### exclude

Type: `String | String[]`

检查变更行数时需要忽略的文件或文件夹（相对路径，如：test）

### include

Type: `String | String[]`

检查变更行数时需要包括的文件或文件夹（相对路径，如：src）

## 在 Docker 上使用

```shell
docker run --rm -it \
  -v $(pwd):$(pwd) -w $(pwd) \
  # 用户传入
  -e PLUGIN_MAXCOMMITCOUNT="10" \
  -e PLUGIN_MAXCOMMITLINE="200" \
  -e PLUGIN_EXCLUDE="test" \
  -e PLUGIN_INCLUDE="src" \
  # CI 环境变量中获取
  # 事件名
  -e CNB_EVENT="pull_request" \
  # 目标分支最新 commitID
  -e CNB_PULL_REQUEST_TARGET_SHA="" \
  # 源分支最新 commitID
  -e CNB_PULL_REQUEST_SHA="" \
  # 源分支名
  -e CNB_PULL_REQUEST_BRANCH="" \
  # 目标分支名
  -e CNB_BRANCH="" \

  tencentcom/git-mr-limit
```

## 在 云原生构建 上使用

```yaml
main:
  pull_request:
  - stages:
    - name: mr 变更检查
      image: tencentcom/git-mr-limit
      settings:
        maxCommitCount: 10
        maxCommitLine: 200
        exclude: 
          - test
          - package-lock.json
        include: src
```

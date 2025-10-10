# 附件插件

支持 commit 和 release 上传或下载附件，目前暂支持 5GB 以内的文件上传。

## 镜像

`cnbcool/attachments:latest`

## 支持的事件

以下事件支持 `上传`、`下载` 附件操作：

- push
- branch.create
- branch.delete
- pull_request
- pull_request.target
- pull_request.approved
- pull_request.changes_requested
- pull_request.mergeable
- pull_request.merged
- tag_push
- vscode
- auto_tag
- tag_deploy.*

## 参数说明

- `type`: 操作类型，默认值为`UPLOAD`，非必填。
  - `UPLOAD`：上传附件
  - `DOWNLOAD`：下载附件
- `tag`: tag 名称，操作指定 tag 时填写，非必填。`tag` 需有对应的 release 才能上传附件。
- `commit`: commit sha，长 `hash`，非必填。
- `attachments`: 文件名，必填。
- `slug`: 仓库路径，非必填。默认为触发当前流水线的仓库，非跨仓上传无需传入。如需跨仓库上传附件，可使用此参数指定仓库路径。
- `endpoint`: 上传或下载附件的地址，非必填。默认为当前平台的 OPENAPI 的地址，非跨平台无需传入。
如需跨平台上传附件，如需要从 CNB PAAS 版 上传到 CNB SAAS 版，可使用此参数指定 `endpoint`，如 `https://api.cnb.cool`。
- `token`: 上传或下载附件的 token，非必填。默认为当前用户的临时票据。如需跨平台上传，或跨仓库上传时当前用户无仓库权限，可传入该参数。

### attachments

当`type`为`UPLOAD`时，须传入上传文件的相对路径；
多个文件时用逗号分隔或使用数组，支持 `glob` 表达式匹配。

当`type`为`DOWNLOAD`时，传入下载文件的文件名。
传路径时，代码中会自动获取文件名，默认下载到工作目录。

### tag/commit

其中 `tag` 或 `commit` 都不传时，操作目标默认为：

`tag_push`和`tag_deploy.*` 事件，`tag` 取环境变量 `$CNB_BRANCH`，
操作对象为 `tag` 对应的 `release` 的附件。

其他事件，`commit` 取环境变量 `$CNB_COMMIT`，操作对象为 `commit` 的附件。

其中 `tag` 和 `commit` 都传时，操作目标默认为 `tag` 对应的 `release`（没有 release 不能上传）。

## 在 云原生构建 中使用

### release 附件

```yaml
$:
  tag_push:
    - stages:
      - name: release 上传附件
        image: cnbcool/attachments:latest
        settings:
          attachments:
            - "./*.txt"
      
      - name: 支持否定模式，下面示例会排除 ./test1.txt 文件
        image: cnbcool/attachments:latest
        settings:
          attachments:
            - "./*.txt"
            - "!./test1.txt"

      - name: release 下载附件
        image: cnbcool/attachments:latest
        settings:
          type: DOWNLOAD
          attachments: "./test1.txt,./test2.txt"
```

### commit 附件

```yaml
main:
  push:
    - stages:
      - name: commit 上传附件，多个文件数组写法
        image: cnbcool/attachments:latest
        settings:
          attachments:
            - "./test1.txt"
            - "./test2.txt"

      - name: commit 上传附件，多个文件逗号分隔
        image: cnbcool/attachments:latest
        settings:
          attachments: "./test1.txt,./test2.txt"
      
      - name: commit 下载附件
        image: cnbcool/attachments:latest
        settings:
          type: DOWNLOAD
          attachments: "test1.txt,test2.txt"
```

### 跨仓库上传附件

使用场景举例：某私有仓库构建产物需要上传到公开仓库（构建产物想公开，但仓库不想公开，可创建公开仓库，专门用于公开构建产物）
，可如下配置可实现跨仓上传附件：

```yaml
$:
  tag_push:
    - stages:
      - name: 跨仓库上传 release 附件
        image: cnbcool/attachments:latest
        settings:
          attachments:
            - "./*.txt"
          # 上传附件到指定仓库
          slug: groupname/reponame
      
      - name: 跨仓库下载 release 附件
        image: cnbcool/attachments:latest
        settings:
          type: DOWNLOAD
          attachments: "./test1.txt,./test2.txt"
          # 从指定仓库下载附件
          slug: groupname/reponame
```

### 跨平台上传附件

只支持 CNB 不同平台之间上传，例如可从 CNB PAAS 版上传到 CNB SAAS 版。

```yaml
$:
  tag_push:
    - stages:
      - name: 跨平台跨仓库上传 release 附件
        image: cnbcool/attachments:latest
        settings:
          attachments:
            - "./*.txt"
          # 跨平台上传附件到指定仓库
          slug: groupname/reponame
          endpoint: https://api.cnb.cool
          # 跨平台一定需要传入 token
          token: xxxxx
      
      - name: 跨平台跨仓库下载 release 附件
        image: cnbcool/attachments:latest
        settings:
          type: DOWNLOAD
          attachments: "./test1.txt,./test2.txt"
          # 跨平台从指定仓库下载附件
          slug: groupname/reponame
          endpoint: https://api.cnb.cool
          # 跨平台一定需要传入 token
          token: xxxxx
```

## 使用 Docker 镜像直接上传/下载附件

参数说明：

- `PLUGIN_ATTACHMENTS` 附件列表，多个文件逗号分隔
- `PLUGIN_TYPE` 上传/下载，`UPLOAD` 或 `DOWNLOAD`，默认为 `UPLOAD`
- `PLUGIN_COMMIT` commit id，如需上传/下载 commit 附件需要传此参数
- `PLUGIN_TAG` tag 名称，如需上传/下载 release 附件需要传此参数
- `CNB_TOKEN` CNB API Token，在 `个人设置 -> 访问令牌` 中创建并获取 token，注意选中 `repo-contents` 读写权限
- `CNB_API_ENDPOINT` CNB API Endpoint，CNB SAAS 版为 `https://api.cnb.cool`
- `CNB_REPO_SLUG` 仓库路径，例如 `groupname/reponame`

### 上传/下载 commit 附件

```shell
# commit 上传附件
docker run --rm \
    -e TZ=Asia/Shanghai \
    -e CNB_TOKEN='xxxx' \
    -e CNB_API_ENDPOINT='https://api.cnb.cool' \
    -e CNB_REPO_SLUG='groupname/reponame' \
    -e PLUGIN_COMMIT='xxx' \
    -e PLUGIN_ATTACHMENTS='./xxx.png' \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    cnbcool/attachments:latest

# commit 下载附件
docker run --rm \
    -e TZ=Asia/Shanghai \
    -e CNB_TOKEN='xxxx' \
    -e CNB_API_ENDPOINT='https://api.cnb.cool' \
    -e CNB_REPO_SLUG='groupname/reponame' \
    -e PLUGIN_COMMIT='xxx' \
    -e PLUGIN_ATTACHMENTS='xxx.png' \
    -e PLUGIN_TYPE='DOWNLOAD' \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    cnbcool/attachments:latest
```

### 上传/下载 release 附件

需要已经存在 tag，并且有创建过 release

```shell
# release 上传附件
docker run --rm \
    -e TZ=Asia/Shanghai \
    -e CNB_TOKEN='xxxx' \
    -e CNB_API_ENDPOINT='https://api.cnb.cool' \
    -e CNB_REPO_SLUG='groupname/reponame' \
    -e PLUGIN_TAG='v1.0.0' \
    -e PLUGIN_ATTACHMENTS='./xxx.png' \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    cnbcool/attachments:latest

# release 下载附件
docker run --rm \
    -e TZ=Asia/Shanghai \
    -e CNB_TOKEN='xxxx' \
    -e CNB_API_ENDPOINT='https://api.cnb.cool' \
    -e CNB_REPO_SLUG='groupname/reponame' \
    -e PLUGIN_TAG='v1.0.0' \
    -e PLUGIN_ATTACHMENTS='xxx.png' \
    -e PLUGIN_TYPE='DOWNLOAD' \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    cnbcool/attachments:latest
```

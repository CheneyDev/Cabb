# 介绍

提交本地代码

## 参数

### commitMessage

- type: string
- 必填: 否
- 默认值: chore: committed by ci

提交时的 `Commit Message`

### branchName

- type: string
- 必填: 否
- 默认值: commit-from-ci/${sha前八位}-${三位随机字符串}

临时分支名

为了防止自动化提交过于随意而导致的风险，会先创建临时分支，提交变动，再推送到 `origin` ，而不是直接在对应分支上提交。

有些项目可能会对分支名有限制，这里可以指定分支名。

### pushCurrent

- type: boolean
- 必填: 否
- 默认值: false

是否送到到当前分支，为 `true` 则忽略 `branchName`

如果在 `push` 流水线中，可能会引起流水线死循环

***慎重选择***

### addPatterns

- type: string | string[]
- 必填: 否
- 默认值: **
  
指定 `git add` 的目录、文件列表，通过 [glob](https://globster.xyz/)
模式匹配 `git status --porcelain` 得到的变更文件。

注意：新增了一个文件夹，如 `dir1/dir2`，不管里面多少文件，变更文件是 `dir1/dir2/`，
可用 `dir1/dir2/**` 匹配其文件夹或文件夹下所有文件。

## 导出

### branch

若 `pushCurrent` 不为 `true`，则导出 `branch` 为 新建分支名。

## 在 云原生构建 上使用

```yaml
test:
  push:
    - stages:
        - name: set token
          imports: https://xxx.com/token.yml # 从私有仓库引入 TOKEN
          image: tencentcom/git-set-credential:latest
          settings:
            userName: ${USER_NAME}
            userEmail: ${USER_EMAIL}
            loginUserName: ${LOGIN_USER_NAME}
            loginPassword: ${LOGIN_PASSWORD}
        - name: update file
          script: echo test > test.txt
        - name: commit
          image: tencentcom/git-commit:latest
          settings:
            add:
              - test.txt
            commitMessage: "commit by ci"
          exports:
            branch: COMMIT_BRANCH
        - name: show branch
          script: echo $COMMIT_BRANCH
```

## 在 Docker 上使用

```shell
docker run --rm \
    -e TZ=Asia/Shanghai \
    -e PLUGIN_COMMITMESSAGE="test: committed by ci" \
    -e PLUGIN_ADDPATTERNS="test.txt" \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    tencentcom/git-commit:latest
```

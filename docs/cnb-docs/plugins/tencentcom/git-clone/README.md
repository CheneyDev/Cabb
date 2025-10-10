# git clone

## 介绍

将其他仓库clone到本地

## 原理

```shell
git clone https://{user}:{password}@{url}
```

## 参数

### user_name

- type: string
- 必填: 否

如设置，插件会执行 `git config user.name`

### user_email

- type: string
- 必填: 否

如设置，插件会执行 `git config user.email`

### dist

- type: string
- 必填: 否

将仓库 clone 到何处

### git_url

- type: string
- 必填: 是

仓库地址, 如 `https://xxx.com/xxx/xxx.git`，不支持 ssh 地址

### git_user

- type: string
- 必填: 是

登陆用的用户名

### git_password

- type: string
- 必填: 是

登陆用的密码

## 在 云原生构建 上使用

```yaml
main:
  push:
    - stages:
        - name: git clone
          image: tencentcom/git-clone
          # 引用环境变量
          imports: https://xxx.com/git-clone.yml
          settings:
            user_name: xxx
            user_email: xxx
            git_url: https://xxx/xxx/demo1.git
            git_user: $GIT_USER
            git_password: $GIT_PASSWORD
            dist: other_repos
        - name: git commit
          script:
            - cd other_repos/demo1
            - echo 1 > tmp.txt
            - git add .
            - git commit -m "test"
            - git push
```

在 `云原生构建` 可用 `CNB_TOKEN_USER_NAME` 和 `CNB_TOKEN` 作
为 `git_user` 和 `git_password` 去 clone `云原生构建` 的仓库。

```yaml
# .cnb.yml
main:
  push:
    - stages:
        - name: git clone
          image: tencentcom/git-clone
          settings:
            user_name: xxx
            user_email: xxx
            git_url: https://xxx/xxx/demo1.git
            git_user: $CNB_TOKEN_USER_NAME
            git_password: $CNB_TOKEN
            dist: other_repos
```

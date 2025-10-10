# 设置 git 账号信息

## 介绍

设置 `git` 的 `user.name`、 `user.email`及登录所需用户密码等信息，供 `git` 提交使用。

## 原理

```shell
git config user.name xxx
git config user.email
echo "https://xxx:xxx@xxx.com" >> .git/.git-credentials
git config credential.helper "store --file=$(pwd)/.git/.git-credentials"
```

## 参数

### userName

- type: string
- 必填: 是

user.name

### userEmail

- type: string
- 必填: 是

user.email

### loginUserName

- type: string
- 必填: 是

登陆用的用户名

### loginPassword

- type: string
- 必填: 是

登陆用的密码

## 在 Docker 上使用

```shell
docker run --rm \
    -e TZ=Asia/Shanghai \
    -e PLUGIN_USERNAME="xxx" \
    -e PLUGIN_USEREMAIL="xxx@xxx.com" \
    -e PLUGIN_LOGINUSERNAME="xxx" \
    -e PLUGIN_LOGINPASSWORD="xxx" \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    tencentcom/git-set-credential:latest
```

## 在 云原生构建 上使用

```yaml
# .cnb.yml
main:
  push:
    - stages:
        - name: set token
          image: tencentcom/git-set-credential:latest
          settings:
            userName: ${CNB_BUILD_USER}
            userEmail: ${CNB_BUILD_USER_EMAIL}
            loginUserName: ${CNB_TOKEN_USER_NAME}
            loginPassword: ${CNB_TOKEN}
```

```yaml
# 从私有仓库引入 TOKEN
main:
  push:
    - stages:
        - name: set token
          imports: https://xxx.com/token.yml
          image: tencentcom/git-set-credential:latest
          settings:
            userName: ${USER_NAME}
            userEmail: ${USER_EMAIL}
            loginUserName: ${LOGIN_USER_NAME}
            loginPassword: ${LOGIN_PASSWORD}
```

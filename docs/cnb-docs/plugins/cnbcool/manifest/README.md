# Manifest

推送 Docker manifest 多系统架构镜像文件的插件。

![badge](https://cnb.cool/cnb/plugins/cnbcool/manifest/-/badge/git/latest/ci/pipeline-as-code)
![badge](https://cnb.cool/cnb/plugins/cnbcool/manifest/-/badge/git/latest/ci/git-clone-yyds)
![badge](https://cnb.cool/cnb/plugins/cnbcool/manifest/-/badge/git/latest/code/vscode-started)
![badge](https://cnb.cool/cnb/plugins/cnbcool/manifest/-/badge/git/latest/ci/status/pull_request)
![badge](https://cnb.cool/cnb/plugins/cnbcool/manifest/-/badge/git/latest/ci/status/auto_merge)

## 镜像

`cnbcool/manifest:latest`

## 在云原生构建中使用

```yaml
main:
  push:
  # 示例1：镜像源为 `hub.docker.com`
  - stages:
    - name: manifest
      image: cnbcool/manifest
      settings:
        # docker 用户名
        username: docker-usernname
        # docker 密码
        password: docker-password
        # docker 镜像仓库地址
        target: foo/bar:v1.0.0
        # 模板，用于生成多系统架构镜像名称
        template: foo/bar:v1.0.0-OS-ARCH
        # 多系统架构
        platforms:
          - linux/amd64
          - linux/arm64
        # 是否跳过 TLS 证书验证，默认 false
        skipVerify: false
        # 忽略缺失的源镜像，即如果有源镜像缺失不报错，默认false
        ignoreMissing: false

  # 示例二：镜像源为 `docker.cnb.cool`
  # 无需传用户名密码，默认使用环境变量中的用户名和密码
  - stages:
    - name: manifest
      image: cnbcool/manifest
      settings:
        # docker 镜像仓库地址
        target: docker.cnb.cool/foo/bar:v1.0.0
        # 模板，用于生成多系统架构镜像名称
        template: docker.cnb.cool/foo/bar:v1.0.0-OS-ARCH
        # 多系统架构
        platforms:
          - linux/amd64
          - linux/arm64
        # 是否跳过 TLS 证书验证，默认 false
        skipVerify: false
        # 忽略缺失的源镜像，即如果有源镜像缺失不报错，默认false
        ignoreMissing: false
```

通过上面的配置，可以在 DockerHub 给 foo/bar 打上 v1.0.0 标签，
并且 v1.0.0 跟 DockerHub 上这两个标签的镜像关联，
注意，DockerHub 上以下两个镜像标签必须存在，即已经推送过：

- foo/bar:v1.0.0-linux-amd64
- foo/bar:v1.0.0-linux-arm64

执行 docker pull foo/bar:v1.0.0 时会基于运行环境的系统架构，拉取对应的镜像。

## 参数说明

- `username`: 非必填，docker 认证用户名，当 `target` 镜像源为 CNB 制品库时，
默认使用环境变量 [`CNB_TOKEN_USER_NAME`](/zh/env.html#CNB_TOKEN_USER_NAME)。
- `password`: 非必填，docker 认证密码，当 `target` 镜像源为 CNB 制品库时，
默认使用环境变量 [`CNB_TOKEN`](/zh/env.html#CNB_TOKEN)。
- `skipVerify`: 非必填，`true` 或 `false`，默认 `false`，是否跳过 TLS 证书验证。
- `platforms`: 必填，多系统架构列表，多个用分号分隔或以上述示例中的数组形式传入。
格式为 `OS/ARCH`(系统/架构) 的系统架构列表，它会基于 template 模板生成 target 要关联的镜像列表。
例如：`linux/amd64,linux/arm64`
- `target`: 必填，目标镜像名称（合并后的 manifest 会推送到这个镜像名）。
`imagename:tag`（`imagename` 替换为镜像名，`tag` 替换为标签名）格式。
例如：`foo/bar:v1.0.0`
- template : 必填，关联镜像的名称模板，模板中的 OS 和 ARCH 字符串会被替换成具体系统和架构（platforms 中传入的系统/架构）。
`imagename:tag-OS-ARCH`（`imagename` 替换为镜像名，`tag` 替换为标签名）格式
例如：`foo/bar:v1.0.0-OS-ARCH`
- ignoreMissing : 非必填，`true` 或 `false`，默认 `false`，是否忽略缺失的源镜像，即需要合并的多系统架构的镜像如果缺失不报错，
上述示例中 `foo/bar:v1.0.0-linux-arm64` 和 `foo/bar:v1.0.0-linux-amd64` 两个镜像的其中任何一个缺失，不会报错。
但如果都缺失仍然会报错。

---
title: Helm 制品库
permalink: https://docs.cnb.cool/zh/artifact/helm.html
summary: 该文档介绍了CNB Helm制品库的使用方法。包括使用CNB访问令牌登录（命令：`helm registry login helm.cnb.cool -u cnb -p {token-value}`），制品路径规则（有同名和非同名两种），推送制品的多种方式（本地推送及云原生构建、开发中推送），使用制品的常见命令（如拉取、查看信息、预览manifest、安装和升级等） 。
---

## 登录 CNB Helm 制品库

CNB Helm 仅支持 oci 格式的 Helm 制品，建议使用 Helm v3.8.0 以上版本，详情请查阅 Helm 文档

您可以使用 CNB 的访问令牌作为登录凭据，登录命令：

```bash
helm registry login helm.cnb.cool -u cnb -p {token-value}
```

## Helm 制品路径规则

helm 制品在发布到某一仓库时，支持两种命名规则

1. 同名制品 - 制品路径与仓库路径一致，如：`helm.cnb.cool/{repository-path}`
2. 非同名制品 - 以仓库路径作为制品的命名空间，制品路径 = 仓库路径/制品名称，如：`helm.cnb.cool/{repository-path}/{artifact-name}`

注意：helm push remote 时，chart name 不出现在 remote-url 中，而是从 chart 中读取。

所以，同名制品的 remote-url 为 `helm.cnb.cool/{group-path}`，非同名制品的 remote-url 为 `helm.cnb.cool/{repository-path}`。

## 推送制品

### 本地推送

同名制品

```bash
# chart name 需与仓库同名
helm package chart-path
helm push chartname-version.tgz oci://helm.cnb.cool/{group-path}
```

非同名制品

```bash
helm package chart-path
helm push chartname-version.tgz oci://helm.cnb.cool/{repository-path}
```

### 云原生构建中推送

```yaml
main:
  push:
    - docker:
        image: alpine/helm
      stages:
        - name: helm login
          script: helm registry login -u ${CNB_TOKEN_USER_NAME} -p "${CNB_TOKEN}" ${CNB_HELM_REGISTRY}
        - name: helm package
          script: helm package ${YOUR_CHART_PATH}
        - name: helm push
          script: helm push oci://${CNB_HELM_REGISTRY}/${CNB_GROUP_SLUG_LOWERCASE}
```

### 云原生开发中推送

同名制品

```bash
# chart name 需与仓库同名
helm package ${YOUR_CHART_PATH}
helm push chartname-version.tgz oci://helm.cnb.cool/{grou-path}
```

非同名制品

```bash
helm package ${YOUR_CHART_PATH}
helm push chartname-version.tgz oci://helm.cnb.cool/{repository-path}
```

## 使用制品

### 本地命令行拉取

```bash
helm pull oci://helm.cnb.cool/{artifact-path} --version {version}

# ...
```

### 其他常用命令

#### 查看 Helm 信息

```bash
helm show all oci://helm.cnb.cool/{artifact-path} --version {version}
```

#### 预览 manifest

```bash
helm template {my-release} oci://helm.cnb.cool/{artifact-path} --version {version}
```

#### 安装 helm chart

```bash
helm install {my-release} oci://helm.cnb.cool/{artifact-path} --version {version}
```

#### 升级 helm chart

```bash
helm upgrade {my-release} oci://helm.cnb.cool/{artifact-path} --version {version}
```

## 更多用法

更多 Helm 用法，请查阅官方文档

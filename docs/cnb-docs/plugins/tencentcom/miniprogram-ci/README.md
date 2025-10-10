# 微信小程序

从微信开发者工具中抽取小程序/小游戏项目代码的编译模块，实现一键**构建 npm /预览/上传小程序**。

可运行在 **云原生构建** 、**GitHub Actions** 。

## 注意事项

使用前需要使用小程序管理员身份访问 "微信公众平台-开发-开发设置"后下载代码上传密钥，
并关闭 IP 白名单，才能使用此插件进行上传、预览操作。

**注意：该插件使用的npm库[miniprogram-ci](https://www.npmjs.com/package/miniprogram-ci)**
**从`1.8.0`升级至`2.0.6`，可使用 `tencentcom/miniprogram-ci:v1.8.0` 来使用旧版的 npm 库**

## 在 云原生构建 上使用

```yaml
main:
  push:
  - stages:
    - name: miniprogram-ci
      image: tencentcom/miniprogram-ci
      settings:
        appid: wxsomeappid
        # example  ./dist
        projectPath: the/project/path           
        # example  ./privateKeyPath: private.wx15949dc6d9b035d2.key
        privateKeyPath: the/path/to/privatekey  
        version: v1.0.0
```

## 在 Github-action 上使用

```yaml
name: CI

on:
  push:
    branches:
    - master

jobs:
  scf-deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout 
        uses: actions/checkout@master
      - name: Build 
        uses: docker://tencentcom/miniprogram-ci:latest
        env: 
          PLUGIN_APPID: wxsomeappid
          PLUGIN_PROJECTPATH: the/project/path
          PLUGIN_PRIVATEKEYPATH: the/path/to/privatekey
          PLUGIN_VERSION: v1.0.0
```

## 证书获取

1. 登录微信公众平台，进入到 开发-开发者工具，下载证书.
2. 如果遇到 ip 白名单问题，可以选择关闭白名单限制

## 更多信息以及配置项

更多信息，请查阅：[package/miniprogram-ci](https://www.npmjs.com/package/miniprogram-ci)

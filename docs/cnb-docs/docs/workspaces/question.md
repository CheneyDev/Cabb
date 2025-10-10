---
title: 常见问题解答
permalink: https://docs.cnb.cool/zh/workspaces/question.html
summary: 文本是关于 code-server 安装的常见问题解答，介绍了安装最新版 code-server 的命令，即通过“curl -fsSL https://code-server.dev/install.sh | sh”安装。同时，针对社区版 code-server 某些版本有 bug 导致 WebIde 无法使用的状况，给出了安装指定版本（如 4.100.3 ，命令为“curl -fsSL https://code-server.dev/install.sh | sh -s -- --version 4.100.3” ）来临时解决问题的方法 ，问题解决后可再安装最新版 。
---

## 如何安装指定版本 code-server

自定义开发环境一般会安装 code-server 来支持 WebIde，以下方式会安装最新版本 code-server：

```shell
curl -fsSL https://code-server.dev/install.sh | sh
```

由于 code-server 来自于社区，有时候某些版本可能会有 bug 导致 WebIde 无法正常使用，
此时可安装指定版本来临时解决问题。等问题解决后再安装最新版

使用如下命令安装指定版本：

```shell
# 安装指定版本，例如 4.100.3
curl -fsSL https://code-server.dev/install.sh | sh -s -- --version 4.100.3
```

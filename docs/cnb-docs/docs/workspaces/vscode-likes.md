---
title: VSCode/Cursor 客户端
permalink: https://docs.cnb.cool/zh/workspaces/vscode-likes.html
summary: 云原生开发默认支持 VSCode 和 Cursor 客户端进行远程开发。若使用自定义开发环境，需安装 `openssh-server` 并配置相应的 SSH 服务，同时在客户端中安装 Remote-SSH 插件。此外，为避免窗口覆盖问题，应将 `Open Folders In New Window` 设置为 `on`。
---

云原生开发默认环境支持 VSCode 客户端、Cursor 客户端远程开发。

但开发者如果选择自定义开发环境，需做以下配置来支持 VSCode 客户端 和 Cursor 客户端远程开发：

### 如何支持 VSCode/Cursor 客户端

- **开发环境中安装 `openssh-server`**

自定义开发环境时可在 Dockerfile 文件中提前安装好 ssh 服务。

```bash
# .ide/Dockerfile
FROM your-image

# 注意：基础镜像不同，安装方式可能存在差异，可根据实际情况采用不同安装方式
apt-get update
apt-get install -y openssh-server
```

- **下载 VSCode/Cursor 客户端，并安装 Remote-SSH 插件**

### 解决 VSCode/Cursor 窗口覆盖问题

点击打开 VSCode 客户端按钮，可能会出现覆盖原有窗口问题。

可修改 VSCode/Cursor 设置解决问题：将 `Open Folders In New Window` 设置为 `on` 来实现每次打开新窗口。

设置路径如下：

- Manage -> Settings -> User -> Window -> New Window -> Open Folders in New Window
- 管理 -> 设置 -> 用户 -> 窗口 -> 新建窗口 -> Open Folders in New Window

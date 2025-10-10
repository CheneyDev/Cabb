---
title: 业务端口预览
permalink: https://docs.cnb.cool/zh/workspaces/business-preview.html
summary: 运行在云原生开发环境的业务可通过 `WebIDE` 或 `VSCode`/`Cursor` 客户端预览，方法一是在 `WebIDE` 控制台 `PORTS` 面板增加端口映射获取可访问 url，也可从环境变量获取，且服务需启动在 `0.0.0.0` ，方法二是通过 `VSCode`/`Cursor` 客户端 `port forward` 转发端口到本地 。
---

运行在云原生开发环境的业务，可以通过 `WebIDE` 或 `VSCode`/`Cursor` 客户端访问业务端口，实现预览。

方法一：

使用 `WebIDE` 时，可在 `WEBIDE` 的控制台的 `PORTS` 面板中增加端口映射，会自动出现业务端口的可访问 url。

业务端口访问 url 可通过如下两种方式获取：

- `WebIDE` 控制台的 `PORTS` 面板获取
- 从环境变量获取：`CNB_VSCODE_PROXY_URI`。例如 `https://fjisdofi21-{{port}}.cnb.run`，需将 `{{port}}` 替换为实际端口

:::tip
注意，服务需启动在 `0.0.0.0` 上才能使用该方法访问。启动在 `localhost` 或 `127.0.0.1` 上的服务，无法使用方法一访问。
:::

方法二：

可通过 `VSCode`/`Cursor` 客户端的 `port forward` 端口转发能力转发需要访问的业务端口到本地

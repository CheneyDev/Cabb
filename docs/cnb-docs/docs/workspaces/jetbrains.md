---
title: Jetbrains 客户端
permalink: https://docs.cnb.cool/zh/workspaces/jetbrains.html
summary: 通过JetBrains Gateway可连接并访问云原生开发环境，需在Mac/Windows下载安装Gateway，在Dockerfile中安装openssh-server和所需IDE版本。环境创建成功后，可从右上角头像下拉菜单进入已创建开发环境，若有安装JetBrains IDE则会有“JetBrains”按钮，点击可打开JetBrains Gateway并打开对应IDE 。
redirectFrom: /zh/vscode/jetbrains.html
---

进行以下配置后，可通过 JetBrains Gateway
（IDEA、GoLand、PhpStorm、PyCharm、RubyMine、WebStorm、Rider、CLion、RustRover）
连接并访问云原生开发环境。

## 准备工作

如果需要使用 Jetbrains 客户端访问远程开发环境，需要做如下准备：

### Mac/windows 下载并安装 JetBrains Gateway

JetBrains Gateway（帮助连接云原生开发环境），[下载地址](https://www.jetbrains.com/remote-development/gateway/)

### Dockerfile 中安装 `openssh-server` 和需要的 ide 版本

- openssh-server: JetBrains Gateway 需要远程开发环境支持 ssh 服务，因此需要安装 `openssh-server`,
- ide: 远程开发环境中需要安装 IDE server（提前在 Dockerfile 中安装好可以节省打开时间）

对于 ide 可安装下方 `.ide/Dockerfile` 中的版本，也可自行获取 IDE 下载路径，获取方式如下：

- [打开 JetBrains 产品页](https://www.jetbrains.com/products/)
- 找到需要的 IDE，点击 download，进入下载详情页
- 切换为 linux 版本，点击下载，会下载并打开提示页面，找到 `direct link`，右键复制链接地址即可得到下载地址

```shell{6-48}
# .ide/Dockerfile
FROM node:22

WORKDIR /root

# 安装 ssh 服务
RUN apt-get update && apt-get install -y wget unzip openssh-server

# 创建 /ide_cnb 目录，用于安装 IDE，注意安装路径必须是这个，便于自动识别环境中支持哪些 ide
RUN mkdir -p /ide_cnb

# 选择安装下方其中一个或多个 IDE

# 安装 GoLand
RUN wget https://download.jetbrains.com/go/goland-2024.3.3.tar.gz
RUN tar -zxvf goland-2024.3.3.tar.gz -C /ide_cnb

# 安装 IntelliJ IDEA
RUN wget https://download.jetbrains.com/idea/ideaIU-2024.3.5.tar.gz
RUN tar -zxvf ideaIU-2024.3.5.tar.gz -C /ide_cnb

# 安装 PhpStorm
RUN wget https://download.jetbrains.com/webide/PhpStorm-2024.3.3.tar.gz
RUN tar -zxvf PhpStorm-2024.3.3.tar.gz -C /ide_cnb

# 安装 PyCharm
RUN wget https://download.jetbrains.com/python/pycharm-professional-2024.3.5.tar.gz
RUN tar -zxvf pycharm-professional-2024.3.5.tar.gz -C /ide_cnb

# 安装 RubyMine
RUN wget https://download.jetbrains.com/ruby/RubyMine-2024.3.3.tar.gz
RUN tar -zxvf RubyMine-2024.3.3.tar.gz -C /ide_cnb

# 安装 WebStorm
RUN wget https://download.jetbrains.com/webstorm/WebStorm-2024.3.3.tar.gz
RUN tar -zxvf WebStorm-2024.3.3.tar.gz -C /ide_cnb

# 安装 CLion
RUN wget https://download.jetbrains.com/cpp/CLion-2024.3.3.tar.gz
RUN tar -zxvf CLion-2024.3.3.tar.gz -C /ide_cnb

# 安装 RustRover
RUN wget https://download.jetbrains.com/rustrover/RustRover-2024.3.5.tar.gz
RUN tar -zxvf RustRover-2024.3.5.tar.gz -C /ide_cnb

# 安装 Rider
RUN wget https://download.jetbrains.com/rider/JetBrains.Rider-2024.3.5.tar.gz
RUN tar -zxvf JetBrains.Rider-2024.3.5.tar.gz -C /ide_cnb

# 安装 code-server(VSCode WebIDE 支持)
RUN curl -fsSL https://code-server.dev/install.sh | sh \
  && code-server --install-extension cnbcool.cnb-welcome \
  && code-server --install-extension redhat.vscode-yaml \
  && code-server --install-extension orta.vscode-jest \
  && code-server --install-extension dbaeumer.vscode-eslint \
  && code-server --install-extension waderyan.gitblame \
  && code-server --install-extension mhutchie.git-graph \
  && code-server --install-extension donjayamanne.githistory

ENV LANG C.UTF-8
```

## 如何访问

点击 `启动云原生开发` 按钮，环境创建成功后，有如下入口可以进入 JetBrains IDE:

- 右上角头像下拉菜单 -> 我的云原生开发列表页 -> 已创建的开发环境中如果安装了 JetBrains IDE，会显示 `JetBrains` 按钮，点击可打开 `JetBrains Gateway`，点击链接即可打开对应的 IDE

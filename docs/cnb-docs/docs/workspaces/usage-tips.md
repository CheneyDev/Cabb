---
title: 使用技巧
permalink: https://docs.cnb.cool/zh/workspaces/usage-tips.html
summary: 该文档介绍了WebIDE的使用技巧，包括避免浏览器剪贴板授权弹出框的方法、解决快捷键冲突的方式，以及VSCode配置文件的路径和漫游策略，还说明了在不同场景下安装VSCode扩展的具体方法 。
redirectFrom: /zh/vscode/use.html
---
## 快捷键

### 浏览器剪贴板授权

浏览器中使用 `ctrl/command + v` 快捷键粘贴，可能会弹出剪贴板授权框，
可使用以下方法避免弹出授权框：

- 使用以下快捷键粘贴，不会弹出 Chrome 的剪贴板授权框。

  - win: shift + insert
  - mac: shift + command + v

- 浏览器设置中添加剪贴板授权

复制 `chrome://settings/content/clipboard` 到 chrome 浏览器中打开剪贴板设置，
添加网站域名到 `允许查看您的剪贴板` 列表中。

具体打开路径：设置 -> 隐私设置和安全 -> 权限（更多权限） -> 剪贴板 -> 允许查看您的剪贴板

### 解决快捷键冲突

当 WebIDE 快捷键跟浏览器快捷键冲突，可以使用本地客户端远程开发

## VSCode 配置文件漫游

### 配置文件路径

**WebIDE 的配置文件（`settings.json`）路径：**

- 用户维度：`~/.local/share/code-server/User/settings.json`。配置将根据用户维度进行漫游，实现保存个人配置
- 机器纬度：`~/.local/share/code-server/Machine/settings.json`。可以在 .ide/Dockerfile 中修改，来实现配置共享
- 仓库维度：`.vscode/settings.json`。相对目录为当前工作区（默认为 `/workspace`），即仓库根目录

**VSCode 客户端（Remote-SSH） 的配置文件（`settings.json`）路径：**

- 用户维度：存在用户本地，可自行配置
- 机器纬度：`~/.vscode-server/data/Machine/settings.json`。可以在 .ide/Dockerfile 中修改，来实现配置共享
- 仓库维度：`.vscode/settings.json`。相对目录为当前工作区（默认为 `/workspace`），即仓库根目录

### 更多共享策略

- **用户维度共享：可实现同一个用户的配置共享，不同用户隔离**

对于 WebIDE：用户纬度的配置文件（`settings.json`）会被自动漫游，且只对用户自己生效。

注意：对于客户端，由于用户纬度配置文件存在本地，无法被漫游，如需修改配置，需自行修改本地配置文件。

- **仓库维度共享：可实现同一个仓库的配置共享，不同仓库隔离**

如何在同一个仓库内共享配置：提交 `.vscode/settings.json` 文件到仓库即可

- **环境维度共享：可实现跨仓库的环境共享**

假设你已经有构建好的开发环境镜像，那么可以直接使用在 `.cnb.yml` 配置中：

```yaml{9}
# .cnb.yml
$:
  vscode:
    - name: vscode
      services:
        - vscode
      docker:
        # 使用自定义镜像作为开发环境
        image: groupname/imagename:latest
```

## 如何安装 VSCode 扩展

自定义开发环境时可在 Dockerfile 中安装常用的 VSCode 扩展，方便复用，有以下两种安装方式。

说明：WebIDE 使用的扩展源是 `open-vsx` (开源)，
非 `微软官方扩展源`（闭源），
如果 `open-vsx` 缺失某些扩展，可在 `微软官方扩展源` 搜索扩展，
在扩展详情页下载 vsx 源文件安装。

注意：当 `open-vsx` 能搜索到需要安装的扩展时，可以通过扩展 id 安装；
当 `open-vsx` 缺失某些扩展时，可通过 vsx 扩展文件安装。

### 通过扩展 id 安装

`code-server --install-extension ${扩展 id}`

扩展 id 如何查看：在 WebIDE 或 `open-vsx` 搜索扩展，在详情页查看扩展 id，例如：安装 python 扩展，id 为 `ms-python.python`

### 通过 vsx 扩展文件安装

将 vsx 安装包下载并提交到仓库，这样就可以通过 vsx 文件来安装：

`code-server --install-extension ms-python.vscode-pylance.vsix`

下面以 python 开发需要安装 `python`（openvsx 有这个扩展） 和 `pylance`（openvsx 没有这个扩展） 扩展为例
（需配置 .cnb.yml 和 .ide/Dockerfile 文件）：

```yaml
# .cnb.yml
$:
  # vscode 事件：专供页面中启动远程开发用
  vscode:
    - docker:
        # 自定义镜像作为开发环境
        build:
          dockerfile: .ide/Dockerfile
          # by: 声明构建镜像需要用到的文件
          by:
            - .ide/ms-python.vscode-pylance.vsix
          # versionBy: 声明版本控制需要用到的文件
          # 当 .ide/Dockerfile 和 versionBy 中的文件有更新时，会重新构建镜像
          versionBy:
            - .ide/ms-python.vscode-pylance.vsix
      services:
        - vscode
        - docker
      stages:
        - name: ls
          script: ls -al
```

注意：.ide 目录下需有 `ms-python.vscode-pylance.vsix` 扩展文件。

```dockerfile{7,9}
# .ide/Dockerfile
FROM python:3.8

COPY .ide/ms-python.vscode-pylance.vsix .

# 安装 code-server 和扩展（使用 id 安装 python 扩展，使用 vsix 安装包安装 pylance 扩展）
RUN curl -fsSL https://code-server.dev/install.sh | sh \
  && code-server --install-extension ms-python.python \
  && code-server --install-extension ./ms-python.vscode-pylance.vsix \
  && echo done

# 安装 ssh 服务，用于支持 VSCode 客户端通过 Remote-SSH 访问开发环境
RUN apt-get update && apt-get install -y wget unzip openssh-server

# 指定字符集支持命令行输入中文（根据需要选择字符集）
ENV LANG C.UTF-8
ENV LANGUAGE C.UTF-8
```

---
title: UI 定制
permalink: https://docs.cnb.cool/zh/repo/settings.html
summary: 通过在仓库根目录添加 `.cnb/settings.yml` 配置文件，可以对页面部分 UI 进行定制，包括定制云原生开发启动按钮、创建 issue 按钮、fork 按钮及“抄了么”按钮的描述和图片等，且各参数多数为可选，配置解析失败或文件过大则不生效，“抄了么”按钮需在仓库简介加 `example` 标签才会显示 。
---

可通过在仓库根目录新增 `.cnb/settings.yml` 配置文件对页面部分 UI 进行定制，解锁更多玩法

## 配置文件说明

需在仓库根目录新增并提交 `.cnb/settings.yml` 配置文件，如下为配置文件示例：

```yaml
# .cnb/settings.yml
# 如下参数均为可选参数

# 云原生开发配置，读取云原生启动按钮所在页面当前分支的 .cnb/settings.yml 配置
workspace:
  launch:
    # 定制云原生开发启动按钮
    button:
      # 按钮名称
      name: 启动云原生开发
      # 按钮描述
      # 如果值为 null，则不显示默认描述
      description: 点击此按钮启动云原生开发环境
      # 鼠标悬浮在按钮上显示的图片（只能用仓库中当前分支的图片，填写相对仓库根目录的路径，如 .cnb/launch-hover.gif）
      # 图片最大 10MB
      hoverImage: .cnb/launch-hover.gif
    # CPU 核心数，默认为：8。仅对默认模版有效，如果有自定义云原生开发启动流水线，则此配置无效
    cpus: 4
    # 是否禁用默认按钮。默认为：false 表示不禁用。true 表示禁用
    disabled: false
    # 环境创建完是否自动打开 WebIDE，默认为 false
    # 当开发环境中未安装 openssh(仅支持 WebIDE)：无论此参数配置为 true 还是 false，环境创建完都将自动打开 WebIDE
    autoOpenWebIDE: false

# issue 配置，读取仓库主干 .cnb/settings.yml 配置
issue:
  # 定制创建 issue 按钮
  button:
    # 按钮描述
    description: ~bug~ 给你!
    # 鼠标悬浮在按钮上显示的图片（只能用仓库中当前分支的图片，填写相对仓库根目录的路径，如 .cnb/issue-hover.png）
    # 图片最大 10MB
    hoverImage: ".cnb/issue-hover.png"

# fork 配置，读取仓库主干 .cnb/settings.yml 配置
fork:
  # 定制 fork 按钮
  button:
    # 按钮描述
    description: 你的仓库不错，现在是我的了
    # 鼠标悬浮在按钮上显示的图片（只能用仓库中当前分支的图片，填写相对仓库根目录的路径，如 .cnb/fork-hover.png）
    # 图片最大 10MB
    hoverImage: ".cnb/fork-hover.png"

# 抄了么 按钮配置，读取仓库主干 .cnb/settings.yml 配置
# 仓库首页默认没有此按钮
# 抄了么按钮出现条件：仓库首页右方简介增加 example 标签（推荐可作为模板或例子的仓库增加此按钮）
copyRepo:
  # 定制 抄了么 按钮
  button:
    # 按钮描述
    description: 你的仓库不错，现在是我的了
    # 鼠标悬浮在按钮上显示的图片（只能用仓库中当前分支的图片，填写相对仓库根目录的路径，如 .cnb/copy-hover.png）
    # 图片最大 10MB
    hoverImage: ".cnb/copy-hover.png"
```

:::tip
注意：当文件解析失败或大小超过限制，将不会使用该配置
:::

## 抄了么按钮

仓库首页默认不会出现 `抄了么` 按钮，推荐可作为模板或例子的仓库增加此按钮。

如何增加此按钮：在仓库首页右侧简介编辑按钮，增加 `example` 标签，即可出现此按钮

点击该按钮，获取快速复制该仓库的具体方法

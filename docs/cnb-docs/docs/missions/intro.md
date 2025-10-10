---
title: 任务集介绍
permalink: https://docs.cnb.cool/zh/missions/intro.html
summary: CNB项目的任务集是一个协同工具，提供可视化看板，用于跨组织和仓库集中管理Issue/PR。用户可自定义视图和管理任务集权限，团队能有效协作处理开发中的问题。
---

任务集是 CNB 项目协同工具，提供可视化任务看板，帮助团队跨组织/跨仓库集中化管理 Issue/PR。

## 创建任务集

进入 [CNB](//cnb.cool)，单击右上角的“**＋**”，选择 **创建任务集**；

或者在组织目录下，进入 **任务集** ，在任务集列表中点击 **创建任务集**。

- 任务集名称：任务集的标识，也会组成任务集的访问路径。

- 数据范围：选择仓库，任务集会自动读取所选仓库的 Issue/PR 作为数据来源。

   ![](https://docs.cnb.cool/images/missions/fb4f5e30f29c11efbda7525400454e06.png)

## 任务集视图

在任务集中，您可以：

- 自由创建/切换表格视图、看板视图，以不同视角管理 Issue/PR。

- 通过筛选、分组、排序来自定义你的视图结构。

- 快速创建指定归属仓库的 Issue。

- 批量修改 Issue/PR 的属性值。

- 通过移动拖放，快速修改事项的优先级。

- 收藏任务集，方便后续快速访问

**表格视图效果如下：**

![](https://docs.cnb.cool/images/missions/4263f8d1f2a011efa823525400e889b2.png)

**看板视图效果如下：**

![](https://docs.cnb.cool/images/missions/511113a0f2a011efb67252540099c741.png)

## 任务集设置

在任务集中单击右侧的![](https://docs.cnb.cool/images/missions/20ada5d5f34211efb25a5254007c27c5.png)，然后选择**设置**，即可进入任务集设置页。可以进行基础设置、高级设置和管理任务集成员。

![](https://docs.cnb.cool/images/missions/14f67a72f34211ef920e5254005ef0f7.png)

### 基础设置

可以修改任务集名称和编辑任务集的数据范围等任务集的基础信息。

:::tip
修改任务集名称将会变更任务集的访问 URL，需谨慎操作。
:::

![](https://docs.cnb.cool/images/missions/d92796edf29e11efa8355254001c06ec.png)

### 高级设置

可以复制任务集、修改任务集的可见性和删除任务集。

:::tip
修改任务集可见性会影响任务集的访问范围，删除任务集后数据无法恢复，需谨慎操作。
:::

![](https://docs.cnb.cool/images/missions/06f68692f29f11ef9b7d525400bf7822.png)

### 任务集成员

默认继承上级组织的所有成员权限关系。此外，还可单击**邀请成员**直接邀请成员加入当前任务集成为“任务集成员”或“外部协作者”。其中，外部协作者适用于一些临时参与协作的用户。

![](https://docs.cnb.cool/images/missions/8ec50f63f29f11efbda7525400454e06.png)

成员角色权限如下：

![](https://docs.cnb.cool/images/missions/role.png)

## 任务集权限

1、用户能访问任务集视图，需要同时满足以下两个权限：

- 用户为该任务集成员

- 用户为该任务集关联仓库的成员

2、用户在任务集下，对Issue/PR 的操作权限，与他所在对应仓库的权限一致。

例：任务集 A 关联了仓库 A 、仓库 B; 小明对仓库 A 的 Issue 有编辑标题的权限，对仓库 B 的 Issue 仅有只读权限。那小明在任务集 A 中， 对仓库 A 的所有 Issue 有编辑标题的权限，对仓库 B 的所有 Issue 仅有只读权限，无法操作。

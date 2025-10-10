---
title: 元数据
permalink: https://docs.cnb.cool/zh/repo/annotations.html
summary: 元数据用于给仓库的 Tag 和 Commit 存储注解，以`key: value`形式存储。可通过仓库页面查看（不支持操作），亦可在云原生构建流水线中通过插件对元数据进行增加、删除和查询操作 。
---

## 什么是元数据？

元数据是仓库提供的数据存储能力，目前支持给 Tag 和 Commit 存储元数据，表示对 Tag 或 Commit 的注解。

## 存储形式

以 `key: value` 形式进行存储。

以 Tag v1.0.0 的元数据为例，其元数据支持 `key: value` 形式存储，如下：

```yaml
key1: value1
key2: value2
# ...
```

## 如何查看元数据？

页面查看路径：

- Tag 元数据：仓库 > Tag 或 Release 列表 > Tag 或 Release 详情 > 元数据（没有元数据时不展示）
- Commit 元数据：仓库 > Commit 列表 > Commit 详情 > 元数据

目前页面仅支持元数据查看，不支持增加、修改和删除。

## 如何操作元数据？

可在云原生构建流水线中，通过插件对元数据进行增加、删除、查询元数据。

详见 [元数据插件](../build/plugin.md#public/cnbcool/annotations)。

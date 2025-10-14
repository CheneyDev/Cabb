---
title: "time"
source_url: https://open.feishu.cn/document/client-docs/block/block-frame/code-components-and-structure/view-layer/ttss/basic-data-type/time
last_remote_update: 2022-07-15
last_remote_update_timestamp: 1657871784000
---
最后更新于 2022-07-15

# time

`<time>`数据类型表达了以秒（`s`）或毫秒（`ms`）为单位的时间的值。于`animation`、`transition`及相关属性中使用。

## 语法

`<time>`数据类型由`<number>`后加下方列出的单位之一组成。前方也可选择添加一个`+`或`-`。单位与数字之间无空格。

### 单位

-   `s` 以秒为单位的时间。 例：`0s`, `1.5s`, `-60s`。

-   `ms` 以毫秒为单位的时间。例：`0ms`, `150.25ms`, `-60000ms`。

## 示例

### 有效时间

```
12s         正整数
-456ms      负整数
4.3ms       非整数
14mS        单位大小写不敏感，虽然不推荐大写
+0s         0以及前置 + 号和单位
-0ms        0以及前置 - 号和单位
```

### 无效时间

```
0           尽管<length>中允许0，其在<time>中无效。
12.0        此为<number>而非<time>。其缺乏单位。
7 ms        数字与单位间不允许有空格。
```

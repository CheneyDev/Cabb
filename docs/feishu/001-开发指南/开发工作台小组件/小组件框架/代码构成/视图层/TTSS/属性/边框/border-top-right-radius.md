---
title: "border-top-right-radius"
source_url: https://open.feishu.cn/document/client-docs/block/block-frame/code-components-and-structure/view-layer/ttss/attributes/border/border-top-right-radius
last_remote_update: 2022-07-15
last_remote_update_timestamp: 1657871781000
---
最后更新于 2022-07-15

# border-top-right-radius

## 介绍

用于添加右上角圆角边框。第一个值是水平半径，第二个值是垂直半径。如果省略第二个值，则复制第一个值。如果长度为零，则边角为方形，而不是圆形。水平半径的百分比值参考边框盒的宽度，而垂直半径的百分比值参考边框盒的高度。

## 语法

```css
/* the corner is a circle */
/* border-top-right-radius: radius */
border-top-right-radius: 3px;
/* the corner is an ellipsis */
/* border-top-right-radius: horizontal vertical */
border-top-right-radius: 0.5em 1em;
```

### 取值

-   `<length>`

定义圆角的形状。

-   `<percentage>`

以百分比定义圆角的形状

## 标准化语法

```css
border-top-right-radius: [<length> | <percentage>] [ / [<length> | <percentage>]]?
```

---
title: "font-weight"
source_url: https://open.feishu.cn/document/client-docs/block/block-frame/code-components-and-structure/view-layer/ttss/attributes/text/font-weight
last_remote_update: 2022-07-15
last_remote_update_timestamp: 1657871780000
---
最后更新于 2022-07-15

# font-weight

## 介绍

设置字体粗细。

## 语法

```css
font-weight: bold;

font-weight: normal;

font-weight: 100;

font-weight: 200;

...

font-weight: 900;
```

### 取值

-   `bold`

粗体，与 700 字重相同。

-   `normal`

普通，与 400 字重相同。

-   `<number>`

只能使用 100、200、300、400、500、600、700、800、900，其他值无效。当取值为 500 - 900 时会加粗。

## 例子

```html
<text style="font-weight: normal;">font-weight: normal; </text>

<text style="font-weight: bold;">font-weight: bold; </text>

<text style="font-weight: 100;">font-weight: 100; </text>

<text style="font-weight: 400;">font-weight: 400; </text>

<text style="font-weight: 500;">font-weight: 500; </text>

<text style="font-weight: 900;">font-weight: 900; </text>
```

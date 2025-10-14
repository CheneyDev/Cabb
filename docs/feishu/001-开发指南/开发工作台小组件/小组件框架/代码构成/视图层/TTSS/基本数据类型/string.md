---
title: "string"
source_url: https://open.feishu.cn/document/client-docs/block/block-frame/code-components-and-structure/view-layer/ttss/basic-data-type/string
last_remote_update: 2022-07-15
last_remote_update_timestamp: 1657871784000
---
最后更新于 2022-07-15

# string

`<string>`是用来表示一串字符的数据类型。

## 语法

`<string>`数据类型是由包含在英文双引号（`"`）或英文单引号（`'`）中的任意数量的`Unicode`字符组成。

大多数字符都可以写成字面量的形式，并且所有字符都可以写成以反斜线（``）开头的十六进制`Unicode`码点的形式。比如说：`\22`表示一个双引号，`\27`表示一个单引号，`\A9`表示版权声明符号（©）。

需要注意的是，在使用某些特殊字符时需要在前面加上反斜线来转义，包括但不限于在双引号字符串中使用双引号字符，在单引号字符串中使用单引号字符，以及反斜线字符，比如`\`表示一个反斜线字符。

## 示例

```html
/* 简单字符串 */
"This string is demarkated by double quotes."
'This string is demarkated by single quotes.'
/* 转义字符 */
"This is a string with " an escaped double quote."
"This string also has \22 an escaped double quote."
'This is a string with ' an escaped single quote.'
'This string also has \27 an escaped single quote.'
"This is a string with \ an escaped backslash."
```

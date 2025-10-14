---
title: "angle"
source_url: https://open.feishu.cn/document/client-docs/block/block-frame/code-components-and-structure/view-layer/ttss/basic-data-type/angle
last_remote_update: 2022-07-15
last_remote_update_timestamp: 1657871783000
---
最后更新于 2022-07-15

# angle

角度`<angle>`用于表示角的大小，单位为度（`degrees`）、 弧度（`radians`）或圈数（`turns`）。在和`transform`的某些方法等场景中有所应用。

## 语法

`<angle>`数据类型由`<number>`和下列单位组成。数字与单位之间没有空格。数字为`0`时，单位可以省略。

`<angle>`可以使用`+`或`-`开头。正数表示顺时针的角，负数表示逆时针的角。对于静态的角，同样的角度可以使用任意等效的值表示。比如 `90deg`等于`-270deg`，`1turn`等于`4turn`。而对于动态的角，比如应用了`animation`或`transition`时，显示效果则不一样。

## 单位

-   `deg`

度。一个完整的圆是`360deg`。例：`0deg`，`90deg`，`14.23deg`。

-   `rad`

弧度。一个完整的圆是`2π`弧度，约等于`6.2832rad`。`1rad`是`180/π`度。例：`0rad`，`1.0708rad`，`6.2832rad`。

-   `turn`

圈数。一个完整的圆是`1turn`。例：`0turn`，`0.25turn`，`1.2turn`。

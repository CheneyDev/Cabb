---
title: "gradient"
source_url: https://open.feishu.cn/document/client-docs/block/block-frame/code-components-and-structure/view-layer/ttss/basic-data-type/gradient
last_remote_update: 2022-07-15
last_remote_update_timestamp: 1657871784000
---
最后更新于 2022-07-15

# gradient

渐变图像类型`<gradient>`用于表现两种或多种颜色的过渡转变。

`<gradient>`没有内在尺寸；即，不具备固有或首选的尺寸，也不具备首选的比率。其实际的大小取决于其填充元素的大小。

该数据类型通常用于`background`，`background-image`，`background-color`等属性之中。

## 语法

`<gradient>`数据类型由下列渐变函数定义。

### 线形渐变

颜色值沿着一条隐式的直线逐渐过渡。由`linear-gradient()`产生。

```css
/* 渐变轴为45度，从蓝色渐变到红色 */
linear-gradient(45deg, blue, red);
```

```css
/* 从右下到左上、从蓝色渐变到红色 */
linear-gradient(to left top, blue, red);
```

```css
/* 从下到上，从蓝色开始渐变、到高度40%位置是绿色渐变开始、最后以红色结束 */
linear-gradient(0deg, blue, green 40%, red);
```

#### 标准语法

```css
linear-gradient([ <angle> | to <side-or-corner> ,]? <color-stop> [ , <color-stop> ]+ )
```

#### 可填值

-   `<side-or-corner>` 描述渐变线的起始点位置。

它包含`to`和**两个关键词**：第一个指出水平位置 left or right，第二个指出垂直位置 top or bottom。关键词的先后顺序无影响，且都是可选的。

`to top`, `to bottom`, `to left`和`to right`这些值会被转换成角度`0`度、`180`度、`270`度和`90`度。其余值会被转换为一个以向顶部中央方向为起点顺时针旋转的角度。渐变线的结束点与其起点中心对称。

-   `<angle>` 用角度值指定渐变的方向（或角度）。角度顺时针增加。

-   `<color-stop>` 由一个`<color>`值组成，并且跟随着一个可选的终点位置（可以是一个百分比值或者是沿着渐变轴的`<length>`）。

### 径向渐变

颜色值由一个中心点（原点）向外扩散并逐渐过渡到其他颜色值。由`radial-gradient()`函数产生。

```css
/* 从中心向四周渐变 */
radial-gradient(#e66465, #9198e5);
```

```css
/* 从边缘相切处开始渐变 */
radial-gradient(closest-side, #3f87a6, #ebf8e1, #f69d3c);
```

```css
/* 圆形渐变 */
radial-gradient(circle at 100%, #333, #333 50%, #eee 75%, #333 75%);
```

```css
/* 椭圆渐变 */
radial-gradient(ellipse at top, #e66465, transparent);
```

#### 标准语法

```css
radial-gradient(
  [ [ circle || <length> ][ at <position> ]? , |
    [ ellipse || [ <length> | <percentage> ]{2} ]  [ at <position> ]? , |
    [ [ circle | ellipse ] || <extent-keyword> ] [at <position> ]? , |
    at <position> ,
  ]?
  <color-stop> [ , <color-stop> ]+
)
```

#### 可填值

-   `<position>` `<position>`与`background-position`或者`transform-origin`类似。如缺省，默认为中心点。

-   `<shape>` 渐变的形状。圆形（`circle`）或椭圆形（`ellipse`）。默认值为椭圆。

-   `<size>` 确定渐变的最终形状的大小。

-   `<extent-keyword>` 当渐变的最终形状大小是在某个方向上的无限延伸，而不是有限的线段时，可以使用关键字来描述渐变形状的大小，默认为`farthest-corner`。

| 关键字             | 描述                                                       |
| --------------- | -------------------------------------------------------- |
| closest-side    | 渐变的边缘形状与容器距离渐变中心点最近的一边相切（圆形）或者至少与距离渐变中心点最近的垂直和水平边相切（椭圆）。 |
| closest-corner  | 渐变的边缘形状与容器距离渐变中心点最近的一个角相交。                               |
| farthest-side   | 与 closest-side 相反，边缘形状与容器距离渐变中心点最远的一边相切（或最远的垂直和水平边）。     |
| farthest-corner | 渐变的边缘形状与容器距离渐变中心点最远的一个角相交。                               |

-   `<color-stop>` 同`linear-gradient`的描述。

---
title: 插件制作
permalink: https://docs.cnb.cool/zh/build/create-plugin.html
summary: 在“云原生构建”中，插件是一个 Docker 镜像，从零开发插件的步骤包括：设计插件参数（这些参数会转化为特定格式的环境变量传给容器，支持多种类型，复杂参数可存文件）、书写脚本、构建插件镜像（通过 Dockerfile）、测试插件（本地通过 docker run 运行并传参）、导出变量（供后续任务使用）、发布插件（可发布到“云原生构建”仓库的制品库或 Docker Hub）。
---

在 `云原生构建` 中，一个插件就是一个 `Docker` 镜像。
下面我们介绍如何使用 `Bash` 从零开发一个镜像插件。
这个插件的功能是打印 `hello world`。
这篇内容应该可以给你创作自己的插件提供一个清晰的思路。
我们这里假设你已经知道 `Docker` 的一些基本知识。

## 设计插件

### 参数设计

第一步应该是去设计插件所需要的参数：

- `text`: 要输出到控制台的文本内容

```yaml
main:
  push:
    - stages:
        - name: hello world
          image: cnbcool/hello-world
          settings:
            text: hello world
```

这些入参将会以环境变量的形式传给容器，不同的是，他们将会变成大写且辅以 `PLUGIN_` 前缀。

上面的入参将会转化为如下环境变量：

```text
PLUGIN_TEXT="hello world"
```

### 支持的参数类型

参数类型支持`字符串`、`数值`、`布尔值`、`一维数组`、`普通对象`

其中：

- 数组在传给容器时将会以英文逗号 `,` 分割
- 普通对象在传给容器时，会转成 `JSON` 字符串

比如：

```yaml
main:
  push:
    - stages:
        - name: hello world
          image: cnbcool/hello-world
          settings:
            boolean: true
            number: 123
            array: [hello, world]
            map:
              key: value
```

上传参数值将会转化为以下环境变量：

```text
PLUGIN_BOOLEAN='true'
PLUGIN_NUMBER='123'
PLUGIN_ARRAY='hello,world'
PLUGIN_MAP='{"key":"value"}'
```

特别复杂的参数值可以存到一个文件中，插件运行时加载。
如果你遇到参数值异常复杂的情况，往往不是格式能解决的，应当简化这些参数，或者将他们做成多个插件。

## 书写脚本

下一步写一个打印参数的 Bash 脚本，如下：

```bash
#!/bin/sh
echo "$PLUGIN_TEXT"
```

## 构建插件镜像

插件将会被打包成 `Docker` 镜像进行分发使用。
因此需要创建一个 `Dockerfile` 把我们之前写好的脚本打包进去，
并且把它设置为 `Entrypoint`。

```bash
FROM alpine

ADD entrypoint.sh /bin/
RUN chmod +x /bin/entrypoint.sh

ENTRYPOINT /bin/entrypoint.sh
```

构建你的镜像:

```bash
docker build -t cnbcool/hello-world .
```

## 测试插件

你应当在本地测试好你的插件，可以使用 `docker run` 来运行插件，并且把参数通过环境变量的方式传进去：

```bash
docker run --rm \
  -e PLUGIN_TEXT="hello world" \
  cnbcool/hello-world
```

### 测试文件系统读取

插件有读取你构建流程工作区目录的权限，它会默认把构建的目录映射到插件的某个目录，然后把这个目录设置为工作区：

```bash
docker run --rm \
  -e PLUGIN_TEXT="hello world" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  cnbcool/hello-world
```

## 导出变量

如果插件执行完后，需要返回结果并导出为变量供后续任务使用，可以参考 [exports](./grammar.md#exports)

## 发布插件

插件是一个 `Docker` 镜像，所以发布一个插件，就意味着需要把镜像发布到一个镜像源。

可以发布到 `云原生构建` 仓库自带的[制品库](../artifact/docker.md)。

对于全球可用的插件，也可以发布到 `Docker Hub`。

### 发布镜像

```bash
docker push cnbcool/hello-world
```

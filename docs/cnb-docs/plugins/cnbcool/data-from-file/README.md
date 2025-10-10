# data-from-file

解析本地文件，指定属性路径导出为指定名称的环境变量

属性路径参考 [lodash.get](https://www.lodashjs.com/docs/lodash.get)

文件格式可为 `.json`、`.yml`、`.yaml`

其他格式按照 `plain text` 处理，该格式不支持深层路径（不推荐）

> 原理是通过输出 `##[set-output key=value]`，让 `云原生构建`解析为对应的环境变量，会暴露到构建日志，不适合敏感信息处理。

## 参数

### file

- type: String
- required: 是

文件路径

### keys

- type: String
- required: 是

文件中对象的属性路径与导出的环境变量名的映射关系

每行文本代表一个映射：key_in_file:key_to_export

示例：

```shell
prop1:key1
prop2.sub3:key2
prop4.sub5[1].name:key3
```

### encoding

- type: String
- required: 否
- default: utf8

以该 encoding 解析文件内容

## 在 云原生构建 上使用

```yaml
# .cnb.yml
main:
  push:
    - stages:
      - name: export env from file
        image: cnbcool/data-from-file
        settings:
          file: env.yml
          keys: |
            prop1:key1
            prop2.sub3:key2
            prop4.sub5[1].name:key3
        exports:
          key1: KEY_1
          key2: KEY_2
          key3: KEY_3
      - name: show
        script: echo "$KEY_1 $KEY_2 $KEY_3"
```

```yaml
# env.yml
prop1: value1
prop2:
  sub3: value2
prop4:
  sub5:
    - name: tom cat
      age: 18
    - name: jerry mouse
      age: 24
```

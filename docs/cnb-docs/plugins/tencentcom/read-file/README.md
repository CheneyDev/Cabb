# read-files

读取本地文件内容。

普通的文件可以直接用 shell 命令直接读取输出，涉及到特殊编码，例如 GBK，可以使用该插件指定编码读取。

## 输入

### file_path

- Type: `String`
- 必填：是

指定要读取的本地文件路径。

### encoding

- Type: `String`
- 必填：否，默认值 `utf8`

指定文件内容以哪种编码解析。

### parse

- Type: `Boolean`
- 必填：否，默认值 false

是否需要解析文件内容。true: 解析文件内容；false：仅原样输出文件内容

## 输出

标准输出流输出内容。

## 在 Docker 上使用

```shell
docker run --rm -it \
  -v $(pwd):$(pwd) -w $(pwd) \
  -e PLUGIN_FILE_PATH="./test.txt" \
  tencentcom/read-file
```

## 在 云原生构建 上使用

### 仅输出文件内容

如下示例会原样输出文件内容：

```yaml
main:
  push:
  - stages:
    - name: 读取文件
      image: tencentcom/read-file
      settings:
        file_path: ./test.txt
```

### 需解析文件内容

参数 `parse: true` 时，会解析文件内容，并输出到环境变量。

支持解析的文件格式:

1、yaml：文件后缀为.yaml或.yml

```yaml
name1: val1
name2: 
  name4: val4
  name5: val5
name3:
  - val1
  - val2
  - val3
```

2、json：文件后缀为 json

```json
{
  "name1": "val1",
  "name2": {
    "name4": "val4",
    "name5": "val5"
  },
  "name3": [
    "val3",
    "val4",
    "val5"
  ]
}
```

3、其他文件，无论什么后缀，均按照 txt 文本格式读取

```txt
name1=1
name2=2
name3=3
```

```yaml
master:
  push:
  - stages:
    - name: 读取文件
      image: tencentcom/read-file
      settings:
        # 假如文件内容为上述 yaml 文件示例
        file_path: ./test.yaml
        parse: true
      exports:
        name1: NAME1
        name2.name4: NAME2
        name2.name5: NAME3
        name3.0: NAME4
        name3.1: NAME5
        name3.2: NAME6
    - name: 输出环境变量
      script:
        - echo $NAME1
        - echo $NAME2
        - echo $NAME3
        - echo $NAME4
        - echo $NAME5
        - echo $NAME6
```

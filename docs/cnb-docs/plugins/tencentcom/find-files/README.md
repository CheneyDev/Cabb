# find-files

根据 `glob` 表达式查询文件，并列出结果

## 输入

### include

Type: `String|String[]`

必填，指定要包含的文件，支持 glob 表达式。

## 输出

### files

Type: `String`

查询到的文件，多个文件用逗号连接。

## 在 Docker 上使用

```shell
docker run --rm -it \
  -v $(pwd):$(pwd) -w $(pwd) \
  -e PLUGIN_INCLUDE="**/*.js,*.js" \
  tencentcom/find-files
```

## 在 云原生构建 上使用

```yaml
main:
  push:
  - stages:
    - name: 查询文件
      image: tencentcom/find-files
      settings:
        include: 
         - "**/*.js"
         - '*.js'
```

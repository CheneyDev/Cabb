# ajv-validator

ajv-validator Docker 插件

检查Yaml文件是否符合JSON schema定义。

## 输入

- `schema`: `String` JSON schema定义文件的url。本地文件可用`file://`协议。
- `strict`: `Boolean` 是否使用strict模式，默认值为`false`。
- `yaml`: `String | String[]` 需要被检查的文件路径。

## 在 云原生构建 上使用

```yaml
main:
  push:
  - stages:
    - name: 检查yml格式是否正确
      image: tencentcom/ajv-validator:latest
      settings:
        schema: https://xxx/conf-schema.json
        yaml: test.yml
```

## 参考资料

- [https://json-schema.org](https://json-schema.org)
- [https://ajv.js.org](https://ajv.js.org)

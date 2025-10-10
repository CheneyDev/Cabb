# plugin-lint

检查插件元数据 plugin-meta.json 是否符合编写规范。

## 在 Docker 上使用

```shell
docker run --rm -t -v $(pwd):$(pwd) -w $(pwd) cnbcool/plugin-lint .
```

## 在 云原生构建 上使用

```yaml
# .cnb.yml
main:
  pull_request:
  - stages:
    - name: plugin-lint
      image: cnbcool/plugin-lint
```

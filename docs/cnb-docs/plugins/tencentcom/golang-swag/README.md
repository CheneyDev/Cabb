# golang-swag

将 golang 的注释转换为 swagger2.0 API 文档。

## 在 Docker 上使用

```shell
docker run --rm -it -v $(pwd):$(pwd) -w $(pwd) \
  --entrypoint="" \
  tencentcom/golang-swag \
  swag -h
```

```shell
docker run --rm -it -v $(pwd):$(pwd) -w $(pwd) \
  -e PLUGIN_DIR="./" \
  -e PLUGIN_EXTRA_ARGS="--parseDependency --parseDepth 3" \
  tencentcom/golang-swag
```

## 在 云原生构建 上使用

```yaml
main:
  push:
  - stages:
    - name: 生成swagger文档
      image: tencentcom/golang-swag
      settings:
        dir: ./
        extra_args: --parseDependency --parseDepth 3 
```

```yaml
main:
  push:
  - stages:
    - name: 生成swagger文档
      image: tencentcom/golang-swag
      commands: |
        swag -h
        swag init
```

## 更多用法

参考 [swaggo/swag](https://github.com/swaggo/swag)

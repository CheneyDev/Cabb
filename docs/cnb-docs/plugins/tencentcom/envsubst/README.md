# envsubst

替换文件中的环境变量

## 镜像

`tencentcom/envsubst:latest`

## 在 Docker 中使用

```shell
docker run --rm -v $(pwd):$(pwd) -w $(pwd) \
  -e PLUGIN_FILE=example.txt \
  tencentcom/envsubst:latest
```

## 在 云原生构建 中使用

```yml
main:
  push:
  - stages:
    - name: 替换文件中的环境变量，仅支持 CNB_ 开头的环境变量
      image: tencentcom/envsubst:latest
      settings:
        # 要替换的文件名
        file: example.txt
```

```yml
main:
  push:
  - stages:
    - name: 替换文件中的环境变量，支持自定义环境变量
      image: tencentcom/envsubst:latest
      env:
        ENV_A: value
      script:
        - envsubst < input.txt > output.txt
```

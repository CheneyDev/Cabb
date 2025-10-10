# coscli 插件

使用 coscli 插件，用户可通过简单的命令行指令实现对对象（Object）的批量上传、下载、删除等操作。

## 镜像

`tencentcom/coscli:latest`

## 在 Docker 中使用

```shell
docker run --rm -it -v $(pwd):$(pwd) -w $(pwd) cnbcool/coscli:latest --version
```

## 在云原生构建中使用

```shell
main:
  push:
  - stages:
    - name: run with tencentcom/coscli
      image: tencentcom/coscli
      script: |
        coscli config set --secret_id $SECRET_ID --secret_key $SECRET_KEY 
        coscli config add --init-skip=true -b $BUCKET -r $REGION
        coscli cp ./build cos://$BUCKET -r
```

## 更多用法

更多用法，参照[coscli 文档](https://cloud.tencent.com/document/product/436/63143)。

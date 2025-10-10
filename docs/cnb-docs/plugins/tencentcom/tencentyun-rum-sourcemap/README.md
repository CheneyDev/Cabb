# rum-sourcemap

通过腾讯云 API 接口上传 Sourcemap 文件到 [腾讯云前端性能监控RUM][yun-url] 服务

[yun-url]:https://cloud.tencent.com/product/rum

## 前置要求

1. 申请 [腾讯云API平台子用户][api]，获得 secretId 和 secretKey，如果用户之前已有子账号可以忽略此步骤。
1. 申请平台子账户必要权限：
   1. rum:CreateReleaseFile
   1. rum:DescribeReleaseFileSign
   1. rum:DescribeReleaseFiles
1. 从 [RUM 应用设置][setting] 获得 RUM 应用 ID（纯数字 ID，例如 12345）。

[api]:https://cloud.tencent.com/document/product/598/13674
[setting]:https://console.cloud.tencent.com/rum/web/group-projects-manage

## 在 Docker 中使用

```shell
docker run --rm \
  -e PLUGIN_SECRETID=xxxxxx \
  -e PLUGIN_SECRETKEY=xxxxxx \
  -e PLUGIN_ENDPOINT=rum.internal.tencentcloudapi.com \
  -e PLUGIN_PROJECTID=123456 \
  -e PLUGIN_VERSION=1.0.0 \
  -e PLUGIN_FILES=./dist/**/*.js.map \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  tencentcom/tencentyun-rum-sourcemap:latest
```

## 在 云原生构建 中使用

```yml
main:
  push:
  - stages:
    - name: 上传 sourcemap 到 RUM
      image: tencentcom/tencentyun-rum-sourcemap:latest
      settings:
        secretId: xxxxxx
        secretKey: xxxxxx
        projectId: 123456
        version: 1.0.0
        files: ./dist/**/*.js.map
        # 内部域名，如需要
        endpoint: rum.internal.tencentcloudapi.com
```

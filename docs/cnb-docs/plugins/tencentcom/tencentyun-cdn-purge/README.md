# tencentyun-cdn-purge

通过腾讯云 API 接口刷新 CDN 缓存，适用于源站资源更新和发布后让用户更快获取到最新内容。

- [刷新 URL](https://cloud.tencent.com/document/api/228/37870)
- [刷新目录](https://cloud.tencent.com/document/api/228/37871)

## 在 Docker 中使用

```shell
docker run --rm \
  -e PLUGIN_SECRETID=****** \
  -e PLUGIN_SECRETKEY=****** \
  -e PLUGIN_ENDPOINT=cdn.internal.tencentcloudapi.com \
  -e PLUGIN_URLS=https://example.com/path/,https://example.com/path/file.js \
  tencentcom/tencentyun-cdn-purge:latest
```

## 在 云原生构建 中使用

```yml
main:
  push:
  - stages:
    - name: 刷新 CDN 缓存
      image: tencentcom/tencentyun-cdn-purge:latest
      settings:
        secretId: xxxxxx
        secretKey: xxxxxx
        urls: 
          - https://cdn.example.com/path/ # 刷新目录需要以 / 结尾
          - https://cdn.example.com/path/file.ext
        flushType: flush
        endpoint: cdn.internal.tencentcloudapi.com # 内部 AKSK 调用需要使用内部域名
```

## 参数说明

| 参数名    | 是否必填 | 默认值  | 说明                                             |
| --------- | -------- | ------- | -------------------------------------------- |
| secretId  | 是       | 空      | 腾讯云 API secretId                           |
| secretKey | 是       | 空      | 腾讯云 API secretKey                           |
| urls      | 是       | 空      | 需要刷新的 CDN 路径，可以是文本或者数组             |
| flushType | 否       | `flush` | 目录的刷新类型                                  |
| endpoint  | 否       | `cdn.tencentcloudapi.com` | 调用域名 |

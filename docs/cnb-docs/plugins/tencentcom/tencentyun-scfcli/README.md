# tencentyun-scfcli

SCF(Serverless Cloud Function)

SCF CLI 是腾讯云无服务器云函数 SCF（Serverless Cloud Function）产品的命令行工具。
通过 scf 命令行工具，您可以方便的实现函数打包、部署、本地调试，
也可以方便的生成云函数的项目并基于 demo 项目进一步的开发。

本项目主要提供了 scf cli 的 Docker 镜像，实现一键部署到云函数 SCF（Serverless Cloud Function）。

## 输入

### appid

账号 ID：即 APPID。

通过访问控制台中的【账号中心】>[【账号信息】][appid-url]，可以查询到您的账号 ID。

[appid-url]:https://console.cloud.tencent.com/developer

### region

地域：产品期望所属的地域。

### secret_id

SecretID 及 SecretKey：指云 API 的密钥 ID 和密钥 Key。
您可以通过登录【[访问管理控制台](https://console.cloud.tencent.com/cam/overview)】，
选择【云 API 密钥】>【[API 密钥管理](https://console.cloud.tencent.com/cam/capi)】，
获取相关密钥或创建相关密钥。

### secret_key

同 secret_id

## 在 Docker 上使用

```shell
docker run --rm -it \
  -e PLUGIN_APPID=12539702XX \
  -e PLUGIN_REGION=ap-guangzhou \
  -e PLUGIN_SECRET_ID=your-secret-id \
  -e PLUGIN_SECRET_KEY=your-secret-key \
  -v $(pwd):$(pwd) -w $(pwd) tencentcom/tencentyun-scfcli
```

## 在 云原生构建 上使用

```yaml
main:
  push:
  - stages:
    - name: scf deploy
      image: tencentcom/tencentyun-scfcli
      settings:
        appid: 12539702XX
        region: ap-guangzhou
        secret_id: your-secret-id
        secret_key: your-secret-key
```

## 更多用法

参见：[tencentyun/scfcli](https://github.com/tencentyun/scfcli)

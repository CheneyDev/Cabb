# tencentyun-wecos

腾讯云 wecos 微信小程序 COS 瘦身解决方案

通过 WeCOS，将小程序图片资源上传到 腾讯云COS，
自动引用线上地址，移除图片资源，解决包大小限制的问题。

## 在 Docker 上使用

```shell
docker run --rm -v $(pwd):$(pwd) -w $(pwd) tencentcom/tencentyun-wecos
```

## 在 云原生构建 上使用

```yaml
main:
  push:
  - stages:
    - name: tencentyun-wecos
      image: tencentcom/tencentyun-wecos
```

## 配置文件

在你的小程序目录同级下创建`wecos.config.json`文件

`wecos.config.json`配置项例子：

```json
{
  "appDir": "./app",
  "cos": {
    "secret_id": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    "secret_key": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "bucket": "wxapp-1251902136",
    "region": "ap-guangzhou", //创建bucket时选择的地域简称
    "folder": "/", //资源存放在bucket的哪个目录下
  }
}
```

### appDir

小程序项目目录，默认为 `./app`。

### cos

必填，填写需要上传到COS对应的配置信息，部分信息可在[COS控制台][url]查看

[url]:https://console.qcloud.com/cos4/secret

## 更多用法

参见： [腾讯云 wecos][yun-url]

[yun-url]:https://github.com/tencentyun/wecos

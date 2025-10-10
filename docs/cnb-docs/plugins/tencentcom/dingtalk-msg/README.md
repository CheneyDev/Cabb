# dingtalk-msg

发送钉钉消息

## 在 云原生构建 上使用

```yml

main:
  push:
    - stages:
      - name: dingtalk-msg
        imports: https://xxx/envs.yml
        image: tencentcom/dingtalk-msg:latest
        settings:
          content: "your message"
          to: "receiver"
          c_type: "text"
          appKey: $APPKEY
          appSecret: $APPSECRET
          agentId: $AGENTID


```

envs.yml文件示例：

 ```yml

APPKEY: xxx
APPSECRET: xxx
AGENTID: xxx
 ```

## 参数

* `content`：消息内容

* `to`：接收人

* `c_type`：消息类型。可选：text,markdown,file

* `appKey`：AppKey [需要管理员身份在钉钉开放平台查看](https://open-dev.dingtalk.com/)

* `appSecret`：AppSecret[需要管理员身份在钉钉开放平台查看](https://open-dev.dingtalk.com/)

* `agentId`：需要企业管理员身份在钉钉开放平台创建企业内钉钉应用，权限管理需要添加【权限范围-全部员工】和
[通讯录管理-根据手机号姓名获取成员信息的接口访问权限](https://open-dev.dingtalk.com/fe/app?spm=a2q3p.21071111.0.0.6d921cfaNXjuer#/corp/app)

## 更多用法

更多用法参考：[钉钉开放平台帮助文档](https://open.dingtalk.com/document/)

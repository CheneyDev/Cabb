# dingtalk-bot-msg

发送钉钉机器人消息

## 在 云原生构建 上使用

```yml

main:
  push:
    - stages:
      - name: dingtalk-bot-msg
        imports: https://xxx/envs.yaml
        image: tencentcom/dingtalk-bot-msg:latest
        settings:
          content: "your message"
          c_type: "text"
          webhook: $WEBHOOK
          at: "199xxxxxx"
          isAtAll: false

```

envs.yaml文件示例：

 ```yml

WEBHOOK: xxx

 ```

## 参数

* `content`：消息内容

* `c_type`：消息类型。可选：text,markdown

* `webhook`：钉钉机器人 WebHook [需要在PC端钉钉客户端创建]。注意创建钉钉机器人需要配置一个自定义关键字："CODING"，
[参考文档](https://developers.dingtalk.com/document/app/custom-robot-access)

* `at`：需要 at 的人。填写需要 at 人的手机号，多个以 ";" 隔开

* `isAtAll`：是否 at 所有人。bool

## 更多用法

更多用法参考：[钉钉开放平台帮助文档](https://open.dingtalk.com/document/)

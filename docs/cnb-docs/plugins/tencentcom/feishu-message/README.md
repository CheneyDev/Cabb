# feishu-message

利用飞书自定义机器人发送消息，详情参考[飞书文档](https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot)。

## 参数

- sign_secret: `String`，飞书自定义机器人签名密钥。若机器人未开启签名，可忽略此参数。
- msg_type: `String`，消息类型，可选值有 `text`、`post`、`interactive`，默认为 `text`。
- content: `String | Object`，消息内容。 `msg_type` 为 `text`、`post` 时，必填。内容格式参考下文格式说明。
- card: `Object`，卡片消息内容。`msg_type` 为 `interactive` 时，必填。内容格式参考下文格式说明。
- filepath: `String`，消息内容文件路径。当 存在该值时，会从文件读取内容，覆盖 `content` 和 `card` 参数。内容格式参考下文格式说明。

## 消息内容格式说明

### content

当 msg_type 为 `text` 时，消息内容格式为：

```json
{
  "text": "新更新提醒"
}
```

若解析失败，会将内容原样发送。即会自动组装成如下格式：

```json
{
  "text": "$content"
}
```

当 msg_type 为 `post` 时，消息内容格式为：

```json
{
  "post": {
    "zh_cn": {
      "title": "项目更新通知",
      "content": [
        [
          {
            "tag": "text",
            "text": "项目有更新: "
          },
          {
            "tag": "a",
            "text": "请查看",
            "href": "http://www.example.com/"
          },
          {
            "tag": "at",
            "user_id": "ou_18eac8********17ad4f02e8bbbb"
          }
        ]
      ]
    }
  }
}
```

### card

卡片消息内容格式为：

```json
{
  "elements": [
    {
      "tag": "div",
      "text": {
        "content": "**西湖**，位于浙江省杭州市西湖区龙井路1号，杭州市区西部，景区总面积49平方千米，汇水面积为21.22平方千米，湖面面积为6.38平方千米。",
        "tag": "lark_md"
      }
    },
    {
      "actions": [
        {
          "tag": "button",
          "text": {
            "content": "更多景点介绍 :玫瑰:",
            "tag": "lark_md"
          },
          "url": "https://www.example.com",
          "type": "default",
          "value": {}
        }
      ],
      "tag": "action"
    }
  ],
  "header": {
    "title": {
      "content": "今日旅游推荐",
      "tag": "plain_text"
    }
  }
}
```

### filepath

filepath 对应文件内容格式为：

```json
{
  "content": {
    "post": {
      "zh_cn": {
        "title": "项目更新通知",
        "content": [
          [
            {
              "tag": "text",
              "text": "项目没有更新: "
            },
            {
              "tag": "a",
              "text": "请查看",
              "href": "http://www.example.com/"
            }
          ],
          [
            {
              "tag": "text",
              "text": "项目没有更新: "
            },
            {
              "tag": "a",
              "text": "请查看",
              "href": "http://www.example.com/"
            }
          ]
        ]
      }
    }
  },
  "card": {
    "elements": [
      {
        "tag": "div",
        "text": {
          "content": "**东湖**，位于浙江省杭州市西湖区龙井路1号，杭州市区西部，景区总面积49平方千米，汇水面积为21.22平方千米，湖面面积为6.38平方千米。",
          "tag": "lark_md"
        }
      },
      {
        "actions": [
          {
            "tag": "button",
            "text": {
              "content": "更多景点介绍 :玫瑰:",
              "tag": "lark_md"
            },
            "url": "https://www.example.com",
            "type": "default",
            "value": {}
          }
        ],
        "tag": "action"
      }
    ],
    "header": {
      "title": {
        "content": "今日旅游推荐",
        "tag": "plain_text"
      }
    }
  }
}
```

按 `msg_type`，存在 `content` 和 `card` 其中一种内容即可。

## 云原生构建 示例

可将 robot 地址和 签名密钥配置在私有仓库中，然后通过 imports 导入环境变量中。

```yaml
# env.yml
ROBOT: https://open.feishu.cn/open-apis/bot/v2/hook/xx
SIGN_SECRET: xxx
```

### 文本消息

```yaml
main:
  pull_request:
    - imports: https://xxx/env.yml
      stages:
        - name: send simple message
          image: tencentcom/feishu-message
          settings:
            robot: $ROBOT
            content: 新更新提醒
        - name: send  message
          image: tencentcom/feishu-message
          settings:
            robot: $ROBOT
            content: |
              {
                  "text": "新更新提醒"
              }
```

### 富文本消息

```yaml
main:
  pull_request:
    - imports: https://xxx/env.yml
      stages:
        - name: send message
          image: tencentcom/feishu-message
          settings:
            robot: $ROBOT
            msg_type: post
            content: |
              {
                  "post": {
                      "zh_cn": {
                          "title": "项目更新通知",
                          "content": [
                              [{
                                  "tag": "text",
                                  "text": "项目有更新: "
                              }, {
                                  "tag": "a",
                                  "text": "请查看",
                                  "href": "http://www.example.com/"
                              }, {
                                  "tag": "at",
                                  "user_id": "ou_18eac8********17ad4f02e8bbbb"
                              }]
                          ]
                      }
                  }
              }
```

### 卡片消息

```yaml
main:
  pull_request:
    - imports: https://xxx/env.yml
      stages:
        - name: send message
          image: tencentcom/feishu-message
          settings:
            robot: $ROBOT
            msg_type: interactive
            card: |
              {
                  "elements": [
                      {
                          "tag": "div",
                          "text": {
                              "content": "**西湖**，位于浙江省杭州市西湖区龙井路1号，杭州市区西部，景区总面积49平方千米，汇水面积为21.22平方千米，湖面面积为6.38平方千米。",
                              "tag": "lark_md"
                          }
                      },
                      {
                          "actions": [
                              {
                                  "tag": "button",
                                  "text": {
                                      "content": "更多景点介绍 :玫瑰:",
                                      "tag": "lark_md"
                                  },
                                  "url": "https://www.example.com",
                                  "type": "default",
                                  "value": {}
                              }
                          ],
                          "tag": "action"
                      }
                  ],
                  "header": {
                      "title": {
                          "content": "今日旅游推荐",
                          "tag": "plain_text"
                      }
                  }
              }
```

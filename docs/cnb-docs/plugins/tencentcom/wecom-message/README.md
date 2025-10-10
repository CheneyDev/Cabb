# wecom-message

企业微信群里发消息的插件，通过机器人，可以向固定的群中推送消息。

## 输入

- `robot`: `String` 必填，值为企业微信机器人Webhook地址，或Webhook地址上的key参数
- `msgType`: `String` 可选，消息类型。`text|markdown|image|file|message`，默认值：`markdown`
- `content`: `String` 可选，消息内容。
`msyType` 为 `text`，`markdown` 时有效，`content` 和 `fromFile` 二选一，优先使用 `content`
- `fromFile`: `String` 可选，消息内容所在文件相对路径。
`msyType` 为 `text`，`markdown` 时有效，`content` 和 `fromFile` 二选一，优先使用 `content`
- `filePath`: `String` 可选，指定要发送的文件相对路径。大小应小于等于 20 MB
- `layouts`: `String` 可选，消息内容 或 消息内容文件路径，`JSON` 格式。 `msyType` 为 `message` 时需要填入
- `chatId`: `String` 可选，企业微信会话 id，支持传多个，用|分隔
- `visibleToUser` : `String` 可选，指定此消息可见的群成员，用 | 或者 , 分隔
- `postId`: `String` 可选，有且只有chatid指定了一个小黑板的时候生效
- `lastMsg`: `String` 可选，企业微信发送消息的长度有限制（2048 字节），如果超过指定长度，用于在结尾添加一行提示
- `attachments`: `String` 可选，json 字符串（对象格式）。定义消息附加信息。详见下方 `attachments` 的详细说明。
- `mentioned_list`: `String` 可选，指定此消息提及的群成员，用 | 或者 , 分隔。`msyType` 为 `text` 时有效。

消息长度、文件大小的限制请参考[企业微信API](https://developer.work.weixin.qq.com/document/path/90236)

### attachments

定义消息附加信息，目前仅支持 button 类型。`msyType` 为 `markdown` 格式时支持传入 `attachments`。

- `callback_id`: `String` 必填，attachments 对应的回调 id，企业微信在回调时会透传该值
- `actions`: `String` 必填，对象数组格式，attachments 动作， 一个 attachment 中最多支持 20 个 action。
  - `type`: `String` 必填，action 类型，目前仅支持按钮
  - `name`: `String` 必填，action 名字，企业微信回调时会透传该值，最长不超过64字节。
    为了区分不同的按钮，建议开发者自行保证 name 的唯一性
  - `text`: `String` 必填，需要展示的文本，最长不超过128字节
  - `border_color`: `String` 可选，按钮边框颜色
  - `text_color`: `String` 可选，按钮文字颜色
  - `value`: `String` 必填，action 的值，企业微信回调时会透传该值，最长不超过 128 字节
  - `replace_text`: `String` 必填，点击按钮后替换的文本，最长不超过 128 字节

## 示例

### PR 消息通知到群

#### 在 云原生构建 上

```yaml
# .cnb.yml
main:
  pull_request:
  - stages:
    - name: send message
      image: tencentcom/wecom-message
      settings:
        robot: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
        msgType: markdown
        # 消息内容所在文件相对路径，fromFile 和 content 二选一，优先使用 content
        # fromFile: ./message.txt
        content: |
          走查代码啦
          　
          ${CNB_PULL_REQUEST_TITLE}
          [${CNB_EVENT_URL}](${CNB_EVENT_URL})
          
          from ${CNB_BUILD_USER}
```

### 自定义按钮

```yaml
main:
  push:
  - stages:
    - name: send message
      image: tencentcom/wecom-message
      settings:
        robot: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
        msgType: markdown
        # 自定义按钮
        attachments: |
          [
              {
                  "callback_id": "button_two_row",
                  "actions": [
                      {
                          "name": "button_1",
                          "text": "S",
                          "type": "button",
                          "value": "S",
                          "replace_text": "你已选择S",
                          "border_color": "2EAB49",
                          "text_color": "2EAB49"
                      },
                      {
                          "name": "button_2",
                          "text": "M",
                          "type": "button",
                          "value": "M",
                          "replace_text": "你已选择M",
                          "border_color": "2EAB49",
                          "text_color": "2EAB49"
                      },
                      {
                          "name": "button_3",
                          "text": "L",
                          "type": "button",
                          "value": "L",
                          "replace_text": "你已选择L",
                          "border_color": "2EAB49",
                          "text_color": "2EAB49"
                      },
                      {
                          "name": "button_4",
                          "text": "不确定",
                          "type": "button",
                          "value": "不确定",
                          "replace_text": "你已选择不确定",
                          "border_color": "2EAB49",
                          "text_color": "2EAB49"
                      },
                      {
                          "name": "button_5",
                          "text": "不参加",
                          "type": "button",
                          "value": "不参加",
                          "replace_text": "你已选择不参加",
                          "border_color": "2EAB49",
                          "text_color": "2EAB49"
                      }
                  ]
              }
          ]
```

### layouts

```yaml
main:
  push:
  - stages:
    - name: send message
      image: tencentcom/wecom-message
      settings:
        robot: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
        msgType: message
        layouts: layouts.json
```

`layouts.json`

```json
[
    {
        "type": "column_layout",
        "components": [
            {
                "type": "plain_text",
                "text": "Checkbox",
                "style": "title"
            },
            {
                "type": "checkbox",
                "key": "checkbox_3",
                "options": [
                    {
                        "id": "1",
                        "text": "选项名称1",
                        "checked": true
                    },
                    {
                        "id": "2",
                        "text": "选项名称2",
                        "checked": false
                    }
                ]
            }
        ]
    }
]
```

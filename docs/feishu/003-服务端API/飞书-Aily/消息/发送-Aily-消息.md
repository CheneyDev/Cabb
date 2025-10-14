---
title: "发送 Aily 消息"
source_url: https://open.feishu.cn/document/aily-v1/aily_session-aily_message/create
last_remote_update: 2025-07-10
last_remote_update_timestamp: 1752154928000
---
最后更新于 2025-07-10

# 发送 Aily 消息

该 API 用于向某个飞书 Aily 应用发送一条消息（Message）；每个消息从属于一个活跃的会话（Session）。
**注意事项**：更多信息及示例代码，可参考 [Aily OpenAPI 接入与接口说明](https://bytedance.larkoffice.com/wiki/UTU6wVTVGigefykjO1acAOOvnNc)。

## 实体概念说明

- **会话**（Session）：管理用户与 Aily 助手之间的交互会话；每次会话记录了用户发送给 Aily 助手的消息以及 Aily 助手的响应。
- **消息**（Message）：消息可以包含文本、表格、图片等多种类型的内容。
- **运行**（Run）：Aily 助手基于会话内消息进行意图判定、调用匹配的技能，并返回技能执行后的结果消息。

## 请求
名称 | 值
---|---
HTTP URL | https://open.feishu.cn/open-apis/aily/v1/sessions/:aily_session_id/messages
HTTP Method | POST
接口频率限制 | [1000 次/分钟、50 次/秒](https://open.feishu.cn/document/ukTMukTMukTM/uUzN04SN3QjL1cDN)
支持的应用类型 | 自建应用
权限要求  
            **调用该 API 所需的权限。开启其中任意一项权限即可调用** | 发送消息(aily:message:write)

### 请求头

名称 | 类型 | 必填 | 描述
--- | --- | --- | ---
Authorization | string | 是 | `tenant_access_token`  
或  
`user_access_token`  
**值格式**："Bearer `access_token`"  
**示例值**："Bearer u-7f1bcd13fc57d46bac21793a18e560"  
[了解更多：如何选择与获取 access token](https://open.feishu.cn/document/uAjLw4CM/ugTN1YjL4UTN24CO1UjN/trouble-shooting/how-to-choose-which-type-of-token-to-use)
Content-Type | string | 是 | **固定值**："application/json; charset=utf-8"

### 路径参数

名称 | 类型 | 描述
--- | --- | ---
aily_session_id | string | 会话 ID；参考 [创建会话](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/aily-v1/aily_session/create) 接口  
**示例值**："session_4dfunz7sp1g8m"  
**数据校验规则**：  
- 长度范围：`9` ～ `32` 字符  
- 正则校验：`session_[0-9a-hjkmnp-z]{1,24}`

### 请求体

名称 | 类型 | 必填 | 描述
--- | --- | --- | ---
idempotent_id | string | 是 | 幂等 ID（如使用 UUID 生成器或时间戳），同一会话下相同的幂等 ID 视为同一个消息（72h）  
**示例值**："idempotent_id_1"  
**数据校验规则**：  
- 长度范围：`0` ～ `64` 字符
content_type | string | 是 | 消息的类型，包括 `MDX` | `TEXT` 等  
- `MDX` 能够表达富文本信息结构，可参考 [Aily 消息节点 markdown 语法](https://bytedance.larkoffice.com/wiki/ZlHYw8jJci3o4dkkEZZcybM0nah)  
- `TEXT` 作为纯文本进行处理  
**示例值**："MDX"  
**可选值有**：  
- MDX：MDX  
- TEXT：TEXT  
- CLIP：GUI 卡片  
- SmartCard：SmartCard  
- JSON：JSON
content | string | 是 | 消息内容  
**示例值**："你好"  
**数据校验规则**：  
- 长度范围：`0` ～ `61440` 字符
file_ids | string\[\] | 否 | 消息中包含的文件 ID 列表  
**示例值**：["file_4d9nu1ev3a2rq"]  
**数据校验规则**：  
- 长度范围：`0` ～ `32`
quote_message_id | string | 否 | 引用的消息 ID  
**示例值**："message_4de9bpg70qskh"  
**数据校验规则**：  
- 长度范围：`0` ～ `32` 字符  
- 正则校验：`message_[0-9a-hjkmnp-z]{1,24}`
mentions | aily_mention\[\] | 否 | 被@的实体  
**数据校验规则**：  
- 长度范围：`0` ～ `32`
entity_id | string | 否 | 实体 ID  
**示例值**："ou_5ad573a6411d72b8305fda3a9c15c70e"  
**数据校验规则**：  
- 长度范围：`0` ～ `64` 字符
identity_provider | string | 否 | 身份提供者  
**示例值**："FEISHU"  
**可选值有**：  
- AILY：Aily 账号体系  
- FEISHU：飞书账号体系
key | string | 否 | 被@实体在消息体中的占位符  
**示例值**："@_user_1"  
**数据校验规则**：  
- 长度范围：`0` ～ `32` 字符
name | string | 否 | 被@实体的名称  
**示例值**："张三"  
**数据校验规则**：  
- 长度范围：`0` ～ `32` 字符
aily_id | string | 否 | 飞书智能伙伴创建平台账号体系下的 ID  
**示例值**："1794840334557292"  
**数据校验规则**：  
- 长度范围：`0` ～ `20` 字符

### 请求体示例
```json
{
    "idempotent_id": "idempotent_id_1",
    "content_type": "MDX",
    "content": "你好",
    "file_ids": [
        "file_4d9nu1ev3a2rq"
    ],
    "quote_message_id": "message_4de9bpg70qskh",
    "mentions": [
        {
            "entity_id": "ou_5ad573a6411d72b8305fda3a9c15c70e",
            "identity_provider": "FEISHU",
            "key": "@_user_1",
            "name": "张三",
            "aily_id": "1794840334557292"
        }
    ]
}
```

## 响应

### 响应体

名称 | 类型 | 描述
--- | --- | ---
code | int | 错误码，非 0 表示失败
msg | string | 错误描述
data | \- | \-
message | aily_message | 消息信息
id | string | 消息 ID
session_id | string | 会话 ID
run_id | string | 运行 ID
content_type | string | 消息内容类型  
**可选值有**：  
- MDX：MDX  
- TEXT：TEXT  
- CLIP：GUI 卡片  
- SmartCard：SmartCard  
- JSON：JSON
content | string | 消息内容
files | aily_message_file\[\] | 消息中包含的文件
id | string | 文件 ID
mime_type | string | 文件类型，参见[MIME 类型（IANA 媒体类型）](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Basics_of_HTTP/MIME_types)
file_name | string | 文件名
metadata | string | 其他透传信息
created_at | string | 文件的创建时间，毫秒时间戳
preview_url | aily_message_file_preview | 文件预览链接
url | string | 文件的 URL
expired_at | string | url 过期时间，秒时间戳
quote_message_id | string | 引用的消息 ID
sender | aily_sender | 发送者
entity_id | string | 实体 ID
identity_provider | string | 身份提供者  
**可选值有**：  
- AILY：Aily 账号体系  
- FEISHU：飞书账号体系
sender_type | string | 发送人类型  
**可选值有**：  
- USER：用户  
- ASSISTANT：应用
aily_id | string | 飞书智能伙伴创建平台账号体系下的 ID
mentions | aily_mention\[\] | 被@的实体
entity_id | string | 实体 ID
identity_provider | string | 身份提供者  
**可选值有**：  
- AILY：Aily 账号体系  
- FEISHU：飞书账号体系
key | string | 被@实体在消息体中的占位符
name | string | 被@实体的名称
aily_id | string | 飞书智能伙伴创建平台账号体系下的 ID
plain_text | string | 消息体的纯文本表达
created_at | string | 消息的创建时间，毫秒时间戳
status | string | 状态  
**可选值有**：  
- IN_PROGRESS：生成中  
- COMPLETED：已完成

### 响应体示例
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "message": {
            "id": "message_4df45f2xknvcc",
            "session_id": "session_4dfunz7sp1g8m",
            "run_id": "run_4dfrxvctjqzzj",
            "content_type": "MDX",
            "content": "你好",
            "files": [
                {
                    "id": "file_4d9nu1ev3a2rq",
                    "mime_type": "image/png",
                    "file_name": "发票.png",
                    "metadata": "{}",
                    "created_at": "1711975665710",
                    "preview_url": {
                        "url": "http://path_to_file",
                        "expired_at": "1719413169"
                    }
                }
            ],
            "quote_message_id": "message_4de9bpg70qskh",
            "sender": {
                "entity_id": "ou_5ad573a6411d72b8305fda3a9c15c70e",
                "identity_provider": "FEISHU",
                "sender_type": "USER",
                "aily_id": "1794840334557292"
            },
            "mentions": [
                {
                    "entity_id": "ou_5ad573a6411d72b8305fda3a9c15c70e",
                    "identity_provider": "FEISHU",
                    "key": "@_user_1",
                    "name": "张三",
                    "aily_id": "1794840334557292"
                }
            ],
            "plain_text": "你好",
            "created_at": "1711975665710",
            "status": "COMPLETED"
        }
    }
}
```

### 错误码

HTTP状态码 | 错误码 | 描述 | 排查建议
--- | --- | --- | ---
400 | 2700001 | param is invalid | 参数错误

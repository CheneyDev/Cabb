---
title: "更新用户 ID"
source_url: https://open.feishu.cn/document/contact-v3/user/update_user_id
last_remote_update: 2024-07-05
last_remote_update_timestamp: 1720166743000
---
最后更新于 2024-07-05

# 更新用户ID

调用该接口更新用户的 user_id。

## 注意事项

更新后的用户 user_id 需要保证在当前租户内未被占用。

## 请求
名称 | 值
---|---
HTTP URL | https://open.feishu.cn/open-apis/contact/v3/users/:user_id/update_user_id
HTTP Method | PATCH
接口频率限制 | [20 次/秒](https://open.feishu.cn/document/ukTMukTMukTM/uUzN04SN3QjL1cDN)
支持的应用类型 | 自建应用
权限要求  
            **调用该 API 所需的权限。开启其中任意一项权限即可调用** | 更新用户 ID(contact:contact:update_user_id)
字段权限要求 | **注意事项**：该接口返回体中存在下列敏感字段，仅当开启对应的权限后才会返回；如果无需获取这些字段，则不建议申请  
        获取用户 user ID(contact:user.employee_id:readonly)

### 请求头

名称 | 类型 | 必填 | 描述
--- | --- | --- | ---
Authorization | string | 是 | `tenant_access_token`  
**值格式**："Bearer `access_token`"  
**示例值**："Bearer t-7f1bcd13fc57d46bac21793a18e560"  
[了解更多：如何选择与获取 access token](https://open.feishu.cn/document/uAjLw4CM/ugTN1YjL4UTN24CO1UjN/trouble-shooting/how-to-choose-which-type-of-token-to-use)
Content-Type | string | 是 | **固定值**："application/json; charset=utf-8"

### 路径参数

名称 | 类型 | 描述
--- | --- | ---
user_id | string | 用户 ID，ID 类型与查询参数 user_id_type 的取值保持一致。  
**示例值**："ou-938e3e4fdc5e1993bee01250076f0cc2"

### 查询参数

名称 | 类型 | 必填 | 描述
--- | --- | --- | ---
user_id_type | string | 否 | 用户 ID 类型  
**示例值**：open_id  
**可选值有**：  
- open_id：标识一个用户在某个应用中的身份。同一个用户在不同应用中的 Open ID 不同。[了解更多：如何获取 Open ID](https://open.feishu.cn/document/uAjLw4CM/ugTN1YjL4UTN24CO1UjN/trouble-shooting/how-to-obtain-openid)  
- union_id：标识一个用户在某个应用开发商下的身份。同一用户在同一开发商下的应用中的 Union ID 是相同的，在不同开发商下的应用中的 Union ID 是不同的。通过 Union ID，应用开发商可以把同个用户在多个应用中的身份关联起来。[了解更多：如何获取 Union ID？](https://open.feishu.cn/document/uAjLw4CM/ugTN1YjL4UTN24CO1UjN/trouble-shooting/how-to-obtain-union-id)  
- user_id：标识一个用户在某个租户内的身份。同一个用户在租户 A 和租户 B 内的 User ID 是不同的。在同一个租户内，一个用户的 User ID 在所有应用（包括商店应用）中都保持一致。User ID 主要用于在不同的应用间打通用户数据。[了解更多：如何获取 User ID？](https://open.feishu.cn/document/uAjLw4CM/ugTN1YjL4UTN24CO1UjN/trouble-shooting/how-to-obtain-user-id)  
**默认值**：`open_id`  
**当值为 `user_id`，字段权限要求**：  
获取用户 user ID(contact:user.employee_id:readonly)

### 请求体

名称 | 类型 | 必填 | 描述
--- | --- | --- | ---
new_user_id | string | 是 | 自定义新的用户 user_id。长度不能超过 64 字符。  
**示例值**："3e3cf96b"

### 请求体示例
```json
{
    "new_user_id": "3e3cf96b"
}
```

## 响应

### 响应体

名称 | 类型 | 描述
--- | --- | ---
code | int | 错误码，非 0 表示失败
msg | string | 错误描述
data | \- | \-

### 响应体示例
```json
{
    "code": 0,
    "msg": "success",
    "data": {}
}
```

### 错误码

HTTP状态码 | 错误码 | 描述 | 排查建议
--- | --- | --- | ---
400 | 40001 | Invalid params | 参数错误。请参考接口文档的参数描述，检查输入参数是否填写有误。
400 | 41009 | no email or mobile error | 无邮箱或者手机号，请添加后重试。
400 | 41011 | userID already exist error | 用户的 user_id 已存在，你需要更换 ID 并重试。
400 | 41013 | exceed userID update limit error | 超过用户 ID 更新次数限制。
400 | 41050 | no user authority error | 无用户权限。当前操作的用户需在应用的通讯录权限范围内。通讯录权限范围的介绍与设置方式，参见[权限范围资源介绍](https://open.feishu.cn/document/ukTMukTMukTM/uETNz4SM1MjLxUzM/v3/guides/scope_authority)。
400 | 41054 | new user id is required  error | 更新后的 user_id 错误。请更换 user_id 值后重试。

更多错误码信息，参见[通用错误码](https://open.feishu.cn/document/ukTMukTMukTM/ugjM14COyUjL4ITN)。

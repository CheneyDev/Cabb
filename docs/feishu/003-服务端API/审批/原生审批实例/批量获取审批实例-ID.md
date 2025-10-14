---
title: "批量获取审批实例 ID"
source_url: https://open.feishu.cn/document/server-docs/approval-v4/instance/list
last_remote_update: 2025-08-29
last_remote_update_timestamp: 1756436658000
---
最后更新于 2025-08-29

# 批量获取审批实例 ID

根据审批定义的 approval_code 批量获取审批实例的 instance_code，用于拉取企业下某个审批定义的全部审批实例。默认以审批创建时间先后顺序排列。

## 请求
名称 | 值
---|---
HTTP URL | https://open.feishu.cn/open-apis/approval/v4/instances
HTTP Method | GET
接口频率限制 | [100 次/分钟](https://open.feishu.cn/document/ukTMukTMukTM/uUzN04SN3QjL1cDN)
支持的应用类型 | 自建应用、商店应用
权限要求  
            **调用该 API 所需的权限。开启其中任意一项权限即可调用**  
            开启任一权限即可 | 查看、创建、更新、删除审批应用相关信息(approval:approval)  
            访问审批应用(approval:approval:readonly)  
            查看、创建、更新、删除原生审批实例相关信息(approval:instance)

### 请求头

名称 | 类型 | 必填 | 描述
--- | --- | --- | ---
Authorization | string | 是 | `tenant_access_token`  
**值格式**："Bearer `access_token`"  
**示例值**："Bearer t-7f1bcd13fc57d46bac21793a18e560"  
[了解更多：如何选择与获取 access token](https://open.feishu.cn/document/uAjLw4CM/ugTN1YjL4UTN24CO1UjN/trouble-shooting/how-to-choose-which-type-of-token-to-use)

### 查询参数

名称 | 类型 | 必填 | 描述
--- | --- | --- | ---
page_size | int | 否 | 分页大小，用于指定一次请求所返回的数据量上限。  
**示例值**：100  
**默认值**：`100`  
**数据校验规则**：  
- 取值范围：`1` ～ `100`
page_token | string | 否 | 分页标记，第一次请求不填，表示从头开始遍历；分页查询结果还有更多项时会同时返回新的 page_token，下次遍历可采用该 page_token 获取查询结果  
**示例值**：nF1ZXJ5VGhlbkZldGNoCgAAAAAA6PZwFmUzSldvTC1yU
approval_code | string | 是 | 审批定义 Code。获取方式：  
- 调用[创建审批定义](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/approval-v4/approval/create)接口后，从响应参数 approval_code 获取。  
- 登录审批管理后台，在指定审批定义的 URL 中获取，具体操作参见[什么是 Approval Code](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/approval-v4/approval/overview-of-approval-resources#8151e0ae)。  
**示例值**：7C468A54-8745-2245-9675-08B7C63E7A85
start_time | string | 是 | 审批实例创建时间的开始区间，毫秒时间戳。  
**说明**：start_time 与 end_time 组成时间区间查询条件，接口会返回在该时间区间内创建的审批实例数据。  
**注意事项**：单次查询时间范围不要超过10小时  
**示例值**：1567690398020
end_time | string | 是 | 审批实例创建时间的结束区间，毫秒时间戳。  
**说明**：start_time 与 end_time 组成时间区间查询条件，接口会返回在该时间区间内创建的审批实例的 Code。  
**注意事项**：单次查询时间范围不要超过10小时  
**示例值**：1567690398020

## 响应

### 响应体

名称 | 类型 | 描述
--- | --- | ---
code | int | 错误码，非 0 表示失败
msg | string | 错误描述
data | \- | \-
instance_code_list | string\[\] | 审批实例 Code 列表，默认以审批实例创建时间的先后顺序排列。
page_token | string | 分页标记，当 has_more 为 true 时，会同时返回新的 page_token，否则不返回 page_token
has_more | boolean | 是否还有更多项

### 响应体示例
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "instance_code_list": [
            "357C21A0-2069-4F6B-955F-1DFBE6710C51"
        ],
        "page_token": "nF1ZXJ5VGhlbkZldGNoCgAAAAAA6PZwFmUzSldvTC1yU",
        "has_more": false
    }
}
```

### 错误码

HTTP状态码 | 错误码 | 描述 | 排查建议
--- | --- | --- | ---
400 | 1390001 | param is invalid | 参数错误。排查方案：  
- 根据接口文档的参数说明，检查请求时传入的参数是否正确。  
- 如果传入的有表单参数（form），则需要检查该参数内传入的表单控件数据是否正确。如果报错信息内包含控件 ID（如 `控件= widget17261088448220001`），可以调用[查看指定审批定义](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/approval-v4/approval/get)或者[获取单个审批实例详情](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/approval-v4/instance/get)接口，获取响应参数 form 值，检索有问题的控件 ID，然后检查该控件的配置是否正确。
400 | 1390002 | approval code not found | 找不到审批定义 Code，检查传入的审批定义 Code 是否正确。  
审批定义 Code 获取方式：  
- 调用[创建审批定义](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/approval-v4/approval/create)接口后，从响应参数 approval_code 获取。  
- 登录审批管理后台，在指定审批定义的 URL 中获取，具体操作参见[什么是 Approval Code](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/approval-v4/approval/overview-of-approval-resources#8151e0ae)。
400 | 1395001 | There have been some errors. Please try again later | 服务出现错误。排查方案：  
1. 参考接口文档的参数说明，检查请求时传入的参数是否正确。如果传入的有表单参数（form），则需要检查传入的表单控件数据是否正确。  
2. 降低请求频率，并重试。如果重试仍然报错，请联系[技术支持](https://applink.feishu.cn/TLJpeNdW)。

更多错误码信息，参见[通用错误码](https://open.feishu.cn/document/ukTMukTMukTM/ugjM14COyUjL4ITN)。

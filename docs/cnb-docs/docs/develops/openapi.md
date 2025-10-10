---
title: Open API
permalink: https://docs.cnb.cool/zh/openapi.html
summary: 该文档介绍了Open API相关信息，API服务请求地址为`api.cnb.cool`，文档地址也是此处，调用时在Header头中用`Authorization`字段检验（格式为`Bearer ${token}`，${token}为访问令牌），还给出了curl请求示例（以获取用户组信息为例）及返回示例 。
---
## API 服务地址

服务请求地址为 `api.cnb.cool`。

## API 文档

API 文档地址 [api.cnb.cool](https://api.cnb.cool)。

## 调用方式

### Header 头信息

- `Authorization`  进行检验，格式为: `Bearer ${token}`

其中`${token}`为[访问令牌](../guide/access-token.md)

### 请求示例

curl 请求示例

``` shell
curl  -X "GET" \
      -H "accept: application/json" \
      -H "Authorization: Bearer 1Z00000000000000000000000vA" \
  "api.cnb.cool/user/groups?page=1&page_size=10"
```

返回示例

```json
[
  {
    "id": 1816756487609032700,
    "name": "test",
    "remark": "测试组织",
    "description": "",
    "site": "",
    "email": "",
    "freeze": false,
    "wechat_mp": "hello-world",
    "created_at": "2024-07-26T08:44:35Z",
    "updated_at": "2024-08-13T07:32:13Z",
    "follow_count": 0,
    "member_count": 4,
    "all_member_count": 4,
    "sub_group_count": 5,
    "sub_repo_count": 7,
    "sub_mission_count": 1,
    "all_sub_group_count": 13,
    "all_sub_repo_count": 12,
    "all_sub_mission_count": 1,
    "has_sub_group": true,
    "path": "test",
    "access_role": "Owner"
  }
]
```

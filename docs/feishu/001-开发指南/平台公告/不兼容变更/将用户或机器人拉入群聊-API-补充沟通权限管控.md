---
title: "将用户或机器人拉入群聊 API 补充沟通权限管控"
source_url: https://open.feishu.cn/document/platform-notices/breaking-change/additional-communication-controls-for-add-to-group-api
last_remote_update: 2025-03-10
last_remote_update_timestamp: 1741608117000
---
最后更新于 2025-03-10

# 将用户或机器人拉入群聊API补充沟通权限管控
## 变更说明

为保障产品功能体验的一致性，我们将于 **2025年3月31日** 对[将用户或机器人拉入群聊 API](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/chat-members/create)补充沟通权限管控，包括沟通协作权限、对外沟通权限，以及屏蔽状态校验。管控生效后，以用户身份调用 API 拉企业内用户进群时将额外受到沟通协作权限限制，拉外部用户进群时将额外要求必须具备对外沟通权限，被外部联系人屏蔽时将无法拉其进群。

在调用 API 拉人进群时，可以通过入参 `succeed_type` 选择出现不可用 ID 后的处理方式，对于沟通权限校验未通过的情况，也会有不同结果返回：
- `succeed_type = 0`：存在已离职 ID 时，会将其他可用 ID 拉入群聊，返回拉群成功，离职以外的不可用 ID 会导致拉群失败。沟通权限校验未通过将返回 `232024` 错误码。
- `succeed_type = 1`：将可用的 ID 全部拉入群聊，返回拉群成功，并展示剩余不可用的 ID 及原因。沟通权限校验未通过的 ID 将在返回值的 `invalid_id_list` 字段中展示。
- `succeed_type = 2`：任一不可用的 ID 会导致拉群失败。沟通权限校验未通过将返回 `232024` 错误码。

> 是否跟随客户端版本：不跟版 <br>
> 预计生效时间：2025-3-31 <br>

## 潜在影响

以用户身份调用[将用户或机器人拉入群聊 API](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/chat-members/create) 时，若对内部用户配置并命中了沟通协作权限拦截，在历史版本可以成功拉进群，增加管控之后将被拦截；若没有外部沟通权限，在历史版本可以成功将外部用户拉进群，增加管控之后也将被拦截；若被外部联系人屏蔽，在历史版本可以成功将其拉进群，管控后无法拉进群。

## 解决方案

开发者可以通过通知企业管理员，在管理后台中为指定用户解除沟通协作权限限制，或者为其开通对外沟通权限的方式避免被管控。被外部联系人屏蔽导致的拉群失败需要联系对应用户解除屏蔽。

若你未能及时确认并调整，管控生效后，可能会导致相关场景受损。开放平台预计在 **2025年3月31日** 完成升级，请于此前确认以上信息，并根据情况做好相应适配。

如需适配协助，请联系开放平台技术支持。

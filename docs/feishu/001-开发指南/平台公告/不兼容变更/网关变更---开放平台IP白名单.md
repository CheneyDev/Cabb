---
title: "网关变更 - 开放平台IP白名单"
source_url: https://open.feishu.cn/document/faq/breaking-change/gateway-change-ip-allowlist
last_remote_update: 2022-05-27
last_remote_update_timestamp: 1653636213000
---
最后更新于 2022-05-27

# 网关变更 - 开放平台IP白名单
### 变更事项
我们正在对开放平台的基础设施进行升级，目前升级已经进行到最后阶段。<br>
在新版本的灰度过程中，平台监测到部分用户的应用在调用 OpenAPI 的时候，有部分请求IP来自于应用已配置的IP白名单之外。

是否跟版：否<br>
预计生效版本：- <br>
预计生效时间：2021-07-27<br>
### 潜在影响
平台升级结束后，这部分 IP 白名单之外的请求将会被网关安全策略拦截并返回错误。

### 解决方案
为了不影响你的正常业务，请及时检查并修改IP白名单配置。<br>
修改方式：
1.  如应用需要继续保留 IP 白名单的安全配置，你可将所有允许访问开放平台 OpenAPI 的服务器出口 IP 填写到 IP 白名单中
2.  如需暂时关闭 IP 白名单配置，将应用 IP 白名单配置项清空即可

关于如何获取服务器出口 IP，请参考[服务端错误码说明](https://open.feishu.cn/document/ukTMukTMukTM/ugjM14COyUjL4ITN)文档中关于99991401错误码的排查建议。

<br> 如需适配协助，请洽客服支持

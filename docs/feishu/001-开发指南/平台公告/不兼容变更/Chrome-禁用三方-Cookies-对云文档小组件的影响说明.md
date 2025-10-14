---
title: "Chrome 禁用三方 Cookies 对云文档小组件的影响说明"
source_url: https://open.feishu.cn/document/platform-notices/breaking-change/chrome-disabling-third-party-cookies-on-document-widgets
last_remote_update: 2024-01-24
last_remote_update_timestamp: 1706078833000
---
最后更新于 2024-01-24

# 变更背景

Chrome 浏览器计划 2024 年全面禁用三方 Cookies，将影响使用三方 Cookies 进行用户身份验证的云文档小组件。

# 变更节奏

Chrome 计划从 2024 年第一季度开始对 1% 的用户禁用第三方 Cookies，并从 2024 年第三季度开始逐步将禁用范围扩大到 100%。

# 自查指南

我们建议开发者提前自查云文档小组件的业务功能是否受影响，自查步骤如下：

1. 浏览器打开 chrome://settings/cookies, 选择阻止第三方 Cookies

2. 确认浏览器在阻止第三方 Cookies 的状态下，对小组件业务进行走查，可以通过：
    1. 自测小组件在阻止第三方 Cookies 状态下，功能是否如预期运行
    1. 观察前端控制台 Network 是否有请求鉴权失败

# 建议方案

若通过对比禁用 Cookies 前后情况，确认是由于禁用 Cookies 后导致的请求失败，我们建议你将原来携带在 Cookies 中的鉴权信息，改造为携带在 request header 中。

请求失败的原因是在未禁用第三方 Cookie 的情况下，请求中包含了第三方 Cookies 信息，业务系统利用 Cookies 中的第三方登录信息进行鉴权操作。

当禁用第三方 Cookies 后，Cookies 中只能携带同域内的信息，这可能导致鉴权信息不足，从而请求失败。

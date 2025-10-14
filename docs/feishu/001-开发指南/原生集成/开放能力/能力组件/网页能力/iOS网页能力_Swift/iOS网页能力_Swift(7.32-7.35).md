---
title: "iOS网页能力_Swift(7.32-7.35)"
source_url: https://open.feishu.cn/document/native-integration/open-capability/capability-components/web-capabilities/ios-web-capability/ios-web-capabil
last_remote_update: 2025-04-21
last_remote_update_timestamp: 1745215563000
---
最后更新于 2025-04-21

# iOS网页能力_Swift(7.32-7.35)
组件类型 | 最低飞书版本 | 生效基线      |
| ---- | ------ | --------- |
| 能力组件 | 7.32  | 7.32.0 及以上 |

## 组件说明
通过`LKWebViewApiExternal`组件，开发者可以实现在飞书的 H5 容器中设置自定义的 cookie，或追加自定义的 User Agent 信息。这样可以为开发者提供更多的灵活性和个性化定制的能力，以满足不同的业务需求和用户体验要求。

## API 列表

### protocol KAWebViewApiDelegate

#### **func** setCookie(cookie: HTTPCookie)

##### 接口说明

| 兼容飞书版本 | 说明               | 生效版本  |
| ------ | ---------------- | ----- |
| 7.32及以上   | 飞书 H5容器设置 Cookie | 1.0.0 |

##### 参数说明

| 名称     | 类型         | 描述         |
| ------ | ---------- | ---------- |
| cookie | HTTPCookie | 设置的Cookie值 |

#### **func** appendUserAgent(customUA: String)

##### 接口说明

| 兼容飞书版本 | 说明                  | 生效版本  |
| ------ | ------------------- | ----- |
| 7.32及以上  | 向 userAgent 中注入额外字段 | 1.0.0 |

##### 参数说明

| 名称       | 类型     | 描述                  |
| -------- | ------ | ------------------- |
| customUA | string | 要注入的 userAgent 相关字段 |
**注意事项**：建议 application 启动直接注入，确保一次性注入完全，多次调用注入不同的值会造成覆盖。

格式要求：新注入的 product 请不要与已有 product 冲突， product 含义请参考 https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent

## 接入说明

通过原生集成工具，选择 2.0.0 及以上基线，依赖 **NativeAppPublicKit** 组件，podspec 增加以下内容：
```javascript
s.dependency 'NativeAppPublicKit', '{你需要依赖的版本号}'
```
调用：
```javascript
[[[NativeAppConnectManager shared] getNativeAppManager] setCookie: ];
```

## 更新日志

版本号 | Release Note
--- | ---
1.0.0 | 新功能：飞书的 H5 容器中设置自定义的 cookie，并追加自定义的 User Agent 信息
1.1.1 | 改进：  
- 组件名更改为 LKWebViewApiExternal  
- 协议名改为 KAWebViewApiDelegate  
- 更改为 KAAPI 调用

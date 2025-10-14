---
title: "getStorageSync"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/cache/getstoragesync
last_remote_update: 2025-01-21
last_remote_update_timestamp: 1737434391000
---
最后更新于 2025-01-21

# getStorageSync

调用 getStorageSync(string key) 获取本地缓存数据。

## 支持说明

该接口仅支持小程序调用，对应的客户端版本支持情况如下所示。

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | **✓** | **✓** | **✓** | V7.35.0+ | 预览
网页应用 | **x** | **x** | **x** | **x** | /

## 输入

名称 | 数据类型 | 是否必填 | 默认值 | 描述
--- | --- | --- | --- | ---
key | string | 是 | \- | 键名。最小长度为 1 字符。示例值：stateKey

## 输出

名称 | 数据类型 | 描述
--- | --- | ---
data | [本地缓存数据类型](https://open.feishu.cn/document/uYjL24iN/uMTOz4yM5MjLzkzM) | 键名对应的数据。

## 示例代码

调用示例：

```js
try {
    let result = tt.getStorageSync("name");
    console.log(`getStorageSync success: ${JSON.stringify(result)}`);
} catch (error) {
    console.log(`getStorageSync fail: ${JSON.stringify(error)}`);
}
```

返回值示例：
```json
"DemoName"
```

---
title: "onMemoryWarning"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/performance/onmemorywarning
last_remote_update: 2024-01-16
last_remote_update_timestamp: 1705404230000
---
最后更新于 2024-01-16

# onMemoryWarning(function callback)

监听内存不足的告警事件。当手机内存占用过高时，触发回调函数。该事件不会杀掉小程序, 建议开发者可以在接受到告警后释放不必要的资源。 在Android平台下有告警等级划分。

## 支持说明

应用能力 | Android | iOS | PC | 预览效果
--- | --- | --- | --- | ---
小程序 | V5.4.0+ | V5.4.0+ | **X** | 预览
网页应用 | V7.10.0+ | V7.10.0+ | **X** | **/**

## 输入

名称 | 数据类型 | 必填 | 默认值 | 描述
--- | --- | --- | --- | ---
callback | function | 是 |  | 该事件的回调函数

## 输出
回调函数返回对象的属性：

名称 | 数据类型 | 描述
--- | --- | ---
level | number | 对应系统内存告警等级宏（Level）定义，数值越高，告警等级越高。  
**可选值**：  
- `5`：TRIM_MEMORY_RUNNING_MODERATE  
- `10`：TRIM_MEMORY_RUNNING_LOW  
- `15`：TRIM_MEMORY_RUNNING_CRITICAL  
**Notice**：仅 Android 端返回该字段

## 示例代码

<div style="display: flex">

</div> 

```js
tt.onMemoryWarning(function () {
  console.log("onMemoryWarning");
});
```

回调函数返回对象示例（仅Android）：

```json
{"level":"5"}

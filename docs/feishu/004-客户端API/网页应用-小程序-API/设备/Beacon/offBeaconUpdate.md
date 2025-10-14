---
title: "offBeaconUpdate"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/ibeacon/offbeaconupdate
last_remote_update: 2025-01-21
last_remote_update_timestamp: 1737432742000
---
最后更新于 2025-01-21

# offBeaconUpdate(function callback)

取消监听 Beacon 设备更新事件

## 支持说明

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | V4.6.0+ | V4.6.0+ | **X** | V7.35.0+ | 预览
网页应用 | V4.6.0+ | V4.6.0+ | **X** | V7.35.0+ | 预览

## 输入

名称 | 数据类型 | 必填 | 默认值 | 描述
--- | --- | --- | --- | ---
callback | function | 是 |  | 该事件的回调函数

## 输出
无

## 示例代码

```js
const callback = (res) => {
	console.log(res);
};
tt.onBeaconUpdate(callback);
tt.offBeaconUpdate(callback);
```

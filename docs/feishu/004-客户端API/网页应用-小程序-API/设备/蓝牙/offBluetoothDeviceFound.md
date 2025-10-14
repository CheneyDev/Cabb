---
title: "offBluetoothDeviceFound"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/bluetooth/offbluetoothdevicefound
last_remote_update: 2025-01-21
last_remote_update_timestamp: 1737432463000
---
最后更新于 2025-01-21

# offBluetoothDeviceFound(function callback)

取消监听寻找到新设备的事件

## 支持说明

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | V3.25.0+ | V3.25.0+ | **X** | V7.35.0+ | 预览
网页应用 | V5.16.0+ | V5.16.0+ | **X** | V7.35.0+ | /

## 输入

名称 | 数据类型 | 必填 | 默认值 | 描述
--- | --- | --- | --- | ---
callback | function | 是 |  | 该事件的回调函数

## 输出
无

## 示例代码

```js
tt.offBluetoothDeviceFound(this.callback);
```

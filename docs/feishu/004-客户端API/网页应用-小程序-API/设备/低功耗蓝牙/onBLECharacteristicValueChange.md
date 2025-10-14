---
title: "onBLECharacteristicValueChange"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/ble/onblecharacteristicvaluechange
last_remote_update: 2025-04-15
last_remote_update_timestamp: 1744683771000
---
最后更新于 2025-04-15

# onBLECharacteristicValueChange(function callback)
监听特征值数据变化
**Notice**：注意事项：
- 为防止多次注册事件监听导致一次事件多次回调，建议每次调用on方法监听事件之前，先调用off方法，关闭之前的事件监听。

## 支持说明

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | V3.25.0+ | V3.25.0+ | **X** | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Fbluetooth%2Fbluetooth)
网页应用 | V7.3.0+ | V7.3.0+ | **X** | V7.35.0+ | /

## 输入

名称 | 数据类型 | 必填 | 默认值 | 描述
--- | --- | --- | --- | ---
callback | function | 是 |  | 该事件的回调函数

## 输出
回调函数返回对象的属性：

名称 | 数据类型 | 描述
--- | --- | ---
deviceId | string | 蓝牙设备 ID，参考 device 对象。  
<md-alert type="tip" icon="none">**注意**：  
- Harmony 字段返回为虚拟值。  
  <md-alert>
serviceId | string | 蓝牙特征值对应 service 的 UUID。
characteristicId | string | 蓝牙特征值的 UUID。
value | hex string | 特征值最新的 16 进制值。

## 示例代码

<md-download-code href="https://open.feishu.cn/document/uYjL24iN/uYDM04iNwQjL2ADN" mobileDisplay="none">下载示例代码</md-download-code>

<div style="display: flex">
          [预览小程序](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Fbluetooth%2Fbluetooth)

</div> 

```js
tt.onBLECharacteristicValueChange(function(res) {
    console.log(JSON.stringify(res));
});
```

返回对象示例：

```json
{
  "deviceId": "E5:66:9F:82:46:61",
  "serviceId": "xxx",
  "characteristicId": "xxx",
  "value": "xxx"
}
```

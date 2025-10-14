---
title: "getBluetoothAdapterState"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/bluetooth/getbluetoothadapterstate
last_remote_update: 2025-04-15
last_remote_update_timestamp: 1744683468000
---
最后更新于 2025-04-15

# getBluetoothAdapterState(Object object)

获取本机蓝牙适配器状态。

## 支持说明

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | V3.25.0+ | V3.25.0+ | **X** | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Fbluetooth%2Fbluetooth)
网页应用 | V3.44.0+ | V3.44.0+ | **X** | V7.35.0+ | 预览

## 输入

继承[标准对象输入](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM)，无扩展属性

## 输出

`success`返回对象的扩展属性：

名称 | 数据类型 | 描述
--- | --- | ---
discovering | boolean | 是否正在搜索设备
available | boolean | 蓝牙适配器是否可用

## 示例代码

<md-download-code href="https://open.feishu.cn/document/uYjL24iN/uYDM04iNwQjL2ADN" mobileDisplay="none">下载示例代码</md-download-code>

<div style="display: flex">
          [预览小程序](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Fbluetooth%2Fbluetooth)
          预览网页应用

</div> 

```js
tt.getBluetoothAdapterState({ 
    success(res) {
      console.log(JSON.stringify(res));
    },
    fail(res) {
      console.log(`getBluetoothAdapterState fail: ${JSON.stringify(res)}`);
    }
});
```

`success`返回对象示例：
```json
{
    "available": true,
    "discovering": false,
    "errMsg": "getBluetoothAdapterState:ok"
}
```
`fail`返回对象示例：
```json
{
    "errMsg": "getBluetoothAdapterState:fail not init",
    "errCode": 10000
}
```

## 错误码
`fail`返回对象中会包含[errorCode属性](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM#a825f4c8)，代表错误码。

通用错误码可参见 [蓝牙 API 错误码](https://open.feishu.cn/document/uYjL24iN/uYzNxYjL2cTM24iN3EjN)

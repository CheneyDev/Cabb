---
title: "getConnectedWifi"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/wi-fi/getconnectedwifi
last_remote_update: 2025-02-19
last_remote_update_timestamp: 1739949163000
---
最后更新于 2025-02-19

# getConnectedWifi(Object object)

getConnectedWifi(Object object) 用于获取设备当前所连的 Wifi。

## 注意事项

调用该接口前，需要用户授权 `scope.userLocation`。了解如何授权，参见 [API 权限](https://open.feishu.cn/document/uYjL24iN/uITMuITMuITM)。

## 支持说明

该接口支持小程序和网页应用调用，对应的客户端版本支持情况如下所示。

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | **✓** | **✓** | **✓** | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=%2Fpage%2FAPI%2Fpages%2Fget-connected-wifi%2Fget-connected-wifi)
网页应用 | V3.44.0+ | V3.44.0+ | V3.47.0+ | V7.35.0+ | 预览

## 输入

该接口继承[标准对象输入](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM)，无扩展属性。

## 输出

该接口继承[标准对象输出](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM#8c92acb8)，`success` 返回对象的扩展属性如下所示。

名称 | 数据类型 | 描述
--- | --- | ---
SSID | string | Wi-Fi 的 SSID。
BSSID | string | Wi-Fi 的 BSSID。  
**Notice**：**注意**：  
  - PC 端暂不支持输出该字段。  
  - Harmony 飞书V7.37.0+版本且Harmony API 14及以上版本返回真实值，否则返回虚拟值。  
  <md-alert>
secure | boolean | Wi-Fi 是否安全。  
<md-alert type="tip" icon="none">  
**注意**：  
- iOS 端暂不支持输出该字段。  
- Android 10 及以上版本不支持输出该字段。
signalStrength | number | Wi-Fi 信号强度。  
**Notice**：**注意**：iOS 端暂不支持输出该字段。

## 示例代码

调用示例：

<md-download-code href="https://open.feishu.cn/document/uYjL24iN/uYDM04iNwQjL2ADN" mobileDisplay="none">下载示例代码</md-download-code>

<div style="display: flex">
          [预览小程序](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=%2Fpage%2FAPI%2Fpages%2Fget-connected-wifi%2Fget-connected-wifi)
          预览网页应用

</div> 

```js
tt.getConnectedWifi({
    success(res) {
      console.log(JSON.stringify(res));
    },
    fail(res) {
      console.log(`getConnectedWifi fail: ${JSON.stringify(res)}`);
    }
});
```

`success`返回对象示例：
```json
{
    "SSID": "Future Inc",
    "errMsg": "getConnectedWifi:ok",
    "BSSID": "b8:4b:5a:d8:28:d2"
}
```

## 错误码

`fail` 返回对象中可能包含 errno 属性，表示错误码。关于 errno 错误码的详细说明以及通用错误码列表，可参见[Errno 错误码](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/errno)。

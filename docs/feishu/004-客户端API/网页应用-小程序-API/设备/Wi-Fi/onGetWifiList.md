---
title: "onGetWifiList"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/wi-fi/ongetwifilist
last_remote_update: 2025-01-21
last_remote_update_timestamp: 1737432871000
---
最后更新于 2025-01-21

# onGetWifiList(function callback)

onGetWifiList(function callback) 用于监听获取到 Wi-Fi 列表数据事件。
**Notice**：该接口需要与 [getWifiList](https://open.feishu.cn/document/uYjL24iN/uUDO4UjL1gDO14SN4gTN) 接口配套使用。

## 支持说明

该接口支持小程序和网页应用调用，对应的客户端版本支持情况如下所示。

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | **✓** | **X** | **X** | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=%2Fpage%2FAPI%2Fpages%2Fwifilist%2Fwifilist)
网页应用 | V3.44.0+ | **X** | **X** | V7.35.0+ | 预览

## 输入

名称 | 数据类型 | 是否必填 | 默认值 | 描述
--- | --- | --- | --- | ---
callback | function | 是 | \- | 该事件的回调函数。

## 输出

回调函数返回对象的属性：

名称 | 数据类型 | 描述
--- | --- | ---
wifiList | object[] | Wi-Fi 列表数据。
&emsp;  
                    ∟  
                &nbsp;  
                    SSID | string | Wi-Fi 的 SSID。
&emsp;  
                    ∟  
                &nbsp;  
                    BSSID | string | Wi-Fi 的 BSSID。(Harmony 字段返回为虚拟值)
&emsp;  
                    ∟  
                &nbsp;  
                    secure | boolean | Wi-Fi 是否安全。
&emsp;  
                    ∟  
                &nbsp;  
                    signalStrength | number | Wi-Fi 信号强度。

## 示例代码

调用示例：

<md-download-code href="https://open.feishu.cn/document/uYjL24iN/uYDM04iNwQjL2ADN" mobileDisplay="none">下载示例代码</md-download-code>

<div style="display: flex">
          [预览小程序](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=%2Fpage%2FAPI%2Fpages%2Fwifilist%2Fwifilist)
          预览网页应用

</div> 

```js
tt.getWifiList();
tt.onGetWifiList(function(res) {
    console.log(JSON.stringify(res));
});
```

回调函数返回对象示例：
```json
{
    "wifiList": [
        {
            "BSSID": "b8:3a:6b:d8:28:d3",
            "SSID": "WIFI Test",
            "innerCode": 0,
            "secure": true,
            "signalStrength": 56
        },
        {
            "BSSID": "b8:8a:8a:d8:28:d2",
            "SSID": "Future Inc",
            "innerCode": 0,
            "secure": true,
            "signalStrength": 56
        }
    ]
}
```

---
title: "getBeacons"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/ibeacon/getbeacons
last_remote_update: 2025-01-21
last_remote_update_timestamp: 1737432721000
---
最后更新于 2025-01-21

# getBeacons(Object object)

getBeacons(Object object) 用于获取所有已搜索到的 Beacon 设备。

## 注意事项

调用该接口前，需要先调用 [startBeaconDiscovery](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/ibeacon/startbeacondiscovery)。

## 支持说明

该接口支持小程序和网页应用调用，对应的客户端版本支持情况如下所示。

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | V4.6.0+ | V4.6.0+ | **X** | V7.35.0+ | 预览
网页应用 | V4.6.0+ | V4.6.0+ | **X** | V7.35.0+ | 预览

## 输入

该接口继承[标准对象输入](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM)，无扩展属性。

## 输出

该接口继承[标准对象输出](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM#8c92acb8)，`success` 返回对象的扩展属性如下所示。

名称 | 数据类型 | 描述
--- | --- | ---
beacons | [BeaconInfo[]](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/ibeacon/ibeaconinfo) | Beacon 设备列表。
&emsp;  
                ∟  
                &nbsp;  
                uuid | string | Beacon 设备广播的 uuid。
&emsp;  
                ∟  
                &nbsp;  
                major | number | Beacon 设备的主 id。  
**Notice**：**注意**：历史存在 Android 中 `major` 以及 `minor` 字段为 string 类型的问题，该问题在飞书 V5.26.0 版本已修复。
&emsp;  
                ∟  
                &nbsp;  
                minor | number | Beacon 设备的次 id。  
**Notice**：**注意**：历史存在 Android 中 `major` 以及 `minor` 字段为 string 类型的问题，该问题在飞书 V5.26.0 版本已修复。
&emsp;  
                ∟  
                &nbsp;  
                proximity | number | 表示设备距离的枚举值。可能值：  
- `0`：无效  
- `1`：非常近  
- `2`：近  
- `3`：远
&emsp;  
                ∟  
                &nbsp;  
                accuracy | number | Beacon 设备的距离。
&emsp;  
                ∟  
                &nbsp;  
                rssi | number | 表示设备的信号强度。

## 示例代码

调用示例：

```js
tt.startBeaconDiscovery({
    uuids: [
        "fda50693-a4e2-4fb1-afcf-c6eb07647825"
    ],
    ignoreBluetoothAvailable: true,
    success(res) {
      tt.getBeacons({ 
    	success(res) {
      	  console.log(JSON.stringify(res));
    	},
    	fail(res) {
      	  console.log(`getBeacons fail: ${JSON.stringify(res)}`);
    	}
	    });
    },
    fail(res) {
      console.log(`startBeaconDiscovery fail: ${JSON.stringify(res)}`);
    }
});
```

`success`返回对象示例：
```json
{
    "beacons": [
        {
            "accuracy": 10.467388920465797,
            "major": 11054,
            "minor": 11389,
            "proximity": 3,
            "rssi": -81,
            "uuid": "FDA50693-A4E2-4FB1-AFCF-C6EB07647825"
        }
    ],
    "errMsg": "getBeacons:ok"
}
```

## 错误码

`fail` 返回对象中可能包含 errCode 属性和 errno 属性，均代表错误码。具体错误码可查阅 [Beacon API 错误码](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/ibeacon/ibeacon-api-error-code) 或 [Errno 错误码](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/errno)。

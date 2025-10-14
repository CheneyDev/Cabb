---
title: "getLocation"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/location/getlocation
last_remote_update: 2025-09-02
last_remote_update_timestamp: 1756779772000
---
最后更新于 2025-09-02

# getLocation

调用 getLocation(Object object) 获取设备当前的地理位置。

## 注意事项
- 调用前需要用户授权 `scope.userLocation`。如何授权可参见 [API 权限](https://open.feishu.cn/document/uYjL24iN/uITMuITMuITM)。
- 该 API 需要用户在手机系统中给飞书客户端授予地理位置权限，位置精度和调用耗时会因设备而异。
- 该 API 有一定性能消耗，请注意不要频繁调用以防设备过热和耗电过快，小程序框架也会做相应的节流处理。

## 支持说明

该接口支持小程序和网页应用调用，对应的客户端版本支持情况如下所示。

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | **✓** | **✓** | **X** | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Fget-location%2Fget-location)
网页应用 | V3.44.0+ | V3.44.0+ | **X** | V7.35.0+ | 预览

## 输入

该接口继承[标准对象输入](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM)，扩展属性如下所示。

名称 | 数据类型 | 是否必填 | 默认值 | 描述
--- | --- | --- | --- | ---
timeout | number | 否 | 5 | 定位超时时间，单位：秒。如果传入允许范围之外的数值，会根据 `accuracy` 字段自动选择，`best` 模式下 10 秒，`high`模式下 3 秒。  
- 最小值：3  
- 最大值：180
cacheTimeout | number | 否 | 0 | 定位缓存超时时间，单位：秒。每次定位缓存当前定位数据，并记下时间戳，当下次调用处于 `cacheTimeout` 之内时，返回缓存数据。如果 `cacheTimeout` 小于 0 或大于 60s，则不使用缓存。  
- 最小值：0  
- 最大值：60
accuracy | string | 否 | high | 指定期望精度，支持 `high`、`best`。当指定 `high` 时，期望精度值为 100 m，当指定 `best` 时，期望精度值为 20 m。当定位得到的精度不符合条件时，在 `timeout` 之前会继续定位，尝试获取符合要求的定位结果。可选值：  
- `high`：期望精度值为 100 m。  
- `best`：期望精度值为 20 m。

## 输出

该接口继承[标准对象输出](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM#8c92acb8)，`success` 返回对象的扩展属性如下所示。

名称 | 数据类型 | 描述
--- | --- | ---
latitude | number | 纬度。范围为 `-90` ~ `90`，正数表示北，负数表示南。
longitude | number | 经度。范围为 `-180` ~ `180`，正数表示东，负数表示西。
accuracy | number | 位置的精确度。  
**Notice**：**注意**：Android/iOS 端均返回水平精度。
verticalAccuracy | number | 垂直精度。单位：m  
**Notice**：**注意**：Android 无法获取该精度，会返回 0。
horizontalAccuracy | number | 水平精度。单位：m
authorizationAccuracy | string | 指示应用程序有权使用的位置准确性级别。可能值：  
- `reduced`：非精确位置授权。  
- `full`：精确位置授权。  
**Notice**：**注意**：仅 iOS 14 且飞书版本为 [V3.36.0](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility) 及以上时，支持输出该字段。
timestamp | number | 定位数据的时间戳。单位：ms
type | string | 坐标系类型。可能值：  
- `wgs84`：wgs84 坐标系。  
- `gcj02`：gcj02 坐标系。

## 示例代码

调用示例：

<md-download-code href="https://open.feishu.cn/document/uYjL24iN/uYDM04iNwQjL2ADN" mobileDisplay="none">下载示例代码</md-download-code>

<div style="display: flex">
          [预览小程序](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Fget-location%2Fget-location)
          预览网页应用

</div> 

```js
tt.getLocation({
    "type": "gcj02",
    "timeout": 5,
    "cacheTimeout": 30,
    "accuracy": "best",
    success(res) {
      console.log(JSON.stringify(res));
    },
    fail(res) {
      console.log(`getLocation fail: ${JSON.stringify(res)}`);
    }
});
```

`success`返回对象示例：
```json
{
    "authorizationAccuracy": "full",
    "latitude": 30.48318898654514,
    "errMsg": "getLocation:ok",
    "longitude": 120.03518184678819,
    "accuracy": 148,
    "horizontalAccuracy": 148,
    "verticalAccuracy": 13.574329376220703,
    "timestamp": 1637490791204
}
``` 

## 错误码

`fail` 返回对象中可能包含 errno 属性，表示错误码。关于 errno 错误码的详细说明以及通用错误码列表，可参见[Errno 错误码](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/errno)。

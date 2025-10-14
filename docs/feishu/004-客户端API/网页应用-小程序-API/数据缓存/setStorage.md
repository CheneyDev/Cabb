---
title: "setStorage"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/cache/setstorage
last_remote_update: 2025-01-21
last_remote_update_timestamp: 1737434326000
---
最后更新于 2025-01-21

# setStorage

调用 setStorage(Object object) 以键值对的形式设置本地缓存数据。

## 注意事项

单个 key 允许存储的最大数据长度为 1 MB，所有数据存储上限为 10 MB。同时，受用户设备存储空间、缓存清理等机制的限制，可能会导致信息丢失，因此请不要将重要数据存放在本地数据缓存中。

## 支持说明

该接口仅支持小程序调用，对应的客户端版本支持情况如下所示。

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | **✓** | **✓** | **✓** | V7.35.0+ | 预览
网页应用 | **x** | **x** | **x** | **x** | /

## 输入

该接口继承[标准对象输入](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM)，扩展属性如下所示。

名称 | 数据类型 | 是否必填 | 默认值 | 描述
--- | --- | --- | --- | ---
key | string | 是 | \- | 键名。最小长度为 1 字符。示例值：name
data | [本地缓存数据类型](https://open.feishu.cn/document/uYjL24iN/uMTOz4yM5MjLzkzM) | 否 | undefined | 键名对应的数据。示例值：DemoName

## 输出

该接口继承[标准对象输出](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM#8c92acb8)，无扩展属性。

## 示例代码

调用示例：

```js
tt.setStorage({
    key: "name",
    data: "DemoName",
    success(res) {
      console.log(JSON.stringify(res));
    },
    fail(res) {
      console.log(`setStorage fail: ${JSON.stringify(res)}`);
    }
});
```

`success`返回对象示例：
```json
{
    "errMsg": "setStorage:ok"
}
```

## 错误码

`fail` 返回对象中可能包含 errno 属性，表示错误码。关于 errno 错误码的详细说明以及通用错误码列表，可参见[Errno 错误码](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/errno)。

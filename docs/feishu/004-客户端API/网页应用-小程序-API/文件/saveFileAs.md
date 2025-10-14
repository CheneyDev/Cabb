---
title: "saveFileAs"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/file/savefileas
last_remote_update: 2025-04-15
last_remote_update_timestamp: 1744683793000
---
最后更新于 2025-04-15

# saveFileAs

saveFileAs(Object object) 用于保存文件到本地指定目录。

## 支持说明

该接口支持小程序和网页应用调用，对应的客户端版本支持情况如下所示。

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | **X** | **X** | V3.9.0+ | **X** | 预览
网页应用 | **X** | **X** | V5.16.0+ | **X** | /

## 输入

该接口继承[标准对象输入](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM)，扩展属性如下所示。

名称 | 数据类型 | 是否必填 | 默认值 | 描述
--- | --- | --- | --- | ---
filePath | string | 是 | \- | 文件路径。  
**示例值**：ttfile://user/feishu.png  
**Notice**：**注意**：不支持网络地址。

## 输出

该接口继承[标准对象输出](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM#8c92acb8)，无扩展属性。

## 示例代码

调用示例：

```js
tt.chooseImage({
  success(res) {
    const filePath = res.tempFilePaths[0];
        tt.saveFileAs({
          filePath,
          success(res) {
            console.log(`${JSON.stringify(res)}`);
          },
          fail(res) {
            console.log(`saveFileAs fail: ${JSON.stringify(res)}`);
          },
        });
  },
});
```

`success`返回对象示例：

```json
{
  "errMsg": "saveFileAs:ok"
}
```

## 错误码

`fail` 返回对象中可能包含 errno 属性，表示错误码。关于 errno 错误码的详细说明以及通用错误码列表，可参见[Errno 错误码](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/errno)。

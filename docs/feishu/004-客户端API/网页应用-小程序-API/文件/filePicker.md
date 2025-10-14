---
title: "filePicker"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/file/filepicker
last_remote_update: 2025-04-15
last_remote_update_timestamp: 1744683815000
---
最后更新于 2025-04-15

# filePicker

filePicker(Object object) 用于打开附件选择列表。

## 支持说明

该接口支持小程序和网页应用调用，对应的客户端版本支持情况如下所示。

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | **✓** | **✓** | **✓** | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Ffile%2Ffile)
网页应用 | V3.44.0+ | V3.44.0+ | V3.47.0+ | V7.35.0+ | 预览

## 输入

该接口继承[标准对象输入](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM)，扩展属性如下所示。

名称 | 数据类型 | 是否必填 | 默认值 | 描述
--- | --- | --- | --- | ---
maxNum | number | 否 | 不限制数量 | 选择附件时的最大选择数量。  
**示例值**：10
isSystem | boolean | 否 | false | 是否使用系统的文件选择器。取值：  
- true：使用系统的文件选择器，可以选择系统文件系统里的文件，并且 `maxNum` 参数会被设置为 `1`。  
- false：使用飞书的文件选择器，可以选择飞书文档。  
**Notice**：**注意**：  
- Android/iOS 端：飞书 [V3.8.0](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility) 及以上版本支持设置该字段。  
- PC 端：暂不支持设置该字段。  
- Harmony 端：暂不支持设置该字段，当前仅支持系统文件选择器。
pickerTitle | string | 否 | Select a file | 自定义文件选择器标题，仅在飞书文件选择器（`isSystem=false`）生效。  
  **Notice**：**注意**：飞书 [V3.37.0](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility) 及以上版本支持设置该字段。
pickerConfirm | string | 否 | Confirm | 自定义组件的选择按钮文案，仅在飞书文件选择器（`isSystem=false`）生效。  
  **Notice**：**注意**：  
- Android/iOS 端：飞书 [V3.37.0](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility) 及以上版本支持设置该字段。  
- PC 端：飞书 [V3.37](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility)~[V3.40](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility) 版本支持设置该字段。
supportFileTypes | string[] | 否 | \- | 文件类型，可根据指定的文件类型过滤文件列表，仅在飞书文件选择器（`isSystem=false`）生效。  
  **Notice**：**注意**：  
- Android：飞书 [V7.7.0](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility) 及以上版本支持设置该字段。  
- iOS 和 PC 端：暂不支持设置该字段  
- 文件类型 包括但不限于**mp4、png、doc、pdf**等    
- 可以传递多个文件类型，如 **['pdf','png']**  
- 如果不传递此参数，默认展示所有的文件

## 输出

该接口继承[标准对象输出](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM#8c92acb8)，`success` 返回对象的扩展属性如下所示。

名称 | 数据类型 | 描述
--- | --- | ---
list | object[] | 文件列表。
&emsp;  
∟  
&nbsp;  
path | string | 文件路径。
&emsp;  
∟  
&nbsp;  
name | string | 文件名。
&emsp;  
∟  
&nbsp;  
size | string | 文件大小，以字节为单位。

## 示例代码

调用示例：

<md-download-code href="https://open.feishu.cn/document/uYjL24iN/uYDM04iNwQjL2ADN" mobileDisplay="none">下载示例代码</md-download-code>
  <div style="display: flex">
    [预览小程序 ](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Ffile%2Ffile)
    预览网页应用

</div> 

```js
tt.filePicker({
  maxNum: 10,
  pickerTitle: "Select a file",
  pickerConfirm: "Confirm",
  isSystem: false,
  success(res) {
    console.log(JSON.stringify(res));
  },
  fail(res) {
    console.log(`filePicker fail: ${JSON.stringify(res)}`);
  },
});
```

`success`返回对象示例：

```json
{
  "errMsg": "filePicker:ok",
  "list": [
    {
      "path": "ttfile://temp/311d60f2-9a56-4f43-8758-92aa7378fad8-app.js",
      "name": "app.js",
      "size": 40
    }
  ]
}
```

## 错误码

`fail` 返回对象中可能包含 errno 属性，表示错误码。关于 errno 错误码的详细说明以及通用错误码列表，可参见[Errno 错误码](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/errno)。

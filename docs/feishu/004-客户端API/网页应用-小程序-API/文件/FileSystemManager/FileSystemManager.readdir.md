---
title: "FileSystemManager.readdir"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/file/file_system_manager/file_system_manager_read_dir
last_remote_update: 2025-01-21
last_remote_update_timestamp: 1737433599000
---
最后更新于 2025-01-21

# FileSystemManager.readdir(Object object)

读取本地目录内文件列表。

## 支持说明

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | V5.23.0+ | V5.23.0+ | V5.23.0+ | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Ffile%2Ffile)
网页应用 | V5.23.0+ | V5.23.0+ | V5.23.0+ | V7.35.0+ | 预览

## 输入
继承[标准对象输入](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM)，扩展属性描述：

名称 | 数据类型 | 必填 | 默认值 | 描述
--- | --- | --- | --- | ---
dirPath | string | 是 |  | 要读取的目录路径

## 输出
继承[标准对象输出](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM#8c92acb8)，`success`返回对象的扩展属性：

名称 | 数据类型 | 描述
--- | --- | ---
files | Array<string\> | 指定目录下的文件名数组

## 代码示例

```js
const fileSystemManager = tt.getFileSystemManager();

fileSystemManager.readdir({
  dirPath: "ttfile://user/",
  success(res) {
    console.log("调用成功", res.files);
  },
  fail(res) {
    console.log("调用失败", res.errMsg);
  },
});
```

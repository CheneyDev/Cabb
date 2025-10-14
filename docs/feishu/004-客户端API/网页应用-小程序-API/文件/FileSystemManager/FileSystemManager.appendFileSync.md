---
title: "FileSystemManager.appendFileSync"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/file/file_system_manager/appendfilesync
last_remote_update: 2025-01-21
last_remote_update_timestamp: 1737433407000
---
最后更新于 2025-01-21

# FileSystemManager.appendFileSync(string filePath, string|ArrayBuffer data, string encoding)

在文件结尾追加内容。

## 注意事项

单次追加文件大小不能超过 10M。文件总大小不能超过用户目录总大小 200M。

## 支持说明

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | V5.23.0+ | V5.23.0+ | **X** | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Ffile%2Ffile)
网页应用 | **X** | **X** | **X** | **X** | 预览

## 输入

名称 | 数据类型 | 必填 | 默认值 | 描述
--- | --- | --- | --- | ---
filePath | string | 是 |  | 本地文件路径
data | string/ArrayBuffer | 是 |  | 要写入的文本或二进制数据
encoding | string | 否 | utf8 | 指定写入文件的字符编码，默认值仅在data为string类型时生效

### encoding 的合法值

| 值                          | 说明         |Android|iOS|
| --------------------------- | ------------ |------|---|
| ascii                       |              |支持|支持|
| base64                      |              |支持|支持|
| binary                      |              |支持|支持|
| hex                         |              |支持|支持|
| ucs2/ucs-2/utf16le/utf-16le | 以小端序读取   |支持|支持|
| utf-8/utf8                  |              |支持|支持|
| latin1                      |              |支持|支持|

## 输出
继承[标准对象输出](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM#8c92acb8)，无扩展属性

## 代码示例

```js
const fileSystemManager = tt.getFileSystemManager();

const filePath = `ttfile://user/example.txt`;

try {
  fileSystemManager.appendFileSync(filePath, "example append content", "utf8");
  console.log("调用成功");
} catch (err) {
  console.log("调用失败", err);
}
```

## 已知问题

- `encoding`参数与`data`的实际编码需要保持一致，否则可能会无法写入成功

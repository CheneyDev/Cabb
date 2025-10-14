---
title: "Stats.isFile"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/file/stat/stats_is_file
last_remote_update: 2025-03-12
last_remote_update_timestamp: 1741768346000
---
最后更新于 2025-03-12

# Stats.isFile()

判断当前文件是否一个普通文件。

## 支持说明

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | V4.11.0+ | V4.11.0+ | **X** | V7.39.0+ | 预览
网页应用 | V4.11.0+ | V4.11.0+ | **X** | V7.39.0+ | 预览

## 输入
无

## 输出

返回值：boolean

## 示例代码

```js
const fileSystemManager = tt.getFileSystemManager();
fileSystemManager.stat({
    path: "ttfile://temp",
    success(res) {
      console.log(res.stat.isFile());
    },
    fail(res) {
      console.log(`stat fail: ${JSON.stringify(res)}`);
    }
});
```
返回值示例：

```json
false
```

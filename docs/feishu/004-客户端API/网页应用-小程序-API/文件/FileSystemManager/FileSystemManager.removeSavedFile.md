---
title: "FileSystemManager.removeSavedFile"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/file/file_system_manager/file_system_manager_remove_saved_file
last_remote_update: 2025-01-21
last_remote_update_timestamp: 1737433664000
---
最后更新于 2025-01-21

# FileSystemManager.removeSavedFile(Object object)

删除保存的本地文件。

## 支持说明

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | V5.23.0+ | V5.23.0+ | V5.23.0+ | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=page%2FAPI%2Fpages%2Ffile%2Ffile)
网页应用 | V5.23.0+ | V5.23.0+ | V5.23.0+ | V7.35.0+ | 预览

## 输入
继承[标准对象输入](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM)，扩展属性描述：

名称 | 数据类型 | 必填 | 默认值 | 描述
--- | --- | --- | --- | ---
filePath | string | 是 |  | 本地文件路径

## 代码示例

```js
const fileSystemManager = tt.getFileSystemManager();

fileSystemManager.getSavedFileList({
  success(res) {
    console.log("当前用户目录下的文件:", res.fileList);
    removeSavedFile(res.fileList);
  },
});

function removeSavedFile(fileList) {
  fileList.forEach((filePath) => {
    fileSystemManager.removeSavedFile({
      filePath,
      success(res) {
        console.log(`${filePath} 删除成功`);
      },
      fail(res) {
        console.error(`${filePath} 删除失败`, res.errMsg);
      },
    });
  });
}
```

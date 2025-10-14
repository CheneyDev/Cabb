---
title: "getNFCAdapter"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/nfc/getnfcadapter
last_remote_update: 2025-08-25
last_remote_update_timestamp: 1756093503000
---
最后更新于 2025-08-25

# getNFCAdapter()

获取客户端NFC适配器

## 支持说明

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | <md-version>V3.38.0+<md-version> | V5.25.0+ | **X** | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=%2Fpage%2FAPI%2Fpages%2Fnfc%2Fnfc)
网页应用 | V3.44.0+ | V5.25.0+ | **X** | V7.35.0+ | 预览

## 输入
无

## 输出

返回值：`NFCAdapter`，该对象的方法列表参见下表：
**Notice**：点击下表中的方法名，查看对应API的支持说明、调用方法

方法 | 介绍
--- | ---
[getNfcA](https://open.feishu.cn/document/uYjL24iN/ugzM4YjL4MDO24COzgjN) | 获取NfcA实例，实例支持NFC-A (ISO 14443-3A)标准的读写
[getNfcB](https://open.feishu.cn/document/uYjL24iN/uUzM4YjL1MDO24SNzgjN/getnfcb) | 获取NfcB实例，实例支持NFC-B (ISO 14443-3B)标准的读写  
**Notice**：飞书[V7.15](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility)及以上版本支持
[getNfcV](https://open.feishu.cn/document/uYjL24iN/uUzM4YjL1MDO24SNzgjN/getNfcV) | 获取NfcV实例，实例支持NFC-V (ISO 15693)标准的读写  
**Notice**：飞书[V5.14](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility)及以上版本支持
[getIsoDep](https://open.feishu.cn/document/uYjL24iN/uUzM4YjL1MDO24SNzgjN/getisodep) | 获取IsoDep实例，实例支持ISO-DEP (ISO 14443-4)标准的读写  
**Notice**：飞书[V7.15](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility)及以上版本支持
[getMifareClassic](https://open.feishu.cn/document/uYjL24iN/uEDN4YjLxQDO24SM0gjN) | 获取 MifareClassic 实例，实例支持 MIFARE Classic 标签的读写
[getNdef](https://open.feishu.cn/document/uYjL24iN/uUzM4YjL1MDO24SNzgjN/nfcadapter_getndef) | 获取Ndef实例，实例支持NDEF标签的读写  
**Notice**：飞书[V7.48](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/version-compatibility)及以上版本支持
[startDiscovery](https://open.feishu.cn/document/uYjL24iN/uIDN4YjLyQDO24iM0gjN) | 开始扫描 NFC 标签
[stopDiscovery](https://open.feishu.cn/document/uYjL24iN/uMDN4YjLzQDO24yM0gjN) | 关闭 NFC 标签扫描
[onDiscovered](https://open.feishu.cn/document/uYjL24iN/uUDN4YjL1QDO24SN0gjN) | 监听 NFC Tag
[offDiscovered](https://open.feishu.cn/document/uYjL24iN/uQDN4YjL0QDO24CN0gjN) | 取消监听 NFC Tag

## 示例代码

<md-download-code href="https://open.feishu.cn/document/uYjL24iN/uYDM04iNwQjL2ADN" mobileDisplay="none">下载示例代码</md-download-code>

<div style="display: flex">
          [预览小程序](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=%2Fpage%2FAPI%2Fpages%2Fnfc%2Fnfc)
              预览网页应用

</div> 

```js
const adapter = tt.getNFCAdapter()
```

## 错误码
`fail`返回对象中会包含[errno属性](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/errno)，代表错误码。

通用错误码可参见 [NFC API 错误码](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/nfc-error-codes)

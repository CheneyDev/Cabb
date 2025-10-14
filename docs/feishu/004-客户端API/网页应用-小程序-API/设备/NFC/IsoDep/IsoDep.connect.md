---
title: "IsoDep.connect"
source_url: https://open.feishu.cn/document/web-app/gadget-api/device/nfc/isodep/connect
last_remote_update: 2025-08-25
last_remote_update_timestamp: 1756093855000
---
最后更新于 2025-08-25

# IsoDep.connect(Object object)

连接 ISO-DEP (ISO 14443-4) 类型的标签

## 支持说明

应用能力 | Android | iOS | PC | Harmony | 预览效果
--- | --- | --- | --- | --- | ---
小程序 | V7.15.0+ | **X** | **X** | V7.35.0+ | [预览](https://applink.feishu.cn/client/mini_program/open?appId=cli_9dff7f6ae02ad104&path=%2Fpage%2FAPI%2Fpages%2Fnfc%2Fnfc)
网页应用 | V7.15.0+ | **X** | **X** | V7.35.0+ | 预览

## 输入
继承[标准对象输入](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM)，无扩展属性

## 输出
继承[标准对象输出](https://open.feishu.cn/document/uYjL24iN/ukzNy4SO3IjL5cjM#8c92acb8)，无扩展属性

## 示例代码
**Notice**：这里仅示例当前 API 使用，实际流程请参考：[接入 NFC](https://open.feishu.cn/document/uYjL24iN/ugTN4YjL4UDO24CO1gjN)

```js
let adapter = tt.getNFCAdapter();
adapter.onDiscovered(
    (res) => {
        console.log('onDiscovered res, ' + JSON.stringify(res));
        // 扫描到NFC标签后，判断是否包含ISO-DEP标签
        if（res.techs.indexOf("ISO-DEP") != -1）{
            let isoDep = adapter.getIsoDep();
            // 连接ISO-DEP标签
            isoDep.connect({
                success(res) {
                    console.log(JSON.stringify(res))
                },
                fail(err) {
                    console.log(`IsoDep.connect fail: ${JSON.stringify(res)}`)
                }
            });   
        }   
    }
);
```
`success`返回对象示例：
```json
{"errMsg":"nfcConnect:ok"}
```

## 错误码
`fail`返回对象中会包含[errno属性](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/errno)，代表错误码。

通用错误码可参见 [NFC API 错误码](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/nfc-error-codes)

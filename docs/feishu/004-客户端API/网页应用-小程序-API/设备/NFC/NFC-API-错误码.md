---
title: "NFC API 错误码"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/nfc/nfc-error-codes
last_remote_update: 2025-08-01
last_remote_update_timestamp: 1754044135000
---
最后更新于 2025-08-01

# NFC API 错误码
|errno               |errString    |说明   |
|------------------|--------|-----|
|102               |Internal error      |内部错误   |
|1504001           |NFC is not avaliable|手机不支持 NFC |
|1504003             |system NFC switch not opened|手机 NFC 开关未打开|
|1504204             |Tech already connected|NFC 已连接|
|1504205             |Tech has not connected|NFC 未连接 |
|1504206             |NFC tag has not been discoveryed|系统未发现 NFC 标签|
|1504208             |unavailable tech|当前 NFC 标签不支持该tech|
|1504209             |function not support|不支持该 NFC 能力|
|1504211             |data is null |数据传输参数错误|
|1504212             |array buffer is empty|数据传输参数错误|
|1504213             |base64 value is empty|数据传输参数错误|
|1504214             |base64 decode failed|数据传输参数错误|
|1504215             |transceive data failed|数据传输command协议错误|
|1504216             |nfc type is empty|NFC type为空|
|1504217             |nfc service dead|NFC 已断开|
|1504219             |Failed to read NdefMessage from the NFC card|NFC卡片读取NdefMessage失败|
|1504220             |The NFC card does not support writing NdefMessage|NFC卡片不支持NdefMessage写操作|
|1504221             |Failed to write NdefMessage to the NFC card|NFC卡片写入NdefMessage失败|

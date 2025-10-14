---
title: "iOS Setting能力_Swift(7.18-7.31)"
source_url: https://open.feishu.cn/document/native-integration/open-capability/capability-components/setting-ability/ios-setting-capa
last_remote_update: 2025-04-29
last_remote_update_timestamp: 1745895220000
---
最后更新于 2025-04-29

# iOS Setting能力_Swift(7.18-7.31) 

|组件名称 | 组件版本 | 组件能力 |
| ---- | ------ | -------- |
| LKSettingExternal | 1.1.3 | 通过该组件，可以获取飞书的动态setting配置，只在SaaS可用，私有化不支持 |

## 示例代码

完整示例请查看 [https://github.com/larksuite/alchemy_app_demo/tree/main/alchemy_app_demo_ios](https://github.com/larksuite/alchemy_app_demo/tree/main/alchemy_app_demo_ios)

```swift
import LKSettingExternal
import LKKABridge

let setting = KAAPI(channel: /* channel_id */).settings
setting?.getConfig( ... )
```

## PROTOCOL

### KASettingProtocol

Setting 能力组件接口协议，用于获取 Setting 配置

```swift
public protocol KASettingProtocol: AnyObject
```

#### Methods
#### `getConfig(space:key:)`

获取 config
- Returns: 远端配置

```swift
func getConfig(space: String, key: String) -> [String: Any]
```

##### Parameters

| 定义名称 | 描述 |
| ---- | -- |
| key | isv key |
| space | isv space |
## EXTENSION

### KAAPI
```swift
extension KAAPI
```

#### Properties
#### `settings`

Setting 能力接口实例

```swift
@objc public var settings: KASettingProtocol?
```

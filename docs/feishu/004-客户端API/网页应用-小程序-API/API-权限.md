---
title: "API 权限"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/api-scopes
last_remote_update: 2023-11-13
last_remote_update_timestamp: 1699843826000
---
最后更新于 2023-11-13

# API 权限

为了严格保护用户信息，应用需要经过用户侧授权同意后，才能调用部分小程序接口。本文将提供对应的权限、授权行为以及有效期等相关说明。

## 权限说明

接口权限按使用范围分成多个 scope ，当用户授权给一个 scope 后，当前应用可以直接调用其对应的所有接口。
**注意**：用户可以拒绝授权，因此你（即应用开发者）需要兼容用户拒绝授权的场景。

Scope | API/组件 | 权限 | 描述 | 是否需要授权给客户端
--- | --- | --- | --- | ---
`scope.userInfo` | [getUserInfo](https://open.feishu.cn/document/uYjL24iN/ucjMx4yNyEjL3ITM) | 用户信息 | 获得你的基本信息（昵称、头像、性别及地区等）。 | ✓（飞书 5.2及以上版本生效）
`scope.userLocation` | [startLocationUpdate](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/startlocationupdate) | 地理位置 | 获取你的地理位置信息。 | ✓
`scope.camera` | [camera](https://open.feishu.cn/document/uYjL24iN/uYTNuYTNuYTN/camera) | 摄像头 | 访问你的摄像头功能。 | ✓（飞书 5.21及以上版本生效）
`scope.record` | [getRecorderManager](https://open.feishu.cn/document/uYjL24iN/uQDOx4CN4EjL0gTM)  
[CameraContext.startRecord](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/camera/cameracontext/startrecord) | 麦克风 | 访问你的麦克风功能。 | ✓
`scope.writePhotosAlbum` | [saveImageToPhotosAlbum](https://open.feishu.cn/document/uYjL24iN/uUTMx4SNxEjL1ETM) | 保存到相册 | 保存图片或视频到你的相册。 | ✓
`scope.clipboard` | [getClipboardData](https://open.feishu.cn/document/uYjL24iN/uczNx4yN3EjL3cTM) | 剪切板 | 访问剪切板。 | ✓
`scope.appBadge` | [updateBadge](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/updagebadge)  
[onServerBadgePush](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/onserverbadgepush)  
[reportBadge](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/reportbadge) | 应用红点 | 展示应用红点的权限。 | ✓（飞书 5.9及以上版本生效）
`scope.bluetooth` | [openBluetoothAdapter](https://open.feishu.cn/document/uYjL24iN/ugzNxYjL4cTM24CO3EjN) | 蓝牙 | 是否授权蓝牙权限。 | ✓
`scope.runData` | [getStepCount](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/sport-and-health/getstepcount) | 运动步数 | 是否授权获取运动步数。 | ✓

## 调用接口后的行为说明

* 如果用户已授权，会继续执行 API 后续逻辑。
* 如果用户未授权过此权限，会弹窗询问用户，用户同意授权后，会继续执行 API 后续逻辑。
* 如果用户已拒绝授权，会进入接口 `fail` 回调，并且弹窗询问用户，引导用户前往权限设置界面。

![图片](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/08e2f27cad9047a8bf988e84cf4ab4f1_UIb05Z2xz8.png?height=1334&lazyload=true&maxWidth=600&width=1580)

## 授权有效期

一旦用户明确同意或拒绝过授权，其授权信息会被记录。在飞书客户端 V3.5 及以后版本，授权信息会同步到服务器，更换设备或者卸载重装后不再需要用户重新授权。

## 系统接口权限

操作系统对于一些系统接口进行了权限控制，当调用这些接口时，需要用户对飞书进行授权。各操作系统对于这些权限的处理方式相似，但也存在不同，你需要自行做好区分。
在调用一个同时需要小程序权限和系统权限的 API 时，可能会先出现小程序 API 授权弹窗，再出现系统接口授权弹窗。

![图片](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/a1f141b5ca65f6f4f6adc3958550c469_0PCGvCc2jB.png?height=1334&lazyload=true&maxWidth=300&width=750)

---
title: "H5 JSAPI 总览"
source_url: https://open.feishu.cn/document/client-docs/h5/
last_remote_update: 2025-08-01
last_remote_update_timestamp: 1754044045000
---
最后更新于 2025-08-01

# H5 JSAPI 总览
[requestAuthCode](https://open.feishu.cn/document/uYjL24iN/uUzMuUzMuUzM/20220308)、[closeWindow](https://open.feishu.cn/document/uYjL24iN/uYTOuYTOuYTO/closewindow)、[requestAccess](https://open.feishu.cn/document/uYjL24iN/uUzMuUzMuUzM/requestaccess) API 无需[网页应用鉴权](https://open.feishu.cn/document/uYjL24iN/uEzM4YjLxMDO24SMzgjN)即可使用；<BR>
其他 JSAPI 均需要 [网页应用鉴权](https://open.feishu.cn/document/uYjL24iN/uEzM4YjLxMDO24SMzgjN) 后才能使用

飞书扫码在线预览 | 示例代码
--- | ---
![H5二维码](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/74ecc6711cfdb760830452202b32de70_NTxyKCfCzu.png?height=155&lazyload=true&width=155) | - 代码包下载：[web_app_with_jssdk.zip](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/8e80b6864660587b68d2b235d18f7ba5_YVXf93orWM.zip)  
- 使用说明：[快速开发网页应用](https://open.feishu.cn/document/home/integrating-web-apps-in-5-minutes/create-app-and-configuration)

## 开放接口

#### 鉴权

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[h5sdk.config](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/authentication/h5sdkconfig) | 对使用到的网页应用JSAPI进行鉴权 | 5.1.0 | 5.1.0 | 5.1.0 | 7.35.0

#### 登录

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[requestAccess](https://open.feishu.cn/document/uYjL24iN/uUzMuUzMuUzM/requestaccess) | 增量授予应用访问权限 | 6.9.0 | 6.9.0 | 6.9.0 | 7.35.0

#### 用户信息

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[enterProfile](https://open.feishu.cn/document/uYjL24iN/ucDM04yNwQjL3ADN) | 进入个人信息主页 | 4.0.0 | 4.0.0 | 4.0.0 | 7.35.0
[getUserInfo](https://open.feishu.cn/document/uYjL24iN/ucjMx4yNyEjL3ITM) | 获取已登录用户的基本信息或特殊信息 | 3.44.0 | 3.44.0 | 3.44.0 | 7.35.0

#### 聊天

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[enterChat](https://open.feishu.cn/document/uYjL24iN/ukDM04SOwQjL5ADN) | 打开指定会话 | **✓** | **✓** | **✓** | 7.35.0
[toggleChat](https://open.feishu.cn/document/uYjL24iN/ugDM04COwQjL4ADN/toggleChat) | 侧边栏形式打开或关闭会话 | **X** | **X** | 4.1.0 | **X**
[chooseChat](https://open.feishu.cn/document/uYjL24iN/uMTN3QjLzUzN04yM1cDN) | 打开用户会话列表选择会话，调用前确保用户已经登入 | 3.44.0 | 3.44.0 | 3.44.0 | 7.35.0
[getChatInfo](https://open.feishu.cn/document/uYjL24iN/uEDN2UjLxQjN14SM0YTN) | 获取某个会话的信息 | 3.44.0 | 3.44.0 | 3.44.0 | 7.35.0
[getBlockActionSourceDetail](https://open.feishu.cn/document/getBlockActionSourceDetail) | 支持从消息快捷操作点击进入应用后，获取应用模块对应业务的详细信息 | **✓** | **✓** | **✓** | **X**
[getTriggerContext](https://open.feishu.cn/document/no_class/plusmenu/gettriggercontext) | 加号菜单场景中，获取当前会话信息 | **✓** | **✓** | **✓** | **X**
[sendMessageCard](https://open.feishu.cn/document/uYjL24iN/uUjN5UjL1YTO14SN2kTN) | 发送消息卡片到指定会话 | 3.44.0 | 3.44.0 | 3.44.0 | 7.39.0
[onChatBadgeChange](https://open.feishu.cn/document/uYjL24iN/uQDN2UjL0QjN14CN0YTN) | 监听某个群未读消息数变化 | 7.10.0 | 7.10.0 | 7.10.0 | **X**
[offChatBadgeChange](https://open.feishu.cn/document/uYjL24iN/ugDM04COwQjL4ADN/offchatbadgechange) | 取消监听某个群未读消息数变化 | 7.10.0 | 7.10.0 | 7.10.0 | **X**

#### 联系人

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[chooseContact](https://open.feishu.cn/document/uYjL24iN/uMTM04yMxQjLzEDN) | 打开用户的联系人选择列表，可以选择用户或部门，并返回选定用户或部门的信息 | 5.7.0 | 5.7.0 | 5.7.0 | 7.35.0

#### 设置

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[openSetting](https://open.feishu.cn/document/uYjL24iN/uUzMx4SNzEjL1MTM) | 打开设置页面，展示用户设置（包括授予和拒绝）过的权限，并返回用户设置过的授权结果 | 3.44.0 | 3.44.0 | 3.44.0 | 7.35.0
[getSetting](https://open.feishu.cn/document/uYjL24iN/uQzMx4CNzEjL0MTM) | 获取用户设置（包括授予和拒绝）过的权限 | 3.44.0 | 3.44.0 | 3.44.0 | 7.35.0

#### 分享

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[shareWebContent](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/share/sharewebcontent) | 分享应用内网页内容 | 5.23.0 | 5.23.0 | 5.23.0 | 7.39.0
[share](https://open.feishu.cn/document/uYjL24iN/ugDM04COwQjL4ADN/thirdShare) | 分享内容到三方应用 | 5.14.0 | 5.14.0 | **X** | **X**

#### 授权

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[authorize](https://open.feishu.cn/document/uYjL24iN/ugzMx4COzEjL4MTM) | 向用户发出设置权限请求。如果该权限用户没有设置过，会弹窗咨询用户是否授予；如果该权限用户拒绝授予，会打开设置页面(appBadge权限除外)；如果该权限用户同意授予，会直接返回成功 | 3.44.0 | 3.44.0 | 3.44.0 | 7.35.0
[getTenantAppScopes](https://open.feishu.cn/document/uYjL24iN/uczMx4yNzEjL3MTM/gettenantappscopes) | 查询租户下当前应用的授权状态 | 3.43.0 | 3.43.0 | 3.43.0 | 7.35.0
[applyTenantAppScope](https://open.feishu.cn/document/uYjL24iN/uczMx4yNzEjL3MTM/applytenantappscope) | 弹窗咨询用户是否向租户管理员申请所有未授予权限(不包括租户敏感权限) | 3.43.0 | 3.43.0 | 3.43.0 | 7.35.0

#### 安全密码验证

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[startPasswordVerify](https://open.feishu.cn/document/uYjL24iN/ugTO3IjL4kzNy4CO5cjM) | 调起二次验证飞书安全密码的输入界面 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0

#### 系统认证

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[startDeviceCredential](https://open.feishu.cn/document/uYjL24iN/uIDN14iM0UjLyQTN) | 打开系统解锁界面 | 3.44.0 | 3.44.0 | **X** | 7.37.0

#### 水印

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[checkWatermark](https://open.feishu.cn/document/uYjL24iN/ukTM1EjL5ETNx4SOxUTM) | 查看宿主是否显示了全局水印 | 3.44 | 3.44 | 3.47 | 7.35.0

#### 邮件

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[mailto](https://open.feishu.cn/document/uYjL24iN/uAjNwEjLwYDMx4CM2ATM) | 调用系统发送邮件 | 3.44 | 3.44 | **X** | 7.35.0

#### 应用角标

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[reportBadge](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/reportbadge) | 上报应用内的角标（badge）数量 | 5.9.0 | 5.9.0 | 5.9.0 | **X**
[updateBadge](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/updagebadge) | 更新当前设备上的应用角标数 | 5.9.0 | 5.9.0 | 5.9.0 | **X**
[onServerBadgePush](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/onserverbadgepush) | 监听应用的红点（Badge）推送 | 7.10.0 | 7.10.0 | 7.10.0 | **X**
[offServerBadgePush](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/offserverbadgepush) | 取消对红点（Badge）推送的监听 | 7.10.0 | 7.10.0 | 7.10.0 | **X**

## 界面
#### 交互反馈

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[showActionSheet](https://open.feishu.cn/document/uYjL24iN/ukDNy4SO0IjL5QjM) | 显示操作菜单 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[showModal](https://open.feishu.cn/document/uYjL24iN/ugDNy4CO0IjL4QjM) | 显示模态弹窗 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[showPrompt](https://open.feishu.cn/document/uYjL24iN/uYTO4UjL2kDO14iN5gTN) | 展示可输入内容的弹窗 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[showLoading](https://open.feishu.cn/document/uYjL24iN/uMDNy4yM0IjLzQjM) | 显示灰色背景的 loading 提示框 | **✓** | **✓** | **✓** | 7.35.0
[hideLoading](https://open.feishu.cn/document/uYjL24iN/uYDNy4iN0IjL2QjM) | 隐藏 loading 提示框 | **✓** | **✓** | **✓** | 7.35.0
[showToast](https://open.feishu.cn/document/uYjL24iN/ugzMy4COzIjL4MjM) | 显示灰色背景的消息提示框 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[hideToast](https://open.feishu.cn/document/uYjL24iN/ukzMy4SOzIjL5MjM) | 隐藏灰色背景的消息提示框 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0

#### 导航栏

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[showMoreItemPanel](https://open.feishu.cn/document/uYjL24iN/uYjMy4iNyIjL2IjM/showmoreitempanel) | 唤起原生“更多”菜单面板，并按需隐藏部分菜单项 | 7.20.0 | 7.20.0 | **X** | **X**
[setNavigationBar](https://open.feishu.cn/document/uYjL24iN/uYjMy4iNyIjL2IjM/setnavigationbar) | 自定义导航栏左侧和右侧的按钮（如“返回”、“确定”） | 5.3.0 | 5.3.0 | 5.3.0 | **X**
[onLeftNavigationBarClick](https://open.feishu.cn/document/uYjL24iN/uYjMy4iNyIjL2IjM/onleftnavigationbarclick) | 监听导航栏左侧点击事件 | 5.3.0 | 5.3.0 | 5.3.0 | **X**
[onRightNavigationBarClick](https://open.feishu.cn/document/uYjL24iN/uYjMy4iNyIjL2IjM/onrightnavigationbarclick) | 监听导航栏右侧点击事件 | 5.3.0 | 5.3.0 | 5.3.0 | **X**
[setNavigationBarColor](https://open.feishu.cn/document/uYjL24iN/ucjMyUjL3IjM14yNyITN) | 设置导航栏颜色 | 6.3.0 | 6.3.0 | **X** | **X**

#### 窗口

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[setWindowSize](https://open.feishu.cn/document/uYjL24iN/uEDO3UjLxgzN14SM4cTN/setwindowsize) | 小程序在 window 和 window-semi 模式下调整独立窗口的大小和位置 | **X** | **X** | 7.10.0 | **X**

#### 页面离开二次确认

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[enableLeaveConfirm](https://open.feishu.cn/document/uYjL24iN/uUTOuUTOuUTO/page_leave_confirm_modal/enable_leave_confirm) | 当某个页面返回或关闭时设置二次弹窗提示 | 5.17.0 | 5.16.0 | 5.13.0  
      7.39.0
[disableLeaveConfirm](https://open.feishu.cn/document/uYjL24iN/uUTOuUTOuUTO/page_leave_confirm_modal/disable_leave_confirm) | 取消当前页面离开二次确认弹框 | 5.17.0 | 5.16.0 | 5.13.0  
      7.39.0

#### Pad

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[togglePadFullScreen](https://open.feishu.cn/document/uYjL24iN/uUTOuUTOuUTO/pad/togglepadfullscreen) | 在当前Pad小程序窗口可以全屏缩放的前提下，进行全屏缩放状态的切换 | 4.10.0 | 4.10.0 | **X** | **X**
[getPadDisplayScaleMode](https://open.feishu.cn/document/uYjL24iN/uUTOuUTOuUTO/pad/getpaddisplayscalemode) | 获取当前Pad的小程序窗口缩放状态，当前显示状态能否进行全屏缩放的切换 | 4.10.0 | 4.10.0 | **X** | **X**

## 设备
#### 系统信息

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getSystemInfo](https://open.feishu.cn/document/uYjL24iN/uQjNx4CN2EjL0YTM) | 获取系统信息 | 3.43.0 | 3.43.0 | 3.47.0 | 7.35.0

#### NFC

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getNFCAdapter](https://open.feishu.cn/document/uYjL24iN/ukzM4YjL5MDO24SOzgjN) | 获取 NFC 实例 | 3.44.0 | 5.25.0 | **X** | 7.35.0
[NFCAdapter.getNfcA](https://open.feishu.cn/document/uYjL24iN/ugzM4YjL4MDO24COzgjN) | 获取NfcA实例，实例支持NFC-A (ISO 14443-3A)标准的读写 | 3.44.0 | 5.25.0 | **X** | 7.35.0
[NFCAdapter.getNfcB](https://open.feishu.cn/document/uYjL24iN/uUzM4YjL1MDO24SNzgjN/getnfcb) | 获取NfcB实例，实例支持NFC-B (ISO 14443-3B)标准的读写 | 7.15.0 | 7.15.0 | **X** | 7.35.0
[NFCAdapter.getNfcV](https://open.feishu.cn/document/uYjL24iN/uUzM4YjL1MDO24SNzgjN/getNfcV) | 获取NfcV实例，实例支持NFC-V (ISO 15693)标准的读写 | 5.14.0 | **X** | **X** | 7.35.0
[NFCAdapter.getIsoDep](https://open.feishu.cn/document/uYjL24iN/uUzM4YjL1MDO24SNzgjN/getisodep) | 获取IsoDep实例，实例支持ISO-DEP (ISO 14443-4)标准的读写 | 7.15.0 | **X** | **X** | 7.35.0
[NFCAdapter.getMifareClassic](https://open.feishu.cn/document/uYjL24iN/uEDN4YjLxQDO24SM0gjN) | 获取MifareClassic实例，实例支持MIFARE Classic标签的读写 | 3.44.0 | **X** | **X** | 7.35.0
[NFCAdapter.startDiscovery](https://open.feishu.cn/document/uYjL24iN/uIDN4YjLyQDO24iM0gjN) | 开始扫描NFC标签 | 3.44.0 | 5.25.0 | **X** | 7.35.0
[NFCAdapter.stopDiscovery](https://open.feishu.cn/document/uYjL24iN/uMDN4YjLzQDO24yM0gjN) | 关闭NFC标签扫描 | 3.44.0 | 5.25.0 | **X** | 7.35.0
[NFCAdapter.onDiscovered](https://open.feishu.cn/document/uYjL24iN/uUDN4YjL1QDO24SN0gjN) | 监听 NFC Tag | 3.44.0 | 5.25.0 | **X** | 7.35.0
[NFCAdapter.offDiscovered](https://open.feishu.cn/document/uYjL24iN/uQDN4YjL0QDO24CN0gjN) | 取消监听 NFC Tag | 3.44.0 | **X** | **X** | 7.35.0
[NfcA.connect](https://open.feishu.cn/document/uYjL24iN/ucDN4YjL3QDO24yN0gjN) | 连接NfcA类型的标签 | 3.44.0 | 5.25.0 | **X** | 7.35.0
[NfcA.transceive](https://open.feishu.cn/document/uYjL24iN/uITN4YjLyUDO24iM1gjN) | 发送数据给NFCA类型的标签 | 3.44.0 | 5.25.0 | **X** | 7.35.0
[NfcA.close](https://open.feishu.cn/document/uYjL24iN/uYDN4YjL2QDO24iN0gjN) | 断开与NFCA标签之间的连接 | 3.44.0 | 5.25.0 | **X** | 7.35.0
[NfcA.getAtqa](https://open.feishu.cn/document/uYjL24iN/ugDN4YjL4QDO24CO0gjN) | 获取ATQA信息 | 3.44.0 | **X** | **X** | 7.35.0
[NfcA.getMaxTransceiveLength](https://open.feishu.cn/document/uYjL24iN/ukDN4YjL5QDO24SO0gjN) | 获取最大传输长度 | 3.44.0 | **X** | **X** | 7.35.0
[NfcA.getSak](https://open.feishu.cn/document/uYjL24iN/uATN4YjLwUDO24CM1gjN) | 获取SAK信息 | 3.44.0 | **X** | **X** | 7.35.0
[NfcA.setTimeout](https://open.feishu.cn/document/uYjL24iN/uETN4YjLxUDO24SM1gjN) | 设置超时时间 | 3.44.0 | **X** | **X** | 7.35.0
[NfcB.connect](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/nfcb/connect) | 连接 NFC-B (ISO 14443-3B) 类型的标签 | 7.15.0 | 7.15.0 | **X** | 7.35.0
[NfcB.transceive](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/nfcb/transceive) | 发送数据给NFC-B类型的标签 | 7.15.0 | 7.15.0 | **X** | 7.35.0
[NfcB.close](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/nfcb/close) | 断开与NFC-B标签之间的连接 | 7.15.0 | 7.15.0 | **X** | 7.35.0
[NfcB.getMaxTransceiveLength](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/nfcb/getmaxtransceivelength) | 获取最大传输长度 | 7.15.0 | **X** | **X** | 7.35.0
[NfcB.setTimeout](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/nfcb/settimeout) | 设置超时时间 | 7.15.0 | **X** | **X** | 7.35.0
[NfcV.connect](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/NfcV/connect) | 连接 NFC-V (ISO 15693) 类型的标签 | 5.14.0 | **X** | **X** | 7.35.0
[NfcV.transceive](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/NfcV/transceive) | 发送数据给NFCV类型的标签 | 5.14.0 | **X** | **X** | 7.35.0
[NfcV.close](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/NfcV/close) | 断开与NFCV标签之间的连接 | 5.14.0 | **X** | **X** | 7.35.0
[NfcV.getMaxTransceiveLength](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/NfcV/getmaxtransceivelength) | 获取最大传输长度 | 5.14.0 | **X** | **X** | 7.35.0
[NfcV.setTimeout](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/NfcV/settimeout) | 设置超时时间 | 5.14.0 | **X** | **X** | 7.35.0
[IsoDep.connect](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/isodep/connect) | 连接 ISO-DEP (ISO 14443-4) 类型的标签 | 7.15.0 | **X** | **X** | 7.35.0
[IsoDep.transceive](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/isodep/transceive) | 发送数据给ISO-DEP类型的标签 | 7.15.0 | **X** | **X** | 7.35.0
[IsoDep.close](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/isodep/close) | 断开与ISO-DEP标签之间的连接 | 7.15.0 | **X** | **X** | 7.35.0
[IsoDep.getMaxTransceiveLength](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/isodep/getmaxtransceivelength) | 获取最大传输长度 | 7.15.0 | **X** | **X** | 7.35.0
[IsoDep.setTimeout](https://open.feishu.cn/document/uYjL24iN/uQzM4YjL0MDO24CNzgjN/isodep/settimeout) | 设置超时时间 | 7.15.0 | **X** | **X** | 7.35.0
[MifareClassic.connect](https://open.feishu.cn/document/uYjL24iN/uQTN4YjL0UDO24CN1gjN) | 连接MifareClassic类型的标签 | 3.44.0 | **X** | **X** | 7.35.0
[MifareClassic.transceive](https://open.feishu.cn/document/uYjL24iN/ucTN4YjL3UDO24yN1gjN) | 发送数据给MifareClassic类型的标签 | 3.44.0 | **X** | **X** | 7.35.0
[MifareClassic.close](https://open.feishu.cn/document/uYjL24iN/uMTN4YjLzUDO24yM1gjN) | 断开与MifareClassic标签之间的连接 | 3.44.0 | **X** | **X** | 7.35.0
[MifareClassic.getMaxTransceiveLength](https://open.feishu.cn/document/uYjL24iN/uUTN4YjL1UDO24SN1gjN) | 获取最大传输长度 | 3.44.0 | **X** | **X** | 7.35.0
[MifareClassic.setTimeout](https://open.feishu.cn/document/uYjL24iN/uYTN4YjL2UDO24iN1gjN) | 设置超时时间 | 3.44.0 | **X** | **X** | 7.35.0

#### 蓝牙

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[openBluetoothAdapter](https://open.feishu.cn/document/uYjL24iN/ugzNxYjL4cTM24CO3EjN) | 初始化蓝牙模块 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[closeBluetoothAdapter](https://open.feishu.cn/document/uYjL24iN/uYDOxYjL2gTM24iN4EjN) | 关闭蓝牙模块。调用该方法将断开所有已建立的连接并释放系统资源。建议在使用蓝牙流程后，与 tt.openBluetoothAdapter 成对调用 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[getBluetoothAdapterState](https://open.feishu.cn/document/uYjL24iN/uUDOxYjL1gTM24SN4EjN) | 获取本机蓝牙适配器状态 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[startBluetoothDevicesDiscovery](https://open.feishu.cn/document/uYjL24iN/uUzNxYjL1cTM24SN3EjN) | 开始搜寻附近的蓝牙外围设备 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[stopBluetoothDevicesDiscovery](https://open.feishu.cn/document/uYjL24iN/uczNxYjL3cTM24yN3EjN) | 停止搜寻附近的蓝牙外围设备。若已经找到需要的蓝牙设备并不需要继续搜索时，建议调用该接口停止蓝牙搜索 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[getConnectedBluetoothDevices](https://open.feishu.cn/document/uYjL24iN/uMDOxYjLzgTM24yM4EjN) | 根据 uuid 获取处于已连接状态的设备 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[getBluetoothDevices](https://open.feishu.cn/document/uYjL24iN/uQDOxYjL0gTM24CN4EjN) | 获取在蓝牙模块生效期间所有已发现的蓝牙设备。包括已经和本机处于连接状态的设备 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[onBluetoothDeviceFound](https://open.feishu.cn/document/uYjL24iN/ukzNxYjL5cTM24SO3EjN) | 监听蓝牙设备获取的状态变化 | 5.16.0 | 5.16.0 | **X** | 7.35.0
[offBluetoothDeviceFound](https://open.feishu.cn/document/uYjL24iN/uEDOxYjLxgTM24SM4EjN) | 取消监听蓝牙设备获取的状态变化 | 5.16.0 | 5.16.0 | **X** | 7.35.0
[onBluetoothAdapterStateChange](https://open.feishu.cn/document/uYjL24iN/uADOxYjLwgTM24CM4EjN) | 监听蓝牙适配器的状态变化 | 5.16.0 | 5.16.0 | **X** | 7.35.0
[offBluetoothAdapterStateChange](https://open.feishu.cn/document/uYjL24iN/uIDOxYjLygTM24iM4EjN) | 取消监听蓝牙适配器的状态变化 | 5.16.0 | 5.16.0 | **X** | 7.35.0

#### 低功耗蓝牙

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[setBLEMTU](https://open.feishu.cn/document/uYjL24iN/uMTMyYjLzEjM24yMxIjN) | 设置蓝牙最大传输单元 | 7.3.0 | **X** | **X** | 7.35.0
[readBLECharacteristicValue](https://open.feishu.cn/document/uYjL24iN/uYTOxYjL2kTM24iN5EjN) | 读取数据 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[writeBLECharacteristicValue](https://open.feishu.cn/document/uYjL24iN/ucTOxYjL3kTM24yN5EjN) | 写入蓝牙数据 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[getBLEDeviceCharacteristics](https://open.feishu.cn/document/uYjL24iN/ukDOxYjL5gTM24SO4EjN) | 获取读写特征 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[connectBLEDevice](https://open.feishu.cn/document/uYjL24iN/ucDOxYjL3gTM24yN4EjN) | 低功耗蓝牙链接外围设备 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[disconnectBLEDevice](https://open.feishu.cn/document/uYjL24iN/ugDOxYjL4gTM24CO4EjN) | 断开设备连接 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[getBLEDeviceServices](https://open.feishu.cn/document/uYjL24iN/uATOxYjLwkTM24CM5EjN) | 低功耗蓝牙获取设备服务 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[onBLEConnectionStateChange](https://open.feishu.cn/document/uYjL24iN/uUTOxYjL1kTM24SN5EjN) | 低功耗蓝牙获取设备服务 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[offBLEConnectionStateChange](https://open.feishu.cn/document/uYjL24iN/uMTOxYjLzkTM24yM5EjN) | 取消监听蓝牙低功耗连接状态的改变事件 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[notifyBLECharacteristicValueChange](https://open.feishu.cn/document/uYjL24iN/uETOxYjLxkTM24SM5EjN) | 监听特征值数据变化 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[onBLECharacteristicValueChange](https://open.feishu.cn/document/uYjL24iN/uQTOxYjL0kTM24CN5EjN) | 监听特征值数据变化 | 7.3.0 | 7.3.0 | **X** | 7.35.0
[offBLECharacteristicValueChange](https://open.feishu.cn/document/uYjL24iN/uITOxYjLykTM24iM5EjN) | 取消监听蓝牙低功耗设备的特征值变化事件 | 7.3.0 | 7.3.0 | **X** | 7.35.0

#### Beacon

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[startBeaconDiscovery](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/ibeacon/startbeacondiscovery) | 开始搜索附近的 iBeacon 设备 | 4.6.0 | 4.6.0 | **X** | 7.35.0
[stopBeaconDiscovery](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/ibeacon/stopbeacondiscovery) | 停止搜索附近的 iBeacon 设备 | 4.6.0 | 4.6.0 | **X** | 7.35.0
[getBeacons](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/ibeacon/getbeacons) | 获取所有已搜索到的 iBeacon 设备 | 4.6.0 | 4.6.0 | **X** | 7.35.0
[onBeaconUpdate](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/ibeacon/onbeaconupdate) | 监听 iBeacon 设备更新事件，仅能注册一个监听 | 4.6.0 | 4.6.0 | **X** | 7.35.0
[offBeaconUpdate](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/ibeacon/offbeaconupdate) | 取消监听 iBeacon 设备更新事件 | 4.6.0 | 4.6.0 | **X** | 7.35.0
[onBeaconServiceChange](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/ibeacon/onbeaconservicechange) | 监听蓝牙适配器状态变化事件 | 4.6.0 | 4.6.0 | **X** | 7.35.0
[offBeaconServiceChange](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/ibeacon/offbeaconservicechange) | 取消监听 iBeacon 服务状态变化事件 | 4.6.0 | 4.6.0 | **X** | 7.35.0

#### 扫码

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[scanCode](https://open.feishu.cn/document/uYjL24iN/uYzNx4iN3EjL2cTM) | 扫描二维码并返回扫描结果 | 3.44.0 | 3.44.0 | **X** | 7.35.0

#### Wi-Fi

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getConnectedWifi](https://open.feishu.cn/document/uYjL24iN/ugjNx4CO2EjL4YTM) | 获取设备当前所连的 Wifi | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[getWifiStatus](https://open.feishu.cn/document/uYjL24iN/uYTN4QjL2UDO04iN1gDN) | 请求获取 Wi-Fi 开关状态 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[getWifiList](https://open.feishu.cn/document/uYjL24iN/uUDO4UjL1gDO14SN4gTN) | 请求获取Wifi 列表 | 3.44.0 | **X** | **X** | 7.35.0
[onGetWifiList](https://open.feishu.cn/document/uYjL24iN/uYDO4UjL2gDO14iN4gTN) | 监听获取到 Wi-Fi 列表数据事件 | 3.44.0 | **X** | **X** | 7.35.0
[offGetWifiList](https://open.feishu.cn/document/uYjL24iN/ucDO4UjL3gDO14yN4gTN) | 取消监听获取到 Wi-Fi 列表数据事件 | 3.44.0 | **X** | **X** | 7.35.0

#### 剪贴板

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[setClipboardData](https://open.feishu.cn/document/uYjL24iN/ugzNx4CO3EjL4cTM) | 设置系统剪贴板内容 | 3.44.0 | 3.44.0 | 3.44.0 | 7.35.0
[getClipboardData](https://open.feishu.cn/document/uYjL24iN/uczNx4yN3EjL3cTM) | 获取系统粘贴板数据 | 3.44.0 | 3.44.0 | 3.44.0 | 7.35.0

#### 网络状态

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getNetworkType](https://open.feishu.cn/document/uYjL24iN/uYjNx4iN2EjL2YTM) | 获取设备当前所处的网络类型 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[onNetworkStatusChange](https://open.feishu.cn/document/uYjL24iN/ucjNx4yN2EjL3YTM) | 监听网络状态变化 | 7.7.0 | 7.7.0 | **X** | 7.35.0
[getNetworkQualityType](https://open.feishu.cn/document/uYjL24iN/uUTNx4SN1EjL1UTM/getnetworkqualitytype) | 网络评级接口，获取当前设备所处的网络状态 | 4.9.0 | 4.9.0 | 5.1.0 | 7.35.0
[onNetworkQualityChange](https://open.feishu.cn/document/uYjL24iN/uUTNx4SN1EjL1UTM/onnetworkqualitychange) | 监听网络质量变化 | 4.9.0 | 4.9.0 | 5.1.0 | 7.35.0

#### 屏幕亮度

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[setScreenBrightness](https://open.feishu.cn/document/uYjL24iN/uIjNx4iM2EjLyYTM/set-screen-brightness) | 设置屏幕亮度 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[getScreenBrightness](https://open.feishu.cn/document/uYjL24iN/uIjNx4iM2EjLyYTM/get-screen-brightness) | 获取屏幕亮度 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[setKeepScreenOn](https://open.feishu.cn/document/uYjL24iN/ukzNx4SO3EjL5cTM) | 设置是否保持常亮状态 | 3.44.0 | 3.44.0 | **X** | 7.35.0

#### 截屏监听

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[grantScreenShot](https://open.feishu.cn/document/uYjL24iN/uIjNwEjLyYDMx4iM2ATM/grantscreenshot) | 用于开启/禁用截屏，仅在Android和Harmony设备上生效。 | 7.7.0 | **X** | **X** | 7.35.0
[onUserCaptureScreen](https://open.feishu.cn/document/uYjL24iN/uMjNwEjLzYDMx4yM2ATM) | 监听用户主动截屏事件。用户使用系统截屏按键截屏时触发 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[offUserCaptureScreen](https://open.feishu.cn/document/uYjL24iN/uQjNwEjL0YDMx4CN2ATM) | 取消监听用户主动截屏事件 | 3.44.0 | 3.44.0 | **X** | 7.35.0

#### 加速度计

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[startAccelerometer](https://open.feishu.cn/document/uYjL24iN/ukjNx4SO2EjL5YTM) | 通知客户端开始监听加速度计数据。具体的数据返回通过注册onAccelerometerChange接口回调方法获取 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[stopAccelerometer](https://open.feishu.cn/document/uYjL24iN/uAzNx4CM3EjLwcTM) | 停止监听加速度计数据 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[onAccelerometerChange](https://open.feishu.cn/document/uYjL24iN/uEzNx4SM3EjLxcTM) | 监听加速度计数据 | 3.44.0 | 3.44.0 | **X** | 7.35.0

#### 拨打电话

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[makePhoneCall](https://open.feishu.cn/document/uYjL24iN/uUzNx4SN3EjL1cTM) | 拨打电话 | 3.44.0 | 3.44.0 | **X** | 7.35.0

#### 震动

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[vibrateShort](https://open.feishu.cn/document/uYjL24iN/uADOx4CM4EjLwgTM) | 使手机发生较短时间的振动 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[vibrateLong](https://open.feishu.cn/document/uYjL24iN/uEDOx4SM4EjLxgTM) | 使手机发生较长时间的振动 | 3.44.0 | 3.44.0 | **X** | 7.35.0

#### 罗盘

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[startCompass](https://open.feishu.cn/document/uYjL24iN/uIzNx4iM3EjLycTM) | 开始监听罗盘数据。 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[stopCompass](https://open.feishu.cn/document/uYjL24iN/uMzNx4yM3EjLzcTM) | 停止监听罗盘数据 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[onCompassChange](https://open.feishu.cn/document/uYjL24iN/uQzNx4CN3EjL0cTM) | 监听罗盘数据状态变化 | 5.16.0 | 5.16.0 | **X** | 7.35.0

#### 性能

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[onMemoryWarning](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/performance/onmemorywarning) | 监听内存不足的告警事件。当手机内存占用过高时，触发回调函数。 | 7.10.0 | 7.10.0 | **X** | **X**

## 文件

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[saveFile](https://open.feishu.cn/document/uYjL24iN/ugDOz4CO4MjL4gzM) | 保存临时文件到本地永久目录 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[saveFileAs](https://open.feishu.cn/document/uYjL24iN/uQjN3UjL0YzN14CN2cTN) | 保存文件到本地指定目录 | **X** | **X** | 5.16.0 | **X**
[filePicker](https://open.feishu.cn/document/uYjL24iN/uETM04SMxQjLxEDN) | 打开附件选择列表 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[docsPicker](https://open.feishu.cn/document/uYjL24iN/ukTN3UjL5UzN14SO1cTN) | 打开云文档选择列表 | 3.44.0 | 3.44.0 | 3.47.0 | 7.49.0
[openDocument](https://open.feishu.cn/document/uYjL24iN/ukTN24SO1YjL5UjN) | 在新页面打开文档 | 3.44.0 | 3.44.0 | 3.44.0 | 7.39.0
[getFileSystemManager](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/tt_get_file_system_manager) | 获取全局唯一的文件管理器 | 4.11.0 | 4.11.0 | 4.11.0 | 7.35.0
[FileSystemManager.appendFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/appendfile) | 在文件结尾追加内容 | 5.23.0 | 5.23.0 | 5.24.0 | 7.35.0
[FileSystemManager.access](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_access) | 判断文件/目录是否存在 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.copyFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_copy_file) | 复制本地文件/目录 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.getFileInfo](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_get_file_info) | 获取本地文件信息 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.getSavedFileList](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_get_saved_file_list) | 获取用户目录内文件列表 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.mkdir](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_mkdir) | 创建本地目录 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.readdir](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_read_dir) | 读取本地目录内文件列表 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.readFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_read_file) | 读取本地文件内容 | 4.11.0 | 4.11.0 | 5.23.0 | 7.35.0
[FileSystemManager.removeSavedFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_remove_saved_file) | 删除保存的本地文件 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.rename](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_rename) | 重命名本地文件/目录 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.rmdir](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_rmdir) | 删除本地目录 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.saveFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_save_file) | 保存临时文件到本地用户目录 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.stat](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_stat) | 获取本地文件 Stats 对象 | 4.11.0 | 4.11.0 | 5.23.0 | 7.39.0
[FileSystemManager.unlink](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_unlink) | 删除本地文件 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.writeFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_write_file) | 写入本地文件 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[Stats.isDirectory](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/stat/stats_is_directory) | 判断当前文件是否一个目录 | 4.11.0 | 4.11.0 | **X** | 7.39.0
[Stats.isFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/stat/stats_is_file) | 判断当前文件是否一个普通文件 | 4.11.0 | 4.11.0 | **X** | 7.39.0

## 媒体
#### 图片

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[saveImageToPhotosAlbum](https://open.feishu.cn/document/uYjL24iN/uUTMx4SNxEjL1ETM) | 保存图片到系统相册 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[chooseImage](https://open.feishu.cn/document/uYjL24iN/uMTMx4yMxEjLzETM) | 从系统相册中选择图片，或使用相机拍摄图片 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[previewImage](https://open.feishu.cn/document/uYjL24iN/uMDOx4yM4EjLzgTM) | 预览一组图片 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[compressImage](https://open.feishu.cn/document/uYjL24iN/uMjN24yM2YjLzYjN) | 压缩图片接口，可选压缩质量 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[getImageInfo](https://open.feishu.cn/document/uYjL24iN/ugjNwEjL4YDMx4CO2ATM) | 获取图片信息 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0

#### 视频

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getVideoInfo](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/getvideoinfo) | 获取视频详细信息 | 6.5.0 | 6.5.0 | 6.5.0 | 7.35.0
[saveVideoToPhotosAlbum](https://open.feishu.cn/document/uYjL24iN/ucDOx4yN4EjL3gTM) | 保存视频到系统相册 | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[chooseVideo](https://open.feishu.cn/document/uYjL24iN/uEjMx4SMyEjLxITM) | 从系统相册中选择视频，或使用相机拍摄视频 | 3.44.0 | 3.44.0 | 3.47.0 | **X**
[chooseMedia](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/choosemedia) | 拍摄或从系统相册中选择图片或视频 | 4.7.0 | 4.7.0 | 4.7.0 | 7.35.0

#### 录音

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getRecorderManager](https://open.feishu.cn/document/uYjL24iN/uQDOx4CN4EjL0gTM) | 获取全局唯一的 recorderManager 。通过 recorderManager 进行录音操作和管理 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.start](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/start) | 开始录音 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.pause](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/pause) | 暂停录音 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.resume](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/resume) | 继续录音 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.stop](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/stop) | 停止录音 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.getState](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/getstate) | 获取当前录音状态 | 7.13.0 | 7.13.0 | **X** | 7.43.0
[RecorderManager.onStart](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onstart) | 监听录音开始事件 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.onPause](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onpause) | 监听录音暂停事件 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.onResume](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onresume) | 监听录音继续事件 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.onStop](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onstop) | 监听录音结束事件 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.onFrameRecorded](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onframerecorded) | 监听已录制完指定帧大小的文件事件 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.onError](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onerror) | 监听录音错误事件 | **✓** | **✓** | **X** | 7.43.0

## 导航

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[openSchema](https://open.feishu.cn/document/uYjL24iN/ukzN4IjL5cDOy4SO3gjM) | 跳转到小程序以外的应用（如云文档、网页等） | 3.44.0 | 3.44.0 | 3.47.0 | 7.35.0
[closeWindow](https://open.feishu.cn/document/uYjL24iN/uYTOuYTOuYTO/closewindow) | 关闭当前窗口 | 5.14.0 | 5.14.0 | 5.14.0 | 7.41.0

## 地理位置

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getLocationStatus](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/getlocationstatus) | 获取 GPS 开关状态 | 6.10.0 | 6.10.0 | **X** | 7.35.0
[reverseGeocode](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/reversegeocode) | 逆地理编码，根据传入的经纬度，获取相应的位置信息 | 6.6.0 | 6.6.0 | **X** | 7.35.0
[offLocationChangeError](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/offlocationchangeerror) | 取消监听持续定位接口返回失败时触发 | 5.21.0 | 5.21.0 | **X** | 7.35.0
[onLocationChangeError](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/onlocationchangeerror) | 监听持续定位接口返回失败时触发 | 5.21.0 | 5.21.0 | **X** | 7.35.0
[openLocation](https://open.feishu.cn/document/uYjL24iN/uQTOz4CN5MjL0kzM) | 使用客户端内置地图查看位置 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[chooseLocation](https://open.feishu.cn/document/uYjL24iN/uUDN1EjL1QTNx4SN0UTM) | 打开地图选择位置 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[getLocation](https://open.feishu.cn/document/uYjL24iN/uUTOz4SN5MjL1kzM) | 获取设备当前的地理位置 | 3.44.0 | 3.44.0 | **X** | 7.35.0
[startLocationUpdate](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/startlocationupdate) | 开启小程序在前台时接收位置更新消息 | 5.2.0 | 5.2.0 | **X** | 7.35.0
[stopLocationUpdate](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/stoplocationupdate) | 关闭监听实时位置变化，停止位置消息接收 | 5.2.0 | 5.2.0 | **X** | 7.35.0
[onLocationChange](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/onlocationchange) | 监听实时地理位置变化事件，需结合startLocationUpdate使用 | 5.2.0 | 5.2.0 | **X** | 7.35.0
[offLocationChange](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/offlocationchange) | 取消监听实时地理位置变化事件 | 5.2.0 | 5.2.0 | **X** | 7.35.0

## 运动与健康

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getStepCount](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/sport-and-health/getstepcount) | 获取设备中的步数信息 | 5.1.0 | 5.1.0 | **X** | 7.35.0

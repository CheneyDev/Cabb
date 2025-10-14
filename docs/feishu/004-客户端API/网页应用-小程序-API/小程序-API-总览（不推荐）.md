---
title: "小程序 API 总览（不推荐）"
source_url: https://open.feishu.cn/document/client-docs/gadget/gadget-api-list
last_remote_update: 2025-08-01
last_remote_update_timestamp: 1754044068000
---
最后更新于 2025-08-01

# 小程序 API 总览

飞书扫码在线预览 | 示例代码
--- | ---
![图片名称](https://sf3-cn.feishucdn.com/obj/website-img/acb2395e3e575b81f2352bd209525b10_dnwbLGxD33.png?height=155&lazyload=true&width=155) | - 下载示例代码：[microapp-demo.zip](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/c87f497f58a7c8b44b39e8112311c567_MzCYmZjOxE.zip)   
- 导入并调试示例代码：[使用说明](https://open.feishu.cn/document/uYjL24iN/uYDM04iNwQjL2ADN)

## 开放接口
#### 登录

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[requestAccess](https://open.feishu.cn/document/uYjL24iN/uUzMuUzMuUzM/requestaccess) | 增量授予应用访问权限 | 6.9.0 | 6.9.0 | 6.9.0 | 7.35.0
[checkSession](https://open.feishu.cn/document/uYjL24iN/ukTMx4SOxEjL5ETM) | 检查用户当前的 session 状态是否有效 | **✓** | **✓** | **✓** | 7.35.0

#### 用户信息

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[enterProfile](https://open.feishu.cn/document/uYjL24iN/ucDM04yNwQjL3ADN) | 打开个人信息主页 | **✓** | **✓** | **✓** | 7.35.0
[getUserInfo](https://open.feishu.cn/document/uYjL24iN/ucjMx4yNyEjL3ITM) | 获取已登录用户的基本信息或特殊信息 | **✓** | **✓** | **✓** | 7.35.0

#### 聊天

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[enterChat](https://open.feishu.cn/document/uYjL24iN/ukDM04SOwQjL5ADN) | 打开指定会话 | **✓** | **✓** | **✓** | 7.35.0
[toggleChat](https://open.feishu.cn/document/uYjL24iN/ugDM04COwQjL4ADN/toggleChat) | 侧边栏形式打开或关闭会话 | **X** | **X** | **✓** | **X**
[chooseChat](https://open.feishu.cn/document/uYjL24iN/uMTN3QjLzUzN04yM1cDN) | 打开用户会话列表选择会话，调用前确保用户已经登入 | 3.1.0 | 3.1.0 | 3.1.0 | 7.35.0
[getChatInfo](https://open.feishu.cn/document/uYjL24iN/uEDN2UjLxQjN14SM0YTN) | 获取某个会话的信息 | 3.10.0 | 3.10.0 | 3.10.0 | 7.35.0
[getBlockActionSourceDetail](https://open.feishu.cn/document/getBlockActionSourceDetail) | 支持从block action点击进入应用后，获取block对应业务的详细信息 | **✓** | **✓** | **✓** | **X**
[getTriggerContext](https://open.feishu.cn/document/no_class/plusmenu/gettriggercontext) | 加号菜单场景中，获取当前会话信息 | **✓** | **✓** | **✓** | **X**
[enterBot](https://open.feishu.cn/document/uYjL24iN/uAjM1EjLwITNx4CMyUTM) | 打开机器人聊天页面 | 2.7.0 | 2.7.0 | 2.7.0 | 7.35.0
[sendMessageCard](https://open.feishu.cn/document/uYjL24iN/uUjN5UjL1YTO14SN2kTN) | 发送消息卡片到指定会话 | 3.19.0 | 3.19.0 | 3.19.0 | 7.39.0
[onChatBadgeChange](https://open.feishu.cn/document/uYjL24iN/uQDN2UjL0QjN14CN0YTN) | 监听某个群未读消息数变化，确保用户已经登录 | 3.10.0 | 3.10.0 | 3.10.0 | **X**
[offChatBadgeChange](https://open.feishu.cn/document/uYjL24iN/ugDM04COwQjL4ADN/offchatbadgechange) | 取消监听某个群未读消息数变化 | 3.10.0 | 3.10.0 | 3.10.0 | **X**

#### 联系人

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[chooseContact](https://open.feishu.cn/document/uYjL24iN/uMTM04yMxQjLzEDN) | 打开用户联系人选择列表 | **✓** | **✓** | **✓** | 7.35.0

#### 设置

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[openSetting](https://open.feishu.cn/document/uYjL24iN/uUzMx4SNzEjL1MTM) | 打开设置页面，展示用户设置（包括授予和拒绝）过的权限，并返回用户设置过的授权结果 | **✓** | **✓** | **✓** | 7.35.0
[getSetting](https://open.feishu.cn/document/uYjL24iN/uQzMx4CNzEjL0MTM) | 获取用户设置（包括授予和拒绝）过的权限 | **✓** | **✓** | **✓** | 7.35.0

#### 分享

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[shareWebContent](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/share/sharewebcontent) | 分享应用内网页内容 | 5.23.0 | 5.23.0 | 5.23.0 | 7.39.0
[share](https://open.feishu.cn/document/uYjL24iN/ugDM04COwQjL4ADN/thirdShare) | 分享内容到三方应用 | 3.47.0 | 3.47.0 | **X** | **X**
[showShareMenu](https://open.feishu.cn/document/uYjL24iN/ugjN24CO2YjL4YjN) | 显示当前页面的分享按钮 | 2.2.0 | 2.2.0 | 2.2.0 | 7.35.0
[hideShareMenu](https://open.feishu.cn/document/uYjL24iN/ukjN24SO2YjL5YjN) | 隐藏当前页面的分享按钮 | 2.2.0 | 2.2.0 | 2.2.0 | 7.35.0

#### 授权

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[authorize](https://open.feishu.cn/document/uYjL24iN/ugzMx4COzEjL4MTM) | 向用户发出设置权限请求。如果该权限用户没有设置过，会弹窗咨询用户是否授予；如果该权限用户拒绝授予，会打开设置页面(appBadge权限除外)；如果该权限用户同意授予，会直接返回成功 | **✓** | **✓** | **✓** | 7.35.0
[getTenantAppScopes](https://open.feishu.cn/document/uYjL24iN/uczMx4yNzEjL3MTM/gettenantappscopes) | 查询租户下当前应用的授权状态 | 3.43.0 | 3.43.0 | 3.43.0 | 7.35.0
[applyTenantAppScope](https://open.feishu.cn/document/uYjL24iN/uczMx4yNzEjL3MTM/applytenantappscope) | 弹窗咨询用户是否向租户管理员申请所有未授予权限(不包括租户敏感权限) | 3.43.0 | 3.43.0 | 3.43.0 | 7.35.0

#### 飞书启动参数

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getLaunchOptionsSync](https://open.feishu.cn/document/uYjL24iN/uAzM1YjLwMTN24CMzUjN) | 获取小程序启动时的参数 | 3.22.0 | 3.22.0 | 3.22.0 | 7.35.0

#### 安全密码验证

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[startPasswordVerify](https://open.feishu.cn/document/uYjL24iN/ugTO3IjL4kzNy4CO5cjM) | 调起二次验证飞书安全密码的输入界面 | 3.1.0 | 3.1.0 | 3.1.0 | 7.35.0

#### 系统认证

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[startDeviceCredential](https://open.feishu.cn/document/uYjL24iN/uIDN14iM0UjLyQTN) | 打开系统解锁界面 | **✓** | **✓** | **X** | 7.37.0

#### 水印

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[checkWatermark](https://open.feishu.cn/document/uYjL24iN/ukTM1EjL5ETNx4SOxUTM) | 查看宿主是否显示了全局水印 | 2.7.0 | 2.7.0 | 2.7.0 | 7.35.0

#### 邮件

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[mailto](https://open.feishu.cn/document/uYjL24iN/uAjNwEjLwYDMx4CM2ATM) | 调用系统发送邮件 | 2.2.0 | 2.2.0 | **X** | 7.35.0

#### 应用角标

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[reportBadge](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/reportbadge) | 上报应用内的角标（badge）数量 | 5.9.0 | 5.9.0 | 5.9.0 | **X**
[updateBadge](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/updagebadge) | 更新当前设备上的应用角标数 | 5.9.0 | 5.9.0 | 5.9.0 | **X**
[onServerBadgePush](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/onserverbadgepush) | 监听应用的红点（Badge）推送 | 5.9.0 | 5.9.0 | 5.9.0 | **X**
[offServerBadgePush](https://open.feishu.cn/document/uYjL24iN/uQjMuQjMuQjM/appbadge/offserverbadgepush) | 取消对红点（Badge）推送的监听 | 5.9.0 | 5.9.0 | 5.9.0 | **X**

## 界面
#### 交互反馈

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[showActionSheet](https://open.feishu.cn/document/uYjL24iN/ukDNy4SO0IjL5QjM) | 显示操作菜单 | **✓** | **✓** | **✓** | 7.35.0
[showModal](https://open.feishu.cn/document/uYjL24iN/ugDNy4CO0IjL4QjM) | 显示模态弹窗 | **✓** | **✓** | **✓** | 7.35.0
[showPrompt](https://open.feishu.cn/document/uYjL24iN/uYTO4UjL2kDO14iN5gTN) | 展示可输入内容的弹窗 | 3.17.0 | 3.17.0 | 3.17.0 | 7.35.0
[showLoading](https://open.feishu.cn/document/uYjL24iN/uMDNy4yM0IjLzQjM) | 显示灰色背景的 loading 提示框 | **✓** | **✓** | **✓** | 7.35.0
[hideLoading](https://open.feishu.cn/document/uYjL24iN/uYDNy4iN0IjL2QjM) | 隐藏 loading 提示框 | **✓** | **✓** | **✓** | 7.35.0
[showToast](https://open.feishu.cn/document/uYjL24iN/ugzMy4COzIjL4MjM) | 显示灰色背景的消息提示框 | **✓** | **✓** | **✓** | 7.35.0
[hideToast](https://open.feishu.cn/document/uYjL24iN/ukzMy4SOzIjL5MjM) | 隐藏灰色背景的消息提示框 | **✓** | **✓** | **✓** | 7.35.0

#### Tab Bar

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[showTabBar](https://open.feishu.cn/document/uYjL24iN/uATN04CM1QjLwUDN) | 显示 tabBar | **✓** | **✓** | **✓** | 7.35.0
[hideTabBar](https://open.feishu.cn/document/uYjL24iN/ukDN04SO0QjL5QDN) | 隐藏 tabBar | **✓** | **✓** | **✓** | 7.35.0
[setTabBarItem](https://open.feishu.cn/document/uYjL24iN/uETN04SM1QjLxUDN) | 动态设置 tabBar 某一项的内容 | **✓** | **✓** | **✓** | 7.35.0
[addTabBarItem](https://open.feishu.cn/document/uYjL24iN/uQjM04CNyQjL0IDN/addtabbaritem) | 当前小程序的tab bar数量进行增加调整 | 5.1.0 | 5.1.0 | 5.1.0 | 7.35.0
[removeTabBarItem](https://open.feishu.cn/document/uYjL24iN/uQjM04CNyQjL0IDN/removetabbaritem) | 删除tab bar的目标item | 4.2.0 | 4.2.0 | 4.2.0 | 7.35.0
[setTabBarStyle](https://open.feishu.cn/document/uYjL24iN/uITN04iM1QjLyUDN) | 动态设置 tabBar 的整体样式 | **✓** | **✓** | **✓** | 7.35.0
[setTabBarBadge](https://open.feishu.cn/document/uYjL24iN/uUjM04SNyQjL1IDN) | 为 tabBar 某一项的右上角添加文本 | **✓** | **✓** | **✓** | 7.35.0
[removeTabBarBadge](https://open.feishu.cn/document/uYjL24iN/ucjM04yNyQjL3IDN) | 移除 tabBar 某一项右上角的文本 | **✓** | **✓** | **✓** | 7.35.0
[showTabBarRedDot](https://open.feishu.cn/document/uYjL24iN/uYjM04iNyQjL2IDN) | 显示 tabBar 某一项的右上角的红点，可以使用底部标签栏红点给予用户提示 | **✓** | **✓** | **✓** | 7.35.0
[hideTabBarRedDot](https://open.feishu.cn/document/uYjL24iN/ugjM04COyQjL4IDN) | 隐藏 tabBar 某一项的右上角的红点 | **✓** | **✓** | **✓** | 7.35.0

#### 导航栏

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[setNavigationBarTitle](https://open.feishu.cn/document/uYjL24iN/uATNy4CM1IjLwUjM) | 设置导航栏标题 | **✓** | **✓** | **✓** | 7.35.0
[setNavigationBarColor](https://open.feishu.cn/document/uYjL24iN/ucjMyUjL3IjM14yNyITN) | 设置导航栏颜色 | 3.7.0 | 3.7.0 | 3.7.0 | 7.35.0

#### 窗口

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[setWindowSize](https://open.feishu.cn/document/uYjL24iN/uEDO3UjLxgzN14SM4cTN/setwindowsize) | 小程序在 window 和 window-semi 模式下调整独立窗口的大小和位置 | **X** | **X** | 4.3.0 | **X**
[onWindowResize](https://open.feishu.cn/document/uYjL24iN/uADO3UjLwgzN14CM4cTN) | 监听窗口尺寸变化事件 | **✓** | **✓** | 3.13.0 | 7.35.0
[offWindowResize](https://open.feishu.cn/document/uYjL24iN/uIDO3UjLygzN14iM4cTN) | 取消监听窗口尺寸变化事件 | **✓** | **✓** | 3.13.0 | 7.35.0
[setPageOrientation](https://open.feishu.cn/document/uYjL24iN/uEDO3UjLxgzN14SM4cTN/setpageorientation) | 设置当前小程序页面方向 | 5.26.0 | 5.26.0 | **X** | 7.35.0
[onResize](https://open.feishu.cn/document/uYjL24iN/uEDO3UjLxgzN14SM4cTN/onresize) | 监听页面尺寸变化事件 | 5.26.0 | 5.26.0 | 5.26.0 | 7.35.0

#### 下拉刷新

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[Page.onPullDownRefresh](https://open.feishu.cn/document/uYjL24iN/uQTNy4CN1IjL0UjM) | 在 Page 中注册下拉刷新的监听方法，当用户触发下拉刷新时会调用 | **✓** | **✓** | **X** | 7.35.0
[startPullDownRefresh](https://open.feishu.cn/document/uYjL24iN/uYTNy4iN1IjL2UjM) | 下拉刷新 | **✓** | **✓** | **X** | 7.35.0
[stopPullDownRefresh](https://open.feishu.cn/document/uYjL24iN/ugTNy4CO1IjL4UjM) | 停止当前页面下拉刷新 | **✓** | **✓** | **X** | 7.35.0

#### 页面离开二次确认

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[enableLeaveConfirm](https://open.feishu.cn/document/uYjL24iN/uUTOuUTOuUTO/page_leave_confirm_modal/enable_leave_confirm) | 当某个页面返回或关闭时设置二次弹窗提示 | 5.17.0 | 5.16.0 | 5.13.0  
      7.39.0
[disableLeaveConfirm](https://open.feishu.cn/document/uYjL24iN/uUTOuUTOuUTO/page_leave_confirm_modal/disable_leave_confirm) | 取消当前页面离开二次确认弹框 | 5.17.0 | 5.16.0 | 5.13.0  
      7.39.0

#### 页面

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[pageScrollTo](https://open.feishu.cn/document/uYjL24iN/uITNy4iM1IjLyUjM) | 滚动页面到目标位置 | **✓** | **✓** | **✓** | 7.35.0

#### Canvas绘图

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[createCanvasContext](https://open.feishu.cn/document/uYjL24iN/uMTNy4yM1IjLzUjM) | 创建并返回对应 canvasId 的绘图上下文 | **✓** | **✓** | **✓** | 7.35.0
[canvasToTempFilePath](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvastotempfilepath) | 导出当前画布指定区域，生成图片并返回文件路径 | **✓** | **✓** | **✓** | 7.35.0
[canvasPutImageData](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvasputimagedata) | 更新画布像素数据 | 3.45.0 | 3.45.0 | 3.45.0 | 7.35.0
[canvasGetImageData](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvasgetimagedata) | 获取画布像素数据 | 3.45.0 | 3.45.0 | 3.45.0 | 7.35.0
[CanvasContext.createPattern](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-createPattern) | 创建径向渐变管理对象 | 3.45.0 | 3.45.0 | 3.45.0 | 7.35.0
[CanvasContext.createCircularGradient](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-createCircularGradient) | 创建圆形渐变管理对象 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.translate](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-translate) | 平移坐标矩阵 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.transform](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-transform) | 坐标转换矩阵叠加，每一次调用会在乘以前一次的变换矩阵 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.strokeText](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-strokeText) | 绘制文字路径 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.strokeRect](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-strokeRect) | 绘制矩形路径，不添加到当前路径中 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.stroke](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-stroke) | 绘制当前路径 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setTransform](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setTransform) | 设置坐标转换矩阵 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setTextBaseline](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setTextBaseline) | 设置字体的对齐基线 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setTextAlign](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setTextAlign) | 设置字体对齐方式 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setStrokeStyle](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setStrokeStyle) | 设置绘制线样式 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setShadow](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setShadow) | 设置阴影 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setMiterLimit](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setMiterLimit) | 设置线连接点渲染的斜面倾斜程度 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setLineWidth](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setLineWidth) | 设置线宽 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setLineJoin](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setLineJoin) | 设置线连接点样式 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setLineDash](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setLineDash) | 设置间断线 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setLineCap](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setLineCap) | 设置线端点样式 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setGlobalAlpha](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setGlobalAlpha) | 设置全局不透明度 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.setFillStyle](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-setFillStyle) | 设置填充样式 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.scale](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-scale) | 缩放坐标点 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.save](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-save) | 创建新的绘图上下文，并将之前的上下文保存在栈中 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.rotate](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-rotate) | 旋转坐标点 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.restore](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-restore) | 恢复栈中存储的上下文 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.rect](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-rect) | 添加矩形到当前路径中 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.quadraticCurveTo](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-quadraticCurveTo) | 添加二次贝塞尔曲线到路径中 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.moveTo](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-moveTo) | 移动绘制点 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.measureText](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-measureText) | 测量文字宽度 | 3.9.0 | 3.9.0 | 3.9.0 | 7.35.0
[CanvasContext.lineTo](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-lineTo) | 移动并添加线段到路径中 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.fillText](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-fillText) | 填充文字 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.fillRect](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-fillRect) | 填充矩形，不添加到当前路径中 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.fill](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-fill) | 填充当前路径 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.drawImage](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-drawImage) | 绘制 Image | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.createLinearGradient](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-createLinearGradient) | 创建线性渐变对象 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.closePath](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-closePath) | 闭合当前路径 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.clip](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-clip) | 剪切当前路径，限制后续的渲染范围 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.clearRect](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-clearRect) | 清空画布矩形区域 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.bezierCurveTo](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-bezierCurveTo) | 添加三次贝塞尔曲线到路径中 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.beginPath](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-beginPath) | 创建新的子路径 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.arcTo](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-arcTo) | 移动并添加弧线到当前路径中 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.arc](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-arc) | 添加圆弧到当前路径中 | **✓** | **✓** | **✓** | 7.35.0
[CanvasContext.draw](https://open.feishu.cn/document/uYjL24iN/ukjMy4SOyIjL5IjM/canvascontext/canvascontext-draw) | 将所有的操作绘制到 Canvas 中 | **✓** | **✓** | **✓** | 7.35.0

#### 动画

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[createAnimation](https://open.feishu.cn/document/uYjL24iN/uETNy4SM1IjLxUjM) | 创建一个动画实例 animation。调用实例的方法来描述动画。最后通过动画实例的 export 方法导出动画数据传递给组件的 animation 属性 | **✓** | **✓** | **✓** | 7.35.0
[Animation.backgroundColor](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_backgroundcolor) | 设置背景色 | **✓** | **✓** | **✓** | 7.35.0
[Animation.bottom](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_bottom) | 设置 bottom 值 | **✓** | **✓** | **✓** | 7.35.0
[Animation.export](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_export) | 导出动画队列。export 方法每次调用后会清掉之前的动画操作 | **✓** | **✓** | **✓** | 7.35.0
[Animation.height](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_height) | 设置高度 | **✓** | **✓** | **✓** | 7.35.0
[Animation.left](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_left) | 设置 left 值 | **✓** | **✓** | **✓** | 7.35.0
[Animation.matrix](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_matrix) | 同 transform-function matrix | **✓** | **✓** | **✓** | 7.35.0
[Animation.matrix3d](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_matrix3d) | 同 transform-function matrix3d | **✓** | **✓** | **✓** | 7.35.0
[Animation.opacity](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_opacity) | 设置透明度 | **✓** | **✓** | **✓** | 7.35.0
[Animation.right](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_right) | 设置 right 值 | **✓** | **✓** | **✓** | 7.35.0
[Animation.rotate](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_rotate) | 从原点顺时针旋转一个角度 | **✓** | **✓** | **✓** | 7.35.0
[Animation.rotate3d](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_rotate3d) | 从 固定 轴顺时针旋转一个角度 | **✓** | **✓** | **✓** | 7.35.0
[Animation.rotateX](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_rotatex) | 从 X 轴顺时针旋转一个角度 | **✓** | **✓** | **✓** | 7.35.0
[Animation.rotateY](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_rotateY) | 从 Y 轴顺时针旋转一个角度 | **✓** | **✓** | **✓** | 7.35.0
[Animation.rotateZ](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_rotatez) | 从 Z 轴顺时针旋转一个角度 | **✓** | **✓** | **✓** | 7.35.0
[Animation.scale](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_scale) | 缩放 | **✓** | **✓** | **✓** | 7.35.0
[Animation.scale3d](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_scale3d) | 缩放 | **✓** | **✓** | **✓** | 7.35.0
[Animation.scaleX](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_scalex) | 缩放 X 轴 | **✓** | **✓** | **✓** | 7.35.0
[Animation.scaleY](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_scaley) | 缩放 Y 轴 | **✓** | **✓** | **✓** | 7.35.0
[Animation.scaleZ](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_scalez) | 缩放 Z 轴 | **✓** | **✓** | **✓** | 7.35.0
[Animation.skew](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_skew) | 对 X、Y 轴坐标进行倾斜 | **✓** | **✓** | **✓** | 7.35.0
[Animation.skewX](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_skewx) | 对 X 轴坐标进行倾斜 | **✓** | **✓** | **✓** | 7.35.0
[Animation.skewY](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_skewy) | 对 Y 轴坐标进行倾斜 | **✓** | **✓** | **✓** | 7.35.0
[Animation.step](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_step) | 表示一组动画完成。可以在一组动画中调用任意多个动画方法，一组动画中的所有动画会同时开始，一组动画完成后才会进行下一组动画 | **✓** | **✓** | **✓** | 7.35.0
[Animation.top](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_top) | 设置 top 值 | **✓** | **✓** | **✓** | 7.35.0
[Animation.translate](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_translate) | 平移变换 | **✓** | **✓** | **✓** | 7.35.0
[Animation.translate3d](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_translate3d) | 对 xyz 坐标进行平移变换 | **✓** | **✓** | **✓** | 7.35.0
[Animation.translateX](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_translatex) | 对 X 轴平移 | **✓** | **✓** | **✓** | 7.35.0
[Animation.translateY](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_translatey) | 对 Y 轴平移 | **✓** | **✓** | **✓** | 7.35.0
[Animation.translateZ](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_translatez) | 对 Z 轴平移 | **✓** | **✓** | **✓** | 7.35.0
[Animation.width](https://open.feishu.cn/document/uYjL24iN/ucjMy4yNyIjL3IjM/animation/animation_width) | 设置宽度 | **✓** | **✓** | **✓** | 7.35.0

#### Customized Input

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getCustomizedInput](https://open.feishu.cn/document/uYjL24iN/uEDN1EjLxQTNx4SM0UTM) | 获取全局唯一的customizedInput实例。通过customizedInput显示一个 可定制化的富文本输入框，支持@联系人、插入图片、插入表情、显示用户头像、切换用户头像状态 | **✓** | **✓** | **X** | V7.47.0+
[CustomizedInput.show](https://open.feishu.cn/document/uYjL24iN/uADN1EjLwQTNx4CM0UTM/customizedinput/show) | 显示输入框 | **✓** | **✓** | **X** | V7.47.0+
[CustomizedInput.update](https://open.feishu.cn/document/uYjL24iN/uADN1EjLwQTNx4CM0UTM/customizedinput/update) | 更新输入框中显示的内容 | **✓** | **✓** | **X** | V7.47.0+
[CustomizedInput.hide](https://open.feishu.cn/document/uYjL24iN/uADN1EjLwQTNx4CM0UTM/customizedinput/hide) | 隐藏输入框 | **✓** | **✓** | **X** | V7.47.0+
[CustomizedInput.onPicSelect](https://open.feishu.cn/document/uYjL24iN/uADN1EjLwQTNx4CM0UTM/customizedinput/onpicselect) | 监听连接成功的事件回调 | **✓** | **✓** | **X** | **X**
[CustomizedInput.onModelSelect](https://open.feishu.cn/document/uYjL24iN/uADN1EjLwQTNx4CM0UTM/customizedinput/onmodelselect) | 选择pickerView之后触发的事件 | **✓** | **✓** | **X** | **X**
[CustomizedInput.onPublish](https://open.feishu.cn/document/uYjL24iN/uADN1EjLwQTNx4CM0UTM/customizedinput/onpublish) | 点击发送按钮触发的事件 | **✓** | **✓** | **X** | V7.47.0+
[CustomizedInput.onHide](https://open.feishu.cn/document/uYjL24iN/uADN1EjLwQTNx4CM0UTM/customizedinput/onhide) | 隐藏输入框之后触发的事件 | **✓** | **✓** | **X** | V7.47.0+
[CustomizedInput.onError](https://open.feishu.cn/document/uYjL24iN/uADN1EjLwQTNx4CM0UTM/customizedinput/onerror) | 监听错误回调 | 6.4.0 | 6.4.0 | **X** | V7.47.0+

#### Pad

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[togglePadFullScreen](https://open.feishu.cn/document/uYjL24iN/uUTOuUTOuUTO/pad/togglepadfullscreen) | 在当前Pad小程序窗口可以全屏缩放的前提下， 进行全屏缩放状态的切换 | 4.10.0 | 4.10.0 | **X** | **X**
[getPadDisplayScaleMode](https://open.feishu.cn/document/uYjL24iN/uUTOuUTOuUTO/pad/getpaddisplayscalemode) | 获取当前Pad的小程序窗口缩放状态，当前显示状态 能否进行全屏缩放的切换 | 4.10.0 | 4.10.0 | **X** | **X**

#### DarkMode

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[onThemeChange](https://open.feishu.cn/document/uYjL24iN/uUTOuUTOuUTO/darkmode/onthemechange) | 监听系统主题变化的事件 | 5.3.0 | 5.3.0 | 5.3.0 | **X**
[offThemeChange](https://open.feishu.cn/document/uYjL24iN/uUTOuUTOuUTO/darkmode/offthemechange) | 取消监听系统主题变化的事件 | 5.3.0 | 5.3.0 | 5.3.0 | **X**

## 设备
#### 系统信息

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getSystemInfo](https://open.feishu.cn/document/uYjL24iN/uQjNx4CN2EjL0YTM) | 获取系统信息 | **✓** | **✓** | **✓** | 7.35.0
[getSystemInfoSync](https://open.feishu.cn/document/uYjL24iN/uUjNx4SN2EjL1YTM) | 获取系统信息 | **✓** | **✓** | **✓** | 7.35.0

#### NFC

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getNFCAdapter](https://open.feishu.cn/document/uYjL24iN/ukzM4YjL5MDO24SOzgjN) | 获取 NFC 实例 | 3.38.0 | 5.25.0 | **X** | 7.35.0
[NFCAdapter.getNfcA](https://open.feishu.cn/document/uYjL24iN/ugzM4YjL4MDO24COzgjN) | 获取NfcA实例，实例支持NFC-A (ISO 14443-3A)标准的读写 | 3.38.0 | 5.25.0 | **X** | 7.35.0
[NFCAdapter.getNfcB](https://open.feishu.cn/document/uYjL24iN/uUzM4YjL1MDO24SNzgjN/getnfcb) | 获取NfcB实例，实例支持NFC-B (ISO 14443-3B)标准的读写 | 7.15.0 | 7.15.0 | **X** | 7.35.0
[NFCAdapter.getNfcV](https://open.feishu.cn/document/uYjL24iN/uUzM4YjL1MDO24SNzgjN/getNfcV) | 获取NfcV实例，实例支持NFC-V (ISO 15693)标准的读写 | 5.14.0 | **X** | **X** | 7.35.0
[NFCAdapter.getIsoDep](https://open.feishu.cn/document/uYjL24iN/uUzM4YjL1MDO24SNzgjN/getisodep) | 获取IsoDep实例，实例支持ISO-DEP (ISO 14443-4)标准的读写 | 7.15.0 | **X** | **X** | 7.35.0
[NFCAdapter.getMifareClassic](https://open.feishu.cn/document/uYjL24iN/uEDN4YjLxQDO24SM0gjN) | 获取MifareClassic实例，实例支持MIFARE Classic标签的读写 | 3.38.0 | **X** | **X** | 7.35.0
[NFCAdapter.startDiscovery](https://open.feishu.cn/document/uYjL24iN/uIDN4YjLyQDO24iM0gjN) | 开始扫描NFC标签 | 3.38.0 | 5.25.0 | **X** | 7.35.0
[NFCAdapter.stopDiscovery](https://open.feishu.cn/document/uYjL24iN/uMDN4YjLzQDO24yM0gjN) | 关闭NFC标签扫描 | 3.38.0 | 5.25.0 | **X** | 7.35.0
[NFCAdapter.onDiscovered](https://open.feishu.cn/document/uYjL24iN/uUDN4YjL1QDO24SN0gjN) | 监听 NFC Tag | 3.38.0 | 5.25.0 | **X** | 7.35.0
[NFCAdapter.offDiscovered](https://open.feishu.cn/document/uYjL24iN/uQDN4YjL0QDO24CN0gjN) | 取消监听 NFC Tag | 3.38.0 | **X** | **X** | 7.35.0
[NfcA.connect](https://open.feishu.cn/document/uYjL24iN/ucDN4YjL3QDO24yN0gjN) | 连接NfcA类型的标签 | 3.38.0 | 5.25.0 | **X** | 7.35.0
[NfcA.transceive](https://open.feishu.cn/document/uYjL24iN/uITN4YjLyUDO24iM1gjN) | 发送数据给NFCA类型的标签 | 3.38.0 | 5.25.0 | **X** | 7.35.0
[NfcA.close](https://open.feishu.cn/document/uYjL24iN/uYDN4YjL2QDO24iN0gjN) | 断开与NFCA标签之间的连接 | 3.38.0 | 5.25.0 | **X** | 7.35.0
[NfcA.getAtqa](https://open.feishu.cn/document/uYjL24iN/ugDN4YjL4QDO24CO0gjN) | 获取ATQA信息 | 3.38.0 | **X** | **X** | 7.35.0
[NfcA.getMaxTransceiveLength](https://open.feishu.cn/document/uYjL24iN/ukDN4YjL5QDO24SO0gjN) | 获取最大传输长度 | 3.38.0 | **X** | **X** | 7.35.0
[NfcA.getSak](https://open.feishu.cn/document/uYjL24iN/uATN4YjLwUDO24CM1gjN) | 获取SAK信息 | 3.38.0 | **X** | **X** | 7.35.0
[NfcA.setTimeout](https://open.feishu.cn/document/uYjL24iN/uETN4YjLxUDO24SM1gjN) | 设置超时时间 | 3.38.0 | **X** | **X** | 7.35.0
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
[MifareClassic.connect](https://open.feishu.cn/document/uYjL24iN/uQTN4YjL0UDO24CN1gjN) | 连接MifareClassic类型的标签 | 3.38.0 | **X** | **X** | 7.35.0
[MifareClassic.transceive](https://open.feishu.cn/document/uYjL24iN/ucTN4YjL3UDO24yN1gjN) | 发送数据给MifareClassic类型的标签 | 3.38.0 | **X** | **X** | 7.35.0
[MifareClassic.close](https://open.feishu.cn/document/uYjL24iN/uMTN4YjLzUDO24yM1gjN) | 断开与MifareClassic标签之间的连接 | 3.38.0 | **X** | **X** | 7.35.0
[MifareClassic.getMaxTransceiveLength](https://open.feishu.cn/document/uYjL24iN/uUTN4YjL1UDO24SN1gjN) | 获取最大传输长度 | 3.38.0 | **X** | **X** | 7.35.0
[MifareClassic.setTimeout](https://open.feishu.cn/document/uYjL24iN/uYTN4YjL2UDO24iN1gjN) | 设置超时时间 | 3.38.0 | **X** | **X** | 7.35.0

#### 蓝牙

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[openBluetoothAdapter](https://open.feishu.cn/document/uYjL24iN/ugzNxYjL4cTM24CO3EjN) | 初始化蓝牙模块 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[closeBluetoothAdapter](https://open.feishu.cn/document/uYjL24iN/uYDOxYjL2gTM24iN4EjN) | 关闭蓝牙模块。调用该方法将断开所有已建立的连接并释放系统资源。建议在使用蓝牙流程后，与 tt.openBluetoothAdapter 成对调用 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[getBluetoothAdapterState](https://open.feishu.cn/document/uYjL24iN/uUDOxYjL1gTM24SN4EjN) | 获取本机蓝牙适配器状态 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[startBluetoothDevicesDiscovery](https://open.feishu.cn/document/uYjL24iN/uUzNxYjL1cTM24SN3EjN) | 开始搜寻附近的蓝牙外围设备 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[stopBluetoothDevicesDiscovery](https://open.feishu.cn/document/uYjL24iN/uczNxYjL3cTM24yN3EjN) | 停止搜寻附近的蓝牙外围设备。若已经找到需要的蓝牙设备并不需要继续搜索时，建议调用该接口停止蓝牙搜索 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[getConnectedBluetoothDevices](https://open.feishu.cn/document/uYjL24iN/uMDOxYjLzgTM24yM4EjN) | 根据 uuid 获取处于已连接状态的设备 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[getBluetoothDevices](https://open.feishu.cn/document/uYjL24iN/uQDOxYjL0gTM24CN4EjN) | 获取在蓝牙模块生效期间所有已发现的蓝牙设备。包括已经和本机处于连接状态的设备 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[onBluetoothDeviceFound](https://open.feishu.cn/document/uYjL24iN/ukzNxYjL5cTM24SO3EjN) | 监听寻找到新设备的事件 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[offBluetoothDeviceFound](https://open.feishu.cn/document/uYjL24iN/uEDOxYjLxgTM24SM4EjN) | 取消监听寻找到新设备的事件 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[onBluetoothAdapterStateChange](https://open.feishu.cn/document/uYjL24iN/uADOxYjLwgTM24CM4EjN) | 监听蓝牙适配器状态变化事件 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[offBluetoothAdapterStateChange](https://open.feishu.cn/document/uYjL24iN/uIDOxYjLygTM24iM4EjN) | 取消监听蓝牙适配器状态变化事件 | 3.25.0 | 3.25.0 | **X** | 7.35.0

#### 低功耗蓝牙

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[setBLEMTU](https://open.feishu.cn/document/uYjL24iN/uMTMyYjLzEjM24yMxIjN) | 设置蓝牙最大传输单元 | 3.26.0 | **X** | **X** | 7.35.0
[readBLECharacteristicValue](https://open.feishu.cn/document/uYjL24iN/uYTOxYjL2kTM24iN5EjN) | 读取数据 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[writeBLECharacteristicValue](https://open.feishu.cn/document/uYjL24iN/ucTOxYjL3kTM24yN5EjN) | 写入数据 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[getBLEDeviceCharacteristics](https://open.feishu.cn/document/uYjL24iN/ukDOxYjL5gTM24SO4EjN) | 获取读写特征 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[connectBLEDevice](https://open.feishu.cn/document/uYjL24iN/ucDOxYjL3gTM24yN4EjN) | 链接外围设备 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[disconnectBLEDevice](https://open.feishu.cn/document/uYjL24iN/ugDOxYjL4gTM24CO4EjN) | 断开设备连接 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[getBLEDeviceServices](https://open.feishu.cn/document/uYjL24iN/uATOxYjLwkTM24CM5EjN) | 获取设备服务 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[onBLEConnectionStateChange](https://open.feishu.cn/document/uYjL24iN/uUTOxYjL1kTM24SN5EjN) | 蓝牙连接状态变化 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[offBLEConnectionStateChange](https://open.feishu.cn/document/uYjL24iN/uMTOxYjLzkTM24yM5EjN) | 取消监听蓝牙低功耗连接状态的改变事件 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[notifyBLECharacteristicValueChange](https://open.feishu.cn/document/uYjL24iN/uETOxYjLxkTM24SM5EjN) | 监听特征值数据变化 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[onBLECharacteristicValueChange](https://open.feishu.cn/document/uYjL24iN/uQTOxYjL0kTM24CN5EjN) | 监听特征值数据变化 | 3.25.0 | 3.25.0 | **X** | 7.35.0
[offBLECharacteristicValueChange](https://open.feishu.cn/document/uYjL24iN/uITOxYjLykTM24iM5EjN) | 取消监听蓝牙低功耗设备的特征值变化事件 | 3.25.0 | 3.25.0 | **X** | 7.35.0

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
[scanCode](https://open.feishu.cn/document/uYjL24iN/uYzNx4iN3EjL2cTM) | 扫描二维码并返回扫描结果 | **✓** | **✓** | **X** | 7.35.0

#### Wi-Fi

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getConnectedWifi](https://open.feishu.cn/document/uYjL24iN/ugjNx4CO2EjL4YTM) | 获取设备当前所连的 Wifi | **✓** | **✓** | **X** | 7.35.0
[getWifiStatus](https://open.feishu.cn/document/uYjL24iN/uYTN4QjL2UDO04iN1gDN) | 请求获取 Wi-Fi 开关状态 | **✓** | **✓** | **X** | 7.35.0
[getWifiList](https://open.feishu.cn/document/uYjL24iN/uUDO4UjL1gDO14SN4gTN) | 请求获取Wifi 列表 | **✓** | **X** | **X** | 7.35.0
[onGetWifiList](https://open.feishu.cn/document/uYjL24iN/uYDO4UjL2gDO14iN4gTN) | 监听获取到 Wi-Fi 列表数据事件 | **✓** | **X** | **X** | 7.35.0
[offGetWifiList](https://open.feishu.cn/document/uYjL24iN/ucDO4UjL3gDO14yN4gTN) | 取消监听获取到 Wi-Fi 列表数据事件 | **✓** | **X** | **X** | 7.35.0

#### 剪贴板

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[setClipboardData](https://open.feishu.cn/document/uYjL24iN/ugzNx4CO3EjL4cTM) | 设置系统剪贴板内容 | **✓** | **✓** | **✓** | 7.35.0
[getClipboardData](https://open.feishu.cn/document/uYjL24iN/uczNx4yN3EjL3cTM) | 获取系统粘贴板数据 | **✓** | **✓** | **✓** | 7.35.0

#### 网络状态

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getNetworkType](https://open.feishu.cn/document/uYjL24iN/uYjNx4iN2EjL2YTM) | 获取设备当前所处的网络类型 | **✓** | **✓** | **✓** | 7.35.0
[onNetworkStatusChange](https://open.feishu.cn/document/uYjL24iN/ucjNx4yN2EjL3YTM) | 监听网络状态变化 | **✓** | **✓** | **X** | 7.35.0
[getNetworkQualityType](https://open.feishu.cn/document/uYjL24iN/uUTNx4SN1EjL1UTM/getnetworkqualitytype) | 网络评级接口，获取当前设备所处的网络状态 | 4.9.0 | 4.9.0 | 5.1.0 | 7.35.0
[onNetworkQualityChange](https://open.feishu.cn/document/uYjL24iN/uUTNx4SN1EjL1UTM/onnetworkqualitychange) | 监听网络质量变化 | 4.9.0 | 4.9.0 | 5.1.0 | 7.35.0

#### 屏幕亮度

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[setScreenBrightness](https://open.feishu.cn/document/uYjL24iN/uIjNx4iM2EjLyYTM/set-screen-brightness) | 设置屏幕亮度 | 3.42.0 | 3.42.0 | **X** | 7.35.0
[getScreenBrightness](https://open.feishu.cn/document/uYjL24iN/uIjNx4iM2EjLyYTM/get-screen-brightness) | 获取屏幕亮度 | 3.42.0 | 3.42.0 | **X** | 7.35.0
[setKeepScreenOn](https://open.feishu.cn/document/uYjL24iN/ukzNx4SO3EjL5cTM) | 设置是否保持常亮状态 | **✓** | **✓** | **X** | 7.35.0

#### 截屏监听

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[grantScreenShot](https://open.feishu.cn/document/uYjL24iN/uIjNwEjLyYDMx4iM2ATM/grantscreenshot) | 用于开启/禁用截屏，仅在Android和Harmony设备上生效。 | 7.7.0 | **X** | **X** | 7.35.0
[onUserCaptureScreen](https://open.feishu.cn/document/uYjL24iN/uMjNwEjLzYDMx4yM2ATM) | 监听用户主动截屏事件。用户使用系统截屏按键截屏时触发 | 2.4.0 | 2.4.0 | **X** | 7.35.0
[offUserCaptureScreen](https://open.feishu.cn/document/uYjL24iN/uQjNwEjL0YDMx4CN2ATM) | 取消监听用户主动截屏事件 | 2.4.0 | 2.4.0 | **X** | 7.35.0

#### 加速度计

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[startAccelerometer](https://open.feishu.cn/document/uYjL24iN/ukjNx4SO2EjL5YTM) | 通知客户端开始监听加速度计数据。具体的数据返回通过注册onAccelerometerChange接口回调方法获取 | **✓** | **✓** | **X** | 7.35.0
[stopAccelerometer](https://open.feishu.cn/document/uYjL24iN/uAzNx4CM3EjLwcTM) | 停止监听加速度计数据 | **✓** | **✓** | **X** | 7.35.0
[onAccelerometerChange](https://open.feishu.cn/document/uYjL24iN/uEzNx4SM3EjLxcTM) | 监听加速度计数据 | **✓** | **✓** | **X** | 7.35.0

#### 拨打电话

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[makePhoneCall](https://open.feishu.cn/document/uYjL24iN/uUzNx4SN3EjL1cTM) | 拨打电话 | **✓** | **✓** | **X** | 7.35.0

#### 震动

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[vibrateShort](https://open.feishu.cn/document/uYjL24iN/uADOx4CM4EjLwgTM) | 使手机发生较短时间的振动 | **✓** | **✓** | **X** | 7.35.0
[vibrateLong](https://open.feishu.cn/document/uYjL24iN/uEDOx4SM4EjLxgTM) | 使手机发生较长时间的振动 | **✓** | **✓** | **X** | 7.35.0

#### 罗盘

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[startCompass](https://open.feishu.cn/document/uYjL24iN/uIzNx4iM3EjLycTM) | 开始监听罗盘数据。 | **✓** | **✓** | **X** | 7.35.0
[stopCompass](https://open.feishu.cn/document/uYjL24iN/uMzNx4yM3EjLzcTM) | 停止监听罗盘数据 | **✓** | **✓** | **X** | 7.35.0
[onCompassChange](https://open.feishu.cn/document/uYjL24iN/uQzNx4CN3EjL0cTM) | 监听罗盘数据变化事件，频率：5 次/秒，接口调用后会自动开始监听，可使用 tt.stopCompass 停止监听 | **✓** | **✓** | **X** | 7.35.0

#### 性能

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[onMemoryWarning](https://open.feishu.cn/document/uYjL24iN/uQTOuQTOuQTO/performance/onmemorywarning) | 监听内存不足的告警事件。当手机内存占用过高时，触发回调函数。 | 5.4.0 | 5.4.0 | **X** | **X**

## 文件

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[saveFile](https://open.feishu.cn/document/uYjL24iN/ugDOz4CO4MjL4gzM) | 保存临时文件到本地永久目录 | **✓** | **✓** | **✓** | 7.35.0
[saveFileAs](https://open.feishu.cn/document/uYjL24iN/uQjN3UjL0YzN14CN2cTN) | 保存文件到本地指定目录 | **X** | **X** | 3.9.0 | **X**
[filePicker](https://open.feishu.cn/document/uYjL24iN/uETM04SMxQjLxEDN) | 打开附件选择列表 | **✓** | **✓** | **✓** | 7.35.0
[docsPicker](https://open.feishu.cn/document/uYjL24iN/ukTN3UjL5UzN14SO1cTN) | 打开云文档选择列表 | 3.12.0 | 3.12.0 | 3.13.0 | 7.49.0
[openDocument](https://open.feishu.cn/document/uYjL24iN/ukTN24SO1YjL5UjN) | 在新页面打开文档 | 2.6.0 | 2.6.0 | 2.6.0 | 7.39.0
[getFileSystemManager](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/tt_get_file_system_manager) | 获取全局唯一的文件管理器 | 4.11.0 | 4.11.0 | 4.11.0 | 7.35.0
[FileSystemManager.appendFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/appendfile) | 在文件结尾追加内容 | 5.23.0 | 5.23.0 | 5.24.0 | 7.35.0
[FileSystemManager.appendFileSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/appendfilesync) | 在文件结尾追加内容 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[FileSystemManager.access](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_access) | 判断文件/目录是否存在 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.accessSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/access_sync) | 判断文件/目录是否存在 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[FileSystemManager.copyFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_copy_file) | 复制本地文件/目录 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.copyFileSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_copy_file_sync) | 复制本地文件/目录 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[FileSystemManager.getFileInfo](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_get_file_info) | 获取本地文件信息 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.getSavedFileList](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_get_saved_file_list) | 获取用户目录内文件列表 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.mkdir](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_mkdir) | 创建本地目录 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.mkdirSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_mkdir_sync) | 创建本地目录 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[FileSystemManager.readdir](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_read_dir) | 读取本地目录内文件列表 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.readdirSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_read_dir_sync) | 读取本地目录内文件列表 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[FileSystemManager.readFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_read_file) | 读取本地文件内容 | 4.11.0 | 4.11.0 | 5.23.0 | 7.35.0
[FileSystemManager.readFileSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_read_file_sync) | 读取本地文件内容 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[FileSystemManager.removeSavedFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_remove_saved_file) | 删除保存的本地文件 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.rename](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_rename) | 重命名本地文件/目录 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.renameSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_rename_sync) | 重命名本地文件/目录 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[FileSystemManager.rmdir](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_rmdir) | 删除本地目录 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.rmdirSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_rmdir_sync) | 删除本地目录 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[FileSystemManager.saveFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_save_file) | 保存临时文件到本地用户目录 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.saveFileSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_save_file_sync) | 保存临时文件到本地用户目录 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[FileSystemManager.stat](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_stat) | 获取本地文件 Stats 对象 | 4.11.0 | 4.11.0 | 5.23.0 | 7.39.0
[FileSystemManager.statSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_stat_sync) | 获取本地文件 Stats 对象 | 5.23.0 | 5.23.0 | **X** | 7.39.0
[FileSystemManager.unlink](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_unlink) | 删除本地文件 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.unlinkSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_unlink_sync) | 删除本地文件 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[FileSystemManager.writeFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_write_file) | 写入本地文件 | 5.23.0 | 5.23.0 | 5.23.0 | 7.35.0
[FileSystemManager.writeFileSync](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/file_system_manager/file_system_manager_write_file_sync) | 写入本地文件 | 5.23.0 | 5.23.0 | **X** | 7.35.0
[Stats.isDirectory](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/stat/stats_is_directory) | 判断当前文件是否一个目录 | 4.11.0 | 4.11.0 | **X** | 7.39.0
[Stats.isFile](https://open.feishu.cn/document/uYjL24iN/uETOuETOuETO/stat/stats_is_file) | 判断当前文件是否一个普通文件 | 4.11.0 | 4.11.0 | **X** | 7.39.0

## 媒体
#### 图片

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[saveImageToPhotosAlbum](https://open.feishu.cn/document/uYjL24iN/uUTMx4SNxEjL1ETM) | 保存图片到系统相册 | **✓** | **✓** | **✓** | 7.35.0
[chooseImage](https://open.feishu.cn/document/uYjL24iN/uMTMx4yMxEjLzETM) | 从系统相册中选择图片，或使用相机拍摄图片 | **✓** | **✓** | **✓** | 7.35.0
[previewImage](https://open.feishu.cn/document/uYjL24iN/uMDOx4yM4EjLzgTM) | 预览一组图片 | **✓** | **✓** | **✓** | 7.35.0
[compressImage](https://open.feishu.cn/document/uYjL24iN/uMjN24yM2YjLzYjN) | 压缩图片接口，可选压缩质量 | 2.0.0 | 2.0.0 | 2.0.0 | 7.35.0
[getImageInfo](https://open.feishu.cn/document/uYjL24iN/ugjNwEjL4YDMx4CO2ATM) | 获取图片信息 | 2.4.0 | 2.4.0 | 2.4.0 | 7.35.0

#### 视频

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getVideoInfo](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/getvideoinfo) | 获取视频详细信息 | 6.5.0 | 6.5.0 | 6.5.0 | 7.35.0
[saveVideoToPhotosAlbum](https://open.feishu.cn/document/uYjL24iN/ucDOx4yN4EjL3gTM) | 保存视频到系统相册 | **✓** | **✓** | **✓** | 7.35.0
[chooseVideo](https://open.feishu.cn/document/uYjL24iN/uEjMx4SMyEjLxITM) | 从系统相册中选择视频，或使用相机拍摄视频 | **✓** | **✓** | **✓** | **X**
[chooseMedia](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/choosemedia) | 拍摄或从系统相册中选择图片或视频 | 4.7.0 | 4.7.0 | 4.7.0 | 7.35.0
[createVideoContext](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/createvideocontext) | 创建 VideoContext 实例，通过 id 跟一个 video 组件绑定，操作对应的 video 组件 | 4.3.0 | 4.3.0 | 6.1.0 | 7.35.0
[VideoContext.play](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/videocontext/play) | 播放视频 | 4.3.0 | 4.3.0 | 6.1.0 | 7.35.0
[VideoContext.pause](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/videocontext/pause) | 暂停视频 | 4.3.0 | 4.3.0 | 6.1.0 | 7.35.0
[VideoContext.stop](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/videocontext/stop) | 停止视频 | 4.3.0 | 4.3.0 | 6.1.0 | 7.35.0
[VideoContext.seek](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/videocontext/seek) | 跳转到指定位置 | 4.3.0 | 4.3.0 | 6.1.0 | 7.35.0
[VideoContext.requestFullScreen](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/videocontext/requestfullscreen) | 进入全屏 | 4.3.0 | 4.3.0 | 6.1.0 | 7.35.0
[VideoContext.exitFullScreen](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/videocontext/exitfullscreen) | 退出全屏 | 4.3.0 | 4.3.0 | 6.1.0 | 7.35.0
[VideoContext.playbackRate](https://open.feishu.cn/document/uYjL24iN/uITMx4iMxEjLyETM/videocontext/playbackrate) | 设置视频倍速播放 | 6.1.0 | 6.1.0 | 6.1.0 | 7.35.0

#### 音频

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[createInnerAudioContext](https://open.feishu.cn/document/uYjL24iN/uUDOx4SN4EjL1gTM) | 创建innerAudioContext实例，通过它能够操作音频播放 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.play](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/play) | 播放音频 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.pause](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/pause) | 暂停。暂停后的音频再播放会从暂停处开始播放 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.stop](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/stop) | 停止。停止后的音频再播放会从头开始播放 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.seek](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/seek) | 跳转到指定位置 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.destroy](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/destroy) | 销毁当前实例 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.onCanplay](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/oncanplay) | 音频进入可以播放状态时触发回调函数 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.offCanplay](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/offcanplay) | 取消监听 Canplay 事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.onPlay](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/onPlay) | 监听音频播放事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.offPlay](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/offplay) | 取消监听音频播放事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.onPause](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/onPause) | 监听音频暂停事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.offPause](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/offpause) | 取消监听音频暂停事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.onStop](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/onstop) | 监听音频停止事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.offStop](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/offstop) | 取消监听音频停止事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.onEnded](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/onended) | 监听音频自然播放至结束的事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.offEnded](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/offEnded) | 取消监听音频自然播放至结束的事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.onTimeUpdate](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/onTimeUpdate) | 监听音频播放进度更新事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.offTimeUpdate](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/offtimeupdate) | 取消监听音频播放进度更新事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.onError](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/onerror) | 监听音频播放错误事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.offError](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/offerror) | 取消监听音频播放错误事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.onWaiting](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/onwaiting) | 监听音频加载中事件。当音频因为数据不足，需要停下来加载时会触发 | **✓** | **✓** | **X** | **X**
[InnerAudioContext.offWaiting](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/offwaiting) | 取消监听音频加载中事件 | **✓** | **✓** | **X** | **X**
[InnerAudioContext.onSeeking](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/onseeking) | 监听音频进行跳转操作的事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.offSeeking](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/offseeking) | 取消监听音频进行跳转操作的事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.onSeeked](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/onSeeked) | 监听音频完成跳转操作的事件 | **✓** | **✓** | **X** | 7.43.0
[InnerAudioContext.offSeeked](https://open.feishu.cn/document/uYjL24iN/uETMx4SMxEjLxETM/inneraudiocontext/offseeked) | 取消监听音频完成跳转操作的事件 | **✓** | **✓** | **X** | 7.43.0

#### 背景音频

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getBackgroundAudioManager](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/getbackgroundaudiomanager) | 创建backgroundAudioManager实例，通过它能够操作背景音频播放 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.play](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/play) | 播放背景音频 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.pause](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/pause) | 背景音频暂停 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.stop](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/stop) | 背景音频停止播放 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.seek](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/seek) | 背景音频跳转到指定位置播放 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.onCanplay](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/oncanplay) | 在背景音频进入可播放状态后触发回调函数 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.onPlay](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/onplay) | 监听背景音频播放事件 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.onPause](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/onpause) | 监听背景音频暂停事件 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.onStop](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/onstop) | 监听背景音频停止事件 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.onEnded](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/onended) | 监听背景音频自然播放至结束的事件 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.onTimeUpdate](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/ontimeupdate) | 监听背景音频播放进度更新事件 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.onError](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/onerror) | 监听背景音频播放错误事件 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.onWaiting](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/onwaiting) | 监听背景音频加载中事件。当音频因为数据不足，需要停下来加载时会触发 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.onSeeking](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/onseeking) | 监听音频进行跳转操作的事件 | 5.20.0 | 5.20.0 | **X** | **X**
[BackgroundAudioManager.onSeeked](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/backgroundaudio/backgroundaudiomanager/onseeked) | 监听背景音频完成跳转操作的事件 | 5.20.0 | 5.20.0 | **X** | **X**

#### 录音

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getRecorderManager](https://open.feishu.cn/document/uYjL24iN/uQDOx4CN4EjL0gTM) | 获取全局唯一的recorderManager。通过recorderManager进行录音操作和管理 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.start](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/start) | 开始录音 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.pause](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/pause) | 暂停录音 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.resume](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/resume) | 继续录音 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.stop](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/stop) | 停止录音 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.getState](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/getstate) | 获取当前录音状态 | 7.13.0 | 7.13.0 | **X** | 7.43.0
[RecorderManager.onStart](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onstart) | 监听录音开始事件 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.onPause](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onpause) | 监听录音暂停事件 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.onResume](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onresume) | 监听录音继续事件 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.onStop](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onstop) | 监听录音结束事件 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.onFrameRecorded](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onframerecorded) | 监听已录制完指定帧大小的文件事件。如果设置了 frameSize，则会回调此事件 | **✓** | **✓** | **X** | 7.43.0
[RecorderManager.onError](https://open.feishu.cn/document/uYjL24iN/uATMx4CMxEjLwETM/recordermanager/onerror) | 监听录音错误事件 | **✓** | **✓** | **X** | 7.43.0

#### 相机

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[createCameraContext](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/camera/createcameracontext) | 创建 camera 上下文 CameraContext 对象 | 5.21.0 | 5.21.0 | **X** | **X**
[CameraContext.setZoom](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/camera/cameracontext/setzoom) | 设置缩放级别 | 5.21.0 | 5.21.0 | **X** | **X**
[CameraContext.takePhoto](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/camera/cameracontext/takephoto) | 拍摄照片 | 5.21.0 | 5.21.0 | **X** | **X**
[CameraContext.startRecord](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/camera/cameracontext/startrecord) | 开始录像 | 5.21.0 | 5.21.0 | **X** | **X**
[CameraContext.stopRecord](https://open.feishu.cn/document/uYjL24iN/ukDOukDOukDO/camera/cameracontext/stoprecord) | 结束录像 | 5.21.0 | 5.21.0 | **X** | **X**

## 导航

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[onNavigateBack](https://open.feishu.cn/document/uYjL24iN/uYTOuYTOuYTO/onnavigateback) | 监听点击页面返回的事件 | 5.19.0 | 5.19.0 | 5.19.0 | 7.35.0
[navigateTo](https://open.feishu.cn/document/uYjL24iN/uYTOz4iN5MjL2kzM) | 跳转到指定页面。跳转后原页面保留 | **✓** | **✓** | **✓** | 7.35.0
[navigateBack](https://open.feishu.cn/document/uYjL24iN/uADM04CMwQjLwADN) | 返回上一级页面（或上N级页面）。可通过 getCurrentPages 获取当前的页面栈，决定需要返回几层 | **✓** | **✓** | **✓** | 7.35.0
[redirectTo](https://open.feishu.cn/document/uYjL24iN/ucTOz4yN5MjL3kzM) | 关闭当前页面，跳转到指定页面 | **✓** | **✓** | **✓** | 7.35.0
[switchTab](https://open.feishu.cn/document/uYjL24iN/ukTOz4SO5MjL5kzM) | 跳转到指定 TabBar 页面，并关闭其他所有非 TabBar 页面 | **✓** | **✓** | **✓** | 7.35.0
[reLaunch](https://open.feishu.cn/document/uYjL24iN/uEDM04SMwQjLxADN) | 关闭所有当前页面，打开指定页面 | **✓** | **✓** | **✓** | 7.35.0
[exitMiniProgram](https://open.feishu.cn/document/uYjL24iN/uATN4IjLwUDOy4CM1gjM) | 退出当前小程序 | **✓** | **✓** | **✓** | 7.35.0
[openSchema](https://open.feishu.cn/document/uYjL24iN/ukzN4IjL5cDOy4SO3gjM) | 跳转到小程序以外的应用（如云文档、网页等） | 3.1.0 | 3.1.0 | 3.1.0 | 7.35.0

## 网络

#### 发起请求

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[request](https://open.feishu.cn/document/uYjL24iN/uIDMx4iMwEjLyATM) | 发起一个 HTTP 请求 | **✓** | **✓** | **✓** | 7.35.0
[RequestTask.abort](https://open.feishu.cn/document/uYjL24iN/ugDNugDNugDN/requesttask/abort) | 中断请求任务 | **✓** | **✓** | **✓** | 7.35.0

#### 下载

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[downloadFile](https://open.feishu.cn/document/uYjL24iN/ucDMx4yNwEjL3ATM) | 下载网络文件到本地临时目录 | **✓** | **✓** | **✓** | 7.35.0
[DownloadTask.onProgressUpdate](https://open.feishu.cn/document/uYjL24iN/ugDNugDNugDN/downloadfile/onprogressupdate) | downloadFile的调用结果在通过回调传递的同时会返回一个downloadTask对象，通过onProgressUpdate方法监听下载进度 | **✓** | **✓** | **✓** | 7.35.0
[DownloadTask.abort](https://open.feishu.cn/document/uYjL24iN/ugDNugDNugDN/downloadfile/abort) | 调用downloadFile时，会返回一个downloadTask对象，可以通过该对象的abort方法中断请求任务 | **✓** | **✓** | **✓** | 7.35.0

#### 上传

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[uploadFile](https://open.feishu.cn/document/uYjL24iN/uYDMx4iNwEjL2ATM) | 将本地文件上传到网络 | **✓** | **✓** | **✓** | 7.35.0
[UploadTask.onProgressUpdate](https://open.feishu.cn/document/uYjL24iN/ugDNugDNugDN/uploadtask/onprogressupdate) | 调用uploadFile时，会返回一个uploadTask对象，通过onProgressUpdate方法监听下载进度 | **✓** | **✓** | **✓** | 7.35.0
[UploadTask.abort](https://open.feishu.cn/document/uYjL24iN/ugDNugDNugDN/uploadtask/abort) | 调用uploadFile时，会返回一个uploadTask对象，可以通过该对象的abort方法中断请求任务 | **✓** | **✓** | **✓** | 7.35.0

#### WebSocket

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[connectSocket](https://open.feishu.cn/document/uYjL24iN/ugDMx4COwEjL4ATM) | 创建一个 WebSocket 连接实例，并通过返回的 socketTask 操作该连接 | **✓** | **✓** | **✓** | 7.39.0
[SocketTask.send](https://open.feishu.cn/document/uYjL24iN/ugDOugDOugDO/sockettask/send) | 通过 WebSocket 连接发送数据 | **✓** | **✓** | **✓** | 7.39.0
[SocketTask.close](https://open.feishu.cn/document/uYjL24iN/ugDOugDOugDO/sockettask/close) | 关闭 WebSocket 连接 | **✓** | **✓** | **✓** | 7.39.0
[SocketTask.onOpen](https://open.feishu.cn/document/uYjL24iN/ugDOugDOugDO/sockettask/onopen) | 监听 WebSocket 连接打开事件 | **✓** | **✓** | **✓** | 7.39.0
[SocketTask.onClose](https://open.feishu.cn/document/uYjL24iN/ugDOugDOugDO/sockettask/onclose) | 监听 WebSocket 连接关闭事件 | **✓** | **✓** | **✓** | 7.39.0
[SocketTask.onMessage](https://open.feishu.cn/document/uYjL24iN/ugDOugDOugDO/sockettask/onmessage) | 监听 WebSocket 接受到服务器的消息事件 | **✓** | **✓** | **✓** | 7.39.0

## 地理位置

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getLocationStatus](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/getlocationstatus) | 获取 GPS 开关状态 | 6.10.0 | 6.10.0 | **X** | 7.35.0
[reverseGeocode](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/reversegeocode) | 逆地理编码，根据传入的经纬度，获取相应的位置信息 | 6.6.0 | 6.6.0 | **X** | 7.35.0
[offLocationChangeError](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/offlocationchangeerror) | 取消监听持续定位接口返回失败时触发 | 5.21.0 | 5.21.0 | **X** | 7.35.0
[onLocationChangeError](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/onlocationchangeerror) | 监听持续定位接口返回失败时触发 | 5.21.0 | 5.21.0 | **X** | 7.35.0
[openLocation](https://open.feishu.cn/document/uYjL24iN/uQTOz4CN5MjL0kzM) | 使用客户端内置地图查看位置 | **✓** | **✓** | **X** | 7.35.0
[chooseLocation](https://open.feishu.cn/document/uYjL24iN/uUDN1EjL1QTNx4SN0UTM) | 打开地图选择位置 | **✓** | **✓** | **X** | 7.35.0
[getLocation](https://open.feishu.cn/document/uYjL24iN/uUTOz4SN5MjL1kzM) | 获取设备当前的地理位置 | **✓** | **✓** | **X** | 7.35.0
[startLocationUpdate](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/startlocationupdate) | 开启小程序在前台时接收位置更新消息 | 4.3.0 | 4.3.0 | **X** | 7.35.0
[stopLocationUpdate](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/stoplocationupdate) | 关闭监听实时位置变化，停止位置消息接收 | 4.3.0 | 4.3.0 | **X** | 7.35.0
[onLocationChange](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/onlocationchange) | 监听实时地理位置变化事件，需结合startLocationUpdate使用 | 4.3.0 | 4.3.0 | **X** | 7.35.0
[offLocationChange](https://open.feishu.cn/document/uYjL24iN/uMTOuMTOuMTO/offlocationchange) | 取消监听实时地理位置变化事件 | 4.9.0 | 4.9.0 | **X** | 7.35.0

## 数据缓存

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[setStorage](https://open.feishu.cn/document/uYjL24iN/uETOx4SM5EjLxkTM) | 以「键值对」的形式设置本地缓存数据 | **✓** | **✓** | **✓** | 7.35.0
[setStorageSync](https://open.feishu.cn/document/uYjL24iN/uITOx4iM5EjLykTM) | 以「键值对」的形式设置本地缓存数据 | **✓** | **✓** | **✓** | 7.35.0
[getStorage](https://open.feishu.cn/document/uYjL24iN/ukDOx4SO4EjL5gTM) | 获取本地缓存数据 | **✓** | **✓** | **✓** | 7.35.0
[getStorageSync](https://open.feishu.cn/document/uYjL24iN/uATOx4CM5EjLwkTM) | 获取本地缓存数据 | **✓** | **✓** | **✓** | 7.35.0
[removeStorage](https://open.feishu.cn/document/uYjL24iN/uMTOx4yM5EjLzkTM) | 删除本地缓存数据 | **✓** | **✓** | **✓** | 7.35.0
[removeStorageSync](https://open.feishu.cn/document/uYjL24iN/uQTOx4CN5EjL0kTM) | 删除本地缓存数据 | **✓** | **✓** | **✓** | 7.35.0
[clearStorage](https://open.feishu.cn/document/uYjL24iN/uUTOx4SN5EjL1kTM) | 清理全部本地缓存数据 | **✓** | **✓** | **✓** | 7.35.0
[clearStorageSync](https://open.feishu.cn/document/uYjL24iN/uYTOx4iN5EjL2kTM) | 清理全部本地缓存数据 | **✓** | **✓** | **✓** | 7.35.0
[getStorageInfo](https://open.feishu.cn/document/uYjL24iN/ucTOx4yN5EjL3kTM) | 获取本地缓存数据的相关信息 | **✓** | **✓** | **✓** | 7.35.0
[getStorageInfoSync](https://open.feishu.cn/document/uYjL24iN/ugTOx4CO5EjL4kTM) | 获取本地缓存数据的相关信息 | **✓** | **✓** | **✓** | 7.35.0

## TTML

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[createSelectorQuery](https://open.feishu.cn/document/uYjL24iN/uYjN24iN2YjL2YjN) | 获取一个 SelectorQuery 对象实例 | 2.2.0 | 2.2.0 | 2.2.0 | 7.35.0
[SelectorQuery.in](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/selectorquery/in) | 将选择器的选取范围更改为自定义组件 component 内（初始时，选择器仅选取页面范围的节点，不会选取任何自定义组件中的节点） | 2.2.0 | 2.2.0 | 2.2.0 | 7.35.0
[SelectorQuery.select](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/selectorquery/select) | 在当前页面下选择第一个匹配选择器 selector 的节点，返回一个 NodesRef 对象实例，可以用于获取节点信息。  selector 类似于 CSS 的选择器，其中移动端只支持 ID 选择器 | 2.2.0 | 2.2.0 | 2.2.0 | 7.35.0
[SelectorQuery.selectAll](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/selectorquery/selectall) | 在当前页面下选择匹配选择器 selector 的所有节点，返回一个 NodesRef 对象实例，可以用于获取节点信息。  selector 类似于 CSS 的选择器，同 select | 2.2.0 | 2.2.0 | 2.2.0 | 7.35.0
[SelectorQuery.selectViewport](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/selectorquery/selectviewport) | 选择显示区域。可用于获取显示区域的尺寸、滚动位置等信息 | 2.2.0 | 2.2.0 | **X** | 7.35.0
[SelectorQuery.exec](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/selectorquery/exec) | 执行所有的请求。请求结果按请求次序构成数组，在callback的第一个参数中返回 | 2.2.0 | 2.2.0 | 2.2.0 | 7.35.0
[NodesRef.boundingClientRect](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/nodesref/boundingclientrect) | 添加节点的布局位置的查询请求 | 2.2.0 | 2.2.0 | 2.2.0 | 7.35.0
[NodesRef.scrollOffset](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/nodesref/scrolloffset) | 添加节点的滚动位置查询请求 | 2.2.0 | 2.2.0 | 2.2.0 | 7.35.0
[NodesRef.fields](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/nodesref/fields) | 获取节点的相关信息 | 2.2.0 | 2.2.0 | 2.2.0 | 7.35.0
[createIntersectionObserver](https://open.feishu.cn/document/uYjL24iN/ucTNwEjL3UDMx4yN1ATM) | 创建并返回一个 IntersectionObserver 对象实例 | 2.3.0 | 2.3.0 | 2.3.0 | 7.35.0
[IntersectionObserver.observe](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/intersectionobserver/observe) | 指定目标节点并开始监听相交状态变化情况 | 2.3.0 | 2.3.0 | 2.3.0 | 7.35.0
[IntersectionObserver.relativeTo](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/intersectionobserver/relativeto) | 使用选择器指定一个节点，作为参照区域之一 | 2.3.0 | 2.3.0 | 2.3.0 | 7.35.0
[IntersectionObserver.relativeToViewport](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/intersectionobserver/relativetoviewport) | 指定页面显示区域作为参照区域之一 | 2.3.0 | 2.3.0 | 2.3.0 | 7.35.0
[IntersectionObserver.disconnect](https://open.feishu.cn/document/uYjL24iN/uUjN24SN2YjL1YjN/intersectionobserver/disconnect) | 停止监听，回调函数将不再触发 | 2.3.0 | 2.3.0 | 2.3.0 | 7.35.0

## 更新

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getUpdateManager](https://open.feishu.cn/document/uYjL24iN/uEzM04SMzQjLxMDN) | 获取全局唯一的版本更新管理器，返回 updateManager 对象，用于管理小程序更新 | **✓** | **✓** | **✓** | **X**
[UpdateManager.applyUpdate](https://open.feishu.cn/document/uYjL24iN/uAzM04CMzQjLwMDN/getUpdateManager/applyupdate) | 强制小程序重启并使用新版本。在小程序新版本下载完成后（即收到 onUpdateReady 回调）调用 | **✓** | **✓** | **✓** | **X**
[UpdateManager.onCheckForUpdate](https://open.feishu.cn/document/uYjL24iN/uAzM04CMzQjLwMDN/getUpdateManager/oncheckforupdate) | 监听向后台请求检查更新结果事件。客户端在小程序冷启动时自动检查更新，不需由开发者主动触发。线上环境在有更新内容时会触发callback回调 | **✓** | **✓** | **✓** | **X**
[UpdateManager.onUpdateFailed](https://open.feishu.cn/document/uYjL24iN/uAzM04CMzQjLwMDN/getUpdateManager/onupdatefailed) | 监听小程序更新失败事件 | **✓** | **✓** | **✓** | **X**
[UpdateManager.onUpdateReady](https://open.feishu.cn/document/uYjL24iN/uAzM04CMzQjLwMDN/getUpdateManager/onupdateready) | 监听小程序有版本更新事件 | **✓** | **✓** | **✓** | **X**
[UpdateManager.triggerCheckUpdate](https://open.feishu.cn/document/uYjL24iN/uAzM04CMzQjLwMDN/getUpdateManager/triggercheckupdate) | 主动触发更新小程序 | **✓** | **✓** | **✓** | **X**

## 性能

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getPerformanceRecorder](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/performance/getperformancerecorder) | 获取全局唯一的性能录制器对象 PerformanceRecorder，用于在 性能分析模式 下进行自定义埋点 | 6.1.0 | 6.1.0 | 6.1.0 | **X**
[PerformanceRecorder.mark](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/performance/PerformanceRecorder/mark) | 记录一个时间点，仅在性能分析模式下生效 | 6.1.0 | 6.1.0 | 6.1.0 | **X**
[PerformanceRecorder.startParallel](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/performance/PerformanceRecorder/startparallel) | 开始一段异步操作的记录，仅在性能分析模式下生效 | 6.1.0 | 6.1.0 | 6.1.0 | **X**
[PerformanceRecorder.endParallel](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/performance/PerformanceRecorder/endparallel) | 完成一段异步操作的记录，仅在性能分析模式下生效 | 6.1.0 | 6.1.0 | 6.1.0 | **X**
[PerformanceRecorder.startStack](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/performance/PerformanceRecorder/startstack) | 开始一段同步操作的记录，仅在性能分析模式下生效 | 6.1.0 | 6.1.0 | 6.1.0 | **X**
[PerformanceRecorder.endStack](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/performance/PerformanceRecorder/endstack) | 完成一段同步操作的记录，仅在性能分析模式下生效 | 6.1.0 | 6.1.0 | 6.1.0 | **X**

## 运动与健康

API名称 | 说明 | Android | iOS | PC | Harmony
--- | --- | --- | --- | --- | ---
[getStepCount](https://open.feishu.cn/document/uYjL24iN/uAjMuAjMuAjM/sport-and-health/getstepcount) | 获取设备中的步数信息 | 5.1.0 | 5.1.0 | **X** | 7.35.0

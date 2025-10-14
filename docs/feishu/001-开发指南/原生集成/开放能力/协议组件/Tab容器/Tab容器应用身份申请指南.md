---
title: "Tab容器应用身份申请指南"
source_url: https://open.feishu.cn/document/native-integration/open-capability/protocol-components/tab-container/tab-container-applic
last_remote_update: 2024-05-10
last_remote_update_timestamp: 1715328683000
---
最后更新于 2024-05-10

# Tab容器应用身份申请指南
使用Tab容器协议组件时，需要在飞书开放平台，创建一个原生集成应用后，才可以正常使用此协议组件，具体操作流程如下：

## 一、创建应用
### 1. 前往[飞书开放平台](https://open.feishu.cn/?lang=zh-CN)创建应用
![](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/f9ba59fe55d89bdea4e4809a2407111f_UdrsjCsE7x.png?height=567&lazyload=true&width=1831)

### 2. 填写应用相关信息

- 选择 **企业自建应用**
- 填写 **应用名称**、**应用描述**、**应用图标**
![](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/ba3ad3761c35cc10f71f9c2c6cfe9515_GlcsHqXvwY.png?height=1044&lazyload=true&width=2676)
![](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/915b47c621e86e056a0ed86d68049267_704w7ZReNV.png?height=1306&lazyload=true&width=2644)

### 3. 完成应用创建

- 创建完成后，在开发者后台可以看到已经创建的企业自建应用
- 点击 **应用**，进入 **应用详情页**
![](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/687dbc666a2d910ef3e4210c108eb840_cNDagO40lE.png?height=1016&lazyload=true&width=2742)
### 4. 应用基础信息完善

- 在 **应用详情页**，选择 **凭证与基础信息**，找到 **综合信息**、**国际化配置**，完善应用图标等相关信息
![](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/befdc6eb99552170aba25040b761e074_tA7O60n6En.png?height=949&lazyload=true&width=1812)
### 5. 开启**原生集成应用**

- 在 **应用详情页**，找到 **添加应用能力** 下的 **原生集成应用**，选择 **添加**。 （若无 **原生集成应用**，联系飞书项目经理开通功能白名单）
- 添加完成后，找到 **原生集成应用** ，编辑 **原生集成应用配置** 中的 **最低兼容飞书版本**
- 配置完成后，添加 **飞书主导航**，并配置**应用名称**、**图标**，勾选 **配置到移动端导航栏**
**注意事项**：此处不要添加任何其他应用能力，如：小程序、网页、机器人等

![image.png](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/a46d9a5ac17f9c950bcbcc8e6208f35d_wpAR0WbpfN.png?height=818&lazyload=true&width=1897)

![image.png](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/24525c9a5292f1055ebdc50ffa8b779b_v8C00oHIuy.png?height=809&lazyload=true&width=1905)

![image.png](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/192d7d38cac1afe71ee4846917439208_JrVk3xYFMU.png?height=840&lazyload=true&width=1898)

### 6. 创建版本并线上发布应用

- 在 **应用详情页**，找到 **应用发布** 下的 **版本管理与发布**，选择 **创建版本** 按钮，填写 **版本号**、**移动端的默认能力**、**更新说明**、**可用范围**、**申请理由**
- 开发阶段，此处的 **可用范围** 请配置仅开发相关人员可见。开发测试完成后，可配置成全员可见
- 填写完成后，点击 **申请线上发布**，等待通过即可
![](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/4de53af9fd40fca69835bdbe6fa8e722_uPwecMFF0I.png?height=1264&lazyload=true&width=2836)
![image.png](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/609f6a984bfbc6dcde648c063cf32e30_m6SvasQnJW.png?height=964&lazyload=true&width=1915)
![image.png](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/3340b561948b8e5deb2eaaff10429ebe_YwM6Dc3ipb.png?height=732&lazyload=true&width=1909)

## 二、管理后台配置应用

### 1. [管理后台](http://admin.feishu.cn) 配置应用

- 在 **飞书管理后台** — **工作台** — **应用管理** 下，按照创建的应用名称，找到企业自建应用
- 在应用配置里，勾选 **应用可用范围**，配置应用**可用成员**，
**注意事项**：在开发阶段，此处的 **可用范围** 请配置仅开发相关人员可见；开发测试完成后，请注意配置成 **全部成员**
![image.png](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/8d0a429ab61121925507250678f6ee96_rupNaCtcFH.png?height=593&lazyload=true&width=1905)

### 2. 将应用配置到客户端导航栏
在 **飞书管理后台** — **企业文化** — **功能配置** — **客户端导航栏配置**下，点击 **编辑导航**，添加企业自建应用到移动端导航栏中
![](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/fda14366bf28533fe072e49654483c4b_TFRp66jIWh.png?height=1130&lazyload=true&width=2840)![](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/910554623679f8dcc7be0f1eea7863f7_q6OgfgrR8U.png?height=1354&lazyload=true&width=2408)![](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/4a103295bc3b94119392476b2ec97daa_uZgp29jPip.png?height=1226&lazyload=true&width=2792)

## 三、应用开发接入

在 **飞书开放平台** - **开发者后台** 进入到应用详情，找到 **凭证与基础信息** 里的 **APP ID**，可基于此应用身份ID进行后续开发。
![](https://sf3-cn.feishucdn.com/obj/open-platform-opendoc/728cd5c80ba33e01099a5f6cd324c041_JuYd5XiHs0.png?height=1264&lazyload=true&width=2706)

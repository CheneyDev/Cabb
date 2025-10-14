---
title: "BeaconInfo"
source_url: https://open.feishu.cn/document/client-docs/gadget/-web-app-api/device/ibeacon/ibeaconinfo
last_remote_update: 2022-12-06
last_remote_update_timestamp: 1670294358000
---
最后更新于 2022-12-06

# BeaconInfo

## Beacon的简介

Beacon是建立在蓝牙低功耗 (BLE) 协议基础上的一种广播协议。

Beacon 设备作为蓝牙低功耗协议中的外围设备，持续向周围广播包含设备标识的特定数据包，但不能和中心设备建立连接。小程序运行的设备作为中心设备，可以收到 Beacon 设备的广播包，实现数据交互。常用于室内定位、消息推送等场景。

Beacon的标准包括信号数据的格式等，其中苹果在2013年推出的iBeacon是目前最早且使用广泛的Beacon标准（参考：[iBeacon](https://developer.apple.com/ibeacon/)）。

## Beacon的用途

可以用Beacon进行室内定位（车库/商场），智能打卡，提醒（离开某物体时，比如离开家）。

## BeaconInfo 数据结构

|属性                |类型      |说明       |
|------------------|--------|---------|
|uuid              |string  |iBeacon 设备广播的 uuid|
|major          |number  |iBeacon 设备的主 id|
|minor              |number  |iBeacon 设备的次 id|
|proximity      |number|表示设备距离的枚举值|
|accuracy|number|iBeacon 设备的距离|
|rssi         |number  |表示设备的信号强度|

## 已知问题

- Android中`major`以及`minor`字段为string类型，在V5.26.0版本已修复。

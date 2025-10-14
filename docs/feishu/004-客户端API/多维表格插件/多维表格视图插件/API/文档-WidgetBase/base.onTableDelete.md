---
title: "base.onTableDelete"
source_url: https://open.feishu.cn/document/base-extension/base-view-extensions/api/base/base_ontabledelete
last_remote_update: 2023-07-26
last_remote_update_timestamp: 1690343201000
---
最后更新于 2023-07-26

# base.onTableDelete
监听数据表删除事件，将返回一个取消监听函数。

## 权限要求
**Notice**：开启以下任一权限
查看、评论、编辑和管理多维表格(bitable:app)
查看、评论和导出多维表格(bitable:app:readonly)

## 输入
```js
const off = base.onTableDelete((event) => {})
```

| 名称     | 数据类型 |  是否必填 | 描述 |
| ----------- | ----------- | ------- | --------- |
| event      | {data:{}}       | 否      |	data为一个空对象，后续将支持更多信息。      |

## 输出
取消监听的函数。
## 示例代码

```js
const off = base.onTableDelete((event) => {
	off(); // 监听一次数据表删除事件
	console.log('删除了一个数据表')
})

```

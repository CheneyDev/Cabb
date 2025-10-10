# 知识库插件

使用插件将仓库的文档导入到 CNB 的知识库中，以支持搜索或 LLM 问答。

目前仅支持 Markdown 文件。

## 使用方法

使用过程分为两步：

1. 使用插件将仓库的文档导入到 CNB 的知识库中
2. 使用 CNB 的 Open API 进行召回

## 插件

### 插件镜像名

cnbcool/knowledge-base

### 参数说明

| 参数名 | 说明 | 默认值 | 备注 |
| ------ | ---- | ------ | ---- |
| `embedding_model` | 嵌入模型 | - | 目前只支持 `hunyuan` |
| `include` | 指定需要包含的文件 | `*`（所有文件） | 使用 glob 模式匹配，默认包含所有文件 |
| `exclude` | 指定需要排除的文件 | 空 | 使用 glob 模式匹配，默认不排除任何文件 |

> 注意：`exclude`的优先级高于`include`，被 exclude 的文件不会再被 include

### 插件在 CNB 中使用

```yaml
main:
  push:
    - stages:
        - name: build knowledge base
          image: cnbcool/knowledge-base
          settings:
            embedding_model: hunyuan
            include: "**/**.md"
            exclude: ""
```

## 使用 CNB 的 Open API 进行知识库召回

此 API 用于查询知识库内容，根据提供的查询关键词返回相关信息。

> CNB Open API 使用教程：[https://docs.cnb.cool/zh/open-api.html](https://docs.cnb.cool/zh/open-api.html)
> 访问令牌需要权限：`repo-code:r`（读取仓库代码）

### 接口信息

- **URL**: `https://api.cnb.cool/{slug}/-/knowledge/base/query`
- **方法**: POST
- **内容类型**: application/json

> 注意: `{slug}` 应替换为仓库 slug，例如 CNB 官方文档知识库为 `cnb/docs`

### 请求参数

请求体应为 JSON 格式，包含以下字段:

| 参数名 | 类型   | 必填 | 描述                 |
| ------ | ------ | ---- | -------------------- |
| query  | string | 是   | 要查询的关键词或问题 |

#### 请求示例

```json
{
    "query": "云原生开发配置自定义按钮"
}
```

### 响应内容

响应为 JSON 格式，包含一个结果数组，每个结果包含以下字段:

| 字段名   | 类型   | 描述                     |
| -------- | ------ | ------------------------ |
| score    | number | 匹配相关性分数，范围 0-1 |
| chunk    | string | 匹配的知识库内容文本     |
| metadata | object | 内容元数据               |

#### metadata 字段详情

| 字段名   | 类型   | 描述                 |
| -------- | ------ | -------------------- |
| hash     | string | 内容的唯一哈希值     |
| name     | string | 文档名称             |
| path     | string | 文档路径             |
| position | number | 内容在原文档中的位置 |
| score    | number | 匹配相关性分数       |

#### 响应示例

```json
[
    {
        "score": 0.8671732,
        "chunk": "该云原生远程开发解决方案基于Docker...",
        "metadata": {
            "hash": "15f7a1fc4420cbe9d81a946c9fc88814",
            "name": "quick-start",
            "path": "vscode/quick-start.md",
            "position": 0,
            "score": 0.8671732
        }
    }
]
```

### 使用示例

#### cURL 请求示例

注意：`{slug}` 是运行知识库插件的仓库 slug，例如 CNB 官方文档知识库为 `cnb/docs`

```bash
curl -X "POST" "https://api.cnb.cool/{slug}/-/knowledge/base/query" \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${token}" \
  -d '{
    "query": "云原生开发配置自定义按钮"
}'
```

### 注意事项

1. 响应中的 `chunk` 字段包含了匹配到的知识库内容片段，为 Markdown 格式的文本。
2. `score` 值越高表示匹配度越高。
3. URL 中的 `slug` 是运行知识库插件的仓库 slug。

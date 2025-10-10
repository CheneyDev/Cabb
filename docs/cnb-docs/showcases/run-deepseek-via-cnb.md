# 在 CNB 中运行 DeepSeek 模型

本文地址：https://cnb.cool/examples/ecosystem/deepseek

本文档摘要：通过在 CNB 中运行 DeepSeek 模型，实现快速体验和生产环境部署。支持 1.5b/7b/8b/14b/32b/70b/671b 模型。

快速体验 DeepSeek-R1，支持 1.5b/7b/8b/14b/32b/70b/671b，无需等待，零帧起步。

## 快速体验

### 通过云原生开发体验 1.5～70b

1. `Fork` 本仓库到自己的组织下：https://cnb.cool/examples/ecosystem/deepseek
1. 选择喜欢的分支，点击 `云原生开发` 启动远程开发环境
1. 约 `5～9` 秒后，进入远程开发命令行，输入以下命令即可体验

```shell
ollama run $ds
```

### 通过云原生开发体验 671b 满血版

体验地址见：[DeepSeek-R1-Q8_0](https://cnb.cool/ai-models/deepseek-ai/DeepSeek-R1-GGUF/DeepSeek-R1-Q8_0)

### 免部署体验

适用于无需部署，仅对话的场景，按下 `/` 键直接提问即可。

### 本地部署并体验

在本地或云主机上运行并体验，以 1.5b 为例，可以这样：

```shell
docker run --rm -it docker.cnb.cool/examples/ecosystem/deepseek/1.5b:latest
ollama serve &
ollama run deepseek-r1:1.5b
```

域名 `docker.cnb.cool` 已对腾讯云全局内网加速，无流量费用，不同模型推荐的资源如下：

- `1.5b` - 4核8G内存
- `7b` - 8核16G内存
- `8b` - 8核16G内存
- `14b` - 16核32G内存
- `32b` - 32核64G内存 或 16核24G显存
- `70b` - 16核44G显存

## 生产环境部署

基于 DeepSeek 开发 AI 应用，以下方式可用于部署：

- [基于 DeepSeek 开发 AI 应用并持续部署到 HAI](https://cnb.cool/examples/ecosystem/deploy-deepseek-hai)
- [基于云应用部署 DeepSeek 到指定VPC](https://cloud.tencent.com/document/product/1689/115961)

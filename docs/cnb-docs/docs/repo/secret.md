---
title: 密钥仓库
permalink: https://docs.cnb.cool/zh/repo/secret.html
summary: 密钥仓库是专为存储敏感信息设计的安全型代码仓库，通过严格的访问控制、页面水印和审计追踪等机制确保敏感数据的安全。它不能克隆到本地或推送代码，但可在Web界面编辑，且所有使用都需经过严格权限控制和审计。最佳实践建议对敏感信息进行分类存储，按需分配权限，并定期轮换密钥，同时审计复盘以确保安全。
---

密钥仓库是 `云原生构建` 提供的**安全型代码仓库**，专为存储敏感信息（如密码、API密钥、证书等）设计。
通过严格的访问控制、页面水印、审计追踪等机制，实现敏感数据的安全存储与合规使用。

## 创建密钥仓库

1. **进入仓库创建页面**  
   [点击此处创建密钥仓库](https://cnb.cool/new/repos)（需登录账号）。

2. **选择仓库类型**  
   在仓库类型中选择 **`密钥仓库`**，填写仓库名称与简介。

   ![创建密钥仓库](https://docs.cnb.cool/images/create-secret.png)  

3. **创建**  
   点击创建按钮，创建密钥仓库

## 密钥仓库核心特点

### 1. **安全限制**

| 能力             | 普通仓库 | 密钥仓库 |
| ---------------- | -------- | -------- |
| Git Clone 到本地 | ✅        | ❌        |
| 本地推送代码     | ✅        | ❌        |
| 页面编辑文件     | ✅        | ✅        |
| 创建分支/Tag/PR  | ✅        | ✅        |
| 被流水线引用     | ✅        | ✅        |

### 2. **安全增强特性**

- **动态水印**：页面自动添加当前用户名的半透明水印，防止截图泄露。
- **引用审计**：记录所有引用此仓库文件的流水线记录，支持溯源追踪。
- **强制页面操作**：仅支持在 Web 界面编辑文件，禁止本地操作。
- **严苛的权限控制**: 参考[权限说明](../guide//role-permissions.md)。
- **声明式使用范围**: 参考流水线文件引用[权限检查](../build/file-reference.md#权限检查)。

## 在流水线中引用密钥仓库文件

### 密钥仓库新增文件

```yaml
# env.yml
DOCKER_USER: "username"
DOCKER_TOKEN: "token"
DOCKER_REGISTRY: "https://xxx/xxx"
```

### 导入为环境变量

在流水线配置中，通过 [imports](../build/grammar.md#Pipeline-imports) 字段引用密钥仓库文件，自动注入为环境变量：

```yaml
# .cnb.yml
main:
  push:
    - services:
        - docker
      imports:
        # 引用密钥仓库文件
        - https://cnb.cool/<your-repo-slug>/-/blob/main/xxx/env.yml
      stages:
        - name: docker push
          script: |
            docker login -u ${DOCKER_USER} -p "${DOCKER_TOKEN}" ${CNB_DOCKER_REGISTRY}
            docker build -t ${DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest .
            docker push ${DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:latest
```

### 流水线引用鉴权

默认情况下密钥仓库文件只能被管理员、负责人触发的流水线引用，参考[角色权限](../guide/role-permissions.md)。

若希望团队普通成员触发的流水线也能引用密钥仓库文件，可额外配置 `allow_slugs`、`allow_events`、`allow_branches`、`allow_images` 等字段控制被访问范围。

此时会忽略触发者是否拥有密钥仓库的权限，转而依次检查声明的 `allow_*` 属性，全部通过才能导入为环境变量。更多信息请参考文件引用[权限检查](../build/file-reference.md#权限检查)。

## 最佳实践

1. **敏感信息分类存储**
   - 按环境（prod/dev）或项目拆分不同密钥仓库。
   - 使用 `yaml`、`json` 等文件格式存储管理密钥。
   - 独立于业务组织外新建组织管理密钥仓库，缩小可访问密钥仓库成员范围。

1. **按需使用**
   - 慎重分配管理员、负责人角色。
   - 检查流水线配置是否有滥用、泄漏密钥仓库文件内容的情况。
   - 密钥仓库文件合理配置 `allow_*` 属性。如：
     - 配置 `allow_slugs` 智能被指定范围的仓库流水线引用。
     - 配置 `allow_events` 只能被 `Tag` 相关事件流水线引用。
     - 配置 `allow_branches` 为特定的保护分支，PR 必须经过评审。
     - 指定或制作用于发布、部署的插件，配置 `allow_images` 只能被这些插件引用。

1. **定期轮换密钥**  
   通过页面编辑功能更新密钥后，所有引用此文件的流水线将自动获取新值。

1. **审计复盘**  
   检查审计日志，清理无效引用，撤销离职成员权限。

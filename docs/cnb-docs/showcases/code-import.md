# CNB Code Import(CNB代码仓库批量导入工具)

本文地址：https://cnb.cool/cnb/plugins/cnbcool/code-import

文档摘要：介绍CNB代码仓库批量导入工具的使用方法，包括从CODING、Github、Gitlab、Gitee、阿里云以及通用第三方代码托管平台导入代码仓库。

## 功能介绍
1. 支持CODING、Github、Gitlab、Gitee、阿里云以及通用第三方代码托管平台的代码仓库批量导入至CNB
2. 自动跳过迁移成功的仓库(⚠️依赖工作目录下的`successful.log`文件)

## 在Docker上使用

### 注意事项
1. ⚠️开始导入前，请确保CNB根组织已存在。  
2. `xxx`为需要用户自行替换的字段，具体含义详见参数介绍-核心参数。  
3. 云原生开发自带 docker 命令，内网运行更快速

### 从 Coding 导入

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

### 从 Github 导入

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://github.com" \
  -e PLUGIN_SOURCE_PLATFORM="github" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

### 从 Gitlab 导入

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://gitlab.com" \
  -e PLUGIN_SOURCE_PLATFORM="gitlab" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

### 从 Gitee 导入

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_TOKEN="xxx"  \
  -e PLUGIN_SOURCE_URL="https://gitee.com" \
  -e PLUGIN_SOURCE_PLATFORM="gitee" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

### 从阿里云导入

```shell
docker run --rm  \
  -e PLUGIN_SOURCE_USERNAME="xxx"  \
  -e PLUGIN_SOURCE_PASSWORD="xxx"  \
  -e PLUGIN_SOURCE_PLATFORM="aliyun" \
  -e PLUGIN_SOURCE_AK="xxx" \
  -e PLUGIN_SOURCE_AS="xxx" \
  -e PLUGIN_SOURCE_ORGANIZATIONID="xxx" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

### 从通用第三方代码平台导入

http协议
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_USERNAME="xxx"  \
  -e PLUGIN_SOURCE_PASSWORD="xxx"  \
  -e PLUGIN_SOURCE_REPO="group1/repo1,group1/repo2,group2/repo3" \
  -e PLUGIN_SOURCE_URL="https://common.example.com" \
  -e PLUGIN_SOURCE_PLATFORM="common" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````
ssh协议

⚠️使用ssh协议时请在当前工作目录确保有对应的私钥文件，文件名固定为`ssh.key`
```shell
docker run --rm  \
  -e PLUGIN_SOURCE_REPO="group1/repo1,group1/repo2,group2/repo3" \
  -e PLUGIN_SOURCE_URL="https://common.example.com" \
  -e PLUGIN_SOURCE_PLATFORM="common" \
  -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
  -e PLUGIN_CNB_TOKEN="xxx"  \
  -e PLUGIN_MIGRATE_SSH="true" \
  -e GIT_SSH_COMMAND='ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null' \  
  -v $(pwd):$(pwd) -w $(pwd) \
  cnbcool/code-import
````

## 参数介绍

- PLUGIN_SOURCE_URL
    - 类型：字符串
    - 必填：是
    - 默认值：https://e.coding.net
    - 说明：源仓库代码托管平台URL
    - Ex:
       - github: https://github.com
       - gitlab: https://gitlab.com
       - gitee: https://gitee.com
       - common: https://common.com

- PLUGIN_SOURCE_TOKEN
    - 类型：字符串
    - 必填：是
    - 默认值：-
    - 说明：调用源代码托管平台 API 的 token (当 source_platform 不为 common 时必填)
        - CODING权限：用户信息-只读、项目信息-只读、代码仓库-只读 https://e.coding.net/user/account/setting/tokens
        - Github权限：repo:all、read:org https://github.com/settings/tokens
        - Gitlab权限：read_api https://gitlab.com/-/user_settings/personal_access_tokens
        - Gitee权限：user_info、projects https://gitee.com/profile/personal_access_tokens

- PLUGIN_SOURCE_PLATFORM
    - 类型：字符串
    - 必填：是
    - 默认值：coding
    - 说明：导入的平台名称，支持 coding/gitlab/github/gitee，其他通用平台填写 common

- PLUGIN_SOURCE_REPO
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：仓库路径，需要与source_url拼接成完整的源仓库http克隆地址,如https://common.com/group1/repo1  
  当source_platform 为 common 或者 source_platform 为 coding 且 migrate_type 为 repo 时必填，多个代码仓库以英文逗号隔开
    - Ex: group1/repo1,group1/repo2,group2/repo3

- PLUGIN_SOURCE_USERNAME
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：当 source_platform 为 aliyun或common 时必填，clone 代码仓库时要用到的用户名，需要确保能够clone所有仓库。
      [阿里云帮助文档](https://help.aliyun.com/zh/yunxiao/user-guide/configure-https-clone-account-password?spm=a2c4g.11186623.0.0.78b240cdASV98n)

- PLUGIN_SOURCE_PASSWORD
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：当 source_platform 为 aliyun或common 时必填，clone 代码仓库时要用到的密码

- PLUGIN_SOURCE_AK
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：AccessKey ID,当 source_platform 为 aliyun时必填，需要有AliyunRDCReadOnlyAccess权限，如果是RAM用户，需要关联至云效账号，并授权管理员角色。
      [阿里云帮助文档](https://help.aliyun.com/zh/yunxiao/user-guide/add-a-ram-user?spm=5176.28366559.console-base_help.dexternal.211e336a7R37d8&scm=20140722.S_help%40%40%E6%96%87%E6%A1%A3%40%40203014.S_RQW%40ag0%2BBB2%40ag0%2BBB1%40ag0%2Bos0.ID_203014-RL_ram%E7%94%A8%E6%88%B7%E5%A6%82%E4%BD%95%E5%85%B3%E8%81%94%E8%87%B3%E4%BA%91%E6%95%88-LOC_console~UND~help-OR_ser-V_4-P0_0-P1_0)

- PLUGIN_SOURCE_AS
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：AccessKey Secret,当 source_platform 为 aliyun时必填。

- PLUGIN_SOURCE_ENDPOINT
    - 类型：字符串
    - 必填：否
    - 默认值：devops.cn-hangzhou.aliyuncs.com
    - 说明：AccessKey 请求的地址。

- PLUGIN_SOURCE_ORGANIZATIONID
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：阿里云云效代码仓库企业ID，可在云效访问链接中获取，如https://devops.aliyun.com/organization/【OrganizationId】

- PLUGIN_CNB_URL
    - 类型：字符串
    - 必填：是
    - 默认值：https://cnb.cool
    - 说明：CNB访问URL

- PLUGIN_CNB_TOKEN
    - 类型：字符串
    - 必填：是
    - 默认值：-
    - 说明：CNB 授权令牌，个人令牌-访问令牌创建。权限要求：
        - repo-code 读写
        - repo-basic-info 只读
        - account-profile 只读
        - account-engage 只读
        - group-resource 读写
        - group-manage 读写
        - repo-content 读写

- PLUGIN_CNB_ROOT_ORGANIZATION
    - 类型：字符串
    - 必填：是
    - 默认值：-
    - 说明：迁移后，CNB对应的根组织名称，请确保根组织已提前创建,不需要带/
    - Ex: root-group

- PLUGIN_MIGRATE_TYPE
    - 类型：字符串
    - 必填：否
    - 默认值：team
    - 说明：要迁移的类型，支持项目(project)、仓库(repo)、团队(team)多维度迁移，只支持 coding 平台

- PLUGIN_MIGRATE_CONCURRENCY
    - 类型：数值
    - 必填：否
    - 默认值：10
    - 说明：仓库迁移并发数，最大10

- PLUGIN_MIGRATE_FORCE_PUSH
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：强制push到CNB仓库

- PLUGIN_MIGRATE_SKIP_EXISTS_REPO
    - 类型：布尔值
    - 必填：否
    - 默认值：true
    - 说明：跳过CNB已存在的仓库

- PLUGIN_MIGRATE_USE_LFS_MIGRATE
    - 类型：字符串
    - 必填：否
    - 默认值：false
    - 说明：是否使用lfs migrate处理历史提交中超过CNB单文件最大限制错误
      ⚠️如开启该配置，迁移后commit ID会与源仓库不一致

- PLUGIN_MIGRATE_ORGANIZATION_MAPPING_LEVEL
    - 类型：字符串
    - 必填：否
    - 默认值：1
    - 说明：CODING与CNB组织映射关系，仅支持 Coding 平台
      1: CODING项目映射为CNB子组织，仓库在子组织下面
      2: CODING项目不会映射为CNB子组织，仓库直接在CNB根组织下面

- PLUGIN_MIGRATE_ALLOW_INCOMPLETE_PUSH
    - 类型：字符串
    - 必填：否
    - 默认值：false
    - 说明：针对LFS源文件丢失的仓库，忽略LFS文件下载报错，LFS推送时忽略丢失的对象报错，继续推送

- PLUGIN_MIGRATE_LOG_LEVEL
    - 类型：字符串
    - 必填：否
    - 默认值：info
    - 说明：日志级别(debug/info/warn/error)

- PLUGIN_MIGRATE_RELEASE
    - 类型：布尔值
    - 必填：否
    - 默认值：false
    - 说明：迁移release（暂时只支持 gitlab release迁移）

- PLUGIN_MIGRATE_FILE_LIMIT_SIZE
    - 类型：数值
    - 必填：否
    - 默认值：100
    - 说明：CNB最大文件大小限制，单位Mib

- PLUGIN_MIGRATE_CODE
    - 类型：布尔值
    - 必填：是
    - 默认值：true
    - 说明：迁移代码

- PLUGIN_MIGRATE_SSH
    - 类型：布尔值
    - 必填：是
    - 默认值：false
    - 说明：使用ssh协议克隆通用第三方平台代码仓库

- PLUGIN_MIGRATE_REBASE
    - 类型：布尔值
    - 必填：是
    - 默认值：false
    - 说明：在源和目标都有变更的情况下,且CNB侧仓库对应分支根目录有`.cnb.yml`文件，使用git rebase保证代码同步，CNB侧提交的流水线配置代码不会被覆盖
      ⚠️如开启该配置，将启用强制推送（PLUGIN_MIGRATE_FORCE_PUSH="true"），并在迁移工具执行的工作目录备份CNB侧代码仓库

- PLUGIN_SOURCE_PROJECT
    - 类型：字符串
    - 必填：否
    - 默认值：-
    - 说明：要迁移的 CODING 项目名称 (当 source_platform 为 coding 且 migrate_type 为 project 时必填)，多个项目以英文逗号隔开

## 常见问题
1. 单个文件超过100MB大小限制  
可以开启`PLUGIN_MIGRATE_USE_LFS_MIGRATE`参数

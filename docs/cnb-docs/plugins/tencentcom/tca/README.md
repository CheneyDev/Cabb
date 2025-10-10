# 腾讯云代码分析 - TCA

- 腾讯云代码分析（Tencent Cloud Code Analysis，简称TCA，内部曾用研发代号CodeDog）起步于 2012 年，是集众多代码分析工具的云原生、分布式、高性能的代码综合分析跟踪管理平台，其主要功能是持续跟踪分析代码，观测项目代码质量，支撑团队传承代码文化。
- 用心关注每行代码迭代，助力维护卓越代码文化！
- 官网地址：[tca.tencent.com](https://tca.tencent.com)
- 开源版GitHub地址: [Tencent/CodeAnalysis](https://github.com/Tencent/CodeAnalysis)

## 前置操作

- 请先在腾讯云代码分析官网([https://tca.tencent.com](https://tca.tencent.com/))上创建团队、项目， 接入代码库， 创建分析方案，详情请参考：[cnb-tca-插件配置文档](https://tca.tencent.com/document/zh/guide/%E6%8F%92%E4%BB%B6%E9%85%8D%E7%BD%AE.html#cnb-tca-%E6%8F%92%E4%BB%B6%E9%85%8D%E7%BD%AE)
- 分析方案-规则配置中，推荐使用以下规则包：
  - `【全语种】开源合规检查`规则包，适用于各语言项目开源合规检查，主要针对敏感信息和开源license的检查，协助开发者规避开源协议风险。
  - `Xcheck安全规则包`,适用于多种语言(Go、Java、Nodejs、PHP、Python)的代码安全检查，覆盖包括SQL注入、代码注入、命令注入、跨站脚本、反序列化漏洞、路径穿越等多种漏洞。（属于TCA增强分析能力，可[申请License体验](https://tca.tencent.com/document/zh/guide/License.html#%E5%AE%98%E7%BD%91%E7%89%88-license)）
- 创建分析方案后，请接入代码分析节点，参考：[TCA节点管理](https://tca.tencent.com/document/zh/guide/%E8%8A%82%E7%82%B9%E7%AE%A1%E7%90%86.html)

## 在 云原生构建 上使用

### 1. 创建密钥仓库

- 由于插件需要用到TCA的个人token（个人敏感信息，请注意保密），需要通过密钥仓库的方式引入。请创建一个密钥仓库（比如`tca-private-config`）,在仓库中创建一个`tca-settings.yml`文件，内容示例如下：
- **请从tca官网代码库页面`CNB - TCA 插件配置`中直接拷贝，参考：[cnb-tca-插件配置文档](https://tca.tencent.com/document/zh/guide/%E6%8F%92%E4%BB%B6%E9%85%8D%E7%BD%AE.html#cnb-tca-%E6%8F%92%E4%BB%B6%E9%85%8D%E7%BD%AE)**

```yaml
token: xxxxxxx
allow_images:
  - "tencentcom/tca-plugin:latest"
allow_slugs:
  - "组织名/仓库名"
allow_events:
  - "**"
```

密钥仓库yaml文件参数说明：

- `token`: TCA个人令牌。
- `allow_images`: 设置可访问密钥的镜像，此处只允许tca-plugin镜像访问。
- `allow_slugs`: 设置可访问密钥的代码仓库，默认设置为当前仓库，如有多个仓库需要访问，可按需调整，支持glob表达式。
- `allow_events`: 设置可访问密钥的触发事件，此处默认设置为全部，可按需修改为特定事件，支持glob表达式。
- 后三个参数为CNB文件引用的权限控制参数，详见: [CNB文件引用文档](https://docs.cnb.cool/zh/file-reference.html)

### 2. 配置.cnb.yml

- **请从tca官网代码库页面`CNB - TCA 插件配置`中直接拷贝，参考：[cnb-tca-插件配置文档](https://tca.tencent.com/document/zh/guide/%E6%8F%92%E4%BB%B6%E9%85%8D%E7%BD%AE.html#cnb-tca-%E6%8F%92%E4%BB%B6%E9%85%8D%E7%BD%AE)**

```yml
# .cnb.yml
main:  # 触发的分支名，按需修改
  push:  # push触发，也可以用merge_request等触发
  - stages:
    # 代码分析
    - name: TCA
      image: tencentcom/tca-plugin:latest
      settings:
        org_sid: xxx  # 团队编号，从TCA官网获取
        team_name: xxx  # 项目名称，从TCA官网获取
        scheme_id: xxx  # 分析方案id，从TCA官网获取
      settingsFrom:
        - https://cnb.cool/xxx/tca-private-config/-/blob/main/tca-settings.yml
```

## 参数说明

### org_sid

- type: String
- required: 是
- 团队编号，从TCA官网获取

### team_name

- type: String
- required: 是
- 项目名称，从TCA官网获取

### scheme_id

- type: String
- required: 是
- 分析方案id，从TCA官网获取

### total_scan

- type: Boolean
- 是否全量扫描，可选值: true, false, 默认为 false, 即增量扫描。
- required: 否

### scan_dir

- type: String
- required: 否
- 填写目录的相对路径，指定代码仓库下的子目录作为扫描目录，适用于大仓场景只扫描某个模块目录，默认不配置，为扫描代码仓库根目录。
- 示例1：sub_dir
- 示例2：sub_dir_1/sub_dir_2

### settingsFrom

- 填写创建好的密钥仓库yaml文件URL，将从该文件中引入token信息。

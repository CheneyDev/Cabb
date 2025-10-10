# 二进制软件成分分析

[二进制软件成分分析](https://cloud.tencent.com/document/product/1483/63351)
Binary Software Composition Analysis，BSCA 是一款基于二进制分析能力的自动化软件成分分析工具。

BSCA 可对二进制构建产物进行分析，例如固件、APK、镜像、jar 包等格式。BSCA 聚焦于已知漏洞扫描、开源软件审计和敏感信息检测。
BSCA 无需源代码，一键上传目标文件，就可以输出安全报告，帮助您高效识别风险，节省安全成本，提升安全竞争力。

## 参数说明

- secret_id：必填，腾讯云 API SecretId，用于标识 API 调用者的身份，可在腾讯云平台的访问管理 > 访问密钥 > API 密钥管理 中获取。
  建议配置环境变量并加密，防止 SecretId 泄露。
- secret_key：必填，腾讯云 API SecretKey，用于加密和验证签名字符串，可在腾讯云平台的访问管理 > 访问密钥 > API 密钥管理 中获取。
  建议配置环境变量并加密，防止 SecretKey 泄露。
- file_path：必填，待分析文件所在路径，示例：/bin/echo 。
- analysis_name：必填，创建分析的名称。
- analysis_type：必填，创建分析的类型，与文件类型直接相关。目前支持：
  - RTOS：RTOS 固件
  - DOCKER：Docker 镜像
  - GENERIC：其他制品文件（各类二进制构建产物），默认为“其他制品文件”
- analysis_param：非必填，分析项参数，目前仅 RTOS 固件类型的分析需要，包括 bin、s19、hex。

## 在 云原生构建 上使用

上传 docker 镜像进行检测：

```yaml
# .cnb.yml
main:
  push:
  - services:
      - docker
    stages:
    - name: docker build & push & save
      script: 
        - docker build -t security-audit:latest .
        - docker save security-audit:latest > security-audit.tar
    # 上传 tar 格式镜像文件到安全审计平台
    - name: 安全审计
      image: tencentcom/bsca-analysis
      imports:
        - https://cnb.xxx.vom/xxx/your-envs.yml
      settings: 
        secret_id: $SECRET_ID
        secret_key: $SECRET_KEY
        file_path: ./security-audit.tar
        analysis_type: DOCKER
        analysis_name: $CNB_BUILD_ID
```

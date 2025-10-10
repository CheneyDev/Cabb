# KSAS 安全审计

T-Sec Keenlab Security Audit Suites(KSAS) 是用于各种级别安全审计的工具集合，其中包括工具：系统基线安全审计、制品安全审计、源码安全审计、安卓应用安全审计。

## 参数说明

1、必选参数

+ `team_name`：分析所属团队，需提前创建。示例：`your-team`。
+ `project_name`：分析所属项目，需提前创建。示例：`your-project`。
+ `analysis_path`：待分析文件所在路径。示例：`/bin/echo`。
+ `analysis_type`：创建分析的类型，目前支持三种，默认为构建包：
  + `ArtifactPackage`(构建包)：对应文档中的分析类型 binAuditor，系统类型 Package
  + `ArtifactDocker`(Docker镜像)：对应文档中的分析类型 binAuditor，系统类型 Docker。
  + `ArtifactAPK`(APK)：安卓应用安全审计
  + `ArtifactSource`(源码)：源码安全审计
+ `analysis_name`：分析名称/版本名称。示例：`your-analysis-name`。需要保持唯一
+ `token`：Audit Suites 令牌，可在 Audit Suites web 页面控制台右上角，用户名-令牌管理-新建，获取到令牌。
+ `website`：Audit Suites 的 URL，用于 API 请求、最终报告链接生成，示例：`https://your-url.com`。

2、可选参数

+ `description`: 分析描述，默认值为 'cnb plugin create'。

3、分析配置相关可选参数

+ `analyze_timeout`: 解包超时设置，默认 10 分钟
+ `analyze_file_type`: 选择分析的文件类型（Text/Binary），默认分析全部类型。
+ `file_skiped`: 要跳过的文件或目录。
例如：/skipped/path/*、*/pom.xml 或 */*_test.go。默认无忽略文件。
+ `smart_sca_enabled`: 是否开启 C/C++ SCA，默认开启。
+ `smart_sca_deep_scan`: 是否开启深度扫描，
默认关闭，仅当 C/C++ SCA 开启时有效。
+ `unpacker_enabled`: 是否开启启发式解包，默认开启。
+ `extract_depth`: 解包递归深度，默认 0 表示动态递归。

4、分析配置中，源码审计相关配置参数，可选参数

+ `source_sca_recommended`: 必填，是否使用分析器推荐配置。
此选项为 true 时，则忽略 `source_sca_min_file_count`、`source_sca_snippet_match_method`、`source_sca_snippet_match_rate`、`source_sca_snippet_match_count`、`source_sca_small_file_line_count`。
+ `source_sca_min_file_count`: 检出一个组件的最小文件数量
（命中同一组件的文件数量，<= MinFileCount 时，不检出组件），
取值范围 minFileCount >= 1。
+ `source_sca_snippet_match_method`: 源码匹配方法。
+ `source_sca_snippet_match_rate`: 最小匹配度（命中特征数/总特征数量），
取值范围 0 <= SnippetMatchRate <= 1（百分比）。
+ `source_sca_snippet_match_count`: 最小匹配行数。取值范围 snippetMatchCount >= 0。
+ `source_sca_small_file_line_count`: 去除注释后，特征数量小于一定值不参与匹配。
取值范围 smallFileLineCount >= 1。

5、高级配置可选参数

+ `advanced_config`: 高级配置（JSON格式）。
该参数不为空时，上方3、4两部分的分析配置参数全部失效。

## 在 云原生构建 上使用

上传源码到安全审计平台检测：

```yaml
# .cnb.yml
main:
  push:
  - services:
      - docker
    stages:
    # 将源码打包
    - name: 打包源码
      script: 
        - tar -czvf securityaudit.tar.gz ./src
    # 上传源码到安全审计平台
    - name: 安全审计
      image: tencentcom/ksas-audit-suites
      imports:
        - https://cnb.xxx.vom/xxx/your-envs.yml
      settings: 
        team_name: teamname
        project_name: projectname
        analysis_path: ./securityaudit.tar.gz
        analysis_type: ArtifactSource
        analysis_name: $CNB_BUILD_ID
        website: https://your-url
        token: $TOKEN
```

上传 docker 镜像到安全审计平台检测：

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
      image: tencentcom/ksas-audit-suites
      imports:
        - https://cnb.xxx.vom/xxx/your-envs.yml
      settings: 
        team_name: teamname
        project_name: projectname
        analysis_path: ./security-audit.tar
        analysis_type: ArtifactDocker
        analysis_name: $CNB_BUILD_ID
        website: https://your-url
        token: $TOKEN
```

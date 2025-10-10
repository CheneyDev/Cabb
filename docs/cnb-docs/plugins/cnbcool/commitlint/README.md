# commitlint

基于约定式提交规范，检查提交注释是否符合给定的规范。

## 配置方式

项目根目录放置自定义配置文件，配置文件查找优先级

- .commitlintrc.yml
- commitlint.config.js
- .commitlintrc.js

 未添加该配置时，将使用以下默认配置：

```yaml
# 默认 .commitlintrc.yml 配置
extends:
  - "@commitlint/config-conventional"
```

以下插件已经被默认支持

- @commitlint/config-conventional
- commitlint-plugin-references

## 约定式提交

与开源社区规范一致，详见 [https://www.conventionalcommits.org/zh/](https://www.conventionalcommits.org/zh/)

## 在 云原生构建 中使用

```yaml
# 插件内部使用 git 命令获取提交注释，无需传入提交注释
master:
  pull_request:
  - stages:
    - name: commitlint
      image: cnbcool/commitlint
```

## 允许中文

```yaml
#.commitlintrc.yml

extends:
  - "@commitlint/config-conventional"

rules:
  #允许中文
  subject-case: [0]  
```

## 要求提交记录中必须有 issue 引用

使用如下配置可以要求提交记录中必须带有#号开头的 issue 单号。

```yaml
#.commitlintrc.yml

extends:
  - "@commitlint/config-conventional"

parserPreset:
  parserOpts:
    # issue 是以#开头的
    issuePrefixes: ["#"]

rules:
  #允许中文
  subject-case: [0]  
  #必须有引用单据
  references-empty: [2, 'never'] 
```

## 只对指定的 type 开启关联验证

```yaml
#.commitlintrc.yml

extends:
  - "@commitlint/config-conventional"

plugins:
  - references

parserPreset: 
  parserOpts: 
    issuePrefixes: ["#",'--bug=']  

rules:
  #允许中文
  subject-case: [0]
  # feat、fix 必须有单关联
  references-empty-enum: [2,"never",["fix", "feat"]]
    
```

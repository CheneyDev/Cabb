---
title: 环境变量
permalink: https://docs.cnb.cool/zh/build/env.html
summary: 在云原生构建过程中，可以通过多种方式声明、导入和管理环境变量，支持在 `Pipeline` 和 `Job` 级别的设置，以及对敏感信息进行加密处理。变量名必须符合 POSIX 标准，且变量值的长度不能超过 100KiB。在使用时，变量可以在脚本任务、配置文件和多个内置及插件任务中进行替换和引用。
---

在构建过程中可以使用环境变量，`云原生构建` 内置了一些[默认环境变量](./build-in-env.md)。

## 声明环境变量

1. 通过 `env` 声明环境变量
1. 在 `Pipeline` 里声明的环境变量对当前 `Pipeline` 有效
1. 在 `Job` 里声明的环境变量对当前 `Job` 有效

```yaml
main:
  push:
    - env:
        RELEASE_MSG: release message
      stages:
        - name: release
          type: git:release
          env:
            RELEASE_NAME: release-1
          options:
            description: ${RELEASE_MSG}
            name: ${RELEASE_NAME}
```

## 导入环境变量

1. 通过 `imports` 导入一个密钥仓库文件，可将敏感信息注入到环境变量，供后续任务使用。
1. 当 `env` 和 `imports` 的 `key` 冲突时，优先使用 `env`

```yaml
main:
  push:
    - services:
        - docker
      # 导入密钥仓库文件为环境变量
      imports: https://cnb.cool/<your-repo-slug>/-/blob/main/xxx/envs.yml
      stages:
        - name: docker info
          script: docker info
        - name: docker login
          # 其中 TEST_DOCKER_DOMAIN、TEST_DOCKER_USER、TEST_DOCKER_PWD 为密钥仓库文件中的变量
          script: docker login $TEST_DOCKER_DOMAIN -u $TEST_DOCKER_USER -p $TEST_DOCKER_PWD
```

`https://cnb.cool/<your-repo-slug>/-/blob/main/xxx/envs.yml` 内容示例

```yaml
# docker 镜像源域名
TEST_DOCKER_DOMAIN: registry.example.com
# docker 用户名
TEST_DOCKER_USER: your_docker_username
# docker 密码
TEST_DOCKER_PWD: your_docker_password
```

## 变量名限制

在 Shell 中，环境变量名的命名规则有一些限制。根据 POSIX 标准，环境变量名应符合以下规则：

1. 只能包含字母（大小写均可）、数字和下划线（_）字符。
1. 第一个字符不能是数字。

不符合上述规则的变量会被忽略

## 导出环境变量

`Job` 执行结束后，有一个 `result` 对象，可通过  [exports](./grammar.md#exports) 将 `result` 中的属性导出到环境变量，生命周期为当前 `Pipeline`。

语法格式为：

```yaml
exports:
  from-key: to-key
```

- `from-key` 要导出的 `Job result` 对象中的属性名，支持环境变量，支持深层取值，参考 `lodash.get`。
- `to-key` 映射到环境变量中的变量名。

有如下三种方式设置 result：

- 脚本任务执行结果

- 解析输出中的自定义变量

- 内置任务的 result

### 脚本任务执行结果

`script` 自定义脚本任务执行完，`Job result`的属性有：

- `code`: 返回码
- `stdout`: 标准输出
- `stderr`: 标准错误
- `info`: 标准输出、标准错误，按时序的混合体

**注意：可使用 `printf "%s" "hello\nworld"` 来输出变量，以消除标准输出流最后的换行符，同时保留 `\n` 等转义字符。**

```yaml
main:
  push:
    - stages:
        - name: set env
          script: echo -n $(date "+%Y-%m-%d %H:%M")
          exports:
            code: CUSTOM_ENV_DATE_CODE
            info: CUSTOM_ENV_DATE_INFO
        - name: echo env
          script:
            - echo $CUSTOM_ENV_DATE_CODE
            - echo $CUSTOM_ENV_DATE_INFO
```

包含 `if`、`ifModify`、`ifNewBranch` 等判断逻辑时，可设置的属性有：

- `skip`: 如果执行上述判断逻辑后跳过任务执行则返回跳过原因，否则为空字符串。

```yaml
- name: use if
  if: exit -1
  exports:
    skip: REASON
- name: tell last
  # $REASON 的值为 if 这个字符串
  script: echo $REASON
```

### 解析输出中的自定义变量

CI 会从标准输出流里按行识别 `##[set-output key=value]` 格式的内容，自动放入 `result` 对象中。

若变量值包含换行符 `\n`，可对变量值进行 `base64` 或 `escape` 编码。

变量值若以 `base64,` 开始，`云原生构建` 会对 `base64,` 后面的内容做 `base64` 解码，否则会对变量值做 `unescape` 解码。

使用 `Node.js` 示例代码如下：

```javascript
// test.js
const value = '测试字符串\ntest string';
// 输出 base64 编码的变量值
console.log(`##[set-output redline_msg_base64=base64,${Buffer.from(value, 'utf-8').toString('base64')}]`);

// 输出 escape 编码的变量值
console.log(`##[set-output redline_msg_escape=${escape(value)}]`)
```

```yaml
main:
  push:
    - docker:
        image: node:20-alpine
      stages:
        - name: set output env
          script: node test.js
          # 将 test.js 输出的变量导出为环境变量
          exports:
            redline_msg_base64: BASE64_KEY
            redline_msg_escape: ESCAPE_KEY
        - name: echo env
          script:
            - echo "BASE64_KEY $BASE64_KEY"
            - echo "ESCAPE_KEY $ESCAPE_KEY"
```

使用 `echo` 示例代码如下：

```yaml
main:
  push:
    - stages:
        - name: set output env
          script: echo "##[set-output redline_msg_base64=base64,$(echo -e "测试字符串\ntest string" | base64 -w 0)]"
          exports:
            redline_msg_base64: BASE64_KEY
        - name: echo env
          script:
            - echo -e "BASE64_KEY $BASE64_KEY"
```

**注意：在 Unix-like 系统中，base64 命令默认会在每76个字符后添加一个换行符，可使用 `-w 0` 选项来禁用换行避免 CI 未能按行解析出变量。**

不包含 `\n` 的变量值可直接输出

```bash
echo "##[set-output redline_msg=some value]"
```

:::tip
受限于系统环境变量值长度限制，过大的变量值无效。

CI 会忽略大于等于 `100KB` 的变量值。可写入文件中自行解析。

对于敏感信息，建议使用 [read-file](./internal-steps/#read-file) 内置任务。
:::

### 内置任务中导出环境变量

一些内置任务会有输出结果，可通过 `exports` 导出为环境变量。

```yaml
main:
  push:
    - stages:
        - name: xxxx
          type: xxx:xxx
          options:
            product: public
            name: cnb
            dist: release/
          exports:
            version: CUSTOM_ENV_VERSION
            url: CUSTOM_ENV_URL
            #支持对象深层取值
            nextRelease.gitTag: CUSTOM_ENV_GIT_TAG
        - name: echo env
          script:
            - echo $CUSTOM_ENV_VERSION
            - echo $CUSTOM_ENV_URL
```

`result` 内容参考各内置任务文档。

### 增删改查环境变量

可以覆盖已有的环境变量，设为空字符串或 null 即为删除。

```yaml
main:
  push:
    - env:
        CUSTOM_ENV_DATE_INFO: default
        CUSTOM_ENV_FOR_DELETE: default
      stages:
        - name: set env
          script: echo -n $(date "+%Y-%m-%d %H:%M")
          exports:
            # 新增
            code: CUSTOM_ENV_DATE_CODE
            # 修改
            info: CUSTOM_ENV_DATE_INFO
            # 删除
            CUSTOM_ENV_FOR_DELETE: null
            # 删除
            # CUSTOM_ENV_FOR_DELETE:
        - name: echo env
          script:
            - echo $CUSTOM_ENV_DATE_CODE
            - echo $CUSTOM_ENV_DATE_INFO
            - echo $CUSTOM_ENV_DATE_STDOUT
            - echo $CUSTOM_ENV_FOR_DELETE
            - echo $CUSTOM_ENV_GIT_TAG
```

## 使用环境变量

### 在 脚本任务 中使用

执行脚本任务时，流水线设置的环境变量作为任务执行时的环境变量

```yaml
main:
  push:
    - stages:
        - name: test internal env
          # CNB_BRANCH 为内置环境变量
          script: echo $CNB_BRANCH
        - name: test self defined env
          env:
            cat_name: tomcat
          script: echo $cat_name
```

### 变量替换

配置文件中的一些属性值会进行变量替换。

如有环境变量 `env_name=env_value`，那么属性值中 `$env_name` 会被替换成 `env_value`，若 `env_name` 无值，则会替换成空字符串。

下面列出了会进行变量替换的属性：

- 内置任务

内置任务 options 内的属性值和 optionsFrom 会进行变量替换。

```yaml
# options.yml
name: Nightly
```

```yaml
main:
  push:
    - env:
        address: options.yml
        description: publish for xx task
      stages:
        - name: git release
          type: git:release
          # $address 会被替换成 env 中的 "options.yml"
          optionsFrom: $address
          # options.yml 中的 name 会合并到 options 中
          options:
            # $description 会被替换成 env 中的 "publish for xx task"
            description: $description
```

options 的最终内容为：

```yaml
name: Nightly
description: publish for xx task
```

- 插件任务

插件任务 settings 内的属性值和 settingsFrom 会进行变量替换。

```yaml
# settings.yml
robot: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx
```

```yaml
main:
  push:
    - env:
        address: settings.yml
        message: pr check
      stages:
        - name: notify
          image: tencentcom/wecom-message
          # $address 会被替换成 env 中的 "settings.yml"
          settingsFrom: $address
          # settings.yml 中的 robot 会合并到 settings 中
          settings:
            # $message 会被替换成 env 中的 "pr check"
            content: $message
```

settings 的最终内容为：

```yaml
robot: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx
message: pr check
```

另外写在 Dockerfile 中的 `LABEL` 指定的 `settingsFrom` 同样可以进行变量替换

```dockerfile
FROM node:20

LABEL cnb.cool/settings-from="$address"
```

- env

env 下声明的属性值可以引用上层 env 中的变量进行替换

```yaml
# .cnb.yml
main:
  push:
    - env:
        cat_name: tomcat
      stages:
        - name: echo env
          env:
            # 使用上层 env 声明的 cat_name 值进行替换
            name: "cat $cat_name"
            # 输出 cat tomcat
            script: echo $name
```

- imports

imports 的属性值以及所声明的文件中的属性值均会进行变量替换。

若 imports 为数组，数组前面文件中声明的变量对数组后面元素有效。

```yaml
# env1.yml
address: env2.yml
platform: amd64
```

```yaml
# env2.yml
# 读取 env1.yml 中的 platform 属性值进行替换
action: build for $platform
```

```yaml
# .cnb.yml
main:
  push:
    - imports:
        - env1.yml
        # env1.yml 中声明了 address，$address 会被替换成 env2.yml
        - $address
      stages:
        - name: echo action
          # 读取 env2.yml 中的 action 属性值进行替换
          script: echo $action
```

- pipeline.runner.tags
  
```yaml
# 构建不同架构下的镜像
.build: &build
  runner:
    tags: cnb:arch:$CNB_PIPELINE_NAME
  services:
    - docker
  stages:
    - name: docker build
      script: echo "docker build for $CNB_PIPELINE_NAME"
main:
  push:
    # 下面 "amd64" 和 "arm64:v8" 会被声明为内置环境变量 CNB_PIPELINE_NAME 的值
    amd64: *build
    "arm64:v8": *build
```

- pipeline.docker.volumes

```yaml
.docker-volume: &docker-volume
  docker:
    image: node:22-alpine
    volumes:
      - $volume_path
main:
  push:
    install:
      env:
        volume_path: node_modules
      <<: *docker-volume
      stages:
        - name: install
          script: npm install axios
        # 通知其他流水线执行
        - name: resolve
          type: cnb:resolve
          options:
            key: install
    build:
      env:
        volume_path: node_modules
      <<: *docker-volume
      stages:
        # 等待 install 流水线
        - name: await
          type: cnb:await
          options:
            key: install
        - name: ls
          script: ls node_modules
```

- ifModify

```yaml
# 不同模块下代码有变更才进行对应模块的编译
.build: &build
  ifModify: $CNB_PIPELINE_NAME/*
  stages:
    - name: build $CNB_PIPELINE_NAME
      script: echo "build $CNB_PIPELINE_NAME"
main:
  push:
    module-a: *build
    module-b: *build
```

- name

pipeline.name 、stage.name 和 job.name 的属性值会进行变量替换。

```yaml
main:
  push:
    - name: build in $CNB_REPO_SLUG
      env:
        platform: amd64
      imports:
        - env1.yml
        - env2.yml
      stages:
        - name: stage_$SOME_ENV
          script: echo "hello world"
```

- lock.key

```yaml
# env.yml
build_key: build key
```

```yaml
.build: &build
  imports: env.yml
  lock:
    key: $build_key
  stages:
    - name: echo
      script: echo "hello world"
main:
  push:
    # 以下两条流水线，一条占用了锁成功执行，另一条位占到锁执行失败
    - *build
    - *build
```

- allowFailure

```yaml
main:
  push:
    - env:
        allow_fail: true
      stages:
        - name: echo
          allowFailure: $allow_fail
          # 脚本执行会报错，但 allowFailure 为 true，任务被认为是成功的
          script: echo1 1
```

### 阻止变量替换

如果不希望 `$env_name` 被替换，可以通过 `\$` 阻止替换

```yaml
main:
  push:
    - stages:
        - name: git release
          type: git:release
          options:
            name: Development
            # 属性值是 "some code update $description"
            description: some code update \$description
```

## 限制

环境变量名需由字母、数字或 `_` 组成，不能以数字开头。

变量值的长度不能超过 `100KiB`。

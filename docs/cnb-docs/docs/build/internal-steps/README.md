---
title: 内置任务
permalink: https://docs.cnb.cool/zh/build/internal-steps.html
summary: 该文本介绍了CNB（Cloud Native Build）的多种内置任务，包括Docker缓存、cnb任务（如await-resolve）、vscode、git和测试相关任务，以及artifact相关任务。每个任务都有其适用事件、参数、输出结果及配置样例，可用于丰富CI/CD流程。
---

在每一阶段（Stage）操作中，除了使用自定义任务，
CNB 还扩展了一些常用的内置任务，以供开发者使用。

## Docker

### cache  {#docker-cache}

`docker:cache`

==Docker 缓存==，构建一个 `Docker` 镜像作为缓存，在未来的构建中重复使用。

可以避免网络资源如 `依赖包` 重复下载。

- [适用事件](./#docker-cache-applicable-events)
- [参数](./#docker-cache-parameters)
- [输出结果](./#docker-cache-output)
- [配置样例](./#docker-cache-configuration-examples)

#### 适用事件 {#docker-cache-applicable-events}

[所有事件](../trigger-rule.md#trigger-event)

#### 参数 {#docker-cache-parameters}

- [dockerfile](./#docker-cache-parameters-dockerfile)
- [by](./#docker-cache-parameters-by)
- [versionBy](./#docker-cache-parameters-versionBy)
- [buildArgs](./#docker-cache-parameters-buildArgs)
- [target](./#docker-cache-parameters-target)
- [sync](./#docker-cache-parameters-sync)
- [ignoreBuildArgsInVersion](./#docker-cache-parameters-ignoreBuildArgsInVersion)

##### dockerfile {#docker-cache-parameters-dockerfile}

- type: `String`
- required: `true`

用于构建缓存镜像的 Dockerfile 路径。

为避免超长时间构建，Docker 镜像构建超时时间受 [job.timeout](../grammar.md#timeout) 参数控制。

##### by {#docker-cache-parameters-by}

- type: `Array<String>` | `String`
- required: `false`

用来声明缓存镜像构建过程中依赖的文件列表。**注意：未出现在 by 列表中的文件，除了 Dockerfile，其他在构建镜像过程中，都当不存在处理。**

- 支持数组格式
- 支持字符串格式，多个文件用英文逗号分隔。

##### versionBy {#docker-cache-parameters-versionBy}

- type: `Array<String>` | `String`
- required: `false`

用来进行版本控制，未传入 `versionBy`，则默认取 `by` 的值进行版本控制。

`versionBy` 所指向的文件内容发生变化，我们就会认为是一个新的版本， 具体的计算逻辑见这个表达式：sha1(Dockerfile + versionBy + buildArgs + target + arch)。

- 支持数组格式。
- 支持字符串格式，多个文件用英文逗号分隔。

##### buildArgs {#docker-cache-parameters-buildArgs}

- type: `Object`
- required: `false`

在 build 时插入额外的构建参数 (`--build-arg $key=$value`), value 值为 null 时只加入 key (`--build-arg $key`)。

##### target {#docker-cache-parameters-target}

- type: `String`
- required: `false`

对应 docker build 中的 --target 参数，可以选择性地构建 Dockerfile 中的特定阶段，而不是构建整个 Dockerfile。

##### sync {#docker-cache-parameters-sync}

- type: `Boolean`
- required: `false`
- default: `false`

是否同步模式，等待缓存镜像 `docker push` 成功后才继续。

##### ignoreBuildArgsInVersion {#docker-cache-parameters-ignoreBuildArgsInVersion}

- type: `Boolean`
- required: `false`
- default: `false`

版本计算是否忽略 `buildArgs`。

详见`版本控制`

#### 输出结果 {#docker-cache-output}

```javascript
{
  // 缓存对应的 docker image name
  name
}
```

#### 配置样例 {#docker-cache-configuration-examples}

```yaml
main:
  push:
    - docker:
        image: node:14
      stages:
        - name: build cache image
          type: docker:cache
          options:
            dockerfile: cache.dockerfile
            # by 支持以下两种形式：数组、字符串
            by:
              - package.json
              - package-lock.json
            # versionBy: package-lock.json
            versionBy:
              - package-lock.json
          exports:
            name: DOCKER_CACHE_IMAGE_NAME
        - name: use cache
          image: $DOCKER_CACHE_IMAGE_NAME
          # 将 cache 中的文件拷贝过来使用
          commands:
            - cp -r "$NODE_PATH" ./node_modules
        - name: build with cache
          script:
            - npm run build
```

其中，`cache.dockerfile` 是一个用于构建缓存镜像的 Dockerfile。示例：

```dockerfile
# 选择一个 Base 镜像
FROM node:14

# 设置工作目录
WORKDIR /space

# 将 by 中的文件列表 COPY 过来
COPY . .

# 根据 COPY 过来的文件进行依赖的安装
RUN npm ci

# 设置好需要的环境变量
ENV NODE_PATH=/space/node_modules
```

## cnb

### await {#cnb-await}

参见 [await-resolve](./#cnb-await-resolve)

### resolve {#cnb-resolve}

参见 [await-resolve](./#cnb-await-resolve)

### await-resolve {#cnb-await-resolve}

- `cnb:await`
- `cnb:resolve`

`await` 会等待 `resolve` 的执行，`resolve` 可以向 `await`传递变量。

通过 `await-resolve`，可以让多个并发的 `pipeline` 相互配合，实现更灵活的顺序控制。

:::tip

`await-resolve` 同 `apply`、`trigger` 的区别：

前者指一个构建中某 `pipeline` 执行到 `await` 任务时等待对应 `key` 的 `resolve` 通知才会继续进行。

后者指一个 `pipeline` 触发新事件，开启新的构建。可以跨仓库。可以异步调用，也可以同步等待。
:::

- [使用限制](./#cnb-await-resolve-usage-limitations)
- [死锁检测](./#cnb-await-resolve-deadlock-detection)
- [适用事件](./#cnb-await-resolve-applicable-events)
- [await 参数](./#cnb-await-resolve-await-parameters)
- [resolve 参数](./#cnb-await-resolve-resolve-parameters)
- [输出结果](./#cnb-await-resolve-output)

#### 使用限制 {#cnb-await-resolve-usage-limitations}

1. 只能对同一个事件触发的 `pipeline` 进行 `await` 和 `resolve` 操作
1. 一个 `key` 仅能 `resolve` 一次，但可以 `await` 多次
1. 通过 `key` 对 `await` 和 `resolve` 进行分组

#### 死锁检测 {#cnb-await-resolve-deadlock-detection}

`await` 和 `resolve` 相互配合可以完成灵活的流程控制，但也会引入更复杂的边界情况，比如：

1. `pipeline-1` 和 `pipeline-2` 相互 `await`，即：死锁
1. 多条 `pipeline` 间存在 `await` 环，即：间接死锁
1. `await` 一个不存在的 `key`，或者 `key` 没有关联 `resolve`，即：无限等待
1. `resolve` 所在 `pipeline` 执行失败，对应的 `await` 陷入无限等待
1. 多个 `resolve` 任务关联同一个 `key`，即重复 `resolve` 抛出异常

`死锁检测` 机制会自动检测以上异常，结束 `await` 的等待状态，抛出 `dead lock found.` 异常。

`await` 和 `resolve` 的在配置文件中顺序不影响运行结果，即最后 `await` 任务一定是会等待相应的 `resolve` 完成，这种情况不会被`死锁检测`机制终止。

#### 适用事件 {#cnb-await-resolve-applicable-events}

[所有事件](../trigger-rule.md#trigger-event)

#### await 参数  {#cnb-await-resolve-await-parameters}

##### key {#cnb-await-resolve-await-parameters-key}

- type: String
- required: true

配对 ID

##### resolve 参数 {#cnb-await-resolve-resolve-parameters}

###### key {#cnb-await-resolve-resolve-parameters-key}

- type: String
- required: true

配对 ID

###### data {#cnb-await-resolve-resolve-parameters-data}

- type: object
- required: false

要传递的对象

`key: value` 格式，支持多级。示例：

```yaml{4-8}
- name: resolve a json
  type: cnb:resolve
  options:
    key: demo
    data:
      a: 1
      b:
        c: 2
```

`await` 任务的结果，是 `resolve` 声明的 `data` 对象。
可以通过 `exports` 访问这个对象，示例：

```yaml{3-7}
- name: await a json
  type: cnb:await
  options:
    key: demo
  exports:
    a: VAR_A
    b.c: VAR_B
- name: show var
  script:
    - echo ${VAR_A} # 1
    - echo ${VAR_B} # 2
```

当然，也可以不传送任务内容，仅仅表示一个等待动作：

```yaml
- name: ready
  type: cnb:resolve
  options:
    key: i-am-ready
```

```yaml
- name: ready
  type: cnb:await
  options:
    key: i-am-ready
```

###### 输出结果 {#cnb-await-resolve-output}

```javascript
{
  // resolve 返回的 data 内容
  data
}
```

### apply {#cnb-apply}

`cnb:apply`

- [适用事件](./#cnb-apply-applicable-events)
- [参数](./#cnb-apply-parameters)
  - [环境变量相关](./#cnb-apply-environment-variables)
  - [config 取值优先级](./#cnb-apply-config-priority)
- [输出结果](./#cnb-apply-output)
- [配置样例](./#cnb-apply-configuration-examples)

#### 适用事件 {#cnb-apply-applicable-events}

- `push`
- `commit.add`
- `branch.create`
- `pull_request.target`
- `pull_request.mergeable`
- `tag_push`
- `pull_request.merged`
- `api_trigger`
- `web_trigger`
- `crontab`
- `tag_deploy`

#### 参数 {#cnb-apply-parameters}

- [config](./#cnb-apply-parameters-config)
- [configFrom](./#cnb-apply-parameters-configFrom)
- [event](./#cnb-apply-parameters-event)
- [sync](./#cnb-apply-parameters-sync)
- [continueOnBuildError](./#cnb-apply-parameters-continueOnBuildError)
- [title](./#cnb-apply-parameters-title)

##### config {#cnb-apply-parameters-config}

- type: `String`
- required: `false`
完整的 CI 配置文件内容

##### configFrom {#cnb-apply-parameters-configFrom}

- type: `String`
- required: `false`
指定一个本地文件作为配置文件。

##### event {#cnb-apply-parameters-event}

- type: `String`
- required: `false`
- default: `api_trigger`
要执行的自定义事件名，必须为 `api_trigger` 或以 `api_trigger_` 开头。

##### sync {#cnb-apply-parameters-sync}

- type: `Boolean`
- required: `false`
- default: `false`
是否同步执行。

同步模式下会等待本次 apply 流水线执行成功，再执行下一个任务。

##### continueOnBuildError {#cnb-apply-parameters-continueOnBuildError}

- type: `Boolean`
- required: `false`
- default: `false`
同步模式下，触发的流水线构建失败时，是否继续执行下个任务。

##### title {#cnb-apply-parameters-title}

- type: `String`
- required: `false`
自定义流水线标题

#### 环境变量相关 {#cnb-apply-environment-variables}

当前 `Job` 可见的，业务定义的环境变量全部传递给新的流水线。

默认值中有如下环境变量，**用户无法覆盖**：

- `APPLY_TRIGGER_BUILD_ID`，含义同 CI 默认环境变量中的 `CNB_BUILD_ID`
- `APPLY_TRIGGER_PIPELINE_ID`，含义同 CI 默认环境变量中的 `CNB_PIPELINE_ID`
- `APPLY_TRIGGER_REPO_SLUG`，含义同 CI 默认环境变量中的 `CNB_REPO_SLUG`
- `APPLY_TRIGGER_REPO_ID`，含义同 CI 默认环境变量中的 `CNB_REPO_ID`
- `APPLY_TRIGGER_USER`，含义同 CI 默认环境变量中的 `CNB_BUILD_USER`
- `APPLY_TRIGGER_BRANCH`，含义同 CI 默认环境变量中的 `CNB_BRANCH`
- `APPLY_TRIGGER_COMMIT`，含义同 CI 默认环境变量中的 `CNB_COMMIT`
- `APPLY_TRIGGER_COMMIT_SHORT`，含义同 CI 默认环境变量中的 `CNB_COMMIT_SHORT`
- `APPLY_TRIGGER_ORIGIN_EVENT`，含义同 CI 默认环境变量中的 `CNB_EVENT`
- `CNB_PULL_REQUEST_ID`，含义同 CI 默认环境变量中的 `CNB_PULL_REQUEST_ID`
- `CNB_PULL_REQUEST_IID`，含义同 CI 默认环境变量中的 `CNB_PULL_REQUEST_IID`
- `CNB_PULL_REQUEST_MERGE_SHA`，含义同 CI 默认环境变量中的 `CNB_PULL_REQUEST_MERGE_SHA`

##### config 取值优先级 {#cnb-apply-config-priority}

按以下顺序依次取值，取到为止：

1. `config`
1. `configFrom`
1. 当前仓库 `.cnb.yml`
1. 若在 `pull_request.merged` 事件中调用 `apply` 内置任务，取合并后的配置文件
1. 若在 `pull_request.target`、`pull_request.mergeable` 事件中调用 `apply` 内置任务，取目标分支的配置文件
1. 其他情况取当前分支配置文件

`configFrom` 只支持本地文件如 `./test/.cnb.yml`，远程文件可先自行下到本地。

#### 输出结果 {#cnb-apply-output}

```json
{
  "sn": "cnb-i5o-1ht8e12hi", // 构建号
  "buildLogUrl": "http://xxx/my-group/my-repo/-/build/logs/cnb-i5o-1ht8e12hi", // 构建日志链接
  "message": "success",
  "buildSuccess": true, // 触发的构建是否成功，此key仅在同步模式下存在
  "lastJobExports": {} // 触发的流水线最后一个job导出的环境变量
}
```

#### 配置样例 {#cnb-apply-configuration-examples}

```yaml
main:
  push:
    - stages:
        - name: trigger
          type: cnb:apply
          options:
            configFrom: ./test/.cnb.yml
            event: api_trigger_test
```

```yaml
main:
  push:
    - stages:
        - name: trigger
          type: cnb:apply
          options:
            config: |
              main:
                api_trigger_test: 
                - stages:
                  - name: test
                    script: echo test

            event: api_trigger_test
```

```yaml
main:
  push:
    - stages:
        - name: trigger
          type: cnb:apply
          options:
            # 执行当前配置文件的其它事件
            event: api_trigger_test
  api_trigger_test:
    - stages:
        - name: test
          script: echo test
```

```yaml
main:
  push:
    - stages:
        - name: trigger
          type: cnb:apply
          options:
            configFrom: .xxx.yml
            event: api_trigger_test
            sync: true
```

### read-file {#cnb-read-file}

`cnb:read-file`

读取文件内容解析输出给后续任务。

用 `##[set-output key=value]` 指令输出变量更便捷，但会将内容输出到日志，不适用于敏感信息。

对于事先确定的敏感信息可以用 `imports` 导入，对于构建过程中生成的敏感信息，可写入文件，用该内置任务读取。

`imports` 只支持远端仓库存在的文件，该内置任务只支持读取本地存在的文件。

文件类型通过后缀判断，目前支持 json、 yml(yaml)、 纯文本（`key=value` 结构）。

- [适用事件](./#cnb-read-file-applicable-events)
- [参数](./#cnb-read-file-parameters)
- [输出结果](./#cnb-read-file-output)
- [配置样例](./#cnb-read-file-configuration-examples)

#### 适用事件 {#cnb-read-file-适用事件}

[所有事件](../trigger-rule.md#trigger-event)

#### 参数 {#cnb-read-file-parameters}

##### filePath

- type: String
- required: true

本地文件路径

#### 输出结果 {#cnb-read-file-output}

```json
{
  // 读到的对象
}
```

#### 配置样例 {#cnb-read-file-configuration-examples}

将一些后续任务需要的变量写入文件，用该内置任务读取文件导出为环境变量，供后续任务使用。

```json title="myJson.json"
{
  "myVar": "myVar",
  "deep": {
    "myVar": "myVar in deep"
  },
  "deepWithEnvKey": {
    "myVar": "myVar in deepWithEnvKey"
  }
}
```

```yaml title="cnb.yml"
main:
  push:
    - env:
        myKey: deepWithEnvKey
      stages:
        - name: write env file
          script: echo "write env to myJson.json"
        - name: export env
          type: cnb:read-file
          options:
            filePath: myJson.json
          exports:
            myVar: ENV_MY_VAR
            deep.myVar: ENV_DEEP_MY_VAR #指定多级的key
            $myKey.myVar: ENV_ENV_KEY_MY_VAR #通过环境变量指定获取的key
        - name: echo env
          script:
            - echo $ENV_MY_VAR
            - echo $ENV_DEEP_MY_VAR
            - echo $ENV_ENV_KEY_MY_VAR
```

### trigger {#cnb-trigger}

`cnb:trigger`

在当前仓库中，触发另外一个仓库的自定义事件流水线。

- [适用事件](./#cnb-trigger-applicable-events)
- [参数](./#cnb-trigger-parameters)
- [输出结果](./#cnb-trigger-output)
- [配置样例](./#cnb-trigger-configuration-examples)
  - [基本使用](./#cnb-trigger-configuration-examples-basic-usage)

#### 适用事件 {#cnb-trigger-applicable-events}

[所有事件](../trigger-rule.md#trigger-event)

#### 参数 {#cnb-trigger-parameters}

- [token](./#cnb-trigger-parameters-token)
- [slug](./#cnb-trigger-parameters-slug)
- [event](./#cnb-trigger-parameters-event)
- [branch](./#cnb-trigger-parameters-branch)
- [sha](./#cnb-trigger-parameters-sha)
- [env](./#cnb-trigger-parameters-env)
- [sync](./#cnb-trigger-parameters-sync)
- [continueOnBuildError](./#cnb-trigger-parameters-continueOnBuildError)
- [title](./#cnb-trigger-parameters-title)

##### token {#cnb-trigger-parameters-token}

- type: `String`
- required: `true`

个人访问令牌。

新流水线触发者为令牌对应用户，会判断有无目标仓库权限。

##### slug {#cnb-trigger-parameters-slug}

- type: `String`
- required: `true`

目标仓库的完整路径，如：group/repo。

##### event {#cnb-trigger-parameters-event}

- type: `String`
- required: `true`

触发的自定义事件名，必须是 `api_trigger` 或以 `api_trigger_` 开头。

需要目标仓库配置了对应事件的流水线，才可以触发。

##### branch {#cnb-trigger-parameters-branch}

- type: `String`
- required: `false`

触发分支，默认为主分支。

##### sha {#cnb-trigger-parameters-sha}

- type: `String`
- required: `false`

触发分支中的 CommitId，默认取 branch 的最新提交记录。

##### env {#cnb-trigger-parameters-env}

- type: [TriggerEnv](./#cnb-trigger-parameters-type-definitions-TriggerEnv)
- required: `false`

触发目标仓库流水线时的环境变量。

默认值中有如下环境变量，用户无法覆盖：

- `API_TRIGGER_BUILD_ID`，含义同 CI 默认环境变量中的 `CNB_BUILD_ID`
- `API_TRIGGER_PIPELINE_ID`，含义同 CI 默认环境变量中的 `CNB_PIPELINE_ID`
- `API_TRIGGER_REPO_SLUG`，含义同 CI 默认环境变量中的 `CNB_REPO_SLUG`
- `API_TRIGGER_REPO_ID`，含义同 CI 默认环境变量中的 `CNB_REPO_ID`
- `API_TRIGGER_USER`，含义同 CI 默认环境变量中的 `CNB_BUILD_USER`
- `API_TRIGGER_BRANCH`，含义同 CI 默认环境变量中的 `CNB_BRANCH`
- `API_TRIGGER_COMMIT`，含义同 CI 默认环境变量中的 `CNB_COMMIT`
- `API_TRIGGER_COMMIT_SHORT`，含义同 CI 默认环境变量中的 `CNB_COMMIT_SHORT`

##### sync {#cnb-trigger-parameters-sync}

- type: Boolean
- required: false
- default: false

是否同步执行，同步模式下会等待本次 trigger 流水线执行成功后，再执行下一个任务。

##### continueOnBuildError {#cnb-trigger-parameters-continueOnBuildError}

- type: Boolean
- required: false
- default: false

同步模式下，触发的流水线构建失败时，是否继续执行下个任务。

##### title {#cnb-trigger-parameters-title}

- type: String
- required: false

自定义流水线标题

##### 类型定义 {#cnb-trigger-parameters-type-definitions}

###### TriggerEnv {#cnb-trigger-parameters-type-definitions-TriggerEnv}

```txt
{
    [key: String]: String | Number | Boolean
}
```

#### 输出结果 {#cnb-trigger-output}

```json
{
  "sn": "cnb-i5o-1ht8e12hi", // 构建号
  "buildLogUrl": "http://xxx/my-group/my-repo/-/build/logs/cnb-i5o-1ht8e12hi", // 构建日志链接
  "message": "success",
  "buildSuccess": true, // 触发的构建是否成功，此key仅在同步模式下存在
  "lastJobExports": {} // 触发的流水线最后一个job导出的环境变量
}
```

#### 配置样例 {#cnb-trigger-configuration-examples}

##### 基本使用 {#cnb-trigger-configuration-examples-basic-usage}

在当前仓库触发一个仓库 `main` 分支的事件名为 `api_trigger_test` 的流水线。

该流水线配置文件使用仓库 `main` 分支的 `.cnb.yml` 文件。

使用访问令牌 `$TOKEN` 查询用户是否有仓库权限。

```yaml title="cnb.yml"
main:
  push:
    - stages:
        - name: trigger
          type: cnb:trigger
          imports: https://cnb.cool/<your-repo-slug>/-/blob/main/xxx/envs.yml
          options:
            token: $TOKEN
            slug: a/b
            branch: main
            event: api_trigger_test
```

## vscode

### go {#vscode-go}

`vscode:go`

远程开发中是否配置该内置任务的区别：

- ==使用该任务==： 启动云原生开发时，需等待该任务执行完，才出现 WebIDE 和 VSCode/Cursor 客户端入口。
- ==不使用该任务==： 流水线 prepare 阶段执行完（code-server 代码服务启动），stages 任务执行前，就会出现 WebIDE 和 VSCode/Cursor 客户端入口。

上述入口出现时机区别：仅指从 `loading` 等待页到跳转入口选择页出现的时机。 实际上无论是否使用该任务，在 `code-server` 代码服务启动后，远程开发已经是可用状态。

注意：使用该任务将增加等待时间。如果需要延迟开发者进入时机，在某些任务执行完才允许进入远程开发环境，可使用该任务。

- [适用事件](./#vscode-go-applicable-events)
- [输出结果](./#vscode-go-output)
- [配置样例](./#vscode-go-configuration-examples)

#### 适用事件 {#vscode-go-applicable-events}

- `vscode`
- `branch.create`
- `api_trigger`
- `web_trigger`

#### 输出结果 {#vscode-go-output}

```json
{
  // webide url
  "url": ""
}
```

#### 配置样例 {#vscode-go-configuration-examples}

```yaml title=".cnb.yml"
$:
  # vscode 事件：专供页面中启动远程开发用
  vscode:
    - docker:
        # 使用自定义镜像作为开发环境，未传入此参数，将使用默认镜像 cnbcool/default-dev-env:latest
        image: cnbcool/default-dev-env:latest
      services:
        - vscode
        - docker
      stages:
        # 希望等该任务执行完再进入开发环境
        - name: ls
          script: ls -al
        - name: vscode go
          type: vscode:go
        # 可以在进入开发环境后再执行的任务
        - name: ls
          script: ls -al
```

## git

### auto-merge {#git-auto-merge}

`git:auto-merge`

==自动合并 Pull Request==，一般用于 PR 通过 `pull_request` 流水线检查 和 `Code Review` 后，在 `pull_request.mergeable` 流水线中自动合并 `PR`。无需人工点击合并按钮。

`pull_request.mergeable` 事件触发条件和时机参考[事件](../trigger-rule.md#trigger-event)。

- [适用事件](./#git-auto-merge-applicable-events)
- [参数](./#git-auto-merge-parameters)
- [输出结果](./#git-auto-merge-output)
- [配置样例](./#git-auto-merge-configuration-examples)
  - [单独使用](./#git-auto-merge-configuration-examples-standalone-usage)
  - [配合目标分支的 push 事件使用](./#git-auto-merge-configuration-examples-used-with-target-branch-push-event)
- [最佳实践](./#git-auto-merge-best-practices)
  - [使用 squash 自动合并](./#git-auto-merge-best-practices-using-squash-auto-merge)
  - [使用 auto 自动选择合并类型](./#git-auto-merge-best-practices-using-auto-automatically-select-merge-type)
  - [建立专用的 CR 群交叉走查自动合并](./#git-auto-merge-best-practices-establishing-dedicated-cr-group-cross-review-auto-merge)

#### 适用事件 {#git-auto-merge-applicable-events}

`pull_request.mergeable`

#### 参数 {#git-auto-merge-parameters}

- [mergeType](./#git-auto-merge-parameters-mergeType)
- [mergeCommitMessage](./#git-auto-merge-parameters-mergeCommitMessage)
- [mergeCommitFooter](./#git-auto-merge-parameters-mergeCommitFooter)
- [removeSourceBranch](./#git-auto-merge-parameters-removeSourceBranch)
- [ignoreAssignee](./#git-auto-merge-parameters-ignoreAssignee)

##### mergeType {#git-auto-merge-parameters-mergeType}

- type: `String`
- required: false
- default: auto

- 合并策略，默认为 auto：如果多人提交走 merge ，否则走 squash 。

##### mergeCommitMessage {#git-auto-merge-parameters-mergeCommitMessage}

- type: String
- required: false
- 合并点提交信息。

当合并策略为 `rebase` 的时候，该信息无效，无需填写。

当合并策略为 `merge` ，默认为 `chore: merge node(merged by CNB)` ，会自动追加 `PR` 引用、`reviewers` 名单、`PR` 包含的提交人名单。举例说明：

```txt
chore: merge node(merged by CNB)

PR-URL: !916
Reviewed-By: tom
Reviewed-By: jerry
Co-authored-by: jack
```

当合并策略为 squash 时，默认值为该 pr 的第一条 commit message。并会自动追加 pr 引用和 reviewers 名单、pr 包含的提交人名单。举例说明：

```yaml title=".cnb.yml"
main:
  pull_request.mergeable:
    - stages:
        - name: automerge
          type: git:auto-merge
          options:
            mergeType: squash
```

该配置会产生如下效果：

某个 `PR` (feat/model-a -> main) 中有两条提交记录：

提交记录 1，2023-10-1 日提交

```txt
feat(model-a): 给模块 A 增加一个新特性

由于某某原因，新增某某特性

close #10
```

提交记录 2，修复了在 cr 时被指出的一些问题，2023-10-2 日提交

```txt
fix(model-a): 修复评审中指出的问题
```

在自动合并后将会在 `main` 分支上产生一个这样的提交节点，即后续的提交记录（也就是提交记录 2）将会被抹掉

```txt
feat(model-a): 给模块 A 增加一个新特性

由于某某原因，新增某某特性

close #10

PR-URL: !3976
Reviewed-By: tom
Reviewed-By: jerry
Co-authored-by: jack
```

##### mergeCommitFooter {#git-auto-merge-parameters-mergeCommitFooter}

- type: `String`
- required: `false`

合并点需要设置的脚注，多个脚注用 `\n` 分隔，仅在 `merge` 和 `squash` 生效。

当合并策略为 `rebase` 的时候，该信息无效，无需填写。

当合并策略为 `merge` 或 `squash` 时，会将传入信息按行加入到脚注中，再追加 `PR` 引用、`reviewers` 名单、`pr` 包含的提交人名单。举例说明：

```yaml title=".cnb.yml"
main:
  pull_request.mergeable:
    - stages:
        - name: automerge
          type: git:auto-merge
          options:
            mergeType: squash
            mergeCommitMessage: "add feature for some jobs"
            mergeCommitFooter: "--story=123\n--story=456"
```

```txt
add feature for some jobs

--story=123
--story=456
PR-URL: !916
Reviewed-By: tom
Reviewed-By: jerry
Co-authored-by: jack
```

##### removeSourceBranch {#git-auto-merge-parameters-removeSourceBranch}

- type: Boolean
- required: false
- default: false

合并后是否删除源分支。

源分支与目标分支不同仓库时，该值无效。

##### ignoreAssignee {#git-auto-merge-parameters-ignoreAssignee}

- type: Boolean
- required: false
- default: false

是否忽略 `assignee`。

当该 `PR` 有指定 `assignee` （指派人）的时候，本任务不会执行自动合并的逻辑。 因为 `assignee` 的本意就是指派某人手动来处理。

当为 `true` 时，可忽略 `assignee` 强行合并。

#### 输出结果 {#git-auto-merge-output}

```txt
{
    reviewedBy, // String，追加在提交信息后面的提交者信息
    reviewers, // Array<String>, 走查者列表
}
```

#### 配置样例 {#git-auto-merge-configuration-examples}

##### 单独使用 {#git-auto-merge-configuration-examples-standalone-usage}

```txt
main:
  pull_request.mergeable:
    - stages:
        - name: automerge
          type: git:auto-merge
          options:
            mergeType: merge
```

当分支 `main` 上 的 `PR` 触发了 `pull_request.mergeable` 事件，那么会将这个 `PR` 以 `merge` 的方式自动合并。

##### 配合目标分支的 push 事件使用 {#git-auto-merge-configuration-examples-used-with-target-branch-push-event}

```yaml title=".cnb.yml"
main:
  push:
    - stages:
        - name: build
          script: npm run build
        - name: publish
          script: npm run publish
  pull_request.mergeable:
    - stages:
        - name: automerge
          type: git:auto-merge
          options:
            mergeType: merge
```

当分支 `main` 上 的 `PR` 触发了 `pull_request.mergeable` 事件，那么会将这个 `PR` 以 `merge` 的方式自动合并。 在自动合并之后，会触发目标分支（main）上的 `push` 事件，继续执行声明的 `build` 和 `publish` 流程。

**这里有一个谁应该对产生的 push 事件流程负责的问题，这里的设定是这样的：**

1. 如果提出 `PR` 者，是目标仓库的成员，那么提出 `PR` 的人对后续产生的 `push` 事件流程负责（即 `push` 构建流水会推送到这里）。
2. 如果提出 `PR` 者，不是目标仓库的成员（比如在开源项目中，从 `Forked Repo` 提出 `PR`）。那么最后一个 `Code Review` 的 `Reviewer` 将对后续产生的 `push` 事件流程负责。

#### 最佳实践 {#git-auto-merge-best-practices}

##### 使用 squash 自动合并 {#git-auto-merge-best-practices-using-squash-auto-merge}

使用 `squash` 合并，一次 `PR` 操作只在目标分支产生一个 `commit` 节点，并且在有权限的情况下删除源分支。

```yaml title=".cnb.yml"
main:
  review:
    - stages:
        - name: automerge
          type: git:auto-merge
          options:
            mergeType: squash
            removeSourceBranch: true
```

##### 使用 auto 自动选择合并类型 {#git-auto-merge-best-practices-using-auto-automatically-select-merge-type}

如果多人提交走 `merge` ，否则走 `squash`。

```yaml title=".cnb.yml"
main:
  review:
    - stages:
        - name: automerge
          type: git:auto-merge
          options:
            mergeType: auto
```

##### 建立专用的 CR 群交叉走查自动合并 {#git-auto-merge-best-practices-establishing-dedicated-cr-group-cross-review-auto-merge}

1. 走查通过后自动合并
2. 记录走查者信息在提交信息中

```yaml title=".cnb.yml"
main:
  pull_request.mergeable:
    - stages:
        - name: CR 通过后自动合并
          type: git:auto-merge
          options:
            mergeType: squash
            mergeCommitMessage: $CNB_LATEST_COMMIT_MESSAGE
          exports:
            reviewedBy: REVIEWED_BY
        - name: notify
          image: tencentcom/wecom-message
          settings:
            robot: "155af237-6041-4125-9340-000000000000"
            msgType: markdown
            content: |
              > CR 通过后自动合并 <@${CNB_BUILD_USER}> 
              > 　
              > ${CNB_PULL_REQUEST_TITLE}
              > [${CNB_EVENT_URL}](${CNB_EVENT_URL})
              > 
              > ${REVIEWED_BY}

  pull_request:
    - stages:
        # ...省略其它任务
        - name: notify
          image: tencentcom/wecom-message
          options:
            robot: "155af237-6041-4125-9340-000000000000"
            msgType: markdown
            content: |
              > ${CURR_REVIEWER_FOR_AT}
              > 
              > ${CNB_PULL_REQUEST_TITLE}
              > [${CNB_EVENT_URL}](${CNB_EVENT_URL})
              > 
              > from ${CNB_BUILD_USER}
```

### issue-update {#git-issue-update}

`git:issue-update`

==更新 Issue 状态==，关闭或打开 `Issue`，修改 `Issue` 标签。

- [适用事件](./#git-issue-update-applicable-events)
- [工作机制](./#git-issue-update-work-mechanism)
- [Issue ID 获取方式](./#git-issue-update-how-to-get-issue-id)
- [提交日志中如何带上 Issue ID](./#git-issue-update-how-to-include-issue-id-in-commit-logs)
- [参数](./#git-issue-update-parameters)
- [输出结果](./#git-issue-update-output)
- [配置样例](./#git-issue-update-configuration-examples)

#### 适用事件 {#git-issue-update-applicable-events}

[所有事件](../trigger-rule.md#trigger-event)

#### 工作机制 {#git-issue-update-work-mechanism}

查找 `Issue` 是否存在 -> 检查是否符合 when 条件（可选） -> 检查是否符合 lint 条件（可选）-> 更新 issue 状态或标签

#### Issue ID 获取方式 {#git-issue-update-how-to-get-issue-id}

- 如果传入了 `fromText` 参数，则从 `fromText` 中获取。其中 `branch.delete` 必传 `fromText`（分支已删除，无法从上下文中获取）。
- 如果没传入 `fromText` 参数，则从上下文中获取。
  - `push` 事件：从本次推送的所有 `commit` 的提交日志中获取。
  - 合并请求类事件：从合并请求中的所有 `commit` 的提交日志中获取。
  - 其他事件：从最新的一个 `commit` 的提交日志中获取。

获取方式：提取上述文本中如下两种格式

- #issueID：表示当前仓库的 issue。例如 #123，表示当前仓库 id 为 123 的 issue。
- groupName/repoName#issueID：表示跨仓库（其他仓库）的 issue。例如：test/test#123，表示 test/test 仓库中 id 为 123 的 issue。

注意 #123 或 test/test#123 前需要有空格。

#### 提交日志中如何带上 Issue ID {#git-issue-update-how-to-include-issue-id-in-commit-logs}

提交代码时，可以在提交日志中加上关联的 `Issue ID`， 可以在使用当前内置任务时可以自动提取到关联 `Issue`，用来更新 `Issue` 标签和状态

推荐在提交日志的 body 中带上 `Issue ID`，命令行操作方式如下：

- 方法一：用 `shift + enter` 换行，建议 title 和 body 之间加上一个空行

```txt
git commit -m "fix(云原生构建): 修复一个错误

cnb/feedback#123"
```

- 方法二：

以下提交方式 title 和 body 之间会产生两个换行

```txt
git commit -m "fix(云原生构建): 修复一个错误" -m "cnb/feedback#123"
```

#### 参数 {#git-issue-update-parameters}

- [fromText](./#git-issue-update-parameters-fromText)
- [state](./#git-issue-update-parameters-state)
- [label](./#git-issue-update-parameters-label)
- [when](./#git-issue-update-parameters-when)
- [lint](./#git-issue-update-parameters-lint)
- [defaultColor](./#git-issue-update-parameters-defaultColor)

##### [fromText](./#git-issue-update-parameters-fromText)

- type: `String`
- required: `false`

从给定的文本中解析 `Issue Id`。

不声明时，自动从上下文里的提交记录解析。

> 可以指定一个包含 issue id 引用的文本来声明操作对象，比如： ${LATEST_CHANGE_LOG} 。

##### [state](./#git-issue-update-parameters-state)

- type: [IssueStateMap](#git-issue-update-parameters-type-definitions-IssueStateMap)
- required: `false`

对应 `state` 属性，为 `close` 时，可关闭 `Issue`。

##### [label](./#git-issue-update-parameters-label)

- type: [IssueUpdateLabel](#git-issue-update-parameters-type-definitions-IssueUpdateLabel)
- required: `false`

对 `label` 的操作描述。

##### [when](./#git-issue-update-parameters-when)

- type: [IssueUpdateStatus](#git-issue-update-parameters-type-definitions-IssueUpdateStatus)
- required: `false`

过滤条件，多个条件之间是 `or` 关系。为空时表示对所有 `Issue` 操作。

##### [lint](./#git-issue-update-parameters-lint)

- type: [IssueUpdateStatus](#git-issue-update-parameters-type-definitions-IssueUpdateStatus)
- required: `false`

检查 `Issue` 是否满足条件，不满足时抛出异常，多个条件之间是 or 关系，为空时表示不做检查。

##### [defaultColor](./#git-issue-update-parameters-defaultColor)

- type: `String`
- required: `false`

添加的标签的默认颜色，当有传入 `label.add` 参数时才有效。

##### 类型定义 {#git-issue-update-parameters-type-definitions}

###### IssueStateMap {#git-issue-update-parameters-type-definitions-IssueStateMap}

- `Enum<String>`

open | close

###### IssueUpdateLabel {#git-issue-update-parameters-type-definitions-IssueUpdateLabel}

- add
  - type: `Array<String>` | `String`
  - required: `false`
  - 要添加的标签列表，标签不存在时，会自动创建。

- remove
  - type: `Array<String>` | `String`
  - required: `false`
  - 要移除的标签列表

###### IssueUpdateStatus {#git-issue-update-parameters-type-definitions-IssueUpdateStatus}

- label
  - type: `Array<String>` | `String`
  - required: `false`
  - 标签，多个值之间是 `or` 关系

#### 输出结果 {#git-issue-update-output}

```txt
{
    issues // issue 列表
}
```

#### 配置样例 {#git-issue-update-configuration-examples}

- 合并到 main 后，更新标签

```yaml title=".cnb.yml"
main:
  push:
    - stages:
        - name: update issue
          type: git:issue-update
          options:
            # 移除 “开发中” 标签，添加 “预发布” 标签
            label:
              add: 预发布
              remove: 开发中
            # 当有 “feature” 或 “bug” 标签才进行上述标签操作
            when:
              label:
                - feature
                - bug
```

- tag_push 时，关闭 issue，更新标签

```yaml title=".cnb.yml"
$:
  tag_push:
    - stages:
        - name: 发布操作
          script: echo "可用发布任务替代当前任务"
        # 发布操作后执行 issue 更新操作
        - name: update issue
          type: git:issue-update
          options:
            # 关闭 issue
            state: close
            # 移除 “预发布” 标签，添加 “已发布” 标签
            label:
              add: 已发布
              remove: 预发布
            # 当有 “feature” 或 “bug” 标签才进行上述操作
            when:
              label:
                - feature
                - bug
```

- 根据 changelog 添加标签

```yaml title=".cnb.yml"
$:
  tag_push:
    - stages:
        - name: changelog
          image: cnbcool/changelog
          exports:
            latestChangeLog: LATEST_CHANGE_LOG
        - name: update issue
          type: git:issue-update
          options:
            fromText: ${LATEST_CHANGE_LOG}
            label:
              add: 需求已接收
            when:
              label: feature
```

### reviewer {#git-reviewer}

`git:reviewer`

==配置评审人或处理人==给 PR 添加、删除评审人/处理人，可指定备选评审人/处理人范围。

- [适用事件](./#git-reviewer-applicable-events)
- [参数](./#git-reviewer-parameters)
- [输出结果](./#git-reviewer-output)
- [配置样例](./#git-reviewer-configuration-examples)
- [最佳实践](./#git-reviewer-best-practices)

#### 适用事件 {#git-reviewer-applicable-events}

- `pull_request`
- `pull_request.target`
- `pull_request.update`

#### 参数 {#git-reviewer-parameters}

- [type](./#git-reviewer-parameters-type)
- [reviewers](./#git-reviewer-parameters-reviewers)
- [count](./#git-reviewer-parameters-count)
- [exclude](./#git-reviewer-parameters-exclude)
- [reviewersConfig](./#git-reviewer-parameters-reviewersConfig)
- [role](./#git-reviewer-parameters-role)

##### type {#git-reviewer-parameters-type}

- type: [REVIEW_OPERATION_TYPE](./#git-reviewer-parameters-type-definitions-REVIEW_OPERATION_TYPE)
- required: `false`
- default: `add-reviewer`

操作类型：

- `add-reviewer`: 添加评审人，会从 `reviewers` 参数中选择指定数量的评审人
- `add-reviewer-from-repo-members`: 从仓库直接成员里选一名，添加为评审人
- `add-reviewer-from-group-members`: 从仓库父组织（直接上级组织）里选一名，添加为评审人
- `add-reviewer-from-outside-collaborator-members`: 从仓库的外部协作者里选一名，添加为评审人
- `remove-reviewer`: 从已有的评审人中删除指定的成员
- `add-assignee`: 添加处理人，添加 `reviewers` 参数传入的成员
- `remove-assignee`: 从已有的处理人中删除指定的成员

##### reviewers {#git-reviewer-parameters-reviewers}

- type: `Array<String>` | `String`
- required: `false`

要 `添加` 或 `删除` 的 `reviewer` 用户名。多个使用 `,` 或 `;` 分隔。

`type` 为 `add-reviewer`、`remove-reviewer`、`add-assignee`、`remove-assignee` 时有效。

若同时配置了 `reviewers` 、 `reviewersConfig` ，则会在两者中随机选指定数量的评审人或处理人

##### count {#git-reviewer-parameters-count}

- type: `Number`
- required: `false`

指定要添加的评审人或处理人数量，随机抽取指定数量的评审人或处理人。

- 当 `type=add-reviewer`，count 缺省值为 `reviewers` 的数量，即全部添加
- 当 `type=add-reviewer-from-repo-members`，count 缺省值为 1
- 当 `type=add-reviewer-from-group-members`，count 缺省值为 1
- 当 `type=add-reviewer-from-outside-collaborator-members`，count 缺省值为 1
- 当 `type=add-assignee`，count 缺省值为 `reviewers` 的数量，即全部添加

如果已有评审人或处理人的数量 `< count` ，那么补齐。

如果已有评审人或处理人的数量 `>= count` ，那么什么也不做。

##### exclude {#git-reviewer-parameters-exclude}

- type: `Array<String>` | `String`
- required: `false`

排除指定的用户。

##### reviewersConfig {#git-reviewer-parameters-reviewersConfig}

- type: [IReviewersConfig](./#git-reviewer-parameters-type-definitions-IReviewersConfig) | `String`
- required: `false`

文件评审人或处理人配置，如果配置中有配置当前变更文件的评审人或处理人，则会被添加为评审人或处理人。

`type` 为 `add-reviewer` 或  `add-assignee` 时有效。

支持两种格式：

- 字符串格式，传入配置文件相对路径(支持 `json` 文件)：

```txt
reviewersConfig: config.json
```

```txt
{
   "./src": "name1,name2",
   ".cnb.yml": "name3"
}
```

- 对象格式，传入文件评审人配置：

```txt
reviewersConfig:
  ./src: name1,name2
  .cnb.yml: name3
```

其中，文件评审人或处理人配置的 `key` 为文件相对路径；`value` 为评审人或处理人英文名，用英文逗号分隔。

若同时配置了 `reviewers` 、`reviewersConfig`，则会在两者中随机选指定数量的评审人或处理人。

##### role {#git-reviewer-parameters-role}

- type: [ROLE_TYPE](./#git-reviewer-parameters-type-definitions-ROLE_TYPE)
- required: `false`

评审人或处理人可以添加的角色可选包括：`Developer`、`Master`、`Owner` 如果选择 `Developer`，则可添加 `Developer` 及以上权限成员，包括 `Developer`、`Master`、`Owner`

##### 类型定义 {#git-reviewer-parameters-type-definitions}

###### REVIEW_OPERATION_TYPE {#git-reviewer-parameters-type-definitions-REVIEW_OPERATION_TYPE}

- `Enum<String>`

add-reviewer | add-reviewer-from-repo-members | add-reviewer-from-group-members | add-reviewer-from-outside-collaborator-members | remove-reviewer | add-assignee | remove-assignee

###### IReviewersConfig {#git-reviewer-parameters-type-definitions-IReviewersConfig}

```txt
{
    [key: String]: String
}
```

###### ROLE_TYPE {#git-reviewer-parameters-type-definitions-ROLE_TYPE}

- `Enum<String>`

Developer | Master | Owner

#### 输出结果 {#git-reviewer-output}

添加评审人的输出结果：

```json
{
  // 当前有效的评审人
  "reviewers": [],

  // reviewers 对应的 at 消息格式，方便发送通知
  "reviewersForAt": []
}
```

添加处理人的输出结果：

```json
{
  // 当前有效的处理人
  "assignees": [],

  // assignees 对应的 at 消息格式，方便发送通知
  "assigneesForAt": []
}
```

#### 配置样例 {#git-reviewer-configuration-examples}

- 添加评审人和处理人

```yaml title=".cnb.yml"
main:
  pull_request:
    - stages:
        - name: 添加评审人
          type: git:reviewer
          options:
            type: add-reviewer
            reviewers: aaa;bbb

        - name: 添加处理人
          type: git:reviewer
          options:
            type: add-assignee
            reviewers: ccc;ddd
```

- 删除指定成员的评审人和处理人

```yaml title=".cnb.yml"
main:
  pull_request:
    - stages:
        - name: 删除评审人
          type: git:reviewer
          options:
            type: remove-reviewer
            reviewers: aaa

        - name: 删除处理人
          type: git:reviewer
          options:
            type: remove-assignee
            reviewers: aaa
```

- 结合 ifModify，指定文件被修改时添加评审人和处理人

```yaml title=".cnb.yml"
main:
  pull_request:
    - stages:
        - name: 改配置？需要特别 review
          ifModify:
            - ".cnb.yml"
            - "configs/**"
          type: git:reviewer
          options:
            type: add-reviewer
            reviewers: bbb

        - name: 改配置？需要特别处理
          ifModify:
            - ".cnb.yml"
            - "configs/**"
          type: git:reviewer
          options:
            type: add-assignee
            reviewers: ccc
```

- 结合 if，在指定条件下将某些人添加为评审人或处理人

```yaml title=".cnb.yml"
main:
  pull_request:
    - stages:
        - name: 下班时间发版本？需指定人评审。
          if: |
            [ $(date +%H) -ge 18 ]
          type: git:reviewer
          options:
            type: add-reviewer
            reviewers: bbb

        - name: 下班时间发版本？需指定人处理。
          if: |
            [ $(date +%H) -ge 18 ]
          type: git:reviewer
          options:
            type: add-assignee
            reviewers: ccc
```

- 随机选择一名负责人走查代码

```yaml title=".cnb.yml"
main:
  pull_request:
    - stages:
        - name: random
          image: tencentcom/random
          settings:
            from:
              - aaa
              - bbb
          exports:
            result: CURR_REVIEWER
        - name: show  CURR_REVIEWER
          script: echo ${CURR_REVIEWER}
        - name: add reviewer
          type: git:reviewer
          options:
            type: add-reviewer
            reviewers: ${CURR_REVIEWER}
```

- 从当前仓库成员里选一名评审人进行 review

```yaml title=".cnb.yml"
main:
  pull_request:
    - stages:
        - name: add reviewer
          type: git:reviewer
          options:
            type: add-reviewer-from-repo-members
```

#### 最佳实践 {#git-reviewer-best-practices}

建立专用的 CR 群，交叉评审自动合并

1. `PR` 随机选择N名评审者，通知到群
2. 评审通过后自动合并
3. 记录评审者信息在提交信息中
4. `Git` 设置满足多人评审自动合并策略实现多人交叉评审

```yaml title=".cnb.yml"
main:
  review:
    - stages:
        - name: CR 通过后自动合并
          type: git:auto-merge
          options:
            mergeType: squash
            removeSourceBranch: true
            mergeCommitMessage: $CNB_LATEST_COMMIT_MESSAGE
          exports:
            reviewedBy: REVIEWED_BY
        # 将评审消息发送到企业微信机器人
        - name: notify
          image: tencentcom/wecom-message
          settings:
            msgType: markdown
            robot: "155af237-6041-4125-9340-000000000000"
            content: |
              > CR 通过后自动合并 <@${CNB_BUILD_USER}> 
              > 　
              > ${CNB_PULL_REQUEST_TITLE}
              > [${CNB_EVENT_URL}](${CNB_EVENT_URL})
              > 
              > ${REVIEWED_BY}

  pull_request:
    - stages:
        # ...省略其它任务

        # 发送到 CR 专用群
        - name: add reviewer
          type: git:reviewer
          options:
            reviewers: aaa,bbb,ccc,ddd
            count: 2
          exports:
            reviewersForAt: CURR_REVIEWER_FOR_AT
        - name: notify
          image: tencentcom/wecom-message
          settings:
            msgType: markdown
            robot: "155af237-6041-4125-9340-000000000000"
            message: |
              > ${CURR_REVIEWER_FOR_AT}
              > 　
              > ${CNB_PULL_REQUEST_TITLE}
              > [${CNB_EVENT_URL}](${CNB_EVENT_URL})
              > 　
              > from ${CNB_BUILD_USER}
```

### release {#git-release}

`git:release`

==仓库发布 Release==

- [适用事件](./#git-release-applicable-events)
- [参数](./#git-release-parameters)
- [输出结果](./#git-release-output)
- [配置样例](./#git-release-configuration-examples)

#### 适用事件 {#git-release-applicable-events}

- `push`
- `commit.add`
- `branch.create`
- `tag_push`
- `pull_request.merged`
- `api_trigger`
- `web_trigger`
- `tag_deploy`

#### 参数 {#git-release-parameters}

- [overlying](./#git-release-parameters-overlying)
- [tag](./#git-release-parameters-tag)
- [title](./#git-release-parameters-title)
- [description](./#git-release-parameters-description)
- [preRelease](./#git-release-parameters-preRelease)
- [latest](./#git-release-parameters-latest)

##### overlying {#git-release-parameters-overlying}

- type: `Boolean`
- required: `false`
- default: `false`

叠加模式

- `true`：叠加模式，同一个 `release` 可以提交多次，最终的 `release` 是多次提交的并集。如果同名附件已经存在，会先删除再上传，达到更新的效果。
- `false`：非叠加模式，以最后一个提交的为准，前面的会被清除。

::: warning
默认情况下，当 release 版本已经存在时，会先删除这个版本，重新生成。
:::

##### tag {#git-release-parameters-tag}

- type: `String`
- required: `false`

`release` 对应的 `tag` 名，非必填

`tag_push` 事件无需传入，直接取触发 `tag_push` 事件的 `tag` 名。

非 `tag_push` 事件必填，用来作为 `release` 对应的 `tag` 名。

##### title {#git-release-parameters-title}

- type: `String`
- required: `false`
- default: `tag` 名

release 的标题

##### description {#git-release-parameters-description}

- type: `String`
- required: `false`

`release` 的描述

##### preRelease {#git-release-parameters-preRelease}

- type: `Boolean`
- required: `false`
- default: `false`

是否将 `release` 设置为 预发布

##### latest {#git-release-parameters-latest}

- type: "true" | "false" | true | false
- required: `false`
- default: `false`

是否将 `release` 设置为最新版本

#### 输出结果 {#git-release-output}

无

#### 配置样例 {#git-release-configuration-examples}

- 生成 changelog 并自动更新 release 描述

```yaml title=".cnb.yml"
$:
  tag_push:
    - stages:
        - name: changelog
          image: cnbcool/changelog
          exports:
            latestChangeLog: LATEST_CHANGE_LOG
        - name: upload release
          type: git:release
          options:
            title: release
            description: ${LATEST_CHANGE_LOG}
```

- 主分支 push 时发布 release

```yaml title=".cnb.yml"
main:
  push:
    - stages:
        - name: git release
          type: git:release
          options:
            tag: Nightly
            description: description
```

## testing

### coverage {#testing-coverage}

`testing:coverage`

==单测覆盖率==，通过单测结果报告，计算单测覆盖率结果，并上报徽章。

- [适用事件](./#testing-coverage-applicable-events)
- [全量覆盖率](./#testing-coverage-full-coverage)
- [增量覆盖率](./#testing-coverage-incremental-coverage)
- [参数](./#testing-coverage-parameters)
- [输出结果](./#testing-coverage-output)
- [配置样例](./#testing-coverage-configuration-examples)

#### 适用事件 {#testing-coverage-applicable-events}

[所有事件](../trigger-rule.md#trigger-event)

#### 全量覆盖率 {#testing-coverage-full-coverage}

![coverage](https://cnb.cool/svg/badge/coverage?message=18.67%25&color=l2)
![coverage](https://cnb.cool/svg/badge/coverage?message=38.67%25&color=l3)
![coverage](https://cnb.cool/svg/badge/coverage?message=58.67%25&color=l4)
![coverage](https://cnb.cool/svg/badge/coverage?message=78.67%25&color=l5)

从本地覆盖率报告文件解析而来，目前可识别以下格式的报告：

- `json(js)`（推荐使用）
- `json-summary`
- `lcov`（推荐使用）
- `jacoco`
- `golang`

#### 增量覆盖率 {#testing-coverage-incremental-coverage}

![coverage-pr](https://cnb.cool/svg/badge/coverage%20%40pr?message=18.67%25&color=l2)
![coverage-pr](https://cnb.cool/svg/badge/coverage%20%40pr?message=38.67%25&color=l3)
![coverage-pr](https://cnb.cool/svg/badge/coverage%20%40pr?message=58.67%25&color=l4)
![coverage-pr](https://cnb.cool/svg/badge/coverage%20%40pr?message=78.67%25&color=l5)

变更行对应的单测覆盖率，仅支持 `pull_request`、`pull_request.update`、`pull_request.target` 事件

#### 参数 {#testing-coverage-parameters}

- [pattern](./#testing-coverage-parameters-pattern)
- [lines](./#testing-coverage-parameters-lines)
- [diffLines](./#testing-coverage-parameters-diffLines)
- [allowExts](./#testing-coverage-parameters-allowExts)
- [lang](./#testing-coverage-parameters-lang)
- [breakIfNoCoverage](./#testing-coverage-parameters-breakIfNoCoverage)

##### pattern {#testing-coverage-parameters-pattern}

- type: `String`
- required: `false`

Glob 格式，指定覆盖率报告文件位置，相对于当前工作目录。 缺省时，将尝试查找当前目录（包括子目录）下的以下文件：coverage.json、jacoco*.xml、lcov.info、*.lcov。

##### lines {#testing-coverage-parameters-lines}

- type: `Number`
- required: `false`

指定全量覆盖率红线，判断如果全量覆盖率百分比小于该值，阻断工作流退出流水线。

##### diffLines {#testing-coverage-parameters-diffLines}

- type: `Number`
- required: `false`

指定增量覆盖率红线，判断如果全量覆盖率百分比小于该值，阻断工作流退出流水线。 `pull_request`、`pull_request.update`、`pull_request.target` 事件支持计算增量覆盖率结果，其他事件只计算全量覆盖率。

##### allowExts {#testing-coverage-parameters-allowExts}

- type: `String`
- required: `false`

参与覆盖率计算的代码文件类型白名单，逗号分隔, 如：`.json,.ts,.js`。 缺省时报告中的文件都会参与计算。

##### lang {#testing-coverage-parameters-lang}

- type: `String`
- required: `false`

当覆盖率结果报告目标格式为 `golang` 时，请指定此参数 为 `go`，否则会出现计算误差。其他情况可忽略该参数。

##### breakIfNoCoverage {#testing-coverage-parameters-breakIfNoCoverage}

- type: `Boolean`
- required: `false`

没有找到覆盖率报告文件时，是否抛出错误终止流程。

#### 输出结果 {#testing-coverage-output}

```json
{
  // 代码行覆盖率，例如 100，计算出错时值为 NA
  "lines": 100,
  // 代码增量行覆盖率，例如 100，计算出错时值为 NA
  "diff_pct": 100
}
```

#### 配置样例 {#testing-coverage-configuration-examples}

```yaml title=".cnb.yml"
main:
  push:
    - stages:
        - name: coverage
          type: testing:coverage
          options:
            breakIfNoCoverage: false
          exports:
            lines: LINES
        - name: result
          script: echo $LINES
```

## artifact

### remove-tag {#artifact-remove-tag}

`artifact:remove-tag`

==删除 CNB 制品标签==，目前仅删除支持 CNB docker 和 helm 标签。需要有仓库写权限。

- [适用事件](./#artifact-remove-tag-applicable-events)
- [参数](./#artifact-remove-tag-parameters)
- [配置样例](./#artifact-remove-tag-configuration-examples)

#### 适用事件 {#artifact-remove-tag-applicable-events}

- `push`
- `commit.add`
- `tag_push`
- `tag_deploy`
- `pull_request.merged`
- `api_trigger`
- `web_trigger`
- `crontab`
- `branch.create`

#### 参数 {#artifact-remove-tag-parameters}

- [name](./#artifact-remove-tag-parameters-name)
- [tags](./#artifact-remove-tag-parameters-tags)
- [type](./#artifact-remove-tag-parameters-type)

##### name {#artifact-remove-tag-parameters-name}

- type: `String`
- required: `true`

制品包名

##### tags {#artifact-remove-tag-parameters-tags}

- type: `Array<string>`
- required: true

##### type {#artifact-remove-tag-parameters-type}

- type: `String`
- required: `false`
- default: `docker`

制品类型，目前仅支持 docker 和 helm

#### 配置样例 {#artifact-remove-tag-configuration-examples}

```yaml title=".cnb.yml"
main:
  push:
    - stages:
        - name: remove tag
          type: artifact:remove-tag
          options:
            # 包名
            # 包名示例1，仓库同名制品：reponame
            # 包名示例2，仓库非同名制品：reponame/name
            name: reponame/name
            tags:
              - tag1
              - tag2
            type: docker
````

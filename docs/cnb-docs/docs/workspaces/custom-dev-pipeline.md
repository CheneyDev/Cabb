---
title: 自定义环境创建流程
permalink: https://docs.cnb.cool/zh/workspaces/custom-dev-pipeline.html
summary: 在仓库根目录下新增 `.cnb.yml` 文件可自定义云原生开发的流水线配置，包括声明资源规格、触发事件及可用时机，并可通过内置任务 `vscode:go` 延迟进入开发环境，还可用 `endStages` 定义销毁环境前的任务 。
---

如果希望点击分支页面 `云原生开发` 能使用自定义的流水线配置, 可在仓库根目录下新增 `.cnb.yml` 文件, 文件内新增如下配置:

```yaml
# .cnb.yml
$:
  # vscode 事件：专供页面中启动远程开发用
  vscode:
    - docker:
        # 自定义镜像作为开发环境
        image: node:20
      services:
        - vscode
        - docker
      stages:
        - name: ls
          script: ls -al
```

## 自定义资源规格

可以通过 `runner.cpus` 声明需要的开发资源，最大支持 `64核`，内存为 cpus x 2(GB)。

```yaml{4}
# .cnb.yml
$:
  vscode:
    - runner:
        cpus: 64
      docker:
        build: .ide/Dockerfile
      services:
        - vscode
        - docker
      stages:
        - name: ls
          script: ls -al
```

## 自定义启动事件

通过 `.cnb.yml`，声明指定事件触发时自动创建开发环境。触发事件推荐使用：

- [`vscode`](../build/grammar.md#vscode)：仓库页面点击`启动云原生开发`按钮时创建开发环境
- [`branch.create`](../build/grammar.md#branch-create)：创建分支时创建开发环境
- [`api_trigger`](../build/grammar.md#api_trigger)：自定义事件触发创建开发环境
- [`web_trigger`](../build/grammar.md#web_trigger)：web 页面自定义事件触发创建开发环境

```yaml{5}
# .cnb.yml
# 匹配所有分支
(**):
  # 创建分支时创建开发环境
  branch.create:
    - name: vscode
      services:
        # 声明使用 vscode 服务
        - vscode
      docker:
        # 自定义开发环境
        build: .ide/Dockerfile
      stages:
        - name: 执行自定义脚本
          script:
            - npm install
            - npm run start
```

## 自定义可用时机

云原生开发可用时机默认为：流水线准备（`prepare`）阶段执行完（`code-server` 代码服务在准备阶段启动），`stages` 任务执行前。

如果希望执行某些任务后再进入开发环境，即延迟进入开发环境时机，可使用 [vscode:go](../build/internal-steps/README.md#go) 内置任务。
使用该任务，启动云原生开发后，`loading` 页将延迟进入云原生开发入口选择页。当 `vscode:go` 任务执行后才能进入入口选择页。

注意，使用 `vscode:go` 任务将增加等待时间。
可将必须在进入开发环境前执行的任务放在 `vscode:go` 前执行，
在进入开发环境后执行的任务放在 `vscode:go` 后。
**如果没有必须在进入开发环境前执行的任务，就无需使用 `vscode:go`**。

当 `stages` 任务执行失败，远程开发是否结束：

- 使用 `vscode:go`: `vscode:go` 前的任务执行失败，开发环境将销毁
- 使用 `vscode:go`：`vscode:go` 后的任务执行失败，开发环境不会销毁
- 未使用 `vscode:go`：`stages` 任务执行失败，开发环境不会销毁

## 自定义环境销毁前任务

可使用 [`endStages`](../build/grammar.md#endstages) 定义开发环境销毁前需要执行的任务.

```yaml{14-18}
# .cnb.yml
$:
  vscode:
    - docker:
        image: node:20
      services:
        - vscode
        - docker
      # 开发环境启动后会执行的任务
      stages:
        - name: ls
          script: ls -al
      # 开发环境销毁前会执行该任务
      endStages:
        - name: end stage 1
          script: echo "end stage 1"
        - name: end stage 2
          script: echo "end stage 2"
```

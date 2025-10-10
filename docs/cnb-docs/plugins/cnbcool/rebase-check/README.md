# rebase-check

检测当前分支是否已合并指定分支。

## 参数说明

检测当前分支(source/submoduleSource)是否已合并过指定分支(target/submoduleTarget)。

* source: string, 非必填, 源分支, 默认为当前分支
* target: string, 非必填, 目标分支, 默认为 main
* submoduleSource: string, 非必填, 子模块源分支, 默认为当前分支
* submoduleTarget: string, 非必填, 子模块目标分支, 默认为 main

## 在 云原生构建 中使用

```yaml
main:
  push:
  - stages:
    - name: rebase-check
      image: cnbcool/rebase-check:latest
      settings:
        # 检测 release 分支是否已合并过 main 分支
        # 即 release 分支是否有同步过 main 分支的修改
        source: release
        target: main 
        submoduleSource: release
        submoduleTarget: main
```

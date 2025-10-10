# pbLint Docker Image

本项目主要提供了 pblint 的 Docker 镜像，主要目标是在 云原生构建 中通过 pblint 对 protobuf 文件进行检查。

## 输入

### patterns

* Type：`String | String[]`

该参数指定需要校验的 protobuf 文件，支持高级通配符。
默认值是当前目录下的所有文件。

## 在 云原生构建 上使用

```yml
main:
  pull_request:
  - stages:
    - name: pblint job
      image: tencentcom/pb-lint:latest
      settings:
        patterns: 
          - "./pb/**/*.proto"
          - "./pb2/**/*.proto"
```

# Docker lint

一个更智能的 Dockerfile linter，帮助您构建最佳实践 Docker 镜像。

## 在 Docker 上使用

```shell
docker run --rm -it -v $(pwd):$(pwd) -w $(pwd) hadolint/hadolint hadolint -h
```

## 在 云原生构建 上使用

```yaml
main:
  pull_request:
  - stages:
    - name: docker-lint
      image: hadolint/hadolint
      commands: |
        hadolint -h
        hadolint Dockerfile
```

更多配置项请查阅：[hadolint](https://github.com/hadolint/hadolint)

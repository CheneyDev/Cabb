# markdown-lint

Command Line Interface for [markdownlint](https://github.com/DavidAnson/markdownlint).

## 在 Docker 上使用

```shell
docker run --rm -v $(pwd):$(pwd) -w $(pwd) tencentcom/markdown-lint
```

## 在 云原生构建 上使用

```yaml
main:
  pull_request:
  - stages:
    - name: tencentcom/markdown-lint
      # 默认 **/*.md
      image: tencentcom/markdown-lint
```

```yaml
master:
  pull_request:
  - stages:
    - name: tencentcom/markdown-lint
      image: tencentcom/markdown-lint
      commands:
      - markdownlint --help
      - markdownlint '**/*.md' --ignore node_modules --disable MD013
```

## 参考资料

更多信息，请查阅：[igorshubovych/markdownlint-cli][more]

[more]:https://github.com/igorshubovych/markdownlint-cli

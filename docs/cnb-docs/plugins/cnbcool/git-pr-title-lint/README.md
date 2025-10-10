# git-pr-title-lint

基于约定式提交规范，检查 PR 标题是否符合给定的规范

## 约定式提交

使用插件前请现在仓库根目录下配置检查规则：[`commitlint.config.js`](https://commitlint.js.org/#/reference-configuration)

可参考开源社区规范，详见 [https://www.conventionalcommits.org/zh/](https://www.conventionalcommits.org/)

## 在云原生构建中使用

```yaml
# 检查 PR 的 title，需要符合 commit mssage 规范
main:
  pull_request:
  - stages:
    - name: git-pr-title-lint
      image: cnbcool/git-pr-title-lint
      settings:
        # title 可选参数，待检查的 message。如果不传入则默认取环境变量 CNB_PULL_REQUEST_TITLE
        title: 'fix: xxx'
```

## 在 Docker 中使用

```shell
docker run --rm \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    -e CNB_PULL_REQUEST_TITLE="feat: add demo for plugins" \
    cnbcool/git-pr-title-lint
```

# git-mr-title-lint

基于约定式提交规范，检查 PR 标题是否符合给定的规范。

## 约定式提交

与开源社区规范一致，详见 [https://www.conventionalcommits.org/zh/](https://www.conventionalcommits.org/zh/)

## 配置文件

仓库根目录需要存在配置文件，可将配置文件示例[commitlintrc.yml](./commitlintrc.yml)复制过来。

## 在 Docker 中使用

```shell
docker run --rm \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    -e CNB_PULL_REQUEST_TITLE="feat: add demo for plugins" \
    tencentcom/git-mr-title-lint
```

## 在 云原生构建 中使用

```yaml
# 检查 MR 的 title，需要符合 commit mssage 规范
main:
  pull_request:
  - stages:
    - name: git-mr-title-lint
      image: tencentcom/git-mr-title-lint
```

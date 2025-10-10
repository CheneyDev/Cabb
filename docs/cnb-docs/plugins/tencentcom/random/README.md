# random

根据选项，随机选择

## 输入

- `from`: `String[]` 必填，随机列表
- `exclude`: `String|String[]` 可选，需要排除的值，指定为`$CNB_BUILD_USER`可用于排除当前构建人
- `by`: `Number` 可选，随机向量，不填时随机生成一个
- `count`: `Number` 可选，随机总数，默认为 1

## 输出

- `result`: `String` 命中的选择

## 在 云原生构建 上使用

```yaml
main:
  pull_request:
  - stages:
    - name: random
      image: tencentcom/random
      settings:
        from: 
          - user1
          - user2
      envExport:
        result: CURR_REVIEWER
    - name: show  CURR_REVIEWER
      script: echo ${CURR_REVIEWER}
```

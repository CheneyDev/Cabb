# CoverageCheck

CoverageCheck Docker 插件

## 输入

- `summary`: `string` 必选，从指定文件读取单元测试coverage-summary.json报告
- `check_files`: `string` 可选，从指定文件读取需要检查的文件列表，为空跳过检查
- `lines`: `number` 可选，要求最低`Lines`覆盖率，默认为0
- `statements`: `number` 可选，要求最低`Statements`覆盖率，默认为0
- `functions`: `number` 可选，要求最低`Functions`覆盖率，默认为0
- `branches`: `number` 可选，要求最低`Branches`覆盖率，默认为0

## 示例

在 云原生构建 上使用

```yaml
# .cnb.yml
main:
  pull_request:
  - stages:
    - name: run unit test
      script: npm run test
    - name: git-change-list
      image: tencentcom/git-change-list:latest
      settings:
        changed: changed.txt
    - name: do check coverage
      image: tencentcom/coverage-check:latest
      settings:
        summary: ./coverage/coverage-summary.json
        check_files: changed.txt
        lines: 80
        statements: 80
        functions: 80
        branches: 80
        final: ./coverage/coverage-final.json
```

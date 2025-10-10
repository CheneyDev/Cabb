# CyclomaticComplexityCheck

圈复杂度计算并检查工具，可输出千行圈复杂度超标率等指标。

本插件对 [Lizard][Lizard] 进行了封装，便于大家在 CI 中直接使用。

[Lizard]:https://github.com/terryyin/lizard

CyclomaticComplexityCheck 支持两种工作模式：

- 计算 C/TLoC 并进行阈值检查

例如你可以限制只有 C25/TLoC 的值低于 3 才可以检查通过。
此时插件会检查整个项目（exclude 文件夹中的文件除外），根据千行圈复杂度25平均超标数进行计算，
如果结果大于3，检查将不会通过。

- 计算指定文件的圈复杂度

如果有圈复杂度超过指定值，检查将不会通过。

## 输入

- `ccn`: `number` 可选，圈复杂度标准超标下限，默认值为25
- `threshold`: `number` 可选，圈复杂度千行平均超标数阈值，默认为3
- `total_exceed_ccn`: `number` 可选，圈复杂度超标总数上限，不设默认值，可为空
- `exclude`: `string` 可选，需要排除目录的`yml`配置地址，详情见下文
- `language`: `string` 可选，需要分析的语言，目前支持的值包括`cpp`、`csharp`、`java`、`javascript`、`objectivec`、`php`、`python`、`ruby`、`swift`、`ttcn`，一次只能传一种语言，默认值为自动检测
- `check_files`: `string` 可选，从指定文件读取需要检查的文件列表，为空则检查整个项目

> 注意：对于前端而言，目前仅支持javascript和jsx。不支持typescript或vue格式的文件。

**`threshold` 和 `check_files` 参数不会同时起作用。**

- 如果你希望检查千行圈复杂度超标，可以传的参数有 `ccn`、`threshold`、`exclude`和`language`。

- 如果你希望检查指定文件的圈复杂度超标，可以传的参数有 `ccn`、`check_files`、`exclude`和`language`。

## 在 云原生构建 上使用

将 C25/TLoC（圈复杂度25千行平均超标数） 的阈值限定为 3，那么需要 `ccn` 的值设置为 25，同时将 `threshold` 的值设置为 3。
这表示“圈复杂度超过25的所有复杂度之和 ÷ 该语言本次被扫描的生产代码总行数 * 1000”的值如果超过 3，检查将会失败。

对应的配置如下：

```yaml
# .cnb.yml
main:
  pull_request:
  - stages:
    - name: check C/TLoC
      image: tencentcom/cyclomatic-complexity:latest
      settings:
        ccn: 25
        threshold: 3
```

在 Merge Request 时检查变更文件的圈复杂度，修改的文件圈复杂度不大于 20，新增的文件圈复杂度不大于 5。

对应的配置如下：

```yaml
# .cnb.yml
main:
  pull_request:
  - stages:
    - name: list changed
      type: git:changeList
      options:
        changed: changed.txt
        added: added.txt
    - name: check added files CCN
      image: tencentcom/cyclomatic-complexity:latest
      settings:
        ccn: 5
        check_files: added.txt
    - name: check changed files CCN
      image: tencentcom/cyclomatic-complexity:latest
      settings:
        ccn: 25
        check_files: changed.txt
            
```

如果代码仓库中有一些需要排除的代码，可以传入一个 `yml` 配置的地址，例如在代码仓库的根目录中存放一个 `.lizard.yml` 配置文件。

配置文件的内容如下：

```yaml
exclude:
  - ./node_modules/*
  - ./dist/*
  - ./lib/externals/*
```

配置文件的每一行可以配置一个需要排除的目录，目录支持 `*` 和 `?` 通配符。

同时，在 Coding-CI 的配置文件中增加 `exclude` 参数。

```yaml
# .cnb.yml
main:
  pull_request:
  - stages:
    - name: check C/TLoC
      image: tencentcom/cyclomatic-complexity:latest
      settings:
        ccn: 25
        threshold: 3
        exclude: .lizard.yml
```

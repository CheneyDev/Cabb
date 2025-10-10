# 介绍

获取从最后一个 `TAG` 到当前的提交信息

## 参数

### toFile

* type: string
* 必填: 否
* 默认值: `commits_from_tag.json`

输出到指定文件，`JSON` 格式

## 结果

写入结构为一个数组，每个元素的结果如下，[详细定义](https://github.com/DefinitelyTyped/DefinitelyTyped/blob/master/types/conventional-commits-parser/index.d.ts)：

```json
{
    "type": "feat",
    "scope": "scope",
    "subject": "broadcast $destroy event on scope destruction",
    "merge": null,
    "header": "feat(scope): broadcast $destroy event on scope destruction",
    "body": null,
    "footer": "Closes #1",
    "notes": [],
    "references": [
        {
            "action": "Closes",
            "owner": null,
            "repository": null,
            "issue": "1",
            "raw": "#1",
            "prefix": "#"
        }
    ],
    "mentions": [],
    "revert": null
}
```

## 示例

```yaml
test:
  push:
    - stages:
        - name: commits from tag
          image: tencentcom/git-commits-from-tag:latest
          settings:
            toFile: commits.json
        - name: print
          script: cat commits.json
```

```shell
docker run --rm \
    -e TZ=Asia/Shanghai \
    -e PLUGIN_TOFILE="commits.json" \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    tencentcom/git-commits-from-tag:latest
```

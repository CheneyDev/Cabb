# fetch-git-file

拉取 cnb git 文件到本地

## 参数

### token

- type: String
- required: 否
- default: `CNB_TOKEN`

调 cnb open api 所需的 token

需要有文件管理读权限 `repo-contents:r`

默认为环境变量中的 `CNB_TOKEN`

### slug

- type: String
- required: 是

仓库路径

### ref

- type: String
- required: 是

仓库分支、tag、sha

### files

- type: String
- required: 是

需要拉取的文件列表

多行文本，一行一个路径，如：

```shell
a/b/c.txt
d/e.ts
f.yml
```

### target

- type: String
- required: 是
- defult: `_tmp_`

文件存放目录

仓库根目录可填 `.`

## 在 云原生构建 上使用

```yaml
main:
  push:
    - stages:
        - name: fetch-git-file
          image: cnbcool/fetch-git-file
          settings:
            slug: xx/xx
            ref: master
            files: |
              a/b/c.txt
              d/e.ts
              f.yml
```

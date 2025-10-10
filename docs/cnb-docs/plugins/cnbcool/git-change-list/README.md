# git-change-list

获取变更文件列表，支持输出变更文件列表到文件

## 参数

### changed

- type: String
- required: 否

变更列表存放文件名，包含修改过的、新增加的文件。

### deleted

- type: String
- required: 否

删除列表存放文件名。

### added

- type: String
- required: 否

新增列表存放文件名。

### edited

- type: String
- required: 否

修改列表存放文件名。

## 输出

```js
{
  // 变更列表存放文件名，包含修改过的、新增加的文件。
  changed,

  // 删除列表存放文件名。
  deleted,

  // 新增列表存放文件名。
  added,

  // 修改列表存放文件名。
  edited,
}
```

## 在 云原生构建 中使用

```yaml
# .cnb.yml
main:
  pull_request:
    - stages:
      - name: git-change-list
        image: cnbcool/git-change-list:latest
        settings:
          # 可选，文件列表输出到文件中
          changed: changed.txt
          deleted: deleted.txt
          added: added.txt
          edited: edited.txt
      - name: 回显
        script: ls -al       
```

```yaml
# .cnb.yml
main:
  pull_request:
    - stages:
      - name: git-change-list
        image: cnbcool/git-change-list:latest
        exports:
          changed: CHANGED
      - name: 显示变更文件列表
        script: echo "$CHANGED"
```

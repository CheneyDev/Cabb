# cpplint

通过 cpplint 对 C++ 代码进行检查。

## 使用方法

### settings 参数说明

#### files

> 默认值：$(find . -type f | grep -v .git)

目前仅需要对指定的文件列表进行检查，默认值是当前目录、以及除了 .git 以外所有子目录下的所有文件。
因为 cpplint 使用 /bin/sh 执行，不支持很多高级通配符，所以这里默认直接使用 find 命令将所有文件全部遍历出来。

### 配置文件 CPPLINT.cfg

在自己的C++项目中可以添加 `CPPLINT.cfg` 配置文件，来对 cpplint 生成针对性配置。

```cfg
# 从默认的 80 改为我们的 120
linelength=120

# 支持c++11和c++14语法
filter=-build/c++11,-build/c++14,-runtime/references

# 过滤掉缓存目录中的文件
exclude_files=.*(ccls|bazel-).*
```

## 在 云原生构建 上使用

### 全量检查

```yaml
# .cnb.yml
main:
  pull_request:
    - stages:
      - name: cpplint job
        image: tencentcom/cpplint:latest
        settings:
          files: $(find . -type f | grep -v .git)
```

### 增量检查

```yaml
# .cnb.yml
main:
  pull_request:
    - stages:
      - name: make changelist
        type: git:changeList
        options:
          changed: changed.txt
      - name: cat all changed files
        script: cat changed.txt
        exports:
          stdout: ALL_CHANGED_FILES
      - name: cpplint job
        image: tencentcom/cpplint:latest
        settings:
          files: ${ALL_CHANGED_FILES}
```

## 更多用法

参考：[cpplint](https://github.com/cpplint/cpplint)

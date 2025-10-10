# release 清理插件

## 用法

每天1:30运行一次清理，只保留最近的5个release

```yaml
main:
  "crontab: 30 1 * * *":
    - name: release clean
      stages:
        - name: release clean
          image: docker.cnb.cool/cnb/plugins/cnbcool/release-clean:latest
          settings:
            filter: "RECENT_N=5"
            debug: false
```

## 参数说明

### 参数 filter

TAGNAME_PREFIX=v1.

> 删除tagname 以 "v1." 开头的release

NAME_PREFIX=v1.

> 删除name 以 "v1." 开头的release

RECENT_N=10

> 保留最近的10个release

RECENT_N_DAYS=10

> 保留最近的10天的release

RECENT_N_DAYS_RETAIN_N=10,5

> 保留最近的10天的release且至少保留5个。

### 参数 debug

true

> 仅输出log，不进行删除

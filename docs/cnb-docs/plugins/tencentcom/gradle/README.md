# gradle

- [Tencent Kona JDK](https://github.com/Tencent/TencentKona-11)
- [Gradle](https://gradle.org)

## 在 云原生构建 上使用

```yaml
main:
  push:
  - stages:
    - name: show gradle version
      image: tencentcom/gradle
      commands:
        - gradle -v
```

# maven

一个包含`Tencent Kona JDK`和 `Maven`的Docker镜像

- [Tencent Kona JDK](https://github.com/Tencent/TencentKona-11)
- [Maven](https://maven.apache.org/)

## 在 云原生构建 上使用

```yaml
main:
  push:
  - stages:
    - name: show maven version
      image: tencentcom/maven
      commands:
        - mvn -v
```

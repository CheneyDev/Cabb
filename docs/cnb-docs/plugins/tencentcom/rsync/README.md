# rsync

通过 rsync 将文件同步到远程主机，并在远程主机上执行任意命令

## 在 云原生构建 上使用

简单示例，将本地 dist 文件夹同步到远程机器的 ~/target 目录:

```yaml
# .cnb.yml
main:
  push:
    - stages:
      - name: rsync
        image: tencentcom/rsync
        # 引用密钥仓库配置文件
        imports: https://your-git.com/group/secret-repo/-/blob/main/env.yml
        settings:
          user: $RSYNC_USER
          key: $RSYNC_KEY
          hosts:
            - ip1
            - ip2
          source: ./dist/
          target: ~/target/
          include:
            - "app.tar.gz"
            - "app.tar.gz.md5"
          exclude:
            - "*"
          prescript:
            - cd ~/packages
            - md5sum -c app.tar.gz.md5
            - tar -xf app.tar.gz -C ~/app 
          script:
            - cd ~/packages
            - md5sum -c app.tar.gz.md5
            - tar -xf app.tar.gz -C ~/app
```

引用密钥仓库配置文件获取 rsync key 和 rsync user:

```yaml
# 密钥仓库 env.yml
RSYNC_KEY: xxxx
RSYNC_USER: xxx
# 声明指定镜像的插件任务能引用该配置文件
allow_images:
  - groupname/imagename
# 声明指定仓库的流水线能引用该配置文件
allow_slugs:
  - groupname/reponame
```

## 参数说明

- `user` 用于登录远程机器的用户，默认为 `root`
- `key` 用于访问远程机器的 ssh 私钥
- `hosts` 远程机器的主机名或 IP 地址
- `port` 远程机器的连接端口，默认为 `22`
- `source` 要同步的源文件夹，默认为 `./`
- `target` 要同步到远程机器上的目标文件夹
- `include` rsync 的包含过滤器
- `exclude` rsync 的排除过滤器
- `recursive` 是否递归同步，默认为 `false`
- `delete` 是否删除目标文件夹的内容，默认为 `false`
- `args` 指定插件使用这些额外的 rsync 命令行参数，例如：`"--blocking-io"`
- `prescript` 在 rsync 执行之前在远程机器上运行的命令列表
- `script` 在 rsync 执行之后在远程机器上运行的命令列表
- `log_level` SSH 日志级别，默认为安静模式（quiet）

## 来源说明

Dockerfile 中使用基础镜像来源：[drone-rsync](https://github.com/Drillster/drone-rsync)

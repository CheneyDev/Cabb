# SCP

通过 ssh 复制文件或制品到远程主机。

- 源文件列表支持通配符模式。
- 支持将文件发送到多个主机。
- 支持将文件发送到主机上的多个目标文件夹。
- 支持从绝对路径或 raw 内容加载 ssh 密钥。
- 支持 SSH ProxyCommand。

## 参数说明

- `host`: 目标主机 hostname 或 IP。支持传入多个，字符串或数组格式
- `port`: 目标主机 ssh 端口，默认为 `22`
- `protocol`: 要使用的 IP 协议：可以是 `tcp`、`tcp4` 或 `tcp6`，默认为 `tcp`
- `username`: 目标主机用户名，默认为 `root`
- `password`: 目标主机密码
- `timeout`: SSH 连接建立的超时时间，默认为 `30s`
- `key`: 私钥文本
- `key_path`: 私钥路径
- `passphrase`: 私钥密码短语
- `ciphers`: 允许的密码算法。如果未指定，则使用合理的默认值，支持传入多个，字符串或数组格式
- `use_insecure_cipher`: 在使用不安全的密码时包含更多密码，`true` 或 `false`
- `fingerprint`: 主机公钥的 SHA256 指纹，默认会跳过验证
- `command_timeout`: 执行命令的超时时间，默认为 10m
- `target`: 目标主机的文件夹路径，支持传入多个，字符串或数组格式
- `source`: 要复制的源文件列表，支持传入多个，字符串或数组格式
- `rm`: 在复制文件和制品前删除目标文件夹，`true` 或 `false`
- `proxy_host`: 代理的 hostname 或 IP
- `proxy_port`: 代理主机的 ssh 端口，默认为 `22`
- `proxy_protocol`: 用于代理的 IP 协议：可以是 `tcp`、`tcp4` 或 `tcp6`，默认为 `tcp`
- `proxy_username`: 代理主机用户名，默认为 root
- `proxy_password`: 代理主机密码
- `proxy_key`: 代理主机私钥明文文本
- `proxy_key_path`: 代理主机私钥的路径
- `proxy_passphrase`: 代理主机私钥密码短语，对私钥进行加密保护
- `proxy_timeout`: 连接代理主机的超时时间
- `proxy_ciphers`: 允许的密码算法。如果未指定，则使用合理的默认值。支持传入多个，字符串或数组格式
- `proxy_use_insecure_cipher`: 在使用不安全的密码时包含更多密码，`true` 或 `false`
- `proxy_fingerprint`: 代理主机公钥的 SHA256 指纹，默认会跳过验证
- `strip_components`: 删除指定数量的路径前缀
- `tar_exec`: 在目标主机上替代 tar 的命令，默认为 `tar`
- `tar_tmp_path`: 目标主机上 tar 文件的临时路径
- `debug`: 在复制数据之前删除目标文件夹，`true` 或 `false`
- `overwrite`: 为 tar 命令增加 `--overwrite` 标志，`true` 或 `false`
- `unlink_first`: 为 tar 命令增加 `--unlink-first` 标志，`true` 或 `false`
- `tar_dereference`:  为 tar 命令增加 `--dereference` 标志，`true` 或 `false`

## 在 云原生构建 上使用

### 简单示例

```yaml
main:
  push:
    - stages:
      - name: echo file
        image: tencentcom/scp
        settings:
          host: xx.xx.xx.xxx
          username: root
          password: xxxx
          port: 22
          target: /data/release
          source:
            - release/*.tar.gz
```

### 多台目标机器例子

```yaml
main:
  push:
    - stages:
      - name: echo file
        image: tencentcom/scp
        settings:
          host: 
            - xx.xx.xx.xxx
            - xx.xx.xx.xxx
          username: root
          password: xxxx
          port: 22
          target: /data/release
          source:
            - release/*.tar.gz
```

### host 带 port 例子

```yaml
main:
  push:
    - stages:
      - name: echo file
        image: tencentcom/scp
        settings:
          host: 
            - xx.xx.xx.xxx:22
          username: root
          password: xxxx
          target: /data/release
          source:
            - release/*.tar.gz
```

### 命令超时例子

默认10m

```yaml
main:
  push:
    - stages:
      - name: echo file
        image: tencentcom/scp
        settings:
          host: 
            - xx.xx.xx.xx:22
          username: root
          password: xxxx
          command_timeout: 10s
          target: /data/release
          source:
            - release/*.tar.gz
```

### 引用密钥仓库配置文件获取密码例子

```yaml
# 密钥仓库 env.yml
PAASWORD: xxxx

# 声明指定镜像的插件任务能引用该配置文件
allow_images:
  - tencentcom/scp
# 声明指定仓库的流水线能引用该配置文件
allow_slugs:
  - group/repo
```

```yaml
main:
  push:
    - stages:
      - name: echo file
        # 引用密钥仓库配置文件
        imports: https://xxx/group/secret-repo/-/blob/main/env.yml
        image: tencentcom/scp
        settings:
          host: 
            - xx.xx.xx.xxx:22
          username: root
          # 引用密钥仓库配置文件中的变量
          password: $PAASWORD
          target: /data/release
          source:
            - release/*.tar.gz
```

### 引用密钥仓库配置文件获取 ssh key 例子

```yaml
# 密钥仓库 env.yml
SSH_KEY: |
  -----BEGIN OPENSSH PRIVATE KEY-----
  xxx
  -----END OPENSSH PRIVATE KEY-----

# 声明指定镜像的插件任务能引用该配置文件
allow_images:
  - tencentcom/scp
# 声明指定仓库的流水线能引用该配置文件
allow_slugs:
  - group/repo
```

```yaml
main:
  push:
    - stages:
      - name: echo file
        # 引用密钥仓库配置文件
        imports: https://xxx/group/secret-repo/-/blob/main/env.yml
        image: tencentcom/scp
        settings:
          host: 
            - xx.xx.xx.xxx:22
          username: root
          key: $SSH_KEY
          target: /data/release
          source:
            - release/*.tar.gz
```

### ssh key 带 passphrase 示例

```yaml
main:
  push:
    - stages:
      - name: copy file
        # 引用密钥仓库配置文件
        imports: http://xxx/-group/secret-repo/-/blob/main/env.yml
        image: tencentcom/scp
        settings:
          host: 
            - xx.xx.xx.xxx:22
          username: root
          key: $SSH_KEY_PHRASE
          passphrase: xxx
          target: /data/release
          source:
            - release/*.tar.gz
```

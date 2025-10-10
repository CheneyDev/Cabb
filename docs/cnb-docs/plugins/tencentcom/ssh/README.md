# SSH

通过 ssh 在远端 host 执行命令。

## 参数说明

- `host`: 目标机器 hostname 或 IP
- `port`: 目标主机 ssh port
- `protocol`: 要使用的 IP 协议：可以是 tcp、tcp4 或 tcp6
- `username`: 目标主机用户名
- `password`: 目标主机密码
- `key`: 私钥文本
- `key_path`: 私钥路径
- `passphrase`: 私钥密码短语
- `script`: 在远端服务器执行的命令
- `script_stop`: 命令失败时停止执行后续命令
- `timeout`: SSH 连接建立的最长时间， 默认为 30 秒
- `command_timeout`: 执行命令的最长时间， 默认为 10 分钟
- `proxy_host`: 代理的 hostname 或 IP
- `proxy_port`: 代理主机的 ssh port
- `proxy_protocol`: 用于代理的 IP 协议：可以是 tcp、tcp4 或 tcp6
- `proxy_username`: 代理主机用户名
- `proxy_password`: 代理主机密码
- `proxy_key`: 代理主机私钥明文文本
- `proxy_key_path`: 代理主机私钥的路径
- `proxy_passphrase`: 代理主机私钥密码短语

## 在 云原生构建 上使用

简单示例:

```yaml
main:
  push:
    - stages:
      - name: echo file
        image: tencentcom/ssh
        settings:
          host: xx.xx.xx.xxx
          username: root
          password: xxxx
          port: 22
          script:
            - echo hello world
            - echo test > ~/test.txt
```

多台目标机器例子：

```yaml
main:
  push:
    - stages:
      - name: echo file
        image: tencentcom/ssh
        settings:
          host: 
            - xx.xx.xx.xxx
            - xx.xx.xx.xxx
          username: root
          password: xxxx
          port: 22
          script:
            - echo hello world
            - echo test > ~/test.txt
```

host 带 port 例子：

```yaml
main:
  push:
    - stages:
      - name: echo file
        image: tencentcom/ssh
        settings:
          host: 
            - xx.xx.xx.xxx:22
          username: root
          password: xxxx
          script:
            - echo hello world
            - echo test > ~/test.txt
```

命令超时例子：

```yaml
main:
  push:
    - stages:
      - name: echo file
        image: tencentcom/ssh
        settings:
          host: 
            - xx.xx.xx.xx:22
          username: root
          password: xxxx
          command_timeout: 10s
          script:
            - sleep 15s
```

引用密钥仓库配置文件获取密码例子:

```yaml
# 密钥仓库 env.yml
PAASWORD: xxxx

# 声明指定镜像的插件任务能引用该配置文件
allow_images:
  - tencentcom/ssh
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
        image: tencentcom/ssh
        settings:
          host: 
            - xx.xx.xx.xxx:22
          username: root
          # 引用密钥仓库配置文件中的变量
          password: $PAASWORD
          script:
            - echo hellworld
```

引用密钥仓库配置文件获取 ssh key 例子：

```yaml
# 密钥仓库 env.yml
SSH_KEY: |
  -----BEGIN OPENSSH PRIVATE KEY-----
  xxx
  -----END OPENSSH PRIVATE KEY-----

# 声明指定镜像的插件任务能引用该配置文件
allow_images:
  - tencentcom/ssh
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
        image: tencentcom/ssh
        settings:
          host: 
            - xx.xx.xx.xxx:22
          username: root
          key: $SSH_KEY
          script:
            - echo hellworld
```

脚本失败后停止执行后续脚本示例:

```yaml
main:
  push:
    - stages:
      - name: echo file
        image: tencentcom/ssh
        settings:
          host: 
            - xx.xx.xx.xxx:22
          username: root
          password: xxxx
          script_stop: true
          script:
            - echo test1 > ~/test.txt
            - echo1 hellworld
            # 该命令不会执行
            - echo test2 > ~/test.txt
```

ssh key 带 passphrase 示例:

```yaml
main:
  push:
    - stages:
      - name: echo file
        # 引用密钥仓库配置文件
        imports: http://xxx/-group/secret-repo/-/blob/main/env.yml
        image: tencentcom/ssh
        settings:
          host: 
            - xx.xx.xx.xxx:22
          username: root
          key: $SSH_KEY_PHRASE
          passphrase: xxx
          script:
            - echo hellworld
```

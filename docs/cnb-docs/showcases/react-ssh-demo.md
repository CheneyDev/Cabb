# CNB 实现 React 构建并将静态资源通过 SSH 部署到目标服务器

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该方案通过在项目配置react-ssh-secret.yml文件存储SSH连接信息，并通过.cnb.yml定义构建流程和部署策略，实现React应用的自动化构建与远程部署。核心流程包含Node.js环境依赖安装、生产构建、目录检查、基于密码或私钥两种SSH认证方式的文件传输，部署目标路径设定为服务器的/opt目录。其中通过Docker镜像缓存加速构建，支持使用SCP插件完成安全传输，关键参数均通过环境变量实现动态注入。

基本逻辑：配置 .cnb.yml 文件，用于登录 -> 构建 -> SSH 上传到目标服务器。

仓库地址：[CNB 实现 React 构建并将静态资源通过 SSH 部署到目标服务器](https://cnb.cool/examples/ecosystem/react-ssh-demo)

## 1、配置目标服务器的配置密钥 react-ssh-secret.yml

在密钥仓库中创建一个 `react-ssh-secret.yml` 文件，用于配置目标服务器的 SSH 连接信息，后续在 `.cnb.yml` 文件中通过 `imports` 引入。

```yaml
# ssh cfg
REMOTE_HOST: xxx
REMOTE_USERNAME: xxx
REMOTE_PASSWORD: xxx
REMOTE_PORT: xxx
PRIVATE_KEY: |
    -----BEGIN RSA PRIVATE KEY-----
    xxxxxxxx 
    -----END RSA PRIVATE KEY-----
```

## 2、配置 .cnb.yml 文件

在项目根目录下创建一个 `.cnb.yml` 文件，用于配置构建任务：安装依赖 -> 编译 -> 通过 SSH 上传到目标服务器。

```yaml
main:
  push:
    #导入环境变量：https://docs.cnb.cool/zh/env.html#dao-ru-huan-jing-bian-liang
    - imports:
        - https://cnb.cool/examples/secrets/-/blob/main/react-ssh-secret.yml
      #声明构建环境：https://docs.cnb.cool/zh/
      docker:
        image: node:20
        #volumes缓存：https://docs.cnb.cool/zh/grammar/pipeline.html#volumes
        volumes:
          # 使用缓存，同时更新
          - /root/.npm:cow
      stages:
        - name: 编译
          script: npm install && npm run build
        - name: 查看目录
          script: ls -a
          # https://docs.cnb.cool/zh/plugins/public/open-source/scp/scp 
          # 使用插件将文件上传到服务器

       - name: 使用SCP插件：账号密码
         # 想查看更多插件及其用法请移步，【插件市场】https://ci.coding.net/docs/plugins/index.html
         image: tencentcom/scp
         settings:
           # 环境变量用法参考，【环境变量】https://ci.coding.net/docs/env.html
           host: $REMOTE_HOST
           username: $REMOTE_USERNAME
           password: $REMOTE_PASSWORD
           port: $REMOTE_PORT
           target: /opt
           source: ./build

        - name: 使用SCP插件：私钥
          # 插件用法请移步，【插件市场】https://ci.coding.net/docs/plugins/index.html
          image: tencentcom/scp
          settings:
            # 环境变量用法参考，【环境变量】https://ci.coding.net/docs/env.html
            host: $REMOTE_HOST
            key: $PRIVATE_KEY
            port: $REMOTE_PORT
            target: /opt
            source: ./build
```

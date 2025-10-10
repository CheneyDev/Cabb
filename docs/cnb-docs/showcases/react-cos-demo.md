# CNB 配置实现 react 构建静态资源上传到腾讯云 cos

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该文本采用.cnb.yml文件实现React静态资源自动部署流程，核心步骤包括：通过node:20镜像创建构建环境，使用volumes卷实现npm缓存优化，分阶段执行「安装依赖→构建→目录检查」任务，通过tencentcom/coscli插件联动COS_SECRET_ID等环境变量配置，最终将本地build目录打包上传到指定腾讯云COS存储桶。实现过程包含构建环境声明、分步骤脚本执行、云资源配置及文件传输命令的嵌套调用。

基本实现逻辑：配置 .cnb.yml 文件，用于登录 -> 构建 -> 上传到腾讯云 cos。

仓库地址：[CNB 实现 react 构建静态资源上传到腾讯云 cos](https://cnb.cool/examples/ecosystem/react-cos-demo)

## 1、配置 .cnb.yml 文件

在项目根目录下创建一个 `.cnb.yml` 文件，用于配置构建任务：安装依赖 -> 编译 -> 上传到腾讯云 cos。

```yaml
main:
  push:
    - 
      #声明构建环境：https://docs.cnb.cool/zh/
      docker:
        image: node:20
        #volumes缓存：https://docs.cnb.cool/zh/grammar/pipeline.html#volumes
        volumes:
          # 使用缓存，同时更新
          - /root/.npm:cow
      
      #导入环境变量：https://docs.cnb.cool/zh/env.html#dao-ru-huan-jing-bian-liang
      # imports: https://cnb.cool/xxx/xxxx/-/blob/main/react_cos_secret.yml
      stages:
        - name: 编译
          script: npm install && npm run build
        - name: 查看目录
          script: ls -a
        # 通过插件上传文件到腾讯云 cos
        - name: tencentcom 上传文件
          image: tencentcom/coscli

          #插件地址: https://docs.cnb.cool/zh/plugins/public/tencentcom/coscli
          #配置coscli的COS_SECRET_ID、COS_SECRET_KEY、COS_BUCKET、COS_REGION
          #执行coscli cp将本地的./dist文件夹上传到$COS_BUCKET中
          commands: |
            coscli config set --secret_id $COS_SECRET_ID --secret_key $COS_SECRET_KEY 
            coscli config add --init-skip=true -b $COS_BUCKET -r $COS_REGION
            coscli cp ./build cos://$COS_BUCKET -r
```

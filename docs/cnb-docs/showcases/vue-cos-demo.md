# CNB 配置实现 vue 构建静态资源上传到腾讯云 cos

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该实现方案基于.cnb.yml配置文件实现Vue静态资源构建及腾讯云COS上传流程。通过定义node:20镜像安装依赖并编译，使用tencentcom/coscli镜像执行coscli命令配置COS凭证和存储桶，将构建后的./dist目录同步至云端存储。配置采用环境变量注入敏感信息，利用缓存卷机制加速依赖管理，并通过分阶段任务实现编译验证、目录检查及最终文件传输操作。

基本实现逻辑：配置 .cnb.yml 文件，用于登录 -> 构建 -> 上传到腾讯云 cos。


仓库地址：[CNB 实现 vue 构建静态资源上传到腾讯云 cos](https://cnb.cool/examples/ecosystem/vue-cos-demo)

## 1、配置 .cnb.yml 文件

在项目根目录下创建一个 `.cnb.yml` 文件，用于配置构建任务：安装依赖 -> 编译 -> 上传到腾讯云 cos。

```yaml
main:
  push:
    - 
      #声明构建环境：https://docs.cnb.cool/zh/
      docker:
        image: node:20
        #volumes缓存文档：https://docs.cnb.cool/zh/grammar/pipeline.html#volumes
        volumes:
          # 使用缓存，同时更新
          - /root/.npm:cow
      stages:
        - name: 编译
          script: npm install && npm run build
        - name: 查看目录
          script: ls -a
        - name: coscli 上传文件
          image: tencentcom/coscli
          # 导入环境变量：https://docs.cnb.cool/zh/env.html#dao-ru-huan-jing-bian-liang
          imports: https://cnb.cool/xxxx/xxxx/-/blob/main/vue_cos_secret.yml
          #插件地址: https://docs.cnb.cool/zh/plugins/public/tencentcom/coscli
          #配置coscli的COS_SECRET_ID、COS_SECRET_KEY、COS_BUCKET、COS_REGION
          #执行coscli cp将本地的./dist文件夹上传到$COS_BUCKET中
          commands: |
            coscli config set --secret_id $COS_SECRET_ID --secret_key $COS_SECRET_KEY 
            coscli config add --init-skip=true -b $COS_BUCKET -r $COS_REGION
            coscli cp ./dist cos://$COS_BUCKET -r
```
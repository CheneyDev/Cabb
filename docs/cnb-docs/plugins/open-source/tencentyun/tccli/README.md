# 腾讯云命令行工具(TCCLI)

通过腾讯云命令行工具，您可以快速轻松的调用腾讯云 API来管理您的腾讯云资源。

## 在 Docker 上使用

```shell
docker run --rm -it tencentcom/tencentcloud-cli --version
docker run --rm -it tencentcom/tencentcloud-cli help
```

## 在 云原生构建 上使用

在私有仓库或密钥仓库中创建一个 env.yml 文件，声明需要使用的环境变量：

```yaml
# 腾讯云 API 密钥 id
T_SECRET_ID: xxxx
# 腾讯云 API 密钥 key
T_SECRET_KEY: xxxx
```

### cvm 创建实例

API 详见[官方文档](https://cloud.tencent.com/document/product/213/15730)

编写 .cnb.yml 文件：

```yaml
main:
  push:
    # 引入 env.yml，导入其声明的环境变量
    - imports: https://xxx/env.yml
      stages:
        - name: create cvm instance
          image: tencentcom/tencentcloud-cli
          script: |
            tccli configure set secretId $T_SECRET_ID
            tccli configure set secretKey $T_SECRET_KEY
            tccli cvm RunInstances \
            --cli-unfold-argument \
            --region ap-guangzhou \
            --Placement.Zone ap-guangzhou-6 \
            --ImageId img-eb30mz89
```

### tsf 部署容器应用

API 详见[官方文档](https://cloud.tencent.com/document/product/649/36071)

在 env.yml 补充 tsf 相关环境变量

```yaml
# tsf 制品管理 镜像仓库 用户名
TSF_DOCKER_USER: xxxx
# tsf 制品管理 镜像仓库 密码
TSF_DOCKER_PASSWORD: xxxx
```

编写 .cnb.yml 文件，具体流程为：

1. 声明流水线中需要使用docker
2. 引入上文中定义的 env.yml，导入其声明的环境变量
3. 以 commit sha 为 tag，导出 环境变量，供后续任务复用
4. docker login
5. docker build
6. docker push
7. tsf deploy

```yaml
main:
  push:
    - services:
        # 声明流水线中需要使用docker
        - docker
      # 引入上文中定义的 env.yml，导入其声明的环境变量
      imports: https://xxx/env.yml
      stages:
        # 导出 image tag 环境变量，供后续任务复用
        - name: generate image tag
          script: echo -e "ccr.ccs.tencentyun.com/tsf_xxx/tsf-demo:$CNB_COMMIT"
          exports:
            info: IMAGE_TAG
        - name: docker login
          script: docker login -u $TSF_DOCKER_USER -p $TSF_DOCKER_PASSWORD ccr.ccs.tencentyun.com
        - name: docker build
          script: docker build -t $IMAGE_TAG .
        # 推送到 tsf 制品库
        - name: docker push
          script: docker push $IMAGE_TAG
        # 部署到 tsf
        - name: tsf deploy
          image: tencentcom/tencentcloud-cli
          script: |
            tccli configure set secretId $T_SECRET_ID
            tccli configure set secretKey $T_SECRET_KEY
            tccli tsf DeployContainerGroup \
            --region ap-guangzhou \
            --GroupId group-xxx \
            --InstanceNum 1 \
            --CpuRequest '0.5' \
            --MemRequest 516 \
            --RepoName tsf_xxx/tsf-demo \
            --TagName $CNB_COMMIT
```

## 更多文档

更多 腾讯云 api 可以参考[官方文档](https://cloud.tencent.com/document/api)

更多 腾讯云 tccli 可以参考[官方文档](https://cloud.tencent.com/document/product/440)

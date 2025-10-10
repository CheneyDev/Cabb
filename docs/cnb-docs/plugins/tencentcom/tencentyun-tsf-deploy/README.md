# tencentyun-tsf-deploy

基于 [TCCLI](https://cloud.tencent.com/document/product/440/6176)
实现一键部署**腾讯云微服务平台**容器应用（DeployContainerGroup）的镜像插件。

## 参数说明

- secretId：云 API 密钥 SecretId。
- secretKey：云 API 密钥 SecretKey。
- args：其余参数可以通过命令行 -- 的方式传入，请参考下面例子。

### 在 Docker 中使用

```sh
docker run --rm \
  -e PLUGIN_SECRETID=xxx \
  -e PLUGIN_SECRETKEY=xxx \
  -e PLUGIN_ARGS="--region ap-guangzhou --GroupId group-xxx \
   --InstanceNum 1 --CpuRequest '0.5' --MemRequest 516 \
    --RepoName tsf_xxx/tsf-demo --TagName v1" \
  -v $(pwd):$(pwd) -w $(pwd) \
  tencentcom/tencentyun-tsf-deploy
```

## 在 云原生构建 中使用

```yaml
# env.yml
TSF_DOCKER_USER: xxx
TSF_DOCKER_PASSWORD: xxx
T_SECRET_ID: xxx
T_SECRET_KEY: xxx
```

```yaml
# .cnb.yml
main:
main:
  push:
    - services:
        - docker
      imports: https://xxx/env.yml
      stages:
        - name: generate image tag
          script: echo -e "ccr.ccs.tencentyun.com/tsf_xxx/tsf-demo:$CNB_COMMIT"
          exports:
            info: IMAGE_TAG
        - name: show image tag
          script: echo -e $IMAGE_TAG
        - name: docker login
          script: docker login -u $TSF_DOCKER_USER -p $TSF_DOCKER_PASSWORD ccr.ccs.tencentyun.com
        - name: docker build
          script: docker build -t $IMAGE_TAG .
        - name: docker push
          script: docker push $IMAGE_TAG
        - name: tsf deploy
          image: tencentcom/tencentyun-tsf-deploy:dev
          settings:
            secretId: $T_SECRET_ID
            secretKey: $T_SECRET_KEY
            args: >
              --region ap-guangzhou
              --GroupId group-xxx
              --InstanceNum 1
              --CpuRequest '0.5'
              --MemRequest 516
              --RepoName tsf_xxx/tsf-demo
              --TagName $CNB_COMMIT
```

更多配置项请查阅：[部署容器应用](https://cloud.tencent.com/document/product/649/36071)

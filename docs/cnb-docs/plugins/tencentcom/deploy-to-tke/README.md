# 腾讯云 TKE 镜像更新插件

场景化封装 腾讯云 [容器服务](https://console.cloud.tencent.com/tke2) API
实现工作负载 Deployment 、StatefulSet 类型的**镜像更新**操作

## 前置条件

* 使用的是腾讯云容器服务(TKE)集群,并且已经部署了对应 Deployment、StatefulSet 资源
* 访问集群用的腾讯云 secret_id 、secret_key, 并且已经设置好账号集群权限。
设置方式见[doc/access.md](https://cnb.cool/cnb/plugins/tencentcom/deploy-to-tke/-/blob/main/doc/access.md)

## 参数说明

* `secret_id` 密钥 ID, 详情见 [doc/access.md](https://cnb.cool/cnb/plugins/tencentcom/deploy-to-tke/-/blob/main/doc/access.md)
* `secret_key` 密钥 Key
* `region` 集群地域,格式如：`ap-nanjing`。详情见 [doc/regions.md](https://cnb.cool/cnb/plugins/tencentcom/deploy-to-tke/-/blob/main/doc/regions.md)
* `cluster_id` 集群ID,格式如： `cls-m9miwj4u`
* `namespace` 工作负载所在的集群命名空间,如 default
* `workload_kind` 工作负载类型。支持 `deployment`、`statefulset`
* `workload_name` 工作负载名称
* `container_names` 容器名称, 如有多个用英文`,`分割
* `container_images` 待更新的容器最新镜像, 如有多个用英文`,`分割

## 运行结果

* 权限、参数配置正确后，如果插件更新镜像成功，则插件成功退出，否则失败退出。机制类似 kubectl set image
* 更新镜像完成后，会监控 10 分钟 pod 滚动状态并打印，方便在流水线日志查看。
此步骤不会影响插件成功或失败退出结果，最终 pod 滚动结果可在集群中查看详情。
* 运行机制：类似 kubectl set image, 此插件不会校验待更新镜像是否存在，镜像是否有效、能正常启动。
会直接更新工作负载 yaml 中的 image 字段，并由 K8s 异步调度后续步骤，如 pod 滚动等。

## 示例

### 在 云原生构建 上使用

#### a.快速用法

```yaml
# 更新示例镜像配置
main:
  push:
    - stages:
      - name: 使用tke插件更新镜像
        image: tencentcom/deploy-to-tke
        settings:
          secret_id: AKID***MpL4
          secret_key: mRH1***wu0C
          region: ap-***
          cluster_id: cls-***
          namespace: default
          workload_kind: deployment
          workload_name: my-***-deployment
          container_names: container-***-1
          # 可以使用变量形式如 container_images: ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:${CNB_COMMIT}
          container_images: nginx-***:v1
          

# a.快速用法 完
```

#### b.高级用法

如果不希望写密钥明文、可以在密钥（私有）仓库增加一份密钥信息文件 `your_secrets.yaml`,内容如下：

```yaml
# your_secrets.yml
secret_id: AKID***MpL4
secret_key: mRH1***wu0C
```

配置 `.cnb.yml`,引用上述文件导入环境变量,内容如下：

```yaml
# 更新示例镜像
main:
  push:
    - stages:
      - name: 使用tke插件更新镜像
        image: tencentcom/deploy-to-tke
        settingsFrom: https://cnb.cool/***/my-secret-repo/-/blob/main/your_secrets.yaml
        settings:
          region: ap-***
          cluster_id: cls-***
          namespace: default
          workload_kind: deployment
          workload_name: my-***-deployment
          container_names: container-***-1
          container_images: ${CNB_DOCKER_REGISTRY}/${CNB_REPO_SLUG_LOWERCASE}:${CNB_COMMIT}
# b.高级用法 完
```

### 在 Docker 上使用

```shell
docker run --rm  \
    -e PLUGIN_SECRET_ID="***" \
    -e PLUGIN_SECRET_KEY="***" \
    -e PLUGIN_REGION="ap-shanghai" \
    -e PLUGIN_CLUSTER_ID="cls-***" \
    -e PLUGIN_NAMESPACE="development" \
    -e PLUGIN_WORKLOAD_KIND="deployment" \
    -e PLUGIN_WORKLOAD_NAME="my-***-deployment" \
    -e PLUGIN_CONTAINER_NAMES="container-***-1,container-***-2" \
    -e PLUGIN_CONTAINER_IMAGES="nginx-***:v1,nginx-***:v2" \
    tencentcom/deploy-to-tke
```

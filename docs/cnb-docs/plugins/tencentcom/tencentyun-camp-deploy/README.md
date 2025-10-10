# 腾讯云 camp 部署

![badge](https://cnb.cool/cnb/plugins/tencentcom/tencentyun-camp-deploy/-/badge/git/latest/ci/pipeline-as-code)
![badge](https://cnb.cool/cnb/plugins/tencentcom/tencentyun-camp-deploy/-/badge/git/latest/ci/status/push)

场景化封装 腾讯云 [应用管理平台](https://console.cloud.tencent.com/camp/app) API 实现更新实例镜像

## 参数说明

- `secret_id` 密钥 ID
- `secret_key` 密钥 Key
- `action` 操作名称，
  - `ModifyComponentImages`: 批量更新实例镜像
  - `DeleteInstance`: 销毁实例 ***慎用***
- `instance_url` 实例 URL，若指定则不使用下列 `project_id`、`application_id`、
`instance_id`、`environment_name` 参数，转而从 url 中解析
- `project_id` 业务 ID
- `application_id` 应用 ID
- `instance_id` 实例 ID
- `environment_name` 环境名称
- `result_path` 若指定则将结果以 JSON 格式字符串输出到指定文件

以下是 `ModifyComponentImages` 需要的参数

- `container_names` 容器名称数组，顺序与 `container_images` 顺序一致
- `container_images` 容器镜像数组，顺序与 `container_names` 顺序一致
- `max_surge` 升级中允许每个工作负载中超出所需规模的最大 Pod 数量，默认值 `0`
- `max_unavailable` 升级中允许每个工作负载中最大不可用的 Pod 数量，默认值 `1`
- `in_place_update_flag` 是否原地更新，默认为 `true`

## 示例

### 在 Docker 上使用

```shell
docker run --rm -v $(pwd):$(pwd) -w $(pwd) \
    -e PLUGIN_SECRET_ID="xxx" \
    -e PLUGIN_SECRET_KEY="xxx" \
    -e PLUGIN_ACTION="ModifyComponentImages" \
    -e PLUGIN_PROJECT_ID="prj-xxx" \
    -e PLUGIN_ENVIRONMENT_NAME="development" \
    -e PLUGIN_APPLICATION_ID="app-xxx" \
    -e PLUGIN_INSTANCE_ID="tad-xxx" \
    -e PLUGIN_CONTAINER_NAMES="xxx-container-1,xxx-container-2" \
    -e PLUGIN_CONTAINER_IMAGES="xxx-demo:v2,xxx-demo:t2" \
    tencentcom/tencentyun-camp-deploy
```

## 在 云原生构建 上使用

在密钥（私有仓库）配置一份密钥信息文件 `your_secrets.yaml`，内容如下：

```yaml
# your_secrets.yml
SECRET_ID: xxxx
SECRET_KEY: xxx
```

配置 `.cnb.yml`，引用上述文件导入环境变量，内容如下：

```yaml
# 更新示例镜像，不使用 实例 URL
main:
  push:
    - imports: https://xxx/your_secrets.yml
      stages:
      - name: camp deploy
        image: tencentcom/tencentyun-camp-deploy
        settings:
          secret_id: $SECRET_ID
          secret_key: $SECRET_KEY
          action: ModifyComponentImages
          project_id: prj-xxx
          environment_name: development
          application_id: app-xxx
          instance_id: tad-xxx
          container_names:
            - xxx-container-1
            - xxx-container-2
          container_images: 
            - xxx-demo:v1
            - xxx-demo:t1
```

```yaml
# 更新示例镜像，使用 实例 URL
main:
  push:
    - imports: https://xxx/your_secrets.yml
      stages:
      - name: camp deploy
        image: tencentcom/tencentyun-camp-deploy
        settings:
          secret_id: $SECRET_ID
          secret_key: $SECRET_KEY
          action: ModifyComponentImages
          instance_url: https://console.cloud.tencent.com/camp/app/instance/info?appId=app-xxx&projectId=prj-xxx&instanceId=tad-xxx&envName=development
          container_names:
            - xxx-container-1
            - xxx-container-2
          container_images: 
            - xxx-demo:v1
            - xxx-demo:t1
```

```yaml
# 销毁示例实例，使用 实例 URL
main:
  push:
    - imports: https://xxx/your_secrets.yml
      stages:
      - name: camp deploy
        image: tencentcom/tencentyun-camp-deploy
        settings:
          secret_id: $SECRET_ID
          secret_key: $SECRET_KEY
          action: DeleteInstance
          instance_url: https://console.cloud.tencent.com/camp/app/instance/info?appId=app-xxx&projectId=prj-xxx&instanceId=tad-xxx&envName=development
```

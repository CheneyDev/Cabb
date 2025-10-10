---
title: ISSUE通知到企业微信群
permalink: https://docs.cnb.cool/zh/build/showcase/issue-notice-group.html
summary: 要在企业微信群中接收 ISSUE 相关通知，需先添加群机器人并复制其 Webhook 地址。配置流水线时，可自定义 issue 打开、重新打开和关闭的通知内容，通过特定镜像和设置中的机器人地址及相应格式内容来实现，变量可从环境变量中获取。
---

## 添加群机器人

于企业微信群中添加机器人，复制得到的 `Webhook` 地址。

## 配置流水线

示例：

```yaml
.issue-open: &issue-open
  - name: issue-notice
    image: tencentcom/wecom-message
    settings:
      robot: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx
      content: |
        > **有人提issue啦**
        > **标  题:** $CNB_ISSUE_TITLE
        > **发起人:** $CNB_ISSUE_OWNER
        > [查看详情]($CNB_EVENT_URL)
.issue-reopen: &issue-reopen
  - name: issue-notice
    image: tencentcom/wecom-message
    settings:
      robot: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx
      content: |
        > **$CNB_BUILD_USER重新打开了一个issue**
        > **标  题:** $CNB_ISSUE_TITLE
        > **发起人:** $CNB_ISSUE_OWNER
        > [查看详情]($CNB_EVENT_URL)
.issue-close: &issue-close
  - name: issue-close
    image: tencentcom/wecom-message
    settings:
      robot: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx
      content: |
        > **$CNB_BUILD_USER关闭了一个issue**
        > **标  题:** $CNB_ISSUE_TITLE
        > **发起人:** $CNB_ISSUE_OWNER
        > [查看详情]($CNB_EVENT_URL)
$:
  issue.close:
    - stages:
        - *issue-close
  issue.reopen:
    - stages:
        - *issue-reopen
  issue.open:
    - stages:
        - *issue-open
```

`robot` 填之前复制的 `Webhook` 地址。

流水线触发者是新建、关闭、重新打开这个 `issue` 的用户。

具体信息格式可自定义，可用变量参考[环境变量](../env.md)。

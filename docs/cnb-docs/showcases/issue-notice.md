# issue 信息通知到群聊示例

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该示例演示了通过GitHub事件触发群聊通知和自动更新Issue标签的功能。主要覆盖三个操作：1.配置issue.open/reopen/close事件在新建、重开、关闭Issue时发送包含标题、发起人、详情链接的Markdown通知到企业微信群，通过$tencentcom/wecom-message`镜像实现；2.主分支push事件使用`git:issue-update`任务自动更新Issue标签（如移除"开发中"、添加"已解决"）；3.所有通知均引用环境变量${REVIEW_ROBOT_ID}指定机器人ID，并通过模板变量动态填充事件详情。

可以用于通知 issue 信息到群聊，用于及时获取 issue 信息，或者用于 issue 标签自动更新。

仓库地址：[issue 信息通知到群聊示例](https://cnb.cool/examples/ecosystem/issue-notice)

## 如何配置

issue.open 事件：issue 创建后，会触发 issue.open 事件，将 issue 信息通知到群聊
issue.reopen 事件：issue 关闭后被重新打开，会触发 issue.reopen 事件，将 issue 信息通知到群聊
issue.close 事件：issue 关闭后，会触发 issue.close 事件，将 issue 关闭信息通知到群聊
push 事件：push 事件触发后，使用 git:issue-update 内置任务自动修改 issue 标签

## 配置 .cnb.yml 文件，用于 issue 信息通知到群聊

```yaml
# 主分支 push 事件，用于更新 issue 标签，例如：开发中 -> 已解决
main:
  push:
    - stages:
      - name: 更新 issue
        type: git:issue-update
        options:
          label:
            add: 已解决
            remove: 开发中
# 事件文档：https://docs.cnb.cool/zh/event.html
$:
  issue.open:
    - stages:
      - name: issue 信息通知到群聊
        # 文档 https://docs.cnb.cool/zh/plugins/public/tencentcom/wecom-message
        image: tencentcom/wecom-message
        imports: https://cnb.cool/xxx/xxx/-/blob/main/envs/wework-robots.yml
        settings:
          robot: $REVIEW_ROBOT_ID
          msgType: markdown
          content: |
            > **有人提issue啦**
            > **标  题:** $CNB_ISSUE_TITLE
            > **发起人:** $CNB_ISSUE_OWNER
            > [查看详情]($CNB_EVENT_URL)

  # issue 重新打开事件
  issue.reopen:
    - stages:
      - name: issue 重新打开信息通知到群聊
        image: tencentcom/wecom-message
        imports: https://cnb.cool/xxx/xxx/-/blob/main/envs/wework-robots.yml
        settings:
          robot: $REVIEW_ROBOT_ID
          content: |
            > **$CNB_BUILD_USER重新打开了一个issue**
            > **标  题:** $CNB_ISSUE_TITLE
            > **发起人:** $CNB_ISSUE_OWNER
            > [查看详情]($CNB_EVENT_URL)

            ${REVIEWED_BY}

  # issue 关闭事件
  issue.close:
    - stages:
      - name: issue 关闭信息通知到群聊
        image: tencentcom/wecom-message
        imports: https://cnb.cool/xxx/xxx/-/blob/main/envs/wework-robots.yml
        settings:
          robot: $REVIEW_ROBOT_ID
          content: |
            > **$CNB_BUILD_USER关闭了一个issue**
            > **标  题:** $CNB_ISSUE_TITLE
            > **发起人:** $CNB_ISSUE_OWNER
            > [查看详情]($CNB_EVENT_URL)
```
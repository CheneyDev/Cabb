# 代码评审通知，CODE REVIEW 通知到群聊示例

本文档地址：https://cnb.cool/examples/showcase

文档摘要：该示例展示了如何通过配置企业微信机器人实现代码评审通知流程：1）在密钥仓库创建`wework-robots.yml`文件定义机器人参数及权限范围；2）在工作流文件`.cnb.yml`中设置了两个通知节点——代码评审阶段自动添加评审人后触发群聊通知，合并前通过 squash 方式自动合并并再次发送通过通知；3）通知内容包含需求标题、链接、评审人及提交者等动态信息，合并通知额外支持Markdown格式且@相关成员。

可以用于通知代码评审信息到群聊，用于及时获取代码评审信息。示例中，配置了如何将 CODE REVIEW 消息通知到群聊

仓库地址：[CODE REVIEW 通知到群聊示例](https://cnb.cool/examples/ecosystem/code-review-notice)


## 1. 在密钥仓库中创建一个 `wework-robots.yml` 文件，用于配置企业微信机器人的信息。
```yaml
# wework-robots.yml

OCI_REVIEW_ROBOT: xxxx-xxxx

allow_slugs:
  - examples/**

allow_events:
  - pull_request
  - pull_request.*

allow_images:
  - tencentcom/wecom-message
```

## 2. 配置 .cnb.yml 文件，用于 CODE REVIEW 消息通知到群聊



```yaml
main:
  pull_request:
    - stages:
      # 添加评审人文档：https://docs.cnb.cool/zh/internal-steps/git/reviewer.html
      - name: 添加评审人
        type: git:reviewer
        options:
          type: add-reviewer-from-group-members
          count: 2
        exports:
          reviewersForAt: CURR_REVIEWER_FOR_AT
      # 通知文档：https://docs.cnb.cool/zh/plugins/public/tencentcom/wecom-message
      - name: 评审信息通知到群聊
        imports: https://cnb.cool/cnb/secrets/-/blob/main/envs/wework-robots.yml
        image: tencentcom/wecom-message
        settings:
          robot: $OCI_REVIEW_ROBOT
          content: |
            ${CNB_PULL_REQUEST_TITLE}
            [${CNB_EVENT_URL}](${CNB_EVENT_URL})

            ${CURR_REVIEWER_FOR_AT}
            from ${CNB_BUILD_USER}

  pull_request.mergeable:
    - stages:
        - name: CR 通过后自动合并
          type: git:auto-merge
          options:
            mergeType: squash
            removeSourceBranch: true
          exports:
            reviewedBy: REVIEWED_BY
        # 通知文档：https://docs.cnb.cool/zh/plugins/public/tencentcom/wecom-message
        - name: CR 通过后通知到群聊
          imports: https://cnb.cool/cnb/secrets/-/blob/main/envs/wework-robots.yml
          image: tencentcom/wecom-message
          settings:
            robot: $OCI_REVIEW_ROBOT
            msgType: markdown
            content: |
              CR 通过后自动合并: <@${CNB_PULL_REQUEST_PROPOSER}>
              　
              ${CNB_PULL_REQUEST_TITLE}
              [${CNB_EVENT_URL}](${CNB_EVENT_URL})
              　
              ${REVIEWED_BY}
```

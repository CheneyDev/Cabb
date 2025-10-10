# web-a11y-scanner

Web 无障碍 UI 测试工具，可以输出页面对无障碍的支持情况。

## 输入

- baseUrl: 必填，入口页面URL
- depth: 选填，扫描深度，建议根据实际情况选择1到3，默认为2
- device: 选填，mobile或desktop，默认为mobile
- whitelist: 选填，hostname白名单，多个hostname用英文半角逗号分隔，默认为空
- ignorelist: 选填，对特定页面URL忽略特定规则，
  多个页面用英文半角逗号分隔；URL与规则之间以及规则与规则之间用英文|间隔，
  不指定规则代表忽略当前页面所有规则，默认为空
- cookieList：选填，是cookieDict的数组，数组中的每个词典必须包含name、value、domain

## 输出

- `report`: 报告输出目录

## 在 Docker 上使用

```shell
docker run --rm \
  -e PLUGIN_BASEURL="https://w3.org/" \
  -e PLUGIN_DEPTH="1" \
  -e PLUGIN_DEVICE="mobile" \
  -e PLUGIN_WHITELIST="w3.org" \
  -e PLUGIN_IGNORELIST="https://w3.org/|duplicate-id" \
  -v $(pwd):$(pwd) -w $(pwd) \
  tencentcom/web-a11y-scanner
```

## 在 云原生构建 上使用

```yaml
main:
  push:
  - stages:
    - name: Accessibility Test Mobile
      image: tencentcom/web-a11y-scanner
      settings:
        baseUrl: https://w3.org/
        depth: 1
        device: mobile
        whitelist: w3.org
    - name: ls report/*
      script: ls -al report/*
```

# hello world

![pipeline-as-code](https://cnb.cool/cnb/plugins/cnbcool/hello-world/-/badge/git/latest/ci/pipeline-as-code)
![git-clone-yyds](https://cnb.cool/cnb/plugins/cnbcool/hello-world/-/badge/git/latest/ci/git-clone-yyds)
![push](https://cnb.cool/cnb/plugins/cnbcool/hello-world/-/badge/git/latest/ci/status/push)

一个打印参数的插件示例。

## 在 Docker 上使用

```shell
docker run --rm -v $(pwd):$(pwd) -w $(pwd) cnbcool/hello-world
```

## 在 云原生构建 上使用

```yml
main:
  push:
    - stages:
        - name: hello-world
          image: cnbcool/hello-world
          settings:
            text: Hello World!
```

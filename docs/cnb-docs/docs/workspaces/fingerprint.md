---
title: 云原生开发 SSH 密钥指纹验证
permalink: https://docs.cnb.cool/zh/workspaces/fingerprint.html
summary: 云原生开发支持通过SSH远程连接，为确保安全，需验证SSH密钥指纹。CNB云原生开发的SSH密钥指纹为SHA256:fnWZvpqd+VAIRJxaZdV1KVMFfDgcCjYrP2VSWQ68T/E，用户可通过查看本地`.ssh/known_hosts`文件来验证指纹，若不正确，可删除对应行后重新连接。
---
云原生开发支持 SSH 方式进行远程连接，为保证安全性，您可以验证 SSH 密钥指纹是否正确。

## SSH 密钥指纹

以下是 CNB 云原生开发的 SSH 密钥指纹：

```shell
SHA256:fnWZvpqd+VAIRJxaZdV1KVMFfDgcCjYrP2VSWQ68T/E
```

## 如何验证 SSH 密钥指纹

可通过查看本地 .ssh/know_hosts 文件，验证已信任远程主机的 SSH 密钥指纹是否正确。
如果 SSH 密钥指纹不正确，可删除 .ssh/know_hosts 文件中的对应行，重新连接即可。

在终端中运行命令，查看密钥指纹：

```bash
# 需先进入 .ssh 文件夹所在目录
ssh-keygen -lf .ssh/known_hosts
```

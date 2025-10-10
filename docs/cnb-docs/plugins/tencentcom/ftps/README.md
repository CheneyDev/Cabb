# FTPS

通过 ftp（lftp 命令） 上传制品

## 参数说明

- `hostname`: FTP 主机，包括端口
- `clean_dir`: 如果设置为 true，将在文件传输之前先删除目标目录（`dest_dir`）。与 `only_newer` 不能同时使用
- `chmod`: 如果为 false，将在文件传输后强制执行文件权限更改。默认为空
- `verify`: 如果为 true，则强制 SSL 证书验证，否则无验证。true 或 false，默认： false
- `secure`: 是否使用 SSL 加密。`true`: 使用 SSL 加密来保护 FTP 数据传输的隐私和安全。
`false` 或未设置: 不使用 SSL 加密
- `dest_dir`: 将文件放在远程服务器上的位置。例如：`/path/to/dest`
- `src_dir`: 用于上传的本地目录。例如：`/path/to/src`
- `exclude`: egrep-like 模式匹配从上传的文件中排除文件，支持数组格式或逗号分隔
- `include`: egrep 类模式匹配要上传的文件，支持数组格式或逗号分隔
- `auto_confirm`: 该参数有值时，表示启用 SFTP 自动确认功能。默认不启用自动确认。
在 SFTP 传输过程中，有时会遇到需要用户确认的操作，例如文件权限更改、文件删除等。
默认情况下，lftp 会在遇到这些操作时提示用户进行确认。
然而，在某些情况下，用户可能希望自动确认这些操作，以避免手动干预。
- `pre_action`: 在文件传输之前在服务器上执行的命令（在`clean_dir`之前执行）
- `post_action`: 在文件传输之后在服务器上执行的命令
- `debug`: 该参数有值时，启用调试模式；当参数为空时不启用调试模式。
lftp 命令以调试模式启动时，将输出详细的调试信息，包括发送和接收的网络数据包、内部状态变化等。

## 在 云原生构建 上使用

### 简单示例

```yaml
main:
  push:
  - stages:
    - name: deploy
      image: tencentcom/ftps
      env: 
        # 必要环境变量，FTP 主机用户名和密码
        FTP_PASSWORD: xxx
        FTP_USERNAME: xxx
      settings:
        hostname: example.com:21
        src_dir: ./
        dest_dir: /data/release
        exclude: ^\.git/$
```

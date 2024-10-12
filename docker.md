# docker

## 可用安全源

`/etc/docker/daemon.json `

```"dockerhub.com",
"hub.docker.com",
"login.docker.com",
"api.dso.docker.co",
"sessions-bugsnag.docker.com",
"docker.io",
"epsilon.6sense.com",
"notify-bugsnag.docker.com",
"desktop.docker.com",
"production.cloudflare.docker.com",
"registry-1.docker.io",
"auth.docker.io",
"resolver.dingtalk.com"
```

## 包含各种数据库

Makefile 命令

```
make buildDb
```

## sslmode 参数的作用

### disable：
`不使用 SSL 连接。
数据在传输过程中不会加密，可能存在安全风险。
适用于本地开发环境或信任网络环境。`
### allow：
`尝试使用 SSL 连接，但如果服务器不支持 SSL，则仍然进行非加密连接。
如果服务器支持 SSL，数据会在传输过程中加密。
可在需要时启用 SSL，但如果服务器不支持，仍可以连接。`
### prefer：
`尝试使用 SSL 连接，但如果服务器不支持 SSL，则仍然进行非加密连接。
如果服务器支持 SSL，数据会在传输过程中加密。
与 allow 类似，但更倾向于使用 SSL。`
### require：
`要求使用 SSL 连接。
如果服务器不支持 SSL，则连接失败。
数据会在传输过程中加密，确保数据安全。`
### verify-ca：
`要求使用 SSL 连接，并验证服务器的 SSL 证书。
如果服务器的 SSL 证书无效，连接会失败。
数据会在传输过程中加密，并验证服务器身份。`
### verify-full：
`要求使用 SSL 连接，并验证服务器的 SSL 证书以及主机名。
如果服务器的 SSL 证书无效或主机名不匹配，连接会失败。
数据会在传输过程中加密，并且要求严格验证服务器身份`



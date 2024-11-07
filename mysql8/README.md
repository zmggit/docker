| 连接地址     | 账号   | 密码       | 端口   | 库       |
| ------------ |------|----------|------| -------- |
| 172.16.2.126 | root | 123456   | 3232 | postgres |
| 172.16.2.126 | zmg  | Apom@321 | 3232 | postgres |
| 172.16.2.126 | zmgtest-1  | Apom@321 | 3232 | postgres |

## ssl 直接连接

`mysql --ssl-ca=/Users/zhumenggang/zhu_work/go/src/docker/mysql8/data/sanssl/ca.crt --ssl-cert=/Users/zhumenggang/zhu_work/go/src/docker/mysql8/data/sanssl/client.pem --ssl-key=/Users/zhumenggang/zhu_work/go/src/docker/mysql8/data/sanssl/client.key -u zmg -pApom@321 -h172.16.2.126 -P3232 --ssl-mode VERIFY_CA --ssl-cipher=AES128-SHA
`
## 普通链接

`mysql -u root -p123456 -h172.16.2.126 -P3232`

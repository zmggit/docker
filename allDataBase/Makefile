#安装docker
buildDb:
	@docker-compose -f ./docker-compose.yaml stop
	@docker-compose -f ./docker-compose.yaml up  -d
#删除数据盘
clearData:
	@rm -r ../mogo/data
	@rm -r ../postsql/data
#删除docker
clearDocker:
	@docker-compose -f ./docker-compose.yaml down
#停止
stopDb:
	@docker-compose -f ./docker-compose.yaml stop

#生成postsql ssl 服务端证书
serverPs:
    #服务端证书
	@openssl genrsa -des3 -out server.key 2048
	@openssl rsa -in server.key -out server.key
	@openssl req -new -key server.key -days 2650 -out server.crt -x509
	@cp server.crt root.crt
    #写入配置
	@echo   "ssl=on" >>../postsql/data/postgresql.conf
	@echo   "ssl_ca_file='root.crt'" >>../postsql/data/postgresql.conf
	@echo   "ssl_key_file='server.key'" >>../postsql/data/postgresql.conf
	@echo   "ssl_cert_file='server.crt'" >>../postsql/data/postgresql.conf
    #改服务器的pg_hba.conf文件禁止用户以非SSL连接数据库
	@echo   "hostssl all all 0.0.0.0/0 cert" >>../postsql/data/pg_hba.conf

##生成客户端证书 postSql
clientPs:
	@openssl genrsa -des3 -out postgresql.key 2048
	@openssl rsa -in postgresql.key -out postgresql.key
	#生成客户端csr文件
	#Common Name (e.g. server FQDN or YOUR name) []:postgres该项必须设置为要连接postgresql数据库的用户名，否则会默认使用当前计算机的用户名，导致证书使用时，认证失败。
	@openssl req -new -key postgresql.key -out postgresql.csr
    #生成客户端证书
	@openssl x509 -req -days 3650 -in postgresql.csr -CA root.crt -CAkey server.key -out postgresql.crt -CAcreateserial



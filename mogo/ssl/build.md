### 生成CA（证书颁发机构）私钥和证书:


`openssl genrsa -out ca.key 2048
openssl req -new -key ca.key -out ca.csr
openssl x509 -req -days 36500 -in ca.csr -signkey ca.key -out ca.pem`

### 生成服务器私钥和证书签名请求(CSR):


`openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr`

### 使用CA证书签署服务器端CSR以生成服务器端证书:


`openssl ca -days 36500 -in server.csr -out server.crt -cert ca.pem -keyfile ca.key`

### 将服务器私钥和证书合并为一个.pem文件:

`cat server.key server.crt > mongodb.pem`
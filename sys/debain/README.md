## 构建Docker镜像:

``bash docker build -t debian-ssh . ``

## 运行Docker容器并映射端口（假设你想映射SSH端口到主机的2222端口）

``docker run --name debian-ssh-server -d -p 2222:22 debian-ssh``

## 使用SSH客户端连接到容器:

`ssh root@localhost -p 2222`




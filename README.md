# 一个集成各种数据库，服务器，中间件的环境配置


```
启动docker-compose
docker-compose -f docker-compose.yaml up -d
```

```
使用dockerFile
# 构建Docker镜像
docker build -t my-php-app php/7.4

# 运行Docker容器
docker run -d -p 8080:80 my-php-app
```

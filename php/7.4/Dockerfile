# 使用官方PHP镜像作为基础镜像
FROM 1551605159/php7.4-swoole

# 设置工作目录
WORKDIR /app

# 将当前目录下的文件复制到工作目录
COPY . /app

# 设置环境变量，例如时区
ENV TZ=Asia/Shanghai

# 安装PHP扩展，例如PDO和MBSTRING

# 对外暴露端口，如果你的应用需要
EXPOSE 80

# 运行你的PHP应用
CMD ["php", "-S", "0.0.0.0:80", "-t", "/app"]
FROM debian:latest

# 安装ssh服务
RUN apt-get update && apt-get install -y openssh-server

# 设置root密码
RUN passwd root

# 启动ssh服务并设置运行级别
CMD ["/usr/sbin/sshd", "-D"]
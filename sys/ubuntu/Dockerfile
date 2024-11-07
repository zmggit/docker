FROM ubuntu:latest

# 安装ssh服务
RUN apt-get update && apt-get install -y openssh-server

RUN mkdir -p  /var/run/sshd

RUN chmod 777  /var/run/sshd

# 设置root密码
RUN echo 'root:123456' | chpasswd

RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /bin/bash
EXPOSE 22
EXPOSE 23
EXPOSE 21

RUN apt-get install telnet -y && apt-get instlal telnetd -y && apt-get install net-tools -y && apt-get install xinetd -y
Run  /etc/init.d/xinetd restart

# 启动ssh服务并设置运行级别
CMD ["/usr/sbin/sshd", "-D"]
# 使用官方 Ubuntu 基础镜像
FROM  ubuntu:20.04

# 安装 VNC 服务器和一个简单的桌面环境
RUN apt-get update && apt-get install -y \
    xfce4 \
    xfce4-terminal \
    tigervnc-standalone-server \
    && rm -rf /var/lib/apt/lists/*

# 设置 VNC 服务器密码
ENV VNC_PASSWORD=Z123456

# 创建一个新用户并设置密码，以便在 VNC 会话中使用
RUN useradd -m -d /home/ubuntu -s /bin/bash ubuntu \
    && echo ubuntu:${VNC_PASSWORD} | chpasswd \
    && usermod -aG sudo ubuntu

# 切换到新用户
USER ubuntu

# 启动 VNC 服务器（在后台）并启动桌面环境
VOLUME ["/home/ubuntu"]
CMD [ "/usr/bin/x11vnc", "-usepw", "-create", "-forever", "-display", ":1", "-bg", "-o", "/home/ubuntu/vnc.log", "2>&1" ]
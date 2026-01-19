# websocket

高可用 websocket

<!-- PROJECT SHIELDS -->

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

### 上手指南

WebSocket 与消息系统的本质认知

WebSocket 解决了什么问题

消息系统与聊天系统的区别

##### 开发前的配置要求

1. OpenResty + Nchan
2. Golang 1.15+
3. Redis 5+
4. Lua 5.1+

##### **安装步骤**

###### 1. 编译安装 OpenResty + Nchan

1. 更新系统

```sh
sudo apt update
sudo apt upgrade -y
```
2. 安装编译依赖

```sh
sudo apt install -y build-essential libreadline-dev libncurses5-dev libpcre3 libpcre3-dev zlib1g zlib1g-dev libssl-dev git perl make build-essential curl wget
```
3. 下载 OpenResty 源码  查看版本: https://openresty.org/en/download.html, 使用 1.21.4.1

```sh
wget https://openresty.org/download/openresty-1.21.4.1.tar.gz
tar -zxvf openresty-1.21.4.1.tar.gz
```
4. 下载 Nchan 模块

```sh
cd ~
git clone https://github.com/slact/nchan.git
```
5. 编译安装 OpenResty + Nchan，注意 --add-module=/home/peter/nchan 中的路径要改成你自己的路径

```sh
cd ~/openresty-1.21.4.1

./configure --prefix=/usr/local/openresty --with-luajit --with-pcre-jit --with-http_ssl_module --with-http_stub_status_module --with-http_gzip_static_module --with-http_realip_module --with-http_v2_module --with-stream --with-stream_ssl_module --with-stream_ssl_preread_module --add-module=/home/peter/nchan
```
6. 开始编译

```sh
make -j$(nproc)  # 编译使用多核加速
```
7. 安装

```sh
sudo make install
```
8. 验证是否成功编译Nchan模块

```sh
/usr/local/openresty/nginx/sbin/nginx -V
# 查看输出中是否包含 --add-module=/home/.../nchan
```
9. 检查 Nchan 模块是否可用

```sh
sudo /usr/local/openresty/nginx/sbin/nginx -t
# 如果输出没有报错，说明模块加载正常
sudo /usr/local/openresty/nginx/sbin/nginx -c  /usr/local/openresty/nginx/conf/nginx.conf
# 启动成功
````
###### 2. 启动一个 Nchan WebSocket Demo 服务
1. nginx.conf 配置内容如下：

```nginx
#
# Subscriber (WebSocket)
# ws://localhost:8080/ws/<channel_id>
#
location ~^/ws/(.+)$ {
    nchan_subscriber;
    nchan_channel_id $1;
}

#
# Publisher (HTTP)
# POST /pub/<channel_id>
#
location ~^/pub/(.+)$ {
    nchan_publisher;
    nchan_channel_id $1;
} 
```
2. 发送一条消息

```sh
curl -X POST http://localhost:8080/pub/chat -H "Content-Type: application/json" -d '{"msg": "hello websocket", "from": "curl"}'
```

### 文件目录说明
eg:

```
filetree 
├── ARCHITECTURE.md
├── LICENSE.txt
├── README.md
├── /account/
├── /bbs/
├── /docs/
│  ├── /rules/
│  │  ├── backend.txt
│  │  └── frontend.txt
├── manage.py
├── /oa/
├── /static/
├── /templates/
├── useless.md
└── /util/

```

### 开发的架构

请阅读 查阅为该项目的架构。

### 部署

暂无

### 使用到的框架

- [xxxxxxx](https://getbootstrap.com)


### 作者

xxx@xxxx


<!-- links -->
[your-project-path]: panjiawan//websocket
[contributors-shield]: https://img.shields.io/github/contributors/panjiawan/websocket.svg?style=flat-square
[contributors-url]: https://github.com/panjiawan/websocket/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/panjiawan/websocket.svg?style=flat-square
[forks-url]: https://github.com/panjiawan/websocket/network/members
[stars-shield]: https://img.shields.io/github/stars/panjiawan/websocket.svg?style=flat-square
[stars-url]: https://github.com/panjiawan/websocket/stargazers
[issues-shield]: https://img.shields.io/github/issues/panjiawan/websocket.svg?style=flat-square
[issues-url]: https://img.shields.io/github/issues/panjiawan/websocket.svg
[license-shield]: https://img.shields.io/github/license/panjiawan/websocket.svg?style=flat-square
[license-url]: https://github.com/panjiawan/websocket/blob/main/LICENSE.txt
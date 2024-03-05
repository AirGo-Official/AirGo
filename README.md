

<img width="200px" src="https://telegraph-image.pages.dev/file/c48a2f45ebf102dd66131.png" align="left"/>

# AirGo 前后端分离，多用户，多协议代理服务管理系统，简单易用

![License](https://img.shields.io/badge/License-GPL_v3.0-red)
![Go](https://img.shields.io/badge/Golang-orange?logo=Go&logoColor=white)
![Gorm](https://img.shields.io/badge/Gorm-yellow&logo=gorm)
![Gin](https://img.shields.io/badge/Gin-green?logo=)
![Vue](https://img.shields.io/badge/Vue.js-00b6ff?logo=vuedotjs&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-blue?logo=TypeScript&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-purple?logo=linux&logoColor=white)

<hr/>

支持：vless，vmess，shadowsocks，hysteria2

<div style="color: darkgray">AirGo, front and rear end separation, multi user, multi protocol proxy service management system, simple and easy to use.</div>
<div style="color: darkgray">support: vless，vmess，shadowsocks，hysteria2</div>

<hr/>



<!-- TOC -->
* [AirGo 前后端分离，多用户，多协议代理服务管理系统，简单易用](#airgo-前后端分离多用户多协议代理服务管理系统简单易用)
* [面板部分功能展示](#面板部分功能展示)
* [1、部署](#1部署)
  * [1-1 安装AirGo核心](#1-1-安装airgo核心)
    * [1-1-1 直接安装](#1-1-1-直接安装)
    * [1-1-2 使用docker安装](#1-1-2-使用docker安装)
  * [1-2 配置ssl（可选）](#1-2-配置ssl可选)
  * [1-3 部署前端静态资源（可选，但推荐）](#1-3-部署前端静态资源可选但推荐)
    * [1-3-1 部署到Vercel](#1-3-1-部署到vercel)
    * [1-3-2 部署到nginx、caddy等](#1-3-2-部署到nginxcaddy等)
  * [1-4 配置文件说明](#1-4-配置文件说明)
  * [1-5 启动](#1-5-启动)
* [2、对接节点](#2对接节点)
  * [2-1 XrayR](#2-1-xrayr)
  * [2-2 hysteria2](#2-2-hysteria2)
* [TG群组：https://t.me/AirGo_Group](#tg群组-httpstmeairgogroup)
<!-- TOC -->

# 面板部分功能展示
<div style="color: darkgray" >Display of panel functions</div>

<table>
<tr>
    <td> <img src="https://github.com/ppoonk/AirGo/raw/main/assets/image/1.png">
    <td> <img src="https://github.com/ppoonk/AirGo/raw/main/assets/image/2.png">
    <td> <img src="https://github.com/ppoonk/AirGo/raw/main/assets/image/3.png">
<tr>
    <td> <img src="https://github.com/ppoonk/AirGo/raw/main/assets/image/4.png">
    <td> <img src="https://github.com/ppoonk/AirGo/raw/main/assets/image/5.png">
    <td> <img src="https://github.com/ppoonk/AirGo/raw/main/assets/image/6.png">
</table>
<table>
<tr>
    <td> <img src="https://github.com/ppoonk/AirGo/raw/main/assets/image/7.png">
</table>

# 1、部署
## 1-1 安装AirGo核心

### 1-1-1 直接安装
- 使用debian，ununtu，centos等系统，执行以下命令，根据提示安装

```
bash <(curl -Ls https://raw.githubusercontent.com/ppoonk/AirGo/main/server/scripts/install.sh)
```


### 1-1-2 使用docker安装


- 在合适的目录新建配置文件，例如：/$PWD/air/config.yaml，配置文件内容如下：

```
system:
  admin-email: admin@oicq.com
  admin-password: adminadmin
  http-port: 80
  https-port: 443
  db-type: sqlite
mysql:
  address: mysql.sql.com
  port: 3306
  config: charset=utf8mb4&parseTime=True&loc=Local
  db-name: imdemo
  username: imdemo
  password: xxxxxx
  max-idle-conns: 10
  max-open-conns: 100
sqlite:
  path: ./air.db

```
- 根据自己的需求，修改配置文件，启动docker命令参考如下：

```
docker run -tid \
  -v $PWD/air/config.yaml:/air/config.yaml \
  -p 80:80 \
  -p 443:443 \
  --name airgo \
  --restart always \
  --privileged=true \
  ppoiuty/airgo:latest
```

docker compose参考如下：
```
version: '3'

services:
  airgo:
    container_name: airgo
    image: ppoiuty/airgo:latest
    ports:
      - "80:80"
      - "443:443"
    restart: "always"
    privileged: true
    volumes:
      - ./config.yaml:/air/config.yaml
```


## 1-2 配置ssl（可选）

- 不需要前后分离的话，进行到这一步就可以了
- 如果使用`宝塔面板`或者`1panel`，请直接自动申请
- 如果您已有证书，只需在安装目录（/usr/local/AirGo/）下，配置 `air.cer`，`air.key`即可


## 1-3 部署前端静态资源（可选，但推荐）

### 1-3-1 部署到Vercel
- fork本项目，修改`项目/web/.env`的`VITE_API_URL`字段为自己的后端地址（由于vercel的限制，请填https接口地址）
- 登录[Vercel](https://vercel.com)，Add New Project，参考下图配置，注意红圈内的设置！
  ![image](https://telegraph-image.pages.dev/file/afe97f45857b988ebd005.png)
- 部署成功后，自定义域名即可（域名解析到76.76.21.21)

### 1-3-2 部署到nginx、caddy等
推荐使用 `github codespaces`编译，这不会在您电脑上安装额外的依赖
- fork本项目，修改`项目/web/.env`的`VITE_API_URL`字段为自己的后端地址
- 在 项目/web/ 下，执行  `npm i && npm run build`
- 打包后的静态资源文件夹为 web，将web文件夹上传到服务器合适位置。新建网站（纯静态），网站位置选择该web文件夹

## 1-4 配置文件说明
```
system:
  admin-email: admin@oicq.com  //管理员账号，初始化之前需要修改！
  admin-password: adminadmin   //管理员密码，初始化之前需要修改！
  http-port: 8899              //核心监听端口
  https-port: 443              //核心监听端口
  db-type: sqlite              //数据库类型，可选值：mysql，mariadb，sqlite
mysql:
  address: xxx.com             //mysql数据库地址
  port: 3306                   //mysql数据库端口
  config: charset=utf8mb4&parseTime=True&loc=Local //保持默认即可
  db-name: xxx                 //mysql数据库名称
  username: xxx                //mysql数据库用户名
  password: xxx                //mysql数据库密码
  max-idle-conns: 10
  max-open-conns: 100
sqlite:
  path: ./air.db               //sqlite数据库文件名
```

## 1-5 启动

注意！如果首次安装，启动核心时，会根据配置文件config.yaml自动初始化数据。务必修改 1-4 中需要修改的部分。

- 启动核心 `systemctl start AirGo`，或者以 docker 方式启动
- 前后端不分离，网站访问地址的端口和配置文件 config.yaml 中的端口保持一致。例如config.yaml中端口为8888，则浏览器需要访问 http://example.com/8888 
- 前后端分离，网站访问地址为vercel或者nginx（caddy）设置的地址


# 2、对接节点

## 2-1 XrayR

```
bash <(curl -Ls https://raw.githubusercontent.com/ppoonk/XrayR-for-AirGo/main/scripts/manage.sh)
```
启动
`systemctl start XrayR`

- docker仓库：[https://hub.docker.com/repository/docker/ppoiuty/xrayr](https://hub.docker.com/repository/docker/ppoiuty/xrayr)

## 2-2 hysteria2
```
bash <(curl -Ls https://raw.githubusercontent.com/ppoonk/shy/main/scripts/install.sh)
```
启动
`systemctl start shy`



# TG群组：[https://t.me/AirGo_Group](https://t.me/AirGo_Group)

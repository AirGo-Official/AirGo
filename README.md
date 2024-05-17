
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

Supported protocols：Vless，Vmess，shadowsocks，Hysteria2


<div style="color: darkgray">AirGo, front and rear end separation, multi user, multi protocol proxy service management system, simple and easy to use.</div>
<div style="color: darkgray">support: vless，vmess，shadowsocks，hysteria2</div>

<hr/>


<hr/>

<!-- TOC -->
* [AirGo 前后端分离，多用户，多协议代理服务管理系统，简单易用](#airgo-前后端分离多用户多协议代理服务管理系统简单易用)
* [面板部分功能展示](#面板部分功能展示)
* [目录：](#目录)
* [1 部署-前后端不分离](#1-部署-前后端不分离)
  * [1-1 直接安装](#1-1-直接安装)
  * [1-2 使用Docker安装](#1-2-使用docker安装)
* [2 部署-前后端分离](#2-部署-前后端分离)
  * [2-1 后端](#2-1-后端)
    * [2-1-1 直接安装](#2-1-1-直接安装)
    * [2-1-2 docker 安装](#2-1-2-docker-安装)
  * [2-2 前端](#2-2-前端)
    * [2-2-1 部署到 Vercel 等云平台](#2-2-1-部署到-vercel-等云平台)
    * [2-2-2 部署到 Nginx、Caddy、OpenResty 等 Web 应用服务器](#2-2-2-部署到-nginxcaddyopenresty-等-web-应用服务器)
* [3 配置ssl（可选）](#3-配置ssl可选)
  * [3-1 给前端设置ssl证书](#3-1-给前端设置ssl证书)
  * [3-2 给后端设置ssl证书](#3-2-给后端设置ssl证书)
* [4 配置文件说明](#4-配置文件说明)
* [5 对接节点](#5-对接节点)
  * [5-1 V2bX](#5-1-v2bx)
    * [5-1-1 直接安装 V2bX](#5-1-1-直接安装-v2bx)
    * [5-1-2 docker 安装 V2bX](#5-1-2-docker-安装-v2bx)
  * [5-2 XrayR](#5-2-xrayr)
    * [5-2-1 直接安装 XrayR](#5-2-1-直接安装-xrayr)
    * [5-2-2 docker 安装 XrayR](#5-2-2-docker-安装-xrayr)
* [6 更新面板](#6-更新面板)
  * [6-1 更新后端](#6-1-更新后端)
  * [6-2 更新前端](#6-2-更新前端)
* [7 命令行](#7-命令行)
* [8 更多说明](#8-更多说明)
<!-- TOC -->

<br>

>TG频道：[https://t.me/Air_Go](https://t.me/Air_Go)
TG群组：[https://t.me/AirGo_Group](https://t.me/AirGo_Group)
文档上次更新日期：2024.4.8

<br>

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

# 目录：

# 1 部署-前后端不分离

## 1-1 直接安装

- 安装核心，使用Ubuntu、Debian、Centos等Linux系统，执行以下命令，然后根据提示安装

```
bash <(curl -Ls https://raw.githubusercontent.com/ppoonk/AirGo/main/server/scripts/install.sh)
```

- 修改配置文件，配置文件目录 `/usr/local/AirGo/config.yaml`，首次安装，会根据配置文件config.yaml自动初始化数据，请务必修改管理员账号和密码
- 启动核心，`systemctl start AirGo`
- 浏览器访问：`http://ip:port`，其中端口为配置文件设定的值

## 1-2 使用Docker安装

- 在合适的目录新建配置文件，例如：/$PWD/air/config.yaml，配置文件内容如下。首次安装，会根据配置文件config.yaml自动初始化数据，请务必修改管理员账号和密码

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

- 启动docker命令参考如下：

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

- docker compose参考如下：

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

- 浏览器访问：`http://ip:port`，其中端口为配置文件设定的值


# 2 部署-前后端分离

## 2-1 后端

### 2-1-1 直接安装
```
bash <(curl -Ls https://raw.githubusercontent.com/ppoonk/AirGo/main/server/scripts/install.sh)
```

- 修改配置文件，配置文件目录 `/usr/local/AirGo/config.yaml`，首次安装，会根据配置文件config.yaml自动初始化数据，请务必修改管理员账号和密码
- 启动核心，`systemctl start AirGo`

### 2-1-2 docker 安装
- 提前准备好配置文件 config.yaml，参考 [config.yaml](https://github.com/ppoonk/AirGo/blob/main/server/config.yaml),首次安装，会根据配置文件config.yaml自动初始化数据，请务必修改管理员账号和密码
- 启动docker命令参考如下：

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

- docker compose参考如下：

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
## 2-2 前端

### 2-2-1 部署到 Vercel 等云平台
- fork本项目，修改`./web/index.html`的`window.httpurl`字段为自己的后端地址，，可以设置多个，以英文符号 `|` 分割。由于vercel的限制，请填https接口地址
- 登录[Vercel](https://vercel.com)，Add New Project，参考下图配置，注意红圈内的设置！
  ![image](https://telegraph-image.pages.dev/file/afe97f45857b988ebd005.png)
- 部署成功后，自定义域名即可（域名解析到76.76.21.21)

### 2-2-2 部署到 Nginx、Caddy、OpenResty 等 Web 应用服务器

- 下载 release 中编译好的静态资源的 `AirGo-web.zip`
- 修改`./web/index.html`的`window.httpurl`字段为自己的后端地址，可以设置多个，以英文符号 `|` 分割
- 在 项目/web/ 下，执行  `npm i && npm run build`
- 打包后的静态资源文件夹为 web，将web文件夹上传到服务器合适位置。新建网站（纯静态），网站位置选择该web文件夹

# 3 配置ssl（可选）

## 3-1 给前端设置ssl证书

通过 `宝塔面板(bt.cn)`，`1panel(1panel.cn)` 等可直接申请、导入证书

## 3-2 给后端设置ssl证书

- 1、通过 `宝塔面板(bt.cn)`，`1panel(1panel.cn)`，先申请或导入证书，再开启反向代理
- 2、如果您已经拥有证书，只需要复制在安装目录（/usr/local/AirGo/）下，将其重命名为 `air.cer`，`air.key`，然后重启 AirGo


# 4 配置文件说明

```
system:
  mode: release                //模式，默认为 release。如果为 dev，即开发模式。控制台会输出更多信息
  admin-email: admin@oicq.com  //管理员账号，初始化之前需要修改！
  admin-password: adminadmin   //管理员密码，初始化之前需要修改！
  http-port: 8899              //核心监听端口
  https-port: 443              //核心监听端口
  db-type: sqlite              //数据库类型，可选值：mysql，mariadb，sqlite
mysql:
  address: xxx.com             //mysql数据库地址
  port: 3306                   //mysql数据库端口
  db-name: xxx                 //mysql数据库名称
  username: xxx                //mysql数据库用户名
  password: xxx                //mysql数据库密码
  config: charset=utf8mb4&parseTime=True&loc=Local //保持默认即可
  max-idle-conns: 10
  max-open-conns: 100
sqlite:
  path: ./air.db               //sqlite数据库文件名
```

# 5 对接节点

**现支持V2bx、XrayR，暂不支持官方版本，请使用下面的版本：**

## 5-1 V2bX

### 5-1-1 直接安装 V2bX

```
bash <(curl -Ls https://raw.githubusercontent.com/ppoonk/V2bX/main/scripts/install.sh)
```
- 安装完成后请根据需要在```/etc/V2bX/config.json```中修改配置文件
- 启动：使用管理脚本```AV```或直接 `systemctl start AV`

### 5-1-2 docker 安装 V2bX
- 提前准备好配置文件 config.json，参考 [config.json](https://github.com/ppoonk/V2bX/blob/main/config.json)
- 启动docker命令参考如下：

```
docker run -tid \
  -v $PWD/av/config.json:/etc/V2bX/config.json \
  --name av \
  --restart always \
  --net=host \
  --privileged=true \
  ppoiuty/av:latest
```

- docker compose参考如下：

```
version: '3'
services:
  AV:
    container_name: AV
    image: ppoiuty/av:latest
    network_mode: "host"
    restart: "always"
    privileged: true
    volumes:
      - ./config.json:/etc/V2bX/config.json
```

## 5-2 XrayR



### 5-2-1 直接安装 XrayR

```
bash <(curl -Ls https://raw.githubusercontent.com/ppoonk/XrayR-for-AirGo/main/scripts/manage.sh)
```

- 安装完成后请根据需要在```/usr/local/XrayR/config.yml```中修改配置文件
- 启动：使用管理脚本```XrayR```或直接 `systemctl start XrayR`

### 5-2-2 docker 安装 XrayR

- 提前准备好配置文件 config.yml，参考 [config.yml](https://github.com/ppoonk/XrayR-for-AirGo/blob/main/config.yml)

- 启动docker命令参考如下：

```
docker run -tid \
  -v $PWD/xrayr/config.yml:/etc/XrayR/config.yml \
  --name xrayr \
  --restart always \
  --net=host \
  --privileged=true \
  ppoiuty/xrayr:latest
```

- docker compose参考如下：

```
version: '3'
services:
  xrayr:
    container_name: xrayr
    image: ppoiuty/xrayr:latest
    network_mode: "host"
    restart: "always"
    privileged: true
    volumes:
      - ./config.yml:/etc/XrayR/config.yml
```

# 6 更新面板
更新时，请检查 `前端版本` 和 `后端核心版本`，它们处在不同位置并且版本号保持一致，如图：
![](https://github.com/ppoonk/AirGo/raw/main/assets/image/8.png)

## 6-1 更新后端

- 方式1: 下载新的二进制文件，替换旧的，然后执行 ./AirGo update 完成更新
- 方式2: 在版本 `v0.2.5`之后，通过`面板-管理员-系统`，可以点击 `升级按钮`完成更新
- 说明：更新核心后，角色绑定的菜单和casbin权限(api权限)会设置为当前核心的默认值


## 6-2 更新前端

按照 [2-1 前端](#2-1-前端)重新部署即可


# 7 命令行

```
./AirGo help                    获取帮助
./AirGo reset --resetAdmin      重置admin password
./AirGo start                   启动AirGo, 指定配置文件路径：./AirGo start --config path2/config.yaml
./AirGo update                  更新数据库相关AirGo数据
./AirGo version                 查看AirGo的当前版本
```

# 8 更多说明

[点击查看更多](https://github.com/ppoonk/AirGo/wiki/Wiki)


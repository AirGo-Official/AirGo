#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
plain='\033[0m'

[[ $EUID -ne 0 ]] && ${red}"非root权限\n"${plain} && exit 1

arch=$(uname -m)
system=$(uname)
country=''
latestVersion=''
downloadPrefix='https://github.com/ppoonk/AirGo/releases/download/'
githubApi="https://api.github.com/repos/ppoonk/AirGo/releases/latest"
installScript="https://raw.githubusercontent.com/ppoonk/AirGo/v2/server/scripts/install.sh"
acmeGit="https://github.com/acmesh-official/acme.sh.git"
yamlFile="/usr/local/AirGo/config.yaml"

ipv4=""
#ipv6=""
ipv4_local=""


get_system_type(){
if [ "$system" == "Darwin" ]; then
  system="darwin-10.14"
else
  system="linux"
fi
}
get_arch(){
  if [[ $arch == "x86_64" || $arch == "x64" ]]; then
      arch="amd64"
  elif [[ $arch == "aarch64" || $arch == "arm64" || $arch == "armv8" || $arch == "armv8l" ]]; then
      arch="arm64"
  elif [[ $arch == "arm"  || $arch == "armv7" || $arch == "armv7l" || $arch == "armv6" ]];then
      arch="arm-7"
  else
      echo -e ${red}"不支持的arch，请自行编译\n"${plain}
      exit 1
  fi
}
get_region() {
    country=$( curl -4 "https://ipinfo.io/country" 2> /dev/null )
    if [ "$country" == "CN" ]; then
      acmeGit="https://gitee.com/neilpang/acme.sh.git"
      downloadPrefix="https://ghproxy.com/${downloadPrefix}"
      installScript="https://ghproxy.com/${installScript}"
    fi
}
open_ports(){
	systemctl stop firewalld.service >/dev/null 2>&1
	systemctl disable firewalld.service >/dev/null 2>&1
	setenforce 0 >/dev/null 2>&1
	ufw disable >/dev/null 2>&1
	iptables -P INPUT ACCEPT
	iptables -P FORWARD ACCEPT
	iptables -P OUTPUT ACCEPT
	iptables -t nat -F
	iptables -t mangle -F
	iptables -F
	iptables -X
}

set_dependences() {
    if [[ $(command -v yum) ]]; then
      if [[ ! $(command -v wget) ]] || [[ ! $(command -v curl) ]] || [[ ! $(command -v git) ]] || [[ ! $(command -v socat) ]] || [[ ! $(command -v unzip) ]] || [[ ! $(command -v gawk) ]] || [[ ! $(command -v lsof) ]]; then
          echo -e ${green}"安装依赖\n"${plain}
          yum update -y
          yum install wget curl git socat unzip gawk lsof -y
      fi
    elif [[ $(command -v apt) ]]; then
      if [[ ! $(command -v wget) ]] || [[ ! $(command -v curl) ]] || [[ ! $(command -v git) ]] || [[ ! $(command -v socat) ]] || [[ ! $(command -v unzip) ]] || [[ ! $(command -v gawk) ]] || [[ ! $(command -v lsof) ]]; then
          echo -e ${green}"安装依赖\n"${plain}
          apt update -y
          apt install wget curl git socat unzip gawk lsof -y
      fi
       echo -e "依赖已安装\n"
    fi
}
get_latest_version() {
          latestVersion=$(curl -Ls $githubApi | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
          if [[ ! -n "$latestVersion" ]]; then
              echo -e "${red}获取最新版本失败，请稍后重试${plain}"
              exit 1
          fi
}
initialize(){
  get_arch
  get_system_type
  set_dependences
  get_region
  get_latest_version

 ipv4=$(curl -4 -s --max-time 5 http://icanhazip.com/ || '你的ip' )
 #ipv6=$(curl -6 -s --max-time 5 http://icanhazip.com/)
 ipv4_local=$( ip addr | awk '/^[0-9]+: / {}; /inet.*global.*eth/ {print gensub(/(.*)\/(.*)/, "\\1", "g", $2)}' || '你的内网ip')

}

confirm_msg() {
    if [[ $# -gt 1 ]]; then
        echo && read -p "$1 [默认$2]: " temp
        if [[ "${temp}x" == ""x ]]; then
            temp=$2
        fi
    else
        read -p "$1 [y/n]: " temp
    fi
    if [[ "${temp}"x == "y"x || x"${temp}" == x"Y" ]]; then
        return 0
    else
        return 1
    fi
}

read_yaml(){
    cat $1 | while read LINE
    do
      if [ "$(echo $LINE | grep $2)" != "" ];then
       #        echo -e "return：$(echo $LINE | grep $2)"
         text=$( echo "$LINE" | awk -F ":" '{print $2}' )
         eval echo $text
      fi
    done
}
installation_status(){
      if [[ ! -f /etc/systemd/system/$1.service ]] || [[ ! -f /usr/local/$1/$1 ]]; then
        return 1
      else
        return 0
      fi
}
run_status() {
      temp=$(systemctl is-active $1)
      if [[ x"${temp}" == x"active" ]]; then
          return 0
      else
          count=$(ps -ef | grep "$1" | grep -v "grep" | wc -l)
          if [[ count -ne 0 ]]; then
              return 0
          else
              return 1
          fi
      fi
}


download(){
  echo -e "开始下载核心，版本：${latestVersion}"
  rm -rf /usr/local/AirGo
  mkdir /usr/local/AirGo

  wget -N --no-check-certificate -O /usr/bin/AirGo ${installScript}
  chmod 777 /usr/bin/AirGo

  wget -N --no-check-certificate -O /usr/local/AirGo/AirGo.zip ${downloadPrefix}${latestVersion}/AirGo-${latestVersion}-${system}-${arch}.zip
  if [[ $? -ne 0 ]]; then
      echo -e "${red}下载失败，请重试${plain}"
      exit 1
  fi
  echo -e "开始解压..."
  cd /usr/local/AirGo/
  unzip AirGo.zip
  chmod 777 -R /usr/local/AirGo
  mv /usr/local/AirGo/AirGo-${latestVersion}-${system}-${arch} /usr/local/AirGo/AirGo

}
add_service(){
  cat >/etc/systemd/system/$1.service <<-EOF
  [Unit]
  Description=$1 Service
  After=network.target
  Wants=network.target
  StartLimitIntervalSec=0

  [Service]
  Restart=always
  RestartSec=1
  Type=simple
  WorkingDirectory=/usr/local/$1/
  ExecStart=/usr/local/$1/$1 -start

  [Install]
  WantedBy=multi-user.target
EOF

}
install(){
  installation_status 'AirGo'
  if [[ $? -eq 0 ]]; then
      echo -e "${red}AirGo已安装${plain}"
      echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
      main
  fi

  run_status "AirGo"
  if [[ $? -eq 0 ]]; then
   echo -e "${red}AirGo正在运行${plain}"
   echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
   main
  fi

  download
  add_service "AirGo"
  systemctl daemon-reload
  systemctl enable AirGo
#  systemctl start AirGo

  echo -e "${green}安装完成，版本：${latestVersion}${plain}"

  echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
  main
}
uninstall(){
  confirm_msg "确定要卸载吗?" "n"
      if [[ $? != 0 ]]; then
          return 0
      fi
  echo -e "开始卸载"
      systemctl stop $1
      systemctl disable $1
      systemctl daemon-reload
      systemctl reset-failed
      rm -rf /etc/systemd/system/$1.service /usr/local/$1

  echo -e "${green}卸载完成${plain}"
  echo
  echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
  main
}

start(){

  installation_status 'AirGo'
  if [[ $? -eq 1 ]]; then
      echo -e "${red}AirGo未安装${plain}"
      echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
      main
  fi

  run_status "AirGo"
  if [[ $? -eq 0 ]]; then
   echo -e "${red}AirGo正在运行${plain}"
   echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
   main
  fi

  httpPort=$(read_yaml $yamlFile "http-port")
  name=$(lsof -i:$httpPort | awk '{print $1}')
  if [[ ! $name == "AirGo" && ！$name == "" ]]; then
    echo -e "${red}端口被占用，请检查端口，或者修改 /usr/local/AirGo/config.yaml 配置文件${plain}"
    echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
    main
  fi

  systemctl start AirGo
  systemctl is-active AirGo

  echo -e "${green}操作完成${plain}"
  echo -e "${green}公网访问：${ipv4}:${httpPort}${plain}"
  echo -e "${green}内网访问：${ipv4_local}:${httpPort}${plain}"
  echo -e "${yellow}个别vps上述地址无法访问请更换端口${plain}"
  echo
  echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
  main
}
stop(){
    installation_status 'AirGo'
    if [[ $? -eq 1 ]]; then
        echo -e "${red}AirGo未安装${plain}"
        echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
        main
    fi

    run_status "AirGo"
    if [[ $? -eq 1 ]]; then
     echo -e "${red}AirGo未运行${plain}"
     echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
     main
    fi

  systemctl stop AirGo
  systemctl is-active AirGo
  echo -e "${green}操作完成${plain}"
  echo
  echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
  main
}

reset_admin(){
  installation_status 'AirGo'
  if [[ $? -eq 1 ]]; then
      echo -e "${red}AirGo未安装${plain}"
      echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
      main
  fi

  cd /usr/local/AirGo/
  ./AirGo -resetAdmin
  echo -e "${green}完成${plain}"
  echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
  main
}

acme(){
  installation_status 'AirGo'
  if [[ $? -eq 1 ]]; then
      echo -e "${red}AirGo未安装${plain}"
      echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
      main
  fi

  installation_status "AirGo"
  if [[ $? -ne 0 ]]; then
   echo -e "${red}AirGo未安装,脚本退出${plain}"
   exit 1
  fi
  cd /usr/local/AirGo


  if [[ ! -f /usr/local/AirGo/acme.sh/acme.sh ]];then
    git clone ${acmeGit}
    chmod 777 -R acme.sh
  fi
  cd acme.sh

  email=''
  domain=''

  echo -e "${yellow}设置Acme邮箱:${plain}"
  read -p "输入您的邮箱:" email
  echo -e "您的邮箱:${email}"

  echo -e "${yellow}设置域名:${plain}"
  read -p "输入您的域名:" domain
  echo -e "您的域名:${domain}"
  domainPrefix=$(echo ${domain%%.*})

  echo -e "${yellow}配置邮箱账户...${plain}"
  ./acme.sh --install -m ${email}
  echo -e "${yellow}正在发起 DNS 申请...${plain}"
  ./acme.sh --issue --dns -d ${domain} --yes-I-know-dns-manual-mode-enough-go-ahead-please

  echo -e "${yellow}请仔细查看命令行显示文本中，有无以下字段：${plain}"
  echo -e "[Tue Sep 12 12:30:59 UTC 2023] TXT value: '**************************************-****"

  echo -e "${yellow}如果存在该字段，请去你的域名 DNS 管理商，完成下面2个重要操作！！！${plain}"
  echo -e "${yellow}1、${plain}添加一个txt记录"
  echo -e "${yellow}2、${plain}将该记录的[名称]设置为：_acme-challenge.${domainPrefix}，完整域名为 _acme-challenge.${domain}"
  echo

  confirm_msg "是否已经添加这条 txt 记录？是否将该记录的[名称]设置为：_acme-challenge.${domainPrefix}？ "
   if [[ $? -ne 0 ]]; then
     echo -e "${red}未添加txt 记录,脚本退出${plain}"
     exit 1
   fi

  echo -e "${green}添加 txt 记录成功，进行下一步${plain}"
  echo -e "${green}开始申请证书...${plain}"

  ./acme.sh --renew -d ${domain} --yes-I-know-dns-manual-mode-enough-go-ahead-please
    if [ $? -ne 0 ]; then
        echo -e "${red}申请失败,脚本退出${plain}"
        exit 1
    fi
  echo -e "${green}申请成功,证书文件在/root/.acme.sh/${domain}文件夹下${plain}"
  echo -e "${green}正在将证书复制到/usr/local/AirGo/${plain}"

  cp /root/.acme.sh/${domain}*/${domain}.cer /usr/local/AirGo/air.cer
  cp /root/.acme.sh/${domain}*/${domain}.key /usr/local/AirGo/air.key

  systemctl stop AirGo
  systemctl start AirGo

  httpPort=$(read_yaml $yamlFile "http-port")
  httpsPort=$(read_yaml $yamlFile "https-port")

  echo -e "${green}操作完成，服务已重启${plain}"
  echo -e "${green}http公网访问：${ipv4}:${domain}${httpPort}${plain}"
  echo -e "${green}https公网访问：${ipv4}:${domain}${httpsPort}${plain}"
  echo -e "${green}内网访问：${ipv4_local}:${httpPort}${plain}"
  echo
  echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
  main

}
main(){
  clear
  installationStatus='未安装'
  runStatus='未运行'
  installation_status 'AirGo'
  if [[ $? -eq 0 ]]; then
    installationStatus='已安装'
  fi
  run_status 'AirGo'
    if [[ $? -eq 0 ]]; then
      runStatus='已运行'
    fi

  echo -e "
  ${green}AirGo-panel 管理脚本${plain}

  状态： ${green}${installationStatus}${plain}    ${green}${runStatus}${plain}
  ${yellow}-------------------------${plain}
  ${green}1.${plain} 安装
  ${green}2.${plain} 卸载
  -${yellow}------------------------${plain}
  ${green}3.${plain} 启动
  ${green}4.${plain} 停止
  ${yellow}-------------------------${plain}
  ${green}5.${plain} 重置管理员密码
  ${yellow}-------------------------${plain}
  ${green}6.${plain} 使用Acme脚本申请ssl证书
  (dns手动模式，无80和443端口也可申请证书)
  ${yellow}-------------------------${plain}
  ${green}7.${plain} 开放所有端口
  ${yellow}-------------------------${plain}
  ${green}0.${plain} 退出
 "

  echo && read -p "请输入序号: " tem
  case "${tem}" in
  1) install;;
  2) uninstall "AirGo";;
  3) start;;
  4) stop;;
  5) reset_admin;;
  6) acme;;
  7) open_ports;;
  *) exit 0;;

  esac

}
initialize
main
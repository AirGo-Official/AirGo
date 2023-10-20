
cd ../../web
npm i
# 前端打包
npm run build
# 将打包文件移动到后端，嵌入到go编译的二进制文件中
rm -rf ../server/web/web
mv web ../server/web/
cd ../server

# 本机编译
CGO_ENABLED=1 go build -o AirGo -ldflags='-s -w --extldflags "-static -fpic"' main.go


# ubuntu/debian 交叉编译 arm64
#sudo apt update -y
#sudo apt install gcc-aarch64-linux-gnu -y
#CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -o AirGo -ldflags='-s -w --extldflags "-static -fpic"' main.go
#

# ubuntu/debian 交叉编译 arm
#sudo apt update -y
#sudo apt install gcc-arm-linux-gnueabihf  -y
#CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc go build -o AirGo -ldflags='-s -w --extldflags "-static -fpic"' main.go


# 使用 xgo 编译多个平台 文档地址：https://github.com/techknowlogick/xgo
# 1 curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun
# 2 docker pull techknowlogick/xgo:latest
# 3 进入项目server目录
#    go >= 1.17
#    go install src.techknowlogick.com/xgo@latest
# 4 默认编译全部平台
#   xgo -ldflags '-s -w' -out AirGo .
#   编译指定平台
#   xgo --targets=linux/arm64,linux/arm-7,windows-8.1/amd64,darwin-10.14/* -ldflags '-s -w' -out AirGo-2.1.1 .

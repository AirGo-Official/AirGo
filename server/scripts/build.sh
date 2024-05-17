version='v0.2.8'

system=$(uname)

if [[ $system == "Darwin" ]]; then
    sed -i "" 's/old-version/'$version'/g' ../../web/src/layout/footer/index.vue
    sed -i "" 's/old-version/'$version'/g' ../constant/index.go
elif [[ $system == "Linux" ]]; then
    sed -i 's/old-version/${version}/g' ../../web/src/layout/footer/index.vue
    sed -i 's/old-version/${version}/g' ../server/constant/index.go
fi


#npm cache clean -f
#npm install npm -g
#npm install n -g
#n v20.9.0
#


## 前端打包
cd ../../web || exit
npm i
npm run build

## 将打包文件移动到后端，嵌入到go编译的二进制文件中
rm -rf ../server/web/web
mv web ../server/web/

## 后端端打包
cd ../server || exit
## 本机编译
##CGO_ENABLED=0 go build -o AirGo -trimpath -ldflags='-s -w'
#
## amd64 linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o AirGo -trimpath -ldflags='-s -w'
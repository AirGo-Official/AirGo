containerName=$( docker ps -a | awk 'NR==2{print $1}' )
imageName=$( docker image ls | awk 'NR==2{print $3}' )
docker stop $containerName
docker rm $containerName
docker rmi $imageName


docker run -ti \
  -v /$PWD/air/config.yaml:/air/config.yaml \
  -p 10081:80 \
  -p 10082:443 \
  --name airgo \
  --privileged=true \
  ppoiuty/airgo:latest


docker run -ti \
  -v /$PWD/air/config.yaml:/air/config.yaml \
  -p 10010:80 \
  --name airgo \
  --privileged=true \
  airgo:latest
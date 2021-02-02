docker build -t build-img .
docker create --name build-cont build-img
docker cp build-cont:/app/agent ./agent

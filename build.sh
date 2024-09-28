docker --debug build -f ./container/Dockerfile . -t api
# docker push "kotohiro:latest"
docker run --env-file ./.env -p 3000:3000 api

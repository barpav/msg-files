up:
	sudo docker-compose up -d --wait
down:
	sudo docker-compose down

build:
	sudo docker image rm -f ghcr.io/barpav/msg-files:v1
	sudo docker build -t ghcr.io/barpav/msg-files:v1 -f docker/service/Dockerfile .
	sudo docker image ls
clear:
	sudo docker image rm -f ghcr.io/barpav/msg-files:v1

up-debug:
	sudo docker-compose -f compose-debug.yaml up -d --wait
down-debug:
	sudo docker-compose -f compose-debug.yaml down

user:
	curl -v -X POST	-H "Content-Type: application/vnd.newUser.v1+json" \
	-d '{"id": "jane", "name": "Jane Doe", "password": "My1stGoodPassword"}' \
	localhost:8081
session:
	curl -v -X POST -H "Authorization: Basic amFuZTpNeTFzdEdvb2RQYXNzd29yZA==" localhost:8082
# make file KEY=session-key
file:
	curl -v -X POST -H "Authorization: Bearer $(KEY)" \
	-H "Content-Type: application/vnd.newPrivateFile.v1+json" \
	-d '{"name": "test.jpg", "mime": "image/jpeg", "access": ["john", "alice", "bob"]}' \
	localhost:8080
# make file-pub KEY=session-key
file-pub:
	curl -v -X POST -H "Authorization: Bearer $(KEY)" \
	-H "Content-Type: application/vnd.newPublicFile.v1+json" \
	-d '{"name": "test.jpg", "mime": "image/jpeg"}' \
	localhost:8080

push:
	sudo docker push ghcr.io/barpav/msg-files:v1

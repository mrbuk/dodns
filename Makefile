version = 0.3
user = mrbuk
project = dodns

all: docker-build docker-save
.PHONY : all

docker-build:
	docker build . -t ${user}/${project}:$(version) -t ${user}/${project}:latest

docker-push: docker-build
	docker push ${user}/${project}:$(version)
	docker push ${user}/${project}:latest

docker-save:
	mkdir -p images
	docker save ${user}/${project}:latest -o images/${user}_${project}_latest.tgz

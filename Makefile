.PHONY:docker_build_processor docker_build_server docker_compose_run

docker_build_processor:
	docker build -t processor -f build/processor.Dockerfile .

docker_build_server:
	docker build -t server -f build/server.Dockerfile .

docker_compose_run:
	docker-compose -f build/docker-compose.yaml up

docker_compose_build:
	docker-compose -f build/docker-compose.yaml up --build


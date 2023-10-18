
build-all:
	cd checkout && make build
	cd loms && make build
	cd notifications && make build

run-all: build-all
	docker compose up --force-recreate --build

run-services: build-all
	docker compose up -d loms-postgres checkout-postgres loms checkout notifications

run-kafka:
	docker-compose up -d zookeeper kafka1 kafka2 kafka3

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit

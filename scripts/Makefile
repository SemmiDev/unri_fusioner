run:
	go run ../cmd/*.go

build-docker:
	docker build -t unri_fusioner -f ../Dockerfile ../

run-compose:
	docker-compose -f ../docker-compose.yaml up -d

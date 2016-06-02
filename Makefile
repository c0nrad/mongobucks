run: build env
	./mongobucks

build:
	go build

docker-run: 
	docker kill `docker ps -q`
	docker rm `docker ps -a -q`

	docker build -t mongobucks .
	docker run  --env-file=".ENV" -p 8080:8080 -t mongobucks

env:
	export $(cat .ENV | xargs)
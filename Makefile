.PHONY:
docker-build:
	docker build . -t gcr.io/small-biz-template/markenshop/mini-score:latest

.PHONY:
docker-run:
	docker run --rm -p 8081:8080 gcr.io/small-biz-template/markenshop/mini-score:latest



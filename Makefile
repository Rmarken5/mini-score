.Phony:
docker-build-service:
	docker build -f ./dockerfiles/service.Dockerfile . -t gcr.io/small-biz-template/markenshop/mini-score:latest

.Phony:
docker-run-service:
	docker run --rm -p 8081:8080 -e POSTGRES_USER=user \
                                 -e POSTGRES_SSL_MODE=disable \
                                 -e POSTGRES_DATABASE=postgres \
                                 -e POSTGRES_HOST=host.docker.internal \
                                 -e POSTGRES_PORT=5432 \
                                 -e POSTGRES_PASSWORD=password \
                                 gcr.io/small-biz-template/markenshop/mini-score:latest

.Phony:
docker-build-database:
	docker build -f ./dockerfiles/database.Dockerfile . -t gcr.io/small-biz-template/markenshop/mini-score-db:latest

.Phony:
docker-run-database:
	docker run --rm -d -p 5432:5432 -v mini-score-db:/var/lib/postgresql/data gcr.io/small-biz-template/markenshop/mini-score-db:latest

.Phony:
docker-build-nflscheduler:
	docker build -f ./dockerfiles/nflscheduler.Dockerfile . -t gcr.io/small-biz-template/markenshop/nflscheduler:latest

.Phony:
docker-run-nflscheduler:
	docker run --rm -e POSTGRES_USER=user \
                               -e POSTGRES_SSL_MODE=disable \
                               -e POSTGRES_DATABASE=postgres \
                               -e POSTGRES_HOST=host.docker.internal \
                               -e POSTGRES_PORT=5432 \
                               -e POSTGRES_PASSWORD=password \
                                gcr.io/small-biz-template/markenshop/nflscheduler:latest

.Phony:
migrate-up:
	migrate -path ./service/internal/nfl/data-access/db/migrations -database "${NFL_CONNECTION_STRING}" up

.Phony:
migrate-down:
	migrate -path ./service/internal/nfl/data-access/db/migrations -database "${NFL_CONNECTION_STRING}" down

.Phony:
test:
	go test ./...

.PHONY: gen
gen:
	find . -name "*_mock.go" -type f -delete
	go generate ./...



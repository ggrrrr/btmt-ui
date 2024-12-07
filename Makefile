DOCKER_REPO ?= "local"
GIT_HASH ?= $(shell git log --format="%h" -n 1)

build_svc_%:
	docker build --build-arg=MAIN_FILE=svc-$*/cmd/main.go \
		--build-arg=GIT_HASH=${GIT_HASH} \
		-f ./be/.deploy/Dockerfile \
		--tag ${DOCKER_REPO}/be/$*:${GIT_HASH} \
		./be

kind_image_svc_%:
	kind load docker-image ${DOCKER_REPO}/be/$*:${GIT_HASH}

tag_svc_%:
	docker tag ${DOCKER_REPO}/be/$*:${GIT_HASH} ${DOCKER_REPO}/be/$*:latest

build_web:
	docker build \
		-f ./ui/web/.deploy/Dockerfile \
		--tag "${DOCKER_REPO}/ui/web:${GIT_HASH}" \
		./ui/web

build_tag_be: build_tag_auth build_tag_people

build_tag_auth: build_svc_auth tag_svc_auth

build_tag_people: build_svc_people tag_svc_people

tag_web:
	docker tag "${DOCKER_REPO}/ui/web:${GIT_HASH}" "${DOCKER_REPO}/ui/web:latest"

clean_docker:
	docker rm $(docker ps -a -q)
	docker image prune -a -f
	docker volume prune -a -f

go_clean:
	# go clean -cache
	go clean -testcache
	# go clean -fuzzcache
	# go clean -modcache

go_run_monolith:
	docker compose up -d localstack mongo nats otel postgres
	go run be/monolith/main.go

go_test:
	docker-compose start postgres mongo localstack nats
	go test -cover  ./be/...
	# go test -coverprofile cover.out
	# go test -cover -coverprofile cover.out ./be/...

lint:
	golangci-lint run -v be/...


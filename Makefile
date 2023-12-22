DOCKER_REPO ?= "local"
GIT_HASH ?= $(shell git log --format="%h" -n 1)

build_svc_%:
	docker build --build-arg=MAIN_FILE=svc-$*/cmd/main.go \
		-f ./be/.deploy/Dockerfile \
		--tag "${DOCKER_REPO}/be/$*:${GIT_HASH}" \
		./be

tag_svc_%:
	docker tag "${DOCKER_REPO}/be/$*:${GIT_HASH}" "${DOCKER_REPO}/be/$*:latest"

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
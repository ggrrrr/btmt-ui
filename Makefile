DOCKER_REPO ?= "local"
GIT_HASH ?= $(shell git log --format="%h" -n 1)

PROTO_DEPS = "./"
PROTO_DST = "./"

# TODO add tooling to fetch latest version
PROTO_EXT_DEPS="./proto"

install_protoc_go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

	# github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    # github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    # google.golang.org/protobuf/cmd/protoc-gen-go \
    # google.golang.org/grpc/cmd/protoc-gen-go-grpc



protoc_common:
	protoc -I ./ -I=${PROTO_EXT_DEPS} \
		-I=./be/common/templ/v1 \
		--go_opt=paths=source_relative \
		--go_out="./" \
		./be/common/templ/v1/templ.proto
	protoc -I ./ -I=${PROTO_EXT_DEPS} \
		-I=./be/common/templ/v1 \
		--go_opt=paths=source_relative \
		--go_out="./" \
		./be/common/msgbus/v1/msg.proto
	protoc -I ./ -I=${PROTO_EXT_DEPS} \
		-I=./be/common/templ/v1 \
		--go_opt=paths=source_relative \
		--go_out="./" \
		./be/common/msgbus/v1/test.proto

protoc_svc:
	protoc -I ./ -I=${PROTO_EXT_DEPS} \
		-I=./be/common \
		--go_opt=paths=source_relative \
		--go_out="./" \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_out="./" \
		--grpc-gateway_out="./" \
		--grpc-gateway_opt=paths=source_relative \
		./be/svc-auth/authpb/v1/passwd.proto
	protoc -I ./ -I=${PROTO_EXT_DEPS} \
		-I=./be/common \
		--go_opt=paths=source_relative \
		--go_out="./" \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_out="./" \
		./be/svc-email/emailpb/v1/sender.proto
	protoc -I ./ -I=${PROTO_EXT_DEPS} \
		-I=./be/common \
		--go_opt=paths=source_relative \
		--go_out="./" \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_out="./" \
		--grpc-gateway_out="./" \
		--grpc-gateway_opt=paths=source_relative \
		./be/svc-people/peoplepb/v1/people.proto
	protoc -I ./ -I=${PROTO_EXT_DEPS} \
		-I=./be/common \
		--go_opt=paths=source_relative \
		--go_out="./" \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_out="./" \
		./be/svc-tmpl/tmplpb/v1/templates.proto

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

go_mod_upgrade:
	cd be; go get -u ./...

go_clean:
	go clean -cache
	go clean -testcache
	# go clean -fuzzcache
	go clean -modcache

go_run_monolith:
	docker compose up -d localstack mongo nats otel postgres
	go run be/monolith/main.go

go_test:
	docker compose up -d localstack mongo nats otel postgres
	go test -cover  ./be/...
	# go test -coverprofile cover.out
	# go test -cover -coverprofile cover.out ./be/...

lint:
	golangci-lint run -v be/...


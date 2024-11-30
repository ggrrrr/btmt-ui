# Install

## Some golang 

* [uber guidelines](https://github.com/uber-go/guide/blob/master/style.md)
* [ardanlabs DDD](https://github.com/ardanlabs/service6-video/tree/main)

## NATS

```sh
brew tap nats-io/nats-tools
brew install nats-io/nats-tools/nats

# add stream



```

## GoTemplate 

https://gotemplate.io/

## gRPC

```sh
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

- buf generator <https://buf.build/docs/installation/>

## Docs

- gRPC Gateway <https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/adding_annotations/>
  - <https://github.com/grpc-ecosystem/grpc-gateway#installation>
- gRPC auth with JWT <https://dev.to/techschoolguru/use-grpc-interceptor-for-authorization-with-jwt-1c5h>

- Rule Engine or DSL or HCL ?

## HTTP standards

- <https://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml>

## LOCALSTACK

### S3

add the following to your /etc/hosts

```sh
127.0.0.1       test-bucket-1.localhost

```

### Create table in AWS dynamodb

```
awslocal dynamodb create-table \
   --table-name module-auth \
   --attribute-definitions \
    AttributeName=email,AttributeType=S \
   --key-schema AttributeName=email,KeyType=HASH \
   --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

```

## List record in table AWS dynamodb

```
awslocal dynamodb list-tables
awslocal dynamodb scan --table-name module-auth
```

```
# Create/update admin user

go get -u ./...


go run svc-auth/cmd/main.go admin --email asd@asd -p asdasd
# list users
go run svc-auth/cmd/main.go list

# Login with user/pass
go run svc-auth/cmd/main.go client --email asd@asd -p asdasd


curl -v -XPOST -d'{"email":"asd@asd","password":"asdasd"}' \
  ${REST_URL}/v1/auth/login/passwd | jq -r '.payload.token' | pbcopy

export T="<CTRL+V>"

curl -v -H"Authorization: Bearer $T" -XPOST  ${REST_URL}/rest/v1/auth/validate

curl -v -H"Authorization: Bearer $T" -XPOST \
  -d'{"password":"1asdasdasd","new_password":"asd"}' \
  ${REST_URL}/v1/auth/update/passwd

curl -v -H"Authorization: Bearer $T" -XPOST \
   -d'{"data":{"id_numbers":{"EGN":"asdasdasd"},"emails":{"main":"asd@asd123"},"full_name":"Varban Krushev","name":"vesko","phones":{"mobile":"0889430425"},"labels":["hike:snow"]}}' \
  ${REST_URL}/rest/v1/people/save

curl -v -H"Authorization: Bearer $T" -XPOST -d"{}" ${REST_URL}/v1/people/list

curl -v -H"Authorization: Bearer $T" -XPOST -d'{"filters":{"phones":{"list":["0889430425"]}}}' \
  ${REST_URL}/rest/v1/people/list | jq


```

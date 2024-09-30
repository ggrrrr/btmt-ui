#!/bin/bash
set -x
# awslocal dynamodb create-table \
#    --table-name module-auth \
#    --attribute-definitions \
#     AttributeName=email,AttributeType=S \
#    --key-schema AttributeName=email,KeyType=HASH \
#    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
awslocal s3 mb s3://test-bucket-1

set +x

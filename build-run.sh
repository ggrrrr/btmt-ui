#!/bin/bash

pwd

ls -al src/be
cd src/be
go mod download
go run $1 server
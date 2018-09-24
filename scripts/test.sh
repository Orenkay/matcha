#!/bin/bash

docker-compose -f $1 build
docker-compose -f $1 up -d postgres
sleep 2
docker-compose -f $1 run --rm matcha
docker-compose -f $1 down
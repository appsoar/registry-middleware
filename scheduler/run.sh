#!/bin/bash


docker run -v /proc:/host/proc -e DBURL="http://192.168.12.112:8080" -p 9000:9090 cloudsoar/registry-scheduler:4.0

#!/bin/bash


docker run -v /proc:/host/proc -v /:/.hidden/root:ro -e DBURL="http://192.168.15.86:9182" -p 9090:9090 -e DEBUG=true cloudsoar/registry-scheduler:2.0

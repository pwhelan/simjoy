#!/bin/bash

# compile windows first
env GOOS=windows GOARCH=amd64 make
mv bin/simjoy-windows-amd64.exe bin/simjoy.exe

make

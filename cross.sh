#!/bin/bash

env GOOS=windows GOARCH=amd64 make
env GOOS=linux GOARCH=amd64 make

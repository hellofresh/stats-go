#!/bin/sh
CWD=$(pwd)

# Goes to the application source code
cd stats-go-master

# Creates the project source on the GOPATH
mkdir -p ${PROJECT_SRC}

# Copies the current source code from the app to the GOPATH
cp -r . ${PROJECT_SRC}

# Goes to the application on the GOPATH
cd ${PROJECT_SRC}

# Build the go application (run tests actually)
# make

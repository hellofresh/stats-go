#!/bin/sh

mkdir -p ${PROJECT_SRC}
cp -r . ${PROJECT_SRC}
cd ${PROJECT_SRC}

make

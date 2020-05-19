#!/bin/sh

VERSION=$(git describe --tags `git rev-list --tags --max-count=1`)
echo $VERSION

COMMITID=$(git rev-parse --short HEAD)
echo $COMMITID

BUILDTIME=$(date '+%Y-%m-%d %H:%M:%S')
echo $BUILDTIME

CURRENT=$PWD

WORKSPACE=$(dirname $(readlink -f $0))

cd $WORKSPACE

rm -f taiji

go build -o taiji -ldflags "-X main.Version=$VERSION -X main.CommitId=$COMMITID -X 'main.BuildTime=$BUILDTIME'"


cd $CURRENT


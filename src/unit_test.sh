#!/bin/sh

WORKSPACE=$(dirname $(readlink -f $0))
echo $WORKSPACE

cd $WORKSPACE/model/account
go test -v
if [ $? -ne 0 ]
then
	exit $?
fi

cd $WORKSPACE/model/process
go test -v
if [ $? -ne 0 ]
then
	exit $?
fi

cd $WORKSPACE/model/token
go test -v
if [ $? -ne 0 ]
then
	exit $?
fi

cd $WORKSPACE/controller/account
go test -v
if [ $? -ne 0 ]
then
	exit $?
fi

cd $WORKSPACE/controller/token
go test -v
if [ $? -ne 0 ]
then
	exit $?
fi


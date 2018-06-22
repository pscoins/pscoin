#!/usr/bin/env bash

WORKDIR='src/pscoin'
APP_NAME='pscoin'
GO_PROCESS='main.go'


export GOPATH=/opt/gocode;
export PATH=$PATH:$GOPATH/bin


for PROCESS in ${GO_PROCESS}
do
pidf=${GOPATH}/${WORKDIR}_${PROCESS}.pid
exec 221>${pidf}
flock --exclusive --nonblock 221 ||
{
	echo "${pidf} already exists. so we will restart it for you ... ";
	kill `pidof $APP_NAME`
}

echo "${pidf} is not running, I am going to start an instance!!!";

go build -o ${GOPATH}/${WORKDIR}/$APP_NAME ${GOPATH}/${WORKDIR}/${GO_PROCESS}

/usr/local/sbin/daemonize -E BUILD_ID=dontKillMe -c ${GOPATH}/${WORKDIR} -p ${pidf} ${GOPATH}/${WORKDIR}/${APP_NAME} &

done



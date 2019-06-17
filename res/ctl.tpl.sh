#!/bin/bash

export GIN_MODE=release

mkdir -p var

app={{.BinName}}
pidFile=var/pid

function check_pid() {
    if [[ -f ${pidFile} ]];then
        pid=`cat ${pidFile}`
        if [[ -n ${pid} ]]; then
            running=`ps -p ${pid}|grep -v "PID TTY" |wc -l`
            return ${running}
        fi
    fi
    return 0
}

function start() {
    check_pid
    running=$?
    if [[ ${running} -gt 0 ]];then
        echo -n "$app now is running already, pid="
        cat ${pidFile}
        return 1
    fi

    nohup ${app} {{.BinArgs}} >> ./nohup.out 2>&1 &
    sleep 1
    running=`ps -p $! | grep -v "PID TTY" | wc -l`
    if [[ ${running} -gt 0 ]];then
        echo $! > ${pidFile}
        echo "$app started..., pid=$!"
    else
        echo "$app failed to start."
        return 1
    fi
}

function stop() {
    pid=`cat ${pidFile}`
    if [[ $? -eq 0 ]]; then
        kill ${pid}
        rm -f ${pidFile}
    fi
    echo "${app} ${pid} stopped..."
}

function status() {
    check_pid
    running=$?
    if [[ ${running} -gt 0 ]];then
        echo "started, pid=`cat ${pidFile}`"
    else
        echo "stopped!"
    fi
}


if [[ "$1" == "stop" ]];then
    stop
elif [[ "$1" == "start" ]];then
    start
elif [[ "$1" == "restart" ]];then
    stop
    sleep 1
    start
elif [[ "$1" == "status" ]];then
    status
elif [[ "$1" == "tail" ]];then
    ba=`basename ${app}`
    tail -f `ls -l var/${ba}.log | awk '{print $NF}'`
else
    echo "$0 start|stop|restart|status|tail"
fi

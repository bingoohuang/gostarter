#!/bin/bash

export GIN_MODE=release
export PID_FILE=var/pid
mkdir -p var

#set -x #echo on
pidFile=${PID_FILE}
app={{.BinName}}

moreArgs="${*:2}"

function check_pid() {
  if [[ -f ${pidFile} ]]; then
    local pid=$(cat ${pidFile})
    if [[ -n ${pid} ]]; then
      if [[ $(ps -p "${pid}" | grep -v "PID TTY" | wc -l) -gt 0 ]]; then
        echo "${pid}"
        return 1
      fi
    fi
  fi

  # remove prefix ./
  local pureAppName=${app#"./"}
  local running=$(ps -ef | grep $pureAppName | grep -v "grep" | awk '{print $2}')
  if [[ -n ${running} ]]; then
    echo "${running}" >${pidFile}
    echo "${running}"
    return 1
  fi

  echo "0"
  return 0
}

function start() {
  local pid=$(check_pid)
  if [[ ${pid} -gt 0 ]]; then
    echo -n "$app now is running already, pid=$pid"
    return 1
  fi

  nohup ${app} {{.BinArgs}} ${moreArgs} >>./nohup.out 2>&1 &
  sleep 1
  if [[ $(ps -p $! | grep -v "PID TTY" | wc -l) -gt 0 ]]; then
    echo $! >${pidFile}
    echo "$app started..., pid=$!"
    return 0
  else
    echo "$app failed to start."
    return 1
  fi
}

function reload() {
  local pid=$(check_pid)
  if [[ $pid -gt 0 ]]; then
    kill -USR2 "${pid}"
  fi
  sleep 1
  local newPid=$(check_pid)
  echo "${app} ${pid} updated to ${newPid}"
}

function stop() {
  local pid=$(check_pid)
  if [[ $pid -gt 0 ]]; then
    kill "${pid}"
    rm -f ${pidFile}
  fi
  echo "${app} ${pid} stopped..."
}

function status() {
  local pid=$(check_pid)
  if [[ ${pid} -gt 0 ]]; then
    echo "${app} started, pid=$pid"
  else
    echo "${app} stopped!"
  fi
}

function tailLog() {
  local ba=$(basename ${app})
  local logfile="var/${ba}.log"
  local realfile=$(readlink "${logfile}")
  tail -f "$realfile"
}

if [[ "$1" == "stop" ]]; then
  stop
elif [[ "$1" == "start" ]]; then
  start
elif [[ "$1" == "restart" ]]; then
  stop
  sleep 1
  start
elif [[ "$1" == "reload" ]]; then
  reload
elif [[ "$1" == "status" ]]; then
  status
elif [[ "$1" == "tail" ]]; then
  tailLog
else
  echo "$0 start|stop|restart|reload|status|tail"
fi

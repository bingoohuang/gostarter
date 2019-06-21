#!/usr/bin/env python

import argparse
import os


# which find that program is executable or not
def which(program):
    def is_exe(f):
        return os.path.isfile(f) and os.access(f, os.X_OK)

    fp, _ = os.path.split(program)
    if fp:
        if is_exe(program):
            return program
    else:
        for path in os.environ['PATH'].split(os.pathsep):
            exe_file = os.path.join(path, program)
            if is_exe(exe_file):
                return exe_file

    return None


def os_exec(program, check_cmd=''):
    if not check_cmd or which(check_cmd) is not None:
        print program
        os.system(program)


parser = argparse.ArgumentParser(description='building go')
parser.add_argument('-t', '--target', help='target OS, default local', default='local', choices=['linux', 'local'])
parser.add_argument('-u', '--upx', help='enable upx compression if it is available', action='store_true')
parser.add_argument('-b', '--binary', help='binary name, default base of current dir')
args = parser.parse_args()

env = 'env GOOS=linux GOARCH=amd64 ' if args.target == 'linux' else ''
binary = args.binary if args.binary else os.path.basename(os.getcwd())

os_exec('./gv.sh')
os_exec('go get github.com/bingoohuang/statiq')
os_exec('rm -fr statiq/statiq.go')
os_exec('statiq -src=res')
os_exec(env + 'go build -o ' + binary)
os_exec('upx ' + binary, check_cmd='upx') if args.upx else 0

## 什么是 gpm
`gpm` (golang process manager) Go 语言版本的进程管理工具。分为 gpmd (管理进程) 和 gpm (客户端工具)。
得益于 Go 优秀的跨平台功能，`gpm` 目前支持 mac、linux、windows 三种平台。

## gpm 支持的功能
`gpm` 支持以下功能，且三个平台一致:
- 服务的创建、安装、删除、启动、停止、重启、升级、版本回滚
- 服务日志的管理: 监听、切分、过期清理
- ftp 功能: 查看远程目录结构、上传下载文件、远程执行命令

## 使用 gpm
gpm 命令详情如下:
```shell
NAME:
   gpm - package manage tools

USAGE:
   gpm [global options] command [command options] [arguments...]

VERSION:
   v1.1.2

AUTHOR:
   lack <598223084@qq.com>

COMMANDS:
   deploy    deploy gpmd and gpm
   health    confirm gpmd status
   info      get the information of gpmd
   run       run gpmd process
   shutdown  stop gpmd process
   tar       create a package for Install subcommand
   update    update gpm and gpmd
   help, h   Shows a list of commands or help for one command
   bash:
     exec      execute command
     ls        list remote directory
     pull      pull file from service
     push      push files
     terminal  start a terminal
     version   list service history versions
   service:
     create    create a service
     delete    delete a service
     edit      update a service parameters
     get       get service by name
     install   install a service
     list      list all local services
     reboot    reboot a service
     rollback  reboot a service
     start     start a service
     stop      stop a service
     tail      tail service logs
     upgrade   upgrade a service

GLOBAL OPTIONS:
   --dial-timeout int64      specify dial timeout for call option (default: 30s) [$GPM_DIAL_TIMEOUT]
   --host string, -H string  the ip address of gpmd (default: "127.0.0.1:7700") [$GPM_HOST]
   --request-timeout int64   specify request timeout for call option (default: 30s) [$GPM_REQUEST_TIMEOUT]
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```
### 安装 gpm
现在相应操作系统版本的 [gpm](https://github.com/vine-io/gpm/releases) 。
解压后执行命令:
```shell
$ ./gpm deploy
install gpm v1.1.2 successfully!
install gpm successfully!
```

> 注：linux,unix 下在 /usr/local/sbin/ 创建软链接。获得 `gpm` 和 `gpmd` 命令。

### gpmd 相关服务命令
#### 启动 gpmd
```bash
$ gpm run --args '--server-address=0.0.0.0:33700' --args '--enable-log'
start gpmd successfully!
```

#### 停止 gpmd
```shell
$ gpm shutdown
gpmd pid=31048
shutdown gpmd successfully!
```

#### 查看 gpmd 信息
```shell
$ gpm info -o wide
   PROPERTY  |           VALUE
-------------+----------------------------
  Pid        | # 31514
  Version    | v1.1.2-a82f75c-1628482913
  OS         | linux
  Arch       | amd64
  Go version | go1.16.5
  CPU        | 15.16%
  Memory     | 49.23 MB/0.6%
  UpTime     | 35s
```
支持三种格式的输出 `-o wide|json|yaml` 。

#### 检测 gpmd 状态
```shell
$ gpm health
OK
```

#### gpmd 升级
下载不同版本的 gpm 包，使用如下命令升级:
```shell
$ ./linux/gpm update
```
远程升级则需要如下命令:
```shell
$ ./linux/gpm --host 192.168.1.10:7700 --package ./gpm
```
> 注: `--package` 选项指定新版本的二进制包，这种方式可以升级远程机器上不同操作系统下的 gpm。

### 服务操作
#### 远程安装命令
`gpm install` 子命令从本地上传 tar.gz 格式的包到远程机器，并安装服务。该格式的包可以使用 `gpm tar` 子命令创建。
创建一个 *.tar.gz 包:
```shell
$ gpm tar --name /tmp/test.tar.gz --target /opt/test/pp/bin
starting tar /tmp/test.tar.gz
compress /opt/test/pp/bin
compress /opt/test/pp/bin/test
tar /tmp/test.tar.gz successfully
```
> 推荐先 cd 到指定目录的上级目标，再执行 tar 子命令。

安装服务 
```shell
$ gpm --host 192.168.1.10:7700 install --package /tmp/test.tar.gz --name test --dir /opt/test --bin /opt/test/bin/test --auto-restart --version v1.0.0
upload [/tmp/test.tar.gz] 100% |████████████████████████████████████████| (14.593 MB/s)
install service test successfully
```

#### 查看所有服务
```shell
$ gpm list
+------+-----------+-----+------+------------+--------+-----------------------+
| NAME |   USER    | PID | CPU  |   MEMORY   | STATUS |        UPTIME         |
+------+-----------+-----+------+------------+--------+-----------------------+
| test | root:root |   0 | 0.0% | 0.0 B/0.0% | init   | 452357h5m46.15905688s |
+------+-----------+-----+------+------------+--------+-----------------------+

Total: 1
```

#### 查看服务详细信息
```shell
$ gpm get --name test
      PROPERTY      |             VALUE
--------------------+--------------------------------
  Name              | test
  Bin               | /opt/test/bin/test
  Args              |
  Pid               | # 32590
  Dir               | /opt/test
  env               |
  User              | user=root, group=root
  Version           | v1.0.0
  AutoRestart       | True
  CPU               | 0.88%
  Memory            | 7.05 MB/0.1%
  log expire        | 30 days
  log chunk         | 50.00 MB
  CreationTimestamp | 2021-08-09 13:05:39 +0800 CST
  UpdateTimestamp   | 2021-08-09 13:32:24 +0800 CST
  StartTimestamp    | 2021-08-09 13:32:24 +0800 CST
  Status            | running
```
支持三种格式的输出 `-o wide|json|yaml`

#### 启动|停止|重启服务
```
$ gpm start --name test
$ gpm stop --name test
$ gpm reboot --name test
```

#### 创建服务
创建服务和安装服务类似，但是服务相关的文件在 `gpmd` 主机上已存在
```shell
$ gpm create --name gtest --dir /opt/test --bin /opt/test/bin/test --auto-restart --version v1.0.0
upload [/tmp/test.tar.gz] 100% |████████████████████████████████████████| (14.593 MB/s)
service gtest init
```

#### 升级服务
```shell
$ gpm upgrade --name test --package /tmp/test.tar.gz --version v2.0.0
upload [/tmp/test.tar.gz] 100% |████████████████████████████████████████| (4.448 MB/s)
upgrade service test v1.2.8 -> v2.0.0
```

#### 查看服务的历史版本
```shell
$ gpm version --name test
+------+---------+-------------------------------+
| NAME | VERSION |             TIME              |
+------+---------+-------------------------------+
| test | v1.0.0  | 2021-08-09 22:12:08 +0800 CST |
| test | v1.2.3  | 2021-08-09 22:12:36 +0800 CST |
| test | v1.2.4  | 2021-08-09 22:15:06 +0800 CST |
| test | v1.2.5  | 2021-08-09 22:15:59 +0800 CST |
| test | v1.2.7  | 2021-08-09 22:19:38 +0800 CST |
| test | v1.2.8  | 2021-08-09 22:20:30 +0800 CST |
| test | v2.0.0  | 2021-08-09 23:05:54 +0800 CST |
+------+---------+-------------------------------+
```

#### 版本回滚
```shell
$ gpm rollback --name test --revision v1.2.8
rollback test v2.0.0 -> v1.2.8
```

#### 修改服务参数
```shell
$ gpm edit --name test --env "a=b"
edit service 'test' successfully!
$ gpm edit --help
...
OPTIONS:
   --name string, -N string    specify the name for service
   --bin string, -B string     specify the bin for service
   --args strings, -A strings  specify the args for service
   --dir string, -D string     specify the root directory for service
   --env strings, -E strings   specify the env for service
   --user string               specify the user for service
   --group string              specify the group for service
   --log-expire int            specify the expire for service log (default: 0)
   --log-max-size int64        specify the max size for service log (default: 0)
   --auto-restart int          Whether auto restart service when it crashing (1,-1) (default: 1)
   --help, -h                  show help (default: false)
```

#### 查看服务日志
```shell
$ gpm tail --name test
2021-08-09 15:18:19  file=vine/service.go:199 level=info Starting [service] go.vine.helloworld
2021-08-09 15:18:19  file=vine/service.go:200 level=info service [version] v1.0.0
2021-08-09 15:18:19  file=grpc/grpc.go:919 level=info Server [grpc] Listening on [::]:57078
2021-08-09 15:18:19  file=grpc/grpc.go:760 level=info Registry [mdns] Registering node: go.vine.helloworld-dd357c33-8cd4-4911-9155-a152c68f46c6
2021-08-09 15:18:19  file=mdns/mdns_registry.go:266 level=info [mdns] registry create new service with ip: 192.168.3.111 for: 192.168.3.111
```
添加 `-f` 选项可以监听服务的日志变化

#### 删除服务
```shell
$ gpm delete --name test
```
> 注: 删除服务的同时也删除对应的日志，软件包和版本信息


### 其他命令
#### 查看远程文件系统信息
`gpm ls` 功能类似 linux 下 `ls` 命令，可以查看文件文件系统指定目录的信息
```shell
$ gpm ls --path /tmp/san/
+--------+------------+---------+-------------------------------+
|  NAME  |    MODE    |  SIZE   |            MODTIME            |
+--------+------------+---------+-------------------------------+
| ca.pem | -rw-r--r-- | 1.18 KB | 2021-07-17 17:26:26 +0800 CST |
+--------+------------+---------+-------------------------------+
```

#### 执行远程命令
```shell
$ go run cmd/gpm/main.go --host 192.168.3.111:7700 exec --cmd "ls" --A "/tmp/san"
ca.pem
```
支持的参数:
```shell
OPTIONS:
   --cmd string, -C string     specify the command for exec
   --args strings, -A strings  specify the args for exec
   --dir string                specify the directory path for exec
   --env strings, -E strings   specify the env for exec
   --user string               specify the user for exec
   --group string              specify the group for exec
   --help, -h                  show help (default: false)
```

#### 上传文件
```shell
$ gpm --host 192.168.3.111:7700 push --src /tmp/1.txt --dst /tmp/1.txt
 100% |████████████████████████████████████████| (1.712 kB/s)
```

#### 下载文件
```shell
$ gpm pull --src /tmp/1.txt --dst /tmp/11.txt
download [      /tmp/1.txt] [total:  1.18 KB] 100% |████████████████████████████████████████| (1.612 MB/s)
```
# Run Service

## 说明

实现将应用程序以操作系统的Service运行

## 安装

### 编译
安装好golang

Windows下运行.\build.bat，Linux下运行./build，即可在bin生成可执行文件。

### 安装

在Windows下 move .\bin\rs.exe C:\Windows

在Linux下 mv ./bin/rs /usr/local/bin

## 使用方法

```shell
rs <command> [service name] options.. args ...

The command should be:
  install       install service
  uninstall     remove service
  start         start service
  stop          stop service
  restart       restart service
  status        remove service

The options are:
  -e env
        specify serivce application env, eg. 'key1=value1,key2=value2'
  -l relative path
        specify service log relative path (default "./")
  -p path
        specify service absolute application path
  -w workspace
        specify service absolute application workspace

Using args to specify args for application

```


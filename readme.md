# Run Service

## 说明

实现将应用程序以操作系统的Service运行

## 编译

windows运行build.bat即可

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


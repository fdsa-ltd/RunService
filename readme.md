# Run Service

## 说明

实现将应用程序以操作系统的Service运行

## 编译

windows运行build.bat即可



## 使用方法

```shell
ws [arguments] <command>

The arguments are:
  -a args
        specify serivce application args using ','
  -d description
        service description
  -e env
        specify serivce application env, eg. 'key1=value1,key2=value2'
  -help
        for more help
  -n name
        service name
  -p path
        specify service application path
  -w workspace
        specify service application  workspace

The command should be:
  install        install service
  remove         remove service
```


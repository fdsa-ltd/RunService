@echo off
setlocal

FOR /F "delims=" %%i IN ("%cd%") DO (
    set name=%%~ni
) 
if not "%1" == "" (
    del go.mod
    del go.sum
)
if exist go.mod goto build
go mod init fdsa.ltd/%name%
goto build

:build
REM set OLDGOPATH=%GOPATH%
REM set GOPATH=%~dp0

REM gofmt -w src
set GOOS=linux
set GOARCH=amd64
go build -o ./bin/%name% fdsa.ltd/%name%/src
echo go: linux version is finished ok

set GOOS=windows
set GOARCH=amd64
go build -o ./bin/%name%.exe fdsa.ltd/%name%/src
echo go: windows version is finished ok


if "%1" == "upx" goto upx
goto end

:upx
   start upx ./bin/%name%.exe
   start upx ./bin/%name%
   goto end
:end
    REM set GOPATH=%OLDGOPATH%
    echo go: all is finished

set GOOS=linux
set GOARCH=amd64
go build -o ./bin/ws ./src
upx ./bin/ws
set GOOS=windows
set GOARCH=amd64
go build -o ./bin/ws.exe ./src
upx ./bin/ws.exe

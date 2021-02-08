@echo off
echo build...
set GOARCH=amd64
set GOOS=linux
cd ./src
go build -o ../bin/CapturerHub ../src
upx -9 ../bin/CapturerHub
pause

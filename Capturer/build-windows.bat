@echo off
echo build...
set GOARCH=amd64
set GOOS=windows
cd ./src
go build -o ../bin/CapturerHub.exe ../src
upx -9 ../bin/CapturerHub.exe
pause

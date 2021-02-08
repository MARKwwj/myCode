@echo off
set APPNAME=Capturer91
set VERSION=v%date:~2,2%.%date:~5,2%%date:~8,2%
set SAVEPATH=%~dp0/bin/%APPNAME%_%VERSION%
set START_MODULE=91.video

set GOARCH=amd64
set GOOS=linux
echo build %APPNAME% %VERSION% for linux
go build -ldflags "-s -w -X 'main.DefaultModule=%START_MODULE%' -X 'main.CurrentVersion=%VERSION%'" -o %SAVEPATH%/%APPNAME%.lnx %~dp0/src

set GOARCH=amd64
set GOOS=windows
echo build %APPNAME% %VERSION% for windows
go build -ldflags "-s -w -X 'main.DefaultModule=%START_MODULE%' -X 'main.CurrentVersion=%VERSION%'" -o %SAVEPATH%/%APPNAME%.exe %~dp0/src

upx -9 %SAVEPATH%/%APPNAME%.lnx %SAVEPATH%/%APPNAME%.exe
pause

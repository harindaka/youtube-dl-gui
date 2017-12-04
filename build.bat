@echo off
go build -ldflags="-H windowsgui" -o .\bin\youtube-dl-gui.exe .\src || exit /b
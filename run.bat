@echo off
call build.bat || exit /b
.\bin\youtube-dl-gui.exe debug 3030 || exit /b
@echo off
call build.bat || exit /b
.\bin\youtube-dl-gui.exe || exit /b
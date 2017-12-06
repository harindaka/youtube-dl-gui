@echo off
go-bindata -debug -o src/bindata.go lib/... src/ui/...
REM go build -ldflags="-H windowsgui" -o .\bin\youtube-dl-gui.exe .\src || exit /b
go build -o .\bin\youtube-dl-gui.exe .\src || exit /b
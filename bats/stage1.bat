@echo off
cd /d "%~dp0\.."
go run main.go --Script "scripts\stage1.txt"
pause
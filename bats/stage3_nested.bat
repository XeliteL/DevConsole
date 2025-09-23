@echo off
cd /d "%~dp0\.."
go run main.go --Script "scripts\stage3_nested.txt"
pause
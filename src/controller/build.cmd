@echo off
setlocal
set go=go
set GOOS=windows
del .\controller.exe
%go% build
endlocal
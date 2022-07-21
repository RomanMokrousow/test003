@echo off
setlocal
set go=go
set GOOS=windows
del .\host.exe
%go% build
endlocal
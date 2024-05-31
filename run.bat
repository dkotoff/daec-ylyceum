@echo off

for /F "tokens=3" %%i in ('findstr /R "^TIME_DIVISION_MS" .config') do set TIME_DIVISION_MS=%%i
for /F "tokens=3" %%i in ('findstr /R "^TIME_MULTIPLICATION_MS" .config') do set TIME_MULTIPLICATION_MS=%%i
for /F "tokens=3" %%i in ('findstr /R "^TIME_SUBTRACTION_MS" .config') do set TIME_SUBTRACTION_MS=%%i
for /F "tokens=3" %%i in ('findstr /R "^TIME_ADDITION_MS" .config') do set TIME_ADDITION_MS=%%i
for /F "tokens=3" %%i in ('findstr /R "^COMPUTING_POWER" .config') do set COMPUTING_POWER=%%i
for /F "tokens=3" %%i in ('findstr /R "^SERVER_PORT" .config') do set SERVER_PORT=%%i



go build -C ./server/cmd  
go build -C ./agent/cmd

start /B cmd /C ".\server\cmd\cmd.exe"
set pid1=%!

start /B cmd /C ".\agent\cmd\cmd.exe"
set pid2=%!


:wait
    timeout /T 99999 >nul
    goto wait

:exit
    call :cleanup

@echo off

for /f "tokens=3" %%a in ('findstr "^TIME_DIVISION_MS" .config') do set TIME_DIVISION_MS=%%a
for /f "tokens=3" %%a in ('findstr "^TIME_MULTIPLICATION_MS" .config') do set TIME_MULTIPLICATION_MS=%%a
for /f "tokens=3" %%a in ('findstr "^TIME_SUBTRACTION_MS" .config') do set TIME_SUBTRACTION_MS=%%a
for /f "tokens=3" %%a in ('findstr "^TIME_ADDITION_MS" .config') do set TIME_ADDITION_MS=%%a
for /f "tokens=3" %%a in ('findstr "^COMPUTING_POWER" .config') do set COMPUTING_POWER=%%a
for /f "tokens=3" %%a in ('findstr "^SERVER_PORT" .config') do set SERVER_PORT=%%a

go build -C server\cmd  
go build -C agent\cmd

cd server
go build cmd  
start cmd\cmd
set pid1=%!

cd ..\agent
go build cmd
start cmd\cmd
set pid2=%!

cd ..

:cleanup
echo Cleaning up
taskkill /F /PID %pid1%
taskkill /F /PID %pid2%
del server\cmd\cmd.exe
del agent\cmd\cmd.exe
exit /b 0

trap "cleanup" INT

:wait
timeout /t -1 /nobreak >nul
goto wait

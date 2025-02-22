@echo off
setlocal
set SERVICE_NAME=muter

rem Stop and delete the service
sc stop "%SERVICE_NAME%"
sc delete "%SERVICE_NAME%"

rem Remove event log registry entries
reg delete "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\%SERVICE_NAME%" /f

endlocal
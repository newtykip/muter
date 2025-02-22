@echo off
setlocal
set SERVICE_NAME=dev.newty.muter
set RELATIVE_PATH=%1

rem Get the current directory (the absolute path of where the batch script is running)
for %%I in ("%RELATIVE_PATH%") do set ABSOLUTE_PATH=%%~fI

rem Check if the absolute path is empty
if "%ABSOLUTE_PATH%"=="" (
    echo Please provide a relative path to the service executable.
    exit /b 1
)

rem Install the service
sc create "%SERVICE_NAME%" binPath= "%ABSOLUTE_PATH%" DisplayName= "muter" start= auto
sc description "%SERVICE_NAME%" "synchronises system mutes to communication platforms."

rem Auto-restart on failure after 1 minute
sc failure "%SERVICE_NAME%" reset= 0 actions= restart/60000

rem Copy docker security descriptor to the service
sc sdset "%SERVICE_NAME%" D:(A;;CCLCSWRPWPDTLOCRRC;;;SY)(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;BA)(A;;CCLCSWLOCRRC;;;IU)(A;;CCLCSWLOCRRC;;;SU)S:(AU;FA;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;WD)

rem Register event log
reg add "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\%SERVICE_NAME%" /f
reg add "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\%SERVICE_NAME%" /v EventMessageFile /t REG_SZ /d "C:\Windows\System32\EventCreate.exe" /f
reg add "HKLM\SYSTEM\CurrentControlSet\Services\EventLog\Application\%SERVICE_NAME%" /v TypesSupported /t REG_DWORD /d 7 /f

endlocal
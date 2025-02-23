Name "muter"
OutFile "muter-installer.exe"

Icon "mic-off.ico"
UninstallIcon "mic-off.ico"

InstallDir "$PROGRAMFILES\muter"
Page Directory
Page InstFiles

Section "Install muter"
  SetOutPath "$INSTDIR"
  File "muter.exe"
  File "mic-off.ico"

  nsExec::ExecToLog 'sc create dev.newty.muter binPath= "$INSTDIR\muter.exe" DisplayName= "muter" start= auto'
  nsExec::ExecToLog 'sc description dev.newty.muter "synchronises system mutes to communication platforms"'
  nsExec::ExecToLog 'sc failure dev.newty.muter reset= 0 actions= restart/60000'
  
  ; Add event log registry entries
  WriteRegStr HKLM "SYSTEM\CurrentControlSet\Services\EventLog\Application\dev.newty.muter" "" ""
  WriteRegStr HKLM "SYSTEM\CurrentControlSet\Services\EventLog\Application\dev.newty.muter" "EventMessageFile" "C:\Windows\System32\EventCreate.exe"
  WriteRegDWORD HKLM "SYSTEM\CurrentControlSet\Services\EventLog\Application\dev.newty.muter" "TypesSupported" 7
  
  nsExec::ExecToLog 'sc start dev.newty.muter'
  
  ; Add uninstall information to registry
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter" "DisplayName" "muter"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter" "UninstallString" "$INSTDIR\uninstall.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter" "InstallLocation" "$INSTDIR"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter" "DisplayIcon" "$INSTDIR\mic-off.ico"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter" "Publisher" "newty.dev"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter" "DisplayVersion" "1.0.0"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter" "VersionMajor" "1"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter" "VersionMinor" "0"
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter" "NoRepair" 1

  WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

Section "Uninstall"
  nsExec::ExecToLog 'sc stop dev.newty.muter'
  nsExec::ExecToLog 'sc delete dev.newty.muter'
  
  ; Remove event log registry entries
  DeleteRegKey HKLM "SYSTEM\CurrentControlSet\Services\EventLog\Application\dev.newty.muter"
  
  ; Remove uninstall information from registry
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\muter"
  
  RMDir /r "$INSTDIR"
SectionEnd
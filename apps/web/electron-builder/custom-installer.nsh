; Electron Builder's default process detector can falsely report a running app
; on some Windows configurations. During install/update, terminate only the
; exact application executable and continue when it is already absent.
!macro customCheckAppRunning
  nsExec::Exec '"$SYSDIR\taskkill.exe" /F /T /IM "${APP_EXECUTABLE_FILENAME}"'
  Pop $0
  Sleep 500
!macroend

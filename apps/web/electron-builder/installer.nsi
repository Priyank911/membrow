; Custom NSIS entrypoint. It keeps Electron Builder's assisted installer flow,
; but creates the uninstaller from this script. This avoids Electron Builder's
; separate temporary-uninstaller launch, which fails with spawn UNKNOWN here.
Var newStartMenuLink
Var oldStartMenuLink
Var newDesktopLink
Var oldDesktopLink
Var oldShortcutName
Var oldMenuDirectory

!include "common.nsh"
!include "MUI2.nsh"
!include "multiUser.nsh"
!include "allowOnlyOneInstallerInstance.nsh"
!include "installer.nsh"

!ifdef INSTALL_MODE_PER_ALL_USERS
  RequestExecutionLevel admin
!else
  RequestExecutionLevel user
!endif

Var appExe
Var launchLink

!ifdef ONE_CLICK
  !include "oneClick.nsh"
!else
  !include "assistedInstaller.nsh"
!endif

!insertmacro addLangs

Function .onInit
  Call setInstallSectionSpaceRequired
  SetOutPath $INSTDIR
  ${LogSet} on
  !insertmacro check64BitAndSetRegView
  !insertmacro initMultiUser
FunctionEnd

!include "installUtil.nsh"

Section "install" INSTALL_SECTION_ID
  InitPluginsDir
  ${IfNot} ${Silent}
    SetDetailsPrint none
  ${endif}

  StrCpy $appExe "$INSTDIR\${APP_EXECUTABLE_FILENAME}"
  !insertmacro setLinkVars

  ${ifNot} ${UAC_IsInnerInstance}
    !insertmacro CHECK_APP_RUNNING
  ${endif}

  Var /GLOBAL keepShortcuts
  StrCpy $keepShortcuts "false"
  !insertmacro setIsTryToKeepShortcuts
  ${if} $isTryToKeepShortcuts == "true"
    ReadRegStr $R1 SHELL_CONTEXT "${INSTALL_REGISTRY_KEY}" KeepShortcuts
    ${if} $R1 == "true"
    ${andIf} ${FileExists} "$appExe"
      StrCpy $keepShortcuts "true"
    ${endIf}
  ${endif}

  !insertmacro uninstallOldVersion SHELL_CONTEXT
  !insertmacro handleUninstallResult SHELL_CONTEXT

  ${if} $installMode == "all"
    !insertmacro uninstallOldVersion HKEY_CURRENT_USER
    !insertmacro handleUninstallResult HKEY_CURRENT_USER
  ${endIf}

  SetOutPath $INSTDIR
  !ifdef UNINSTALLER_ICON
    File /oname=uninstallerIcon.ico "${UNINSTALLER_ICON}"
  !endif

  !insertmacro extractEmbeddedAppPackage
  WriteUninstaller "$INSTDIR\${UNINSTALL_FILENAME}"
  !insertmacro registryAddInstallInfo
  !insertmacro addStartMenuLink $keepShortcuts
  !insertmacro addDesktopLink $keepShortcuts

  ${if} ${FileExists} "$newStartMenuLink"
    StrCpy $launchLink "$newStartMenuLink"
  ${else}
    StrCpy $launchLink "$INSTDIR\${APP_EXECUTABLE_FILENAME}"
  ${endIf}
SectionEnd

Function setInstallSectionSpaceRequired
  !insertmacro setSpaceRequired ${INSTALL_SECTION_ID}
FunctionEnd

; The normal Electron Builder template compiles this code in a second
; executable. Keeping it here lets WriteUninstaller create a standard NSIS
; uninstaller without launching that blocked temporary executable at build time.
Function un.onInit
  SetOutPath $INSTDIR
  ${LogSet} on
  !insertmacro check64BitAndSetRegView
  nsExec::Exec '"$SYSDIR\taskkill.exe" /F /T /IM "${APP_EXECUTABLE_FILENAME}"'
  Pop $0
  Sleep 500
FunctionEnd

Section "un.Uninstall"
  SetShellVarContext current
  SetOutPath $TEMP
  RMDir /r $INSTDIR
  Delete "$DESKTOP\${SHORTCUT_NAME}.lnk"
  Delete "$SMPROGRAMS\${SHORTCUT_NAME}.lnk"
  DeleteRegKey HKCU "${UNINSTALL_REGISTRY_KEY}"
  DeleteRegKey HKCU "${INSTALL_REGISTRY_KEY}"
SectionEnd

; Inno Setup script for rgit (https://jrsoftware.org/isinfo.php)
;
; Build with scripts\build.ps1, or directly:
;   iscc /DAppVersion=2.0.0 installer\rgit.iss
; The script expects the executable at dist\rgit-windows-x64.exe and writes
; the installer to dist\rgit-windows-x64-setup.exe.

#ifndef AppVersion
  #define AppVersion "2.0.0"
#endif

#define AppName      "rgit"
#define AppPublisher "Ettisaf Rup"
#define AppURL       "https://github.com/ettisafxrup/rgit"
#define AppExe       "rgit.exe"

; The 1.x shell-based installer used this id. Setup removes that version
; because its extensionless "rgit" script would shadow rgit.exe in Git Bash.
#define LegacyAppId  "{E3A6AC5C-CCC8-41D4-A6D5-5EEF62965525}"

[Setup]
AppId={{9B6F1D2E-4C7A-4E5B-8F3D-2A1C0E7B5D94}
AppName={#AppName}
AppVersion={#AppVersion}
AppVerName={#AppName} {#AppVersion}
AppPublisher={#AppPublisher}
AppPublisherURL={#AppURL}
AppSupportURL={#AppURL}/issues
AppUpdatesURL={#AppURL}/releases
VersionInfoVersion={#AppVersion}
VersionInfoDescription={#AppName} setup

DefaultDirName={autopf}\{#AppName}
DefaultGroupName={#AppName}
DisableProgramGroupPage=yes
DisableDirPage=auto

; Let the user choose between "only for me" (no admin rights) and "all users".
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog commandline

ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

; Broadcast PATH changes so new terminals pick them up without a reboot.
ChangesEnvironment=yes

LicenseFile=..\LICENSE
OutputDir=..\dist
OutputBaseFilename=rgit-windows-x64-setup
SetupIconFile=..\assets\rgit.ico
UninstallDisplayIcon={app}\{#AppExe}
UninstallDisplayName={#AppName} {#AppVersion}
WizardStyle=modern dynamic
WizardImageFile=wizard-large-100.bmp,wizard-large-200.bmp
WizardSmallImageFile=wizard-small-100.bmp,wizard-small-200.bmp
Compression=lzma2/max
SolidCompression=yes

; build.ps1 -SignCert defines Sign and the "rgitsign" tool (see scripts\windows\sign.ps1).
#ifdef Sign
SignTool=rgitsign
SignedUninstaller=yes
#endif

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Messages]
FinishedLabel=rgit is installed.%n%nOpen a new terminal (Command Prompt, PowerShell or Git Bash) and run:%n%n    rgit help

[Tasks]
Name: "addtopath"; Description: "Add rgit to the PATH (recommended, lets you run rgit from any terminal)"
Name: "startmenu"; Description: "Create a Start menu shortcut to an rgit terminal"; Flags: unchecked

[Files]
Source: "..\dist\rgit-windows-x64.exe"; DestDir: "{app}"; DestName: "{#AppExe}"; Flags: ignoreversion
Source: "..\README.md";       DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENSE";         DestDir: "{app}"; DestName: "LICENSE.txt"; Flags: ignoreversion
Source: "..\assets\rgit.ico"; DestDir: "{app}"; Flags: ignoreversion

[InstallDelete]
; Leftovers of the 1.x shell version, in case it was installed to this folder.
Type: files;          Name: "{app}\rgit"
Type: files;          Name: "{app}\ARCHITECTURE"
Type: filesandordirs; Name: "{app}\prompts"
Type: filesandordirs; Name: "{app}\configs"
Type: filesandordirs; Name: "{app}\installer"
Type: filesandordirs; Name: "{app}\assets"

[Icons]
Name: "{group}\rgit terminal"; Filename: "{cmd}"; Parameters: "/k ""{app}\{#AppExe}"" help"; \
    WorkingDir: "{userdocs}"; IconFilename: "{app}\rgit.ico"; Tasks: startmenu
Name: "{group}\Uninstall rgit"; Filename: "{uninstallexe}"; Tasks: startmenu

[Code]
const
  SystemEnvKey = 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment';
  UserEnvKey   = 'Environment';
  UninstallRoot = 'Software\Microsoft\Windows\CurrentVersion\Uninstall\';

{ ---------- PATH handling ---------- }

function EnvRoot: Integer;
begin
  if IsAdminInstallMode then
    Result := HKEY_LOCAL_MACHINE
  else
    Result := HKEY_CURRENT_USER;
end;

function EnvKey: String;
begin
  if IsAdminInstallMode then
    Result := SystemEnvKey
  else
    Result := UserEnvKey;
end;

{ Position of Dir inside a ';'-wrapped, upper-cased PATH, or 0. }
function FindInPath(const WrappedPath, Dir: String): Integer;
begin
  Result := Pos(';' + Uppercase(Dir) + ';', Uppercase(WrappedPath));
end;

procedure AddToPath(const Dir: String);
var
  Paths: String;
begin
  if not RegQueryStringValue(EnvRoot, EnvKey, 'Path', Paths) then
    Paths := '';
  if FindInPath(';' + Paths + ';', Dir) > 0 then
    exit;
  if (Paths <> '') and (Paths[Length(Paths)] <> ';') then
    Paths := Paths + ';';
  if RegWriteExpandStringValue(EnvRoot, EnvKey, 'Path', Paths + Dir) then
    Log('Added to PATH: ' + Dir)
  else
    Log('Could not update PATH');
end;

procedure RemoveFromPath(const Dir: String);
var
  Paths: String;
  Position: Integer;
begin
  if not RegQueryStringValue(EnvRoot, EnvKey, 'Path', Paths) then
    exit;
  Paths := ';' + Paths + ';';
  Position := FindInPath(Paths, Dir);
  if Position = 0 then
    exit;
  { Remove ";Dir" and then the ';' we wrapped the value in. }
  Delete(Paths, Position, Length(Dir) + 1);
  Paths := Copy(Paths, 2, Length(Paths) - 2);
  if RegWriteExpandStringValue(EnvRoot, EnvKey, 'Path', Paths) then
    Log('Removed from PATH: ' + Dir);
end;

{ ---------- Upgrading from the 1.x shell version ---------- }

// Finds the 1.x uninstaller. Its script had a typo in the AppId (a doubled
// closing brace), so its registry key is spelled that way; the correct
// spelling is checked as well.
function LegacyUninstaller(var RootKey: Integer): String;
var
  Keys: array of String;
  Roots: array of Integer;
  i, j: Integer;
begin
  Result := '';
  SetArrayLength(Keys, 2);
  Keys[0] := UninstallRoot + '{#LegacyAppId}}_is1';
  Keys[1] := UninstallRoot + '{#LegacyAppId}_is1';
  SetArrayLength(Roots, 2);
  Roots[0] := HKLM32;
  Roots[1] := HKEY_CURRENT_USER;
  for i := 0 to GetArrayLength(Roots) - 1 do
    for j := 0 to GetArrayLength(Keys) - 1 do
      if RegQueryStringValue(Roots[i], Keys[j], 'UninstallString', Result) then
      begin
        RootKey := Roots[i];
        Result := RemoveQuotes(Result);
        exit;
      end;
end;

function PrepareToInstall(var NeedsRestart: Boolean): String;
var
  Uninstaller: String;
  RootKey, ResultCode: Integer;
begin
  Result := '';
  Uninstaller := LegacyUninstaller(RootKey);
  if (Uninstaller = '') or not FileExists(Uninstaller) then
    exit;

  { A machine-wide legacy install can only be removed with admin rights. }
  if (RootKey = HKLM32) and not IsAdminInstallMode then
  begin
    Log('Legacy rgit found but setup is not elevated: ' + Uninstaller);
    SuppressibleMsgBox('An older, shell-based version of rgit is installed for all users.' + #13#10 +
      'Please remove it from Settings > Apps after this setup finishes, ' +
      'otherwise it may take precedence over this version in Git Bash.',
      mbInformation, MB_OK, IDOK);
    exit;
  end;

  Log('Removing legacy rgit: ' + Uninstaller);
  if not Exec(Uninstaller, '/VERYSILENT /SUPPRESSMSGBOXES /NORESTART', '', SW_HIDE,
              ewWaitUntilTerminated, ResultCode) or (ResultCode <> 0) then
    Log('Legacy uninstaller did not finish cleanly (code ' + IntToStr(ResultCode) + ')');
end;

{ ---------- Install / uninstall steps ---------- }

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if (CurStep = ssPostInstall) and WizardIsTaskSelected('addtopath') then
    AddToPath(ExpandConstant('{app}'));
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usUninstall then
    RemoveFromPath(ExpandConstant('{app}'));
end;

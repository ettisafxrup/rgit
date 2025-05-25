#define MyAppName "rgit"
#define MyAppVersion "1.0"
#define MyAppPublisher "Qxtend"
#define MyAppURL "https://rgit.netlify.app"

[Setup]
AppId={{E3A6AC5C-CCC8-41D4-A6D5-5EEF62965525}}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
LicenseFile=D:\Codes\terminal\rgit\license\LICENSE.txt
InfoBeforeFile=D:\Codes\terminal\rgit\license\README.txt
InfoAfterFile=D:\Codes\terminal\rgit\license\THANKS.txt
PrivilegesRequired=admin
OutputDir=D:\Output
OutputBaseFilename=rgit
SetupIconFile=C:\Users\HP\OneDrive\ডেস্কটপ\favicon (1).ico
SolidCompression=yes
WizardStyle=modern

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "D:\Codes\terminal\rgit\README.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "D:\Codes\terminal\rgit\rgit"; DestDir: "{app}"; Flags: ignoreversion
Source: "D:\Codes\terminal\rgit\*"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs; Excludes: ".git\*;node_modules\*"

[Registry]
// Add {app} folder to system PATH if not already present
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; \
    ValueType: expandsz; ValueName: "Path"; \
    ValueData: "{code:GetNewPath|{olddata}}"; Flags: preservestringtype uninsdeletevalue

[Code]

const
  // HWND_BROADCAST = $FFFF;
  WM_SETTINGCHANGE = $001A;
  SMTO_ABORTIFHUNG = $0002;

function SendMessageTimeout(hWnd: Integer; Msg: Integer; wParam: Integer;
  lParam: AnsiString; fuFlags: Integer; uTimeout: Integer;
  var lpdwResult: Integer): Integer;
  external 'SendMessageTimeoutA@user32.dll stdcall';

// Add {app} to PATH if not already present
function GetNewPath(OldPathParam: string): string;
var
  OldPath, AppPath: string;
begin
  OldPath := OldPathParam;
  AppPath := ExpandConstant('{app}');
  if Pos(';' + LowerCase(AppPath) + ';', ';' + LowerCase(OldPath) + ';') = 0 then
  begin
    if (OldPath = '') or (Copy(OldPath, Length(OldPath), 1) = ';') then
      Result := OldPath + AppPath
    else
      Result := OldPath + ';' + AppPath;
  end
  else
    Result := OldPath;
end;

procedure NotifyEnvironmentChange;
var
  ResultCode: Integer;
begin
  SendMessageTimeout(HWND_BROADCAST, WM_SETTINGCHANGE, 0, 'Environment', SMTO_ABORTIFHUNG, 5000, ResultCode);
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
    NotifyEnvironmentChange;
end;

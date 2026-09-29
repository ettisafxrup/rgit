<#
.SYNOPSIS
    Builds rgit: runs the tests, compiles the binaries and the Windows installer.

.EXAMPLE
    .\scripts\build.ps1                 # test, build dist\rgit.exe and the installer
    .\scripts\build.ps1 -All            # also build Linux and macOS binaries
    .\scripts\build.ps1 -SkipTests -SkipInstaller
#>
param(
    [string]$Version,
    [switch]$All,
    [switch]$SkipTests,
    [switch]$SkipInstaller,
    [switch]$Assets   # regenerate the logo, icon and installer artwork
)

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

function Step([string]$message) { Write-Host "`n> $message" -ForegroundColor Cyan }

function Invoke-Checked([string]$exe, [string[]]$arguments) {
    & $exe @arguments
    if ($LASTEXITCODE -ne 0) { throw "$exe $($arguments -join ' ') failed with exit code $LASTEXITCODE" }
}

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw 'Go is not installed or not on PATH (https://go.dev/dl).'
}

if (-not $Version) {
    $match = Select-String -Path 'internal\version\version.go' -Pattern 'Version = "([^"]+)"'
    $Version = $match.Matches[0].Groups[1].Value
}
Write-Host "Building rgit $Version" -ForegroundColor Green

if ($Assets) {
    Step 'Generating logo and installer artwork'
    Invoke-Checked go @('run', './tools/genlogo')
}

if (-not $SkipTests) {
    Step 'Running tests'
    Invoke-Checked go @('vet', './...')
    Invoke-Checked go @('test', './...')
}

# Embed the icon and version details into rgit.exe (needs windres from MinGW;
# without it the committed .syso file is used as is).
$windres = Get-Command windres -ErrorAction SilentlyContinue
if ($windres) {
    Step 'Compiling Windows resources'
    $numeric = (($Version -split '[.-]')[0..2] + '0') -join ','
    Push-Location 'scripts\windows'
    try {
        Invoke-Checked $windres.Source @("-DRGIT_VERSION=$Version", "-DRGIT_VERSION_NUM=$numeric",
            '-i', 'rgit.rc', '-O', 'coff', '-o', '..\..\cmd\rgit\rsrc_windows_amd64.syso')
    } finally { Pop-Location }
}

$ldflags = "-s -w -X github.com/ettisafxrup/rgit/internal/version.Version=$Version"

# dist\ holds exactly the files of a release, so start from an empty folder.
$dist = Join-Path $root 'dist'
if (Test-Path $dist) { Get-ChildItem $dist -File | Remove-Item -Force }
New-Item -ItemType Directory -Force $dist | Out-Null

# Release file names carry no version, so "latest release" download links
# (used by the website and install.sh) never change.
$targets = @(@{ os = 'windows'; arch = 'amd64'; out = 'rgit-windows-x64.exe' })
if ($All) {
    $targets += @(
        @{ os = 'linux';  arch = 'amd64'; out = 'rgit-linux-x64' },
        @{ os = 'linux';  arch = 'arm64'; out = 'rgit-linux-arm64' },
        @{ os = 'darwin'; arch = 'amd64'; out = 'rgit-macos-x64' },
        @{ os = 'darwin'; arch = 'arm64'; out = 'rgit-macos-arm64' }
    )
}

foreach ($t in $targets) {
    Step "Building $($t.os)/$($t.arch) -> dist\$($t.out)"
    $env:GOOS = $t.os; $env:GOARCH = $t.arch; $env:CGO_ENABLED = '0'
    try {
        Invoke-Checked go @('build', '-trimpath', '-ldflags', $ldflags, '-o', "dist\$($t.out)", './cmd/rgit')
    } finally {
        Remove-Item Env:GOOS, Env:GOARCH, Env:CGO_ENABLED -ErrorAction SilentlyContinue
    }
}

if (-not $SkipInstaller) {
    $iscc = Get-Command iscc -ErrorAction SilentlyContinue
    if (-not $iscc) {
        $iscc = Get-ChildItem 'C:\Program Files*\Inno Setup *\ISCC.exe' -ErrorAction SilentlyContinue |
            Sort-Object FullName -Descending | Select-Object -First 1
    }
    if (-not $iscc) { throw 'Inno Setup (ISCC.exe) was not found. Install it from https://jrsoftware.org or pass -SkipInstaller.' }

    Step 'Building the installer -> dist\rgit-windows-x64-setup.exe'
    $isccPath = if ($iscc.Source) { $iscc.Source } else { $iscc.FullName }
    Invoke-Checked $isccPath @('/Q', "/DAppVersion=$Version", 'installer\rgit.iss')
}

Step 'Writing checksums -> dist\SHA256SUMS.txt'
$sums = Get-ChildItem $dist -File | Where-Object Name -ne 'SHA256SUMS.txt' | Sort-Object Name | ForEach-Object {
    '{0}  {1}' -f (Get-FileHash $_.FullName -Algorithm SHA256).Hash.ToLower(), $_.Name
}
# Unix line endings, so "sha256sum -c SHA256SUMS.txt" works on Linux and macOS.
[IO.File]::WriteAllText((Join-Path $dist 'SHA256SUMS.txt'), (($sums -join "`n") + "`n"))

Write-Host "`nRelease files in dist\:" -ForegroundColor Green
Get-ChildItem $dist -File | ForEach-Object { Write-Host ('  {0,-28} {1,8:N1} MB' -f $_.Name, ($_.Length / 1MB)) }
Write-Host "`nDone." -ForegroundColor Green

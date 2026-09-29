<#
.SYNOPSIS
    Code-signs a Windows file with a certificate from the Windows certificate store.

.DESCRIPTION
    Used by build.ps1 for rgit.exe and by Inno Setup for the installer and
    uninstaller. Signed files show your name instead of "Unknown publisher"
    and build SmartScreen reputation over time.

.EXAMPLE
    .\scripts\windows\sign.ps1 -Path dist\rgit-windows-x64.exe -Thumbprint 0123ABCD...
#>
param(
    [Parameter(Mandatory)] [string]$Path,
    [Parameter(Mandatory)] [string]$Thumbprint,
    [string]$TimestampServer = 'http://timestamp.digicert.com'
)

$ErrorActionPreference = 'Stop'

$cert = Get-ChildItem Cert:\CurrentUser\My, Cert:\LocalMachine\My |
    Where-Object { $_.Thumbprint -eq $Thumbprint -and $_.HasPrivateKey } |
    Select-Object -First 1
if (-not $cert) { throw "No code-signing certificate with thumbprint $Thumbprint was found." }

$options = @{ FilePath = $Path; Certificate = $cert; HashAlgorithm = 'SHA256' }
# The timestamp keeps the signature valid after the certificate expires.
if ($TimestampServer) { $options.TimestampServer = $TimestampServer }

$result = Set-AuthenticodeSignature @options
if (-not $result.SignerCertificate) { throw "Signing $Path failed: $($result.StatusMessage)" }
Write-Host "Signed $Path ($($result.Status))"

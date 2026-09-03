[CmdletBinding()]
param(
  [switch]$SkipFrontend,
  [string]$GoArch = 'amd64'
)

$ErrorActionPreference = 'Stop'
$Root = Split-Path -Parent $PSScriptRoot
$Frontend = Join-Path $Root 'frontend'
$Deploy = Join-Path $Root 'deploy'
$Release = Join-Path $Root 'release'
$FrontendDist = Join-Path $Deploy 'frontend-dist'

if (-not $SkipFrontend) {
  Write-Host 'Building frontend...'
  Push-Location $Frontend
  try { pnpm.cmd build } finally { Pop-Location }
}

$BuiltDist = Join-Path $Frontend 'dist'
if (-not (Test-Path (Join-Path $BuiltDist 'index.html'))) {
  throw "Frontend build output is missing: $BuiltDist\index.html"
}
New-Item -ItemType Directory -Path $FrontendDist -Force | Out-Null
Get-ChildItem -LiteralPath $FrontendDist -Force -ErrorAction SilentlyContinue |
  Remove-Item -Recurse -Force
Copy-Item -Path (Join-Path $BuiltDist '*') -Destination $FrontendDist -Recurse -Force

New-Item -ItemType Directory -Path $Release -Force | Out-Null
$oldGoOS = $env:GOOS
$oldGoArch = $env:GOARCH
$oldCgo = $env:CGO_ENABLED
try {
  $env:GOOS = 'linux'
  $env:GOARCH = $GoArch
  $env:CGO_ENABLED = '0'
  Push-Location (Join-Path $Root 'backend')
  try {
    Write-Host "Building Linux/$GoArch central binary..."
    go build -trimpath -o (Join-Path $Release 'xpanel-central') ./cmd/server
  } finally { Pop-Location }

  Push-Location (Join-Path $Root 'agent')
  try {
    Write-Host "Building Linux/$GoArch agent binary..."
    go build -trimpath -o (Join-Path $Release 'xpanel-agent') ./cmd/agent
  } finally { Pop-Location }
} finally {
  if ($null -eq $oldGoOS) { Remove-Item Env:GOOS -ErrorAction SilentlyContinue } else { $env:GOOS = $oldGoOS }
  if ($null -eq $oldGoArch) { Remove-Item Env:GOARCH -ErrorAction SilentlyContinue } else { $env:GOARCH = $oldGoArch }
  if ($null -eq $oldCgo) { Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue } else { $env:CGO_ENABLED = $oldCgo }
}

Write-Host 'Bundle ready:'
Get-Item (Join-Path $FrontendDist 'index.html'), (Join-Path $Release 'xpanel-central') |
  Select-Object FullName, Length

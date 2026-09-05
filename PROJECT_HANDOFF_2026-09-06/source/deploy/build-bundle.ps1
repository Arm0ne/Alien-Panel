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
$AgentVersionFile = Join-Path $Root 'agent/VERSION'
if (-not (Test-Path -LiteralPath $AgentVersionFile)) {
  throw "Agent version file is missing: $AgentVersionFile"
}
$AgentVersion = (Get-Content -LiteralPath $AgentVersionFile -Raw).Trim()
if ($AgentVersion -notmatch '^v\d+\.\d+\.\d+(?:-[0-9A-Za-z]+(?:\.[0-9A-Za-z]+)*)?(?:\+[0-9A-Za-z]+(?:\.[0-9A-Za-z]+)*)?$') {
  throw "Agent version must use semantic versioning, for example v1.0.1: $AgentVersion"
}
$AgentCommit = (git -C $Root rev-parse --short=12 HEAD 2>$null).Trim()
$AgentBuildTime = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')
$AgentLdFlags = "-X xpanel-central/agent/internal/buildinfo.Version=$AgentVersion -X xpanel-central/agent/internal/buildinfo.Commit=$AgentCommit -X xpanel-central/agent/internal/buildinfo.BuildTime=$AgentBuildTime"

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
    go build -trimpath -ldflags $AgentLdFlags -o (Join-Path $Release 'xpanel-agent') ./cmd/agent
  } finally { Pop-Location }
} finally {
  if ($null -eq $oldGoOS) { Remove-Item Env:GOOS -ErrorAction SilentlyContinue } else { $env:GOOS = $oldGoOS }
  if ($null -eq $oldGoArch) { Remove-Item Env:GOARCH -ErrorAction SilentlyContinue } else { $env:GOARCH = $oldGoArch }
  if ($null -eq $oldCgo) { Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue } else { $env:CGO_ENABLED = $oldCgo }
}

$AgentChecksum = Join-Path $Release 'xpanel-agent.sha256'
(Get-FileHash -LiteralPath (Join-Path $Release 'xpanel-agent') -Algorithm SHA256).Hash.ToLowerInvariant() |
  Set-Content -LiteralPath $AgentChecksum -NoNewline -Encoding ascii

Write-Host 'Bundle ready:'
Write-Host "Embedded Agent version: $AgentVersion"
Get-Item (Join-Path $FrontendDist 'index.html'), (Join-Path $Release 'xpanel-central'), (Join-Path $Release 'xpanel-agent.sha256') |
  Select-Object FullName, Length

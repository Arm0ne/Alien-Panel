[CmdletBinding(SupportsShouldProcess)]
param(
    [Parameter(Mandatory = $true)] [string]$Source,
    [string]$Database = $(if ($env:XPANEL_DATABASE) { $env:XPANEL_DATABASE } else { [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\backend\data\panel.db')) }),
    [switch]$Yes,
    [string]$Binary = $(if ($env:XPANEL_DB_MAINTENANCE_BIN) { $env:XPANEL_DB_MAINTENANCE_BIN } else { [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\backend\db-maintenance.exe')) })
)

$ErrorActionPreference = 'Stop'
if (-not $Yes) { throw 'Restore replaces the database. Re-run with -Yes after stopping the central service.' }
if (-not (Test-Path -LiteralPath $Binary -PathType Leaf)) { throw "db-maintenance binary not found: $Binary" }
if (-not (Test-Path -LiteralPath $Source -PathType Leaf)) { throw "backup file not found: $Source" }
if ($PSCmdlet.ShouldProcess($Database, "restore $Source")) {
    & $Binary restore --source $Source --database $Database --yes
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
}

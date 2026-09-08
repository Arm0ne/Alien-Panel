[CmdletBinding()]
param(
    [string]$Database = $(if ($env:XPANEL_DATABASE) { $env:XPANEL_DATABASE } else { [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\backend\data\panel.db')) }),
    [string]$BackupDirectory = $(if ($env:XPANEL_BACKUP_DIR) { $env:XPANEL_BACKUP_DIR } else { [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\backend\data\backups')) }),
    [int]$Retention = $(if ($env:XPANEL_BACKUP_RETENTION) { [int]$env:XPANEL_BACKUP_RETENTION } else { 14 }),
    [string]$Binary = $(if ($env:XPANEL_DB_MAINTENANCE_BIN) { $env:XPANEL_DB_MAINTENANCE_BIN } else { [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\backend\db-maintenance.exe')) })
)

$ErrorActionPreference = 'Stop'
if (-not (Test-Path -LiteralPath $Binary -PathType Leaf)) {
    throw "db-maintenance binary not found: $Binary"
}
& $Binary backup --database $Database --backup-dir $BackupDirectory --retention $Retention
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

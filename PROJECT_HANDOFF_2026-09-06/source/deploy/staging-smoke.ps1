[CmdletBinding()]
param(
    [string]$BaseUrl = $(if ($env:XPANEL_STAGING_URL) { $env:XPANEL_STAGING_URL } else { 'http://localhost:8090' }),
    [string]$AdminUser = $(if ($env:XPANEL_ADMIN_USER) { $env:XPANEL_ADMIN_USER } else { 'admin' }),
    [string]$AdminPassword = $env:XPANEL_ADMIN_PASSWORD
)

$ErrorActionPreference = 'Stop'
$BaseUrl = $BaseUrl.TrimEnd('/')

function Assert-Envelope([object]$Response, [string]$Name, [string]$Status) {
    if ($Response.code -ne '0000' -or $Response.data.status -ne $Status) {
        throw "$Name returned an unexpected response"
    }
}

$live = Invoke-RestMethod -Uri "$BaseUrl/health/live" -Method Get
$ready = Invoke-RestMethod -Uri "$BaseUrl/health/ready" -Method Get
Assert-Envelope $live 'health/live' 'ok'
Assert-Envelope $ready 'health/ready' 'ready'

$headers = (Invoke-WebRequest -UseBasicParsing -Uri "$BaseUrl/health/live" -Method Get).Headers
if ($headers['X-Content-Type-Options'] -ne 'nosniff') { throw 'missing X-Content-Type-Options' }
if ($headers['X-Frame-Options'] -ne 'DENY') { throw 'missing X-Frame-Options' }
if ($headers['Cache-Control'] -ne 'no-store') { throw 'missing Cache-Control: no-store' }

if (-not [string]::IsNullOrWhiteSpace($AdminPassword)) {
    $loginBody = @{ userName = $AdminUser; password = $AdminPassword } | ConvertTo-Json -Compress
    $login = Invoke-RestMethod -Uri "$BaseUrl/api/auth/login" -Method Post -ContentType 'application/json' -Body $loginBody
    if ($login.code -ne '0000' -or [string]::IsNullOrWhiteSpace($login.data.token)) { throw 'administrator login failed' }
    $dashboard = Invoke-RestMethod -Uri "$BaseUrl/api/dashboard" -Method Get -Headers @{ Authorization = "Bearer $($login.data.token)" }
    if ($dashboard.code -ne '0000') { throw 'protected dashboard request failed' }
}

Write-Output "staging smoke check passed: $BaseUrl"

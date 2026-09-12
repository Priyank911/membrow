$Root = Split-Path -Parent $PSScriptRoot

$ApiPath = Join-Path $Root "apps/api"
$WebPath = Join-Path $Root "apps/web"

$api = Start-Process -FilePath "go" -ArgumentList "run", "./cmd/api" -WorkingDirectory $ApiPath -PassThru
$web = Start-Process -FilePath "npm" -ArgumentList "run", "dev" -WorkingDirectory $WebPath -PassThru

Write-Host "Membrow API PID: $($api.Id)"
Write-Host "Membrow Web PID: $($web.Id)"
Write-Host "Press Enter to stop..."
[void][System.Console]::ReadLine()

Stop-Process -Id $api.Id, $web.Id -ErrorAction SilentlyContinue

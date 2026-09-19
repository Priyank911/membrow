$Root = Split-Path -Parent $PSScriptRoot
$WebPath = Join-Path $Root "apps/web"

Write-Host "Starting Membrow Vite Dev Server..."
$web = Start-Process -FilePath "npm.cmd" -ArgumentList "run", "dev" -WorkingDirectory $WebPath -PassThru

Write-Host "Waiting for Vite dev server on http://localhost:5173..."
Start-Sleep -Seconds 3

Write-Host "Launching Membrow Desktop Browser (Electron)..."
$electron = Start-Process -FilePath "npx.cmd" -ArgumentList "electron", "." -WorkingDirectory $WebPath -PassThru

Write-Host "Membrow Desktop Browser PID: $($electron.Id)"
Write-Host "Press Enter to exit..."
[void][System.Console]::ReadLine()

Stop-Process -Id $web.Id, $electron.Id -ErrorAction SilentlyContinue

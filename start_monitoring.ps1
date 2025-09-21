# SOA Team 10 - Monitoring Stack Startup Script
# PowerShell version with proper signal handling

Write-Host "Starting Docker services..." -ForegroundColor Green
docker-compose up -d --build

Write-Host "Starting WMI Exporter..." -ForegroundColor Green
$wmiProcess = Start-Process -FilePath ".\windows_exporter.exe" -ArgumentList "--collectors.enabled=cpu,memory,logical_disk,net,system" -PassThru -WindowStyle Hidden

Write-Host "`nAll services started!" -ForegroundColor Green
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   SOA Team 10 - Monitoring Stack" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   Frontend:     http://localhost:4200" -ForegroundColor Yellow
Write-Host "   API Gateway:  http://localhost:8080" -ForegroundColor Yellow
Write-Host "   Prometheus:   http://localhost:9090" -ForegroundColor Yellow
Write-Host "   Grafana:      http://localhost:3000" -ForegroundColor Yellow
Write-Host "   WMI Exporter: http://localhost:9182" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Press Ctrl+C to stop all services..." -ForegroundColor Red

# Register cleanup handler for Ctrl+C
$null = Register-EngineEvent PowerShell.Exiting -Action {
    Write-Host "`nStopping services..." -ForegroundColor Red
    
    # Stop WMI Exporter
    try {
        if ($wmiProcess -and !$wmiProcess.HasExited) {
            Write-Host "Stopping WMI Exporter..." -ForegroundColor Yellow
            $wmiProcess.Kill()
        }
    } catch {
        # Fallback: kill by process name
        Get-Process windows_exporter -ErrorAction SilentlyContinue | Stop-Process -Force
    }
    
    # Stop Docker services
    Write-Host "Stopping Docker services..." -ForegroundColor Yellow
    docker-compose down
    
    Write-Host "All services stopped!" -ForegroundColor Green
}

# Keep the script running
try {
    while ($true) {
        Start-Sleep -Seconds 5
        
        # Check if WMI exporter is still running
        if ($wmiProcess.HasExited) {
            Write-Host "WMI Exporter has stopped unexpectedly!" -ForegroundColor Red
            break
        }
    }
} catch {
    # Cleanup will be handled by the exit event
}
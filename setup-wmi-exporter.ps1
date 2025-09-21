# PowerShell script to install and configure WMI Exporter on Windows
# This will allow Prometheus to collect Windows host metrics

Write-Host "Setting up WMI Exporter for Windows..." -ForegroundColor Green

# Check if running as administrator
if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Host "This script requires Administrator privileges. Please run as Administrator." -ForegroundColor Red
    exit 1
}

# Download WMI Exporter
$wmiExporterUrl = "https://github.com/prometheus-community/windows_exporter/releases/latest/download/windows_exporter.exe"
$downloadPath = "$env:TEMP\windows_exporter.exe"
$installPath = "C:\Program Files\windows_exporter\windows_exporter.exe"

Write-Host "Downloading WMI Exporter..." -ForegroundColor Yellow
try {
    Invoke-WebRequest -Uri $wmiExporterUrl -OutFile $downloadPath -UseBasicParsing
    Write-Host "Download completed!" -ForegroundColor Green
} catch {
    Write-Host "Failed to download WMI Exporter: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Create installation directory
$installDir = "C:\Program Files\windows_exporter"
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force
}

# Copy executable
Copy-Item $downloadPath $installPath -Force

# Create Windows service
Write-Host "Creating Windows service..." -ForegroundColor Yellow
$serviceName = "windows_exporter"

# Remove existing service if it exists
if (Get-Service $serviceName -ErrorAction SilentlyContinue) {
    Stop-Service $serviceName -Force
    sc.exe delete $serviceName
}

# Create new service
$serviceArgs = "--collectors.enabled=cpu,memory,logical_disk,physical_disk,net,os,system,process"
sc.exe create $serviceName binpath= "`"$installPath`" $serviceArgs" start= auto

# Start the service
Start-Service $serviceName

# Check if service is running
$service = Get-Service $serviceName
if ($service.Status -eq "Running") {
    Write-Host "WMI Exporter service is running successfully!" -ForegroundColor Green
    Write-Host "Metrics are available at: http://localhost:9182/metrics" -ForegroundColor Cyan
} else {
    Write-Host "Failed to start WMI Exporter service" -ForegroundColor Red
}

# Add firewall rule to allow port 9182
Write-Host "Adding firewall rule for port 9182..." -ForegroundColor Yellow
try {
    New-NetFirewallRule -DisplayName "WMI Exporter" -Direction Inbound -Port 9182 -Protocol TCP -Action Allow -ErrorAction SilentlyContinue
    Write-Host "Firewall rule added successfully!" -ForegroundColor Green
} catch {
    Write-Host "Warning: Could not add firewall rule. You may need to manually allow port 9182." -ForegroundColor Yellow
}

Write-Host "`nSetup completed!" -ForegroundColor Green
Write-Host "You can now restart your Docker Compose stack to collect Windows host metrics." -ForegroundColor Cyan
Write-Host "Test the exporter at: http://localhost:9182/metrics" -ForegroundColor Cyan
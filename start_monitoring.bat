@echo off
echo Starting Docker services...
docker-compose up -d --build

echo Starting WMI Exporter...
start /B windows_exporter.exe --collectors.enabled="cpu,memory,logical_disk,net,system"

echo All services started!
echo Press Ctrl+C to stop all services...

:wait
timeout /t 5 >nul
goto wait

:cleanup
echo.
echo Stopping WMI Exporter...
taskkill /f /im windows_exporter.exe 2>nul
echo Stopping Docker services...
docker-compose down
echo All services stopped!
pause
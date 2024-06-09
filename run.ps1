$exePath = "C:\Users\Administrator\Desktop\Monitoring\app-amd64.exe"
$parameters = "<server-ip or hostname> <port>"

while ($true) {
    Start-Process -FilePath $exePath -ArgumentList $parameters -WindowStyle Hidden
    Start-Sleep -Seconds 30
}

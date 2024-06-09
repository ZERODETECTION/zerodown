# Path to the agent
$exePath = "C:\Users\Administrator\Desktop\Monitoring\app-amd64.exe"
$parameters = "<server> <port>"

# Trigger every 30 seconds
$trigger = New-ScheduledTaskTrigger -Once -At (Get-Date).AddSeconds(30) -RepetitionDuration ([TimeSpan]::MaxValue)

# Create action
$action = New-ScheduledTaskAction -Execute "$exePath" -Argument "$parameters"

# Create Task
Register-ScheduledTask -TaskName "AppMonitor" -Description "Monitoring-App ausführen alle 30 Sekunden" -Trigger $trigger -Action $action -RunLevel Highest

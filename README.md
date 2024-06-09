# zerodown
zerodown is a straightforward, clean, and minimalist uptime monitor for your infrastructure.

It has no complex dependencies, system requirements, or convoluted setups; just simple monitoring.

The code is written in Go and utilizes Bootstrap. It consists of an agent that sends monitored data over HTTP and JSON. The server monitors the clients and displays them in a list for easy overview.


<img src="https://raw.githubusercontent.com/DeinBenutzername/DeinRepositoryName/master/images/example.jpg" alt="Bildbeschreibung">


## Installation (server)
```
git clone https://github.com/ZERODETECTION/zerodown.git
```

```
sudo apt install golang-go
```

```
go mod init zerodown
```

```
go mod tidy
```

```
go run server.go
```

## Setup (agent)

```
zerodown_agent_amd64.exe <server-ip> <port>
```
You need to create a task or cronjob to execute the client every X seconds. Ensure that X is set to a value under 60 seconds; otherwise, the host will be marked as down.

## Dashboard

```
http://<ip>:8080/view
```

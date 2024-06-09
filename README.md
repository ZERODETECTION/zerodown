# zerodown
zerodown is a straightforward, clean, and minimalist uptime monitor for your infrastructure.

It has no complex dependencies, system requirements, or convoluted setups; just simple monitoring.

The code is written in Go and utilizes Bootstrap. It consists of an agent that sends monitored data over HTTP and JSON. The server monitors the clients and displays them in a list for easy overview.


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

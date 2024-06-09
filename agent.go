package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"math"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	gopsutilnet "github.com/shirou/gopsutil/net" // Alias the import to avoid conflict
)

type Stats struct {
	Timestamp        int64    `json:"timestamp"`
	CPUUsage         float64  `json:"cpu_usage"`
	RAMUsage         float64  `json:"ram_usage"`
	HDDUsage         uint64   `json:"hdd_usage"`
	NetworkSent      uint64   `json:"network_sent"`
	NetworkReceived  uint64   `json:"network_received"`
	Hostname         string   `json:"hostname"`
	OS               string   `json:"os"`
	IPAddresses      []string `json:"ip_addresses"`
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run main.go <server_address> <port>")
	}

	serverAddress := os.Args[1]
	port := os.Args[2]

	// Get current timestamp
	timestamp := time.Now().Unix()

	// CPU Auslastung abrufen (mit Wartezeit)
	time.Sleep(1 * time.Second) // Kurze Wartezeit, um die CPU-Nutzung genau zu messen
	cpuPercent, _ := cpu.Percent(0, false)
	cpuPercentInt := int(math.Round(cpuPercent[0]))

	// RAM-Auslastung auslesen
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Fehler beim Lesen der RAM-Auslastung:", err)
	}
	ramUsage := vmStat.UsedPercent

	// HDD-Auslastung auslesen
	diskStat, err := disk.Usage("/")
	if err != nil {
		log.Println("Fehler beim Lesen der HDD-Auslastung:", err)
	}
	hddUsage := uint64(diskStat.UsedPercent)

	// Netzwerkstatistiken auslesen
	netIO, err := gopsutilnet.IOCounters(false) // Use the aliased import
	if err != nil {
		log.Println("Fehler beim Lesen der Netzwerkstatistiken:", err)
	}

	// Hostname auslesen
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Fehler beim Lesen des Hostnamens:", err)
	}

	// Betriebssystem auslesen
	osStat, err := host.Info()
	if err != nil {
		log.Println("Fehler beim Lesen des Betriebssystems:", err)
	}

	// IP-Adressen abrufen
	ipAddresses, err := getLocalIPs()
	if err != nil {
		log.Println("Fehler beim Abrufen der IP-Adressen:", err)
	}

	// Create a struct to hold all stats
	stats := Stats{
		Timestamp:        timestamp,
		CPUUsage:         float64(cpuPercentInt), // Hier wird der Integer-Wert in einen float64-Wert umgewandelt
		RAMUsage:         ramUsage,
		HDDUsage:         hddUsage,
		NetworkSent:      netIO[0].BytesSent,
		NetworkReceived:  netIO[0].BytesRecv,
		Hostname:         hostname,
		OS:               osStat.OS,
		IPAddresses:      ipAddresses,
	}

	// Log the stats
	log.Printf("Collected stats: %+v\n", stats)

	// Convert the stats to JSON
	jsonData, err := json.Marshal(stats)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Send the data to the server
	resp, err := http.Post("http://" + serverAddress + ":" + port + "/stats", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to send data: %v", err)
	}
	defer resp.Body.Close()

	log.Println("Data sent successfully:", resp.Status)
}

func getLocalIPs() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		// Check if the address is not loopback and is an IP address
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}

	return ips, nil
}

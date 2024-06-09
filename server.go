package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Struct to hold the system stats
type SystemStats struct {
	Timestamp       int64   `json:"timestamp"`
	CPUUsage        float64 `json:"cpu_usage"`
	RAMUsage        float64 `json:"ram_usage"`
	HDDUsage        float64 `json:"hdd_usage"`
	NetworkSent     int64   `json:"network_sent"`
	NetworkReceived int64   `json:"network_received"`
	Hostname        string  `json:"hostname"`
	OS              string  `json:"os"`
	State           string  `json:"state"` // New field for state
	IPAddresses     []string `json:"ip_addresses"` // New field for IP addresses
}

var (
	statsList       []SystemStats
	mu              sync.Mutex
	lastSentTimeMap map[string]time.Time // Map to store last sent time of each host
)

func init() {
	lastSentTimeMap = make(map[string]time.Time)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var stats SystemStats
		if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		
		// Update existing stats or add new stats
		found := false
		for i, existingStats := range statsList {
			if existingStats.Hostname == stats.Hostname {
				statsList[i] = stats // Update existing stats
				statsList[i].State = "up" // Set state to up on receiving new stats
				found = true
				break
			}
		}
		if !found {
			stats.State = "up" // Set state to up for new stats
			statsList = append(statsList, stats)
		}
		lastSentTimeMap[stats.Hostname] = time.Now() // Update last sent time
		log.Printf("Received stats from %s: %+v\n", stats.Hostname, stats)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func checkInactiveHosts() {
	for {
		time.Sleep(60 * time.Second) // Check every 60 seconds
		mu.Lock()
		for i, stats := range statsList {
			if time.Since(lastSentTimeMap[stats.Hostname]) > 60*time.Second {
				log.Printf("Alert: Host %s has been inactive for more than 60 seconds", stats.Hostname)
				statsList[i].State = "down" // Set state to down if inactive
				// Trigger alert mechanism here (send email, notification, etc.)
			} else {
				statsList[i].State = "up" // Ensure state is up if active
			}
		}
		mu.Unlock()
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>System Stats</title>
	<link href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css" rel="stylesheet">
	<link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.6.1/font/bootstrap-icons.css" rel="stylesheet">
	<style>
		.state-up {
			color: green;
			font-size: 14px; /* You can specify the size in pixels, ems, rems, etc. */
			font-weight: bold;
		}
		.state-down {
			color: red;
			font-size: 14px; /* You can specify the size in pixels, ems, rems, etc. */
			font-weight: bold;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1 class="mt-5">System Stats</h1>
		<p>Total number of hosts: {{.TotalHosts}}</p>
		<form action="/view" method="GET" class="form-inline mb-3">
			<div class="form-group mr-2">
				<label for="hostname" class="mr-2">Filter by Hostname:</label>
				<input type="text" id="hostname" name="hostname" class="form-control">
			</div>
			<button type="submit" class="btn btn-primary">Filter</button>
		</form>
		<table class="table table-striped mt-3">
			<thead>
				<tr>
					<th>Timestamp</th>
					<th>Hostname</th>
					<th>IP Addresses</th>
					<th>OS</th>
					<th>CPU Usage</th>
					<th>RAM Usage</th>
					<th>HDD Usage</th>
					<th>Network Sent</th>
					<th>Network Received</th>
					<th>State</th>
				</tr>
			</thead>
			<tbody>
				{{range .Stats}}
				<tr>
					<td>{{.Timestamp}}</td>
					<td>{{.Hostname}}</td>
					<td>{{range .IPAddresses}}{{.}} {{end}}</td>
					<td>{{if eq .OS "windows"}}<i class="bi bi-windows"></i>{{else}}{{.OS}}{{end}}</td>
					<td>
						<div class="progress">
							<div class="progress-bar" role="progressbar" style="width: {{.CPUUsage}}%;" aria-valuenow="{{.CPUUsage}}" aria-valuemin="0" aria-valuemax="100">{{.CPUUsage}}%</div>
						</div>
					</td>
					<td>
						<div class="progress">
							<div class="progress-bar" role="progressbar" style="width: {{.RAMUsage}}%;" aria-valuenow="{{.RAMUsage}}" aria-valuemin="0" aria-valuemax="100">{{.RAMUsage}}%</div>
						</div>
					</td>
					<td>
						<div class="progress">
							<div class="progress-bar" role="progressbar" style="width: {{.HDDUsage}}%;" aria-valuenow="{{.HDDUsage}}" aria-valuemin="0" aria-valuemax="100">{{.HDDUsage}}%</div>
						</div>
					</td>
					<td>{{.NetworkSent}}</td>
					<td>{{.NetworkReceived}}</td>
					<td class="{{if eq .State "up"}}state-up{{else}}state-down{{end}}">{{.State}}</td>
				</tr>
				{{end}}
			</tbody>
		</table>
	</div>
</body>
</html>
	`
	t := template.Must(template.New("stats").Parse(tmpl))

	mu.Lock()
	defer mu.Unlock()

	// Sort statsList by Timestamp descending
	sort.Slice(statsList, func(i, j int) bool {
		return statsList[i].Timestamp > statsList[j].Timestamp
	})

	// Filter by Hostname if present in the URL query
	if hostname := r.URL.Query().Get("hostname"); hostname != "" {
		var filteredStatsList []SystemStats
		for _, stats := range statsList {
			if stats.Hostname == hostname {
				filteredStatsList = append(filteredStatsList, stats)
			}
		}
		// Calculate the total number of unique hosts
		uniqueHosts := make(map[string]bool)
		for _, stats := range filteredStatsList {
			uniqueHosts[stats.Hostname] = true
		}
		totalHosts := len(uniqueHosts)

		// Add the total number of unique hosts to the data model
		data := struct {
			TotalHosts int
			Stats      []SystemStats
		}{
			TotalHosts: totalHosts,
			Stats:      filteredStatsList,
		}

		if err := t.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Otherwise, display statsList without filtering
	// Calculate the total number of unique hosts
	uniqueHosts := make(map[string]bool)
	for _, stats := range statsList {
		uniqueHosts[stats.Hostname] = true
	}
	totalHosts := len(uniqueHosts)

	// Add the total number of unique hosts to the data model
	data := struct {
		TotalHosts int
		Stats      []SystemStats
	}{
		TotalHosts: totalHosts,
		Stats:      statsList,
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Set up HTTP handlers
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/view", viewHandler)

	go checkInactiveHosts() // Start the goroutine to check inactive hosts

	log.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

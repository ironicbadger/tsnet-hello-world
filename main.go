package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"tailscale.com/tsnet"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var (
	hostname     = envOrDefault("TS_HOSTNAME", "ts-hello-world")
	stateDir     = envOrDefault("TS_STATE_DIR", "/var/lib/ts-hello-world")
	authKey      = os.Getenv("TS_AUTHKEY")
	enableServe  = envOrDefault("TS_ENABLE_SERVE", "false") == "true"
	enableFunnel = envOrDefault("TS_ENABLE_FUNNEL", "false") == "true"
)

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type PageData struct {
	Title        string
	UserLogin    string
	UserName     string
	NodeName     string
	NodeID       string
	TailnetName  string
	RemoteAddr   string
	EnableServe  bool
	EnableFunnel bool
	Headers      map[string]string
	Nodes        []NodeInfo
}

type NodeInfo struct {
	Name      string
	ID        string
	Addresses []string
	Tags      []string
	Online    bool
}

func main() {
	flag.Parse()

	srv := &tsnet.Server{
		Hostname: hostname,
		Dir:      stateDir,
		AuthKey:  authKey,
	}

	if authKey == "" {
		log.Println("Warning: TS_AUTHKEY not set. The server will print an auth URL to complete authentication.")
	}

	log.Printf("Starting tsnet server with hostname: %s", hostname)
	log.Printf("State directory: %s", stateDir)
	log.Printf("Serve enabled: %v", enableServe)
	log.Printf("Funnel enabled: %v", enableFunnel)

	var ln net.Listener
	var err error

	listenAddr := ":443"

	if enableFunnel {
		log.Printf("Starting with Funnel (publicly accessible via Tailscale)")
		ln, err = srv.ListenFunnel("tcp", listenAddr)
	} else if enableServe {
		log.Printf("Starting with HTTPS on port 443 (tailnet only)")
		ln, err = srv.ListenTLS("tcp", listenAddr)
	} else {
		log.Printf("Starting on port 443 (direct connection)")
		ln, err = srv.Listen("tcp", listenAddr)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	lc, err := srv.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))

	// Load templates
	tmpl := template.Must(template.ParseFS(templateFS, "templates/*.html"))

	// Main page handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		who, err := lc.WhoIs(r.Context(), r.RemoteAddr)
		if err != nil {
			http.Error(w, fmt.Sprintf("WhoIs error: %v", err), 500)
			return
		}

		status, err := lc.Status(r.Context())
		if err != nil {
			log.Printf("Error getting status: %v", err)
		}

		var nodes []NodeInfo
		if status != nil {
			// Add self node first
			if status.Self != nil {
				self := NodeInfo{
					Name:   status.Self.HostName + " (this node)",
					ID:     string(status.Self.ID),
					Online: status.Self.Online,
				}
				for _, addr := range status.Self.TailscaleIPs {
					self.Addresses = append(self.Addresses, addr.String())
				}
				// Add tags if present
				if status.Self.Tags != nil {
					for i := 0; i < status.Self.Tags.Len(); i++ {
						self.Tags = append(self.Tags, status.Self.Tags.At(i))
					}
				}
				nodes = append(nodes, self)
			}

			// Add peer nodes
			if status.Peer != nil {
				for _, peer := range status.Peer {
					// Skip Mullvad exit nodes
					if strings.HasSuffix(peer.DNSName, ".mullvad.ts.net.") {
						continue
					}
					node := NodeInfo{
						Name:   peer.HostName,
						ID:     string(peer.ID),
						Online: peer.Online,
					}
					for _, addr := range peer.TailscaleIPs {
						node.Addresses = append(node.Addresses, addr.String())
					}
					// Add tags if present
					if peer.Tags != nil {
						for i := 0; i < peer.Tags.Len(); i++ {
							node.Tags = append(node.Tags, peer.Tags.At(i))
						}
					}
					nodes = append(nodes, node)
				}
			}
		}

		// Extract Tailscale headers
		headers := make(map[string]string)
		for key, values := range r.Header {
			if len(key) >= 9 && key[:9] == "Tailscale" {
				headers[key] = values[0]
			}
		}

		data := PageData{
			Title:        "Tailscale tsnet Demo",
			UserLogin:    who.UserProfile.LoginName,
			UserName:     who.UserProfile.DisplayName,
			NodeName:     who.Node.ComputedName,
			NodeID:       who.Node.ID.String(),
			TailnetName:  who.Node.Name,
			RemoteAddr:   r.RemoteAddr,
			EnableServe:  enableServe,
			EnableFunnel: enableFunnel,
			Headers:      headers,
			Nodes:        nodes,
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	// Status API endpoint
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		status, err := lc.Status(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Status error: %v", err), 500)
			return
		}

		response := map[string]interface{}{
			"hostname":       hostname,
			"serve_enabled":  enableServe,
			"funnel_enabled": enableFunnel,
			"state_dir":      stateDir,
		}

		if status != nil {
			response["version"] = status.Version
			response["backend_state"] = status.BackendState
			
			if status.Self != nil {
				response["self"] = map[string]interface{}{
					"name":      status.Self.HostName,
					"id":        string(status.Self.ID),
					"addresses": status.Self.TailscaleIPs,
					"online":    status.Self.Online,
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Printf("Server listening on %s", ln.Addr())
	log.Fatal(http.Serve(ln, mux))
}
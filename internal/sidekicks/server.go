package sidekicks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const defaultServerPort = 7777

// HTTPServer handles webhook requests for sidekicks.
type HTTPServer struct {
	root       string
	port       int
	routeTable *RouteTable
	server     *http.Server
}

// NewHTTPServer creates a new sidekick HTTP server.
func NewHTTPServer(root string, routeTable *RouteTable) *HTTPServer {
	s := &HTTPServer{
		root:       root,
		port:       defaultServerPort,
		routeTable: routeTable,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/status", s.handleStatus)
	mux.HandleFunc("/outbox", s.handleOutbox)
	mux.HandleFunc("/complete", s.handleComplete)
	mux.HandleFunc("/send-report", s.handleSendReport)
	mux.HandleFunc("/", s.handleWebhook)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	return s
}

// Start starts the HTTP server.
func (s *HTTPServer) Start() error {
	log.Printf("Sidekick server listening on :%d", s.port)
	return s.server.ListenAndServe()
}

// Stop stops the HTTP server.
func (s *HTTPServer) Stop() error {
	return s.server.Close()
}

func (s *HTTPServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	active, _ := ListActive(s.root)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"active":  len(active),
		"routes":  s.routeTable.ListRoutes(),
	})
}

func (s *HTTPServer) handleOutbox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}

	var msg struct {
		Message string `json:"message"`
		Channel string `json:"channel"`
	}
	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "parse body", http.StatusBadRequest)
		return
	}

	log.Printf("Outbox message: %s (channel: %s)", msg.Message, msg.Channel)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

func (s *HTTPServer) handleComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}

	var msg struct {
		Summary string `json:"summary"`
	}
	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "parse body", http.StatusBadRequest)
		return
	}

	if err := Complete(s.root, msg.Summary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "completed"})
}

func (s *HTTPServer) handleSendReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}

	var msg struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "parse body", http.StatusBadRequest)
		return
	}

	Report(s.root, msg.Message)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

func (s *HTTPServer) handleWebhook(w http.ResponseWriter, r *http.Request) {
	skill, handler, ok := s.routeTable.Match(r.URL.Path)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Verify HMAC if secret is configured
	secret := s.routeTable.GetSecret(r.URL.Path)
	if secret != "" {
		body, _ := io.ReadAll(r.Body)
		sig := r.Header.Get("X-Hub-Signature-256")
		if sig == "" {
			sig = r.Header.Get("X-Signature-256")
		}
		if !VerifyHMAC(string(body), sig, secret) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
	}

	log.Printf("Webhook received: %s → %s/%s", r.URL.Path, skill, handler)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "received"})
}

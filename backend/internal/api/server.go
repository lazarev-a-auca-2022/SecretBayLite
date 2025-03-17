// Package api provides HTTP server and request handling functionality.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"secretbay/internal/auth"
	"secretbay/internal/config"
	"secretbay/internal/vpn"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server and its dependencies
type Server struct {
	config *config.Config
	logger *log.Logger
	vpn    *vpn.Service
	auth   *auth.Auth
}

// NewServer creates a new API server instance
func NewServer(cfg *config.Config, logger *log.Logger) *Server {
	return &Server{
		config: cfg,
		logger: logger,
		vpn:    vpn.NewService(logger.Writer()),
		auth:   auth.NewAuth(cfg.JWTSecret),
	}
}

// Router returns the configured HTTP router
func (s *Server) Router() *mux.Router {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(s.authMiddleware)

	api.HandleFunc("/vpn/configure", s.handleVPNConfiguration).Methods(http.MethodPost)
	api.HandleFunc("/health", s.handleHealthCheck).Methods(http.MethodGet)

	return r
}

// handleVPNConfiguration processes VPN configuration requests
func (s *Server) handleVPNConfiguration(w http.ResponseWriter, r *http.Request) {
	var req vpn.ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate request
	if err := s.validateVPNRequest(&req); err != nil {
		s.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Process VPN configuration
	resp, err := s.vpn.ConfigureVPN(&req)
	if err != nil {
		s.logger.Printf("VPN configuration failed: %v", err)
		s.respondWithError(w, http.StatusInternalServerError, "Failed to configure VPN")
		return
	}

	s.respondWithJSON(w, http.StatusOK, resp)
}

// handleHealthCheck handles health check requests
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	s.respondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// validateVPNRequest validates the VPN configuration request
func (s *Server) validateVPNRequest(req *vpn.ConfigRequest) error {
	if req.ServerIP == "" {
		return fmt.Errorf("server IP is required")
	}
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	if req.AuthMethod != "password" && req.AuthMethod != "key" {
		return fmt.Errorf("invalid authentication method")
	}
	if req.VPNType != "openvpn" && req.VPNType != "ios" {
		return fmt.Errorf("invalid VPN type")
	}
	if req.AuthCredential == "" {
		return fmt.Errorf("authentication credential is required")
	}
	return nil
}

// respondWithError sends an error response
func (s *Server) respondWithError(w http.ResponseWriter, code int, message string) {
	s.respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sends a JSON response
func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// authMiddleware handles JWT authentication
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health check
		if r.URL.Path == "/api/health" {
			next.ServeHTTP(w, r)
			return
		}

		token, err := auth.ExtractTokenFromHeader(r)
		if err != nil {
			s.respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		claims, err := s.auth.ValidateToken(token)
		if err != nil {
			s.respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add claims to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

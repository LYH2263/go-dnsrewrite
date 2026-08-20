package api

import (
	"net/http"
	"strings"

	"example.com/dnsrewrite"
)

func (s *Server) handleListRules(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.eng.ListRules())
}

func (s *Server) handleAddRule(w http.ResponseWriter, r *http.Request) {
	var spec dnsrewrite.RuleSpec
	if err := readJSON(r, &spec); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	spec.Enabled = true
	if err := s.eng.AddRule(spec); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]string{"ok": "1"})
}

func (s *Server) handleDeleteRule(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/rules/")
	if id == "" {
		writeJSON(w, 400, map[string]string{"error": "missing id"})
		return
	}
	if r.Method != http.MethodDelete {
		w.WriteHeader(405)
		return
	}
	if err := s.eng.RemoveRule(id); err != nil {
		writeJSON(w, 404, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]string{"ok": "1"})
}

func (s *Server) handleTry(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
		Type uint16 `json:"type"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	if body.Type == 0 {
		body.Type = 1
	}
	res, err := s.eng.TryResolve(r.Context(), body.Name, dnsrewrite.RRType(body.Type))
	if err != nil {
		writeJSON(w, 502, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, res)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.eng.Stats())
}

func (s *Server) handleEnable(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	en := r.URL.Query().Get("enabled") == "1" || r.URL.Query().Get("enabled") == "true"
	if err := s.eng.EnableRule(id, en); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"id": id, "enabled": en})
}

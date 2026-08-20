package api

import (
	"net/http"

	"example.com/dnsrewrite"
)

// Options HTTP 服务选项。
type Options struct {
	WebDir    string
	AllowCORS bool
}

// Server 管理 API。
type Server struct {
	eng *dnsrewrite.Engine
	mux *http.ServeMux
	opt Options
}

// New 构造。
func New(eng *dnsrewrite.Engine, opt Options) http.Handler {
	s := &Server{eng: eng, mux: http.NewServeMux(), opt: opt}
	s.routes()
	var h http.Handler = s.mux
	if opt.AllowCORS {
		h = withCORS(h)
	}
	return h
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/rules", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleListRules(w, r)
		case http.MethodPost:
			s.handleAddRule(w, r)
		default:
			w.WriteHeader(405)
		}
	})
	s.mux.HandleFunc("/api/rules/", s.handleDeleteRule)
	s.mux.HandleFunc("/api/try", s.handleTry)
	s.mux.HandleFunc("/api/stats", s.handleStats)
	s.mux.HandleFunc("/api/enable", s.handleEnable)
	s.mux.Handle("/", staticHandler(s.opt.WebDir))
}

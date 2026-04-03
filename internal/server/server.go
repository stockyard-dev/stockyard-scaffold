package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-scaffold/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/templates",s.list);s.mux.HandleFunc("POST /api/templates",s.create);s.mux.HandleFunc("GET /api/templates/{id}",s.get);s.mux.HandleFunc("DELETE /api/templates/{id}",s.del)
s.mux.HandleFunc("POST /api/templates/{id}/generate",s.generate)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"templates":oe(s.db.List())})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){var t store.Template;json.NewDecoder(r.Body).Decode(&t);if t.Name==""{we(w,400,"name required");return};s.db.Create(&t);wj(w,201,s.db.Get(t.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){t:=s.db.Get(r.PathValue("id"));if t==nil{we(w,404,"not found");return};wj(w,200,t)}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)generate(w http.ResponseWriter,r *http.Request){var req struct{Variables map[string]string `json:"variables"`};json.NewDecoder(r.Body).Decode(&req);files:=s.db.Generate(r.PathValue("id"),req.Variables);if files==nil{we(w,404,"not found");return};wj(w,200,map[string]any{"files":files})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]int{"templates":s.db.Count()})}
func(s *Server)health(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"status":"ok","service":"scaffold","templates":s.db.Count()})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}

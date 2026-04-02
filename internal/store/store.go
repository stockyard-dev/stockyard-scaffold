package store
import ("database/sql";"encoding/json";"fmt";"os";"path/filepath";"strings";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Template struct{ID string `json:"id"`;Name string `json:"name"`;Description string `json:"description,omitempty"`;Language string `json:"language,omitempty"`;Files []TemplateFile `json:"files"`;Variables []string `json:"variables,omitempty"`;CreatedAt string `json:"created_at"`;UseCount int `json:"use_count"`}
type TemplateFile struct{Path string `json:"path"`;Content string `json:"content"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"scaffold.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS templates(id TEXT PRIMARY KEY,name TEXT NOT NULL,description TEXT DEFAULT '',language TEXT DEFAULT '',files_json TEXT DEFAULT '[]',variables_json TEXT DEFAULT '[]',use_count INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(t *Template)error{t.ID=genID();t.CreatedAt=now();if t.Files==nil{t.Files=[]TemplateFile{}};if t.Variables==nil{t.Variables=[]string{}}
fj,_:=json.Marshal(t.Files);vj,_:=json.Marshal(t.Variables)
_,err:=d.db.Exec(`INSERT INTO templates(id,name,description,language,files_json,variables_json,created_at)VALUES(?,?,?,?,?,?,?)`,t.ID,t.Name,t.Description,t.Language,string(fj),string(vj),t.CreatedAt);return err}
func(d *DB)Get(id string)*Template{var t Template;var fj,vj string
if d.db.QueryRow(`SELECT id,name,description,language,files_json,variables_json,use_count,created_at FROM templates WHERE id=?`,id).Scan(&t.ID,&t.Name,&t.Description,&t.Language,&fj,&vj,&t.UseCount,&t.CreatedAt)!=nil{return nil}
json.Unmarshal([]byte(fj),&t.Files);json.Unmarshal([]byte(vj),&t.Variables);return &t}
func(d *DB)List()[]Template{rows,_:=d.db.Query(`SELECT id,name,description,language,files_json,variables_json,use_count,created_at FROM templates ORDER BY name`);if rows==nil{return nil};defer rows.Close()
var o []Template;for rows.Next(){var t Template;var fj,vj string;rows.Scan(&t.ID,&t.Name,&t.Description,&t.Language,&fj,&vj,&t.UseCount,&t.CreatedAt);json.Unmarshal([]byte(fj),&t.Files);json.Unmarshal([]byte(vj),&t.Variables);o=append(o,t)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM templates WHERE id=?`,id);return err}
func(d *DB)Generate(id string,vars map[string]string)[]TemplateFile{t:=d.Get(id);if t==nil{return nil};d.db.Exec(`UPDATE templates SET use_count=use_count+1 WHERE id=?`,id)
var out []TemplateFile;for _,f:=range t.Files{path:=f.Path;content:=f.Content;for k,v:=range vars{path=strings.ReplaceAll(path,"{{"+k+"}}",v);content=strings.ReplaceAll(content,"{{"+k+"}}",v)};out=append(out,TemplateFile{Path:path,Content:content})};return out}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM templates`).Scan(&n);return n}

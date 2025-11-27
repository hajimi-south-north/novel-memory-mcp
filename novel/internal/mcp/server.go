package mcp

import (
	"bufio"
	"encoding/json"
	"mcpnovel/internal/conflict"
	"mcpnovel/internal/helpers"
	"mcpnovel/internal/models"
	"mcpnovel/internal/outline"
	"mcpnovel/internal/storage"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Server struct {
    DB        *gorm.DB
    Services  *helpers.Services
    Detector  *conflict.Detector
    Generator *outline.Generator
    DBPath    string
}

type jsonrpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type jsonrpcResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      any           `json:"id"`
	Result  any           `json:"result,omitempty"`
	Error   *jsonrpcError `json:"error,omitempty"`
}

type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

func Run() {
    p := os.Getenv("MCP_NOVEL_DB")
    if p == "" {
        wd, _ := os.Getwd()
        p = wd + "/novel.db"
    }
    db, err := storage.Open(p)
    if err != nil {
        os.Exit(1)
    }
    autoMigrate(db)
    s := &Server{DB: db, DBPath: p}
    s.Services = &helpers.Services{DB: db}
    s.Detector = &conflict.Detector{DB: db}
    s.Generator = &outline.Generator{DB: db}
    s.loop()
}

func (s *Server) loop() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var req jsonrpcRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			resp := jsonrpcResponse{JSONRPC: "2.0", ID: nil, Error: &jsonrpcError{Code: -32700, Message: "Parse error"}}
			writeJSON(w, resp)
			continue
		}
		s.handle(w, req)
	}
}

func (s *Server) handle(w *bufio.Writer, req jsonrpcRequest) {
	switch req.Method {
	case "initialize":
		writeJSON(w, jsonrpcResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools": map[string]any{
					"list": true,
					"call": true,
				},
			},
			"serverInfo": map[string]any{
				"name":    "mcp-novel",
				"version": "0.1.0",
			},
		}})
    case "tools/list":
        writeJSON(w, jsonrpcResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{"tools": s.tools()}})
	case "tools/call":
		var p map[string]any
		_ = json.Unmarshal(req.Params, &p)
		name := stringField(p, "name")
		argsRaw := p["arguments"]
		var args map[string]any
		if argsRaw != nil {
			b, _ := json.Marshal(argsRaw)
			_ = json.Unmarshal(b, &args)
		}
		res, err := s.call(name, args)
		if err != nil {
			writeJSON(w, jsonrpcResponse{JSONRPC: "2.0", ID: req.ID, Error: &jsonrpcError{Code: -32000, Message: err.Error()}})
			return
		}
		writeJSON(w, jsonrpcResponse{JSONRPC: "2.0", ID: req.ID, Result: res})
	default:
		writeJSON(w, jsonrpcResponse{JSONRPC: "2.0", ID: req.ID, Error: &jsonrpcError{Code: -32601, Message: "Method not found"}})
	}
}

func writeJSON(w *bufio.Writer, v any) {
	b, _ := json.Marshal(v)
	w.WriteString(string(b))
	w.WriteByte('\n')
	w.Flush()
}

func (s *Server) tools() []Tool {
    return []Tool{
        {Name: "dbHelper", Description: "数据库管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "path": map[string]any{"type": "string"}}}},
        {Name: "sqlHelper", Description: "SQL操作", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"entity": map[string]any{"type": "string"}, "action": map[string]any{"type": "string"}, "data": map[string]any{"type": "object"}}}},
        {Name: "novelHelper", Description: "小说管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "title": map[string]any{"type": "string"}, "description": map[string]any{"type": "string"}, "id": map[string]any{"type": "number"}}}},
        {Name: "volumeHelper", Description: "分卷管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "novelID": map[string]any{"type": "number"}, "title": map[string]any{"type": "string"}, "index": map[string]any{"type": "number"}}}},
        {Name: "chapterHelper", Description: "章节管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "volumeID": map[string]any{"type": "number"}, "title": map[string]any{"type": "string"}, "index": map[string]any{"type": "number"}, "status": map[string]any{"type": "string"}, "id": map[string]any{"type": "number"}, "content": map[string]any{"type": "string"}}}},
        {Name: "eventHelper", Description: "事件管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "chapterID": map[string]any{"type": "number"}, "worldID": map[string]any{"type": "number"}, "locationID": map[string]any{"type": "number"}, "timeSegmentID": map[string]any{"type": "number"}, "description": map[string]any{"type": "string"}, "characters": map[string]any{"type": "array", "items": map[string]any{"type": "number"}}, "items": map[string]any{"type": "array", "items": map[string]any{"type": "number"}}}}},
        {Name: "worldHelper", Description: "世界管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "name": map[string]any{"type": "string"}, "description": map[string]any{"type": "string"}}}},
        {Name: "periodHelper", Description: "时期管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "worldID": map[string]any{"type": "number"}, "name": map[string]any{"type": "string"}, "index": map[string]any{"type": "number"}}}},
        {Name: "timeSegmentHelper", Description: "时间段管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "periodID": map[string]any{"type": "number"}, "name": map[string]any{"type": "string"}, "start": map[string]any{"type": "string"}, "end": map[string]any{"type": "string"}}}},
        {Name: "characterHelper", Description: "人物管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "name": map[string]any{"type": "string"}, "bio": map[string]any{"type": "string"}}}},
        {Name: "characterRelationshipHelper", Description: "人物关系管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "aid": map[string]any{"type": "number"}, "bid": map[string]any{"type": "number"}, "type": map[string]any{"type": "string"}, "intimacy": map[string]any{"type": "number"}}}},
        {Name: "locationHelper", Description: "地点管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "worldID": map[string]any{"type": "number"}, "name": map[string]any{"type": "string"}, "description": map[string]any{"type": "string"}}}},
        {Name: "itemHelper", Description: "物品管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "name": map[string]any{"type": "string"}, "ownerID": map[string]any{"type": "number"}, "locationID": map[string]any{"type": "number"}, "status": map[string]any{"type": "string"}, "itemID": map[string]any{"type": "number"}, "fromID": map[string]any{"type": "number"}, "toID": map[string]any{"type": "number"}, "eventID": map[string]any{"type": "number"}}}},
        {Name: "characterAbilityHelper", Description: "人物能力管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "characterID": map[string]any{"type": "number"}, "name": map[string]any{"type": "string"}, "level": map[string]any{"type": "number"}, "abilityID": map[string]any{"type": "number"}, "eventID": map[string]any{"type": "number"}, "note": map[string]any{"type": "string"}}}},
        {Name: "plotThreadHelper", Description: "情节线索管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "novelID": map[string]any{"type": "number"}, "name": map[string]any{"type": "string"}, "stage": map[string]any{"type": "string"}, "plotID": map[string]any{"type": "number"}}}},
        {Name: "characterMemoryHelper", Description: "人物记忆管理", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "characterID": map[string]any{"type": "number"}, "eventID": map[string]any{"type": "number"}, "content": map[string]any{"type": "string"}, "trigger": map[string]any{"type": "string"}}}},
        {Name: "conflictDetectionHelper", Description: "冲突检测", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}}}},
        {Name: "outlineGeneratorHelper", Description: "纲要生成", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "id": map[string]any{"type": "number"}}}},
        {Name: "articleExportHelper", Description: "文章导出", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "id": map[string]any{"type": "number"}}}},
        {Name: "styleHelper", Description: "文笔风格参考", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"action": map[string]any{"type": "string"}, "novelID": map[string]any{"type": "number"}, "content": map[string]any{"type": "string"}}}},
    }
}

func (s *Server) call(name string, args map[string]any) (any, error) {
	switch name {
    case "dbHelper":
        act := stringField(args, "action")
        if act == "init" {
            path := stringField(args, "path")
            if path != "" {
                ndb, err := storage.Open(path)
                if err != nil { return nil, err }
                autoMigrate(ndb)
                s.DB = ndb
                s.DBPath = path
                s.Services = &helpers.Services{DB: ndb}
                s.Detector = &conflict.Detector{DB: ndb}
                s.Generator = &outline.Generator{DB: ndb}
            } else {
                autoMigrate(s.DB)
            }
            return map[string]any{"ok": true, "path": s.DBPath}, nil
        }
        if act == "export" {
            return map[string]any{"path": s.DBPath}, nil
        }
	case "novelHelper":
		act := stringField(args, "action")
		if act == "create" {
			n, err := s.Services.CreateNovel(stringField(args, "title"), stringField(args, "description"))
			if err != nil {
				return nil, err
			}
			return n, nil
		}
		if act == "export" {
			id := uintField(args, "id")
			res, err := s.Services.ExportNovel(id)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
		if act == "outline" {
			id := uintField(args, "id")
			o, err := s.Generator.NovelOutline(id)
			if err != nil {
				return nil, err
			}
			return map[string]any{"outline": o}, nil
		}
	case "volumeHelper":
		v, err := s.Services.CreateVolume(uintField(args, "novelID"), stringField(args, "title"), intField(args, "index"))
		if err != nil {
			return nil, err
		}
		return v, nil
	case "chapterHelper":
		act := stringField(args, "action")
		if act == "create" {
			c, err := s.Services.CreateChapter(uintField(args, "volumeID"), stringField(args, "title"), intField(args, "index"), stringField(args, "status"))
			if err != nil {
				return nil, err
			}
			return c, nil
		}
		if act == "update" {
			err := s.Services.UpsertChapterContent(uintField(args, "id"), stringField(args, "content"))
			if err != nil {
				return nil, err
			}
			return map[string]any{"ok": true}, nil
		}
		if act == "export" {
			res, err := s.Services.ExportChapter(uintField(args, "id"))
			if err != nil {
				return nil, err
			}
			return res, nil
		}
		if act == "outline" {
			o, err := s.Generator.ChapterOutline(uintField(args, "id"))
			if err != nil {
				return nil, err
			}
			return map[string]any{"outline": o}, nil
		}
	case "eventHelper":
		chars := uintSliceField(args, "characters")
		items := uintSliceField(args, "items")
		e, err := s.Services.CreateEvent(uintField(args, "chapterID"), uintField(args, "worldID"), uintField(args, "locationID"), uintField(args, "timeSegmentID"), stringField(args, "description"), chars, items)
		if err != nil {
			return nil, err
		}
		return e, nil
	case "worldHelper":
		w, err := s.Services.CreateWorld(stringField(args, "name"), stringField(args, "description"))
		if err != nil {
			return nil, err
		}
		return w, nil
	case "periodHelper":
		p, err := s.Services.CreatePeriod(uintField(args, "worldID"), stringField(args, "name"), intField(args, "index"))
		if err != nil {
			return nil, err
		}
		return p, nil
	case "timeSegmentHelper":
		start, _ := time.Parse(time.RFC3339, stringField(args, "start"))
		end, _ := time.Parse(time.RFC3339, stringField(args, "end"))
		ts, err := s.Services.CreateTimeSegment(uintField(args, "periodID"), stringField(args, "name"), start, end)
		if err != nil {
			return nil, err
		}
		return ts, nil
	case "characterHelper":
		c, err := s.Services.CreateCharacter(stringField(args, "name"), stringField(args, "bio"))
		if err != nil {
			return nil, err
		}
		return c, nil
	case "characterRelationshipHelper":
		r, err := s.Services.SetCharacterRelationship(uintField(args, "aid"), uintField(args, "bid"), stringField(args, "type"), floatField(args, "intimacy"))
		if err != nil {
			return nil, err
		}
		return r, nil
	case "locationHelper":
		l, err := s.Services.CreateLocation(uintField(args, "worldID"), stringField(args, "name"), stringField(args, "description"))
		if err != nil {
			return nil, err
		}
		return l, nil
	case "itemHelper":
		act := stringField(args, "action")
		if act == "create" {
			it, err := s.Services.CreateItem(stringField(args, "name"), uintField(args, "ownerID"), uintField(args, "locationID"), stringField(args, "status"))
			if err != nil {
				return nil, err
			}
			return it, nil
		}
		if act == "transfer" {
			t, err := s.Services.TransferItem(uintField(args, "itemID"), uintField(args, "fromID"), uintField(args, "toID"), uintField(args, "eventID"))
			if err != nil {
				return nil, err
			}
			return t, nil
		}
	case "characterAbilityHelper":
		act := stringField(args, "action")
		if act == "create" {
			ab, err := s.Services.CreateAbility(uintField(args, "characterID"), stringField(args, "name"), intField(args, "level"))
			if err != nil {
				return nil, err
			}
			return ab, nil
		}
		if act == "upgrade" {
			ab, err := s.Services.UpgradeAbility(uintField(args, "abilityID"), intField(args, "level"))
			if err != nil {
				return nil, err
			}
			return ab, nil
		}
		if act == "use" {
			u, err := s.Services.UseAbility(uintField(args, "abilityID"), uintField(args, "eventID"), stringField(args, "note"))
			if err != nil {
				return nil, err
			}
			return u, nil
		}
	case "plotThreadHelper":
		act := stringField(args, "action")
		if act == "create" {
			pt, err := s.Services.CreatePlotThread(uintField(args, "novelID"), stringField(args, "name"), stringField(args, "stage"))
			if err != nil {
				return nil, err
			}
			return pt, nil
		}
		if act == "update" {
			pt, err := s.Services.UpdatePlotStage(uintField(args, "plotID"), stringField(args, "stage"))
			if err != nil {
				return nil, err
			}
			return pt, nil
		}
	case "characterMemoryHelper":
		m, err := s.Services.CreateMemory(uintField(args, "characterID"), uintField(args, "eventID"), stringField(args, "content"), stringField(args, "trigger"))
		if err != nil {
			return nil, err
		}
		return m, nil
	case "conflictDetectionHelper":
		cs, err := s.Detector.DetectAll()
		if err != nil {
			return nil, err
		}
		return cs, nil
	case "outlineGeneratorHelper":
		act := stringField(args, "action")
		id := uintField(args, "id")
		if act == "chapter" {
			o, err := s.Generator.ChapterOutline(id)
			if err != nil {
				return nil, err
			}
			return map[string]any{"outline": o}, nil
		}
		if act == "volume" {
			o, err := s.Generator.VolumeOutline(id)
			if err != nil {
				return nil, err
			}
			return map[string]any{"outline": o}, nil
		}
		if act == "novel" {
			o, err := s.Generator.NovelOutline(id)
			if err != nil {
				return nil, err
			}
			return map[string]any{"outline": o}, nil
		}
	case "articleExportHelper":
		act := stringField(args, "action")
		id := uintField(args, "id")
		if act == "chapter" {
			res, err := s.Services.ExportChapter(id)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
		if act == "volume" {
			res, err := s.Services.ExportVolume(id)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
		if act == "novel" {
			res, err := s.Services.ExportNovel(id)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	case "styleHelper":
		act := stringField(args, "action")
		if act == "set" {
			sr, err := s.Services.SetStyleRef(uintField(args, "novelID"), stringField(args, "content"))
			if err != nil {
				return nil, err
			}
			return sr, nil
		}
		if act == "get" {
			sr, err := s.Services.GetStyleRef(uintField(args, "novelID"))
			if err != nil {
				return nil, err
			}
			return sr, nil
		}
	}
	return map[string]any{"ok": false}, nil
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.Novel{},
		&models.Volume{},
		&models.Chapter{},
		&models.World{},
		&models.Period{},
		&models.TimeSegment{},
		&models.Location{},
		&models.Character{},
		&models.CharacterRelationship{},
		&models.LocationRelationship{},
		&models.Item{},
		&models.ItemTransfer{},
		&models.Ability{},
		&models.AbilityUsage{},
		&models.PlotThread{},
		&models.Event{},
		&models.Memory{},
		&models.StyleRef{},
	)
}

func stringField(m map[string]any, k string) string {
	v, _ := m[k]
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

func intField(m map[string]any, k string) int {
	v, _ := m[k]
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	default:
		return 0
	}
}

func floatField(m map[string]any, k string) float64 {
	v, _ := m[k]
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	default:
		return 0
	}
}

func uintField(m map[string]any, k string) uint {
	v := intField(m, k)
	if v < 0 {
		return 0
	}
	return uint(v)
}

func uintSliceField(m map[string]any, k string) []uint {
	v, _ := m[k]
	var out []uint
	a, ok := v.([]any)
	if !ok {
		return out
	}
	for _, e := range a {
		switch t := e.(type) {
		case float64:
			if t >= 0 {
				out = append(out, uint(t))
			}
		case int:
			if t >= 0 {
				out = append(out, uint(t))
			}
		}
	}
	return out
}

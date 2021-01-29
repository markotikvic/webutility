package webutility

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	mu       = &sync.Mutex{}
	metadata = make(map[string]Payload)

	updateQue = make(map[string][]byte)

	metadataDB    *sql.DB
	activeProject string

	inited     bool
	metaDriver string
)

// LangMap ...
type LangMap map[string]map[string]string

// Field ...
type Field struct {
	Parameter string `json:"param"`
	Type      string `json:"type"`
	Visible   bool   `json:"visible"`
	Editable  bool   `json:"editable"`
}

// CorrelationField ...
type CorrelationField struct {
	Result   string   `json:"result"`
	Elements []string `json:"elements"`
	Type     string   `json:"type"`
}

// Translation ...
type Translation struct {
	Language     string            `json:"language"`
	FieldsLabels map[string]string `json:"fieldsLabels"`
}

// Payload ...
type Payload struct {
	Method       string             `json:"method"`
	Params       map[string]string  `json:"params"`
	Lang         []Translation      `json:"lang"`
	Fields       []Field            `json:"fields"`
	Correlations []CorrelationField `json:"correlationFields"`
	IDField      string             `json:"idField"`

	Links PaginationLinks `json:"_links"`

	// Data holds JSON payload. It can't be used for itteration.
	Data interface{} `json:"data"`
}

// NewPayload returs a payload sceleton for entity described with key.
func NewPayload(r *http.Request, key string) Payload {
	p := metadata[key]
	p.Method = r.Method + " " + r.RequestURI
	return p
}

func (p *Payload) addLang(code string, labels map[string]string) {
	t := Translation{
		Language:     code,
		FieldsLabels: labels,
	}
	p.Lang = append(p.Lang, t)
}

// SetData ...
func (p *Payload) SetData(data interface{}) {
	p.Data = data
}

// SetPaginationInfo ...
func (p *Payload) SetPaginationInfo(count, total int64, params PaginationParams) {
	p.Links.Count = count
	p.Links.Total = total
	p.Links = params.links()
}

// InitPayloadsMetadata loads all payloads' information into 'metadata' variable.
func InitPayloadsMetadata(drv string, db *sql.DB, project string) error {
	var err error
	if drv != "ora" && drv != "mysql" {
		err = errors.New("driver not supported")
		return err
	}

	metaDriver = drv
	metadataDB = db
	activeProject = project

	mu.Lock()
	defer mu.Unlock()
	err = initMetadata(project)
	if err != nil {
		return err
	}
	inited = true

	return nil
}

// EnableHotloading ...
func EnableHotloading(interval int) {
	if interval > 0 {
		go hotload(interval)
	}
}

// GetMetadataForAllEntities ...
func GetMetadataForAllEntities() map[string]Payload {
	return metadata
}

// GetMetadataForEntity ...
func GetMetadataForEntity(t string) (Payload, bool) {
	p, ok := metadata[t]
	return p, ok
}

// QueEntityModelUpdate ...
func QueEntityModelUpdate(entityType string, v interface{}) {
	updateQue[entityType], _ = json.Marshal(v)
}

// UpdateEntityModels ...
func UpdateEntityModels(command string) (total, upd, add int, err error) {
	if command != "force" && command != "missing" {
		return total, 0, 0, errors.New("webutility: unknown command: " + command)
	}

	if !inited {
		return 0, 0, 0, errors.New("webutility: metadata not initialized but update was tried")
	}

	total = len(updateQue)

	toUpdate := make([]string, 0)
	toAdd := make([]string, 0)

	for k := range updateQue {
		if _, exists := metadata[k]; exists {
			if command == "force" {
				toUpdate = append(toUpdate, k)
			}
		} else {
			toAdd = append(toAdd, k)
		}
	}

	var uStmt *sql.Stmt
	if metaDriver == "ora" {
		uStmt, err = metadataDB.Prepare("update entities set entity_model = :1 where entity_type = :2")
		if err != nil {
			return
		}
	} else if metaDriver == "mysql" {
		uStmt, err = metadataDB.Prepare("update entities set entity_model = ? where entity_type = ?")
		if err != nil {
			return
		}
	}
	for _, k := range toUpdate {
		_, err = uStmt.Exec(string(updateQue[k]), k)
		if err != nil {
			return
		}
		upd++
	}

	blankPayload, _ := json.Marshal(Payload{})
	var iStmt *sql.Stmt
	if metaDriver == "ora" {
		iStmt, err = metadataDB.Prepare("insert into entities(projekat, metadata, entity_type, entity_model) values(:1, :2, :3, :4)")
		if err != nil {
			return
		}
	} else if metaDriver == "mysql" {
		iStmt, err = metadataDB.Prepare("insert into entities(projekat, metadata, entity_type, entity_model) values(?, ?, ?, ?)")
		if err != nil {
			return
		}
	}
	for _, k := range toAdd {
		_, err = iStmt.Exec(activeProject, string(blankPayload), k, string(updateQue[k]))
		if err != nil {
			return
		}
		metadata[k] = Payload{}
		add++
	}

	return total, upd, add, nil
}

func initMetadata(project string) error {
	rows, err := metadataDB.Query(`select
		entity_type,
		metadata
		from entities
		where projekat = ` + fmt.Sprintf("'%s'", project))
	if err != nil {
		return err
	}
	defer rows.Close()

	if len(metadata) > 0 {
		metadata = nil
	}
	metadata = make(map[string]Payload)
	for rows.Next() {
		var name, load string
		rows.Scan(&name, &load)

		p := Payload{}
		err := json.Unmarshal([]byte(load), &p)
		if err != nil {
			fmt.Printf("webutility: couldn't init: '%s' metadata: %s:\n%s\n", name, err.Error(), load)
		} else {
			metadata[name] = p
		}
	}

	return nil
}

// LoadMetadataFromFile expects file in format:
//
// [ payload A identifier ]
// key1 = value1
// key2 = value2
// ...
// [ payload B identifier ]
// key1 = value1
// key2 = value2
// ...
//
// TODO(marko): Currently supports only one hardcoded language...
func LoadMetadataFromFile(path string) error {
	lines, err := ReadFileLines(path)
	if err != nil {
		return err
	}

	metadata = make(map[string]Payload)

	var name string
	for i, l := range lines {
		// skip empty lines
		if l = strings.TrimSpace(l); len(l) == 0 {
			continue
		}

		if IsWrappedWith(l, "[", "]") {
			name = strings.Trim(l, "[]")
			p := Payload{}
			p.addLang("sr", make(map[string]string))
			metadata[name] = p
			continue
		}

		if name == "" {
			return fmt.Errorf("webutility: LoadMetadataFromFile: error on line %d: [no header] [%s]", i+1, l)
		}

		parts := strings.Split(l, "=")
		if len(parts) != 2 {
			return fmt.Errorf("webutility: LoadMetadataFromFile: error on line %d: [invalid format] [%s]", i+1, l)
		}

		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if v != "-" {
			metadata[name].Lang[0].FieldsLabels[k] = v
		}
	}

	return nil
}

func hotload(n int) {
	entityScan := make(map[string]int64)
	firstCheck := true
	for {
		time.Sleep(time.Duration(n) * time.Second)
		rows, err := metadataDB.Query(`select
			ora_rowscn,
			entity_type
			from entities where projekat = ` + fmt.Sprintf("'%s'", activeProject))
		if err != nil {
			fmt.Printf("webutility: hotload failed: %v\n", err)
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}

		var toRefresh []string
		for rows.Next() {
			var scanID int64
			var entity string
			rows.Scan(&scanID, &entity)
			oldID, ok := entityScan[entity]
			if !ok || oldID != scanID {
				entityScan[entity] = scanID
				toRefresh = append(toRefresh, entity)
			}
		}
		rows.Close()

		if rows.Err() != nil {
			fmt.Printf("webutility: hotload rset error: %v\n", rows.Err())
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}

		if len(toRefresh) > 0 && !firstCheck {
			mu.Lock()
			refreshMetadata(toRefresh)
			mu.Unlock()
		}
		if firstCheck {
			firstCheck = false
		}
	}
}

func refreshMetadata(entities []string) {
	for _, e := range entities {
		fmt.Printf("refreshing %s\n", e)
		rows, err := metadataDB.Query(`select
			metadata
			from entities
			where projekat = ` + fmt.Sprintf("'%s'", activeProject) +
			` and entity_type = ` + fmt.Sprintf("'%s'", e))

		if err != nil {
			fmt.Printf("webutility: refresh: prep: %v\n", err)
			rows.Close()
			continue
		}

		for rows.Next() {
			var load string
			rows.Scan(&load)
			p := Payload{}
			err := json.Unmarshal([]byte(load), &p)
			if err != nil {
				fmt.Printf("webutility: couldn't refresh: '%s' metadata: %s\n%s\n", e, err.Error(), load)
			} else {
				metadata[e] = p
			}
		}
		rows.Close()
	}
}

/*
func ModifyMetadataForEntity(entityType string, p *Payload) error {
	md, err := json.Marshal(*p)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	_, err = metadataDB.PrepAndExe(`update entities set
		metadata = :1
		where projekat = :2
		and entity_type = :3`,
		string(md),
		activeProject,
		entityType)
	if err != nil {
		return err
	}
	return nil
}

func DeleteEntityModel(entityType string) error {
	_, err := metadataDB.PrepAndExe("delete from entities where entity_type = :1", entityType)
	if err == nil {
		mu.Lock()
		delete(metadata, entityType)
		mu.Unlock()
	}
	return err
}
*/

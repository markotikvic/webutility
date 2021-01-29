package webutility

import (
	"database/sql"
	"fmt"
)

// ListOptions ...
type ListOptions struct {
	GlobalFilter  bool   `json:"globalFilter"`
	LocalFilters  bool   `json:"localFilters"`
	RemoteFilters bool   `json:"remoteFilters"`
	Pagination    bool   `json:"pagination"`
	PageSize      uint32 `json:"pageSize"`
	Pivot         bool   `json:"pivot"`
	Detail        bool   `json:"detail"`
	Total         bool   `json:"total"`
}

// ListFilter ...
type ListFilter struct {
	Position       uint32   `json:"-"`
	ObjectType     string   `json:"-"`
	FiltersField   string   `json:"filtersField"`
	DefaultValues  string   `json:"defaultValues"`
	FiltersType    string   `json:"filtersType"`
	FiltersLabel   string   `json:"filtersLabel"`
	DropdownConfig Dropdown `json:"dropdownConfig"`
}

// Dropdown ...
type Dropdown struct {
	ObjectType   string `json:"objectType"`
	FiltersField string `json:"filtersField"`
	IDField      string `json:"idField"`
	LabelField   string `json:"labelField"`
}

// ListGraph ...
type ListGraph struct {
	ObjectType string `json:"objectType"`
	X          string `json:"xField"`
	Y          string `json:"yField"`
	GroupField string `json:"groupField"`
	Label      string `json:"label"`
}

// ListActions ...
type ListActions struct {
	Create    bool `json:"create"`
	Update    bool `json:"update"`
	Delete    bool `json:"delete"`
	Export    bool `json:"export"`
	Print     bool `json:"print"`
	Graph     bool `json:"graph"`
	LiveGraph bool `json:"liveGraph"`
	SaveFile  bool `json:"saveFile"`
	ShowFile  bool `json:"showFile"`
}

// ListNavNode ...
type ListNavNode struct {
	ObjectType        string `json:"objectType"`
	LabelField        string `json:"label"`
	Icon              string `json:"icon"`
	ParentObjectType  string `json:"parentObjectType"`
	ParentIDField     string `json:"parentIdField"`
	ParentFilterField string `json:"parentFilterField"`
}

// ListParentNode ...
type ListParentNode struct {
	ObjectType  string `json:"objectType"`
	LabelField  string `json:"labelField"`
	FilterField string `json:"filterField"`
}

// ListPivot ...
type ListPivot struct {
	ObjectType    string `json:"objectType"`
	GroupField    string `json:"groupField"`
	DistinctField string `json:"distinctField"`
	Value         string `json:"valueField"`
}

// ListDetails ...
type ListDetails struct {
	ObjectType        string `json:"objectType"`
	ParentObjectType  string `json:"parentObjectType"`
	ParentFilterField string `json:"parentFilterField"`
	SingleDetail      bool   `json:"singleDetail"`
}

// ListLiveGraph ...
type ListLiveGraph struct {
	ObjectType  string `json:"objectType"`
	ValueFields string `json:"valueFields"`
	LabelFields string `json:"labelFields"`
}

// ListConfig ...
type ListConfig struct {
	ObjectType string           `json:"objectType"`
	Title      string           `json:"title"`
	LazyLoad   bool             `json:"lazyLoad"`
	InlineEdit bool             `json:"inlineEdit"`
	Options    ListOptions      `json:"options"`
	Filters    []ListFilter     `json:"defaultFilters"`
	Graphs     []ListGraph      `json:"graphs"`
	Actions    ListActions      `json:"actions"`
	Parent     []ListParentNode `json:"parent"`
	Navigation []ListNavNode    `json:"navigation"`
	Pivots     []ListPivot      `json:"pivots"`
	Details    ListDetails      `json:"details"`
	LiveGraph  ListLiveGraph    `json:"liveGraphs"`
}

// GetListConfig returns list configuration for the provided object type for the front-end application
// or an error if it fails.
func GetListConfig(db *sql.DB, objType string) (ListConfig, error) {
	list := NewListConfig(objType)

	err := list.setParams(db, objType)
	err = list.SetNavigation(db, objType)
	err = list.SetActions(db, objType)
	err = list.SetFilters(db, objType)
	err = list.SetOptions(db, objType)
	err = list.SetParent(db, objType)
	err = list.SetPivot(db, objType)
	err = list.SetGraph(db, objType)
	err = list.SetDetails(db, objType)
	err = list.SetLiveGraph(db, objType)

	if err != nil {
		return list, err
	}

	return list, nil
}

// GetListConfigObjectIDField takes in database connection and an object type and it returns the
// ID field name for the provided object type.
func GetListConfigObjectIDField(db *sql.DB, otype string) string {
	var resp string

	rows, err := db.Query(`SELECT
		ID_FIELD
		FROM LIST_CONFIG_ID_FIELD
		WHERE OBJECT_TYPE = ` + fmt.Sprintf("'%s'", otype))
	if err != nil {
		return ""
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&resp)
	}

	if rows.Err() != nil {
		return ""
	}

	return resp
}

// NewListConfig returns default configuration for the provided object type.
func NewListConfig(objType string) ListConfig {
	list := ListConfig{
		ObjectType: objType,
		Title:      objType,
		LazyLoad:   false,
		Options: ListOptions{
			GlobalFilter:  true,
			LocalFilters:  true,
			RemoteFilters: false,
			Pagination:    true,
			PageSize:      20,
		},
		Filters: nil,
		Actions: ListActions{
			Create:    false,
			Update:    false,
			Delete:    false,
			Export:    false,
			Print:     false,
			Graph:     false,
			LiveGraph: false,
		},
		Parent:     nil,
		Navigation: nil,
	}

	return list
}

// setParams sets the default parameters of the provided configuration list for the provided object type.
func (list *ListConfig) setParams(db *sql.DB, objType string) error {
	rows, err := db.Query(`SELECT
		OBJECT_TYPE,
		TITLE,
		LAZY_LOAD,
		INLINE_EDIT
		FROM LIST_CONFIG
		WHERE OBJECT_TYPE = ` + fmt.Sprintf("'%s'", objType))
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		otype, title := "", ""
		lazyLoad, inlineEdit := 0, 0
		rows.Scan(&otype, &title, &lazyLoad, &inlineEdit)

		if otype != "" {
			list.ObjectType = otype
		}
		if title != "" {
			list.Title = title
		}
		list.LazyLoad = lazyLoad != 0
		list.InlineEdit = inlineEdit != 0
	}
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}

// SetNavigation returns set's navigation nodes for listObjType object type.
func (list *ListConfig) SetNavigation(db *sql.DB, listObjType string) error {
	list.Navigation = make([]ListNavNode, 0)
	rows, err := db.Query(`SELECT
		a.OBJECT_TYPE,
		a.PARENT_OBJECT_TYPE,
		a.LABEL,
		a.ICON,
		a.PARENT_FILTER_FIELD,
		b.PARENT_ID_FIELD
		FROM LIST_CONFIG_NAVIGATION b
		JOIN LIST_CONFIG_CHILD a ON b.PARENT_CHILD_ID = a.PARENT_CHILD_ID
		WHERE b.LIST_OBJECT_TYPE = ` + fmt.Sprintf("'%s'", listObjType) +
		` ORDER BY b.RB ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var node ListNavNode
	for rows.Next() {
		rows.Scan(&node.ObjectType, &node.ParentObjectType, &node.LabelField, &node.Icon,
			&node.ParentFilterField, &node.ParentIDField)
		list.Navigation = append(list.Navigation, node)
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

// SetActions sets list's actions based for objType object type.
func (list *ListConfig) SetActions(db *sql.DB, objType string) error {
	rows, err := db.Query(`SELECT
		ACTION_CREATE,
		ACTION_UPDATE,
		ACTION_DELETE,
		ACTION_EXPORT,
		ACTION_PRINT,
		ACTION_GRAPH,
		ACTION_LIVE_GRAPH,
		ACTION_SAVE_FILE,
		ACTION_SHOW_FILE
		FROM LIST_CONFIG
		WHERE OBJECT_TYPE = ` + fmt.Sprintf("'%s'", objType))
	if err != nil {
		return err
	}
	defer rows.Close()

	var create, update, delete, export, print, graph, liveGraph, saveFile, showFile uint32
	if rows.Next() {
		rows.Scan(&create, &update, &delete, &export, &print, &graph, &liveGraph, &saveFile, &showFile)
		list.Actions.Create = create != 0
		list.Actions.Update = update != 0
		list.Actions.Delete = delete != 0
		list.Actions.Export = export != 0
		list.Actions.Print = print != 0
		list.Actions.Graph = graph != 0
		list.Actions.LiveGraph = liveGraph != 0
		list.Actions.SaveFile = saveFile != 0
		list.Actions.ShowFile = showFile != 0
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

// SetFilters ...
func (list *ListConfig) SetFilters(db *sql.DB, objType string) error {
	list.Filters = make([]ListFilter, 0)
	filtersFields, err := list.GetFilterFieldsAndPosition(db, objType)
	if err != nil {
		return err
	}
	for field, pos := range filtersFields {
		filters, _ := list.getFiltersByFilterField(db, field)
		for _, filter := range filters {
			var f ListFilter
			f.Position = pos
			f.ObjectType = objType
			f.FiltersField = field
			f.DefaultValues = filter.DefaultValues
			f.FiltersLabel = filter.Label
			f.FiltersType = filter.Type
			if filter.Type == "dropdown" {
				err := f.SetDropdownConfig(db, field)
				if err != nil {
					return err
				}
			}
			list.Filters = append(list.Filters, f)
		}
	}

	list.sortFilters()

	return nil
}

// GetFilterFieldsAndPosition returns a map of filter fields and their respective position in the menu.
func (list *ListConfig) GetFilterFieldsAndPosition(db *sql.DB, objType string) (map[string]uint32, error) {
	filtersField := make(map[string]uint32, 0)
	rows, err := db.Query(`SELECT
		FILTERS_FIELD,
		RB
		FROM LIST_CONFIG_FILTERS
		WHERE OBJECT_TYPE = ` + fmt.Sprintf("'%s'", objType))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var field string
		var rb uint32
		rows.Scan(&field, &rb)
		filtersField[field] = rb
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return filtersField, nil
}

type _filter struct {
	DefaultValues string
	Label         string
	Type          string
}

// getFiltersByFilterField ...
func (list *ListConfig) getFiltersByFilterField(db *sql.DB, filtersField string) ([]_filter, error) {
	resp := make([]_filter, 0)
	rows, err := db.Query(`SELECT
		FILTERS_TYPE,
		FILTERS_LABEL,
		DEFAULT_VALUES
		FROM LIST_FILTERS_FIELD
		WHERE FILTERS_FIELD = ` + fmt.Sprintf("'%s'", filtersField))
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	var f _filter
	for rows.Next() {
		rows.Scan(&f.Type, &f.Label, &f.DefaultValues)
		resp = append(resp, f)
	}
	if rows.Err() != nil {
		return resp, rows.Err()
	}
	return resp, nil
}

// SetDropdownConfig ...
func (f *ListFilter) SetDropdownConfig(db *sql.DB, filtersField string) error {
	var resp Dropdown
	rows, err := db.Query(`SELECT
		FILTERS_FIELD,
		OBJECT_TYPE,
		ID_FIELD,
		LABEL_FIELD
		FROM LIST_DROPDOWN_FILTER
		WHERE FILTERS_FIELD = ` + fmt.Sprintf("'%s'", filtersField))
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&resp.FiltersField, &resp.ObjectType, &resp.IDField, &resp.LabelField)
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	f.DropdownConfig = resp

	return nil
}

// sortFilters bubble sorts provided filters slice by position field.
func (list *ListConfig) sortFilters() {
	done := false
	var temp ListFilter
	for !done {
		done = true
		for i := 0; i < len(list.Filters)-1; i++ {
			if list.Filters[i].Position > list.Filters[i+1].Position {
				done = false
				temp = list.Filters[i]
				list.Filters[i] = list.Filters[i+1]
				list.Filters[i+1] = temp
			}
		}
	}
}

// SetGraph ...
func (list *ListConfig) SetGraph(db *sql.DB, objType string) error {
	list.Graphs = make([]ListGraph, 0)
	rows, err := db.Query(`SELECT
		OBJECT_TYPE,
		X_FIELD,
		Y_FIELD,
		GROUP_FIELD,
		LABEL
		FROM LIST_GRAPHS
		WHERE OBJECT_TYPE = ` + fmt.Sprintf("'%s'", objType))
	if err != nil {
		return err
	}
	defer rows.Close()

	var lg ListGraph
	for rows.Next() {
		rows.Scan(&lg.ObjectType, &lg.X, &lg.Y, &lg.GroupField, &lg.Label)
		list.Graphs = append(list.Graphs, lg)
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

// SetOptions ...
func (list *ListConfig) SetOptions(db *sql.DB, objType string) error {
	rows, err := db.Query(`SELECT
		GLOBAL_FILTER,
		LOCAL_FILTER,
		REMOTE_FILTER,
		PAGINATION,
		PAGE_SIZE,
		PIVOT,
		DETAIL,
		TOTAL
		FROM LIST_CONFIG
		WHERE OBJECT_TYPE = ` + fmt.Sprintf("'%s'", objType))
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		var gfilter, lfilters, rfilters, pagination, pageSize, pivot, detail, total uint32
		rows.Scan(&gfilter, &lfilters, &rfilters, &pagination, &pageSize, &pivot, &detail, &total)
		list.Options.GlobalFilter = gfilter != 0
		list.Options.LocalFilters = lfilters != 0
		list.Options.RemoteFilters = rfilters != 0
		list.Options.Pagination = pagination != 0
		list.Options.PageSize = pageSize
		list.Options.Pivot = pivot != 0
		list.Options.Detail = detail != 0
		list.Options.Total = total != 0
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

// SetParent ...
func (list *ListConfig) SetParent(db *sql.DB, objType string) error {
	list.Parent = make([]ListParentNode, 0)
	rows, err := db.Query(`SELECT
		PARENT_OBJECT_TYPE,
		PARENT_LABEL_FIELD,
		PARENT_FILTER_FIELD
		FROM LIST_CONFIG_CHILD
		WHERE OBJECT_TYPE = ` + fmt.Sprintf("'%s'", objType))
	if err != nil {
		return err
	}
	defer rows.Close()

	var pnode ListParentNode
	for rows.Next() {
		rows.Scan(&pnode.ObjectType, &pnode.LabelField, &pnode.FilterField)
		list.Parent = append(list.Parent, pnode)
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

// SetPivot ...
func (list *ListConfig) SetPivot(db *sql.DB, objType string) error {
	list.Pivots = make([]ListPivot, 0)
	rows, err := db.Query(`SELECT
		OBJECT_TYPE,
		GROUP_FIELD,
		DISTINCT_FIELD,
		VALUE_FIELD
		FROM LIST_PIVOTS
		WHERE OBJECT_TYPE = ` + fmt.Sprintf("'%s'", objType))
	if err != nil {
		return err
	}
	defer rows.Close()

	var p ListPivot
	for rows.Next() {
		rows.Scan(&p.ObjectType, &p.GroupField, &p.DistinctField, &p.Value)
		list.Pivots = append(list.Pivots, p)
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

// SetDetails ...
func (list *ListConfig) SetDetails(db *sql.DB, objType string) error {
	var resp ListDetails
	rows, err := db.Query(`SELECT
		OBJECT_TYPE,
		PARENT_OBJECT_TYPE,
		PARENT_FILTER_FIELD,
		SINGLE_DETAIL
		FROM LIST_CONFIG_DETAIL
		WHERE PARENT_OBJECT_TYPE = ` + fmt.Sprintf("'%s'", objType))
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		var singleDetail uint32
		rows.Scan(&resp.ObjectType, &resp.ParentObjectType, &resp.ParentFilterField, &singleDetail)
		resp.SingleDetail = singleDetail != 0
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	list.Details = resp

	return nil
}

// SetLiveGraph ...
func (list *ListConfig) SetLiveGraph(db *sql.DB, objType string) error {
	var resp ListLiveGraph
	rows, err := db.Query(`SELECT
		OBJECT_TYPE,
		VALUE_FIELDS,
		LABEL_FIELDS
		FROM LIST_LIVE_GRAPH
		WHERE OBJECT_TYPE = ` + fmt.Sprintf("'%s'", objType))
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&resp.ObjectType, &resp.ValueFields, &resp.LabelFields)
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	list.LiveGraph = resp

	return nil
}

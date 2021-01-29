package webutility

import "database/sql"

// SelectConfig ...
type SelectConfig struct {
	ListObjType string `json:"listObjectType"`
	ObjType     string `json:"objectType"`
	Type        string `json:"type"`
	IDField     string `json:"idField"`
	LabelField  string `json:"labelField"`
	ValueField  string `json:"valueField"`
}

// GetSelectConfig returns select configuration slice for the given object type.
func GetSelectConfig(db *sql.DB, otype string) ([]SelectConfig, error) {
	resp := make([]SelectConfig, 0)
	rows, err := db.Query(`SELECT
		a.LIST_OBJECT_TYPE,
		a.OBJECT_TYPE,
		a.ID_FIELD,
		a.LABEL_FIELD,
		a.TYPE,
		b.FIELD
		FROM LIST_SELECT_CONFIG a, LIST_VALUE_FIELD b
		WHERE a.LIST_OBJECT_TYPE` + otype + `
		AND b.LIST_TYPE = a.LIST_OBJECT_TYPE
		AND b.OBJECT_TYPE = a.OBJECT_TYPE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sc SelectConfig
	for rows.Next() {
		rows.Scan(&sc.ListObjType, &sc.ObjType, &sc.IDField, &sc.LabelField, &sc.Type, &sc.ValueField)
		resp = append(resp, sc)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return resp, nil
}

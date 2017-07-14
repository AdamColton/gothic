package sqlmodel

import (
	"text/template"
)

var templates = template.Must(template.New("templates").Parse(`
{{define "insert"}}	res, err := {{.Conn}}.Exec("INSERT INTO {{.QName}} ({{.Fields}}) VALUES ({{.QM}})", {{.Args}})
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	{{.Receiver}}.{{.Primary}} = {{.PrimaryType}}(id)
	return nil{{end}}
{{define "update"}}	_, err := {{.Conn}}.Exec("UPDATE {{.QName}} SET ({{.Set}}) WHERE {{.QPrimary}}=?", {{.Args}}, {{.PrimaryArg}})
	return err{{end}}
{{define "createTable"}}func init() {
	gsql.AddMigration("{{.Migration}}",
	{{.BackTick}}CREATE TABLE IF NOT EXISTS "{{.Name}}" (
			{{.DefineTable}}
		);{{.BackTick}},
	"DROP TABLE {{.QName}};"){{end}}
}
{{define "scan"}}	{{.Receiver}} := &{{.Name}}{}{{range .ConvertFields}}
	var {{.Name}} {{.Type}}{{end}}
	err := rows.Scan({{.ScanFields}})
	if err != nil {
		return nil, err
	}{{range .ConvertFields}}
	{{.R}}.{{.Name}} = {{.FromDB}}({{.Name}}){{end}}
	return {{.Receiver}}, nil{{end}}
{{define "select"}}rows, err := {{.Conn}}.Query("SELECT {{.Fields}} FROM {{.QName}} "+where, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var {{.Receiver}}s []*{{.Name}}
	for rows.Next() {
		{{.Receiver}},err := {{.Scanner}}(rows)
		if err != nil {
			return nil, err
		}	
		{{.Receiver}}s = append({{.Receiver}}s, {{.Receiver}})
	}
	return {{.Receiver}}s, nil{{end}}
{{define "upsert"}}	if {{.Receiver}}.{{.Primary}} != {{.PrimaryZeroVal}} {
		_, err := {{.Conn}}.Exec("UPDATE {{.QName}} SET ({{.Set}}) WHERE {{.QPrimary}}=?", {{.Args}}, {{.PrimaryArg}})
		return err
	}
	res, err := {{.Conn}}.Exec("INSERT INTO {{.QName}} ({{.Fields}}) VALUES ({{.QM}})", {{.Args}})
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	{{.Receiver}}.{{.Primary}} = {{.PrimaryType}}(id)
	return nil{{end}}
`))

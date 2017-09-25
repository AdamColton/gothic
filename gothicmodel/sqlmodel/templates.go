package sqlmodel

import (
	"text/template"
)

// Templates for generating SQL methods and functions.
var Templates = template.Must(template.New("templates").Parse(`
{{define "insert"}}	res, err := {{.Conn}}.Exec("INSERT INTO {{.TableNameQ}} ({{.FieldsQ}}) VALUES ({{.QM}})", {{.Args}})
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	{{.Receiver}}.{{.Primary}} = {{.PrimaryType}}(id)
	return nil{{end}}
{{define "update"}}	_, err := {{.Conn}}.Exec("UPDATE {{.TableNameQ}} SET ({{.Set}}) WHERE {{.PrimaryQ}}=?", {{.Args}}, {{.PrimaryArg}})
	return err{{end}}
{{define "createTable"}}func init() {
	gsql.AddMigration("{{.Migration}}",
	{{.BackTick}}CREATE TABLE IF NOT EXISTS "{{.Name}}" (
			{{.DefineTable}}
		);{{.BackTick}},
	"DROP TABLE {{.TableNameQ}};")
}{{end}}
{{define "scan"}}	{{.Receiver}} := &{{.Name}}{}{{range .FieldConverters}}
	var {{.Name}} {{.GoType}}{{end}}
	err := rows.Scan({{.ScanFields}})
	if err != nil {
		return nil, err
	}{{range .FieldConverters}}
	{{.Receiver}}.{{.Name}} = {{.FromDB}}{{end}}
	return {{.Receiver}}, nil{{end}}
{{define "select"}}	rows, err := {{.Conn}}.Query("SELECT {{.FieldsQ}} FROM {{.TableNameQ}} "+where, args...)
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
	{{template "update" .}}
	}
{{template "insert" .}}{{end}}
`))

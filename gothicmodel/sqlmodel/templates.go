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
{{define "createTable"}}
	gsql.AddMigration("{{.Migration}}",
		{{.BackTick}}CREATE TABLE IF NOT EXISTS "{{.Name}}" (
				{{.DefineTable}}
			);{{.BackTick}},
		"DROP TABLE {{.QName}};")
{{end}}
`))

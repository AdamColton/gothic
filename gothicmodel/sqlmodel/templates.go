package sqlmodel

import (
	"text/template"
)

// Templates for generating SQL methods and functions.
var Templates = template.Must(template.New("templates").Parse(`
{{define "insert" -}}
	{{"\t"}}res, err := {{.Conn}}.Exec("INSERT INTO {{.TableNameQ}} ({{.FieldsQ}}) VALUES ({{.QM}})", {{.Args}})
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	{{.Receiver}}.{{.Primary}} = {{.PrimaryGoType}}(id)
	return nil
{{- end}}
{{define "update" -}}
	{{"\t"}}_, err := {{.Conn}}.Exec("UPDATE {{.TableNameQ}} SET {{.Set}} WHERE {{.PrimaryQ}}=?", {{.Args}}, {{.PrimaryArg}})
	return err
{{- end}}
{{define "createTable"}}
	gsql.AddMigration("{{.Migration}}",
	{{.BackTick}}CREATE TABLE IF NOT EXISTS {{.TableNameLQ}} (
			{{.DefineTable}}
		);{{.BackTick}},
	"DROP TABLE {{.TableNameQ}};")
{{end}}
{{define "scan" -}}
	{{"\t"}}{{.Receiver}} := &{{.Name}}{}
	{{- range .FieldConverters}}
	var {{.Name}} {{.GoType}}
	{{- end}}
	err := row.Scan({{.ScanFields}})
	if err != nil {
		return nil, err
	}
	{{- range .FieldConverters}}
	{{.Receiver}}.{{.Name}} = {{.FromDB}}
	{{- end}}
	return {{.Receiver}}, nil
{{- end}}
{{define "select" -}}
	{{"\t"}}rows, err := {{.Conn}}.Query("SELECT {{.FieldsQ}} FROM {{.TableNameQ}} "+where, args...)
	if err != nil || rows==nil{
		return nil, err
	}
	defer rows.Close()
	var {{.Receiver}}s []*{{.Name}}
	for rows.Next() {
		{{.Receiver}},err := {{.Scanner}}(rows)
		if err != nil {
			return nil, err
		}	
		{{.Receiver}}s = append({{.Receiver}}s, {{.Receiver}})
	}
	return {{.Receiver}}s, nil{{end}}
{{define "upsert" -}}
	if {{.Receiver}}.{{.Primary}} != {{.PrimaryZeroVal}} {
	{{template "update" .}}
	}
{{template "insert" .}}
{{- end}}
{{define "selectsingle" -}}
	{{"\t"}}{{.Receiver}}, err := {{.Select}}(where+" LIMIT 1", args...)
	if err != nil || len({{.Receiver}}) == 0 {
		return nil, err
	}
	return {{.Receiver}}[0], nil
{{- end}}
{{define "delete" -}}
	{{"\t"}}res, err := {{.Conn}}.Exec("DELETE FROM {{.TableNameQ}} WHERE {{.AndConditions}}", {{.Args}})
	if err != nil || res == nil {
		return 0, err
	}
	return res.RowsAffected()
{{- end}}
{{define "whereEqual" -}}
{{"\t"}}return {{.Select}}("WHERE {{.AndConditions}}", {{.Custom}})
{{- end}}
{{define "whereEqualSingle" -}}
{{"\t"}}{{.Receiver}}, err := {{.Select}}("WHERE {{.AndConditions}} LIMIT 1", {{.Custom}})
	if err != nil || len({{.Receiver}}) == 0 {
		return nil, err
	}
	return {{.Receiver}}[0], nil
{{- end}}
{{define "mustWhereEqual" -}}
{{"\t"}} {{.Receiver}}, err := {{.Select}}("WHERE {{.AndConditions}}", {{.Custom}})
	if err != nil {
		panic(err)
	}
	return {{.Receiver}}
{{- end}}
{{define "mustWhereEqualSingle" -}}
{{"\t"}}{{.Receiver}}, err := {{.Select}}("WHERE {{.AndConditions}} LIMIT 1", {{.Custom}})
	if err != nil {
		panic(err)
	}
	if len({{.Receiver}}) == 0 {
		return nil
	}
	return {{.Receiver}}[0]
{{- end}}
`))

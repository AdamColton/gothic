package htmlmodel

import (
	"fmt"
	"github.com/adamcolton/gothic/gothichtml"
	"github.com/adamcolton/gothic/gothicmodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() *gothicmodel.Model {
	InputTypes["primary"] = "hidden"
	m := gothicmodel.New("test").
		AddPrimary("ID", "primary").
		AddField("Name", "string").
		AddField("Age", "int").
		AddField("LastLogin", "datetime")

	return m
}

func TestGenerateLapizForm(t *testing.T) {
	generateRowWrapper := func(field, labelStr string, children ...gothichtml.Node) *gothichtml.Tag {
		div := gothichtml.NewTag("div", "class", "row")
		label := div.Tag("label")
		label.Text(labelStr)
		if field != "" {
			label.AddAttributes("for", field)
		}
		div.AddChildren(children...)
		return div
	}

	var generateRow FieldGenerator = func(field, kind string, model *gothicmodel.Model) gothichtml.Node {
		vt := gothichtml.NewVoidTag("input", "type", kind, "id", field, "name", field, "value", "$"+field)
		if kind == "hidden" {
			return vt
		}
		return generateRowWrapper(field, field, vt)
	}

	var generateForm = func(m *gothicmodel.Model) gothichtml.Node {
		form := gothichtml.NewTag("form", "l-view", "editPerson", "submit", "form.saveTask")
		InputTypes.GenerateFields(m, generateRow, form)
		save := generateRowWrapper("", "&nbsp;")
		save.VoidTag("input", "type", "submit", "value", "Save")
		form.AddChildren(save)
		return form
	}

	expected := `<form l-view="editPerson" submit="form.saveTask">
	<input id="ID" name="ID" type="hidden" value="$ID" />
	<div class="row">
		<label for="Name">Name</label>
		<input id="Name" name="Name" type="text" value="$Name" />
	</div>
	<div class="row">
		<label for="Age">Age</label>
		<input id="Age" name="Age" type="number" value="$Age" />
	</div>
	<div class="row">
		<label>&nbsp;</label>
		<input type="submit" value="Save" />
	</div>
</form>`
	m := setup()
	form := generateForm(m)
	assert.Equal(t, expected, gothichtml.String(form))
}

func TestGenerateTable(t *testing.T) {
	var generateHeader FieldGenerator = func(field, kind string, model *gothicmodel.Model) gothichtml.Node {
		th := gothichtml.NewTag("th")
		th.Text(field)
		return th
	}

	var generateRow FieldGenerator = func(field, kind string, model *gothicmodel.Model) gothichtml.Node {
		td := gothichtml.NewTag("td")
		td.Text("{{." + field + "}}")
		return td
	}

	var generateTable = func(m *gothicmodel.Model, fields ...string) gothichtml.Node {
		table := gothichtml.NewTag("table")
		header := table.Tag("tr")
		InputTypes.GenerateFields(m, generateHeader, header, fields...)
		header.Tag("th").Text("Actions")
		table.Text("{{range .}}")
		row := table.Tag("tr")
		InputTypes.GenerateFields(m, generateRow, row, fields...)
		actions := row.Tag("td")
		actions.Tag("a", "href", fmt.Sprintf("/%s/edit/{{.ID}}", m.Name())).Text("Edit")
		actions.Tag("a", "href", fmt.Sprintf("/%s/delete/{{.ID}}", m.Name())).Text("Delete")
		table.Text("{{end}}")
		return table
	}

	expected := `<table>
	<tr>
		<th>Name</th>
		<th>Age</th>
		<th>Actions</th>
	</tr>
	{{range .}}
	<tr>
		<td>{{.Name}}</td>
		<td>{{.Age}}</td>
		<td>
			<a href="/test/edit/{{.ID}}">Edit</a>
			<a href="/test/delete/{{.ID}}">Delete</a>
		</td>
	</tr>
	{{end}}
</table>`

	m := setup()
	table := generateTable(m, "Name", "Age")
	assert.Equal(t, expected, gothichtml.String(table))
}

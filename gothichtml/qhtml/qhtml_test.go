package qhtml

import (
	"github.com/adamcolton/gothic/gothichtml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDocument(t *testing.T) {
	example := `
    <title="This is the title"/>
    <css="/css/foo.css /css/bar.css"/>
    <scripts="/js/foo.js /js/bar.js"/>
    <div.row>
      <label=name foo="this is \"a test">Name</>
      <text#name=$name />
    </>
    <hr />
    <div.row>
      <label=age>Age</>
      <text#age=$age />
    </>
    <div.row>
      <label>&nbsp;</>
      <button>Save</button>
    </>
`
	html, err := ParseDocument(example)
	assert.NoError(t, err)

	expected := `<!DOCTYPE html>
<header>
  <title>This is the title</title>
  <link href="/css/foo.css" rel="stylesheet" />
  <link href="/css/bar.css" rel="stylesheet" />
</header>
<body>
  <div class="row">
    <label foo="this is \"a test" for="name">Name</label>
    <text id="name" value="$name" />
  </div>
  <hr />
  <div class="row">
    <label for="age">Age</label>
    <text id="age" value="$age" />
  </div>
  <div class="row">
    <label>&nbsp;</label>
    <button>Save</button>
  </div>
  <script src="/js/foo.js"></script>
  <script src="/js/bar.js"></script>
</body>`
	gothichtml.Padding = "  "
	assert.Equal(t, expected, gothichtml.String(html))
}

func TestParseFragment(t *testing.T) {
	example := `
    <div.row>
      <label=name foo="this is \"a test">Name</>
      <text#name=$name />
    </>
    <hr />
    <div.row>
      <label=age>Age</>
      <text#age=$age />
    </>
    <div.row>
      <label>&nbsp;</>
      <button>Save</button>
    </>
`
	html, err := ParseFragment(example)
	assert.NoError(t, err)

	expected := `<div class="row">
  <label foo="this is \"a test" for="name">Name</label>
  <text id="name" value="$name" />
</div>
<hr />
<div class="row">
  <label for="age">Age</label>
  <text id="age" value="$age" />
</div>
<div class="row">
  <label>&nbsp;</label>
  <button>Save</button>
</div>`
	gothichtml.Padding = "  "
	assert.Equal(t, expected, gothichtml.String(html))
}

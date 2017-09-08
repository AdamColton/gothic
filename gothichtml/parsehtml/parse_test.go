package parsehtml

import (
	"github.com/adamcolton/gothic/gothichtml"
	"github.com/stretchr/testify/assert"
	"testing"
)

var doc = `<!DOCTYPE html><html><header><title>Test Doc</title><link href="/css/foo.css" rel=stylesheet type="text/css"/><link href="/css/bar.css" rel=stylesheet type="text/css"/></header><body id=body><p>I am a paragraph</p><ul><li>Item 1</li><li>Item 2</li><li>Item 3  </li></ul><script src="/js/foo.js"></script><script src="/js/bar.js"></script></body></html>`

func TestPieces(t *testing.T) {
	lxms := lxr.Lex(doc)
	assert.NotNil(t, lxms)
	pn := prsr.Parse(lxms)
	pn = rdcr.Reduce(pn)
	assert.NotNil(t, pn)
}

func TestParse(t *testing.T) {
	expected := `<!DOCTYPE html>
<html>
  <header>
    <title>Test Doc</title>
    <link href="/css/foo.css" rel="stylesheet" type="text/css" />
    <link href="/css/bar.css" rel="stylesheet" type="text/css" />
  </header>
  <body id="body">
    <p>I am a paragraph</p>
    <ul>
      <li>Item 1</li>
      <li>Item 2</li>
      <li>Item 3</li>
    </ul>
    <script src="/js/foo.js"></script>
    <script src="/js/bar.js"></script>
  </body>
</html>`
	nodes, err := Parse(doc)
	if assert.NoError(t, err) {
		str := gothichtml.String(nodes)
		assert.Equal(t, expected, str)
	}
}

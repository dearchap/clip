package clip

import (
	"bytes"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

func TestHelpContextFullName(t *testing.T) {
	var hctx *helpContext

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		hctx = newHelpContext(ctx)
		return nil
	}

	grandchild := NewCommand("grandchild", WithAction(action))
	child := NewCommand("child", WithCommand(grandchild))
	root := NewCommand("root", WithCommand(child))

	args := []string{root.Name(), child.Name(), grandchild.Name()}
	assert.NilError(t, root.Execute(args))
	assert.Assert(t, wasCalled)
	assert.Check(t, hctx.FullName() == "root child grandchild")
}

func TestHelpCommands(t *testing.T) {
	buf := new(bytes.Buffer)
	root := NewCommand(
		"root",
		WithWriter(buf),
		WithCommand(NewCommand("child-one", WithSummary("1"))),
		WithCommand(NewCommand("child-two", WithSummary("2"))),
		WithCommand(NewCommand("child-three", WithSummary("3"))),
	)

	args := []string{root.Name()}
	assert.NilError(t, root.Execute(args))

	output := buf.String()
	assert.Check(t, cmp.Contains(output, "child-one    1"))
	assert.Check(t, cmp.Contains(output, "child-two    2"))
	assert.Check(t, cmp.Contains(output, "child-three  3"))
}

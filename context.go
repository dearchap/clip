package clip

// Context is a command context with runtime metadata.
type Context struct {
	*Command

	args   []string
	parent *Context
}

// Args returns the list of arguments.
func (ctx *Context) Args() []string { return ctx.args }

// Parent is the context's parent context.
func (ctx *Context) Parent() *Context { return ctx.parent }

// Root is the context's root context.
func (ctx *Context) Root() *Context {
	cur := ctx
	for cur.parent != nil {
		cur = cur.parent
	}
	return cur
}

// run runs the command with a given context.
func (ctx *Context) run() error {
	// No sub commands
	if len(ctx.commands) == 0 {
		return ctx.action(ctx)
	}

	// Sub commands, but nothing passed
	if len(ctx.args) == 0 {
		return newUsageError(ctx, "no sub-command passed")
	}

	// Sub commands, something passed
	subCmdName := ctx.args[0]
	if subCmd, ok := ctx.subCommandMap[subCmdName]; ok {
		subCtx := Context{
			Command: subCmd,
			args:    ctx.args[1:],
			parent:  ctx,
		}
		return subCtx.run()
	}

	return newUsageErrorf(ctx, "undefined sub-command %q", subCmdName)
}

package pulumi

import (
	"regexp"

	p "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Context struct {
	ctx *p.Context
}

func (c *Context) Export(name string, val string) {
	c.ctx.Export(name, TokenToInput(val))
}

func (c *Context) RegisterResource(
	t, name string, props p.Input, resource p.Resource, opts ...p.ResourceOption) error {
	return c.ctx.RegisterResource(t, name, props, resource, opts...)
}

func Run(body func(ctx *Context) error, opts ...p.RunOption) {
	p.Run(func(ctx *p.Context) error {
		return body(&Context{ctx})
	}, opts...)
}

func TokenToInput(s string) p.StringInput {
	// Find `UNKNOWN` and replace it
	r := regexp.MustCompile(`UNKNOWN`)
	matches := r.FindAllString(s, -1)
	var tokenOutputs []interface{}
	for range matches {
		unknown := p.UnsafeUnknownOutput(nil)
		tokenOutputs = append(tokenOutputs, unknown)
	}
	fmtString := r.ReplaceAllString(s, "%s")
	return p.Sprintf(fmtString, tokenOutputs...)
}

func PtrTokenToInput(s *string) p.StringPtrInput {
	if s == nil {
		return nil
	}
	return TokenToInput(*s)
}

func Ptr[T any](x T) *T {
	return &x
}

func Printf(format string, args ...string) {
	var innerArgs []interface{}
	for _, arg := range args {
		innerArgs = append(innerArgs, TokenToInput(arg))
	}
	p.Printf(format, innerArgs)
}

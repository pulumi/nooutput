package pulumi

import (
	"fmt"
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

var tokens map[string]p.StringOutput = map[string]p.StringOutput{}
var nextToken = 0

func StringPtr(v string) *string {
	return &v
}

func OutputToToken(v p.StringOutput) string {
	token := fmt.Sprintf("#{$Token%d}#", nextToken)
	nextToken++
	tokens[token] = v
	return token
}

func TokenToInput(s string) p.StringInput {
	r := regexp.MustCompile(`#\{\$Token\d+\}#`)
	matches := r.FindAllStringIndex(s, -1)
	var tokenOutputs []interface{}
	for _, loc := range matches {
		token := s[loc[0]:loc[1]]
		tokenOutput, ok := tokens[token]
		if !ok {
			panic(fmt.Sprintf("invalid token: %s", token))
		}
		tokenOutputs = append(tokenOutputs, tokenOutput)
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

func Apply(token string, callback func(s string) string) string {
	return OutputToToken(TokenToInput(token).ToStringOutput().ApplyT(callback).(p.StringOutput))
}

type StringMap p.StringMapOutput

func (o StringMap) Lookup(s string) string {
	return OutputToToken(p.StringMapOutput(o).MapIndex(p.String(s)))
}

package pulumi

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"

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

var tokens map[int]p.StringOutput = map[int]p.StringOutput{}
var nextToken = 0

func StringPtr(v string) *string {
	return &v
}

func awaitValOrUnknown(o p.StringOutput) string {
	ch := make(chan string)
	go func() {
		v, known, _, _, _ := o.Await(context.Background())
		if !known {
			ch <- "UNKNOWN"
		} else {
			ch <- v.(string)
		}
	}()
	return <-ch
}

func OutputToToken(v p.StringOutput) string {
	val := awaitValOrUnknown(v)
	token := fmt.Sprintf("#{$Token%d:%s}#", nextToken, val)
	tokens[nextToken] = v
	nextToken++
	return token
}

func TokenToInput(s string) p.StringInput {
	r := regexp.MustCompile(`#\{\$Token(\d+):.*\}#`)
	matches := r.FindAllStringSubmatchIndex(s, -1)
	var tokenOutputs []interface{}
	for _, loc := range matches {
		token := s[loc[2]:loc[3]]
		i, err := strconv.Atoi(token)
		contract.AssertNoError(err)
		tokenOutput, ok := tokens[i]
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
	input := TokenToInput(token).ToStringOutput()
	output := input.ApplyT(callback).(p.StringOutput)
	return OutputToToken(output)
}

type StringMap p.StringMapOutput

func (o StringMap) Lookup(s string) string {
	return OutputToToken(p.StringMapOutput(o).MapIndex(p.String(s)))
}

package htmx

import (
	"context"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// Response contain data that the server can communicate back to HTMX
type Response struct {
	Push               string
	Redirect           string
	Refresh            bool
	Trigger            string
	TriggerAfterSwap   string
	TriggerAfterSettle string
	NoContent          bool
}
type Context struct {
	*gin.Context
}

const (
	HeaderRequest            = "HX-Request"
	HeaderBoosted            = "HX-Boosted"
	HeaderTrigger            = "HX-Trigger"
	HeaderTriggerName        = "HX-Trigger-Name"
	HeaderTriggerAfterSwap   = "HX-Trigger-After-Swap"
	HeaderTriggerAfterSettle = "HX-Trigger-After-Settle"
	HeaderTarget             = "HX-Target"
	HeaderPrompt             = "HX-Prompt"
	HeaderPush               = "HX-Push"
	HeaderRedirect           = "HX-Redirect"
	HeaderRefresh            = "HX-Refresh"
)

func NewWrap(ginctx *gin.Context) *Context {
	return &Context{ginctx}
}

func (c *Context) Templ(status int, component templ.Component) {
	newContext := context.WithValue(c.Request.Context(), "ahtmx", c)
	r := NewRenderer(newContext, status, component)
	c.Render(status, r)
}
func Ctx(ctx context.Context) *Context {
	return ctx.Value("ahtmx").(*Context)
}
func (c *Context) RequestIsBoosted() bool {
	return c.Request.Header.Get(HeaderBoosted) == "true"
}
func (c *Context) RequestIsEnabled() bool {
	return c.Request.Header.Get(HeaderRequest) == "true"
}
func (c *Context) RequestTrigger() string {
	return c.Request.Header.Get(HeaderTrigger)
}
func (c *Context) RequestTriggerName() string {
	return c.Request.Header.Get(HeaderTriggerName)
}
func (c *Context) RequestTarget() string {
	return c.Request.Header.Get(HeaderTarget)
}
func (c *Context) RequestPrompt() string {
	return c.Request.Header.Get(HeaderPrompt)
}
func (c *Context) RequestPath() string {
	return c.Request.URL.Path
}
func (c *Context) ResponsePush(push string) {
	c.Header(HeaderPush, push)
}
func (c *Context) ResponseRedirect(path string) {
	c.Header(HeaderRedirect, path)
}
func (c *Context) ResponseRefresh(refresh bool) {
	var str string
	if refresh {
		str = "true"
	} else {
		str = "false"
	}
	c.Header(HeaderRefresh, str)
}
func (c *Context) ResponseTrigger(trigger string) {
	c.Header(HeaderTrigger, trigger)
}
func (c *Context) ResponseTriggerAfterSwap(trigger string) {
	c.Header(HeaderTriggerAfterSwap, trigger)
}
func (c *Context) ResponseTriggerAfterSettle(trigger string) {
	c.Header(HeaderTriggerAfterSettle, trigger)
}
func (c *Context) ResponseError(status int, err error) {
	c.Error(err)
	c.Status(status)
}

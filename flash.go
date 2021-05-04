package flash

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"regexp"
)

type Flash struct {
	CookiePrefix string
	data         fiber.Map
}

var DefaultFlash = &Flash{
	CookiePrefix: "Fiber-App",
	data:         fiber.Map{},
}

var cookieKeyValueParser = regexp.MustCompile("\x00([^:]*):([^\x00]*)\x00")

func New(CookiePrefix string) {
	DefaultFlash = &Flash{
		CookiePrefix: CookiePrefix,
		data:         fiber.Map{},
	}
}

func (f *Flash) Get(c *fiber.Ctx) fiber.Map {
	t := fiber.Map{}
	f.data = nil
	cookieValue := c.Cookies(f.CookiePrefix + "-Flash")
	if cookieValue != "" {
		parseKeyValueCookie(cookieValue, func(key string, val interface{}) {
			t[key] = val
		})
		f.data = t
	}
	c.Set("Set-Cookie", f.CookiePrefix+"-Flash=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; HttpOnly")
	return f.data
}

func (f *Flash) Redirect(c *fiber.Ctx, location string, data interface{}, status ...int) error {

	f.data = data.(fiber.Map)
	if len(status) > 0 {
		return c.Redirect(location, status[0])
	} else {
		return c.Redirect(location, fiber.StatusFound)
	}
}

func (f *Flash) WithError(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.data = data
	f.error(c)
	return c
}

func (f *Flash) WithSuccess(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.data = data
	f.success(c)
	return c
}

func (f *Flash) error(c *fiber.Ctx) {
	f.data["error"] = true
	f.setCookie(c)
}

func (f *Flash) success(c *fiber.Ctx) {
	f.data["success"] = true
	f.setCookie(c)
}

func (f *Flash) setCookie(c *fiber.Ctx) {
	var flashValue string
	for key, value := range f.data {
		flashValue += "\x00" + key + ":" + fmt.Sprintf("%v", value) + "\x00"
	}
	c.Cookie(&fiber.Cookie{
		Name:  f.CookiePrefix + "-Flash",
		Value: url.QueryEscape(flashValue),
	})
}

func Get(c *fiber.Ctx) fiber.Map {
	return DefaultFlash.Get(c)
}

func Redirect(c *fiber.Ctx, location string, data interface{}, status ...int) error {
	return DefaultFlash.Redirect(c, location, data)
}

func WithError(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithError(c, data)
}

func WithSuccess(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithSuccess(c, data)
}

func WithData(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	DefaultFlash.data = data
	return c
}

// parseKeyValueCookie takes the raw (escaped) cookie value and parses out key values.
func parseKeyValueCookie(val string, cb func(key string, val interface{})) {
	val, _ = url.QueryUnescape(val)
	if matches := cookieKeyValueParser.FindAllStringSubmatch(val, -1); matches != nil {
		for _, match := range matches {
			cb(match[1], match[2])
		}
	}
}

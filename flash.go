package flash

import (
	"fmt"
	"github.com/gofiber/fiber"
	"net/url"
	"regexp"
)

type Flash struct {
	CookiePrefix string
	Data         fiber.Map
}

var cookieKeyValueParser = regexp.MustCompile("\x00([^:]*):([^\x00]*)\x00")

func (f *Flash) Error(c *fiber.Ctx) {
	var flashValue string
	f.Data["error"] = true
	for key, value := range f.Data {
		flashValue += "\x00" + key + ":" + fmt.Sprintf("%v", value) + "\x00"
	}
	c.Cookie(&fiber.Cookie{
		Name:  f.CookiePrefix + "-Flash",
		Value: url.QueryEscape(flashValue),
	})
}

func (f *Flash) Success(c *fiber.Ctx) {
	var flashValue string
	f.Data["success"] = true
	for key, value := range f.Data {
		flashValue += "\x00" + key + ":" + fmt.Sprintf("%v", value) + "\x00"
	}
	c.Cookie(&fiber.Cookie{
		Name:  f.CookiePrefix + "-Flash",
		Value: url.QueryEscape(flashValue),
	})
}

func (f *Flash) Get(c *fiber.Ctx) {
	t := fiber.Map{}
	f.Data = nil
	cookieValue := c.Cookies(f.CookiePrefix + "-Flash")
	if cookieValue != "" {
		ParseKeyValueCookie(cookieValue, func(key string, val interface{}) {
			t[key] = val
		})
		f.Data = t
	}
	c.Set("Set-Cookie", f.CookiePrefix+"-Flash=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; HttpOnly")
}

// ParseKeyValueCookie takes the raw (escaped) cookie value and parses out key values.
func ParseKeyValueCookie(val string, cb func(key string, val interface{})) {
	val, _ = url.QueryUnescape(val)
	if matches := cookieKeyValueParser.FindAllStringSubmatch(val, -1); matches != nil {
		for _, match := range matches {
			cb(match[1], match[2])
		}
	}
}

func (f *Flash) Redirect(c *fiber.Ctx, location string, data interface{}, status ...int) {

	f.Data = data.(fiber.Map)
	if len(status) > 0 {
		c.Redirect(location, status[0])
	} else {
		c.Redirect(location, fiber.StatusFound)
	}

}

func (f *Flash) WithError(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.Data = data
	f.Error(c)
	return c
}

func (f *Flash) WithSuccess(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.Data = data
	f.Success(c)
	return c
}

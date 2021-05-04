package flash

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"regexp"
)

type Flash struct {
	CookiePrefix string
	Data         fiber.Map
}

var DefaultFlash = &Flash{
	CookiePrefix: "Fiber-App",
	Data:         fiber.Map{},
}

var cookieKeyValueParser = regexp.MustCompile("\x00([^:]*):([^\x00]*)\x00")

func New(cookiePrefix string) {
	DefaultFlash = &Flash{
		CookiePrefix: cookiePrefix,
		Data:         fiber.Map{},
	}
}

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

func (f *Flash) Get(c *fiber.Ctx) fiber.Map {
	t := fiber.Map{}
	f.Data = nil
	cookieValue := c.Cookies(f.CookiePrefix + "-Flash")
	if cookieValue != "" {
		parseKeyValueCookie(cookieValue, func(key string, val interface{}) {
			t[key] = val
		})
		f.Data = t
	}
	c.Set("Set-Cookie", f.CookiePrefix+"-Flash=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; HttpOnly")
	return f.Data
}

func (f *Flash) Redirect(c *fiber.Ctx, location string, data interface{}, status ...int) error {

	f.Data = data.(fiber.Map)
	if len(status) > 0 {
		return c.Redirect(location, status[0])
	} else {
		return c.Redirect(location, fiber.StatusFound)
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

func Error(c *fiber.Ctx) {
	var flashValue string
	DefaultFlash.Data["error"] = true
	for key, value := range DefaultFlash.Data {
		flashValue += "\x00" + key + ":" + fmt.Sprintf("%v", value) + "\x00"
	}
	c.Cookie(&fiber.Cookie{
		Name:  DefaultFlash.CookiePrefix + "-Flash",
		Value: url.QueryEscape(flashValue),
	})
}

func Success(c *fiber.Ctx) {
	var flashValue string
	DefaultFlash.Data["success"] = true
	for key, value := range DefaultFlash.Data {
		flashValue += "\x00" + key + ":" + fmt.Sprintf("%v", value) + "\x00"
	}
	c.Cookie(&fiber.Cookie{
		Name:  DefaultFlash.CookiePrefix + "-Flash",
		Value: url.QueryEscape(flashValue),
	})
}

func Get(c *fiber.Ctx) fiber.Map {
	t := fiber.Map{}
	DefaultFlash.Data = nil
	cookieValue := c.Cookies(DefaultFlash.CookiePrefix + "-Flash")
	if cookieValue != "" {
		parseKeyValueCookie(cookieValue, func(key string, val interface{}) {
			t[key] = val
		})
		DefaultFlash.Data = t
	}
	c.Set("Set-Cookie", DefaultFlash.CookiePrefix+"-Flash=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; HttpOnly")
	return DefaultFlash.Data
}

func Redirect(c *fiber.Ctx, location string, data interface{}, status ...int) error {

	DefaultFlash.Data = data.(fiber.Map)
	if len(status) > 0 {
		return c.Redirect(location, status[0])
	} else {
		return c.Redirect(location, fiber.StatusFound)
	}
}

func WithError(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	DefaultFlash.Data = data
	DefaultFlash.Error(c)
	return c
}

func WithSuccess(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	DefaultFlash.Data = data
	DefaultFlash.Success(c)
	return c
}

func SetData(c *fiber.Ctx, data fiber.Map) fiber.Map {
	DefaultFlash.Data = data
	return data
}

func WithData(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	DefaultFlash.Data = data
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

package flash

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Flash struct {
	data   fiber.Map
	config Config
}

type Config struct {
	Name        string    `json:"name"`
	Value       string    `json:"value"`
	Path        string    `json:"path"`
	Domain      string    `json:"domain"`
	MaxAge      int       `json:"max_age"`
	Expires     time.Time `json:"expires"`
	Secure      bool      `json:"secure"`
	HTTPOnly    bool      `json:"http_only"`
	SameSite    string    `json:"same_site"`
	SessionOnly bool      `json:"session_only"`
}

var DefaultFlash *Flash

func init() {
	Default(Config{
		Name: "fiber-app-flash",
	})
}

var cookieKeyValueParser = regexp.MustCompile("\x00([^:]*):([^\x00]*)\x00")

func Default(config Config) {
	DefaultFlash = New(config)
}

func New(config Config) *Flash {
	if config.SameSite == "" {
		config.SameSite = "Lax"
	}
	return &Flash{
		config: config,
		data:   fiber.Map{},
	}
}

func (f *Flash) Get(c *fiber.Ctx) fiber.Map {
	t := fiber.Map{}
	f.data = nil
	cookieValue := c.Cookies(f.config.Name)
	if cookieValue != "" {
		parseKeyValueCookie(cookieValue, func(key string, val interface{}) {
			t[key] = val
		})
		f.data = t
	}
	c.Set("Set-Cookie", f.config.Name+"=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; HttpOnly; SameSite="+f.config.SameSite)
	if f.data == nil {
		f.data = fiber.Map{}
	}
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

func (f *Flash) RedirectToRoute(c *fiber.Ctx, routeName string, data fiber.Map, status ...int) error {
	f.data = data
	if len(status) > 0 {
		return c.RedirectToRoute(routeName, data, status[0])
	} else {
		return c.RedirectToRoute(routeName, data, fiber.StatusFound)
	}
}

func (f *Flash) RedirectBack(c *fiber.Ctx, fallback string, data fiber.Map, status ...int) error {
	f.data = data
	if len(status) > 0 {
		return c.RedirectBack(fallback, status[0])
	} else {
		return c.RedirectBack(fallback, fiber.StatusFound)
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

func (f *Flash) WithWarn(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.data = data
	f.warn(c)
	return c
}

func (f *Flash) WithInfo(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.data = data
	f.info(c)
	return c
}

func (f *Flash) WithData(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	f.data = data
	f.setCookie(c)
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

func (f *Flash) warn(c *fiber.Ctx) {
	f.data["warn"] = true
	f.setCookie(c)
}

func (f *Flash) info(c *fiber.Ctx) {
	f.data["info"] = true
	f.setCookie(c)
}

func (f *Flash) setCookie(c *fiber.Ctx) {
	var flashValue string
	for key, value := range f.data {
		flashValue += "\x00" + key + ":" + fmt.Sprintf("%v", value) + "\x00"
	}
	c.Cookie(&fiber.Cookie{
		Name:        f.config.Name,
		Value:       url.QueryEscape(flashValue),
		SameSite:    f.config.SameSite,
		Secure:      f.config.Secure,
		Path:        f.config.Path,
		Domain:      f.config.Domain,
		MaxAge:      f.config.MaxAge,
		Expires:     f.config.Expires,
		HTTPOnly:    f.config.HTTPOnly,
		SessionOnly: f.config.SessionOnly,
	})
}

func Get(c *fiber.Ctx) fiber.Map {
	return DefaultFlash.Get(c)
}

func Redirect(c *fiber.Ctx, location string, data interface{}, status ...int) error {
	return DefaultFlash.Redirect(c, location, data, status...)
}

func RedirectToRoute(c *fiber.Ctx, routeName string, data fiber.Map, status ...int) error {
	return DefaultFlash.RedirectToRoute(c, routeName, data, status...)
}

func RedirectBack(c *fiber.Ctx, fallback string, data fiber.Map, status ...int) error {
	return DefaultFlash.RedirectBack(c, fallback, data, status...)
}

func WithError(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithError(c, data)
}

func WithSuccess(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithSuccess(c, data)
}

func WithWarn(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithWarn(c, data)
}

func WithInfo(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithInfo(c, data)
}

func WithData(c *fiber.Ctx, data fiber.Map) *fiber.Ctx {
	return DefaultFlash.WithData(c, data)
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

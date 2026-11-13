package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/language"
)

const (
	DefaultLanguage = "en"
	CookieName      = "lang"
	cookieMaxAge    = 60 * 60 * 24 * 365
)

//go:embed locales/*.json
var localeFS embed.FS

var (
	supportedTags = []language.Tag{language.English, language.Arabic}
	matcher       = language.NewMatcher(supportedTags)
	catalogs      = loadCatalogs()
)

type Translator struct {
	Lang     string
	Dir      string
	messages map[string]string
	fallback map[string]string
}

func New(lang string) Translator {
	lang = Normalize(lang)
	messages := catalogs[lang]
	if messages == nil {
		lang = DefaultLanguage
		messages = catalogs[lang]
	}

	return Translator{
		Lang:     lang,
		Dir:      Direction(lang),
		messages: messages,
		fallback: catalogs[DefaultLanguage],
	}
}

func (t Translator) T(key string, args ...any) string {
	message, ok := t.messages[key]
	if !ok {
		message, ok = t.fallback[key]
	}
	if !ok {
		return key
	}

	if len(args) == 0 {
		return message
	}

	return fmt.Sprintf(message, args...)
}

func (t Translator) Date(created time.Time) string {
	month := t.T(fmt.Sprintf("date.month.%d", int(created.Month())))

	return fmt.Sprintf(
		"%02d %s %d · %s",
		created.Day(),
		month,
		created.Year(),
		created.Format("03:04 PM"),
	)
}

func Match(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}

	tag, err := language.Parse(value)
	if err != nil {
		return "", false
	}

	matched, _, confidence := matcher.Match(tag)
	if matched == language.Und || confidence == language.No {
		return "", false
	}

	base, _ := matched.Base()

	return base.String(), true
}

func Normalize(value string) string {
	if lang, ok := Match(value); ok {
		return lang
	}

	return DefaultLanguage
}

func Resolve(r *http.Request) string {
	if lang, ok := Match(r.URL.Query().Get("lang")); ok {
		return lang
	}

	if cookie, err := r.Cookie(CookieName); err == nil {
		if lang, ok := Match(cookie.Value); ok {
			return lang
		}
	}

	if tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language")); err == nil && len(tags) > 0 {
		matched, _, confidence := matcher.Match(tags...)
		if matched != language.Und && confidence != language.No {
			base, _ := matched.Base()

			return base.String()
		}
	}

	return DefaultLanguage
}

func Direction(lang string) string {
	tag, err := language.Parse(lang)
	if err != nil {
		return "ltr"
	}

	base, _ := tag.Base()

	switch base.String() {
	case "ar", "he", "fa", "ur":
		return "rtl"
	default:
		return "ltr"
	}
}

func Cookie(lang string) *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    Normalize(lang),
		Path:     "/",
		MaxAge:   cookieMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func loadCatalogs() map[string]map[string]string {
	entries, err := fs.ReadDir(localeFS, "locales")
	if err != nil {
		panic(fmt.Sprintf("read locale directory: %v", err))
	}

	result := make(map[string]map[string]string, len(entries))

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := fs.ReadFile(localeFS, "locales/"+entry.Name())
		if err != nil {
			panic(fmt.Sprintf("read locale %s: %v", entry.Name(), err))
		}

		var raw map[string]any
		if err := json.Unmarshal(data, &raw); err != nil {
			panic(fmt.Sprintf("parse locale %s: %v", entry.Name(), err))
		}

		messages := make(map[string]string)
		flattenJSON("", raw, messages)

		lang := strings.TrimSuffix(entry.Name(), ".json")
		result[lang] = messages
	}

	if result[DefaultLanguage] == nil {
		panic("missing default locale")
	}

	return result
}

func flattenJSON(prefix string, value any, out map[string]string) {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			next := key
			if prefix != "" {
				next = prefix + "." + key
			}
			flattenJSON(next, child, out)
		}
	case string:
		out[prefix] = typed
	default:
		panic(fmt.Sprintf("invalid translation value at %q", prefix))
	}
}

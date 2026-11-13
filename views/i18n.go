package views

import (
	"net/http"

	"to-do/internal/i18n"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

var errorKeys = map[string]string{
	"Cannot be blank.":                             "error.cannot_be_blank",
	"Failed to read request.":                      "error.read_request",
	"Failed to load todos.":                        "error.load_todos",
	"Failed to render todos.":                      "error.render_todos",
	"Failed to build list rule.":                   "error.build_list_rule",
	"Failed to apply list rule.":                   "error.apply_list_rule",
	"Failed to render home page.":                  "error.render_home",
	"Missing todos collection.":                    "error.missing_todos",
	"You are not allowed to perform this request.": "error.forbidden",
	"Only superusers can perform this action.":     "error.superuser_only",
}

func translatorFor(e *core.RequestEvent) i18n.Translator {
	if raw := e.Request.URL.Query().Get("lang"); raw != "" {
		if lang, ok := i18n.Match(raw); ok {
			e.SetCookie(i18n.Cookie(lang))
		}
	}

	return i18n.New(i18n.Resolve(e.Request))
}

func (v *Views) SetLanguage(e *core.RequestEvent) error {
	lang := i18n.Normalize(e.Request.FormValue("lang"))

	e.SetCookie(i18n.Cookie(lang))

	return e.NoContent(http.StatusNoContent)
}

func localizeError(message string, tr i18n.Translator) string {
	if key, ok := errorKeys[message]; ok {
		return tr.T(key)
	}

	if tr.Lang != i18n.DefaultLanguage {
		return tr.T("error.generic")
	}

	return message
}

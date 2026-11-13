package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTFormatsArguments(t *testing.T) {
	en := New("en")
	if got, want := en.T("todo.delete_confirm", "Write docs"), `Delete "Write docs"?`; got != want {
		t.Fatalf("unexpected English translation: %q", got)
	}

	ar := New("ar")
	if got, want := ar.T("todo.delete_confirm", "كتابة الوثائق"), `حذف "كتابة الوثائق"؟`; got != want {
		t.Fatalf("unexpected Arabic translation: %q", got)
	}
}

func TestTFallsBackToKey(t *testing.T) {
	if got := New("en").T("missing.key"); got != "missing.key" {
		t.Fatalf("expected missing key fallback, got %q", got)
	}
}

func TestResolvePriority(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		cookie string
		accept string
		want   string
	}{
		{name: "query wins", query: "en", cookie: "ar", accept: "ar", want: "en"},
		{name: "cookie wins over header", cookie: "ar", accept: "en", want: "ar"},
		{name: "accept language", accept: "ar-EG,ar;q=0.9,en;q=0.5", want: "ar"},
		{name: "default", accept: "fr-FR,fr;q=0.9", want: "en"},
		{name: "invalid query falls through", query: "fr", cookie: "ar", want: "ar"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := "/"
			if test.query != "" {
				target += "?lang=" + test.query
			}

			req := httptest.NewRequest("GET", target, nil)
			if test.cookie != "" {
				req.AddCookie(&http.Cookie{Name: CookieName, Value: test.cookie})
			}
			if test.accept != "" {
				req.Header.Set("Accept-Language", test.accept)
			}

			if got := Resolve(req); got != test.want {
				t.Fatalf("Resolve() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestNormalizeAndDirection(t *testing.T) {
	if got := Normalize("ar-EG"); got != "ar" {
		t.Fatalf("Normalize(ar-EG) = %q", got)
	}
	if got := Normalize("fr"); got != "en" {
		t.Fatalf("Normalize(fr) = %q", got)
	}
	if got := Direction("ar"); got != "rtl" {
		t.Fatalf("Direction(ar) = %q", got)
	}
	if got := Direction("en"); got != "ltr" {
		t.Fatalf("Direction(en) = %q", got)
	}
}

func TestDate(t *testing.T) {
	created := time.Date(2026, time.November, 13, 16, 5, 0, 0, time.UTC)

	if got, want := New("en").Date(created), "13 Nov 2026 · 04:05 PM"; got != want {
		t.Fatalf("English date = %q, want %q", got, want)
	}
	if got, want := New("ar").Date(created), "13 نوفمبر 2026 · 04:05 PM"; got != want {
		t.Fatalf("Arabic date = %q, want %q", got, want)
	}
}

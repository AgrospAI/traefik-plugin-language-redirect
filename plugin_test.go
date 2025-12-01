package traefik_plugin_language_redirect

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

var langRedirectTests = []struct {
	name                 string
	url                  string
	cookieLanguage       string
	headerAcceptLanguage string
	config               *Config
	expectedLocation     string
}{
	{
		name:                 "Redirect to cookie language",
		url:                  "http://example.com/page",
		cookieLanguage:       "fr",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/fr/page",
	},
	{
		name:                 "Redirect to header accept language",
		url:                  "http://example.com/page",
		cookieLanguage:       "",
		headerAcceptLanguage: "de",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/de/page",
	},
	{
		name:                 "Redirect to default language",
		url:                  "http://example.com/page",
		cookieLanguage:       "",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/en/page",
	},
	{
		name:                 "No redirect needed, same as cookie language",
		url:                  "http://example.com/fr/page",
		cookieLanguage:       "fr",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/fr/page",
	},
	{
		name:                 "No redirect needed, same as header accept language",
		url:                  "http://example.com/de/page",
		cookieLanguage:       "",
		headerAcceptLanguage: "de",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/de/page",
	},
	{
		name:                 "No redirect needed, same as default language",
		url:                  "http://example.com/en/page",
		cookieLanguage:       "",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/en/page",
	},
	{
		name:                 "Redirect keep query parameters",
		url:                  "http://example.com/page?param=value",
		cookieLanguage:       "fr",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/fr/page?param=value",
	},
	{
		name:                 "Redirect keep fragment",
		url:                  "http://example.com/page#section1",
		cookieLanguage:       "de",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/de/page#section1",
	},
	{
		name:                 "No redirect needed, trailing slash in URL",
		url:                  "http://example.com/en/",
		cookieLanguage:       "",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/en/",
	},
	{
		name:                 "Redirect to root language",
		url:                  "http://example.com/es/page",
		cookieLanguage:       "en",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "es", "de"},
		},
		expectedLocation: "http://example.com/page",
	},
	{
		name:                 "Works with localhost",
		url:                  "http://localhost/page",
		cookieLanguage:       "fr",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://localhost/fr/page",
	},
	{
		name:                 "Works with port in URL",
		url:                  "http://example.com:8080/page",
		cookieLanguage:       "de",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com:8080/de/page",
	},
	{
		name:                 "Works with root path",
		url:                  "http://example.com/",
		cookieLanguage:       "fr",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/fr/",
	},
	{
		name:                 "Redirect keeps trailing slash",
		url:                  "http://example.com/docs/page/",
		cookieLanguage:       "de",
		headerAcceptLanguage: "",
		config: &Config{
			CookieName:         "lang",
			DefaultLanguage:    "en",
			RootLanguage:       "en",
			SupportedLanguages: []string{"en", "fr", "de"},
		},
		expectedLocation: "http://example.com/de/docs/page/",
	},
}

func TestLangRedirect(t *testing.T) {
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	for _, tt := range langRedirectTests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new handler for this test's config
			handler, err := New(ctx, next, tt.config, "lang-redirect")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set cookie if specified
			if tt.cookieLanguage != "" {
				req.AddCookie(&http.Cookie{
					Name:  tt.config.CookieName,
					Value: tt.cookieLanguage,
				})
			}

			// Set Accept-Language header if specified
			if tt.headerAcceptLanguage != "" {
				req.Header.Set("Accept-Language", tt.headerAcceptLanguage)
			}

			handler.ServeHTTP(recorder, req)

			// Determine if redirect is expected
			if tt.expectedLocation != "" && tt.expectedLocation != tt.url {
				assertRedirection(t, recorder, tt.expectedLocation)
			} else {
				assertNoRedirection(t, recorder)
			}
		})
	}
}

func assertRedirection(t *testing.T, recorder *httptest.ResponseRecorder, location string) {
	assertStatusCode(t, recorder, 302)
	assertHeader(t, recorder, "Location", location)
}

func assertNoRedirection(t *testing.T, recorder *httptest.ResponseRecorder) {
	assertStatusCode(t, recorder, 200)
	assertHeader(t, recorder, "Location", "")
}

func assertStatusCode(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if recorder.Code != expected {
		t.Errorf("expected status code %d, got %d", expected, recorder.Code)
	}
}

func assertHeader(t *testing.T, recorder *httptest.ResponseRecorder, key, expected string) {
	t.Helper()

	actual := recorder.Header().Get(key)
	if actual != expected {
		t.Errorf("expected header %s to be %q, got %q", key, expected, actual)
	}
}

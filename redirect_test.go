package traefik_plugin_language_redirect

import (
	"net/url"
	"testing"
)

func TestGetLangFromPath(t *testing.T) {
	supportedLangs := []string{"en", "fr", "es"}

	tests := []struct {
		name      string
		inputPath string
		wantLang  string
		wantPath  string
	}{
		{"root path", "/", "", "/"},
		{"root path empty", "", "", ""},
		{"language only", "/en", "en", ""},
		{"language only with slash", "/en/", "en", "/"},
		{"language with page", "/en/page", "en", "/page"},
		{"language with page and slash", "/en/page/", "en", "/page/"},
		{"unsupported language", "/de", "", "/de"},
		{"unsupported language with page", "/de/page", "", "/de/page"},
		{"nested path without language", "/page/subpage", "", "/page/subpage"},
		{"nested path with language in middle", "/page/en", "", "/page/en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := url.URL{Path: tt.inputPath}
			gotLang, gotURL := getLangFromPath(u, supportedLangs)

			if gotLang != tt.wantLang {
				t.Errorf("gotLang = %q, want %q", gotLang, tt.wantLang)
			}
			if gotURL.Path != tt.wantPath {
				t.Errorf("gotURL.Path = %q, want %q", gotURL.Path, tt.wantPath)
			}
		})
	}
}

func TestPrependLangToPath(t *testing.T) {
	tests := []struct {
		name      string
		lang      string
		inputPath string
		wantPath  string
	}{
		{"root path", "en", "/", "/en/"},
		{"empty path", "en", "", "/en"},
		{"single segment", "en", "/page", "/en/page"},
		{"single segment with slash", "en", "/page/", "/en/page/"},
		{"nested path", "fr", "/page/subpage", "/fr/page/subpage"},
		{"nested path with slash", "fr", "/page/subpage/", "/fr/page/subpage/"},
		{"already has language", "es", "/es/page", "/es/es/page"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := url.URL{Path: tt.inputPath}
			got := prependLangToPath(tt.lang, u)

			if got.Path != tt.wantPath {
				t.Errorf("prependLangToPath(%q, %q) = %q, want %q",
					tt.lang, tt.inputPath, got.Path, tt.wantPath)
			}
		})
	}
}

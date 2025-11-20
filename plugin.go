package plugin

import (
	"context"
	"fmt"
	"net/http"
	"slices"
)

type Config struct {
	cookieName         string
	defaultLanguage    string
	rootLanguage       string
	supportedLanguages []string
}

type LanguageRedirect struct {
	next   http.Handler
	config *Config
}

func CreateConfig() *Config {
	return &Config{
		cookieName:         "",
		defaultLanguage:    "",
		rootLanguage:       "",
		supportedLanguages: []string{},
	}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.supportedLanguages) == 0 {
		return nil, fmt.Errorf("supportedLanguages cannot be empty")
	}

	if config.defaultLanguage == "" {
		return nil, fmt.Errorf("defaultLanguage cannot be empty")
	}

	if config.cookieName == "" {
		return nil, fmt.Errorf("cookieName cannot be empty")
	}

	if config.rootLanguage != "" && !slices.Contains(config.supportedLanguages, config.rootLanguage) {
		return nil, fmt.Errorf("rootLanguage %s is not in supportedLanguages", config.rootLanguage)
	}

	if !slices.Contains(config.supportedLanguages, config.defaultLanguage) {
		return nil, fmt.Errorf("defaultLanguage %s is not in supportedLanguages", config.defaultLanguage)
	}

	return &LanguageRedirect{
		next:   next,
		config: config,
	}, nil
}

func (a *LanguageRedirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	cookie, err := req.Cookie(a.config.cookieName)
	preferredLang := req.Header.Get("Accept-Language")

	redirectUrl, err := RedirectURL(RedirectOptions{
		URL: url,
		CookieLang: func() string {
			if err == nil {
				return cookie.Value
			} else {
				return ""
			}
		}(),
		AcceptLang:     preferredLang,
		DefaultLang:    a.config.defaultLanguage,
		RootLang:       a.config.rootLanguage,
		SupportedLangs: a.config.supportedLanguages,
	})
	if err == nil && redirectUrl != url {
		http.Redirect(rw, req, redirectUrl, http.StatusFound)
		return
	}

	a.next.ServeHTTP(rw, req)
}

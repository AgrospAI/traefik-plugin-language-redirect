package traefik_plugin_language_redirect

import (
	"context"
	"fmt"
	"net/http"
	"slices"
)

type Config struct {
	CookieName         string
	DefaultLanguage    string
	RootLanguage       string
	SupportedLanguages []string
}

type LanguageRedirect struct {
	next   http.Handler
	config *Config
}

func CreateConfig() *Config {
	return &Config{
		CookieName:         "",
		DefaultLanguage:    "",
		RootLanguage:       "",
		SupportedLanguages: []string{},
	}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// Print configuration for debugging
	// fmt.Fprintf(os.Stdout, "Configuration: %+v\n", config)

	if len(config.SupportedLanguages) == 0 {
		return nil, fmt.Errorf("supportedLanguages cannot be empty")
	}

	if config.DefaultLanguage == "" {
		return nil, fmt.Errorf("defaultLanguage cannot be empty")
	}

	if config.CookieName == "" {
		return nil, fmt.Errorf("cookieName cannot be empty")
	}

	if config.RootLanguage != "" && !slices.Contains(config.SupportedLanguages, config.RootLanguage) {
		return nil, fmt.Errorf("rootLanguage %s is not in supportedLanguages", config.RootLanguage)
	}

	if !slices.Contains(config.SupportedLanguages, config.DefaultLanguage) {
		return nil, fmt.Errorf("defaultLanguage %s is not in supportedLanguages", config.DefaultLanguage)
	}

	return &LanguageRedirect{
		next:   next,
		config: config,
	}, nil
}

func (a *LanguageRedirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}

	url := fmt.Sprintf("%s://%s%s", scheme, req.Host, req.URL.RequestURI())
	cookie, err := req.Cookie(a.config.CookieName)
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
		DefaultLang:    a.config.DefaultLanguage,
		RootLang:       a.config.RootLanguage,
		SupportedLangs: a.config.SupportedLanguages,
	})
	if err == nil && redirectUrl != url {
		http.Redirect(rw, req, redirectUrl, http.StatusFound)
		return
	}

	a.next.ServeHTTP(rw, req)
}

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

func (lr *LanguageRedirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	result, err := RedirectURL(RedirectOptions{
		URL:            *req.URL,
		CookieLang:     getCookieValue(req, lr.config.CookieName),
		AcceptLang:     req.Header.Get("Accept-Language"),
		DefaultLang:    lr.config.DefaultLanguage,
		RootLang:       lr.config.RootLanguage,
		SupportedLangs: lr.config.SupportedLanguages,
	})

	if err == nil && result.ShouldRedirect {
		http.Redirect(rw, req, result.Target.String(), http.StatusFound)
		return
	}

	lr.next.ServeHTTP(rw, req)
}

func getCookieValue(req *http.Request, name string) string {
	c, _ := req.Cookie(name)
	if c != nil {
		return c.Value
	}
	return ""
}

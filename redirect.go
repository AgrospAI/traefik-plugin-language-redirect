package traefik_plugin_language_redirect

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

type RedirectOptions struct {
	URL            string
	CookieLang     string
	AcceptLang     string
	DefaultLang    string
	RootLang       string
	SupportedLangs []string
}

func RedirectURL(opts RedirectOptions) (string, error) {
	url := opts.URL
	cookieLang := opts.CookieLang
	acceptLang := opts.AcceptLang
	defaultLang := opts.DefaultLang
	rootLang := opts.RootLang
	supportedLangs := opts.SupportedLangs

	if cookieLang != "" && !slices.Contains(supportedLangs, cookieLang) {
		cookieLang = ""
	}

	if acceptLang != "" && !slices.Contains(supportedLangs, acceptLang) {
		acceptLang = ""
	}

	if !slices.Contains(supportedLangs, defaultLang) {
		return "", fmt.Errorf("defaultLang %s is not in supportedLangs", defaultLang)
	}

	if rootLang != "" && !slices.Contains(supportedLangs, rootLang) {
		return "", fmt.Errorf("rootLang %s is not in supportedLangs", rootLang)
	}

	root, path, err := GetRootAndPath(url)
	if err != nil {
		return "", err
	}

	_, restOfPath := GetLangPath(path, supportedLangs)

	if cookieLang != "" {
		if cookieLang == rootLang {
			return root + restOfPath, nil
		}
		return root + "/" + cookieLang + restOfPath, nil
	}

	if acceptLang != "" {
		if acceptLang == rootLang {
			return root + restOfPath, nil
		}
		return root + "/" + acceptLang + restOfPath, nil
	}

	if defaultLang == rootLang {
		return root + restOfPath, nil
	}
	return root + "/" + defaultLang + restOfPath, nil
}

func GetRootAndPath(raw string) (string, string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}

	root := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	path := u.Path

	return root, path, nil
}

func GetLangPath(path string, supportedLangs []string) (string, string) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	if len(parts) == 0 {
		return "", path
	}

	firstPart := parts[0]

	for _, lang := range supportedLangs {
		if firstPart == lang {
			rest := "/" + strings.Join(parts[1:], "/")
			if rest == "/" {
				rest = ""
			}
			return lang, rest
		}
	}

	return "", path
}

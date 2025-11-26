package traefik_plugin_language_redirect

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

type RedirectOptions struct {
	URL            url.URL
	CookieLang     string
	AcceptLang     string
	DefaultLang    string
	RootLang       string
	SupportedLangs []string
}

type RedirectResult struct {
	ShouldRedirect bool
	Target         url.URL
}

func RedirectURL(opts RedirectOptions) (RedirectResult, error) {
	// Validate inputs
	if opts.CookieLang != "" && !slices.Contains(opts.SupportedLangs, opts.CookieLang) {
		opts.CookieLang = ""
	}

	if opts.AcceptLang != "" && !slices.Contains(opts.SupportedLangs, opts.AcceptLang) {
		opts.AcceptLang = ""
	}

	if !slices.Contains(opts.SupportedLangs, opts.DefaultLang) {
		return RedirectResult{false, opts.URL}, fmt.Errorf("defaultLang %s is not in supportedLangs", opts.DefaultLang)
	}

	if opts.RootLang != "" && !slices.Contains(opts.SupportedLangs, opts.RootLang) {
		return RedirectResult{false, opts.URL}, fmt.Errorf("rootLang %s is not in supportedLangs", opts.RootLang)
	}

	// Determine preferred language
	preferredLang := ""
	if opts.CookieLang != "" {
		preferredLang = opts.CookieLang
	} else if opts.AcceptLang != "" {
		preferredLang = opts.AcceptLang
	} else {
		preferredLang = opts.DefaultLang
	}

	// Check current path language
	pathLang, rootUrl := getLangFromPath(opts.URL, opts.SupportedLangs)
	if pathLang == "" {
		pathLang = opts.RootLang
	}

	// Redirect if needed
	if preferredLang != pathLang {
		if preferredLang == opts.RootLang {
			return RedirectResult{true, rootUrl}, nil
		}
		return RedirectResult{true, prependLangToPath(preferredLang, rootUrl)}, nil
	}

	return RedirectResult{false, opts.URL}, nil
}

func prependLangToPath(lang string, url url.URL) url.URL {
	path := url.Path
	segments := strings.Split(strings.Trim(path, "/"), "/")
	segments = append([]string{lang}, segments...)
	url.Path = "/" + strings.Join(segments, "/")
	return url
}

func getLangFromPath(url url.URL, supportedLangs []string) (string, url.URL) {
	path := url.Path
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) > 0 {
		firstSegment := segments[0]
		if slices.Contains(supportedLangs, firstSegment) {
			restOfPath := "/" + strings.Join(segments[1:], "/")
			url.Path = restOfPath
			return firstSegment, url
		}
	}
	return "", url
}

# Language Redirect Plugin for Traefik

## Installation

To install the plugin in Traefik, you can use the following static configuration (given here in YAML):

```yaml
experimental:
  plugins:
    language-redirect:
      moduleName: github.com/AgrospAI/traefik-plugin-language-redirect
      version: v0.0.1
```

## Usage

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the `http.middlewares` section:

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - my-plugin

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  
  middlewares:
    language-redirect:
      plugin:
        language-redirect:
          cookieName: lang
          defaultLanguage: en
          rootLanguage: es # The / path is in Spanish (there is no /es prefix)
          supportedLanguages:
            - en
            - es
            - de
```

## Configuration

The plugin can be configured with the following parameters:

- `cookieName` (string): Name of the cookie to check for language preference.
- `defaultLanguage` (string): Default language to redirect to if no preference is found.
- `supportedLanguages` ([]string): List of supported languages for redirection.
- `rootLanguage` (string, optional): Language to use for root URL redirection.

## Development

You can use Visual Studio Code Dev Containers for a consistent development environment. Open the command palette and select `Dev Containers: Reopen in Container`.

> Check the [Makefile](./Makefile) for useful commands to lint, test, and format the code.

## References

- <https://plugins.traefik.io/create>
- <https://github.com/dcasia/plugin-cond-redirect>

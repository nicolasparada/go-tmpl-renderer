[![Go Reference](https://pkg.go.dev/badge/github.com/nicolasparada/go-tmpl-renderer.svg)](https://pkg.go.dev/github.com/nicolasparada/go-tmpl-renderer)

# Golang Template Renderer

An opinionated HTML template renderer for Golang.

## Install

```bash
go get github.com/nicolasparada/go-tmpl-renderer
```

## Usage

`templates/includes/layout.tmpl`

```handlebars
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Golang Template Renderer Demo</title>
</head>
<body>
    {{template "header.tmpl" .}}
    {{block "content" .}}{{end}}
</body>
</html>
```

`templates/includes/header.tmpl`

```handlebars
<header>
  <a href="/"><h1>Home</h1></a>
</header>
```

`templates/welcome.tmpl`

```handlebars
{{template "layout.tmpl" .}}

{{define "content"}}
  <main>
    <h1>Welcome</h1>
  </main>
{{end}}

```

`main.go`

```go
package main

import (
    // other imports
    tmplrenderer "github.com/nicolasparada/go-tmpl-render"
)

//go:embed templates/includes/*.tmpl templates/*.tmpl
var templatesFS embed.FS

func main() {
    renderer := &tmplrenderer.Renderer{
        FS:             templatesFS,
        BaseDir:        "templates",
        IncludePatters: []string{"includes/*.tmpl"},
    }

    http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
        renderer.Render(w, "welcome.tmpl", nil)
    })
    http.ListenAndServe(":4000")
}
```

## Included Functions

HTML templates will include two functions from the start: `dict` and `list`.

- `dict` given a list of key-value pairs, it will construct a map `map[string]any` out of it.<br>
  Example:

```
{{$_ := dict "foo" "bar"}} => map[string]any{"foo": "bar"}
```

It's very usefull for passing named parameters to other templates, for example:

```handlebars
{{define "greet"}}
    Hello, {{.Name}}!
{{end}}

{{template "greet" dict
  "Name" "john"
}}
```

- `list`: given a list of items, it will construct an slice `[]any` out of it.<br>
  Example:

```
{{$_ := list "foo" "bar"}} => []any{"foo", "bar"}
```

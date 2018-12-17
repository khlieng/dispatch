package server

type indexTemplateData struct {
	InlineScript string
	CSSPath      string
	Scripts      []string
}

const indexTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="theme-color" content="#f0f0f0">

  <title>Dispatch</title>
  <meta name="description" content="Web-based IRC client.">

  <link rel="preload" href="/init" as="fetch" crossorigin>

	{{if .InlineScript}}
  <script>{{.InlineScript}}</script>
	{{end}}

  <link rel="preload" href="/font/fontello.woff2?48901973" as="font" type="font/woff2" crossorigin>
  <link rel="preload" href="/font/RobotoMono-Regular.woff2" as="font" type="font/woff2" crossorigin>
  <link rel="preload" href="/font/Montserrat-Regular.woff2" as="font" type="font/woff2" crossorigin>
  <link rel="preload" href="/font/Montserrat-Bold.woff2" as="font" type="font/woff2" crossorigin>
  <link rel="preload" href="/font/RobotoMono-Bold.woff2" as="font" type="font/woff2" crossorigin>

  {{if .CSSPath}}
  <link href="/{{.CSSPath}}" rel="stylesheet">
  {{end}}

  <link rel="manifest" href="/manifest.json">
</head>

<body>
  <div id="root"></div>

  {{range .Scripts}}
  <script src="/{{.}}"></script>
  {{end}}

  <noscript>This page needs JavaScript enabled to function.</noscript>
</body>

</html>`

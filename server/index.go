package server

type indexTemplateData struct {
	InlineScript string
	Stylesheet   string
	Scripts      []string
}

const indexTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="theme-color" content="#222">

  <title>Dispatch</title>
  <meta name="description" content="Web-based IRC client.">

  <link rel="preload" href="/init" as="fetch" crossorigin>

  {{if .InlineScript}}
  <script>{{.InlineScript}}</script>
  {{end}}

  {{range .Scripts}}
  <script src="{{.}}" defer></script>
  {{end}}

  <link rel="preload" href="/font/RobotoMono-Regular.woff2" as="font" type="font/woff2" crossorigin>
  <link rel="preload" href="/font/Montserrat-Regular.woff2" as="font" type="font/woff2" crossorigin>
  <link rel="preload" href="/font/Montserrat-Bold.woff2" as="font" type="font/woff2" crossorigin>
  <link rel="preload" href="/font/RobotoMono-Bold.woff2" as="font" type="font/woff2" crossorigin>

  {{if .Stylesheet}}
  <link href="{{.Stylesheet}}" rel="stylesheet">
  {{end}}

  <link rel="manifest" href="/manifest.json">
  <link rel="apple-touch-icon" href="/icon_192.png">
</head>

<body>
  <noscript>This page needs JavaScript enabled to function.</noscript>
  <div id="root"></div>
</body>

</html>`

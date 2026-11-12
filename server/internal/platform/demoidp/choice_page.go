package demoidp

import "html/template"

var choicePage = template.Must(template.New("choice").Parse(`<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Демо-вход Lodestar</title>
<style>
body{font-family:system-ui,sans-serif;background:#f4f4f5;color:#18181b;margin:0;padding:24px}
main{max-width:420px;margin:40px auto;background:#fff;border:1px solid #e4e4e7;border-radius:16px;padding:24px}
h1{font-size:20px;margin:0 0 4px}p{color:#52525b;margin:0 0 20px}
button{display:block;width:100%;text-align:left;padding:12px 16px;margin:0 0 8px;border:1px solid #e4e4e7;border-radius:10px;background:#fff;font:inherit;cursor:pointer}
button:hover{border-color:#0b786e}small{display:block;color:#52525b}
</style>
</head>
<body>
<main>
<h1>Демо-вход</h1>
<p>Тестовый вуз. Выберите, кем войти. Пароль не нужен.</p>
{{range .Users}}
<form method="post" action="authorize">
{{range $key, $values := $.Query}}{{range $values}}<input type="hidden" name="{{$key}}" value="{{.}}">{{end}}{{end}}
<input type="hidden" name="subject" value="{{.Subject}}">
<button type="submit">{{.Name}}<small>{{.Title}} · {{.Email}}</small></button>
</form>
{{end}}
</main>
</body>
</html>
`))

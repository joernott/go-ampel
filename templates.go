package main

type Page struct {
	Title string
	Body  string
}

const html_page string = `
<!DOCTYPE html>
<html>
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <title>{{ .Title }}</title>
  <link type="text/css" rel="stylesheet" href="/static/style.css">
 </head>
 <body>
{{ .Body }}
 </body>
</html>
`
const list_services string = `
   <table class="list">
    <tr>
	 <th>Service</th>
	 <th>Status</th>
	 <th>User</th>
	 <th>Until</th>
	 <th>Action</th>
	</tr>
	{{ range . }}
	<tr>
	 <td><a href="/service/{{ .Name }}">{{ .Name }}</a></td>
	 {{ if .Status }}
	 <td><img src="/static/red.png" alt="in use" /></td>
	<td>{{ .Who }}</td>
	<td>{{ .FUntil }}</td>
	 {{ else }}
	 <td><img src="/static/green.png" alt="free" /></td>
	<td>&nbsp;</td>
	<td>&nbsp;</td>
	 {{ end }}
	 <td>
	 {{ if .Status }}
	  <form action="/service/{{ .Name }}/go" method="post" >
	   Who: <input type="text" name="who" size="10" />
	   <button type="submit">Release</button>
	  </form>
	{{ else }}
	  <form action="/service/{{ .Name }}/stop" method="post" >
	   Who: <input type="text" name="who" size="10" />
	   Until:<input type="text" name="until" size="20"/>
	   <button type="submit">Lock</button>
	  </form>
	{{ end }}
	 </td>
	</tr>
	{{ end }}
</table>
`
const svc_not_found string = `
   <p class="error">Could not find the service {{ . }}.</p>
`

const svc_locked string = `
   <p class="locked">
    <img src="/static/red.png" alt="in use" class="status" />
    The service {{ .Name }} is currently in use by {{ .Who }}. This lock is valid until {{ .FUntil }}.
   </p>
   <form action="/service/{{ .Name }}/go" method="post" >
    Who: <input type="text" name="who" size="10" />
	<button type="submit">Release</button>
   </form>
`
const svc_free string = `
   <p class="locked">
    <img src="/static/green.png" alt="free" class="status" />
    The service {{ .Name }} is currently free.
   </p>
   <form action="/service/{{ .Name }}/stop" method="post" >
    Who: <input type="text" name="who" size="10" />
    Until:<input type="text" name="until" size="20"/>
    <button type="submit">Lock</button>
   </form>
`

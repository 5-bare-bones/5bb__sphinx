For more information about a command, its flags and examples please visit the [commands folder](https://github.com/5-bare-bones/5bb__sphinx/tree/master/docs/commands).

{{ range .Commands -}}
- [{{ .Name }}](#{{ .Name }})
{{ end }}
{{ range .Commands }}
### {{ .Name }}

```
{{ cmdAndFlags . -}}
{{ subCmds . "" }}
```

---
{{ end }}

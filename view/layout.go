package view

import (
	"html/template"
	"net/http"

	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/internal"
)

func RenderCustomPage(reg *admin.Registry, w http.ResponseWriter, r *http.Request, title string, content template.HTML) {
	user, _ := internal.GetUserFromRequest(reg, r)
	styleContent, _ := admin.TemplateFS.ReadFile("templates/style.css")
	tmpl := template.Must(template.ParseFS(admin.TemplateFS, "templates/layout.html"))
	tmpl = template.Must(tmpl.New("title").Parse(title))
	tmpl = template.Must(tmpl.New("content").Parse(`<div style="padding: 2rem;">` + string(content) + `</div>`))
	pd := PageData{
		SiteTitle: reg.Config.SiteTitle, GroupedResources: reg.GetGroupedResources(), GroupedPages: reg.GetGroupedPages(),
		User: user, CSS: template.CSS(styleContent),
		Flash: reg.GetFlash(w, r),
	}
	if err := tmpl.ExecuteTemplate(w, "layout", pd); err != nil {
		http.Error(w, "Template error", 500)
		return
	}
}

package view

import (
	"fmt"
	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
	"html/template"
	"net/http"
)

func RenderDashboard(reg *admin.Registry, w http.ResponseWriter, r *http.Request, user *models.AdminUser) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var stats []Stat
	for name, res := range reg.Resources {
		var count int64
		reg.DB.Model(res.Model).Count(&count)
		stats = append(stats, Stat{Label: name, Value: count})
	}
	var widgets []ChartWidget
	for i, c := range reg.Charts {
		l, v := c.Data(reg.DB)
		widgets = append(widgets, ChartWidget{ID: fmt.Sprintf("chart-%d", i), Label: c.Label, Type: c.Type, Labels: l, Values: v})
	}
	styleContent, _ := admin.TemplateFS.ReadFile("templates/style.css")
	tmpl := LoadTemplates("templates/dashboard.html")
	pd := PageData{
		SiteTitle: reg.Config.SiteTitle, GroupedResources: reg.GetGroupedResources(), GroupedPages: reg.GetGroupedPages(), 
		User: user, Stats: stats, CSS: template.CSS(styleContent), ChartData: widgets,
		Flash: reg.GetFlash(w, r),
	}
	tmpl.ExecuteTemplate(w, "dashboard.html", pd)
}

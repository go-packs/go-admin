package resource

import (
	"html/template"
	"testing"
)

type MockModel struct {
	ID   uint
	Name string
}

func TestResource(t *testing.T) {
	t.Run("NewResource", func(t *testing.T) {
		res := NewResource(MockModel{})
		if res.Name != "MockModel" {
			t.Errorf("Expected 'MockModel', got %s", res.Name)
		}
		if res.Path != "/MockModel" {
			t.Errorf("Expected '/MockModel', got %s", res.Path)
		}
	})

	t.Run("RegisterField", func(t *testing.T) {
		res := NewResource(MockModel{})
		res.RegisterField("Name", "Display Name", false)
		if len(res.Fields) != 1 {
			t.Fatal("Expected 1 field")
		}
		if res.Fields[0].Name != "Name" || res.Fields[0].Label != "Display Name" {
			t.Errorf("Field registration failed")
		}
	})

	t.Run("GetFieldsFor", func(t *testing.T) {
		res := NewResource(MockModel{})
		res.RegisterField("ID", "ID", true)
		res.RegisterField("Name", "Name", false)
		res.SetIndexFields("Name")

		fields := res.GetFieldsFor("index")
		if len(fields) != 1 || fields[0].Name != "Name" {
			t.Errorf("GetFieldsFor 'index' failed")
		}

		fields = res.GetFieldsFor("show") // Should return all if not set
		if len(fields) != 2 {
			t.Errorf("GetFieldsFor 'show' should return all fields, got %d", len(fields))
		}
	})
}

func TestResourceConfiguration(t *testing.T) {
	res := NewResource(MockModel{})
	res.RegisterField("ID", "ID", true)
	res.RegisterField("Name", "Display Name", false)

	t.Run("FieldOptions", func(t *testing.T) {
		res.SetFieldType("Name", "select", "A", "B").
			SetSearchable("Name", "MockModel").
			SetSortable("Name", false)

		field := res.Fields[1]
		if field.Type != "select" {
			t.Fatalf("Expected select field type, got %s", field.Type)
		}
		if len(field.Options) != 2 || field.Options[0] != "A" || field.Options[1] != "B" {
			t.Fatalf("Unexpected field options: %#v", field.Options)
		}
		if !field.Searchable || field.SearchResource != "MockModel" {
			t.Fatalf("Searchable field was not configured: %#v", field)
		}
		if field.Sortable {
			t.Fatal("Expected Name field to be non-sortable")
		}
	})

	t.Run("ViewSpecificFieldsIgnoreUnknownNames", func(t *testing.T) {
		res.SetShowFields("Missing", "Name")
		fields := res.GetFieldsFor("show")
		if len(fields) != 1 || fields[0].Name != "Name" {
			t.Fatalf("Expected only known show field Name, got %#v", fields)
		}
	})

	t.Run("DecoratorsActionsAssociationsAndSidebars", func(t *testing.T) {
		res.SetDecorator("Name", func(val interface{}) template.HTML {
			return template.HTML("decorated-" + val.(string))
		})
		res.AddMemberAction("preview", "Preview", nil)
		res.AddCollectionAction("publish", "Publish", nil)
		res.AddBatchAction("archive", "Archive", nil)
		res.AddScope("active", "Active", nil)
		res.HasMany("Children", "Children", "MockModel", "ParentID")
		res.BelongsTo("ParentID", "Parent", "MockModel", "ID")
		res.AddSidebar("Details", nil)

		if res.Fields[1].Decorator == nil {
			t.Fatal("Expected decorator to be set")
		}
		if len(res.MemberActions) != 1 || res.MemberActions[0].Name != "preview" {
			t.Fatalf("Member action not registered: %#v", res.MemberActions)
		}
		if len(res.CollectionActions) != 1 || res.CollectionActions[0].Name != "publish" {
			t.Fatalf("Collection action not registered: %#v", res.CollectionActions)
		}
		if len(res.BatchActions) != 1 || res.BatchActions[0].Name != "archive" {
			t.Fatalf("Batch action not registered: %#v", res.BatchActions)
		}
		if len(res.Scopes) != 1 || res.Scopes[0].Name != "active" {
			t.Fatalf("Scope not registered: %#v", res.Scopes)
		}
		if len(res.Associations) != 2 {
			t.Fatalf("Expected 2 associations, got %d", len(res.Associations))
		}
		if res.Associations[0].Type != "HasMany" || res.Associations[1].Type != "BelongsTo" {
			t.Fatalf("Unexpected associations: %#v", res.Associations)
		}
		if len(res.Sidebars) != 1 || res.Sidebars[0].Label != "Details" {
			t.Fatalf("Sidebar not registered: %#v", res.Sidebars)
		}
	})
}

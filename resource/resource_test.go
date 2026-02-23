package resource

import (
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

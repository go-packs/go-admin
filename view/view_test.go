package view

import (
	"github.com/go-packs/go-admin/resource"
	"reflect"
	"testing"
)

type TestItem struct {
	ID   uint
	Name string
}

func TestViewHelpers(t *testing.T) {
	res := resource.NewResource(TestItem{})
	fields := []resource.Field{
		{Name: "ID", Label: "ID"},
		{Name: "Name", Label: "Name"},
	}

	t.Run("ItemToMap", func(t *testing.T) {
		item := TestItem{ID: 1, Name: "Test"}
		m := ItemToMap(res, fields, reflect.ValueOf(item))
		if m["ID"] != uint(1) {
			t.Errorf("Expected ID 1, got %v", m["ID"])
		}
		if m["Name"] != "Test" {
			t.Errorf("Expected Name 'Test', got %v", m["Name"])
		}
	})

	t.Run("SliceToMap", func(t *testing.T) {
		items := []TestItem{
			{ID: 1, Name: "A"},
			{ID: 2, Name: "B"},
		}
		maps := SliceToMap(res, fields, reflect.ValueOf(items))
		if len(maps) != 2 {
			t.Errorf("Expected 2 items, got %d", len(maps))
		}
		if maps[0]["Name"] != "A" || maps[1]["Name"] != "B" {
			t.Errorf("SliceToMap failed")
		}
	})
}

package admin

import (
	"testing"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCRUDOperations(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&TestModel{})
	reg := NewRegistry(db)
	reg.Register(TestModel{})

	// 1. Create
	item := &TestModel{Name: "Initial Name"}
	err := reg.Create("TestModel", item)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// 2. List
	list, err := reg.List("TestModel")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	items := list.([]TestModel)
	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}

	// 3. Get
	fetched, err := reg.Get("TestModel", item.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if fetched.(*TestModel).Name != "Initial Name" {
		t.Errorf("Get returned wrong data")
	}

	// 4. Update
	item.Name = "Updated Name"
	err = reg.Update("TestModel", item)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	
	fetched, _ = reg.Get("TestModel", item.ID)
	if fetched.(*TestModel).Name != "Updated Name" {
		t.Errorf("Update data not persisted")
	}

	// 5. Delete
	err = reg.Delete("TestModel", item.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	
	list, _ = reg.List("TestModel")
	if len(list.([]TestModel)) != 0 {
		t.Errorf("Delete failed to remove item")
	}
}

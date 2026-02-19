package admin

import (
	"testing"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestModel struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func TestNewRegistry(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	reg := NewRegistry(db)

	if reg.DB != db {
		t.Errorf("Expected DB to be set correctly")
	}
	if len(reg.Resources) != 0 {
		t.Errorf("Expected initial resources to be 0")
	}
}

func TestRegister(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	reg := NewRegistry(db)

	res := reg.Register(TestModel{})
	if res.Name != "TestModel" {
		t.Errorf("Expected resource name to be TestModel, got %s", res.Name)
	}

	if _, ok := reg.GetResource("TestModel"); !ok {
		t.Errorf("Expected resource to be found in registry")
	}
}

func TestResourceGrouping(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	reg := NewRegistry(db)

	reg.Register(TestModel{}).SetGroup("Testing")
	
	groups := reg.getGroupedResources()
	if len(groups["Testing"]) != 1 {
		t.Errorf("Expected 1 resource in Testing group")
	}
}

func TestResourceNames(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	reg := NewRegistry(db)

	reg.Register(TestModel{})
	names := reg.ResourceNames()

	if len(names) != 1 || names[0] != "TestModel" {
		t.Errorf("Expected ResourceNames to return ['TestModel']")
	}
}

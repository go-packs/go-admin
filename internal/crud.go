// Package internal contains internal helpers for auth, CRUD and auditing.
package internal

import (
	"reflect"

	"github.com/go-packs/go-admin"
)

// List returns all records for a registered resource.
func List(reg *admin.Registry, resourceName string) (interface{}, error) {
	res, ok := reg.GetResource(resourceName)
	if !ok {
		return nil, nil
	}
	modelType := reflect.TypeOf(res.Model)
	destSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0)
	dest := reflect.New(destSlice.Type())
	err := reg.DB.Find(dest.Interface()).Error
	return dest.Elem().Interface(), err
}

// Create persists a new record for the given model.
func Create(reg *admin.Registry, data interface{}) error {
	return reg.DB.Create(data).Error
}

// Get fetches a single record by ID for the named resource.
func Get(reg *admin.Registry, resourceName string, id interface{}) (interface{}, error) {
	res, ok := reg.GetResource(resourceName)
	if !ok {
		return nil, nil
	}
	model := reflect.New(reflect.TypeOf(res.Model)).Interface()
	err := reg.DB.First(model, id).Error
	return model, err
}

// Update saves changes to an existing record.
func Update(reg *admin.Registry, data interface{}) error {
	return reg.DB.Save(data).Error
}

// Delete removes a record by ID for the named resource.
func Delete(reg *admin.Registry, resourceName string, id interface{}) error {
	res, ok := reg.GetResource(resourceName)
	if !ok {
		return nil
	}
	model := reflect.New(reflect.TypeOf(res.Model)).Interface()
	return reg.DB.Delete(model, id).Error
}

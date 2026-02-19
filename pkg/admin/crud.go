package admin

import (
	"reflect"
)

// List fetches all records for a resource.
func (reg *Registry) List(resourceName string) (interface{}, error) {
	res, ok := reg.GetResource(resourceName)
	if !ok {
		return nil, nil
	}

	// Create a slice of the model type using reflection: []Model
	modelType := reflect.TypeOf(res.Model)
	destSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0)
	dest := reflect.New(destSlice.Type())
	dest.Elem().Set(destSlice)

	// Query the database
	err := reg.DB.Find(dest.Interface()).Error
	return dest.Elem().Interface(), err
}

// Create inserts a new record into the database.
func (reg *Registry) Create(resourceName string, data interface{}) error {
	return reg.DB.Create(data).Error
}

// Get fetches a single record by ID.
func (reg *Registry) Get(resourceName string, id interface{}) (interface{}, error) {
	res, ok := reg.GetResource(resourceName)
	if !ok {
		return nil, nil
	}

	model := reflect.New(reflect.TypeOf(res.Model)).Interface()
	err := reg.DB.First(model, id).Error
	return model, err
}

// Update modifies an existing record.
func (reg *Registry) Update(resourceName string, data interface{}) error {
	return reg.DB.Save(data).Error
}

// Delete removes a record from the database.
func (reg *Registry) Delete(resourceName string, id interface{}) error {
	res, ok := reg.GetResource(resourceName)
	if !ok {
		return nil
	}

	model := reflect.New(reflect.TypeOf(res.Model)).Interface()
	return reg.DB.Delete(model, id).Error
}

package admin

import (
	"reflect"
)

func (reg *Registry) List(resourceName string) (interface{}, error) {
	res, ok := reg.GetResource(resourceName)
	if !ok { return nil, nil }
	modelType := reflect.TypeOf(res.Model)
	destSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0); dest := reflect.New(destSlice.Type())
	err := reg.DB.Find(dest.Interface()).Error
	return dest.Elem().Interface(), err
}

func (reg *Registry) Create(resourceName string, data interface{}) error {
	return reg.DB.Create(data).Error
}

func (reg *Registry) Get(resourceName string, id interface{}) (interface{}, error) {
	res, ok := reg.GetResource(resourceName)
	if !ok { return nil, nil }
	model := reflect.New(reflect.TypeOf(res.Model)).Interface()
	err := reg.DB.First(model, id).Error
	return model, err
}

func (reg *Registry) Update(resourceName string, data interface{}) error {
	return reg.DB.Save(data).Error
}

func (reg *Registry) Delete(resourceName string, id interface{}) error {
	res, ok := reg.GetResource(resourceName)
	if !ok { return nil }
	model := reflect.New(reflect.TypeOf(res.Model)).Interface()
	return reg.DB.Delete(model, id).Error
}

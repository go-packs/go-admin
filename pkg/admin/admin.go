package admin

import (
	"fmt"
	"gorm.io/gorm"
)

// Registry keeps track of all registered admin resources.
type Registry struct {
	DB        *gorm.DB
	Resources map[string]*Resource
	Charts    []Chart
	Config    *Config
}

// Chart represents a visual widget on the dashboard.
type Chart struct {
	Label string
	Type  string // bar, line, pie
	Data  func(db *gorm.DB) (labels []string, values []float64)
}

// NewRegistry initializes a new admin registry with a database connection.
func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		DB:        db,
		Resources: make(map[string]*Resource),
		Charts:    []Chart{},
		Config:    DefaultConfig(),
	}
}

// SetConfig updates the registry configuration.
func (reg *Registry) SetConfig(config *Config) {
	reg.Config = config
}

// AddChart adds a chart widget to the dashboard.
func (reg *Registry) AddChart(label, chartType string, provider func(db *gorm.DB) ([]string, []float64)) {
	reg.Charts = append(reg.Charts, Chart{Label: label, Type: chartType, Data: provider})
}

// Register adds a model to the admin interface.
func (reg *Registry) Register(model interface{}) *Resource {
	resource := NewResource(model)
	reg.Resources[resource.Name] = resource
	fmt.Printf("Registered resource: %s\n", resource.Name)
	return resource
}

// GetResource returns a registered resource by name.
func (reg *Registry) GetResource(name string) (*Resource, bool) {
	res, ok := reg.Resources[name]
	return res, ok
}

// ResourceNames returns a list of all registered resource names.
func (reg *Registry) ResourceNames() []string {
	names := make([]string, 0, len(reg.Resources))
	for name := range reg.Resources {
		names = append(names, name)
	}
	return names
}

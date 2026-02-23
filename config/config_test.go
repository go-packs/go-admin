package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		c := DefaultConfig()
		if c.SiteTitle != "Go Admin" {
			t.Errorf("Expected 'Go Admin', got %s", c.SiteTitle)
		}
		if c.DefaultPerPage != 10 {
			t.Errorf("Expected 10, got %d", c.DefaultPerPage)
		}
	})

	t.Run("LoadConfig", func(t *testing.T) {
		yaml := `site_title: 'Test Site'
default_per_page: 10
`
		err := os.WriteFile("test_config.yml", []byte(yaml), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove("test_config.yml")

		c, err := LoadConfig("test_config.yml")
		if err != nil {
			t.Fatal(err)
		}
		if c.SiteTitle != "Test Site" {
			t.Errorf("Expected 'Test Site', got %s", c.SiteTitle)
		}
		if c.DefaultPerPage != 10 {
			t.Errorf("Expected 10, got %d", c.DefaultPerPage)
		}
	})
}

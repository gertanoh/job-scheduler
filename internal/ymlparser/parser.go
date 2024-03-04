package ymlparser

import (
	"errors"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// Job represents a scheduled job.
type Job struct {
	Name     string `yaml:"name"`
	Schedule string `yaml:"schedule"`
	RunOnce  bool   `yaml:"run_once"`
	Steps    []struct {
		Name string `yaml:"name"`
		Run  string `yaml:"run"`
	} `yaml:"steps"`
}

// ParseYAMLFile parses a YAML file and returns a slice of Job structs.
func ParseYAML(yamlData []byte) ([]Job, error) {

	var jobs []Job

	if err := yaml.Unmarshal([]byte(yamlData), &jobs); err != nil {
		return nil, err
	}

	// Validate each job for required fields using reflection
	for _, job := range jobs {
		if err := validateRequiredFields(job); err != nil {
			return nil, err
		}
	}
	return jobs, nil
}

// validateRequiredFields checks if required fields in the struct are present.
func validateRequiredFields(obj interface{}) error {
	val := reflect.ValueOf(obj)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("yaml")
		if strings.Contains(tag, "required") {
			value := val.Field(i).Interface()
			if reflect.DeepEqual(value, reflect.Zero(field.Type).Interface()) {
				return errors.New("required field missing: " + field.Name)
			}
		}
	}

	return nil
}

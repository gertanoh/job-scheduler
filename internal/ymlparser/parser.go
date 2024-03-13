package ymlparser

import (
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

	// Define a struct to match the structure of the YAML data
	var yamlStruct struct {
		Jobs []Job `yaml:"jobs"`
	}

	// Unmarshal the YAML data into the temporary struct
	if err := yaml.Unmarshal(yamlData, &yamlStruct); err != nil {
		return nil, err
	}

	return yamlStruct.Jobs, nil
}

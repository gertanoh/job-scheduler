package ymlparser

import (
	"gopkg.in/yaml.v3"
)

// Job represents a scheduled job.
type Step struct {
	Name string `json:"name" yaml:"name"`
	Run  string `json:"run" yaml:"run"`
}

type Job struct {
	Name     string `json:"name" yaml:"name"`
	Schedule string `json:"schedule" yaml:"schedule"`
	RunOnce  bool   `json:"run_once" yaml:"run_once"`
	Steps    []Step `json:"steps" yaml:"steps"`
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

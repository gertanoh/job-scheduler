package ymlparser_test

import (
	"testing"

	. "gertanoh.job-scheduler/internal/ymlparser"
)

func TestParseYAML(t *testing.T) {
	tests := []struct {
		name     string
		yamlData []byte
		expected []Job
		wantErr  bool
	}{
		{
			name: "Valid YAML",
			yamlData: []byte(`
- name: Every30SecondsJob
  schedule: "*/30 * * * * *"
  run_once: false
  steps:
    - name: YourStep
      run: your_command_here
`),
			expected: []Job{
				{
					Name:     "Every30SecondsJob",
					Schedule: "*/30 * * * * *",
					RunOnce:  false,
					Steps: []struct {
						Name string `yaml:"name"`
						Run  string `yaml:"run"`
					}{
						{
							Name: "YourStep",
							Run:  "your_command_here",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "Invalid YAML",
			yamlData: []byte("invalid: -yaml: data"),
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseYAML(tt.yamlData)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.expected) {
					t.Errorf("ParseYAML() got = %v, want %v", got, tt.expected)
					return
				}

				for i := range got {
					if got[i].Name != tt.expected[i].Name ||
						got[i].Schedule != tt.expected[i].Schedule ||
						got[i].RunOnce != tt.expected[i].RunOnce ||
						len(got[i].Steps) != len(tt.expected[i].Steps) {
						t.Errorf("ParseYAML() got = %v, want %v", got, tt.expected)
						return
					}
					for j := range got[i].Steps {
						if got[i].Steps[j].Name != tt.expected[i].Steps[j].Name ||
							got[i].Steps[j].Run != tt.expected[i].Steps[j].Run {
							t.Errorf("ParseYAML() got = %v, want %v", got, tt.expected)
							return
						}
					}
				}
			}
		})
	}
}

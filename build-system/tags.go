package builder

import (
	"github.com/griffincancode/polyglot.js/core"
)

// GenerateTags generates build tags based on enabled runtimes
func GenerateTags(config *core.Config) []string {
	tags := []string{}

	for name, rtConfig := range config.Languages {
		if rtConfig.Enabled {
			tags = append(tags, "runtime_"+name)
		}
	}

	return tags
}

// IsRuntimeEnabled checks if a runtime build tag is set
func IsRuntimeEnabled(tags []string, runtime string) bool {
	for _, tag := range tags {
		if tag == "runtime_"+runtime {
			return true
		}
	}
	return false
}

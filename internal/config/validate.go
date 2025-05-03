package config

import "fmt"

func validateFilterConfig(c *FilterConfig) error {
	const configSectionIsRequiredFmt = `config section "%s" is required, but was not set`

	if c.Paths == nil {
		return fmt.Errorf(configSectionIsRequiredFmt, "paths")
	}
	return nil
}

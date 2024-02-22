package cli

import "slices"

type ArgDefinition struct {
	Name        string
	Description string
	Required    bool
	ValuesEnum  []string
}

func (d *ArgDefinition) HasValuesEnum() bool {
	return len(d.ValuesEnum) > 0
}

func (d *ArgDefinition) ValueInputs(val string) bool {
	return slices.Contains(d.ValuesEnum, val)
}

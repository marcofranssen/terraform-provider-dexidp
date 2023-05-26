package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

func ListStringValuesToSlice(values types.List) []string {
	elements := make([]string, len(values.Elements()))
	for i, e := range values.Elements() {
		if s, ok := e.(types.String); ok {
			elements[i] = s.ValueString()
		} else {
			elements[i] = e.String()
		}
	}

	return elements
}

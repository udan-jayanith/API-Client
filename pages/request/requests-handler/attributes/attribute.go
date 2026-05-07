package attr

type Attribute struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type AttrCheck struct {
	Checked bool
	Key     string `json:"Key"`
	Value   string `json:"Value"`
}

// MergeAttrList merges a into b
func MergeAttrList(a []Attribute, b []Attribute) []Attribute {
	query_mapped := make(map[string]string, len(a))
	for _, attr := range a {
		query_mapped[attr.Key] = attr.Value
	}

	for i, attr := range b {
		val, ok := query_mapped[attr.Key]
		if ok {
			attr.Value = val
			b[i] = attr
		}
	}

	return b
}

type attr_check struct {
	val     string
	checked bool
}

func bidirectional_attr_check_merge(a []AttrCheck, b []AttrCheck) []AttrCheck {
	merged := make([]AttrCheck, 0, len(a)+len(b))
	mapped := make(map[string]int, len(a)+len(b))

	for _, attr := range a {
		i, ok := mapped[attr.Key]
		if ok {
			merged[i] = attr
		} else {
			merged = append(merged, attr)
			mapped[attr.Key] = len(merged) - 1
		}
	}

	for _, attr := range b {
		i, ok := mapped[attr.Key]
		if ok {
			merged[i] = attr
		} else {
			mapped[attr.Key] = len(merged)
			merged = append(merged, attr)
		}
	}

	return merged
}

// MergeAttrCheckList merges a into b
func MergeAttrCheckList(a []AttrCheck, b []AttrCheck, bidirectional bool) []AttrCheck {
	if bidirectional {
		return bidirectional_attr_check_merge(a, b)
	}

	query_mapped := make(map[string]attr_check, len(a))
	for _, attr := range a {
		query_mapped[attr.Key] = attr_check{
			val:     attr.Value,
			checked: attr.Checked,
		}
	}

	for i, attr := range b {
		attr_check, ok := query_mapped[attr.Key]
		if ok {
			attr.Value = attr_check.val
			attr.Checked = attr_check.checked
			b[i] = attr
		}
	}

	return b
}

package v1alpha1

type DynamicNetwork struct {
	NetworkFetcher    NetworkFetcher     `json:"fetchNetwork"`
	SubnetAllocations []SubnetAllocation `json:"subnetAllocations"`
	ResourceMutations []ResourceMutation `json:"resourceMutations"`
}

// ------------------------------------ NETWORK FETCHER ------------------------------------ //

type NetworkFetcher struct {
	Type       string `json:"type"` // only webhook is supported.
	WebhookURL struct {
		URL string `json:"url"`
	} `json:"webhook"`
	PostTransformations []PostTransformation `json:"postTransformations"`
}

type PostTransformation struct {
	Type     string             `json:"type"` // one of ["jsonpath", "format"]
	Format   string             `json:"format"`
	JSONPath *jsonpath.JSONPath `json:"jsonpath"`
}

var errUnmarshallingPostTransformation = errors.New("unmarshalling postTransformation struct")

func (t *PostTransformation) UnmarshalJSON(data []byte) error {
	var v struct {
		Type     string `json:"type"`
		Format   string `json:"format"`
		JSONPath string `json:"jsonpath"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return errors.Join(err, errUnmarshallingPostTransformation)
	}

	switch v.Type {
	default:
		return errors.Join(fmt.Errorf("type %q not supported", v.Type))
	case "format":
		t.Format = v.Format
	case "jsonpath":
		jp := jsonpath.New("")
		if err := jp.Parse(v.JSONPath); err != nil {
			return errors.Join(err, errUnmarshallingPostTransformation)
		}

		t.JSONPath = jp
	}

	t.Type = v.Type

	return nil
}

// ------------------------------------ SUBNET ALLOCATOR ------------------------------------ //

type SubnetAllocation struct {
	Name    string `json:"name"`
	Netmask string `json:"netmask"`
}

// ------------------------------------ RESOURCE MUTATOR ------------------------------------ //

type ResourceMutation struct {
	SubjectQuery  SubjectQuery       `json:"subjectQuery"`
	FieldSelector *jsonpath.JSONPath `json:"fieldSelector"`
	SubnetName    string             `json:"subnetName"`
}

type SubjectQuery struct {
	Group    string            `json:"group"`
	Version  string            `json:"version"`
	Kind     string            `json:"kind"`
	Selector map[string]string `json:"selector"`
}

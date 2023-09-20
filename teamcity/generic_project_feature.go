package teamcity

import (
	"encoding/json"
)

type GenericProjectFeature struct {
	id          string
	featureType string
	projectID   string
	disabled    bool
	properties  *Properties
}

func (pf *GenericProjectFeature) ID() string {
	return pf.id
}

func (pf *GenericProjectFeature) SetID(value string) {
	pf.id = value
}

func (pf *GenericProjectFeature) Type() string {
	return pf.featureType
}

func (pf *GenericProjectFeature) Properties() *Properties {
	return pf.properties
}

func (pf *GenericProjectFeature) ProjectID() string {
	return pf.projectID
}

func (pf *GenericProjectFeature) SetProjectID(value string) {
	pf.projectID = value
}

func (pf *GenericProjectFeature) Disabled() bool {
	return pf.disabled
}

func (pf *GenericProjectFeature) SetDisabled(value bool) {
	pf.disabled = value
}

func (pf *GenericProjectFeature) MarshalJSON() ([]byte, error) {
	out := &projectFeatureJSON{
		ID:         pf.id,
		Disabled:   NewBool(pf.disabled),
		Properties: pf.properties,
		Inherited:  NewFalse(),
		Type:       pf.Type(),
	}

	return json.Marshal(out)
}

func (pf *GenericProjectFeature) UnmarshalJSON(data []byte) error {
	var aux projectFeatureJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	pf.id = aux.ID
	pf.featureType = aux.Type

	disabled := aux.Disabled
	if disabled == nil {
		disabled = NewFalse()
	}
	pf.disabled = *disabled

	if aux.Properties != nil {
		pf.properties = NewProperties(aux.Properties.Items...)
	}

	return nil
}

func NewGenericProjectFeature(featureType string, propertiesRaw map[string]interface{}) (*GenericProjectFeature, error) {
	properties := NewPropertiesEmpty()
	for name, value := range propertiesRaw {
		value := value.(string)
		properties.Add(&Property{
			Name:  name,
			Value: value,
		})
	}

	return &GenericProjectFeature{
		featureType: featureType,
		properties:  properties,
	}, nil
}

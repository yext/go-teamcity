package teamcity

type BuildTriggerCustomization struct {
	EnforceCleanCheckout                bool        `json:"enforceCleanCheckout"`
	EnforceCleanCheckoutForDependencies bool        `json:"enforceCleanCheckoutForDependencies"`
	Parameters                          *Parameters `json:"parameters,omitempty"`
}

func NewBuildTriggerCustomization() *BuildTriggerCustomization {
	return &BuildTriggerCustomization{}
}

package profiles

var includeWayTags = []string{
	"highway",
}


var excludeWayTagVals = map[string][]string{
	"highway": {"motorway", "trunk", "escape", "raceway"},
	"access": {"private"},
}

// Client for github API
type OffroadProfile struct {}

// NewProfile instantiate a Profile
func NewOffroadProfile() *OffroadProfile {
	return &OffroadProfile{}
}

func (p *OffroadProfile) GetIncludedWayTags() []string {
	return includeWayTags
}

func (p *OffroadProfile) GetExcludedWayTagVals() map[string][]string {
	return excludeWayTagVals
}


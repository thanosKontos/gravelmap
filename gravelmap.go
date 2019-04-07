package gravelmap

type Profile interface {
	GetIncludedWayTags() []string
	GetExcludedWayTagVals() map[string][]string
}

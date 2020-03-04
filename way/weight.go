package way

const (
	// vehicleAcceptanceExclusively defines a way specifically made for the vehicle (e.g. cycleway for bicycles)
	vehicleAcceptanceExclusively = iota

	// vehicleAcceptanceNo defines a way that can be used by vehicle (a small city road for bicycles)
	vehicleAcceptanceYes

	// vehicleAcceptanceNo defines a way that can be used by vehicle but it not recommended (a larger road for bicycles)
	vehicleAcceptancePartially

	// vehicleAcceptanceNo defines a way that cannot be used by vehicle (e.g. footway for bicycles with no bike designation tags)
	vehicleAcceptanceMaybe

	// vehicleAcceptanceNo defines a way that cannot be used vehicle (e.g. path for SUVs)
	vehicleAcceptanceNo
)

const (
	// wayAcceptanceYes defines a way that is allowed to follow in a specific direction (e.g. a 2 way road)
	wayAcceptanceYes = iota

	// wayAcceptanceNo defines a way that is not allowed to follow in a specific direction (e.g. a direction off a one way road)
	wayAcceptanceNo
)

type wayAcceptance struct {
	normal  int32
	reverse int32
}

package commands

type GroundcrewCommand struct {
	Drain DrainCommand `command:"drain"`
}

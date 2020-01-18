package event

// Loaded represents the event for the wallet finishing its loading sequence
type Loaded struct {
	WalletsLoadedCount int32
}

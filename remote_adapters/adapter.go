package remote_adapters

// Adapter for remote repository
type Adapter interface {
	CanHandleCommand(command string) bool
	HandleCommand(command string)
}

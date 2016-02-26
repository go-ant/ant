package startup

var startupHooks []func()

// Register register a function to be run at app startup.
func Register(fn func()) {
	startupHooks = append(startupHooks, fn)
}

func Run() {
	for _, hook := range startupHooks {
		hook()
	}
}

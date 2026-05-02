package configurator

func LaunchConfigurator(live bool) {
	if UIStates.Active {
		return
	}
	LaunchUI(live)
}
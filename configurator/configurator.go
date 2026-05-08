package configurator

import "fmt"

func LaunchConfigurator(live bool, ch chan struct{}) {
	if UIStates.Active {
		fmt.Println("ui active, sorry boss")
		return
	}
	LaunchUI(live, ch)
}
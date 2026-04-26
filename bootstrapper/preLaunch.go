package bootstrapper

func RunPreLaunchProcedures() {
	StartProxy()
	go FindAndOpenLog()

}
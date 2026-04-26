package bootstrapper

func RunPreLaunchProcedures() {
	ModifyHostsFile(true)
	go FindAndOpenLog()

}
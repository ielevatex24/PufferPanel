package comms

type StartServer struct {
	Type   string `json:"type"`
	Id     string `json:"id"`
	Server string `json:"server"`
}

func StartServerType() string { return "server start" }
func NewStartServer(server string) StartServer {
	return StartServer{
		Type:   StartServerType(),
		Server: server,
		Id:     NewId(),
	}
}

type RestartServer struct {
	Type   string `json:"type"`
	Id     string `json:"id"`
	Server string `json:"server"`
}

func RestartServerType() string { return "server restart" }
func NewRestartServer(server string) RestartServer {
	return RestartServer{
		Type:   RestartServerType(),
		Server: server,
		Id:     NewId(),
	}
}

type StopServer struct {
	Type   string `json:"type"`
	Id     string `json:"id"`
	Server string `json:"server"`
}

func StopServerType() string { return "server stop" }
func NewStopServer(server string) StopServer {
	return StopServer{
		Type:   StopServerType(),
		Server: server,
		Id:     NewId(),
	}
}

type KillServer struct {
	Type   string `json:"type"`
	Id     string `json:"id"`
	Server string `json:"server"`
}

func KillServerType() string { return "server kill" }
func NewKillServer(server string) KillServer {
	return KillServer{
		Type:   KillServerType(),
		Server: server,
		Id:     NewId(),
	}
}

type InstallServer struct {
	Type   string `json:"type"`
	Id     string `json:"id"`
	Server string `json:"server"`
}

func InstallServerType() string { return "server install" }
func NewInstallServer(server string) InstallServer {
	return InstallServer{
		Type:   InstallServerType(),
		Server: server,
		Id:     NewId(),
	}
}

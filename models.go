package mikrotik

type Secret struct {
	Name          string
	CallerID      string
	Profile       string
	Comment       string
	RemoteAddress string
	Bts           string
	Host          string
}

type AddressList struct {
	ID           string
	Address      string
	Comment      string
	CreationTime string
	List         string
	Status       string
}

type ActiveConnection struct {
	Name     string
	CallerID string
	Address  string
	Comment  string
	Uptime   string
}

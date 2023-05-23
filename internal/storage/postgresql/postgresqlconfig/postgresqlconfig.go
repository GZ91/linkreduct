package postgresqlconfig

type ConfigDB struct {
	Address  string
	User     string
	Password string
	Dbname   string
}

func New(address, user, password, dbname string) *ConfigDB {
	return &ConfigDB{Address: address, User: user, Password: password, Dbname: dbname}
}

func (d ConfigDB) Empty() bool {
	return d.Address == ""
}

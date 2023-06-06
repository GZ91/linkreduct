package postgresqlconfig

type ConfigDB struct {
	StringServer string
}

func New(StringServer string) *ConfigDB {
	return &ConfigDB{StringServer: StringServer}
}

func (d ConfigDB) Empty() bool {
	return d.StringServer == ""
}

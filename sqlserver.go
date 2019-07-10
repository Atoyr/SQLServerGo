package database

type SqlConnection struct {
	Database string
	Instance string
	Server   string
	User     string
	Password string
}

func (sc *SqlConnection) ConnectionString() string {
	var ret = make([]byte, 0, 512)
	ret = append(ret, "server="...)
	ret = append(ret, sc.Server...)
	if sc.Instance != "" {
		ret = append(ret, "\\"...)
		ret = append(ret, sc.Instance...)
	}
	ret = append(ret, ";user id="...)
	ret = append(ret, sc.User...)
	ret = append(ret, ";password="...)
	ret = append(ret, sc.Password...)
	ret = append(ret, ";database="...)
	ret = append(ret, sc.Database...)
	return string(ret)
}

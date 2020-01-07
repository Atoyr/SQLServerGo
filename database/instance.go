package database

func GetVersion(majorVersion string) string {
	ret := "Microsoft SQL Server "
	switch majorVersion {
	case "14":
		return ret + "2017"
	case "13":
		return ret + "2016"
	case "12":
		return ret + "2014"
	case "11":
		return ret + "2012"
	case "10":
		return ret + "2008R2"
	case "9":
		return ret + "2008"
	}
	return ret
}

package wzry

const (
	token = "free"
)

var (
	Zone    = zone{"aqq", "awx", "ios_qq", "ios_wx"}
	API_URL = "https://www.hive-net.cn/funtools/heroPower/getPower?hero=%s&type=%s&token=%s"
)

type (
	zone struct {
		aqq    string
		awx    string
		ios_qq string
		ios_wx string
	}
)

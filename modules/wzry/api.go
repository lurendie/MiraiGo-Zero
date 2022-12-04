package wzry

const (
	token = "free"
)

var (
	Zone = zone{"qq", "wx", "ios_qq", "ios_wx"}
	API_URL  = "https://www.hive-net.cn/heropower/?token=%s&hero=%s&type=%s"
)

type (
	zone struct {
		qq     string
		wx     string
		ios_qq string
		ios_wx string
	}
)

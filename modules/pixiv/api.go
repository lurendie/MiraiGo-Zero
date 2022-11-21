package pixiv

const (
	POST      = "POST"
	GET       = "GET"
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.24"
	Group     = 1
	Private   = 0
)

var config pixivConfig

// 全局配置结构体
type pixivConfig struct {
	Master int64   `json:"master"`
	Groups []int64 `json:"groups"`
	Proxy  string  `json:"proxy"`
}

var (
	//图片
	IllustURL = "https://api.acgmx.com/illusts/detail?illustId=%v"
	//图片URL
	ImgDataURL = "https://api.acgmx.com/illusts/urlLook?url=%v&cache=false"
	//排行
	RankURL = "https://api.acgmx.com/illusts/ranking?mode=day&date=%v"
	//涩图
	SetuURL = "https://api.acgmx.com/public/setu"
	//画师
	UserURL = "https://api.acgmx.com/public/search/users/illusts?id=%v"
	//以图搜图KEY
	API_KEY = "4f8076a50dfd9964a6b6a0f5dd49c44ef730cc76"
	//以图搜图
	searchImgURL = "https://saucenao.com/search.php?api_key=%v&output_type=2&url=%v"
)

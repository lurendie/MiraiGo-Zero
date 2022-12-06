package pixiv

const (

	//以图搜图KEY
	API_KEY = "4f8076a50dfd9964a6b6a0f5dd49c44ef730cc76"
	//涩图
	SetuURL = "https://api.lolicon.app/setu/v2?r18=1"
)

var (
	//图片
	IllustURL = "https://api.acgmx.com/illusts/detail?illustId=%v"
	//图片URL
	ImgDataURL = "https://api.acgmx.com/illusts/urlLook?url=%v&cache=false"
	//排行
	RankURL = "https://api.acgmx.com/illusts/ranking?mode=day&date=%v"
	//画师
	UserURL = "https://api.acgmx.com/public/search/users/illusts?id=%v"
	//以图搜图
	searchImgURL = "https://saucenao.com/search.php?api_key=%v&output_type=2&url=%v"
)

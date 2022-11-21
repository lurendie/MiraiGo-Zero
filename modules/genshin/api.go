package genshin

const ACT_ID = "e202009291139501"
const BH3_ACT_ID = "e202207181446311"
const APP_VERSION = "2.28.1"
const USER_AGENT = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) miHoYoBBS/" + APP_VERSION
const REFERER_URL = "https://webstatic.mihoyo.com/bbs/event/signin-ys/index.html?bbs_auth_required=true&act_id=" + ACT_ID + "&utm_source=bbs&utm_medium=mys&utm_campaign=icon"
const SIGN_URL = "https://api-takumi.mihoyo.com/event/bbs_sign_reward/sign"
const ROLE_URL = "https://api-takumi.mihoyo.com/binding/api/getUserGameRolesByCookie?game_biz="
const YS_ROLE_URL = ROLE_URL + "hk4e_cn"
const BH3_ROLE_URL = ROLE_URL + "bh3_cn"
const AWARD_URL = "https://api-takumi.mihoyo.com/event/bbs_sign_reward/home?act_id=${ACT_ID}"
const BH3_AWARD_URL = "https://api-takumi.mihoyo.com/event/luna/home?act_id=${BH3_ACT_ID}"
const INFO_URL = "https://api-takumi.mihoyo.com/event/bbs_sign_reward/info?region={region}&act_id=${ACT_ID}&uid={uid}"
const USER_INFO = "https://bbs-api.mihoyo.com/user/wapi/getUserFullInfo?gids=3"

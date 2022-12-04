package genshin

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 原神签到测试
func sign(uid string, region string, cookie string) string {
	client := &http.Client{Timeout: 0}
	//uid := "113821428"
	act_id := "e202009291139501"
	App_Ver := "2.28.1"
	paramStr := "{\"act_id\":\"" + act_id + "\",\"region\":\"" + region + "\",\"uid\":\"" + uid + "\"}"
	DS := getDs()
	UUID := getUUID()
	param := bytes.NewBuffer([]byte(paramStr))
	req, _ := http.NewRequest("POST", SIGN_URL, param)
	req.Header.Add("User-agent", USER_AGENT)
	req.Header.Add("Referer", REFERER_URL)
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Accept-Encoding", "deflate, br")
	req.Header.Add("x-rpc-client_type", "5")
	req.Header.Add("x-rpc-app_version", App_Ver)
	req.Header.Add("x-rpc-device_id", UUID)
	req.Header.Add("DS", DS)
	resp, erro := client.Do(req)
	if erro != nil {
		fmt.Println("访问错误!!!")
	}
	body, _ := io.ReadAll(resp.Body)

	res := make(map[string]interface{})
	json.Unmarshal(body, &res)

	fmt.Printf("body: %v\n", res)
	if len(res["data"].(map[string]interface{})["challenge"].(string)) != 0 {
		req.Header.Add("x-rpc-challenge", res["data"].(map[string]interface{})["challenge"].(string))
		resp, erro := client.Do(req)
		if erro != nil {
			fmt.Println("访问错误!!!")
		}
		body, _ := io.ReadAll(resp.Body)

		res := make(map[string]interface{})
		json.Unmarshal(body, &res)
		fmt.Printf("res1: %v\n", res)
		return "签到成功!"
	} else {
		return "账号风控,需验证码验证,签到失败,请进入米游社APP手动启动!"
	}
}

func getDs() string {
	var n string = "ulInCDohgEs557j0VsPDYnQaaz6KJcv5"
	//时间戳
	t := time.Now().Unix()
	//获取6位随机字符串
	r := randomStr()
	md5 := Md5Crypt("salt=" + n + "&t=" + strconv.FormatInt(t, 10) + "&r=" + r)
	return strconv.FormatInt(t, 10) + "," + r + "," + md5
}

// 给字符串生成md5
// @params str 需要加密的字符串
// @params salt interface{} 加密的盐
// @return str 返回md5码
func Md5Crypt(str string, salt ...interface{}) (CryptStr string) {
	if l := len(salt); l > 0 {
		slice := make([]string, l+1)
		str = fmt.Sprintf(str+strings.Join(slice, "%v"), salt...)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func randomStr() string {
	CONSTANTS := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	reStr := ""
	var result *big.Int
	for {
		// 生成 10 个 [0, 128) 范围的真随机数。
		for i := 0; i < 1; i++ {
			result, _ = rand.Int(rand.Reader, big.NewInt(61))
		}
		reStr += CONSTANTS[result.Int64() : result.Int64()+1]
		if len(reStr) == 6 {
			break
		}
	}
	return reStr
}

func getUUID() string {
	id, _ := uuid.NewUUID()
	return id.String()
}

func GetInfo(cookie string) {

}

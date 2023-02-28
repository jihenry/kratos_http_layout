package sign

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"gitlab.yeahka.com/gaas/pkg/util"

	"gitlab.yeahka.com/gaas/pkg/log"

	ji "github.com/json-iterator/go"
)

//签名秘钥分为未登录秘钥和已登录秘钥。
//未登录秘钥：通过请求参数计算获得。公式为：signKey=upperCase(hex(md5(timestamp + "#" + nonce)))
//已登录秘钥：通过请求参数和登录返回的签名盐计算得到，若返回签名盐为空时，计算方式与未登录秘钥计算方式一致。
//签名盐公式为：signKey=upperCase(hex(md5(timestamp + "#" + nonce + "#" + saltKey)))

func CheckSign(data []byte, saltKey string, signSlice ...string) (ans string, timestamp, nonce string) {
	var (
		maps = make(map[string]interface{})
	)
	//参数解析错误
	if err := util.JSON.Unmarshal(data, &maps); err != nil {
		return
	}
	keys := make([]string, 0, len(maps))
	//获取data 中获取参数的key 进行 ascii码 从小到大排序
	if len(maps) <= 0 {
		return
	}
	var (
		slices  = make([]interface{}, 0, len(signSlice))
		hashKey string
	)
	for _, v := range signSlice {
		if value, ok := maps[v]; ok {
			slices = append(slices, value)
			continue
		}
	}

	//用户签名
	if len(slices) > 0 {
		var (
			buf1 strings.Builder
		)
		for _, v := range slices {
			val := ToString(v)
			if buf1.Len() > 0 {
				buf1.WriteByte('#')
			}
			buf1.WriteString(val)
		}
		buf1.WriteByte('#')
		buf1.WriteString(saltKey)
		hashKey = strings.ToUpper(Md5(buf1.String()))
	}
	for k := range maps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf strings.Builder
	for _, key := range keys {
		if value, ok := maps[key]; ok {
			val := ToString(value)
			switch key {
			case "timestamp":
				timestamp = val
			case "nonce":
				nonce = val
			}
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(key)
			buf.WriteByte('=')
			buf.WriteString(val)
		}
	}
	queryString := buf.String()
	if len(hashKey) > 0 {
		queryString += "&key=" + hashKey
	}
	log.Infof("[queryString] queryString:%s", queryString)
	ans = Hmac(queryString, hashKey)
	return
}

func Hmac(data, key string) string {
	hm := hmac.New(md5.New, []byte(key))
	hm.Write([]byte(data))
	return hex.EncodeToString(hm.Sum([]byte("")))
}

func Md5(data string) string {
	md := md5.New()
	md.Write([]byte(data))
	return hex.EncodeToString(md.Sum([]byte("")))
}

func mapToSting(value interface{}) map[string]interface{} {
	var (
		keys []string
	)
	for k := range value.(map[string]interface{}) {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	mm := make(map[string]interface{}, 0)
	for _, kk := range keys {
		if val, ok := value.(map[string]interface{})[kk]; ok {
			switch val.(type) {
			case map[string]interface{}:
				mm[kk] = mapToSting(val)
			default:
				mm[kk] = val
			}
		}
	}
	return mm
}

func ToString(value interface{}) (vv string) {
	switch value.(type) {
	case string:
		vv = fmt.Sprintf("%s", value)
		//重新排序
	case map[string]interface{}:
		mm := mapToSting(value)
		json := ji.Config{
			EscapeHTML:                    false,
			ObjectFieldMustBeSimpleString: true,
			UseNumber:                     true,
			SortMapKeys:                   true,
		}.Froze()
		vvv, _ := json.Marshal(mm)
		vv = string(vvv)
	default:
		vv, _ = util.JSON.MarshalToString(value)
	}
	return
}

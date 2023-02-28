package sign

import (
	"testing"
)

func TestSign(t *testing.T) {
	t.Log(Hmac("ActivityID=c6s6b5tk1e6ovfkj7qlg&nonce=yeahkagamesgxk5qq4jsg&timestamp=1642745591448&key=2F3A6CD6CBC1C37C03401CC3FBE699FA", "2F3A6CD6CBC1C37C03401CC3FBE699FA"))
}

func TestMd5(t *testing.T) {
	t.Log(Hmac("ActivityID=c44an8dk1e6u8up6dhbg&nonce=yeahkagameje8i9zmaak&timestamp=1643112551310&key=65E06177DCC1228EFE8E5BB5E46579EC", "65E06177DCC1228EFE8E5BB5E46579EC"))
	data := "wxa50dcaae208f9ac9" + "c7ijh3lk1e6pu2vcre7g"
	t.Log(data[10:26])
}

func TestInitSign(t *testing.T) {
	//should := require.New(t)
	//ua := `Mozilla/5.0 (iPhone; CPU iPhone OS 15_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x18001234) NetType/WIFI Language/zh_CN`
	for i := 0; i < 100; i++ {
		data := `{"ActivityID":"c3h84ie0desb2thug5dg","Score":0,"StartTime":1656600080,"EndTime":1656600097,"WatermelonNum":0,"OperNum":1,"Status":2,"EndType":"game_end","RunID":"17d1c801-404b-4a93-8f6f-2050e86d78df","Items":{"2001":0,"2002":0,"2003":0},"timestamp":1656600097136,"nonce":"yeahkagamebbcatd0ieql"}`
		//saltkey := "83bd80c2-dfbb-47fc-9d0e-3666243fa40f" + util.SplitUserAgent(ua)
		_, timestamp, nonce := CheckSign([]byte(data), "4465620d9dbf2bb3ad2a102dfcf88340", []string{"timestamp", "nonce"}...)
		t.Log("timestamp", timestamp, "nonce", nonce)
		//should.Equal("573b779c718f288ca5c47a4b75147c76", ans)
	}

	//
}

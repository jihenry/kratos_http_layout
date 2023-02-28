package sign

import (
	"testing"
)

func TestXOREncode(t *testing.T) {
	sign := XOREncode("ecode: abcdefgabcdefgabcdefgabcdefg", "pocket")
	t.Log("sign:", sign)
	sign = "BQwpJTUHVA8GIlUvKhQUOwcKLyc3AFENACRXLQ=="
	data := XORDecode(sign, "pocket")
	t.Log("decode:", data)
	//maps := make(map[string]uint64, 0)
	//for _, item := range strings.Split(data, "&") {
	//	slices := strings.Split(item, ":")
	//	for i := 0; i < len(slices); i += 2 {
	//		id, err := strconv.Atoi(slices[i+1])
	//		if err != nil {
	//			continue
	//		}
	//		maps[slices[i]] = uint64(id)
	//	}
	//}
	//t.Log(maps)
}

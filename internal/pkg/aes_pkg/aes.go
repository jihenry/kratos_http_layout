package aesPkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"hlayout/internal/biz/common"
	"sync"
)

type (
	aesPkg struct {
		key  []byte
		iv   []byte
		data []byte
		code string
	}
	WithOption func(a *aesPkg)
	//drawPlugin 抽奖接口
	aesPlugin interface {
		// AesDeCrypt 实现解密
		AesDeCrypt() (string, error)
		// AesEncrypt 实现加密
		AesEncrypt() (string, error)
	}
)

var (
	keys = map[string]string{
		"wxec31287c4d9b9853": "dnJAPa3ndA1JLsuY", //口袋ar
		"wxa50dcaae208f9ac9": "3e8ec605eafe8cae", //跃动女孩
		"wx4e443d9613e66610": "7104edc47f3ec13e", //消消乐
		"wxe50aaf462a3041b0": "96838c6b5ab6e2fd", //大西瓜
		"wx3cdaa9114157de80": "NzVMPFQiRtJJnhg2", // 神奇板子
		"wxb6557f4ad72b0adf": "23dfas12dayhjkre", // 贪吃蛇
		"h5_00001":           "NzVMPFQiRtJSqbz9",
	}
)

var gameKeys map[int32]string

func Init() error {
	gameKeys = map[int32]string{
		common.GetGameBoardID():     "NzVMPFQiRtJJnhg2", //神奇板子
		common.GetGameSGHHHID():     "96838c6b5ab6e2fd", //大西瓜
		common.GetGameDfgirlID():    "3e8ec605eafe8cae",
		common.GetGameCrushID():     "7104edc47f3ec13e",
		common.GetGameSnakeID():     "23dfas12dayhjkre",
		common.GetGameARID():        "dnJAPa3ndA1JLsuY",
		common.GetGameSGHHHIDH5():   "96838c6b5ab6e2fd", // 大西瓜h5
		common.GetGameMidAutumnID(): "d8212c5626a375cc",
		common.GetGamePocketID():    "KCusI0BtE7OSUz8Y", // 极速跑酷
		common.GetGamePTID():        "Z5Fmj6aqU4xEUNBe", //拼图H5
		common.GetNBCBGameH5ID():    "tkLXQzdINpeDMNQx", //蛋黄小镇H5
		common.GetNBCBGameAppID():   "660ENKtaJteBCFjj", //蛋黄小镇小程序
		common.GetGameBocH5ID():     "idEi3HlUBab5JuaA", //中国银行H5
		common.GetGameARDLL():       "Rqfi2T4S4MREZ6Dt", //AR叠叠乐
	}
	return nil
}

// WithData 引入加密参数
func WithData(data string) WithOption {
	return func(a *aesPkg) {
		a.data = []byte(data)
	}
}

// WithKey 引入加密串
func WithKey(appId string) WithOption {
	return func(a *aesPkg) {
		if k, ok := keys[appId]; ok {
			a.key = []byte(k)
		}
	}
}

// WithKey 引入加密串
func WithGameKey(gid int32) WithOption {
	return func(a *aesPkg) {
		if k, ok := gameKeys[gid]; ok {
			a.key = []byte(k)
		}
	}
}

// WithKeyAfterGameKey 补充检查有没有key
func WithKeyAfterGameKey(appId string) WithOption {
	return func(a *aesPkg) {
		if len(a.key) == 0 {
			if k, ok := keys[appId]; ok {
				a.key = []byte(k)
			}
		}
	}
}

func WithIv(iv string) WithOption {
	return func(a *aesPkg) {
		a.iv = []byte(iv)
	}
}

func WithCode(code string) WithOption {
	return func(a *aesPkg) {
		a.code = code
	}
}

func NewAes(opt ...WithOption) aesPlugin {
	pool := sync.Pool{
		New: func() interface{} {
			return &aesPkg{}
		},
	}
	obj := pool.Get().(*aesPkg)
	for _, op := range opt {
		op(obj)
	}
	pool.Put(obj)
	return obj
}

//pKCS7Padding  填充模式
func (a *aesPkg) pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

var (
	iv = []byte("tdrdadq59tbss5n7")
)

//pKCS7UnPadding 填充的反向操作，删除填充字符串
func (a *aesPkg) pKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		padding := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - padding)], nil
	}
}

func (a *aesPkg) aesEncrypt() ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData := a.pKCS7Padding(a.data, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, a.iv[:blockSize])
	crypt := make([]byte, len(origData))
	//执行加密
	blocMode.CryptBlocks(crypt, origData)
	return crypt, nil
}

func (a *aesPkg) aesDeCrypt() (string, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, a.iv[:blockSize])
	origData := make([]byte, len(a.data))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, a.data)
	//去除填充字符串
	origData, err = a.pKCS7UnPadding(origData)
	if err != nil {
		return "", err
	}
	return string(origData), err
}

func (a *aesPkg) AesDeCrypt() (string, error) {
	//解密base64字符串
	if len(a.code) == 0 {
		return "", errors.New("code is nil")
	}
	pwdByte, err := base64.StdEncoding.DecodeString(string(a.code))
	if err != nil {
		return "", err
	}
	//执行AES解密
	a.data = pwdByte
	return a.aesDeCrypt()
}

func (a *aesPkg) AesEncrypt() (string, error) {
	if len(a.key) == 0 {
		return "", errors.New("key is nil")
	}
	result, err := a.aesEncrypt()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

func GenerateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 26)
	copy(genKey, key)
	for i := 0; i < len(key); {
		for j := 20; j < 26 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

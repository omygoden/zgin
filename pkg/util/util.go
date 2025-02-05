package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/omygoden/gotools/sfrand"
	"golang.org/x/crypto/bcrypt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// ip/域名检测
// count表示ping次数，timeout表示ping超时时间
func PingCheck(domain string, count int, timeout time.Duration) (int, error) {
	pinger, err := ping.NewPinger(domain)
	if err != nil {
		return 0, errors.New("初始化ping对象失败：" + err.Error())
	}
	pinger.Count = count
	pinger.Timeout = time.Second * timeout
	err = pinger.Run()
	if err != nil {
		return 0, errors.New("ping操作失败：" + err.Error())
	}
	stats := pinger.Statistics()

	//接收包大于0表示ping成功,否则表示失败
	return stats.PacketsRecv, nil
}

func GetProjectName() string {
	pwd, _ := os.Getwd()
	reg := regexp.MustCompile("/bin|/test")

	pwd = reg.ReplaceAllString(pwd, "")
	projects := strings.Split(pwd, "/")

	return projects[len(projects)-1]
}

func GetProjectPath() string {
	pwd, _ := os.Getwd()
	reg := regexp.MustCompile("/bin|/test")

	pwd = reg.ReplaceAllString(pwd, "")
	return pwd
}

// 生成订单号
func GenerateOrderNo() string {
	m := time.Now().UnixMicro() - time.Now().Unix()*1000000
	s := "Y" + time.Now().Format("20060102150405") + fmt.Sprintf("%06d", m) + fmt.Sprintf("%d", sfrand.RandRange(1000, 9999))
	return s
}

func HashPwdGenerate(pwd string) string {
	hashPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)

	return string(hashPwd)
}

func HashPwdVerify(hashPwd string, pwd string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd))
	if res != nil {
		return false
	}
	return true
}

func MoneyFormatMul(money float64) int64 {
	return int64(money * 100)
}

func MoneyFormatDiv(money int64) float64 {
	return float64(money) / 100
}

func GetMyFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncrypt 加密函数
func AesEncrypt(plaintext []byte, key []byte) (crypted []byte, err error) {
	defer func() {
		if errs := recover(); errs != nil {
			err = errors.New(errs.(string))
			return
		}
		return
	}()

	c := make([]byte, aes.BlockSize+len(plaintext))
	iv := c[:aes.BlockSize]

	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted = make([]byte, len(plaintext))
	blockMode.CryptBlocks(crypted, plaintext)
	return
}

// AesDecrypt 解密函数
func AesDecrypt(ciphertext []byte, key []byte) (origData []byte, err error) {
	defer func() {
		if errs := recover(); errs != nil {
			err = errors.New(errs.(string))
			return
		}
		return
	}()

	c := make([]byte, aes.BlockSize+len(ciphertext))
	iv := c[:aes.BlockSize]

	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData = make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origData, ciphertext)
	origData = PKCS7UnPadding(origData)
	return
}

// 获取redis幂等有效期
// 5点之前，默认设置一个小时，如果是5点之后则将有效期设置到隔天的1-5点
func GetRedisMdExpire() time.Duration {
	if time.Now().Hour() <= 5 {
		return time.Hour
	} else {
		return time.Hour*time.Duration(23-time.Now().Hour()+int(sfrand.RandRange(1, 5))) + time.Minute*time.Duration(sfrand.RandRange(1, 59))
	}
}

func MapFilter(data map[string]interface{}) {
	var f float64
	for k, v := range data {
		if v == f || v == "" || v == nil || v == "0" {
			delete(data, k)
		}
	}
}

func GetRandStr(num int) string {
	var m = map[string]string{"0": "A", "1": "B", "2": "C", "3": "D", "4": "E", "5": "F", "6": "G", "7": "H", "8": "I", "9": "J", "10": "K", "11": "L", "12": "M", "13": "N", "14": "O", "15": "P", "16": "Q", "17": "R", "18": "S", "19": "T", "20": "U", "21": "V", "22": "W", "23": "X", "24": "Y", "25": "Z", "26": "0", "27": "1", "28": "2", "29": "3", "30": "4", "31": "5", "32": "6", "33": "7", "34": "8", "35": "9", "36": "a", "37": "b", "38": "c", "39": "d", "40": "e", "41": "f", "42": "g", "43": "h", "44": "i", "45": "j", "46": "k", "47": "l", "48": "m", "49": "n", "50": "o", "51": "p", "52": "q", "53": "r", "54": "s", "55": "t", "56": "u", "57": "v", "58": "w", "59": "x", "60": "y", "61": "z"}
	var str string
	for _, v := range m {
		str += v
		num--
		if num == 0 {
			break
		}
	}
	return str
}

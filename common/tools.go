package common

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

type Tools struct {
}

var (
	t     *Tools
	once  sync.Once
	tools *Tools
)

func init() {
	tools = NewTools()
}

/**
 * 返回单例实例
 * @method New
 */
func NewTools() *Tools {
	once.Do(func() { //只执行一次
		t = &Tools{}
	})

	return t
}

/**
 * md5 加密
 * @method MD5
 * @param  {[type]} data string [description]
 */
func (t *Tools) MD5(data string) string {
	m := md5.New()
	io.WriteString(m, data)

	return fmt.Sprintf("%x", m.Sum(nil))
}

/**
 * string转换int
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) ParseInt64(b string, defInt int64) int64 {
	id, err := strconv.ParseInt(b, 10, 64)
	if err != nil {
		return defInt
	} else {
		return id
	}
}

/**
 * 结构体转换成map对象
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) GetDateNowString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

/**
 * 结构体转换成map对象
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) StructToMap(obj interface{}) map[string]interface{} {
	return structs.Map(obj)
}

func (t *Tools) CopyStruct(dst, src interface{}) error {
	jsonStr, err := json.Marshal(dst)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonStr, src); err != nil {
		return err
	}

	return nil

	//dstValue := reflect.ValueOf(dst)
	//if dstValue.Kind() != reflect.Ptr {
	//	err := errors.New("dst isn't a pointer to struct")
	//	return err
	//}
	//
	//dstElem := dstValue.Elem()
	//if dstElem.Kind() != reflect.Struct {
	//	err := errors.New("pointer doesn't point to struct")
	//	return err
	//}
	//
	//srcValue := reflect.ValueOf(src)
	//srcType := reflect.TypeOf(src)
	//if srcType.Kind() != reflect.Struct {
	//	err := errors.New("src isn't struct")
	//	return err
	//}
	//
	//for i := 0; i < srcType.NumField(); i++ {
	//	sf := srcType.Field(i)
	//	sv := srcValue.FieldByName(sf.Name)
	//	if dv := dstElem.FieldByName(sf.Name); dv.IsValid() && dv.CanSet() {
	//		dv.Set(sv)
	//	}
	//}
	//
	//return nil
}

/**
 * 判断手机号码
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) IsMobile(mobile string) bool {

	reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`

	rgx := regexp.MustCompile(reg)

	return rgx.MatchString(mobile)
}

/**
 * 验证密码
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) CheckPassword(password, metaPassword string) bool {

	return strings.EqualFold(password, metaPassword)
}

/**
 * 生成随机字符串
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) GetRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}

/**
 * 生成用户Redis key
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) UserRedisKey(userId int64) string {
	userKey := fmt.Sprintf("user_login_%d", userId)

	return userKey
}

/**
 * 生成用户Token
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) GenerateToken(n int) (string, error) {
	token, err := t.GenerateRandomString(n)
	return token, err
}

func (t *Tools) GenerateRandomString(s int) (string, error) {
	b, err := t.GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func (t *Tools) GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *Tools) GenerateFileName(path, fileName string) string {
	now := time.Now().UnixNano()
	random := t.MD5(t.GetRandomString(12))

	number := len(random)

	path = fmt.Sprintf("%s/%s/%s",
		strings.Trim(path, "/"),
		string([]byte(random)[:6]),
		string([]byte(random)[number-6:number]))

	t.CreatedDir(path, os.ModePerm)

	return fmt.Sprintf("%s/%d_%s%s",
		path,
		now,
		string([]byte(random)[10:20]),
		filepath.Ext(fileName))
}

func (t *Tools) CreatedDir(dir string, mode os.FileMode) {
	ok, err := t.PathExists(dir)
	if err == nil && !ok {
		os.MkdirAll(dir, mode)
	}
}

func (t *Tools) SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)

	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}

		data = append(data, d)
	}

	return strings.ToLower(string(data[:]))
}

func (t *Tools) GetNowMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (t *Tools) GetAddMillisecond(add time.Duration) int64 {
	return time.Now().Add(add).UnixNano() / int64(time.Millisecond)
}

func (t *Tools) GetAcceptLanguage(acceptLanguage string) string {
	language := "zh-CN"

	lang := strings.Split(acceptLanguage, ";")
	if len(lang) >= 1 {
		langs := strings.Split(lang[0], ",")
		language = langs[0]
	}

	return language
}

func (t *Tools) PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func (t *Tools) GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max - min)
	randNum = randNum + min

	return randNum
}

func (t *Tools) CreateUUID() string {
	return uuid.NewV4().String()
}

func (t *Tools) SetOffset(page, size int64) int64 {
	offset := (page - 1) * size

	return offset
}

func (t *Tools) Merge(data []string, dist ...[]string) []string {
	for _, list := range dist {
		for _, v := range list {
			data = append(data, v)
		}
	}

	return data
}

func (t *Tools) GetStatus(err error) string {
	if err != nil {
		return "errors"
	}

	return "success"
}

func (t *Tools) LoopCall(structs ...interface{}) {
	for _, v := range structs {
		classType := reflect.TypeOf(v)
		classValue := reflect.ValueOf(v)

		for i := 0; i < classType.NumMethod(); i++ {
			m := classValue.MethodByName(classType.Method(i).Name)
			if m.IsValid() {
				var params []reflect.Value
				m.Call(params)
			}
		}
	}
}

func (t *Tools) UniqueInt(elements []int) []int {
	encountered := map[int]bool{}
	var result []int

	for v := range elements {
		if !encountered[elements[v]] == true {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}

func (t *Tools) UniqueString(elements []string) []string {
	encountered := map[string]bool{}
	var result []string

	for v := range elements {
		if !encountered[elements[v]] == true {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}

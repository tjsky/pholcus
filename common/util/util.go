package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"io"
	r "math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/net/html/charset"

	"github.com/andeya/gust/result"
	"github.com/andeya/pholcus/logs"
)

const (
	// USE_KEYIN is the initial value for enabling Keyin in Spider.
	USE_KEYIN = "\r\t\n"
)

var (
	re             = regexp.MustCompile(">[ \t\n\v\f\r]+<")
	jsonpKeyRegexp = regexp.MustCompile(`([^\s\:\{\,\d"]+|[a-z][a-z\d]*)\s*\:`)
	isNumRegexp    = regexp.MustCompile(`^\d+$`)
)

// JsonpToJson modify jsonp string to json string
// Example: forbar({a:"1",b:2}) to {"a":"1","b":2}
func JsonpToJson(json string) string {
	start := strings.Index(json, "{")
	end := strings.LastIndex(json, "}")
	start1 := strings.Index(json, "[")
	if start1 > 0 && start > start1 {
		start = start1
		end = strings.LastIndex(json, "]")
	}
	if end > start && end != -1 && start != -1 {
		json = json[start : end+1]
	}
	json = strings.ReplaceAll(json, "\\'", "")
	return jsonpKeyRegexp.ReplaceAllString(json, "\"$1\":")
}

// Mkdir creates the directory for the given path.
func Mkdir(Path string) {
	p, _ := path.Split(Path)
	if p == "" {
		return
	}
	d, err := os.Stat(p)
	if err != nil || !d.IsDir() {
		if err = os.MkdirAll(p, 0777); err != nil {
			logs.Log.Error("failed to create path [%v]: %v\n", Path, err)
		}
	}
}

// The GetWDPath gets the work directory path.
func GetWDPath() string {
	wd := os.Getenv("GOPATH")
	if wd == "" {
		panic("GOPATH is not setted in env.")
	}
	return wd
}

// The IsDirExists judges path is directory or not.
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	}
	return fi.IsDir()
}

// The IsFileExists judges path is file or not.
func IsFileExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	}
	return !fi.IsDir()
}

// walkPath resolves targpath to an absolute path. Internal helper using gust.Result.
func walkPath(targpath string) result.Result[string] {
	if filepath.IsAbs(targpath) {
		return result.Ok(targpath)
	}
	abs, err := filepath.Abs(targpath)
	if err != nil {
		return result.TryErr[string](err)
	}
	return result.Ok(abs)
}

// WalkFiles walks files under targpath, optionally filtered by suffixes.
func WalkFiles(targpath string, suffixes ...string) (filelist []string) {
	r := walkPath(targpath)
	if r.IsErr() {
		logs.Log.Error("util.WalkFiles: %v\n", r.UnwrapErr())
		return
	}
	targpath = r.Unwrap()
	err := filepath.Walk(targpath, func(retpath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if len(suffixes) == 0 {
			filelist = append(filelist, retpath)
			return nil
		}
		for _, suffix := range suffixes {
			if strings.HasSuffix(retpath, suffix) {
				filelist = append(filelist, retpath)
			}
		}
		return nil
	})

	if err != nil {
		logs.Log.Error("util.WalkFiles: %v\n", err)
		return
	}

	return
}

// WalkDir walks directories under targpath, optionally filtered by suffixes.
func WalkDir(targpath string, suffixes ...string) (dirlist []string) {
	r := walkPath(targpath)
	if r.IsErr() {
		logs.Log.Error("util.WalkDir: %v\n", r.UnwrapErr())
		return
	}
	targpath = r.Unwrap()
	err := filepath.Walk(targpath, func(retpath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			return nil
		}
		if len(suffixes) == 0 {
			dirlist = append(dirlist, retpath)
			return nil
		}
		for _, suffix := range suffixes {
			if strings.HasSuffix(retpath, suffix) {
				dirlist = append(dirlist, retpath)
			}
		}
		return nil
	})

	if err != nil {
		logs.Log.Error("util.WalkDir: %v\n", err)
		return
	}

	return
}

// WalkRelFiles walks files under targpath and returns relative paths, optionally filtered by suffixes.
func WalkRelFiles(targpath string, suffixes ...string) (filelist []string) {
	r := walkPath(targpath)
	if r.IsErr() {
		logs.Log.Error("util.WalkRelFiles: %v\n", r.UnwrapErr())
		return
	}
	targpath = r.Unwrap()
	err := filepath.Walk(targpath, func(retpath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if len(suffixes) == 0 {
			filelist = append(filelist, RelPath(retpath))
			return nil
		}
		_retpath := RelPath(retpath)
		for _, suffix := range suffixes {
			if strings.HasSuffix(_retpath, suffix) {
				filelist = append(filelist, _retpath)
			}
		}
		return nil
	})

	if err != nil {
		logs.Log.Error("util.WalkRelFiles: %v\n", err)
		return
	}

	return
}

// WalkRelDir walks directories under targpath and returns relative paths, optionally filtered by suffixes.
func WalkRelDir(targpath string, suffixes ...string) (dirlist []string) {
	r := walkPath(targpath)
	if r.IsErr() {
		logs.Log.Error("util.WalkRelDir: %v\n", r.UnwrapErr())
		return
	}
	targpath = r.Unwrap()
	err := filepath.Walk(targpath, func(retpath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			return nil
		}
		if len(suffixes) == 0 {
			dirlist = append(dirlist, RelPath(retpath))
			return nil
		}
		_retpath := RelPath(retpath)
		for _, suffix := range suffixes {
			if strings.HasSuffix(_retpath, suffix) {
				dirlist = append(dirlist, _retpath)
			}
		}
		return nil
	})

	if err != nil {
		logs.Log.Error("util.WalkRelDir: %v\n", err)
		return
	}

	return
}

// RelPath converts targpath to a path relative to the current working directory.
func RelPath(targpath string) string {
	basepath, err := filepath.Abs("./")
	if err != nil {
		logs.Log.Error("util.RelPath: filepath.Abs: %v\n", err)
		return targpath
	}
	rel, err := filepath.Rel(basepath, targpath)
	if err != nil {
		logs.Log.Error("util.RelPath: filepath.Rel(%q, %q): %v\n", basepath, targpath, err)
		return targpath
	}
	return strings.ReplaceAll(rel, `\`, `/`)
}

// The IsNum judges string is number or not.
func IsNum(a string) bool {
	return isNumRegexp.MatchString(a)
}

// XML2mapstr converts simple XML to a string map (supports UTF-8).
func XML2mapstr(xmldoc string) map[string]string {
	var t xml.Token
	var err error
	inputReader := strings.NewReader(xmldoc)
	decoder := xml.NewDecoder(inputReader)
	decoder.CharsetReader = func(s string, r io.Reader) (io.Reader, error) {
		return charset.NewReader(r, s)
	}
	m := make(map[string]string, 32)
	key := ""
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			key = token.Name.Local
		case xml.CharData:
			content := Bytes2String([]byte(token))
			m[key] = content
		default:
		}
	}

	return m
}

// MakeHash converts a string to a CRC32 hash hex string.
func MakeHash(s string) string {
	const IEEE = 0xedb88320
	var IEEETable = crc32.MakeTable(IEEE)
	hash := fmt.Sprintf("%x", crc32.Checksum([]byte(s), IEEETable))
	return hash
}

func HashString(encode string) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(encode))
	return hash.Sum64()
}

// MakeUnique creates a unique fingerprint for obj (method 1: FNV-64).
func MakeUnique(obj interface{}) string {
	b, _ := json.Marshal(obj)
	hash := fnv.New64()
	hash.Write(b)
	return strconv.FormatUint(hash.Sum64(), 10)
}

// MakeMd5 creates an MD5 fingerprint for obj (method 2).
func MakeMd5(obj interface{}, length int) string {
	if length > 32 {
		length = 32
	}
	h := md5.New()
	baseString, _ := json.Marshal(obj)
	h.Write([]byte(baseString))
	s := hex.EncodeToString(h.Sum(nil))
	return s[:length]
}

// JsonString converts obj to a JSON string.
func JsonString(obj interface{}) string {
	b, _ := json.Marshal(obj)
	s := fmt.Sprintf("%+v", Bytes2String(b))
	r := strings.ReplaceAll(s, `\u003c`, "<")
	r = strings.ReplaceAll(r, `\u003e`, ">")
	return r
}

// CheckErr checks and logs the error if non-nil.
func CheckErr(err error) {
	if err != nil {
		logs.Log.Error("%v", err)
	}
}
func CheckErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// FileNameReplace replaces invalid filename characters with similar alternatives.
func FileNameReplace(fileName string) string {
	var q = 1
	r := []rune(fileName)
	size := len(r)
	for i := 0; i < size; i++ {
		switch r[i] {
		case '"':
			if q%2 == 1 {
				r[i] = '“'
			} else {
				r[i] = '”'
			}
			q++
		case ':':
			r[i] = '：'
		case '*':
			r[i] = '×'
		case '<':
			r[i] = '＜'
		case '>':
			r[i] = '＞'
		case '?':
			r[i] = '？'
		case '/':
			r[i] = '／'
		case '|':
			r[i] = '∣'
		case '\\':
			r[i] = '╲'
		}
	}
	return strings.ReplaceAll(string(r), USE_KEYIN, ``)
}

// ExcelSheetNameReplace replaces invalid Excel sheet name characters with underscores.
func ExcelSheetNameReplace(fileName string) string {
	r := []rune(fileName)
	size := len(r)
	for i := 0; i < size; i++ {
		switch r[i] {
		case ':', '：', '*', '?', '？', '/', '／', '\\', '╲', ']', '[':
			r[i] = '_'
		}
	}
	return strings.ReplaceAll(string(r), USE_KEYIN, ``)
}

func Atoa(str interface{}) string {
	if str == nil {
		return ""
	}
	return strings.Trim(str.(string), " ")
}

func Atoi(str interface{}) int {
	if str == nil {
		return 0
	}
	i, _ := strconv.Atoi(strings.Trim(str.(string), " "))
	return i
}

func Atoui(str interface{}) uint {
	if str == nil {
		return 0
	}
	u, _ := strconv.Atoi(strings.Trim(str.(string), " "))
	return uint(u)
}

// RandomCreateBytes generate random []byte by specify chars.
func RandomCreateBytes(n int, alphabets ...byte) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return bytes
}

// KeyinsParse splits user-provided custom keyins into unique tokens.
func KeyinsParse(keyins string) []string {
	keyins = strings.TrimSpace(keyins)
	if keyins == "" {
		return []string{}
	}
	for _, v := range re.FindAllString(keyins, -1) {
		keyins = strings.ReplaceAll(keyins, v, "><")
	}
	m := map[string]bool{}
	for _, v := range strings.Split(keyins, "><") {
		v = strings.TrimPrefix(v, "<")
		v = strings.TrimSuffix(v, ">")
		if v == "" {
			continue
		}
		m[v] = true
	}
	s := make([]string, len(m))
	i := 0
	for k := range m {
		s[i] = k
		i++
	}
	return s
}

// Bytes2String converts []byte to string via direct pointer conversion.
// Both share the same underlying memory; modifying one affects the other.
// Much faster than string([]byte{}) for large conversions.
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Bytes converts string to []byte via direct pointer conversion.
// Both share the same underlying memory; modifying one affects the other.
// Do not mutate the returned slice directly (e.g. b[1]='d') or the program may panic.
func String2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

package filesystem

import (
	"cmsApp/internal/constant"
	stringsx "cmsApp/pkg/utils/strings"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

/*
*
获取项目根目录
*/
func RootPath() (path string, err error) {
	path = GetCurrentAbPathByExecutable()

	if strings.Contains(path, GetTmpDir()) {
		path = GetCurrentAbPathByCaller()
	}
	path = strings.Replace(path, "pkg/utils/filesystem", "", 1)
	return
}

// 获取系统临时目录，兼容go run
func GetTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	if dir == "" {
		dir = "tmp"
	}

	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func GetCurrentAbPath() string {
	path := GetCurrentAbPathByCaller()
	return strings.Replace(path, "pkg/utils/filesystem", "", 1)
}

// 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

/**
* 打开文件句柄
**/
func OpenFile(filepath string) (file *os.File, err error) {

	file, err = os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err == nil {
		return
	}

	dir := path.Dir(filepath)
	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, fs.FileMode(os.O_CREATE))
			if err != nil {
				return
			}
		}
	}
	file, err = os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	return
}

/**
* 过滤非法访问的路径
 */
func FilterPath(root, path string) (string, error) {

	newPath := fmt.Sprintf("%s%s", root, path)
	absPath, err := filepath.Abs(newPath)
	if err != nil {
		return "", err
	}

	absPath = filepath.FromSlash(absPath)
	ifOver := strings.HasPrefix(absPath, filepath.FromSlash(root))
	if !ifOver {
		return "", errors.New("access to the path is prohibited")
	}

	return absPath, nil
}

// GetSize 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	return len(content), err
}

// GetExt 获取文件后缀
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// CheckNotExist 检查文件是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// CheckPermission 检查文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// IsNotExistMkDir 如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir 新建文件夹
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Open 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		//return nil, fmt.Errorf("os.Getwd err: %v", err)
		return nil, errors.New(constant.GET_CURRENT_PATH_ERR)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		//return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
		return nil, errors.New(constant.FILE_PERMISSION_ERR)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		//return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
		return nil, errors.New(constant.CREATE_DIR_ERR)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		//return nil, fmt.Errorf("Fail to OpenFile :%v", err)
		return nil, errors.New(constant.OPEN_FILE_ERR)
	}

	return f, nil
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

// 创建文件下载目录
func GetUploadDirs(name string) (string, string) {
	if name == "" {
		return "", ""
	}
	code := stringsx.GetHashCode(name)
	// 第一层目录 清除后四个bit位  也可以写成 x &= 0xf
	x := int(code) & 0xf
	dir1 := fmt.Sprintf("%x", x)
	// 第二层目录
	y := (int(code) >> 4) & 0xf
	dir2 := fmt.Sprintf("%x", y)
	return dir1, dir2
}

func GetUploadDirs2(name1 string, name2 string) (string, string) {
	if name1 == "" || name2 == "" {
		return "", ""
	}
	code1 := stringsx.GetHashCode(name1)
	code2 := stringsx.GetHashCode(name2)
	// 第一层目录 清除后四个bit位  也可以写成 x &= 0xf
	x := (int(code1) >> 4) & 0xf
	dir1 := fmt.Sprintf("%x", x)
	// 第二层目录
	y := int(code2) & 0xf
	dir2 := fmt.Sprintf("%x", y)
	return dir1, dir2
}

func IsDirEmpty(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil // 目录不存在，可以认为是空的
		}
		return false, err
	}
	return len(entries) == 0, nil
}

func main() {
	fmt.Println(FileExists("/Users/vega/workspace/codes/golang_space/gopath/src/work/cmsApp/uploadfile/2/0/1720020024615243776.png"))
}

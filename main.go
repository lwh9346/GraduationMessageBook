package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	//下面是初始化部分
	if !isFileOrDirectoryExists(path.Join(".", "users")) {
		os.Mkdir(path.Join(".", "users"), 664)
	}
	if !isFileOrDirectoryExists(path.Join(".", "received")) {
		os.Mkdir(path.Join(".", "received"), 664)
	}
	var conf config
	json.Unmarshal([]byte(readStringFromFile("config.json")), &conf)
	//下面是路由配置
	router := gin.Default()
	router.POST("/u/:usernameshort/submit", receiveData)
	router.GET("/u/:usernameshort/backtable", backtable)
	router.GET("/u/:usernameshort/backtable/:filename", backtablefile)
	router.GET("/u/:usernameshort", userIndexPage)
	router.GET("/u/:usernameshort/personal-page", pp)
	router.LoadHTMLGlob("./template/*")
	router.Static("/static", "./static")
	router.Run(":" + strconv.Itoa(conf.Port))
}

func backtable(context *gin.Context) {
	usernameShort := context.Param("usernameshort")
	userPassword := context.Query("password")
	userFilePath := path.Join(".", "users", usernameShort+".json")
	if isFileOrDirectoryExists(userFilePath) {
		var userInfo userInfo
		json.Unmarshal([]byte(readStringFromFile(userFilePath)), &userInfo)
		if userPassword == userInfo.Password {
			var filenames []string
			filesReceived, _ := ioutil.ReadDir(path.Join(".", "received", usernameShort))
			for _, f := range filesReceived {
				if !f.IsDir() {
					filenames = append(filenames, f.Name())
				}
			}
			context.HTML(200, "filelist.html", gin.H{"filelist": filenames, "password": userPassword})
			return
		}
	}
	context.JSON(404, gin.H{})
}

func backtablefile(context *gin.Context) {
	usernameShort := context.Param("usernameshort")
	userPassword := context.Query("password")
	filename := context.Param("filename")
	userFilePath := path.Join(".", "users", usernameShort+".json")
	if isFileOrDirectoryExists(userFilePath) {
		var userInfo userInfo
		json.Unmarshal([]byte(readStringFromFile(userFilePath)), &userInfo)
		if userPassword == userInfo.Password {
			context.String(200, readStringFromFile(path.Join(".", "received", usernameShort, filename)))
			return
		}
	}
	context.JSON(404, gin.H{})
}

func userIndexPage(context *gin.Context) {
	usernameShort := context.Param("usernameshort")
	userFilePath := path.Join(".", "users", usernameShort+".json")
	if isFileOrDirectoryExists(userFilePath) {
		var userInfo userInfo
		json.Unmarshal([]byte(readStringFromFile(userFilePath)), &userInfo)
		context.HTML(http.StatusOK, "index.html", gin.H{"fullname": userInfo.Fullname, "shortname": usernameShort, "vuereplacement": "{{name}}"})
		return
	}
	context.JSON(404, gin.H{})
}

func pp(context *gin.Context) {
	usernameShort := context.Param("usernameshort")
	name := context.Query("name")
	context.HTML(http.StatusOK, "form.html", gin.H{"name": strings.ReplaceAll(name, "\"", "\\\""), "usernameshort": usernameShort})
}

func receiveData(context *gin.Context) {
	name := context.PostForm("name")
	university := context.PostForm("university")
	major := context.PostForm("major")
	city := context.PostForm("city")
	contact := context.PostForm("contact")
	more := context.PostForm("more")
	index := 0
	usernameShort := context.Param("usernameshort")
	if !isFileOrDirectoryExists(path.Join(".", "received", usernameShort)) {
		os.Mkdir(path.Join(".", "received", usernameShort), 664)
	}
	for isFileOrDirectoryExists(path.Join(".", "received", usernameShort, name+"_"+strconv.Itoa(index)+".txt")) {
		index++
	}
	str := ""
	str += name + "\n"
	str += "大学：\n"
	str += university + "\n"
	str += "专业：\n"
	str += major + "\n"
	str += "城市：\n"
	str += city + "\n"
	str += "联系方式：\n"
	str += contact + "\n"
	str += "留言：\n"
	str += more
	writeStringToFile(path.Join(".", "received", usernameShort, name+"_"+strconv.Itoa(index)+".txt"), str)
	context.JSON(200, gin.H{"code": 200, "msg": "succeed", "name": name})
}

func isFileOrDirectoryExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func writeStringToFile(file, s string) error {
	var e error
	if isFileOrDirectoryExists(file) {
		e = os.Remove(file)
		if e != nil {
			return e
		}
	}
	f, e := os.Create(file)
	defer f.Close()
	if e != nil {
		return e
	}
	io.WriteString(f, s)
	return nil
}

func readStringFromFile(inFile string) string {
	b, err := ioutil.ReadFile(inFile)
	if err != nil {
		return ""
	}
	return string(b)
}

type config struct {
	Port int
}

type userInfo struct {
	Fullname string
	Password string
}

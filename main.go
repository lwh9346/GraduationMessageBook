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
	var conf config
	json.Unmarshal([]byte(readStringFromFile("config.json")), &conf)
	router := gin.Default()
	router.POST("/submit", receiveData)
	//router.StaticFile("/", "./static/index.html")
	router.GET("/u/:usernameshort", userIndexPage)
	router.GET("/u/:usernameshort/personal-page", pp)
	//router.StaticFS("/blackdoor", http.Dir("./received"))
	router.LoadHTMLGlob("./template/*")
	router.Static("/static", "./static")
	router.Run(":" + strconv.Itoa(conf.Port))
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
	for isFileOrDirectoryExists(path.Join(".", "received", name+"_"+strconv.Itoa(index)+".txt")) {
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
	writeStringToFile(path.Join(".", "received", name+"_"+strconv.Itoa(index)+".txt"), str)
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

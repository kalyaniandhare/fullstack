package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/icza/backscanner"
	"github.com/jinzhu/gorm"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kalyaniandhare/fullstack/api/models"
	"github.com/kalyaniandhare/fullstack/api/responses"
	//"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

func (server *Server) CreateLogConfig(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.LogConfig{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	//err = user.Validate("")
	//if err != nil {
	//	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	//	return
	//}
	userCreated, err := user.SaveLogConfig(server.DB)

	if err != nil {

		//formattedError := formaterror.FormatError(err.Error())
		//
		//responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	CreateConfig(userCreated, server.DB)
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)


}


func predefineLogLevels(currentLog string) []string {
	var x []string
	currentLog = strings.ToLower(currentLog)
	if currentLog == "debug" {
		x = append(x, "debug","trace")
		return x

	}
	if currentLog == "trace" {
		x = append(x, "trace")
		return x

	}
	if currentLog == "fatal" {
		x = append(x, "fatal","error","warn","info","debug","trace")
		return x

	}
	if currentLog == "error" {
		x = append(x,"error","warn","info","debug","trace")
		return x

	}
	if currentLog == "warn" {
		x = append(x,"warn","info","debug","trace")
		return x

	}
	if currentLog == "info" {
		x = append(x,"info","debug","trace")
		return x

	}
	if currentLog == "off" {
		x = append(x, "fatal","error","warn","info","debug","trace","off")
		return x

	}
	if currentLog == "all" {
		x = append(x,"fatal","error","warn","info","debug","trace","off")
		return x

	}
	return x

}
//func testscanner() {
//	src, err := os.Open("app.txt")
//	scanner := NewScanner(strings.NewReader(src), len(src))
//	for {
//		line, pos, err := scanner.Line()
//		if err != nil {
//			fmt.Println("Error:", err)
//			break
//		}
//		fmt.Printf("Line start: %2d, line: %s\n", pos, line)
//	}
//
//}
//func getSeekLocation() int64 {
//	start := int64(0)
//	input, err := os.Open("app.txt")
//	if err != nil {
//		fmt.Println(err)
//	}
//	if _, err := input.Seek(start, io.SeekStart); err != nil {
//		fmt.Println(err)
//	}
//	scanner := bufio.NewScanner(input)
//
//	pos := start
//	scanLines := func(data []byte, atEOF bool) (advance int, token []byte,
//		err error) {
//		advance, token, err = bufio.ScanLines(data, atEOF)
//		pos += int64(advance)
//		return
//	}
//
//	scanner.Split(scanLines)
//	for scanner.Scan(){
//		fmt.Println(scanner.Text(),"@@@@@@@@@@@@@")
//	}
//	fmt.Println(scanner.Scan(),"#################")
//	fi, err :=input.Stat()
//	size := int64(fi.Size())
//	return size - pos
//}


func CreateConfig(obj *models.LogConfig, db *gorm.DB) {
	f, err := os.Open("app.txt")
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	defer f.Close()

	listOfLevels := predefineLogLevels(obj.LogLevel)

	intervalInSec := obj.Interval
	currentIntervalTime := intervalInSec /60
	scanner := backscanner.New(f, int(fi.Size()))

	for {
		line, pos, err := scanner.Line()
		var result map[string]interface{}
		json.Unmarshal([]byte(line), &result)


		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		currentDate := time.Now()
		currentDate.In(time.UTC)


		datetime := fmt.Sprint(result["time"])
		level := fmt.Sprint(result["level"])
		message := fmt.Sprint(result["message"])

		//dd := time.
		fileDate, err := time.Parse(time.RFC3339, datetime)
		now := time.Now()
		currentTime := now.Add(time.Duration(-currentIntervalTime) * time.Minute)

		if fileDate.After(currentTime) == true || fileDate.Equal(currentTime) == true {
			fmt.Println(level, listOfLevels)
			if stringInSlice(level, listOfLevels) {
				post := models.Post{}

				post.DateTime=datetime
				post.AlertMessage = message
				post.AlertLogLevel = level
				post.LogConfigID = obj.ID
				post.SavePostNEW(db)
				fmt.Println("Will save time", fileDate, pos)
			}
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (server *Server) GetLogsConfig(w http.ResponseWriter, r *http.Request) {

	user := models.LogConfig{}

	users, err := user.FindAllLogs(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) GetLogDetail(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	user := models.Post{}

	users, err := user.FindLogByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)

}


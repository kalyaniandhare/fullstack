package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kalyaniandhare/fullstack/api/models"
	"github.com/kalyaniandhare/fullstack/api/responses"
)

func (server *Server) CreateLog(w http.ResponseWriter, r *http.Request) {

	fmt.Println("REQUEST00000000000000000000000000", r.Body)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	fmt.Println(reflect.TypeOf(body),reflect.TypeOf(post), reflect.TypeOf(&post),"ZZZZZZZZZZZZZZZZZZZZZZZZZ")
	fmt.Println(body, "BBBBBBBBBBBBBBBBBB",post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	//uid, err := auth.ExtractTokenID(r)
	//if err != nil {
	//	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	//	return
	//}
	//if uid != post.LogConfigID {
	//	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	//	return
	//}
	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		fmt.Println(err.Error(),"GGGGGGGGG")
		//formattedError := formaterror.FormatError(err.Error())
		//responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	fmt.Println(postCreated,"GGGGGGGGGgg")
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)
}

func (server *Server) GetAllLogs(w http.ResponseWriter, r *http.Request) {

	post := models.Post{}

	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (server *Server) GetLog(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	post := models.Post{}

	postReceived, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}
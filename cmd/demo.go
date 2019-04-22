package main

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/db"
	"dm/fieldtype"
	"dm/handler"
	"dm/query"
	"dm/util"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//go:generate go run gen_contenttypes/gen.go
func BootStrap() {
	if len(os.Args) >= 2 && os.Args[1] != "" {
		path := os.Args[1]
		util.SetConfigPath(path + "/configs")
	}
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

}

//This is a initial try which use template to do basic feature.

func Display(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	tpl := template.Must(template.ParseFiles("../web/template/view.html"))
	rmdb := db.DBHanlder()
	article := entity.Article{}
	id, _ := strconv.Atoi(vars["id"])

	err := rmdb.GetByID("article", id, &article)

	if err != nil {
		fmt.Println(err)
	}

	//List of folder
	folders, _ := handler.Querier().List("folder", query.Cond("parent_id", 0))

	//Get current Folder
	currentFolder, _ := handler.Querier().Fetch("folder", query.Cond("location.id", id))

	var variables map[string]interface{}
	if currentFolder != nil {
		//Get list of article
		articles, _ := handler.Querier().List("article", query.Cond("parent_id", id))

		variables = map[string]interface{}{"current": currentFolder,
			"list":    articles,
			"folders": folders}
	} else {
		currentArticle, _ := handler.Querier().Fetch("article", query.Cond("location.id", id))
		variables = map[string]interface{}{"current": currentArticle,
			"list":    nil,
			"folders": folders}
	}

	folderList, _ := handler.Querier().List("folder", query.Cond("parent_id", id))
	variables["folder_list"] = folderList
	tpl.Execute(w, variables)
}

func New(w http.ResponseWriter, r *http.Request) {
	// handler := handler.ContentHandler{}

	vars := mux.Vars(r)

	variables := map[string]interface{}{}
	variables["id"] = vars["id"]
	variables["type"] = vars["type"]
	variables["posted"] = false
	if r.Method == "POST" {
		variables["posted"] = true
		if vars["type"] == "article" {
			parentID, _ := strconv.Atoi(vars["id"])
			title := r.FormValue("title")
			body := r.FormValue("body")
			handler := handler.ContentHandler{}
			success, result, error := handler.Create(parentID, "article", map[string]interface{}{"title": title, "body": body})
			if !success {
				variables["success"] = false
				if error != nil {
					variables["error"] = error.Error()
				}
				variables["validation"] = result
			} else {
				variables["success"] = true
			}
		}

	}
	tpl := template.Must(template.ParseFiles("../web/template/new_" + vars["type"] + ".html"))
	//variables := map[string]interface{}{}
	tpl.Execute(w, variables)
}

func Publish(w http.ResponseWriter, r *http.Request) {

}

func main() {

	BootStrap()
	r := mux.NewRouter()
	r.HandleFunc("/content/view/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Display(w, r, vars)
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Display(w, r, map[string]string{"id": "1"})
	})
	// http.HandleFunc("/content/view/", func(w http.ResponseWriter, r *http.Request) {
	// 	Display(w, r)
	// })

	r.HandleFunc("/content/new/{type}/{id}", func(w http.ResponseWriter, r *http.Request) {
		New(w, r)
	})

	r.HandleFunc("/content/publish", func(w http.ResponseWriter, r *http.Request) {
		Publish(w, r)
	})

	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "User-agent: * \nDisallow /")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../web"))))

	http.Handle("/", r)
	http.ListenAndServe(":8089", nil)
}

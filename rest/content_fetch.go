//Author xc, Created on 2019-08-11 16:49
//{COPYRIGHTS}

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/core/handler"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/permission"
	"github.com/xc/digimaker/core/util"

	"github.com/gorilla/mux"
)

func GetContent(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	querier := handler.Querier()
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		HandleError(errors.New("Invalid id"), w)
		return
	}
	var content contenttype.ContentTyper
	contentType := params["contenttype"]
	if contentType != "" {
		content, err = querier.FetchByContentID(contentType, id)
	} else {
		content, err = querier.FetchByID(id)
	}
	if err != nil {
		HandleError(err, w)
		return
	} else {
		if !permission.CanRead(userID, content, r.Context()) {
			HandleError(errors.New("Doesn't have permission."), w, 403)
			return
		}
		w.Header().Set("content-type", "application/json")
		data, _ := contenttype.ContentToJson(content) //todo: use export for same serilization?
		w.Write([]byte(data))
	}

}

func GetVersion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	querier := handler.Querier()
	id, _ := strconv.Atoi(params["id"])

	versionNo, _ := strconv.Atoi(params["version"])

	querier = handler.Querier()
	content, err := querier.FetchByID(id)
	if err != nil {
		HandleError(errors.New("Can not find content. error: "+err.Error()), w)
		return
	}

	if content == nil {
		HandleError(errors.New("Content doesn't exist"), w)
		return
	}

	if !permission.CanRead(userID, content, r.Context()) {
		HandleError(errors.New("No permisison to the content."), w)
		return
	}

	maxVersion := content.Value("version").(int)

	if versionNo > maxVersion {
		HandleError(errors.New("version doesn't exist."), w)
		return
	}

	dbHandler := db.DBHanlder()
	version := contenttype.Version{}
	dbHandler.GetEntity(version.TableName(),
		db.Cond("content_id", content.GetCID()).Cond("content_type", content.ContentType()).Cond("version", versionNo),
		[]string{},
		nil,
		&version)
	if version.ID == 0 {
		HandleError(errors.New("version doesn't exist."), w)
		return
	}

	data, _ := json.Marshal(version)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(data))
}

func buildCondition(userid int, def contenttype.ContentType, query url.Values) (db.Condition, error) {
	author := query.Get("author")
	condition := db.EmptyCond()
	if author != "" {
		if author == "self" {
			condition = condition.Cond("author", userid)
		} else {
			authorInt, err := strconv.Atoi(author)
			if err != nil {
				return db.EmptyCond(), errors.New("wrong author format")
			}
			condition = condition.Cond("author", authorInt)
		}
	}

	//id
	idStr := query.Get("id")
	if idStr != "" {
		ids, err := util.ArrayStrToInt(strings.Split(idStr, ","))
		if err != nil {
			return db.EmptyCond(), errors.New("Wrong id format")
		}
		condition = condition.And("location.id", ids)
	}

	//cid
	cidStr := query.Get("cid")
	if cidStr != "" {
		cids, err := util.ArrayStrToInt(strings.Split(cidStr, ","))
		if err != nil {
			return db.EmptyCond(), errors.New("Wrong cid format")
		}
		condition = condition.And("c.id", cids)
	}
	//contain savan
	nameStr := query.Get("name")

	if nameStr != "" {
		contains := util.Split(nameStr, ":")
		if len(contains) > 0 && contains[0] == "contain" {
			fmt.Printf("savan" + contains[1])
			//todo: esc % to inside in condition
			cValue := "%" + contains[1] + "%"
			condition = condition.And("location.name like", cValue)
		}
	}
	//filter

	for field := range def.FieldMap {
		value := query.Get("field." + field)
		if value != "" {
			condition = condition.And("c."+field, value) //todo: support operator //todo: support more value type(eg. array)
		}
	}

	return condition, nil
}

func BuildSortby(r *http.Request) []string {
	getParams := r.URL.Query()
	sortbyStr := getParams.Get("sortby")
	sortbyArr := util.Split(sortbyStr, ";")
	return sortbyArr
}

func BuildLimit(r *http.Request) ([]int, error) {
	getParams := r.URL.Query()
	offsetStr := getParams.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if offsetStr != "" && err != nil {
		return nil, errors.New("Invalid offset")
	}

	limitStr := getParams.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if limitStr != "" && err != nil {
		return nil, errors.New("Invalid limit")
	}

	return []int{offset, limit}, nil
}

//List
func List(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	getParams := r.URL.Query()

	limit, err := BuildLimit(r)
	if err != nil {
		HandleError(err, w)
		return
	}

	sortby := BuildSortby(r)

	ctype := params["contenttype"]
	def, err := contenttype.GetDefinition(ctype)
	if err != nil {
		HandleError(errors.New("Cann't get content type"), w)
		return
	}

	querier := handler.Querier()

	ctx := r.Context()
	userID := CheckUserID(ctx, w)
	if userID == 0 {
		return
	}

	//filter
	condition, err := buildCondition(userID, def, r.URL.Query())
	if err != nil {
		HandleError(err, w, 410)
		return
	}

	if limit[0] == 0 && limit[1] == 0 {
		limit = []int{0, 10} //todo: use configuration
	}

	rootStr := getParams.Get("parent")
	var list []contenttype.ContentTyper
	var count int
	if rootStr != "" {
		var rootContent contenttype.ContentTyper
		rootID, err := strconv.Atoi(rootStr)
		if err != nil {
			HandleError(errors.New("Invalid parent"), w)
			return
		}
		rootContent, err = querier.FetchByID(rootID)
		if err != nil {
			log.Error(err.Error(), "", r.Context())
			HandleError(errors.New("Can't get parent"), w, 410)
			return
		}

		//level
		levelStr := getParams.Get("level")
		level := 0
		if levelStr != "" {
			level, err = strconv.Atoi(levelStr)
			if err != nil {
				HandleError(errors.New("Invalid level"), w)
				return
			}
		}

		list, count, err = querier.SubList(rootContent, ctype, level, userID, condition, limit, sortby, true, ctx)
		if err != nil {
			HandleError(err, w)
			return
		}
	} else {
		list, count, err = querier.ListForUser(ctx, userID, ctype, condition, limit, sortby, true)
		if err != nil {
			HandleError(err, w)
			return
		}
	}

	result := struct {
		List  interface{} `json:"list"`
		Count int         `json:"count"`
	}{Count: count}

	result.List = list

	data, _ := json.Marshal(result)
	w.Write([]byte(data))
}

//Get content list from relation definition
func RelationOptionList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	contenttype := params["contenttype"]
	fieldIdentifier := params["field"]
	if contenttype == "" || fieldIdentifier == "" {
		HandleError(errors.New("Need type and field"), w, 410)
		return
	}

	cxt := r.Context()
	userid := CheckUserID(cxt, w)
	if userid == 0 {
		return
	}

	querier := handler.Querier()
	list, count, err := querier.RelationOptions(contenttype, fieldIdentifier, nil, nil, false)
	if err != nil {
		HandleError(err, w)
		return
	}

	result := struct {
		List  interface{} `json:"list"`
		Count int         `json:"count"`
	}{Count: count}

	result.List = list

	data, _ := json.Marshal(result)
	w.Write([]byte(data))
}

//Get tree menu under a node
func TreeMenu(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	types := r.FormValue("type")
	if types == "" {
		HandleError(errors.New("Wrong parameters"), w, StatusWrongParams)
		return
	}
	typeList := util.Split(types)
	contenttypeList := contenttype.GetDefinitionList()["default"]
	for _, ctype := range typeList {
		if _, ok := contenttypeList[ctype]; !ok {
			HandleError(errors.New("Invalid type"), w, StatusWrongParams)
			return
		}
	}

	querier := handler.Querier()
	rootContent, err := querier.FetchByID(id)

	if err != nil {
		HandleError(err, w)
		return
	}

	if rootContent == nil {
		HandleError(errors.New("content Not found"), w, 403)
		return
	}

	if !permission.CanRead(userID, rootContent, r.Context()) {
		HandleError(errors.New("No permission"), w)
		return
	}

	tree, err := querier.SubTree(rootContent, 5, strings.Join(typeList, ","), userID, []string{"priority desc", "id"}, r.Context())
	if err != nil {
		HandleError(err, w)
		return
	}

	//todo: make this configurable
	tree.Iterate(func(node *handler.TreeNode) {
		if node.ContentType == "folder" {
			node.Fields = map[string]interface{}{"subtype": node.Content.Value("folder_type")}
		}
	})

	data, _ := json.Marshal(tree)
	w.Write([]byte(data))
}

func init() {
	RegisterRoute("/content/get/{id:[0-9]+}", GetContent)
	RegisterRoute("/content/get/{contenttype}/{id:[0-9]+}", GetContent)
	RegisterRoute("/content/version/{id:[0-9]+}/{version:[0-9]+}", GetVersion)

	RegisterRoute("/content/treemenu/{id:[0-9]+}", TreeMenu)
	RegisterRoute("/content/list/{contenttype}", List)
	RegisterRoute("/relation/optionlist/{contenttype}/{field}", RelationOptionList)
}

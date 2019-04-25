package handler

import (
	"dm/contenttype"
	"dm/db"
	"fmt"
)

type RelationHandler struct {
}

//Add a content to current content(toContent)
func (handler *RelationHandler) AddTo(toContent contenttype.ContentTyper, from contenttype.ContentTyper, identifier string, priority int, description string) error {
	db := db.DBHanlder()
	fmt.Println("FROM: ")
	fmt.Println(from.ToMap()["id"])
	fmt.Println("TO:")
	fmt.Println(toContent)
	values := map[string]interface{}{
		"from_location": from.ToMap()["id"],
		"to_content":    toContent.Value("cid"),
		"relation_type": contenttype.GetContentDefinition(toContent.ContentType()).Fields[identifier].FieldType,
		"priority":      0,
		"identifier":    identifier} //todo: get data from value pattern
	id, err := db.Insert("dm_relation", values)
	fmt.Println(id)
	fmt.Println(err)
	return err
}

//Update all contents which is related to current content(fromContent)
func (handler *RelationHandler) UpdateValues(fromContent contenttype.ContentTyper) {

}
package fieldtype

//All of the fields will implements this interface
type Fieldtyper interface {
	//Get value of
	//Value() string

	//Create()
	//Validate()
	//SetStoreData()
}

type FieldtypeHandler interface {
	ToStorage(input interface{}) interface{}
	Validate(input interface{}) (bool, string)
	IsEmpty(input interface{}) bool
}

//Relation field handler can convert relations into RelationField
type RelationFieldHandler interface {
	ToStorage(contents interface{}) interface{}
	UpdateOne(toContent interface{}, identifier string, from interface{})
}
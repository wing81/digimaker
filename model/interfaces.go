package model

//All the content type(eg. article, folder) will implement this interface.
type ContentTyper interface {
	// Since go's embeded struct can't really inherit well from BaseContentType(eg. ID)
	// (It refers to embeded struct instance instead of override all fields by default)
	// Interface in go is like a kind of 'abstract class' when it comes to generic with data.
	// We use property of instance to declear a general ContentType. This will integrate well with orm enitty.
	ID() int
	Published() int
	Modified() int
	RemoteID() string

	//Return all fields
	Fields() map[string]Field

	//Visit  field dynamically
	Field(name string) Field

	//Visit all attribute dynamically including Fields + internal attribute eg. id, parent_id.
	Attr(name string) interface{}
}

//All of the fields will implements this interface
type Fielder interface {
	//Get value of
	Value()
	Create()
	Validate()
	SetStoreData()
}

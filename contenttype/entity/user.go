//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
    "dm/contenttype"
	"dm/fieldtype"
	. "dm/query"
)



type User struct{
     ContentCommon `boil:",bind"`
    
     
     Firstname fieldtype.TextField `boil:"firstname" json:"firstname" toml:"firstname" yaml:"firstname"`
    
     
     Lastname fieldtype.TextField `boil:"lastname" json:"lastname" toml:"lastname" yaml:"lastname"`
    
     
     Login fieldtype.TextField `boil:"login" json:"login" toml:"login" yaml:"login"`
    
     
     Password fieldtype.TextField `boil:"password" json:"password" toml:"password" yaml:"password"`
    
     Location `boil:"location,bind"`
}

func ( *User ) TableName() string{
	 return "dm_user"
}

func (c *User) contentValues() map[string]interface{} {
	result := make(map[string]interface{})
    
        result["firstname"]=c.Firstname
    
        result["lastname"]=c.Lastname
    
        result["login"]=c.Login
    
        result["password"]=c.Password
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *User) Values() map[string]interface{} {
    result := c.contentValues()

	for key, value := range c.Location.Values() {
		result[key] = value
	}
	return result
}

func (c *User) Value(identifier string) interface{} {
	var result interface{}
	switch identifier {
    
    case "firstname":
        result = c.Firstname
    
    case "lastname":
        result = c.Lastname
    
    case "login":
        result = c.Login
    
    case "password":
        result = c.Password
    
	case "cid":
		result = c.ContentCommon.CID
    default:
    	result = c.ContentCommon.Value( identifier )
    }
	return result
}


func (c *User) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
             
            case "firstname":
            c.Firstname = value.(fieldtype.TextField)
        
             
            case "lastname":
            c.Lastname = value.(fieldtype.TextField)
        
             
            case "login":
            c.Login = value.(fieldtype.TextField)
        
             
            case "password":
            c.Password = value.(fieldtype.TextField)
        
	default:
		err := c.ContentCommon.SetValue(identifier, value)
        if err != nil{
            return err
        }
	}
	//todo: check if identifier exist
	return nil
}

//Store content.
//Note: it will set id to CID after success
func (c *User) Store() error {
	handler := db.DBHanlder()
	if c.CID == 0 {
		id, err := handler.Insert(c.TableName(), c.contentValues())
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.contentValues(), Cond("id", c.CID))
		return err
	}
	return nil
}


func init() {
	new := func() contenttype.ContentTyper {
		return &User{}
	}

	newList := func() interface{} {
		return &[]User{}
	}

	Register("user",
		ContentTypeRegister{
			New:            new,
			NewList:        newList})
}

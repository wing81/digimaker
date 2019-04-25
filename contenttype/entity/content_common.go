package entity

type ContentCommon struct {
	CID       int              `boil:"id" json:"id" toml:"id" yaml:"id"`
	Published int              `boil:"published" json:"published" toml:"published" yaml:"published"`
	Modified  int              `boil:"modified" json:"modified" toml:"modified" yaml:"modified"`
	CUID      string           `boil:"cuid" json:"cuid" toml:"cuid" yaml:"cuid"`
	Relations ContentRelations `boil:"relations" json:"relations" toml:"relations" yaml:"relations"`
}

func (c ContentCommon) IdentifierList() []string {
	return []string{"cid", "published", "modified", "cuid"}
}

func (c ContentCommon) Values() map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = c.CID
	result["published"] = c.Published
	result["modified"] = c.Modified
	result["cuid"] = c.CUID
	return result
}

func (c *ContentCommon) Value(identifier string) interface{} {
	var result interface{}
	switch identifier {
	case "cid":
		result = c.CID
	case "modified":
		result = c.Modified
	case "published":
		result = c.Published
	case "cuid":
		result = c.CUID
	}
	return result
}

func (c *ContentCommon) SetValue(identifier string, value interface{}) error {
	switch identifier {
	case "id":
		c.CID = value.(int)
	case "published":
		c.Published = value.(int)
	case "modified":
		c.Modified = value.(int)
	case "cuid":
		c.CUID = value.(string)
	}
	return nil
}

func GetCID(c *ContentCommon) int {
	return c.CID
}

//TODO: add more common methods related to content here.

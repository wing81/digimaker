//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"

	_ "github.com/go-sql-driver/mysql" //todo: move this to loader
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries"
)

// Implement DBEntitier
type MysqlHandler struct {
}

//Query by ID
func (rmdb *MysqlHandler) GetByID(contentType string, tableName string, id int, content interface{}) error {
	_, err := rmdb.GetByFields(contentType, tableName, Cond("location.id", id), []int{}, []string{}, content, false)
	return err
}

//Query to fill in contentTyper. Use reference in content parameter.
//It fill in with nil if nothing found(no error returned in this case)
//  var content contenttype.Article
//  rmdb.GetByFields("article", map[string]interface{}{"id": 12}, {{"name","asc"}} content)
//
//todo: possible to have more joins between content/entities(relations or others), or ingegrate with ORM
func (r *MysqlHandler) GetByFields(contentType string, tableName string, condition Condition, limit []int, sortby []string, content interface{}, count bool) (int, error) {
	db, err := DB()
	if err != nil {
		return -1, errors.Wrap(err, "[MysqlHandler.GetByFields]Error when connecting db.")
	}

	columns := util.GetInternalSettings("location_columns")
	columnsWithPrefix := util.Iterate(columns, func(s string) string {
		return `location.` + s + ` AS "location.` + s + `"`
	})
	locationColumns := strings.Join(columnsWithPrefix, ",")

	//get condition string for fields
	conditionStr, values := BuildCondition(condition, columns)
	where := ""
	if conditionStr != "" {
		where = "WHERE " + conditionStr
	}

	relationQuery := r.getRelationQuery()

	//limit
	limitStr := ""
	if len(limit) > 0 {
		if len(limit) != 2 {
			return -1, errors.New("limit should be array with only 2 int. There are: " + strconv.Itoa(len(limit)))
		}
		limitStr = " LIMIT " + strconv.Itoa(limit[0]) + "," + strconv.Itoa(limit[1])
	}

	//sort by
	sortbyStr, err := r.getSortBy(sortby, columns)
	if err != nil {
		return -1, err
	}

	sqlStr := `SELECT c.*, c.id AS cid, location_user.name AS author_name, ` + locationColumns + relationQuery + `
                   FROM (` + tableName + ` c INNER JOIN dm_location location ON location.content_type = '` + contentType + `' AND location.content_id=c.id)
                     LEFT JOIN dm_relation relation ON c.id=relation.to_content_id AND relation.to_type='` + contentType + `'
										 LEFT JOIN dm_location location_user ON location_user.content_type='user' AND location_user.content_id=c.author
                    ` + where + `
                     GROUP BY location.id, author_name
										 ` + sortbyStr + " " + limitStr

	log.Debug(sqlStr+","+fmt.Sprintln(values), "db")
	err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, content)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug(err.Error(), "GetByFields")
		} else {
			message := "[MysqlHandler.GetByFields]Error when query. sql - " + sqlStr
			return -1, errors.Wrap(err, message)
		}
	}

	//count if there is
	countResult := 0
	if count {
		countSqlStr := `SELECT COUNT(*) AS count
									 FROM ( ` + tableName + ` c
										 INNER JOIN dm_location location ON location.content_type = '` + contentType + `' AND location.content_id=c.id )
										 ` + where

		rows, err := queries.Raw(countSqlStr, values...).QueryContext(context.Background(), db)
		if err != nil {
			message := "[MysqlHandler.GetByFields]Error when query count. sql - " + countSqlStr
			return -1, errors.Wrap(err, message)
		}
		rows.Next()
		rows.Scan(&countResult)
		rows.Close()
	}

	return countResult, nil
}

//Get non-location content
//todo: possible to have more joins between entities, or ingegrate with ORM
//todo: support select multiple entity once.
//todo: support query without involing location at all.
func (r *MysqlHandler) GetEntityContent(contentType string, tableName string, condition Condition, limit []int, sortby []string, content interface{}, count bool) (int, error) {
	db, err := DB()
	if err != nil {
		return -1, errors.Wrap(err, "[MysqlHandler.GetByFields]Error when connecting db.")
	}
	//get condition string for fields
	conditionStr, values := BuildCondition(condition)
	where := ""
	if conditionStr != "" {
		where = "WHERE " + conditionStr
	}

	relationQuery := r.getRelationQuery()

	//limit
	limitStr := ""
	if len(limit) > 0 {
		if len(limit) != 2 {
			return -1, errors.New("limit should be array with only 2 int. There are: " + strconv.Itoa(len(limit)))
		}
		limitStr = " LIMIT " + strconv.Itoa(limit[0]) + "," + strconv.Itoa(limit[1])
	}

	//sort by
	sortbyStr, err := r.getSortBy(sortby)
	if err != nil {
		return -1, err
	}

	sqlStr := `SELECT c.*, c.id as cid, '` + contentType + `' as content_type` + relationQuery + `
                   FROM (` + tableName + ` c INNER JOIN dm_location location ON c.location_id = location.id )
                     LEFT JOIN dm_relation relation ON c.id=relation.to_content_id AND relation.to_type='` + contentType + `'
                    ` + where + `
                     GROUP BY c.id
										 ` + sortbyStr + " " + limitStr

	log.Debug(sqlStr+","+fmt.Sprintln(values), "db")
	err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, content)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Warning(err.Error(), "GetByFields")
		} else {
			message := "[MysqlHandler.GetByFields]Error when query. sql - " + sqlStr
			return -1, errors.Wrap(err, message)
		}
	}

	//count if there is
	countResult := 0
	if count {
		countSqlStr := `SELECT COUNT(*) AS count FROM ` + tableName + ` c INNER JOIN dm_location location ON c.location_id = location.id ` + where

		rows, err := queries.Raw(countSqlStr, values...).QueryContext(context.Background(), db)
		if err != nil {
			message := "[MysqlHandler.GetByFields]Error when query count. sql - " + countSqlStr
			return -1, errors.Wrap(err, message)
		}
		rows.Next()
		rows.Scan(&countResult)
		rows.Close()
	}

	return countResult, nil
}

func (r *MysqlHandler) getRelationQuery() string {
	relationQuery := `,JSON_ARRAYAGG( JSON_OBJECT( 'identifier', relation.identifier,
                                      'to_content_id', relation.to_content_id,
                                      'to_type', relation.to_type,
                                      'from_content_id', relation.from_content_id,
                                      'from_type', relation.from_type,
                                      'from_location', relation.from_location,
                                      'priority', relation.priority,
                                      'uid', relation.uid,
                                      'description',relation.description,
                                      'data' ,relation.data ) ) AS relations`
	return relationQuery
}

//Get sort by sql based on sortby pattern(eg.[]string{"name asc", "id desc"})
func (r *MysqlHandler) getSortBy(sortby []string, locationColumns ...[]string) (string, error) {
	//sort by
	sortbyArr := []string{}
	for _, item := range sortby {
		if strings.TrimSpace(item) != "" {
			itemArr := util.Split(item, " ")
			sortByField := itemArr[0]
			if len(locationColumns) > 0 && util.Contains(locationColumns[0], sortByField) {
				sortByField = "location." + sortByField
			}
			sortByOrder := "ASC"

			if len(itemArr) == 2 {
				sortByOrder = strings.ToUpper(itemArr[1])
				if sortByOrder != "ASC" && sortByOrder != "DESC" {
					return "", errors.New("Invalid sorting string: " + sortByOrder)
				}
			}
			sortbyItem := sortByField + " " + sortByOrder
			sortbyArr = append(sortbyArr, sortbyItem)
		}
	}

	sortbyStr := ""
	if len(sortbyArr) > 0 {
		sortbyStr = "ORDER BY " + strings.Join(sortbyArr, ",")
		sortbyStr = util.StripSQLPhrase(sortbyStr)
	}
	return sortbyStr, nil
}

// Count based on condition
func (*MysqlHandler) Count(tablename string, condition Condition) (int, error) {
	conditions, values := BuildCondition(condition)
	sqlStr := "SELECT COUNT(*) AS count FROM " + tablename + " WHERE " + conditions
	log.Debug(sqlStr, "db")
	db, err := DB()
	if err != nil {
		return 0, errors.Wrap(err, "[MysqlHandler.Count]Error when connecting db.")
	}
	rows, err := queries.Raw(sqlStr, values...).QueryContext(context.Background(), db)
	if err != nil {
		return 0, errors.Wrap(err, "[MysqlHandler.Count]Error when querying.")
	}
	rows.Next()
	var count int
	rows.Scan(&count)
	rows.Close()
	return count, nil
}

//todo: support limit.
func (r *MysqlHandler) GetEntity(tablename string, condition Condition, sortby []string, limit []int, entity interface{}) error {
	conditions, values := BuildCondition(condition)
	sortbyStr, err := r.getSortBy(sortby)
	if err != nil {
		return err
	}

	limitStr := ""
	if limit != nil && len(limit) == 2 {
		limitStr = " LIMIT " + strconv.Itoa(limit[0]) + "," + strconv.Itoa(limit[1])
	}
	sqlStr := "SELECT * FROM " + tablename + " WHERE " + conditions + " " + sortbyStr + limitStr
	log.Debug(sqlStr, "db")
	db, err := DB()
	if err != nil {
		return errors.Wrap(err, "[MysqlHandler.GetEntity]Error when connecting db.")
	}
	err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, entity)
	if err == sql.ErrNoRows {
		log.Warning(err.Error(), "GetEntity")
	} else {
		return errors.Wrap(err, "[MysqlHandler.GetEntity]Error when query.")
	}
	return nil
}

//Fetch multiple enities
func (*MysqlHandler) GetMultiEntities(tablenames []string, condition Condition, entity interface{}) {

}

func (MysqlHandler) Insert(tablename string, values map[string]interface{}, transation ...*sql.Tx) (int, error) {
	sqlStr := "INSERT INTO " + tablename + " ("
	valuesString := "VALUES("
	var valueParameters []interface{}
	if len(values) > 0 {
		for name, value := range values {
			if name != "id" {
				sqlStr += name + ","
				valuesString += "?,"
				valueParameters = append(valueParameters, value)
			}
		}
		sqlStr = sqlStr[:len(sqlStr)-1]
		valuesString = valuesString[:len(valuesString)-1]
	}
	sqlStr += ")"
	valuesString += ")"
	sqlStr = sqlStr + " " + valuesString
	log.Debug(sqlStr, "db")

	var result sql.Result
	var error error
	//execute using and without using transaction
	if len(transation) == 0 {
		db, err := DB()
		if err != nil {
			return 0, errors.Wrap(err, "MysqlHandler.Insert] Error when getting db connection.")
		}
		//todo: create context to isolate queries.
		result, error = db.ExecContext(context.Background(), sqlStr, valueParameters...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, valueParameters...)
	}
	//execution error
	if error != nil {
		return 0, errors.Wrap(error, "MysqlHandler.Insert]Error when executing. sql - "+sqlStr)
	}
	id, err := result.LastInsertId()
	//Get id error
	if err != nil {
		return 0, errors.Wrap(err, "MysqlHandler.Insert]Error when inserting. sql - "+sqlStr)
	}

	log.Debug("Insert results in id: "+strconv.FormatInt(id, 10), "db")

	return int(id), nil
}

//Generic update an entity
func (MysqlHandler) Update(tablename string, values map[string]interface{}, condition Condition, transation ...*sql.Tx) error {
	sqlStr := "UPDATE " + tablename + " SET "
	var valueParameters []interface{}
	for name, value := range values {
		if name != "id" {
			sqlStr += name + "=?,"
			valueParameters = append(valueParameters, value)
		}
	}

	sqlStr = sqlStr[:len(sqlStr)-1]
	conditionString, conditionValues := BuildCondition(condition)
	valueParameters = append(valueParameters, conditionValues...)
	sqlStr += " WHERE " + conditionString

	log.Debug(sqlStr, "db")

	var result sql.Result
	var error error
	if len(transation) == 0 {
		db, err := DB()
		if err != nil {
			return errors.Wrap(err, "MysqlHandler.Update] Error when getting db connection.")
		}
		result, error = db.ExecContext(context.Background(), sqlStr, valueParameters...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, valueParameters...)
	}
	if error != nil {
		return errors.Wrap(error, "[MysqlHandler.Update]Error when updating. sql - "+sqlStr)
	}
	resultRows, _ := result.RowsAffected()
	log.Debug("Updated rows:"+strconv.FormatInt(resultRows, 10), "db")
	return nil
}

//Delete based on condition
func (*MysqlHandler) Delete(tableName string, condition Condition, transation ...*sql.Tx) error {
	conditionString, conditionValues := BuildCondition(condition)
	sqlStr := "DELETE FROM " + tableName + " WHERE " + conditionString

	log.Debug(sqlStr, "db")

	var result sql.Result
	var error error

	if len(transation) == 0 {
		db, err := DB()
		if err != nil {
			return errors.Wrap(err, "MysqlHandler.Delete] Error when getting db connection.")
		}
		result, error = db.ExecContext(context.Background(), sqlStr, conditionValues...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, conditionValues...)
	}
	if error != nil {
		return errors.Wrap(error, "[MysqlHandler.Delete]Error when deleting. sql - "+sqlStr)
	}
	resultRows, _ := result.RowsAffected()
	log.Debug("Deleted rows:"+strconv.FormatInt(resultRows, 10), "db")
	return nil
}

var dbObject = MysqlHandler{}

func DBHanlder() MysqlHandler {
	return dbObject
}

// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package contenttype

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xc/digimaker/core/db"
	. "github.com/xc/digimaker/core/db"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Relation is an object representing the database table.
// Implement dm.contenttype.ContentTyper interface
type Relation struct {
	ID            int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	ToContentID   int    `boil:"to_content_id" json:"to_content_id" toml:"to_content_id" yaml:"to_content_id"`
	ToType        string `boil:"to_type" json:"to_type" toml:"to_type" yaml:"to_type"`
	FromContentID int    `boil:"from_content_id" json:"from_content_id" toml:"from_content_id" yaml:"from_content_id"`
	FromType      string `boil:"from_type" json:"from_type" toml:"from_type" yaml:"from_type"`
	FromLocation  int    `boil:"from_location" json:"from_location" toml:"from_location" yaml:"from_location"`
	Priority      int    `boil:"priority" json:"priority" toml:"priority" yaml:"priority"`
	Identifier    string `boil:"identifier" json:"identifier" toml:"identifier" yaml:"identifier"`
	Description   string `boil:"description" json:"description" toml:"description" yaml:"description"`
	Data          string `boil:"data" json:"data" toml:"data" yaml:"data"`
	UID           string `boil:"uid" json:"uid" toml:"uid" yaml:"uid"`

	R *relationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L relationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

func (c *Relation) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = c.ID
	result["to_content_id"] = c.ToContentID
	result["to_type"] = c.ToType
	result["from_content_id"] = c.FromContentID
	result["from_type"] = c.FromType
	result["from_location"] = c.FromLocation
	result["priority"] = c.Priority
	result["identifier"] = c.Identifier
	result["description"] = c.Description
	result["data"] = c.Data
	result["uid"] = c.UID
	return result
}

func (c *Relation) TableName() string {
	return "dm_relation"
}

func (c *Relation) Field(name string) interface{} {
	var result interface{}
	switch name {
	case "id", "ID":
		result = c.ID
	case "to_content_id", "ToContentID":
		result = c.ToContentID
	case "to_type", "ToType":
		result = c.ToType
	case "from_content_id", "FromContentID":
		result = c.FromContentID
	case "from_type", "FromType":
		result = c.FromType
	case "from_location", "FromLocation":
		result = c.FromLocation
	case "priority", "Priority":
		result = c.Priority
	case "identifier", "Identifier":
		result = c.Identifier
	case "description", "Description":
		result = c.Description
	case "data", "Data":
		result = c.Data
	case "uid", "UID":
		result = c.UID
	default:
	}
	return result
}

func (c Relation) Store(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	if c.ID == 0 {
		id, err := handler.Insert(c.TableName(), c.ToMap(), transaction...)
		c.ID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.ToMap(), Cond("id", c.ID), transaction...)
		return err
	}
	return nil
}

var RelationColumns = struct {
	ID            string
	ToContentID   string
	ToType        string
	FromContentID string
	FromType      string
	FromLocation  string
	Priority      string
	Identifier    string
	Description   string
	Data          string
	UID           string
}{
	ID:            "id",
	ToContentID:   "to_content_id",
	ToType:        "to_type",
	FromContentID: "from_content_id",
	FromType:      "from_type",
	FromLocation:  "from_location",
	Priority:      "priority",
	Identifier:    "identifier",
	Description:   "description",
	Data:          "data",
	UID:           "uid",
}

// RelationRels is where relationship names are stored.
var RelationRels = struct {
}{}

// relationR is where relationships are stored.
type relationR struct {
}

// NewStruct creates a new relationship struct
func (*relationR) NewStruct() *relationR {
	return &relationR{}
}

// relationL is where Load methods for each relationship are stored.
type relationL struct{}

var (
	relationColumns               = []string{"id", "to_content_id", "to_type", "from_content_id", "from_type", "from_location", "priority", "identifier", "description", "data", "uid"}
	relationColumnsWithoutDefault = []string{"to_type", "from_type", "identifier", "description", "data", "uid"}
	relationColumnsWithDefault    = []string{"id", "to_content_id", "from_content_id", "from_location", "priority"}
	relationPrimaryKeyColumns     = []string{"id"}
)

type (
	// RelationSlice is an alias for a slice of pointers to Relation.
	// This should generally be used opposed to []Relation.
	RelationSlice []*Relation
	// RelationHook is the signature for custom Relation hook methods
	RelationHook func(context.Context, boil.ContextExecutor, *Relation) error

	relationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	relationType                 = reflect.TypeOf(&Relation{})
	relationMapping              = queries.MakeStructMapping(relationType)
	relationPrimaryKeyMapping, _ = queries.BindMapping(relationType, relationMapping, relationPrimaryKeyColumns)
	relationInsertCacheMut       sync.RWMutex
	relationInsertCache          = make(map[string]insertCache)
	relationUpdateCacheMut       sync.RWMutex
	relationUpdateCache          = make(map[string]updateCache)
	relationUpsertCacheMut       sync.RWMutex
	relationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var relationBeforeInsertHooks []RelationHook
var relationBeforeUpdateHooks []RelationHook
var relationBeforeDeleteHooks []RelationHook
var relationBeforeUpsertHooks []RelationHook

var relationAfterInsertHooks []RelationHook
var relationAfterSelectHooks []RelationHook
var relationAfterUpdateHooks []RelationHook
var relationAfterDeleteHooks []RelationHook
var relationAfterUpsertHooks []RelationHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Relation) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range relationBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Relation) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range relationBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Relation) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range relationBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Relation) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range relationBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Relation) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range relationAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Relation) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range relationAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Relation) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range relationAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Relation) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range relationAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Relation) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range relationAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddRelationHook registers your hook function for all future operations.
func AddRelationHook(hookPoint boil.HookPoint, relationHook RelationHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		relationBeforeInsertHooks = append(relationBeforeInsertHooks, relationHook)
	case boil.BeforeUpdateHook:
		relationBeforeUpdateHooks = append(relationBeforeUpdateHooks, relationHook)
	case boil.BeforeDeleteHook:
		relationBeforeDeleteHooks = append(relationBeforeDeleteHooks, relationHook)
	case boil.BeforeUpsertHook:
		relationBeforeUpsertHooks = append(relationBeforeUpsertHooks, relationHook)
	case boil.AfterInsertHook:
		relationAfterInsertHooks = append(relationAfterInsertHooks, relationHook)
	case boil.AfterSelectHook:
		relationAfterSelectHooks = append(relationAfterSelectHooks, relationHook)
	case boil.AfterUpdateHook:
		relationAfterUpdateHooks = append(relationAfterUpdateHooks, relationHook)
	case boil.AfterDeleteHook:
		relationAfterDeleteHooks = append(relationAfterDeleteHooks, relationHook)
	case boil.AfterUpsertHook:
		relationAfterUpsertHooks = append(relationAfterUpsertHooks, relationHook)
	}
}

// One returns a single relation record from the query.
func (q relationQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Relation, error) {
	o := &Relation{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "contenttype: failed to execute a one query for dm_relation")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Relation records from the query.
func (q relationQuery) All(ctx context.Context, exec boil.ContextExecutor) (RelationSlice, error) {
	var o []*Relation

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "contenttype: failed to assign all query results to Relation slice")
	}

	if len(relationAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Relation records in the query.
func (q relationQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "contenttype: failed to count dm_relation rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q relationQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "contenttype: failed to check if dm_relation exists")
	}

	return count > 0, nil
}

var mySQLRelationUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Relation) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("contenttype: no dm_relation provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(relationColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLRelationUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	relationUpsertCacheMut.RLock()
	cache, cached := relationUpsertCache[key]
	relationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			relationColumns,
			relationColumnsWithDefault,
			relationColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			relationColumns,
			relationPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("contenttype: unable to upsert dm_relation, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "dm_relation", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `dm_relation` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(relationType, relationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(relationType, relationMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "contenttype: unable to upsert for dm_relation")
	}

	var lastID int64
	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == relationMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(relationType, relationMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "contenttype: unable to retrieve unique values for dm_relation")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}

	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "contenttype: unable to populate default values for dm_relation")
	}

CacheNoHooks:
	if !cached {
		relationUpsertCacheMut.Lock()
		relationUpsertCache[key] = cache
		relationUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

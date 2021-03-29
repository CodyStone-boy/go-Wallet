// Code generated by SQLBoiler 4.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// AuthFullpubkey is an object representing the database table.
type AuthFullpubkey struct {
	ID            int16     `boil:"id" json:"id" toml:"id" yaml:"id"`
	Coin          string    `boil:"coin" json:"coin" toml:"coin" yaml:"coin"`
	AuthAccount   string    `boil:"auth_account" json:"auth_account" toml:"auth_account" yaml:"auth_account"`
	FullPublicKey string    `boil:"full_public_key" json:"full_public_key" toml:"full_public_key" yaml:"full_public_key"`
	UpdatedAt     null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`

	R *authFullpubkeyR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L authFullpubkeyL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AuthFullpubkeyColumns = struct {
	ID            string
	Coin          string
	AuthAccount   string
	FullPublicKey string
	UpdatedAt     string
}{
	ID:            "id",
	Coin:          "coin",
	AuthAccount:   "auth_account",
	FullPublicKey: "full_public_key",
	UpdatedAt:     "updated_at",
}

// Generated where

var AuthFullpubkeyWhere = struct {
	ID            whereHelperint16
	Coin          whereHelperstring
	AuthAccount   whereHelperstring
	FullPublicKey whereHelperstring
	UpdatedAt     whereHelpernull_Time
}{
	ID:            whereHelperint16{field: "`auth_fullpubkey`.`id`"},
	Coin:          whereHelperstring{field: "`auth_fullpubkey`.`coin`"},
	AuthAccount:   whereHelperstring{field: "`auth_fullpubkey`.`auth_account`"},
	FullPublicKey: whereHelperstring{field: "`auth_fullpubkey`.`full_public_key`"},
	UpdatedAt:     whereHelpernull_Time{field: "`auth_fullpubkey`.`updated_at`"},
}

// AuthFullpubkeyRels is where relationship names are stored.
var AuthFullpubkeyRels = struct {
}{}

// authFullpubkeyR is where relationships are stored.
type authFullpubkeyR struct {
}

// NewStruct creates a new relationship struct
func (*authFullpubkeyR) NewStruct() *authFullpubkeyR {
	return &authFullpubkeyR{}
}

// authFullpubkeyL is where Load methods for each relationship are stored.
type authFullpubkeyL struct{}

var (
	authFullpubkeyAllColumns            = []string{"id", "coin", "auth_account", "full_public_key", "updated_at"}
	authFullpubkeyColumnsWithoutDefault = []string{"coin", "auth_account", "full_public_key"}
	authFullpubkeyColumnsWithDefault    = []string{"id", "updated_at"}
	authFullpubkeyPrimaryKeyColumns     = []string{"id"}
)

type (
	// AuthFullpubkeySlice is an alias for a slice of pointers to AuthFullpubkey.
	// This should generally be used opposed to []AuthFullpubkey.
	AuthFullpubkeySlice []*AuthFullpubkey

	authFullpubkeyQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	authFullpubkeyType                 = reflect.TypeOf(&AuthFullpubkey{})
	authFullpubkeyMapping              = queries.MakeStructMapping(authFullpubkeyType)
	authFullpubkeyPrimaryKeyMapping, _ = queries.BindMapping(authFullpubkeyType, authFullpubkeyMapping, authFullpubkeyPrimaryKeyColumns)
	authFullpubkeyInsertCacheMut       sync.RWMutex
	authFullpubkeyInsertCache          = make(map[string]insertCache)
	authFullpubkeyUpdateCacheMut       sync.RWMutex
	authFullpubkeyUpdateCache          = make(map[string]updateCache)
	authFullpubkeyUpsertCacheMut       sync.RWMutex
	authFullpubkeyUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single authFullpubkey record from the query.
func (q authFullpubkeyQuery) One(ctx context.Context, exec boil.ContextExecutor) (*AuthFullpubkey, error) {
	o := &AuthFullpubkey{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for auth_fullpubkey")
	}

	return o, nil
}

// All returns all AuthFullpubkey records from the query.
func (q authFullpubkeyQuery) All(ctx context.Context, exec boil.ContextExecutor) (AuthFullpubkeySlice, error) {
	var o []*AuthFullpubkey

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to AuthFullpubkey slice")
	}

	return o, nil
}

// Count returns the count of all AuthFullpubkey records in the query.
func (q authFullpubkeyQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count auth_fullpubkey rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q authFullpubkeyQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if auth_fullpubkey exists")
	}

	return count > 0, nil
}

// AuthFullpubkeys retrieves all the records using an executor.
func AuthFullpubkeys(mods ...qm.QueryMod) authFullpubkeyQuery {
	mods = append(mods, qm.From("`auth_fullpubkey`"))
	return authFullpubkeyQuery{NewQuery(mods...)}
}

// FindAuthFullpubkey retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAuthFullpubkey(ctx context.Context, exec boil.ContextExecutor, iD int16, selectCols ...string) (*AuthFullpubkey, error) {
	authFullpubkeyObj := &AuthFullpubkey{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `auth_fullpubkey` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, authFullpubkeyObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from auth_fullpubkey")
	}

	return authFullpubkeyObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *AuthFullpubkey) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no auth_fullpubkey provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.UpdatedAt).IsZero() {
			queries.SetScanner(&o.UpdatedAt, currTime)
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(authFullpubkeyColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	authFullpubkeyInsertCacheMut.RLock()
	cache, cached := authFullpubkeyInsertCache[key]
	authFullpubkeyInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			authFullpubkeyAllColumns,
			authFullpubkeyColumnsWithDefault,
			authFullpubkeyColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(authFullpubkeyType, authFullpubkeyMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(authFullpubkeyType, authFullpubkeyMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `auth_fullpubkey` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `auth_fullpubkey` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `auth_fullpubkey` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, authFullpubkeyPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into auth_fullpubkey")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int16(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == authFullpubkeyMapping["id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for auth_fullpubkey")
	}

CacheNoHooks:
	if !cached {
		authFullpubkeyInsertCacheMut.Lock()
		authFullpubkeyInsertCache[key] = cache
		authFullpubkeyInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the AuthFullpubkey.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *AuthFullpubkey) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	var err error
	key := makeCacheKey(columns, nil)
	authFullpubkeyUpdateCacheMut.RLock()
	cache, cached := authFullpubkeyUpdateCache[key]
	authFullpubkeyUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			authFullpubkeyAllColumns,
			authFullpubkeyPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update auth_fullpubkey, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `auth_fullpubkey` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, authFullpubkeyPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(authFullpubkeyType, authFullpubkeyMapping, append(wl, authFullpubkeyPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update auth_fullpubkey row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for auth_fullpubkey")
	}

	if !cached {
		authFullpubkeyUpdateCacheMut.Lock()
		authFullpubkeyUpdateCache[key] = cache
		authFullpubkeyUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q authFullpubkeyQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for auth_fullpubkey")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for auth_fullpubkey")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AuthFullpubkeySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authFullpubkeyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `auth_fullpubkey` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, authFullpubkeyPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in authFullpubkey slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all authFullpubkey")
	}
	return rowsAff, nil
}

var mySQLAuthFullpubkeyUniqueColumns = []string{
	"id",
	"full_public_key",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *AuthFullpubkey) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no auth_fullpubkey provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	nzDefaults := queries.NonZeroDefaultSet(authFullpubkeyColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLAuthFullpubkeyUniqueColumns, o)

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

	authFullpubkeyUpsertCacheMut.RLock()
	cache, cached := authFullpubkeyUpsertCache[key]
	authFullpubkeyUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			authFullpubkeyAllColumns,
			authFullpubkeyColumnsWithDefault,
			authFullpubkeyColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			authFullpubkeyAllColumns,
			authFullpubkeyPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("models: unable to upsert auth_fullpubkey, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`auth_fullpubkey`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `auth_fullpubkey` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(authFullpubkeyType, authFullpubkeyMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(authFullpubkeyType, authFullpubkeyMapping, ret)
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

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for auth_fullpubkey")
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

	o.ID = int16(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == authFullpubkeyMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(authFullpubkeyType, authFullpubkeyMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for auth_fullpubkey")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for auth_fullpubkey")
	}

CacheNoHooks:
	if !cached {
		authFullpubkeyUpsertCacheMut.Lock()
		authFullpubkeyUpsertCache[key] = cache
		authFullpubkeyUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single AuthFullpubkey record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *AuthFullpubkey) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no AuthFullpubkey provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), authFullpubkeyPrimaryKeyMapping)
	sql := "DELETE FROM `auth_fullpubkey` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from auth_fullpubkey")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for auth_fullpubkey")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q authFullpubkeyQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no authFullpubkeyQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from auth_fullpubkey")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for auth_fullpubkey")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AuthFullpubkeySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authFullpubkeyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `auth_fullpubkey` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, authFullpubkeyPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from authFullpubkey slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for auth_fullpubkey")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *AuthFullpubkey) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindAuthFullpubkey(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AuthFullpubkeySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AuthFullpubkeySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authFullpubkeyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `auth_fullpubkey`.* FROM `auth_fullpubkey` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, authFullpubkeyPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in AuthFullpubkeySlice")
	}

	*o = slice

	return nil
}

// AuthFullpubkeyExists checks if the AuthFullpubkey row exists.
func AuthFullpubkeyExists(ctx context.Context, exec boil.ContextExecutor, iD int16) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `auth_fullpubkey` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if auth_fullpubkey exists")
	}

	return exists, nil
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o AuthFullpubkeySlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {
		if !boil.TimestampsAreSkipped(ctx) {
			currTime := time.Now().In(boil.GetLocation())

			if queries.MustTime(row.UpdatedAt).IsZero() {
				queries.SetScanner(&row.UpdatedAt, currTime)
			}
		}

		nzDefaults := queries.NonZeroDefaultSet(authFullpubkeyColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			authFullpubkeyAllColumns,
			authFullpubkeyColumnsWithDefault,
			authFullpubkeyColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `auth_fullpubkey` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(authFullpubkeyType, authFullpubkeyMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "models: unable to insert into auth_fullpubkey")
	}

	return nil
}

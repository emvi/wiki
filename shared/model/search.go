package model

import (
	"emviwiki/shared/db"
	"fmt"
	"github.com/emvi/hide"
	"strings"
	"time"
)

const (
	filterMaxLimit    = 20
	sortDirectionASC  = "ASC"
	sortDirectionDESC = "DESC"
)

// Basic information used to search/filter articles, groups, ...
// Provides methods to add the parameters to a database query.
type BaseSearch struct {
	CreatedStart time.Time `json:"created_start"` // zero time means don't filter
	CreatedEnd   time.Time `json:"created_end"`
	UpdatedStart time.Time `json:"updated_start"`
	UpdatedEnd   time.Time `json:"updated_end"`
	SortCreated  string    `json:"sort_created"` // values: asc/desc, empty means no sorting
	SortUpdated  string    `json:"sort_updated"`
	Offset       int       `json:"offset"`
	Limit        int       `json:"limit"` // <= 0 means unlimited
}

type SortValue struct {
	Field     string
	Direction string
}

// Adds the given fields to the filter using fuzzy search.
// The index is increased by one for each parameter and the new index is returned
// together with the query and field values.
func (search *BaseSearch) addFieldFilter(table string, index int, params []interface{}, values []string, fields ...string) (string, int, []interface{}) {
	table = getTableName(table)
	filterQuery := make([]string, 0)

	for i := range fields {
		values[i] = strings.TrimSpace(values[i])

		if values[i] != "" {
			filterQuery = append(filterQuery, fmt.Sprintf(`(SIMILARITY(%v"%v", $%v) > 0.2 OR LOWER(%v"%v") LIKE LOWER('%%'||$%v||'%%'))`, table, fields[i], index, table, fields[i], index))
			params = append(params, values[i])
			index++
		}
	}

	if len(filterQuery) == 0 {
		return "", index, params
	}

	return fmt.Sprintf("AND %s ", strings.Join(filterQuery, " AND ")), index, params
}

// Adds the given fields to the filter using a tsvector query.
// The index is increased by one for each parameter and the new index is returned
// together with the query and field values.
func (search *BaseSearch) addTSVectorFieldFilter(table string, index int, params []interface{}, values []string, fields ...string) (string, int, []interface{}) {
	table = getTableName(table)
	filterQuery := make([]string, 0)

	for i := range fields {
		values[i] = strings.TrimSpace(values[i])

		if values[i] != "" {
			filterQuery = append(filterQuery, fmt.Sprintf(`%v"%v" @@ to_tsquery($%v)`, table, fields[i], index))
			params = append(params, db.ToTSVector(values[i]))
			index++
		}
	}

	if len(filterQuery) == 0 {
		return "", index, params
	}

	return fmt.Sprintf("AND %s ", strings.Join(filterQuery, " AND ")), index, params
}

func getTableName(table string) string {
	if table != "" {
		return fmt.Sprintf(`"%s".`, table)
	}

	return ""
}

// Adds the def/mod time fields to the query if set.
// The index is increased by one for each field and returned
// together with the query and field values.
func (search *BaseSearch) addDateFilter(table string, index int, params []interface{}) (string, int, []interface{}) {
	query := ""

	if !search.CreatedStart.IsZero() {
		query += fmt.Sprintf(`AND "%v".def_time > $%v `, table, index)
		params = append(params, search.CreatedStart)
		index++
	}

	if !search.CreatedEnd.IsZero() {
		query += fmt.Sprintf(`AND "%v".def_time-INTERVAL '1 DAY' < $%v `, table, index)
		params = append(params, search.CreatedEnd)
		index++
	}

	if !search.UpdatedStart.IsZero() {
		query += fmt.Sprintf(`AND "%v".mod_time > $%v `, table, index)
		params = append(params, search.UpdatedStart)
		index++
	}

	if !search.UpdatedEnd.IsZero() {
		query += fmt.Sprintf(`AND "%v".mod_time-INTERVAL '1 DAY' < $%v `, table, index)
		params = append(params, search.UpdatedEnd)
		index++
	}

	return query, index, params
}

// Returns all columns used to order the results.
func (search *BaseSearch) getSortColumns(table string, defaultFields []SortValue, fields ...SortValue) string {
	sortFields := search.getSortFields(table, fields, defaultFields)

	if len(sortFields) == 0 {
		return ""
	}

	sortQueries := make([]string, 0, len(sortFields))

	for _, field := range sortFields {
		sortQueries = append(sortQueries, field.Field)
	}

	return strings.Join(sortQueries, ",")
}

// Adds sorting to the query. The values must be an empty string, asc or desc.
func (search *BaseSearch) addSorting(table string, defaultFields []SortValue, fields ...SortValue) string {
	sortFields := search.getSortFields(table, fields, defaultFields)

	if len(sortFields) == 0 {
		return ""
	}

	sortQueries := make([]string, 0, len(sortFields))

	for _, field := range sortFields {
		sortQueries = append(sortQueries, fmt.Sprintf("%v %v", field.Field, field.Direction))
	}

	return fmt.Sprintf("ORDER BY %s ", strings.Join(sortQueries, ","))
}

func (search *BaseSearch) getSortFields(table string, customFields, defaultFields []SortValue) []SortValue {
	sortFields := make([]SortValue, 0)
	sortFields = search.addCustomSortFields(sortFields, customFields, defaultFields)
	sortFields = search.addModDefTimeSortFields(table, sortFields)
	sortFields = append(sortFields, defaultFields...)
	return sortFields
}

func (search *BaseSearch) addCustomSortFields(sortFields, customFields, defaultFields []SortValue) []SortValue {
	for _, field := range customFields {
		value := strings.ToUpper(strings.TrimSpace(field.Direction))

		if search.isValidSortParam(value) {
			defaultFieldIndex := search.indexOfSortValue(defaultFields, field.Field)

			if defaultFieldIndex == -1 {
				sortFields = append(sortFields, SortValue{field.Field, value})
			} else {
				defaultFields[defaultFieldIndex].Direction = value
			}
		}
	}

	return sortFields
}

func (search *BaseSearch) addModDefTimeSortFields(table string, sortFields []SortValue) []SortValue {
	search.SortCreated = strings.ToUpper(strings.TrimSpace(search.SortCreated))
	search.SortUpdated = strings.ToUpper(strings.TrimSpace(search.SortUpdated))

	if search.isValidSortParam(search.SortCreated) {
		sortFields = append(sortFields, SortValue{fmt.Sprintf(`"%v"."def_time"`, table), search.SortCreated})
	}

	if search.isValidSortParam(search.SortUpdated) {
		sortFields = append(sortFields, SortValue{fmt.Sprintf(`"%v"."mod_time"`, table), search.SortUpdated})
	}

	return sortFields
}

func (search *BaseSearch) indexOfSortValue(haystack []SortValue, needle string) int {
	for i, value := range haystack {
		if value.Field == needle {
			return i
		}
	}

	return -1
}

func (search *BaseSearch) isValidSortParam(param string) bool {
	return param != "" && (param == "ASC" || param == "DESC")
}

// Adds the limit and offset to the query if set.
// The index is increased by two if set and returned
// together with the query and field values.
func (search *BaseSearch) addLimit(index int, params []interface{}) (string, int, []interface{}) {
	query := ""

	if search.Limit > filterMaxLimit || search.Limit == 0 {
		search.Limit = filterMaxLimit
	} else if search.Limit < 0 {
		search.Limit = 0
	}

	if search.Offset < 0 {
		search.Offset = 0
	}

	if search.Limit > 0 && search.Offset >= 0 {
		query = fmt.Sprintf("LIMIT $%v OFFSET $%v ", index, index+1)
		params = append(params, search.Limit, search.Offset)
		index += 2
	}

	return query, index, params
}

func (search *BaseSearch) filterInIds(ids []hide.ID, index int, params []interface{}) (string, int, []interface{}) {
	query := ` IN (`

	for i := range ids {
		query += fmt.Sprintf(`$%v,`, index)
		params = append(params, ids[i])
		index++
	}

	query = query[:len(query)-1]
	query += ") "

	return query, index, params
}

func (search *BaseSearch) filterInStrings(strs []string, index int, params []interface{}) (string, int, []interface{}) {
	query := ` IN (`

	for i := range strs {
		query += fmt.Sprintf(`$%v,`, index)
		params = append(params, strs[i])
		index++
	}

	query = query[:len(query)-1]
	query += ") "

	return query, index, params
}

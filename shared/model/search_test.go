package model

import (
	"github.com/emvi/hide"
	"testing"
)

func TestAddFieldFilter(t *testing.T) {
	base := BaseSearch{}
	params := make([]interface{}, 0)
	values := []string{"val1", "val2"}
	fields := []string{"field1", "field2"}
	query, index, params := base.addFieldFilter("table", 1, params, values, fields...)

	if query != `AND (SIMILARITY("table"."field1", $1) > 0.2 OR LOWER("table"."field1") LIKE LOWER('%'||$1||'%')) AND (SIMILARITY("table"."field2", $2) > 0.2 OR LOWER("table"."field2") LIKE LOWER('%'||$2||'%')) ` {
		t.Fatalf("Query not as expected: %v", query)
	}

	if index != 3 {
		t.Fatalf("Index not as expected: %v", index)
	}

	if len(params) != 2 {
		t.Fatalf("Params not as expected: %v", params)
	}
}

func TestAddSorting(t *testing.T) {
	base := BaseSearch{}
	input := []struct {
		defaultFields []SortValue
		fields        []SortValue
	}{
		{nil, []SortValue{{"a", ""}}},
		{nil, []SortValue{{"a", "asc"}, {"b", "desc"}}},
		{[]SortValue{}, []SortValue{{"a", "asc"}, {"b", "desc"}}},
		{[]SortValue{{"default", sortDirectionASC}}, nil},
		{[]SortValue{{"default", sortDirectionASC}}, []SortValue{{"a", "asc"}, {"b", "desc"}}},
		{[]SortValue{{"default", sortDirectionDESC}}, []SortValue{{"default", "asc"}}},
		{[]SortValue{{"default", sortDirectionDESC}, {"default2", sortDirectionASC}}, []SortValue{{"default", "asc"}, {"test", "desc"}}},
	}
	expected := []string{
		``,
		`ORDER BY a ASC,b DESC `,
		`ORDER BY a ASC,b DESC `,
		`ORDER BY default ASC `,
		`ORDER BY a ASC,b DESC,default ASC `,
		`ORDER BY default ASC `,
		`ORDER BY test DESC,default ASC,default2 ASC `,
	}

	for i, in := range input {
		if res := base.addSorting("test", in.defaultFields, in.fields...); res != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], res)
		}
	}
}

func TestGetSortColumns(t *testing.T) {
	base := BaseSearch{}
	input := []struct {
		defaultFields []SortValue
		fields        []SortValue
	}{
		{nil, []SortValue{{"a", ""}}},
		{nil, []SortValue{{"a", "asc"}}},
		{nil, []SortValue{{"a", "asc"}, {"b", "desc"}}},
		{[]SortValue{{"default", sortDirectionASC}}, nil},
		{[]SortValue{{"default", sortDirectionASC}}, []SortValue{{"a", "asc"}, {"b", "desc"}}},
		{[]SortValue{{"default", sortDirectionDESC}}, []SortValue{{"default", "asc"}}},
		{[]SortValue{{"default", sortDirectionDESC}, {"default2", sortDirectionASC}}, []SortValue{{"default", "asc"}, {"test", "desc"}}},
	}
	expected := []string{
		``,
		`a`,
		`a,b`,
		`default`,
		`a,b,default`,
		`default`,
		`test,default,default2`,
	}

	for i, in := range input {
		if res := base.getSortColumns("test", in.defaultFields, in.fields...); res != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], res)
		}
	}
}

func TestFilterInIds(t *testing.T) {
	base := BaseSearch{}
	index := 1
	params := make([]interface{}, 0)
	var query string
	ids := []hide.ID{1, 2, 3}

	query, index, params = base.filterInIds(ids, index, params)

	if query != " IN ($1,$2,$3) " {
		t.Fatalf("Unexpected query: %v", query)
	}

	if index != 4 || len(params) != 3 {
		t.Fatalf("Expected index to be 4 and param length to be 3, but was: %v %v", index, len(params))
	}
}

func TestFilterInStrings(t *testing.T) {
	base := BaseSearch{}
	index := 1
	params := make([]interface{}, 0)
	var query string
	strs := []string{"one", "two", "three"}

	query, index, params = base.filterInStrings(strs, index, params)

	if query != " IN ($1,$2,$3) " {
		t.Fatalf("Unexpected query: %v", query)
	}

	if index != 4 || len(params) != 3 {
		t.Fatalf("Expected index to be 4 and param length to be 3, but was: %v %v", index, len(params))
	}
}

// Package util
// @author Daud Valentino
package util

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
)

const (
	layoutDateTimeFormat = `2006-01-02 15:04:05`
)

type KeyValue struct {
	Key   string
	Value interface{}
}

// StructToMap converts a struct to a map using the struct's tags.
// StructToMap uses tags on struct fields to decide which fields to add to the
// returned map.
func StructToMap(src interface{}, tag string) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return out, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		//field := reflectValue.Field(i).Interface()
		if !fi.IsExported() {
			continue
		}

		tagsv := strings.Split(fi.Tag.Get(tag), ",")

		if tagsv[0] != "" && fi.PkgPath == "" {

			if isNil(v.Field(i).Interface()) {
				continue
			}

			// skip if omitempty
			if (len(tagsv) > 1 && tagsv[1] == "omitempty") && IsEmptyValue(v.Field(i).Interface()) {
				continue
			}

			if isTime(v.Field(i)) {
				if timeIsZero(v.Field(i)) && (len(tagsv) > 1 && tagsv[1] == "omitempty") {
					continue
				}
			}

			if v.Field(i).Kind() == reflect.Struct {
				continue
			}

			col := tagsv[0]

			if InArray("ne", tagsv) {
				col = fmt.Sprintf("%s !", col)
			}
			// set key value of map interface output
			out[col] = v.Field(i).Interface()
		}

		if tagsv[0] == "" && v.Field(i).Kind() == reflect.Struct {
			x, err := StructToMap(v.Field(i).Interface(), tag)
			if err != nil {
				return out, err
			}

			for y, z := range x {
				out[y] = z
			}
		}
	}

	return out, nil
}

// ToColumnsValues iterate struct to separate key field and value
func ToColumnsValues(src interface{}, tag string) ([]string, []interface{}, error) {
	var columns []string
	var values []interface{}

	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		tagsv := strings.Split(fi.Tag.Get(tag), ",")

		if tagsv[0] != "" && fi.PkgPath == "" {

			if tagsv[0] == "-" {
				continue
			}

			if isNil(v.Field(i).Interface()) {
				continue
			}

			// skip if omitempty
			if (len(tagsv) > 1 && InArray("omitempty", tagsv)) && IsEmptyValue(v.Field(i).Interface()) {
				continue
			}

			if isTime(v.Field(i)) {
				if timeIsZero(v.Field(i)) && (len(tagsv) > 1 && InArray("omitempty", tagsv)) {
					continue
				}
			}

			col := tagsv[0]

			if InArray("ne", tagsv) {
				col = fmt.Sprintf("%s !", col)
			}

			// set value of string slice to value in struct field
			columns = append(columns, col)

			// set value interface of value struct field
			values = append(values, v.Field(i).Interface())

		}
	}

	return columns, values, nil
}

func isTime(obj reflect.Value) bool {
	_, ok := obj.Interface().(time.Time)
	if ok {
		return ok
	}

	_, ok = obj.Interface().(*time.Time)

	return ok
}

func timeIsZero(obj reflect.Value) bool {
	t, ok := obj.Interface().(time.Time)
	if ok {
		return t.IsZero()
	}

	t2, ok := obj.Interface().(*time.Time)
	if ok {
		return false
	}

	return t2 == nil
}

func isNil(i interface{}) bool {
	if i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil()) {
		return true
	}

	return false
}

func SliceStructToBulkInsert(src interface{}, tag string) ([]string, []interface{}, []string, error) {
	var columns []string
	var replacers []string
	var values []interface{}

	v := reflect.Indirect(reflect.ValueOf(src))
	t := reflect.TypeOf(src)
	if t.Kind() == reflect.Ptr {
		//v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return columns, values, replacers, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), t.Kind().String())
	}

	for i := 0; i < v.Len(); i++ {

		item := v.Index(i)
		if !item.IsValid() {
			continue
		}

		cols, val, err := ToColumnsValues(item.Interface(), tag)
		if err != nil {
			return columns, values, replacers, err
		}

		if len(columns) == 0 {
			columns = cols
		}

		pattern := fmt.Sprintf(`(%s)`, strings.TrimRight(strings.Repeat("?,", len(columns)), `,`))
		replacers = append(replacers, pattern)
		values = append(values, val...)
	}

	return columns, values, replacers, nil
}

// StructToQueryUpdate struct to query update builder
// this method doesn't support nested struct
func StructToQueryUpdate(input interface{}, where interface{}, tableName, tag string) (string, []interface{}, error) {

	cols, vals, err := ToColumnsValues(input, tag)
	if err != nil {
		return "", vals, err
	}

	cu, vu, err := ToColumnsValues(where, tag)
	if err != nil {
		return "", vals, err
	}

	q := fmt.Sprintf(`UPDATE %s SET %s`, tableName, StringJoin(cols, "=?, ", "=?"))
	if len(cu) > 0 {
		q = fmt.Sprintf(`%s WHERE %s`, q, StringJoin(cu, "=? AND ", "=?"))
		vals = append(vals, vu...)
	}

	return q, vals, nil
}

// StructToKeyValue converts a struct to a key value the struct's tags.
// StructToKeyValue uses tags on struct fields to decide which fields to add to the
// returned slice struct.
func StructToKeyValue(src interface{}, tag string) ([]KeyValue, error) {
	var out []KeyValue
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return out, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		//field := reflectValue.Field(i).Interface()
		if !fi.IsExported() {
			continue
		}

		tagsv := strings.Split(fi.Tag.Get(tag), ",")

		if tagsv[0] != "" && fi.PkgPath == "" {

			if isNil(v.Field(i).Interface()) {
				continue
			}

			// skip if omitempty
			if (len(tagsv) > 1 && tagsv[1] == "omitempty") && IsEmptyValue(v.Field(i).Interface()) {
				continue
			}

			if isTime(v.Field(i)) {
				if timeIsZero(v.Field(i)) && (len(tagsv) > 1 && tagsv[1] == "omitempty") {
					continue
				}
			}

			if v.Field(i).Kind() == reflect.Struct {
				continue
			}

			col := tagsv[0]

			if InArray("ne", tagsv) {
				col = fmt.Sprintf("%s !", col)
			}
			// set key value of struct key value interface output
			out = append(out, KeyValue{
				Key:   col,
				Value: v.Field(i).Interface(),
			})
		}

		if tagsv[0] == "" && v.Field(i).Kind() == reflect.Struct {
			x, err := StructToKeyValue(v.Field(i).Interface(), tag)
			if err != nil {
				return out, err
			}

			out = append(out, x...)
		}
	}

	return out, nil
}

// StructToKeyValueWithSkipOmitEmpty converts a struct to a key value the struct's tags.
// StructToKeyValueWithSkipOmitEmpty uses tags on struct fields to decide which fields to add to the
// returned slice struct.
func StructToKeyValueWithSkipOmitEmpty(src interface{}, tag string, columns []string, skipOmitEmpty bool) ([]KeyValue, error) {
	var out []KeyValue
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return out, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		//field := reflectValue.Field(i).Interface()
		if !fi.IsExported() {
			continue
		}

		tagsv := strings.Split(fi.Tag.Get(tag), ",")
		col := tagsv[0]
		if col != "" && fi.PkgPath == "" {

			if isNil(v.Field(i).Interface()) {
				continue
			}

			if !InArray(col, columns) {
				continue
			}

			// skip if omitempty
			if (len(tagsv) > 1 && tagsv[1] == "omitempty") && IsEmptyValue(v.Field(i).Interface()) && skipOmitEmpty {
				continue
			}

			if isTime(v.Field(i)) {
				if timeIsZero(v.Field(i)) && (len(tagsv) > 1 && tagsv[1] == "omitempty") && skipOmitEmpty {
					continue
				}
			}

			if v.Field(i).Kind() == reflect.Struct {
				continue
			}

			// set key value of struct key value interface output
			out = append(out, KeyValue{
				Key:   col,
				Value: v.Field(i).Interface(),
			})
		}

		if col == "" && v.Field(i).Kind() == reflect.Struct {
			x, err := StructToKeyValue(v.Field(i).Interface(), tag)
			if err != nil {
				return out, err
			}

			out = append(out, x...)
		}
	}

	return out, nil
}



// StructQueryWhere create query builder from struct
func StructQueryWhere(iStruct interface{}, hideDeleted bool, tag string) (q string, vals []interface{}, limit, page uint64, err error) {

	var cols []string
	var startDate, endDate, sortOrder, groupBy string

	if iStruct == nil {
		return q, vals, limit, page, nil
	}

	data, err := StructToKeyValue(iStruct, tag)
	if err != nil {
		return q, vals, limit, page, err
	}

	if len(data) == 0 {
		return q, vals, limit, page, err
	}

	for i := 0; i < len(data); i++ {
		if data[i].Key == "page" {
			page = cast.ToUint64(data[i].Value)
			continue
		}

		if data[i].Key == "limit" {
			limit = cast.ToUint64(data[i].Value)
			continue
		}

		if data[i].Key == "start_date" {
			startDate = cast.ToString(data[i].Value)
			continue
		}

		if data[i].Key == "end_date" {
			endDate = cast.ToString(data[i].Value)
			continue
		}

		if data[i].Key == "sort_order" {
			sortOrder = cast.ToString(data[i].Value)
			continue
		}

		if data[i].Key == "group_by" {
			groupBy = cast.ToString(data[i].Value)
			continue
		}

		vals = append(vals, data[i].Value)
		cols = append(cols, data[i].Key)
	}

	if len(cols) > 0 && !hideDeleted {
		q = fmt.Sprintf(`WHERE %s`, StringJoin(cols, "=? AND ", "=?"))
	}

	if len(cols) > 0 && hideDeleted {
		q = fmt.Sprintf(`WHERE %s AND deleted_at <= '1970-01-02 00:00:01'`, StringJoin(cols, "=? AND ", "=?"))
	}

	if len(cols) == 0 && hideDeleted {
		q = fmt.Sprint(`WHERE deleted_at <= '1970-01-02 00:00:01'`)
	}

	if startDate != "" && endDate != "" {
		q = fmt.Sprintf(`%s AND ( created_at >= ?  AND  created_at <= ? )`, q)
		if hideDeleted {
			q = fmt.Sprintf(`%s AND ( created_at >= ?  AND created_at  <= ? ) AND deleted_at <= '1970-01-02 00:00:01'`, q)
		}

		if len(cols) == 0 && !hideDeleted {
			q = fmt.Sprint(`WHERE (created_at >= ?  AND created_at <= ? )`)
		}

		if len(cols) == 0 && hideDeleted {
			q = fmt.Sprint(`WHERE (created_at >= ?  AND created_at <= ? ) AND deleted_at <= '1970-01-02 00:00:01'`)
		}

		vals = append(vals, startDate, endDate)
	}

	if groupBy != "" {
		q = fmt.Sprintf("%s GROUP BY %s", q, groupBy)
	}

	if sortOrder != "" {
		q = fmt.Sprintf("%s ORDER BY created_at %s", q, sortOrder)
	}

	return q, vals, limit, page, err
}

// StructToQueryInsert struct to query insert builder
// this method doesn't support nested struct
func StructToQueryInsert(s interface{}, tableName, tag string) (string, []interface{}, error) {

	cols, vals, err := ToColumnsValues(s, tag)
	if err != nil {
		return "", vals, err
	}

	q := `INSERT INTO %s (%s) VALUES(%s)`

	if len(cols) == 0 {
		return "", nil, fmt.Errorf("no column available")
	}

	pattern := "?" + strings.Repeat(", ?", len(cols)-1)

	q = fmt.Sprintf(q, tableName, strings.Join(cols, ","), pattern)

	return q, vals, nil
}

// StructToQueryDelete struct to query delete builder
// this method doesn't support nested struct
func StructToQueryDelete(where interface{}, tableName, tag string, soft bool) (string, []interface{}, error) {

	cols, vals, err := ToColumnsValues(where, tag)
	if err != nil {
		return "", vals, err
	}

	q := fmt.Sprintf(`DELETE FROM %s `, tableName)

	if soft {
		fmt.Sprintf(`UPDATE %s SET deleted_at = ?`, tableName)
		vals = append([]interface{}{time.Now().Format(layoutDateTimeFormat)}, vals...)
	}

	if len(cols) > 0 {
		q = fmt.Sprintf(`%s WHERE %s`, q, StringJoin(cols, "=?, ", "=?"))
	}

	return q, vals, nil
}

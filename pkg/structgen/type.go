// Package structgen
// @author Daud Valentino
package structgen

import "errors"

func layerPresentationType(col *ColumnSchema) (string, bool, error) {
	var gt string

	switch col.DataType {
	case "char", "varchar", "enum", "set", "text", "longtext", "mediumtext", "tinytext":
		gt = "string"
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		gt = "[]byte"
	case "date", "time", "datetime", "timestamp":
		gt = "string"
	case "bit", "tinyint":
		gt = "int8"
		if col.IsNullable == "YES" {
			gt = "*int8"
		}
	case "smallint":
		gt = "int16"
		if col.IsNullable == "YES" {
			gt = "*int16"
		}
	case "mediumint":
		gt = "int32"
		if col.IsNullable == "YES" {
			gt = "*int32"
		}
	case "int":
		gt = "int"
		if col.IsNullable == "YES" {
			gt = "*int"
		}

	case "bigint":
		gt = "int64"
		if col.IsNullable == "YES" {
			gt = "*int64"
		}

	case "float", "decimal", "double":
		gt = "float64"
		if col.IsNullable == "YES" {
			gt = "*float64"
		}

	case "year":
		gt = "int"
		if col.IsNullable == "YES" {
			gt = "*int"
		}
	}

	if gt == "" {
		n := col.TableName + "." + col.ColumnName
		return "", false, errors.New("No compatible datatype (" + col.DataType + ") for " + n + " found")
	}

	return gt, col.IsNullable == "YES", nil
}

// Package structgen
// @author Daud Valentino
package structgen

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"sharefood/pkg/util"
)

var defaults = Configuration{
	DbUser:     "db_user",
	DbPassword: "db_pw",
	DbName:     "bd_name",
	PkgName:    "DbStructs",
	TagLabel:   "db",
}
var (
	flags = flag.NewFlagSet("gen", flag.ExitOnError)
)

const (
	entityPath       = "internal/entity"
	repoPath         = "internal/repositories"
	presentationPath = "internal/presentations"
	uCasePath        = "internal/ucase"
)

var config Configuration

type Configuration struct {
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbHost     string `json:"db_host"`
	DbName     string `json:"db_name"`
	// PkgName gives name of the package using the stucts
	PkgName string `json:"pkg_name"`
	// TagLabel produces tags commonly used to match database field names with Go struct members
	TagLabel string `json:"tag_label"`
}

type ColumnSchema struct {
	TableName              string
	ColumnName             string
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	ColumnKey              string
}

type ParserTemplate struct {
	TableName  string
	StructName string
	ObjectName string
	EntityName string
	FileName   string
	Query      string
}

type UseCaseTemplate struct {
	TableName        string
	StructName       string
	PackageName      string
	EntityName       string
	FileName         string
	RepoContractName string
}

var configFile = flag.String("json", "", "Config file")

func getSchema(tableName string) []ColumnSchema {
	conn, err := sql.Open("mysql", config.DbUser+":"+config.DbPassword+"@tcp("+config.DbHost+")/information_schema")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, " +
		"CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, COLUMN_TYPE, " +
		"COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY TABLE_NAME, ORDINAL_POSITION"
	rows, err := conn.Query(q, config.DbName, tableName)
	if err != nil {
		log.Fatal(err)
	}

	columns := []ColumnSchema{}
	for rows.Next() {
		cs := ColumnSchema{}
		err := rows.Scan(&cs.TableName, &cs.ColumnName, &cs.IsNullable, &cs.DataType,
			&cs.CharacterMaximumLength, &cs.NumericPrecision, &cs.NumericScale,
			&cs.ColumnType, &cs.ColumnKey)
		if err != nil {
			log.Fatal(err)
		}
		columns = append(columns, cs)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return columns
}

func formatName(name string) string {
	parts := strings.Split(name, "_")
	newName := ""
	for _, p := range parts {
		if len(p) < 1 {
			continue
		}
		newName = newName + strings.Replace(p, string(p[0]), strings.ToUpper(string(p[0])), 1)
	}

	up := map[string]string{
		"Id":  "ID",
		"Qr":  "QR",
		"Fcm": "FCM",
	}

	for k, v := range up {
		if strings.HasSuffix(newName, k) {
			newName = fmt.Sprintf("%s%s", strings.TrimSuffix(newName, k), v)

		}
	}

	return newName
}

func goType(col *ColumnSchema) (string, string, error) {
	requiredImport := ""
	if col.IsNullable == "YES" {
		requiredImport = "database/sql"
	}
	var gt string = ""
	switch col.DataType {
	case "char", "varchar", "enum", "set", "text", "longtext", "mediumtext", "tinytext":
		if col.IsNullable == "YES" {
			gt = "sql.NullString"
		} else {
			gt = "string"
		}
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		gt = "[]byte"
	case "date", "time", "datetime", "timestamp":
		gt, requiredImport = "time.Time", "time"
	case "bit", "tinyint":
		if col.IsNullable == "YES" {
			gt = "sql.NullInt64"
		} else {
			gt = "int8"
		}
	case "smallint":
		if col.IsNullable == "YES" {
			gt = "sql.NullInt64"
		} else {
			gt = "int16"
		}
	case "mediumint":
		if col.IsNullable == "YES" {
			gt = "sql.NullInt64"
		} else {
			gt = "int32"
		}
	case "int":
		if col.IsNullable == "YES" {
			gt = "sql.NullInt64"
		} else {
			gt = "int"
		}

	case "bigint":
		if col.IsNullable == "YES" {
			gt = "sql.NullInt64"
		} else {
			gt = "int64"
		}

	case "float", "decimal", "double":
		if col.IsNullable == "YES" {
			gt = "sql.NullFloat64"
		} else {
			gt = "float64"
		}
	case "year":
		if col.IsNullable == "YES" {
			gt = "sql.NullInt"
		} else {
			gt = "int"
		}

	}

	if gt == "" {
		n := col.TableName + "." + col.ColumnName
		return "", "", errors.New("No compatible datatype (" + col.DataType + ") for " + n + " found")
	}
	return gt, requiredImport, nil
}

func goTypeNotNull(col *ColumnSchema) (string, string, error) {
	requiredImport := ""
	//if col.IsNullable == "YES" {
	//	requiredImport = "database/sql"
	//}

	var gt string

	switch col.DataType {
	case "char", "varchar", "enum", "set", "text", "longtext", "mediumtext", "tinytext":
		gt = "string"
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		gt = "[]byte"
	case "date", "time", "datetime", "timestamp":
		gt, requiredImport = "time.Time", "time"
		if col.IsNullable == "YES" {
			gt = "*time.Time"
		}
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

	//case "bit", "tinyint", "smallint", "int", "mediumint", "bigint":
	//	if col.IsNullable == "YES" {
	//		gt = "sql.NullInt64"
	//	} else {
	//		gt = "int64"
	//	}
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
		return "", "", errors.New("No compatible datatype (" + col.DataType + ") for " + n + " found")
	}

	return gt, requiredImport, nil
}

func writeStructs(schemas []ColumnSchema, tableName string) (int, error) {
	//file, err := os.Create(fmt.Sprintf("%s.go", config.DbName))
	eFile, err := os.Create(fmt.Sprintf("%s/%s.go", entityPath, tableName))
	if err != nil {
		log.Fatal(err)
	}
	defer eFile.Close()

	currentTable := ""

	neededImports := make(map[string]bool)

	// First, get body text into var out
	out := ""
	for _, cs := range schemas {
		if cs.TableName != currentTable {
			if currentTable != "" {
				out = out + "}\n\n"
			}
			out = out + "type " + formatName(cs.TableName) + " struct{\n"
		}

		goType, requiredImport, err := goTypeNotNull(&cs)
		if requiredImport != "" {
			neededImports[requiredImport] = true
		}

		if err != nil {
			log.Fatal(err)
		}

		out = out + "\t" + formatName(cs.ColumnName) + " " + goType
		tags := strings.Split(config.TagLabel, ",")

		tag := ""

		ct := len(tags)
		for i, t := range tags {

			omitEmpty := ""
			if t == "qb" || t == "db" {
				omitEmpty = ",omitempty"
			}

			sp := " "
			if ct == (i + 1) {
				sp = ""
			}
			tag += t + ":\"" + cs.ColumnName + omitEmpty + "\"" + sp
		}

		out += "\t`" + tag + "`"
		out = out + "\n"
		currentTable = cs.TableName

	}

	if out != "" {
		out = out + "}"
	}

	// Now add the header section
	header := "package entity\n\n"
	if len(neededImports) > 0 {
		header = header + "import (\n"
		for imp := range neededImports {
			header = header + "\t\"" + imp + "\"\n"
		}
		header = header + ")\n\n"
	}

	totalBytes, err := fmt.Fprint(eFile, header+out)
	if err != nil {
		log.Fatal(err)
	}

	return totalBytes, nil
}

func writeStruct(sc []ColumnSchema, tableName string) {
	tName := tableName

	if util.SubStringRight(tName, 1) == "s" {
		tName = util.SubStringLeft(tName, len(tName)-1)
	}

	tpl := ParserTemplate{
		TableName:  tableName,
		StructName: util.ToCamelCase(tName),
		ObjectName: util.UpperFirst(util.ToCamelCase(tName)),
		EntityName: util.UpperFirst(util.ToCamelCase(tName)),
		FileName:   util.ToSnakeCase(tName),
	}

	neededImports := make(map[string]bool)

	query := "SELECT \n"

	out := "// " + formatName(tName) + " entity\n" + "type " + formatName(tName) + " struct{\n"
	for i := 0; i < len(sc); i++ {
		goType, requiredImport, err := goTypeNotNull(&sc[i])
		if err != nil {
			log.Fatal(err)
		}

		if requiredImport != "" {
			neededImports[requiredImport] = true
		}

		out = out + "\t" + formatName(sc[i].ColumnName) + " " + goType

		tags := strings.Split(config.TagLabel, ",")

		query += "\t\t\t" + sc[i].ColumnName
		comma := ",\n"
		if len(sc)-1 == i {
			comma = "\n"
		}
		query += comma

		tag := ""

		ct := len(tags)
		for z, t := range tags {

			omitEmpty := ""
			if t == "qb" || t == "db" || t == "query" || t == "url" || t == "form" {
				omitEmpty = ",omitempty"
			}

			sp := " "
			if ct == (z + 1) {
				sp = ""
			}
			tag += t + ":\"" + sc[i].ColumnName + omitEmpty + "\"" + sp
		}

		out += "\t`" + tag + "`"
		out = out + "\n"
	}

	if out != "" {
		out = out + "}"
	}

	// Now add the header section
	header := "// Package entity\n// Automatic generated\npackage entity\n\n"
	if len(neededImports) > 0 {
		header = header + "import (\n"
		for imp := range neededImports {
			header = header + "\t\"" + imp + "\"\n"
		}
		header = header + ")\n\n"
	}

	fName := fmt.Sprintf("%s/%s.go", entityPath, tName)

	exists := fileExist(fName)
	if exists {
		fmt.Println(fmt.Sprintf("file entity already exist %s", fName))
	}

	var eFile *os.File

	if !exists {
		fl, err := os.Create(fName)
		if err != nil {
			fmt.Println(fmt.Sprintf("error create file %s: %v", fName, err))
		}

		eFile = fl
		if eFile != nil {
			_, err = fmt.Fprint(eFile, header+out)
			if err != nil {
				fmt.Println(err)
			}

			eFile.Close()

			fmt.Println("success created entity ", fName)
		}
	}

	query += fmt.Sprintf("\t\tFROM %s", tableName)

	tpl.Query = query

	createRepo(tpl)

}

func Create(cfg Configuration) {
	config = cfg

	flags.Parse(os.Args[2:])

	args := flags.Args()

	if len(args) == 0 {
		log.Fatal(fmt.Errorf("required argument"))
		return
	}

	columns := getSchema(args[0])
	//bytes, err := writeStructs(columns, args[0])
	writeStruct(columns, args[0])
	writePresentation(columns, args[0])
	createUseCaseList(args[0])
}

func CreateLogic() {
	//config = cfg

	flags.Parse(os.Args[2:])

	args := flags.Args()

	if len(args) == 0 {
		log.Fatal(fmt.Errorf("required argument"))
		return
	}
	if len(args) < 2 {
		log.Fatal(fmt.Errorf("required 2 argument, example: gen:logic packageName, fileName"))
		return
	}
	createUseCaseStorer(args[0], args[1])
}

// Package structgen
// @author Daud Valentino
package structgen

import (
	"fmt"
	"log"
	"os"

	"sharefood/pkg/util"
)

func writePresentation(sc []ColumnSchema, tableName string) {
	tName := tableName

	if util.SubStringRight(tName, 1) == "s" {
		tName = util.SubStringLeft(tName, len(tName)-1)
	}

	//tpl := ParserTemplate{
	//	TableName:  tableName,
	//	StructName: util.ToCamelCase(tName),
	//	ObjectName: util.UpperFirst(util.ToCamelCase(tName)),
	//	EntityName: util.UpperFirst(util.ToCamelCase(tName)),
	//	FileName:   util.ToSnakeCase(tName),
	//}
	tags := []string{
		"url",
		"json",
		"db",
	}
	queryParam := "type(\n"
	queryParam += "\t// " + formatName(tName) + "Query parameter\n"
	queryParam += "\t" + formatName(tName) + "Query struct{\n"
	storeParam := "\t// " + formatName(tName) + "Param input param\n"
	storeParam += "\t" + formatName(tName) + "Param struct{\n"
	listData := "\t// " + formatName(tName) + "Detail detail response\n"
	listData += "\t" + formatName(tName) + "Detail struct{\n"
	for i := 0; i < len(sc); i++ {
		gType, isNull, err := layerPresentationType(&sc[i])
		if err != nil {
			log.Fatal(err)
		}

		tagQuery := ""

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
			tagQuery += t + ":\"" + sc[i].ColumnName + omitEmpty + "\"" + sp

		}

		if !util.InArray(sc[i].ColumnName, []string{"created_at", "updated_at", "deleted_at"}) && sc[i].ColumnKey != "PRI" {
			if !isNull {
				queryParam = queryParam + "\t\t" + formatName(sc[i].ColumnName) + " " + gType
				queryParam += "\t`" + tagQuery + "`"
				queryParam = queryParam + "\n"
			}

			storeParam = storeParam + "\t\t" + formatName(sc[i].ColumnName) + " " + gType
			storeParam += "\t`json:\"" + sc[i].ColumnName + ",omitempty\"" + "`\n"
		}

		listData = listData + "\t\t" + formatName(sc[i].ColumnName) + " " + gType
		listData += "\t`json:\"" + sc[i].ColumnName + "\"" + "`\n"

	}

	if queryParam != "" {
		queryParam += "\t\tPaging\n"
		queryParam += "\t\tPeriodRange\n"
		queryParam = queryParam + "\t}\n\n"
	}

	if storeParam != "" {
		storeParam = storeParam + "\t}\n\n"
	}

	if listData != "" {
		listData = listData + "\t}\n\n"
	}

	queryParam += storeParam
	queryParam += listData
	queryParam += "\n)\n"

	fName := fmt.Sprintf("%s/%s.go", presentationPath, util.ToSnakeCase(tName))

	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}
	fl, err := os.Create(fName)
	if err != nil {
		fmt.Println(fmt.Sprintf("error create file %s: %v", fName, err))
	}

	header := "// Package presentations \n"
	header += "// Automatic generated\n"
	header += "package presentations\n\n"
	_, err = fmt.Fprint(fl, header+queryParam)
	if err != nil {
		log.Fatal(err)
	}

	fl.Close()

	fmt.Println("success create presentation ", fName)
}

// Package structgen
// @author Daud Valentino
package structgen

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"sharefood/pkg/util"
)

func fileExist(fName string) bool {
	_, err := os.Stat(fName)
	return errors.Is(err, os.ErrNotExist) == false
}

func contractName(v string) string {
	if util.SubStringRight(v, 1) == "e" {
		return v + "r"
	}
	return v + "er"
}

func createRepo(p ParserTemplate) {
	fName := fmt.Sprintf("%s/%s.go", repoPath, p.FileName)
	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}

	eFile, err := os.Create(fName)
	defer eFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	tpl, err := template.ParseFiles(`./pkg/structgen/repo.tpl`)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.Execute(eFile, p)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success created repo ", fName)
}

func createUseCaseList(tableName string) {
	tName := tableName
	if util.SubStringRight(tName, 1) == "s" {
		tName = util.SubStringLeft(tName, len(tName)-1)
	}

	pkgName := strings.ReplaceAll(util.ToSnakeCase(tableName), "_", "")
	pkgName = fmt.Sprintf("%s/%s", uCasePath, pkgName)
	if !util.PathExist(pkgName) {
		os.MkdirAll(pkgName, os.ModePerm)
	}

	fName := fmt.Sprintf("%s/list.go", pkgName)
	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}

	eFile, err := os.Create(fName)
	defer eFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	tpl, err := template.ParseFiles(`./pkg/structgen/ucase_list.tpl`)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.Execute(eFile, UseCaseTemplate{
		FileName:    "list",
		TableName:   tableName,
		StructName:  util.ToCamelCase(tName),
		EntityName:  util.UpperFirst(util.ToCamelCase(tName)),
		PackageName: strings.ReplaceAll(util.ToSnakeCase(tableName), "_", ""),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success create use case list ", fName)
}

func createUseCaseStorer(packageName, tableName string) {
	tName := tableName
	if util.SubStringRight(tName, 1) == "s" {
		tName = util.SubStringLeft(tName, len(tName)-1)
	}

	//pkgName := strings.ReplaceAll(util.ToSnakeCase(tableName), "_", "")
	pkgName := fmt.Sprintf("%s/%s", uCasePath, packageName)
	if !util.PathExist(pkgName) {
		os.MkdirAll(pkgName, os.ModePerm)
	}

	fln := util.ToSnakeCase(tName)
	fName := fmt.Sprintf("%s/%s.go", pkgName, fln)
	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}

	eFile, err := os.Create(fName)
	defer eFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	tpl, err := template.ParseFiles(`./pkg/structgen/ucase_store.tpl`)
	if err != nil {
		log.Fatal(err)
	}

	en := util.UpperFirst(util.ToCamelCase(tName))
	err = tpl.Execute(eFile, UseCaseTemplate{
		FileName:         fln,
		TableName:        tableName,
		StructName:       util.ToCamelCase(tName),
		EntityName:       util.UpperFirst(en),
		PackageName:      packageName,
		RepoContractName: contractName(util.UpperFirst(en)),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success create use case storer ", fName)
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xwb1989/sqlparser"
)

type config struct {
	input  *os.File
	output *os.File
}

func parseArgs() (*config, error) {
	config := &config{
		input:  os.Stdin,
		output: os.Stdout,
	}

	var (
		input  = flag.String("i", "", "input file, default = stdin")
		output = flag.String("o", "", "output file, default = stdout")
		err    error
	)
	flag.Parse()

	if *input != "" {
		config.input, err = os.OpenFile(*input, os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
	}

	if *output != "" {
		config.output, err = os.OpenFile(*output, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func main() {
	config, err := parseArgs()
	if err != nil {
		panic(err)
	}

	sql, err := ioutil.ReadAll(config.input)
	if err != nil {
		log.Fatalf("can not read")
	}

	stmt, err := sqlparser.Parse(string(sql))
	if err != nil {
		log.Fatalf("can not parse query: %+v", err)
	}

	ddl, ok := stmt.(*sqlparser.DDL)
	if !ok {
		panic("query is not DDL")
	}

	result := [][]string{}
	for _, col := range ddl.TableSpec.Columns {
		row := []string{
			col.Name.String(),
		}

		length := ""
		if col.Type.Length != nil {
			length = "(" + string(col.Type.Length.Val) + ")"
		}

		signed := ""
		if col.Type.Unsigned {
			signed = " unsigned"
		}

		enum := ""
		if len(col.Type.EnumValues) > 0 {
			enum = "("
			for _, e := range col.Type.EnumValues {
				enum += e + ", "
			}
			enum = enum[:len(enum)-2] + ")"
		}

		chartype := ""
		if col.Type.Charset != "" {
			chartype += col.Type.Charset
		}
		if col.Type.Collate != "" {
			chartype += col.Type.Collate
		}

		row = append(
			row,
			fmt.Sprintf("%s%s%s%s%s",
				col.Type.Type,
				enum,
				length,
				signed,
				chartype,
			),
		)

		if col.Type.NotNull {
			row = append(row, "not null")
		} else {
			row = append(row, "")
		}

		if col.Type.Autoincrement {
			row = append(row, "auto_increment")
		} else {
			row = append(row, "")
		}

		if col.Type.Default != nil {
			row = append(row, string(col.Type.Default.Val))
		} else {
			row = append(row, "")
		}

		if col.Type.Comment != nil {
			row = append(row, string(col.Type.Comment.Val))
		} else {
			row = append(row, "")
		}

		result = append(result, row)
	}

	indexResult := [][]string{}
	for _, index := range ddl.TableSpec.Indexes {
		row := []string{
			index.Info.Name.String(),
			index.Info.Type,
		}

		cols := ""
		for _, c := range index.Columns {
			cols += c.Column.String() + ", "
		}
		row = append(row, cols[:len(cols)-2])

		indexResult = append(indexResult, row)
	}

	resultHeader := []string{
		"name",
		"type",
		"not null",
		"auto_increment",
		"default",
		"comment",
	}

	buf := ddl.NewName.Name.String() + " Table's Definition\n\n|"
	for _, h := range resultHeader {
		buf += h + "|"
	}
	buf += "\n"

	buf += "|"
	for range resultHeader {
		buf += "---|"
	}
	buf += "\n"

	for _, row := range result {
		buf += "|"
		for _, col := range row {
			buf += col + "|"
		}
		buf += "\n"
	}

	buf += "\n" + ddl.NewName.Name.String() + " Table's Indexes\n"

	indexResultHeader := []string{
		"name",
		"type",
		"columns",
	}

	buf += "\n|"
	for _, h := range indexResultHeader {
		buf += h + "|"
	}
	buf += "\n"

	buf += "|"
	for range indexResultHeader {
		buf += "---|"
	}
	buf += "\n"

	for _, row := range indexResult {
		buf += "|"
		for _, col := range row {
			buf += col + "|"
		}
		buf += "\n"
	}

	fmt.Println(buf)
}

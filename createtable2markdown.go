package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/xwb1989/sqlparser"
)

type config struct {
	input     *os.File
	output    *os.File
	verbosity bool
}

var logger *log.Logger

type nopWriter struct{}

func (*nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func parseArgs() (*config, error) {
	config := &config{
		input:     os.Stdin,
		output:    os.Stdout,
		verbosity: false,
	}

	var (
		input     = flag.String("i", "", "input file, default = stdin")
		output    = flag.String("o", "", "output file, default = stdout")
		verbosity = flag.Bool("v", false, "verbosity, default = false. if true, say logs")
		err       error
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

	if *verbosity != false {
		config.verbosity = *verbosity
	}

	return config, nil
}

func main() {
	config, err := parseArgs()
	if err != nil {
		panic(err)
	}

	var writer io.Writer
	if config.verbosity {
		writer = os.Stderr
	} else {
		writer = &nopWriter{}
	}
	logger = log.New(writer, "", log.LstdFlags)

	sql, err := ioutil.ReadAll(config.input)
	if err != nil {
		logger.Fatalf("can not read")
	}

	buf := strings.Builder{}

	for _, sqlstring := range strings.Split(string(sql), ";") {
		tableName, tableSpec, ok := parseCreateTabe(sqlstring)
		if !ok {
			continue
		}

		columns := [][]string{}
		for _, col := range tableSpec.Columns {
			columns = append(columns, columnToMarkdown(col))
		}

		indexes := [][]string{}
		for _, index := range tableSpec.Indexes {
			indexes = append(indexes, indexToMarkdown(index))
		}

		buf.WriteString(buildOutput(tableName, columns, indexes))
	}

	fmt.Println(buf.String())
}

func parseCreateTabe(sql string) (string, *sqlparser.TableSpec, bool) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		logger.Printf("can not parse query: %+v\n", err)
		return "", nil, false
	}

	ddl, ok := stmt.(*sqlparser.DDL)
	if !ok {
		logger.Println("query is not DDL")
		return "", nil, false
	}

	if ddl.TableSpec == nil {
		logger.Println("query is not CREATE TABLE")
		return "", nil, false
	}

	return ddl.NewName.Name.String(), ddl.TableSpec, true
}

func columnToMarkdown(col *sqlparser.ColumnDefinition) []string {
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

	return row
}

func indexToMarkdown(index *sqlparser.IndexDefinition) []string {
	row := []string{
		index.Info.Name.String(),
		index.Info.Type,
	}

	cols := ""
	for _, c := range index.Columns {
		cols += c.Column.String() + ", "
	}
	row = append(row, cols[:len(cols)-2])

	return row
}

func buildOutput(tableName string, columns [][]string, indexes [][]string) string {
	buf := strings.Builder{}

	resultHeader := []string{
		"name",
		"type",
		"not null",
		"auto_increment",
		"default",
		"comment",
	}

	buf.WriteString(tableName + " Table's Definition\n\n|")
	for _, h := range resultHeader {
		buf.WriteString(h + "|")
	}
	buf.WriteString("\n")

	buf.WriteString("|")
	for range resultHeader {
		buf.WriteString("---|")
	}
	buf.WriteString("\n")

	for _, row := range columns {
		buf.WriteString("|")
		for _, col := range row {
			buf.WriteString(col + "|")
		}
		buf.WriteString("\n")
	}

	buf.WriteString("\n" + tableName + " Table's Indexes\n")

	indexResultHeader := []string{
		"name",
		"type",
		"columns",
	}

	buf.WriteString("\n|")
	for _, h := range indexResultHeader {
		buf.WriteString(h + "|")
	}
	buf.WriteString("\n")

	buf.WriteString("|")
	for range indexResultHeader {
		buf.WriteString("---|")
	}
	buf.WriteString("\n")

	for _, row := range indexes {
		buf.WriteString("|")
		for _, col := range row {
			buf.WriteString(col + "|")
		}
		buf.WriteString("\n")
	}
	buf.WriteString("\n")

	return buf.String()
}

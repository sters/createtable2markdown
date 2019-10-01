package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_parseCreateTable(t *testing.T) {
	tests := []struct {
		name          string
		sql           string
		wantTablename string
		wantColumns   []string
		wantIndexes   []string
	}{
		{
			"create table",
			`CREATE TABLE foo (
				id int(10) unsigned NOT NULL AUTO_INCREMENT,
				aaa int(10),
				bbb varchar(10),
				ccc varchar(10),
				PRIMARY KEY (id),
				UNIQUE KEY aaa (aaa),
				KEY bbb_ccc (bbb, ccc)
			)`,
			"foo",
			[]string{
				"id",
				"aaa",
				"bbb",
				"ccc",
			},
			[]string{
				"PRIMARY",
				"aaa",
				"bbb_ccc",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotTableName, gotSpec, _ := parseCreateTable(test.sql)
			if gotTableName != test.wantTablename {
				t.Errorf("table name want = %s, got = %s", test.wantTablename, gotTableName)
			}

			gotColumns := make([]string, 0, len(gotSpec.Columns))
			for _, c := range gotSpec.Columns {
				gotColumns = append(gotColumns, c.Name.String())
			}
			if !reflect.DeepEqual(test.wantColumns, gotColumns) {
				t.Errorf("table name want = %s, got = %s", test.wantColumns, gotColumns)
			}

			gotIndexes := make([]string, 0, len(gotSpec.Indexes))
			for _, i := range gotSpec.Indexes {
				gotIndexes = append(gotIndexes, i.Info.Name.String())
			}
			if !reflect.DeepEqual(test.wantIndexes, gotIndexes) {
				t.Errorf("table name want = %s, got = %s", test.wantIndexes, gotIndexes)
			}
		})
	}
}

func Test_columnToMarkdown(t *testing.T) {
	tests := []struct {
		name         string
		sql          string
		wantMarkdown [][]string
	}{
		{
			"create table",
			`CREATE TABLE foo (
				id int(10) unsigned NOT NULL AUTO_INCREMENT,
				aaa int(10),
				bbb varchar(10),
				ccc varchar(10),
				PRIMARY KEY (id),
				UNIQUE KEY aaa (aaa),
				KEY bbb_ccc (bbb, ccc)
			)`,
			[][]string{
				[]string{"id", "int(10) unsigned", "not null", "auto_increment", "", ""},
				[]string{"aaa", "int(10)", "", "", "", ""},
				[]string{"bbb", "varchar(10)", "", "", "", ""},
				[]string{"ccc", "varchar(10)", "", "", "", ""},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotSpec, _ := parseCreateTable(test.sql)
			gotMarkdown := [][]string{}
			for _, c := range gotSpec.Columns {
				gotMarkdown = append(gotMarkdown, columnToMarkdown(c))
			}
			if !reflect.DeepEqual(test.wantMarkdown, gotMarkdown) {
				t.Errorf("table name want = %s, got = %s", test.wantMarkdown, gotMarkdown)
			}
		})
	}
}

func Test_indexToMarkdown(t *testing.T) {
	tests := []struct {
		name         string
		sql          string
		wantMarkdown [][]string
	}{
		{
			"create table",
			`CREATE TABLE foo (
				id int(10) unsigned NOT NULL AUTO_INCREMENT,
				aaa int(10),
				bbb varchar(10),
				ccc varchar(10),
				PRIMARY KEY (id),
				UNIQUE KEY aaa (aaa),
				KEY bbb_ccc (bbb, ccc)
			)`,
			[][]string{
				[]string{"PRIMARY", "primary key", "id"},
				[]string{"aaa", "unique key", "aaa"},
				[]string{"bbb_ccc", "key", "bbb, ccc"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotSpec, _ := parseCreateTable(test.sql)
			gotMarkdown := [][]string{}
			for _, i := range gotSpec.Indexes {
				gotMarkdown = append(gotMarkdown, indexToMarkdown(i))
			}
			if !reflect.DeepEqual(test.wantMarkdown, gotMarkdown) {
				t.Errorf("table name want = %s, got = %s", test.wantMarkdown, gotMarkdown)
			}
		})
	}
}

func Test_buildOutput(t *testing.T) {
	tests := []struct {
		name      string
		tablename string
		columns   [][]string
		indexes   [][]string
		want      string
	}{
		{
			"output",
			"test",
			[][]string{
				[]string{"id", "int(10) unsigned", "not null", "auto_increment", "", ""},
			},
			[][]string{
				[]string{"PRIMARY", "primary key", "id"},
			},
			`
test Table's Definition

|name|type|not null|auto_increment|default|comment|
|---|---|---|---|---|---|
|id|int(10) unsigned|not null|auto_increment|||

test Table's Indexes

|name|type|columns|
|---|---|---|
|PRIMARY|primary key|id|
			`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := buildOutput(test.tablename, test.columns, test.indexes)

			want := strings.TrimSpace(test.want)
			got = strings.TrimSpace(got)
			if !reflect.DeepEqual(want, got) {
				t.Errorf("table name want = %s, got = %s", want, got)
			}
		})
	}
}

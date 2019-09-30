# createtable2markdown

## Installation

```
go get -u github.com/sters/createtable2markdown
```

## Usage

```
cat foo.sql | createtable2markdown
```

also

```
mysqldump -u database_user -p -d -n --compact database_name | createtable2markdown
```

You can get markdown table, like this:

```
foo Table's Definition

|name|type|not null|auto_increment|default|comment|
|---|---|---|---|---|---|
|id|int(10) unsigned|not null||||
|body|text|not null||||

foo Table's Indexes

|name|type|columns|
|---|---|---|
|PRIMARY|primary key|id|
```

If you want beautiful markdown table, see [sters/markdown-table-formatter](https://github.com/sters/markdown-table-formatter)

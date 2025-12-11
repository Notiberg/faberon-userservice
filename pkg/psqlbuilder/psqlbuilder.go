package psqlbuilder

import "github.com/Masterminds/squirrel"

var placeholder = squirrel.Dollar

func Update(table string) squirrel.UpdateBuilder {
	return squirrel.Update(table).PlaceholderFormat(placeholder)
}

func Insert(table string) squirrel.InsertBuilder {
	return squirrel.Insert(table).PlaceholderFormat(placeholder)
}

func Delete(table string) squirrel.DeleteBuilder {
	return squirrel.Delete(table).PlaceholderFormat(placeholder)
}

func Select(columns ...string) squirrel.SelectBuilder {
	return squirrel.Select(columns...).PlaceholderFormat(placeholder)
}

package mypg

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/xwb1989/sqlparser"
)

type QueryData struct {
	Type       string
	Table      string
	Columns    map[string]string // SET or INSERT columns
	WherePairs map[string]string // WHERE condition as key:value
	SelectCols []string
	Data       pgx.Row
}

func ParseWhere(expr *sqlparser.Where) map[string]string {
	result := map[string]string{}
	if expr == nil {
		return result
	}

	cond := expr.Expr

	switch node := cond.(type) {

	case *sqlparser.ComparisonExpr:
		left := sqlparser.String(node.Left)
		right := sqlparser.String(node.Right)
		result[left] = right

	case *sqlparser.AndExpr:
		// LEFT side
		if leftExpr, ok := node.Left.(*sqlparser.ComparisonExpr); ok {
			left := sqlparser.String(leftExpr.Left)
			right := sqlparser.String(leftExpr.Right)
			result[left] = right
		}

		// RIGHT side
		if rightExpr, ok := node.Right.(*sqlparser.ComparisonExpr); ok {
			left := sqlparser.String(rightExpr.Left)
			right := sqlparser.String(rightExpr.Right)
			result[left] = right
		}
	}
	return result
}

func ReactiveToPublish(sql string, data pgx.Row, args ...interface{}) (*QueryData, error) {
	q := Interpolate(sql, args...)
	stmt, err := sqlparser.Parse(q)
	if err != nil {
		return nil, err
	}

	res := &QueryData{
		Columns:    map[string]string{},
		WherePairs: map[string]string{},
		Data:       data,
	}

	switch s := stmt.(type) {

	case *sqlparser.Insert:
		res.Type = "INSERT"
		res.Table = s.Table.Name.String()

		// columns
		cols := []string{}
		for _, col := range s.Columns {
			cols = append(cols, col.String())
		}

		// values
		values := s.Rows.(sqlparser.Values)[0]
		for i, val := range values {
			res.Columns[cols[i]] = sqlparser.String(val)
		}

	// ----------------------------------------
	// UPDATE
	// ----------------------------------------
	case *sqlparser.Update:
		res.Type = "UPDATE"
		table := s.TableExprs[0].(*sqlparser.AliasedTableExpr)
		res.Table = sqlparser.String(table.Expr)

		// SET column = value
		for _, expr := range s.Exprs {
			res.Columns[expr.Name.Name.String()] = sqlparser.String(expr.Expr)
		}

		// WHERE
		res.WherePairs = ParseWhere(s.Where)

	// ----------------------------------------
	// DELETE
	// ----------------------------------------
	case *sqlparser.Delete:
		res.Type = "DELETE"
		table := s.TableExprs[0].(*sqlparser.AliasedTableExpr)
		res.Table = sqlparser.String(table.Expr)

		res.WherePairs = ParseWhere(s.Where)

	// ----------------------------------------
	// SELECT
	// ----------------------------------------
	case *sqlparser.Select:
		res.Type = "SELECT"
		table := s.From[0].(*sqlparser.AliasedTableExpr)
		res.Table = sqlparser.String(table.Expr)

		for _, col := range s.SelectExprs {
			if c, ok := col.(*sqlparser.AliasedExpr); ok {
				res.SelectCols = append(res.SelectCols, sqlparser.String(c.Expr))
			}
		}

		res.WherePairs = ParseWhere(s.Where)

	default:
		return nil, fmt.Errorf("unsupported query")
	}

	return res, nil
}
func Parser(sql string, args ...interface{}) (*QueryData, error) {
	q := Interpolate(sql, args...)
	stmt, err := sqlparser.Parse(q)
	if err != nil {
		return nil, err
	}

	res := &QueryData{
		Columns:    map[string]string{},
		WherePairs: map[string]string{},
	}

	switch s := stmt.(type) {

	case *sqlparser.Insert:
		res.Type = "INSERT"
		res.Table = s.Table.Name.String()

		// columns
		cols := []string{}
		for _, col := range s.Columns {
			cols = append(cols, col.String())
		}

		// values
		values := s.Rows.(sqlparser.Values)[0]
		for i, val := range values {
			res.Columns[cols[i]] = sqlparser.String(val)
		}

	// ----------------------------------------
	// UPDATE
	// ----------------------------------------
	case *sqlparser.Update:
		res.Type = "UPDATE"
		table := s.TableExprs[0].(*sqlparser.AliasedTableExpr)
		res.Table = sqlparser.String(table.Expr)

		// SET column = value
		for _, expr := range s.Exprs {
			res.Columns[expr.Name.Name.String()] = sqlparser.String(expr.Expr)
		}

		// WHERE
		res.WherePairs = ParseWhere(s.Where)

	// ----------------------------------------
	// DELETE
	// ----------------------------------------
	case *sqlparser.Delete:
		res.Type = "DELETE"
		table := s.TableExprs[0].(*sqlparser.AliasedTableExpr)
		res.Table = sqlparser.String(table.Expr)

		res.WherePairs = ParseWhere(s.Where)

	// ----------------------------------------
	// SELECT
	// ----------------------------------------
	case *sqlparser.Select:
		res.Type = "SELECT"
		table := s.From[0].(*sqlparser.AliasedTableExpr)
		res.Table = sqlparser.String(table.Expr)

		for _, col := range s.SelectExprs {
			if c, ok := col.(*sqlparser.AliasedExpr); ok {
				res.SelectCols = append(res.SelectCols, sqlparser.String(c.Expr))
			}
		}

		res.WherePairs = ParseWhere(s.Where)

	default:
		return nil, fmt.Errorf("unsupported query")
	}

	return res, nil

}
func Interpolate(sql string, args ...interface{}) string {

	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)

		switch v := arg.(type) {
		case string:
			sql = strings.Replace(sql, placeholder, "'"+v+"'", 1)
		default:
			sql = strings.Replace(sql, placeholder, fmt.Sprint(v), 1)
		}
	}
	return sql
}

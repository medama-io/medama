package qb

import (
	"strings"
)

type CTE struct {
	Name     string
	Subquery *QueryBuilder
}

type QueryBuilder struct {
	cteClauses       []CTE
	selectClauses    []string
	fromClause       []string
	leftJoinClause   string
	whereClause      string
	groupByClause    []string
	havingClause     []string
	orderByClause    []string
	paginationClause string
}

func New() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) WithMaterialized(name string, subquery *QueryBuilder) *QueryBuilder {
	qb.cteClauses = append(qb.cteClauses, CTE{Name: name, Subquery: subquery})
	return qb
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.selectClauses = append(qb.selectClauses, columns...)
	return qb
}

func (qb *QueryBuilder) From(tables ...string) *QueryBuilder {
	qb.fromClause = append(qb.fromClause, tables...)
	return qb
}

func (qb *QueryBuilder) LeftJoin(query string) *QueryBuilder {
	qb.leftJoinClause = query
	return qb
}

func (qb *QueryBuilder) Where(query string) *QueryBuilder {
	qb.whereClause = query
	return qb
}

func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groupByClause = append(qb.groupByClause, columns...)
	return qb
}

func (qb *QueryBuilder) Having(query string) *QueryBuilder {
	qb.havingClause = append(qb.havingClause, query)
	return qb
}

func (qb *QueryBuilder) OrderBy(columns ...string) *QueryBuilder {
	qb.orderByClause = append(qb.orderByClause, columns...)
	return qb
}

func (qb *QueryBuilder) Pagination(query string) *QueryBuilder {
	qb.paginationClause = query
	return qb
}

func (qb *QueryBuilder) Build() string {
	var query strings.Builder

	for idx, cte := range qb.cteClauses {
		// First CTE starts with WITH, while subsequent CTEs start with a comma
		if idx == 0 {
			query.WriteString("WITH ")
		} else {
			query.WriteString(", ")
		}

		query.WriteString(cte.Name)
		query.WriteString(" AS MATERIALIZED (")
		query.WriteString(cte.Subquery.Build())
		query.WriteString(")")
	}

	query.WriteString("SELECT ")
	query.WriteString(strings.Join(qb.selectClauses, ", "))

	query.WriteString(" FROM ")
	query.WriteString(strings.Join(qb.fromClause, ", "))

	if qb.leftJoinClause != "" {
		query.WriteString(" LEFT JOIN ")
		query.WriteString(qb.leftJoinClause)
	}

	if len(qb.whereClause) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(qb.whereClause)
	}

	if len(qb.groupByClause) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(qb.groupByClause, ", "))
	}

	if len(qb.havingClause) > 0 {
		query.WriteString(" HAVING ")
		query.WriteString(strings.Join(qb.havingClause, ", "))
	}

	if len(qb.orderByClause) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(qb.orderByClause, ", "))
	}

	if qb.paginationClause != "" {
		query.WriteString(qb.paginationClause)
	}

	return query.String()
}

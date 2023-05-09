package arangodb

import "fmt"

type QueryOpt func(*Query)

func WithPaging(offset, size uint) QueryOpt {
	return func(q *Query) {
		if size == 0 {
			size = 25
		}

		q.limitStmt = "LIMIT @offset, @size"
		q.params["offset"] = offset
		q.params["size"] = size
	}
}

func WithSort(field string, ascending bool) QueryOpt {
	return func(q *Query) {
		order := "DESC"
		if ascending {
			order = "ASC"
		}

		q.sortStmt = fmt.Sprintf("SORT d.%s %s", field, order)
	}
}

func WithField(field string, paramKey string, value any) QueryOpt {
	return func(q *Query) {
		if q.filterStmt != "" {
			q.filterStmt += " && "
		}

		q.filterStmt += fmt.Sprintf("d.%s == @%s", field, paramKey)
		q.params[paramKey] = value
	}
}

func WithFilter(filter string, params map[string]interface{}) QueryOpt {
	return func(q *Query) {
		q.filterStmt = fmt.Sprintf(" %s", filter)
		for k, v := range params {
			q.params[k] = v
		}
	}
}

func NewQuery(col string, opts ...QueryOpt) *Query {
	q := &Query{
		collection: col,
		params:     make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(q)
	}

	return q
}

type Query struct {
	collection string
	filterStmt string
	limitStmt  string
	sortStmt   string
	params     map[string]interface{}
}

func (q *Query) Statement() string {
	result := fmt.Sprintf("FOR d IN %s", q.collection)

	if q.filterStmt != "" {
		result += " FILTER " + q.filterStmt
	}

	if q.sortStmt != "" {
		result += " " + q.sortStmt
	}

	if q.limitStmt != "" {
		result += " " + q.limitStmt
	}

	result += " RETURN d"

	return result
}

func (q *Query) Params() map[string]interface{} {
	return q.params
}

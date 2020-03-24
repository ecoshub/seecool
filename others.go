package seecool

import "strconv"

// using
func (q *query) Using(tables ...string) *query {
	q.using = tables
	return q
}

// orders
func (q *query) Order(col string) *query {
	q.order = " ORDER BY " + col
	return q
}

func (q *query) OrderDesc(col string) *query {
	q.order = " ORDER BY " + col + " DESC"
	return q
}

// group
func (q *query) Group(cols ...string) *query {
	q.groups = cols
	for _, c := range cols {
		exs := false
		for _, g := range q.columns {
			if c == g {
				exs = true
				break
			}
		}
		if !exs {
			q.columns = append(q.columns, c)
		}
	}
	return q
}

// limits
func (q *query) Limit(lim string) *query {
	q.limit = " LIMIT " + lim
	return q
}

func (q *query) LimitInt(lim int) *query {
	q.limit = " LIMIT " + strconv.Itoa(lim)
	return q
}

func (q *query) LimitOff(lim, off string) *query {
	q.limit = " LIMIT " + lim + " OFFSET " + off
	return q
}

func (q *query) LimitOffInt(lim, off int) *query {
	q.limit = " LIMIT " + strconv.Itoa(lim) + " OFFSET " + strconv.Itoa(off)
	return q
}

// insert element
func (q *query) Keys(keys ...string) *query {
	q.keys = keys
	return q
}

func (q *query) Values(values ...string) *query {
	q.values = values
	return q
}

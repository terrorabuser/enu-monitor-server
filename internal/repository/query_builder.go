package repository

import (
	"fmt"
	"strings"
	"time"
)

type ContentQueryBuilder struct {
	query strings.Builder
	args []interface{}
}


func NewContentQueryBuilder() *ContentQueryBuilder{
	qb := &ContentQueryBuilder{}
	qb.query.WriteString(`SELECT 
    c.id, c.user_id, c.macaddress, c.file_name, 
    c.file_path, c.start_time, c.end_time,
    ch.id AS history_id, ch.content_id AS history_content_id, 
    ch.status_id, ch.created_at, 
    ch.user_id, ch.reason AS history_user_id
	FROM content c
	LEFT JOIN content_history ch ON ch.content_id = c.id
		AND ch.id = (SELECT MAX(id) FROM content_history WHERE content_id = c.id)
	WHERE 1=1
	`)
	return qb
}


func (qb *ContentQueryBuilder) ApplyUserId(userId *int64) *ContentQueryBuilder {
	if userId != nil && *userId != 0 {
		qb.query.WriteString(fmt.Sprintf(" AND c.user_id = $%d", len(qb.args)+1))
		qb.args = append(qb.args, *userId)
	}
	return qb
}

func (qb *ContentQueryBuilder) ApplyStatusId(statusId *int32) *ContentQueryBuilder {
	if statusId != nil && *statusId != 0 {
		qb.query.WriteString(fmt.Sprintf(" AND ch.status_id = $%d", len(qb.args)+1))
		qb.args = append(qb.args, *statusId)
	}
	return qb
}

func (qb *ContentQueryBuilder) ApplyTimeFilters(startTime, endTime *time.Time) *ContentQueryBuilder {
    if startTime != nil && !startTime.IsZero() {
        qb.query.WriteString(fmt.Sprintf(" AND c.start_time >= $%d", len(qb.args)+1))
        qb.args = append(qb.args, *startTime)
    }

    if endTime != nil && !endTime.IsZero() {
        qb.query.WriteString(fmt.Sprintf(" AND c.end_time <= $%d", len(qb.args)+1))
        qb.args = append(qb.args, *endTime)
    }
    return qb
}


func (qb *ContentQueryBuilder) Build() (string, []interface{}) {
	qb.query.WriteString(" ORDER BY c.id ASC")
	return qb.query.String(), qb.args
}

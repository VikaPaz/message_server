package message

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/VikaPaz/message_server/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type MessageRepository struct {
	conn *sql.DB
	log  *logrus.Logger
}

func NewRepository(conn *sql.DB, logger *logrus.Logger) *MessageRepository {
	return &MessageRepository{
		conn: conn,
		log:  logger,
	}
}

func (r *MessageRepository) Create(m models.CreateRequest, status models.Status) (models.Message, error) {
	builder := sq.Insert(`messages (message, status, created_at)`)
	builder = builder.Values(m.Message, status, time.Now())
	builder = builder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return models.Message{}, err
	}

	r.log.Debugf("Executing query: %v", query)
	row := r.conn.QueryRow(query, args...)
	if err = row.Err(); err != nil {
		return models.Message{}, err
	}

	me := models.Message{}
	err = row.Scan(&me.ID, &me.Message, &me.Status, &me.CreatedAt, &me.UpdatedAt)
	if err != nil {
		return models.Message{}, err
	}

	return me, err
}

func (r *MessageRepository) Get(m models.Message, limit uint64, offset uint64) (models.FilterResponse, error) {
	resp := models.FilterResponse{}

	builder := sq.Select("count(*) over () ", "id", "message", "status", "created_at", "updated_at").From("messages")
	builder = builder.PlaceholderFormat(sq.Dollar)

	if m.ID != nil {
		builder = builder.Where(sq.Eq{"id": *m.ID})
	}
	if m.Message != nil {
		builder = builder.Where(sq.ILike{"message": fmt.Sprintf("%%%v%%", *m.Message)})
	}
	if m.Status != nil {
		builder = builder.Where(sq.Eq{"status": *m.Status})
	}
	if m.CreatedAt != nil {
		builder = builder.Where(sq.GtOrEq{"created_at": *m.CreatedAt})
	}
	if m.UpdatedAt != nil {
		builder = builder.Where(sq.LtOrEq{"updated_at": *m.UpdatedAt})
	}
	if limit > 0 {
		builder = builder.Limit(limit)
	}
	if offset > 0 {
		builder = builder.Offset(offset)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return models.FilterResponse{}, err
	}

	r.log.Debugf("Executing query: %v", query)
	rows, err := r.conn.Query(query, args...)
	if err != nil {
		return models.FilterResponse{}, err
	}

	for rows.Next() {
		me := models.Message{}
		err = rows.Scan(&resp.Total, &me.ID, &me.Message, &me.Status, &me.CreatedAt, &me.UpdatedAt)
		if err != nil {
			return models.FilterResponse{}, err
		}
		resp.Messages = append(resp.Messages, me)
	}
	r.log.Debugf("Returning results: %v", resp)
	return resp, nil
}

func (r *MessageRepository) Update(id uuid.UUID, status models.Status) (models.Message, error) {
	builder := sq.Update("messages").Set("status", status).Set("updated_at", time.Now().UTC())
	builder = builder.Where(sq.Eq{"id": id})
	builder = builder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return models.Message{}, err
	}

	r.log.Debugf("Executing query: %v", query)
	row := r.conn.QueryRow(query, args...)
	if err = row.Err(); err != nil {
		return models.Message{}, err
	}

	me := models.Message{}
	err = row.Scan(&me.ID, &me.Message, &me.Status, &me.CreatedAt, &me.UpdatedAt)
	if err != nil {
		return models.Message{}, err
	}

	return me, err

}

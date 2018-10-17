package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type MessageService struct {
	db *sql.DB
	ns store.NotificationService
}

const (
	createMessageTableSQL = `
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			sender int NOT NULL,
			reciever int NOT NULL,
			message varchar(512) NOT NULL,
			date int NOT NULL
		);
	`
)

func NewMessageService(db *sql.DB, ns store.NotificationService) store.MessageService {
	if _, err := db.Exec(createMessageTableSQL); err != nil {
		panic(err)
	}
	return &MessageService{db, ns}
}

func (s *MessageService) Add(msg *store.Message) error {
	err := s.db.QueryRow("INSERT INTO messages(sender, reciever, message, date) VALUES($1,$2,$3,$4) RETURNING id", &msg.Sender, &msg.Reciever, &msg.Message, &msg.Date).Scan(&msg.ID)
	{
		if err != nil {
			return err
		}
	}
	s.ns.Push(msg.Reciever, msg.Sender, "message", msg)
	return nil
}

func (s *MessageService) Remove(id int64) error {
	_, err := s.db.Exec("DELETE FROM messages WHERE id=$1", id)
	return err
}

func (s *MessageService) Messages(userId int64) ([]*store.Message, error) {
	var messages []*store.Message
	rows, err := s.db.Query("SELECT * FROM messages WHERE sender=$1 OR reciever=$1", userId)
	{
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		m := &store.Message{}
		err := rows.Scan(&m.ID, &m.Sender, &m.Reciever, &m.Message, &m.Date)
		{
			if err != nil {
				return nil, err
			}
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (s *MessageService) MessagesBetween(u1, u2 int64) ([]*store.Message, error) {
	var messages []*store.Message
	rows, err := s.db.Query("SELECT * FROM messages WHERE (sender=$1 AND reciever=$2) OR (sender=$2 AND reciever=$1)", u1, u2)
	{
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		m := &store.Message{}
		err := rows.Scan(&m.ID, &m.Sender, &m.Reciever, &m.Message, &m.Date)
		{
			if err != nil {
				return nil, err
			}
		}
		messages = append(messages, m)
	}
	return messages, nil
}

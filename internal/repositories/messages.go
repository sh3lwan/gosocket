package repositories

import (
	"database/sql"
	"fmt"

	. "github.com/sh3lwan/gosocket/internal/db"
	. "github.com/sh3lwan/gosocket/internal/models/messages"
)

func GetMessages(receiver string) ([]ReceivedMessage, error) {
	var messages []ReceivedMessage = make([]ReceivedMessage, 0)
	var rows *sql.Rows
	var err error

	if receiver == "" {
		rows, err = DB().Query("SELECT id, message, username, is_new FROM messages WHERE receiver = '' OR receiver IS NULL;")
	} else {
		rows, err = DB().Query("SELECT id, message, username, is_new FROM messages WHERE receiver IS NOT NULL AND receiver <> '' AND (receiver = ? OR username = ?);", receiver, receiver)
	}

	if err != nil {
		fmt.Printf("Error querying database: %v", err)
		return messages, err
	}

	defer rows.Close()

	for rows.Next() {
		var message ReceivedMessage

		err = rows.Scan(&message.Id, &message.Body, &message.Username, &message.IsNew)

		if err != nil {
			fmt.Printf("Error while reading messages rows: %v\n", err)
			continue
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error after iterating over rows: %v\n", err)
		return messages, err
	}

	return messages, nil
}

func InsertMessage(received ReceivedMessage) (int64, error) {
	stmt, err := DB().Prepare("INSERT INTO `messages`(message, username, receiver, is_new) VALUES(?, ?, ?, ?);")

	if err != nil {
		fmt.Printf("Error querying database: %v", err)
		return 0, err
	}
	result, err := stmt.Exec(received.Body, received.Username, received.Receiver, received.IsNew)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

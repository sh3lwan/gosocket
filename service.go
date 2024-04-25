package main

import "fmt"

func getMessages() []ReceivedMessage {
	var messages []ReceivedMessage

	rows, err := DB().Query("SELECT id, message, username FROM messages;")

	if err != nil {
		fmt.Printf("Error querying database: %v", err)
	}

	for rows.Next() {
		var message ReceivedMessage

		err = rows.Scan(&message.Id, &message.Message, &message.Username)

		if err != nil {
			fmt.Printf("Error while reading messages rows: %v", err)
		}

		messages = append(messages, message)
	}

	return messages
}

func insertMessage(received ReceivedMessage) (int64, error) {

	stmt, err := DB().Prepare("INSERT INTO `messages`(message, username, is_new) VALUES(?, ?, ?);")

	if err != nil {
		fmt.Printf("Error querying database: %v", err)
	}
	result, err := stmt.Exec(received.Message, received.Username, received.IsNew)

	return result.LastInsertId()
}

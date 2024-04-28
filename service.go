package main

import "fmt"

func getMessages() ([]ReceivedMessage, error) {
	var messages []ReceivedMessage = make([]ReceivedMessage, 0)

	rows, err := DB().Query("SELECT id, message, username, is_new FROM messages;")

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

func insertMessage(received ReceivedMessage) (int64, error) {
	fmt.Printf("Inserted Message: %v\n", received)

	stmt, err := DB().Prepare("INSERT INTO `messages`(message, username, is_new) VALUES(?, ?, ?);")

	if err != nil {
		fmt.Printf("Error querying database: %v", err)
		return 0, err
	}

	result, err := stmt.Exec(received.Body, received.Username, received.IsNew)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

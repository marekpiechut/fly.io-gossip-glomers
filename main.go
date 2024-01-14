package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	idGenerator := NewIdGenerator()
	broadcaster := NewBroadcaster()

	node.Handle("echo", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "echo_ok"

		return node.Reply(msg, body)

	})

	node.Handle("generate", func(msg maelstrom.Message) error {
		var body = map[string]any{}
		body["type"] = "generate_ok"
		body["id"] = idGenerator.Next()
		return node.Reply(msg, body)
	})

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if err := broadcaster.Add(body["message"].(float64)); err != nil {
			return err
		}

		return node.Reply(msg, map[string]any{
			"type": "broadcast_ok",
		})
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		var messages = broadcaster.Get()
		return node.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": messages,
		})
	})

	node.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		return node.Reply(msg, map[string]any{
			"type": "topology_ok",
		})
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	idGenerator := NewSnowflakeIdGenerator(-1, -1)
	broadcaster := NewBroadcaster(node)

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
		type msgBody struct {
			Message int `json:"message"`
		}
		var body msgBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if err := broadcaster.Add(body.Message); err != nil {
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
		type msgBody struct {
			Topology map[string][]string `json:"topology"`
		}

		var body msgBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		broadcaster.topology.Update(body.Topology)

		return node.Reply(msg, map[string]any{
			"type": "topology_ok",
		})
	})

	node.Handle("propagate", func(msg maelstrom.Message) error {
		var body PropagateBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		return broadcaster.Add(body.Value)
	})
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}

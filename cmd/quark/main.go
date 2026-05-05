// ABOUTME: CLI entry point for running a Quark node and interacting
// ABOUTME: with a running node over HTTP.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"quark"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "node":
		runNode(os.Args[2:])
	case "mine":
		runMine(os.Args[2:])
	case "send":
		runSend(os.Args[2:])
	case "balance":
		runBalance(os.Args[2:])
	case "address":
		runAddress(os.Args[2:])
	case "sync":
		runSync(os.Args[2:])
	case "peer":
		runPeer(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  quark node    --listen :8080 --data node.json [--peers url1,url2]
  quark mine    --node http://host:port
  quark send    --node http://host:port --to ADDR --amount N [--nonce N]
  quark balance --node http://host:port --address ADDR
  quark address --node http://host:port
  quark sync    --node http://host:port
  quark peer    --node http://host:port --url http://other:port`)
}

func runNode(args []string) {
	fs := flag.NewFlagSet("node", flag.ExitOnError)
	listen := fs.String("listen", ":8080", "address to listen on")
	dataPath := fs.String("data", "", "node state file (loaded if exists, written on shutdown)")
	peers := fs.String("peers", "", "comma-separated peer URLs")
	_ = fs.Parse(args)

	var node *quark.Node
	var err error
	if *dataPath != "" {
		node, err = quark.LoadOrCreateNode(*dataPath)
	} else {
		node, err = quark.NewNode()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create node:", err)
		os.Exit(1)
	}

	srv := quark.NewServer(node)
	if *peers != "" {
		for _, p := range strings.Split(*peers, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				srv.AddPeer(p)
			}
		}
	}

	fmt.Printf("quark node listening on %s (address %s)\n", *listen, node.Address())
	fmt.Println("peers:", srv.Peers())
	if err := http.ListenAndServe(*listen, srv.Handler()); err != nil {
		fmt.Fprintln(os.Stderr, "server error:", err)
		if *dataPath != "" {
			_ = node.Save(*dataPath)
		}
		os.Exit(1)
	}
}

func runMine(args []string) {
	fs := flag.NewFlagSet("mine", flag.ExitOnError)
	node := fs.String("node", "http://localhost:8080", "node URL")
	_ = fs.Parse(args)
	post(*node+"/mine", nil)
}

func runSend(args []string) {
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	node := fs.String("node", "http://localhost:8080", "node URL")
	to := fs.String("to", "", "recipient address")
	amount := fs.Int64("amount", 0, "amount to send")
	nonce := fs.Int64("nonce", 0, "transaction nonce")
	_ = fs.Parse(args)
	if *to == "" || *amount <= 0 {
		fmt.Fprintln(os.Stderr, "send requires --to and positive --amount")
		os.Exit(2)
	}
	post(*node+"/send", map[string]any{
		"recipient": *to,
		"amount":    *amount,
		"nonce":     *nonce,
	})
}

func runBalance(args []string) {
	fs := flag.NewFlagSet("balance", flag.ExitOnError)
	node := fs.String("node", "http://localhost:8080", "node URL")
	address := fs.String("address", "", "address to query")
	_ = fs.Parse(args)
	if *address == "" {
		fmt.Fprintln(os.Stderr, "balance requires --address")
		os.Exit(2)
	}
	get(*node + "/balance?address=" + *address)
}

func runAddress(args []string) {
	fs := flag.NewFlagSet("address", flag.ExitOnError)
	node := fs.String("node", "http://localhost:8080", "node URL")
	_ = fs.Parse(args)
	get(*node + "/address")
}

func runSync(args []string) {
	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	node := fs.String("node", "http://localhost:8080", "node URL")
	_ = fs.Parse(args)
	post(*node+"/sync", nil)
}

func runPeer(args []string) {
	fs := flag.NewFlagSet("peer", flag.ExitOnError)
	node := fs.String("node", "http://localhost:8080", "node URL")
	url := fs.String("url", "", "peer URL to add")
	_ = fs.Parse(args)
	if *url == "" {
		fmt.Fprintln(os.Stderr, "peer requires --url")
		os.Exit(2)
	}
	post(*node+"/peers", map[string]string{"url": *url})
}

func post(url string, body any) {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	resp, err := http.Post(url, "application/json", &buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "request failed:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	fmt.Println(string(out))
	if resp.StatusCode >= 400 {
		os.Exit(1)
	}
}

func get(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "request failed:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	fmt.Println(string(out))
	if resp.StatusCode >= 400 {
		os.Exit(1)
	}
}

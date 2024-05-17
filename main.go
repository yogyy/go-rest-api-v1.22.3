package main

func main() {
	server := NewAPIServer("127.0.0.1:8080")
	server.Run()
}

package main

import (
	"juego-websocket/game"
	"log"
	"net/http"
)

func main() {
	server := game.NewServer()

	// Iniciar bucle del juego en goroutine
	go server.GameLoop()

	// Configurar rutas
	http.HandleFunc("/ws", server.HandleWebSocket)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Servidor de Juego WebSocket funcionando"))
	})

	// Servir archivos estáticos (cliente)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
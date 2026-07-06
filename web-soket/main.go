package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

// Cliente representa un usuario conectado
type Client struct {
    conn *websocket.Conn
    send chan []byte
    id   string
}

// Hub mantiene todos los clientes conectados
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mutex      sync.Mutex
}

// Nuevo hub
func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

// Función principal del hub - maneja todas las conexiones
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mutex.Lock()
            h.clients[client] = true
            h.mutex.Unlock()
            log.Printf("Cliente conectado: %s", client.id)
            
        case client := <-h.unregister:
            h.mutex.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                log.Printf("Cliente desconectado: %s", client.id)
            }
            h.mutex.Unlock()
            
        case message := <-h.broadcast:
            h.mutex.Lock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mutex.Unlock()
        }
    }
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // Permitir todas las conexiones (solo para desarrollo)
    },
}

// Maneja cada conexión WebSocket
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    client := &Client{
        conn: conn,
        send: make(chan []byte, 256),
        id:   r.RemoteAddr,
    }

    h.register <- client

    // Leer mensajes entrantes
    go client.readPump(h)

    // Enviar mensajes salientes
    go client.writePump()
}

// Lee mensajes del cliente
func (c *Client) readPump(h *Hub) {
    defer func() {
        h.unregister <- c
        c.conn.Close()
    }()

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Error de lectura: %v", err)
            }
            break
        }

        // Reenviar mensaje a todos los clientes
        h.broadcast <- message
    }
}

// Envía mensajes al cliente
func (c *Client) writePump() {
    defer c.conn.Close()

    for {
        message, ok := <-c.send
        if !ok {
            c.conn.WriteMessage(websocket.CloseMessage, []byte{})
            return
        }

        w, err := c.conn.NextWriter(websocket.TextMessage)
        if err != nil {
            return
        }
        w.Write(message)

        // Enviar mensajes en cola
        n := len(c.send)
        for i := 0; i < n; i++ {
            w.Write([]byte{'\n'})
            w.Write(<-c.send)
        }

        if err := w.Close(); err != nil {
            return
        }
    }
}

func main() {
    hub := NewHub()
    go hub.Run()

    // Servir archivos estáticos
    http.Handle("/", http.FileServer(http.Dir("./static")))

    // Endpoint WebSocket
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        hub.HandleWebSocket(w, r)
    })

    // Iniciar servidor
    fmt.Println("Servidor iniciado en http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
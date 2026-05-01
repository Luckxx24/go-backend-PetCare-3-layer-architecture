package ws

import (
	"log"
	"sync"

	"github.com/google/uuid"
)


type Room struct {
	clients map[*Client]bool 
	mu      sync.RWMutex    
}


type Hub struct {
	rooms      map[uuid.UUID]*Room 
	Register   chan *Client        
	Unregister chan *Client        
	Broadcast  chan BroadcastMsg   
	mu         sync.RWMutex
}


/
type BroadcastMsg struct {
	RoomID  uuid.UUID  
	Message WsResponse 
}


func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[uuid.UUID]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan BroadcastMsg),
	}
}


func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.mu.Lock()

			
			if _, exists := h.rooms[client.BookingID]; !exists {
				h.rooms[client.BookingID] = &Room{
					clients: make(map[*Client]bool),
				}
			}

			h.rooms[client.BookingID].clients[client] = true
			h.mu.Unlock()

			log.Printf("[HUB] Client %v bergabung ke room %v", client.UserID, client.BookingID)

		
		case client := <-h.Unregister:
			h.mu.Lock()

			if room, exists := h.rooms[client.BookingID]; exists {
				
				delete(room.clients, client)
				close(client.Send) 

				
				if len(room.clients) == 0 {
					delete(h.rooms, client.BookingID)
					log.Printf("[HUB] Room %v dihapus (kosong)", client.BookingID)
				}
			}

			h.mu.Unlock()
			log.Printf("[HUB] Client %v keluar dari room %v", client.UserID, client.BookingID)

		
		case msg := <-h.Broadcast:
			h.mu.RLock()

			room, exists := h.rooms[msg.RoomID]

			if !exists {
				h.mu.RUnlock()
				continue
			}

			for client := range room.clients {
				select {
				case client.Send <- msg.Message:
					
				default:
				
					
					close(client.Send)
					delete(room.clients, client)
					log.Printf("[HUB] Client %v dihapus paksa (buffer penuh)", client.UserID)
				}
			}

			h.mu.RUnlock()
		}
	}
}
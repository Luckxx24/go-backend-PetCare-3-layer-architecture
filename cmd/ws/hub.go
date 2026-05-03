package ws

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

type ChatRoom struct {
	activeUsers map[*UserConnection]bool
	lock        sync.RWMutex
}

type MessageHub struct {
	chatRooms  map[uuid.UUID]*ChatRoom
	UserJoined chan *UserConnection
	UserLeft   chan *UserConnection
	NewMessage chan IncomingBroadcast
	lock       sync.RWMutex
}

type IncomingBroadcast struct {
	TargetRoomID uuid.UUID
	Payload      WsResponse
}

func NewMessageHub() *MessageHub {
	return &MessageHub{
		chatRooms:  make(map[uuid.UUID]*ChatRoom),
		UserJoined: make(chan *UserConnection),
		UserLeft:   make(chan *UserConnection),
		NewMessage: make(chan IncomingBroadcast),
	}
}

func (hub *MessageHub) Run() {
	for {
		select {

		case newUser := <-hub.UserJoined:
			hub.lock.Lock()

			roomBelumAda := hub.chatRooms[newUser.BookingID] == nil
			if roomBelumAda {
				hub.chatRooms[newUser.BookingID] = &ChatRoom{
					activeUsers: make(map[*UserConnection]bool),
				}
			}

			hub.chatRooms[newUser.BookingID].activeUsers[newUser] = true

			hub.lock.Unlock()
			log.Printf("[HUB] Pengguna %v bergabung ke room %v", newUser.UserID, newUser.BookingID)

		case disconnectedUser := <-hub.UserLeft:
			hub.lock.Lock()

			room, roomDitemukan := hub.chatRooms[disconnectedUser.BookingID]
			if roomDitemukan {

				delete(room.activeUsers, disconnectedUser)

				close(disconnectedUser.Send)

				roomSudahKosong := len(room.activeUsers) == 0
				if roomSudahKosong {
					delete(hub.chatRooms, disconnectedUser.BookingID)
					log.Printf("[HUB] Room %v dihapus karena sudah kosong", disconnectedUser.BookingID)
				}
			}

			hub.lock.Unlock()
			log.Printf("[HUB] Pengguna %v keluar dari room %v", disconnectedUser.UserID, disconnectedUser.BookingID)

		case broadcast := <-hub.NewMessage:
			hub.lock.RLock()

			targetRoom, roomDitemukan := hub.chatRooms[broadcast.TargetRoomID]
			if !roomDitemukan {
				hub.lock.RUnlock()
				continue
			}

			for user := range targetRoom.activeUsers {
				pesanTerkirim := kirimPesanKeUser(user, broadcast.Payload)
				if !pesanTerkirim {

					close(user.Send)
					delete(targetRoom.activeUsers, user)
					log.Printf("[HUB] Pengguna %v dihapus paksa (buffer penuh)", user.UserID)
				}
			}

			hub.lock.RUnlock()
		}
	}
}

func kirimPesanKeUser(user *UserConnection, pesan WsResponse) bool {
	select {
	case user.Send <- pesan:
		return true
	default:
		return false
	}
}

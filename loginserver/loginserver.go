package loginserver

import (
//"bytes"
	"log"
	"net"
//"strconv"

//"golang.org/x/crypto/bcrypt"
//"golang.org/x/crypto/bcrypt"

//"./serverpackets"
)

type LoginServer struct {
	//clients             []*models.Client
	//gameservers         []*models.GameServer
	//database            *mgo.Database
	//config              config.ConfigObject
	internalServersList []byte
	externalServersList []byte
	status              loginServerStatus
	//databaseSession     *mgo.Session
	clientsListener     net.Listener
	gameServersListener net.Listener
}

type loginServerStatus struct {
	successfulAccountCreation uint32
	failedAccountCreation     uint32
	successfulLogins          uint32
	failedLogins              uint32
	hackAttempts              uint32
}

//func New(cfg config.ConfigObject) *LoginServer {
//	return &LoginServer{config: cfg}
//}

func (l *LoginServer) init() {
	var err error

	// Connect to our database
	//l.databaseSession, err = mgo.Dial(l.config.LoginServer.Database.Host + ":" + strconv.Itoa(l.config.LoginServer.Database.Port))
	if err != nil {
		panic("Couldn't connect to the database server.")
	} else {
		log.Print("Successfully connected to the database server.")
	}

	// Select the appropriate database
	//l.database = l.databaseSession.DB(l.config.LoginServer.Database.Name)

	// Listen for client connections
	l.clientsListener, err = net.Listen("tcp", ":2106")
	if err != nil {
		log.Print("Couldn't initialize the Login Server (Clients listener).")
	} else {
		log.Print("Login Server listening for clients connections on port 2106.")
	}

	// Listen for game servers connections
	//l.gameServersListener, err = net.Listen("tcp", ":9413")
	//if err != nil {
	//	log.Println("Couldn't initialize the Login Server (Gameservers listener)")
	//} else {
	//	log.Println("Login Server listening for gameservers connections on port 9413")
	//}
}

func (l *LoginServer) Start() {
	//defer l.databaseSession.Close()
	defer l.clientsListener.Close()
	//defer l.gameServersListener.Close()

	done := make(chan bool)

	go func() {
		for {
			var err error
			//client := models.NewClient()
			//client.Socket, err = l.clientsListener.Accept()
			//l.clients = append(l.clients, client)
			if err != nil {
				log.Print("Couldn't accept the incoming connection.")
				continue
			} else {
				//go l.handleClientPackets(client)
			}
		}
		done <- true
	}()

	go func() {
		for {
			var err error
			//gameserver := models.NewGameServer()
			//gameserver.Socket, err = l.gameServersListener.Accept()
			//l.gameservers = append(l.gameservers, gameserver)
			if err != nil {
				log.Print("Couldn't accept the incoming connection.")
				continue
			} else {
				//go l.handleGameServerPackets(gameserver)
			}
		}

		done <- true
	}()

	for i := 0; i < 2; i++ {
		<-done
	}

}

//func (l *LoginServer) kickClient(client *models.Client) {
//	client.Socket.Close()
//
//	for i, item := range l.clients {
//		if bytes.Equal(item.SessionID, client.SessionID) {
//			copy(l.clients[i:], l.clients[i + 1:])
//			l.clients[len(l.clients) - 1] = nil
//			l.clients = l.clients[:len(l.clients) - 1]
//			break
//		}
//	}
//
//	log.Println("The client has been successfully kicked from the server.")
//}

//func (l *LoginServer) handleGameServerPackets(gameserver *models.GameServer) {
//	defer gameserver.Socket.Close()
//
//	for {
//		opcode, _, err := gameserver.Receive()
//
//		if err != nil {
//			log.Println(err)
//			log.Println("Closing the connection...")
//			break
//		}
//
//		switch opcode {
//		case 00:
//			log.Println("A game server sent a request to register")
//		default:
//			log.Println("Can't recognize the packet sent by the gameserver")
//		}
//	}
//}

//func (l *LoginServer) handleClientPackets(client *models.Client) {
//	log.Println("A client is trying to connect...")
//	defer l.kickClient(client)
//
//	buffer := serverpackets.NewInitPacket()
//	err := client.Send(buffer, false, false)
//
//	if err != nil {
//		log.Println(err)
//		return
//	} else {
//		log.Println("Init packet sent.")
//	}
//
//	for {
//		opcode, data, err := client.Receive()
//
//		if err != nil {
//			log.Println(err)
//			log.Println("Closing the connection...")
//			break
//		}
//
//		switch opcode {
//		case 00:
//			// response buffer
//			var buffer []byte
//
//			requestAuthLogin := clientpackets.NewRequestAuthLogin(data)
//
//			log.Printf("User %s is trying to login\n", requestAuthLogin.Username)
//
//			accounts := l.database.C("accounts")
//			err := accounts.Find(bson.M{"username": requestAuthLogin.Username}).One(&client.Account)
//
//			if err != nil {
//				if l.config.LoginServer.AutoCreate == true {
//					hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestAuthLogin.Password), 10)
//					if err != nil {
//						log.Println("An error occured while trying to generate the password")
//						l.status.failedAccountCreation += 1
//
//						buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_SYSTEM_ERROR)
//					} else {
//						client.Account = models.Account{
//							Id:          bson.NewObjectId(),
//							Username:    requestAuthLogin.Username,
//							Password:    string(hashedPassword),
//							AccessLevel: ACCESS_LEVEL_PLAYER}
//
//						err = accounts.Insert(&client.Account)
//						if err != nil {
//							log.Printf("Couldn't create an account for the user %s\n", requestAuthLogin.Username)
//							l.status.failedAccountCreation += 1
//
//							buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_SYSTEM_ERROR)
//						} else {
//							log.Printf("Account successfully created for the user %s\n", requestAuthLogin.Username)
//							l.status.successfulAccountCreation += 1
//
//							buffer = serverpackets.NewLoginOkPacket(client.SessionID)
//						}
//					}
//				} else {
//					log.Println("Account not found !")
//					l.status.failedLogins += 1
//
//					buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_USER_OR_PASS_WRONG)
//				}
//			} else {
//				// Account exists; Is the password ok?
//				err = bcrypt.CompareHashAndPassword([]byte(client.Account.Password), []byte(requestAuthLogin.Password))
//
//				if err != nil {
//					log.Printf("Wrong password for the account %s\n", requestAuthLogin.Username)
//					l.status.failedLogins += 1
//
//					buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_USER_OR_PASS_WRONG)
//				} else {
//
//					if client.Account.AccessLevel >= ACCESS_LEVEL_PLAYER {
//						l.status.successfulLogins += 1
//
//						buffer = serverpackets.NewLoginOkPacket(client.SessionID)
//					} else {
//						l.status.failedLogins += 1
//
//						buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_ACCESS_FAILED)
//					}
//
//				}
//			}
//
//			err = client.Send(buffer)
//
//			if err != nil {
//				log.Println(err)
//			}
//
//		case 02:
//			requestPlay := clientpackets.NewRequestPlay(data)
//
//			log.Printf("The client wants to connect to the server : %d\n", requestPlay.ServerID)
//
//			var buffer []byte
//			if len(l.config.GameServers) >= int(requestPlay.ServerID) && (l.config.GameServers[requestPlay.ServerID - 1].Options.Testing == false || client.Account.AccessLevel > ACCESS_LEVEL_PLAYER) {
//				if !bytes.Equal(client.SessionID[:8], requestPlay.SessionID) {
//					l.status.hackAttempts += 1
//
//					buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_ACCESS_FAILED)
//				} else {
//					buffer = serverpackets.NewPlayOkPacket()
//				}
//			} else {
//				l.status.hackAttempts += 1
//
//				buffer = serverpackets.NewPlayFailPacket(serverpackets.REASON_ACCESS_FAILED)
//			}
//			err := client.Send(buffer)
//
//			if err != nil {
//				log.Println(err)
//			}
//
//		case 05:
//			requestServerList := clientpackets.NewRequestServerList(data)
//
//			var buffer []byte
//			if !bytes.Equal(client.SessionID[:8], requestServerList.SessionID) {
//				l.status.hackAttempts += 1
//
//				buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_ACCESS_FAILED)
//			} else {
//				buffer = serverpackets.NewServerListPacket(l.config.GameServers, client.Socket.RemoteAddr().String())
//			}
//			err := client.Send(buffer)
//
//			if err != nil {
//				log.Println(err)
//			}
//
//		default:
//			log.Println("Couldn't detect the packet type.")
//		}
//	}
//}

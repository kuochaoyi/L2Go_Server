package loginserver

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net"
	"os"

	//"golang.org/x/crypto/bcrypt"

	"./crypt/blowfish"
)

func init() {
	modulus, key := generateRSA()

	ln, err := net.Listen("tcp", ":2106")
	defer ln.Close()

	if err != nil {
		log.Print("Couldn't initialize the Login Server")
	} else {
		log.Print("Login Server initialized.")
		log.Print("Listening on 127.0.0.1:2106.")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Couldn't accept the incoming connection.")
			continue
		}

		go handleConnection(conn, modulus, key)
	}
}

func handleConnection(conn net.Conn, modulus []byte, key *rsa.PrivateKey) {

	// Create the packet wrapper
	packet := []byte{0x00,
		0xfd, 0x8a, 0x22, 0x00, 0x5a, 0x78, 0x00, 0x00, // Header
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, // Fake RSA key modulus
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Unknown
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Inject our modulus
	for i := 0; i < len(modulus); i++ {
		packet[9 + i] = modulus[i]
	}

	length := len(packet) + 2
	buffer := make([]byte, length)

	buffer[0] = byte(length & 0xff)
	buffer[1] = byte((length >> 8) & 0xff)
	copy(buffer[2:], packet)

	log.Print("A client is trying to connect...")
	log.Printf("Created a an init packet[%d] = %X\n", len(buffer), buffer)

	log.Print("Sending the Init packet...")
	conn.Write([]byte(buffer))

	log.Print("Receiving the Init response")
	for {

		// Receive the packet header (size)
		header := make([]byte, 2)

		_, _ = conn.Read(header)

		// Calculate the packet size
		var size int = 0
		size = size + int(header[0])
		size = size + (int(header[1]) * 256)

		if size > 0 {

			log.Printf("Received a packet header...\n")
			log.Printf("Expected packet length: %d\n", size - 2)

			// Receive the content packet
			data := make([]byte, size - 2)

			n, _ := conn.Read(data)

			log.Printf("Actual packet length: %d\n", n)

			if n != size - 2 {
				log.Print("Packet size error !!")
			}

			log.Printf("Packet content : %X%X\n", header, data)

			decrypted := blowfishDecrypt(data, []byte(";5.]94-31==-%xT!^[$\000"), size - 2)
			log.Printf("Decrypted packet content : %X\n", decrypted)

			//decoded, _ := rsa.DecryptPKCS1v15(rand.Reader, key, decrypted)
			//log.Println(decoded)
		}

	}

}

func blowfishDecrypt(encrypted, key []byte, size int) []byte {
	// create the cipher
	dcipher, err := blowfish.NewCipher(key)
	if err != nil {
		// fix this. its okay for this tester program, but...
		panic(err)
	}

	count := len(encrypted) / 8

	decrypted := make([]byte, size)

	for i := 0; i < count; i++ {
		dcipher.Decrypt(decrypted[i * 8:i * 8 + 8], encrypted[i * 8:i * 8 + 8])
	}

	return decrypted
}

func generateRSA() ([]byte, *rsa.PrivateKey) {
	privatekey, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		log.Print(err.Error)
		os.Exit(1)
	}

	var publickey *rsa.PublicKey
	publickey = &privatekey.PublicKey
	scrambledModulus := publickey.N.Bytes() // modulus to bytes

	for i := 0; i < 4; i++ {
		temp := scrambledModulus[0x00 + i]
		scrambledModulus[0x00 + i] = scrambledModulus[0x4d + i]
		scrambledModulus[0x4d + i] = temp
	}

	// step 2 xor first 0x40 bytes with last 0x40 bytes
	for i := 0; i < 0x40; i++ {
		scrambledModulus[i] = byte(scrambledModulus[i] ^ scrambledModulus[0x40 + i])
	}

	// step 3 xor bytes 0x0d-0x10 with bytes 0x34-0x38
	for i := 0; i < 4; i++ {
		scrambledModulus[0x0d + i] = byte(scrambledModulus[0x0d + i] ^ scrambledModulus[0x34 + i])
	}

	// step 4 xor last 0x40 bytes with first 0x40 bytes
	for i := 0; i < 0x40; i++ {
		scrambledModulus[0x40 + i] = byte(scrambledModulus[0x40 + i] ^ scrambledModulus[i])
	}

	return scrambledModulus, privatekey
}

package test

import (
	"JCVM/core"
	"JCVM/jcre"
	"JCVM/jcre/nativeMethods"
	"encoding/hex"
	"fmt"
	"net"
	"testing"
)

func TestJcrePINBlocage(t *testing.T) {
	server, client := net.Pipe()
	go func() {
		// Do some stuff
		args := []string{ /*`../test/framework.ijc`,*/ `../test/lang.ijc`}
		for i := 0; i < len(args); i++ {
			dataBuffer := core.ReadInBuffer(args[i])
			core.Lst.PushBack(core.BuildApplet(dataBuffer))
		}
		nativeMethods.PowerUP(server)
		jcre.MainLoop()
	}()

	/*********************************************************************************/
	//powerUP the card
	fmt.Println("PowerUping the card ")
	sendPUp := []byte{0x01, 0x01}
	_, err := client.Write(sendPUp)
	if err != nil {
		t.Error(err)
	}
	respATR := make([]byte, 20)
	n, err := client.Read(respATR)
	if err != nil {
		t.Error(err)

	}
	dst := make([]byte, hex.EncodedLen(len(respATR[0:n])))
	hex.Encode(dst, respATR[0:n])
	fmt.Printf("Card ATR:%s\r\n", dst)

	//Install the applet
	fmt.Println("Installing the applet")
	installBuf := []byte{0x00, 0x01, 0x00, 0x00}
	_, err = client.Write(installBuf)

	if err != nil {
		t.Error(err)

	}
	SWInstall := make([]byte, 4)
	n, err = client.Read(SWInstall)
	if err != nil {
		t.Error(err)
	}
	if SWInstall[0] != 0x90 && SWInstall[1] != 0x00 {
		t.Error("Install failed")
	}
	dst = make([]byte, hex.EncodedLen(len(SWInstall[0:n])))
	hex.Encode(dst, SWInstall[0:n])
	fmt.Printf("Applet installation succeeded, SW: %s \r\n", dst)
	// write a message to server
	fmt.Println("Selecting applet ...")
	//0xD3:0x3E:0x1:0x16:0x50
	message := []byte{0x0, 0xA4, 0x00, 0x00, 0x5, 0xD3, 0x3E, 0x1, 0x16, 0x50}
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive message from server
	buffer := make([]byte, 20)
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x90 && buffer[1] != 0x00 {
		t.Error("Selection failed")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("Applet selection succeeded; SW: %s\n ", dst)

	/*********************************************************************************/
	//testing
	fmt.Println()

	//do verif
	//verify code PIN
	//first false code pin
	message = []byte{0x80, 0x20, 0x01, 0x00, 0x04, 0x31, 0x31, 0x31, 0x30, 0x7F}
	// Verifier le code PIN de l'utilisateur avec un faux code PIN
	fmt.Printf("verify user code PIN with a false PIN\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] == 0x90 && buffer[1] == 0x00 {
		t.Error("pin failed")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()

	//second false code pin
	message = []byte{0x80, 0x20, 0x01, 0x00, 0x04, 0x31, 0x31, 0x31, 0x30, 0x7F}
	fmt.Printf("verify user code PIN with a false PIN\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] == 0x90 && buffer[1] == 0x00 {
		t.Error("pin failed")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()

	//third false code pin
	message = []byte{0x80, 0x20, 0x01, 0x00, 0x04, 0x31, 0x31, 0x31, 0x30, 0x7F}
	fmt.Printf("verify user code PIN with a false PIN\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] == 0x90 && buffer[1] == 0x00 {
		t.Error("pin failed")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()

	//good  code pin
	message = []byte{0x80, 0x20, 0x01, 0x00, 0x04, 0x31, 0x31, 0x31, 0x31, 0x7F}
	fmt.Printf("verify user code PIN with the good PIN after 3 tries\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] == 0x90 && buffer[1] == 0x00 {
		t.Error("pin failed")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()

	//Verifier le code de l'emetteur
	message = []byte{0x80, 0x20, 0x02, 0x00, 0x08, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x7F}
	fmt.Printf("enter manufacturer code PIN \t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x90 && buffer[1] != 0x00 {
		t.Error("pin failed")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()

	// debloquer et mettre a jour le code PIN utilisateur
	message = []byte{0x80, 0x22, 0x00, 0x00, 0x04, 0x31, 0x32, 0x33, 0x34, 0x7F}
	fmt.Printf("deblock and update code pin (nbr tries)\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x90 && buffer[1] != 0x00 {
		t.Error("pin failed")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()

	// Verifier avec le nouveau code PIN
	message = []byte{0x80, 0x20, 0x01, 0x00, 0x04, 0x31, 0x32, 0x33, 0x34, 0x7F}
	fmt.Printf("verify code PIN with the new one\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x90 && buffer[1] != 0x00 {
		t.Error("pin failed")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()

	// Do some stuff
	client.Close()
}

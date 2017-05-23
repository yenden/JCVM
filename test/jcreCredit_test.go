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

func TestJcreCreditCard(t *testing.T) {
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
	//testing credit
	fmt.Println()
	//consult solde
	message = []byte{0x80, 0x04, 0x00, 0x00, 0x00, 0x7F}
	fmt.Printf("consult solde\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer) //0x61xx
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x61 {
		fmt.Println(buffer[:n])
		t.Error("Not 61xxvvv")
	}
	message = []byte{0x80, 0x04, 0x00, 0x00, 0x01} //send get response
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	dst1 := make([]byte, hex.EncodedLen(len(buffer[:n-2])))
	hex.Encode(dst1, buffer[:n-2])
	fmt.Printf("Solde: \t %s\t ", dst1)
	//status
	/*n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x90 && buffer[1] != 00 {
		t.Error("Error consult solde")
	}*/
	dst = make([]byte, hex.EncodedLen(len(buffer[n-2:n])))
	hex.Encode(dst, buffer[n-2:n])
	fmt.Printf("SW: %s\r\n ", dst)
	fmt.Println()
	//do credit
	// Crediter le porte monnaie avant envoie du code PIN
	message = []byte{0x80, 0x00, 0x00, 0x00, 0x01, 0x50, 0x7F}
	fmt.Printf("try to credit card with 0x50 before verifiying PIN\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] == 0x90 && buffer[1] == 00 {
		t.Error("Error consult solde")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()
	//do verif
	//verify code PIN
	message = []byte{0x80, 0x20, 0x02, 0x00, 0x08, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x7F}
	fmt.Printf("verify code PIN\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x90 && buffer[1] != 00 {
		t.Error("Error verify code")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()
	// do credit after verifying pin code
	message = []byte{0x80, 0x00, 0x00, 0x00, 0x01, 0x50, 0x7F}
	fmt.Printf("redo credit card with 0x50 \t ")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x90 && buffer[1] != 00 {
		t.Error("Error redo credit")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()
	//read solde after credit
	message = []byte{0x80, 0x04, 0x00, 0x00, 0x00, 0x7F}
	fmt.Printf("consult solde\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer) //61xx
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x61 {
		t.Error("Not 61xx")
	}
	message = []byte{0x80, 0x04, 0x00, 0x00, 0x01} //send get response
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n-2])))
	hex.Encode(dst, buffer[:n-2])
	fmt.Printf("Solde: \t %s\t ", dst)
	/*//status
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x90 && buffer[1] != 00 {
		t.Error("Error consult solde")
	}*/
	dst = make([]byte, hex.EncodedLen(len(buffer[n-2:n])))
	hex.Encode(dst, buffer[n-2:n])
	fmt.Printf("SW: %s\r\n ", dst)
	fmt.Println()
	// Crediter et lecture du porte monnaie apres code PIN
	// Test du depassement de la valeur Max 0x7F
	message = []byte{0x80, 0x00, 0x00, 0x00, 0x01, 0xB0, 0x7F}
	fmt.Printf("credit card with 0xB0 to test capacity overflow\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] == 0x90 && buffer[1] == 00 {
		t.Error("Error capcity overflow")
	}
	dst = make([]byte, hex.EncodedLen(len(buffer[:n])))
	hex.Encode(dst, buffer[:n])
	fmt.Printf("SW :%s\n ", dst)
	fmt.Println()
	//Lecture
	//read solde after credit
	message = []byte{0x80, 0x04, 0x00, 0x00, 0x00, 0x7F}
	fmt.Printf("consult solde\t")
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	// receive
	n, err = client.Read(buffer) //0x61xx
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x61 {
		t.Error("Error consult solde")
	}
	message = []byte{0x80, 0x04, 0x00, 0x00, 0x01} //send get response
	_, err = client.Write(message)
	if err != nil {
		t.Error(err)
	}
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	dst = make([]byte, hex.EncodedLen(len(buffer[:n-2])))
	hex.Encode(dst, buffer[:n-2])
	fmt.Printf("Solde: \t %s\t ", dst)
	/*//status
	n, err = client.Read(buffer)
	if err != nil {
		t.Error(err)
	}
	if buffer[0] != 0x90 && buffer[1] != 00 {
		t.Error("Error consult solde")
	}*/
	dst = make([]byte, hex.EncodedLen(len(buffer[n-2:n])))
	hex.Encode(dst, buffer[n-2:n])
	fmt.Printf("SW: %s\r\n ", dst)
	server.Close()
}

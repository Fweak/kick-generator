package main

import (
	Kick "app/internal/kick"
	"fmt"
	"time"
)

func main() {
	kick := Kick.CreateClient("temp@gmail.com", "password")
	kick.GetCookies()
	kick.StartSocket()
	defer kick.Conn.Close()
	kick.RequestTokenProvider()
	kick.SendEmail()
	kick.SendEmailCode("15231")
	username, err := kick.RegisterAccount("dr4inG4ng")
	if err != nil {
		fmt.Println("an error occured men...")
		fmt.Println(err.Error())
		return
	}
	fmt.Printf(`%s:%s:%s|%s`, username, kick.Email, kick.Password, time.Now().Format("01-02-2006 15:04:05"))
}

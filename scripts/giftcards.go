package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/c0nrad/mongobucks/models"
	"github.com/c0nrad/mongobucks/ticket"
	"github.com/pborman/uuid"
)

func main() {

	if models.MongoUri != "" {
		fmt.Println("[+] We are using prod!")
	}
	ticketCount := 1

	reward, err := models.NewReward("20 Mongobucks", "", true, models.RedeemReward, 20)
	if err != nil {
		panic(err)
	}

	for i := 0; i < ticketCount; i++ {
		t, err := models.NewTicket(reward.ID, "MongoDB Employee", reward.Name, uuid.New())
		if err != nil {
			panic(err)
		}

		img := ticket.GenerateTicketImage(t)
		SaveImage(img, "img"+strconv.Itoa(i)+".png")
	}

}

func SaveImage(baseImg image.Image, filename string) {
	outFile, err := os.Create("./ticket/" + filename)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, baseImg)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")

}

// func NewTicket(rewardId bson.ObjectId, username, rewardName, redemption string) (*Ticket, error) {

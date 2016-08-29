package ticket

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/c0nrad/mongobucks/models"
	"github.com/golang/freetype/truetype"

	"github.com/golang/freetype"
)

// func main() {
// 	t := models.Ticket{TS: time.Now(), Name: "1 Highfive from Andrew Erlichson", Redemption: "87cc0aa7-2155-42f6-b6e5-f5ee66b93b2e", IsUsed: false}
// 	GenerateTicket(&t)
// }

func GenerateTicketImage(ticket *models.Ticket) image.Image {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile("./ticket/luxisr.ttf")
	if err != nil {
		panic(err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	base, err := os.Open("./ticket/base.png")
	if err != nil {
		panic(err)
	}
	defer base.Close()
	baseImg, err := png.Decode(base)
	if err != nil {
		panic(err)
	}

	imgSet := baseImg.(draw.Image)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(32)
	c.SetClip(baseImg.Bounds())
	c.SetSrc(image.Black)
	c.SetDst(imgSet)

	fmt.Println(baseImg.Bounds())

	nameWidth := CalculateWidth(ticket.Name, font)
	pt := freetype.Pt(CenterStart(nameWidth), 214)
	_, err = c.DrawString(ticket.Name, pt)
	if err != nil {
		panic(err)
	}

	c.SetFontSize(18)
	pt = freetype.Pt(95, 300)
	_, err = c.DrawString("Mongobucks", pt)
	if err != nil {
		panic(err)
	}

	pt = freetype.Pt(285, 298)
	_, err = c.DrawString(ticket.TS.Format("1/2/06"), pt)
	if err != nil {
		panic(err)
	}

	redemptionUrl := "http://mongobucks.mongodb.cc/#/r/" + ticket.Redemption

	c.SetFontSize(12)
	pt = freetype.Pt(25, 369)
	_, err = c.DrawString(time.Now().Format(redemptionUrl), pt)
	if err != nil {
		panic(err)
	}

	// add barcode
	qrImg := GenerateQR(redemptionUrl)

	draw.Draw(imgSet, image.Rect(515, 275, 615, 375), qrImg, image.Point{X: 0, Y: 0}, draw.Src)

	// Save that RGBA image to disk.
	// outFile, err := os.Create("./ticket/out.png")
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// defer outFile.Close()
	// b := bufio.NewWriter(outFile)
	// err = png.Encode(b, baseImg)
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// err = b.Flush()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Wrote out.png OK.")

	return imgSet
}

func CalculateWidth(message string, font *truetype.Font) int {
	opts := truetype.Options{}
	opts.Size = 32
	face := truetype.NewFace(font, &opts)

	// Calculate the widths and print to image
	totalWidth := 0
	for _, x := range message {
		awidth, _ := face.GlyphAdvance(rune(x))
		totalWidth += awidth.Round()
	}

	return totalWidth
}

func GenerateQR(url string) image.Image {
	code, err := qr.Encode(url, qr.L, qr.Unicode)
	if err != nil {
		panic(err)
	}

	code, err = barcode.Scale(code, 100, 100)
	if err != nil {
		panic(err)
	}

	return code
}

func CenterStart(width int) int {
	return (640 - width) / 2
}

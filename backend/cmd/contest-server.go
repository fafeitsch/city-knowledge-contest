package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fafeitsch/city-knowledge-contest/backend/geodata"
	"github.com/fafeitsch/city-knowledge-contest/backend/keygen"
	"github.com/fafeitsch/city-knowledge-contest/backend/webapi"
	"github.com/urfave/cli/v2"
)

func rangeValidation(min int, max int, name string) func(*cli.Context, int) error {
	return func(ctx *cli.Context, value int) error {
		if value < min || value > max {
			return fmt.Errorf("value %d for flag \"%s\" is not withing range %d – %d", value, name, min, max)
		}
		return nil
	}
}

const keyLengthMessage = "Must be between 2 and 255. Lower values improve debugging but increase risk of key collisions (which the app does not handle well)."

var port int
var portFlag = &cli.IntFlag{
	Name:        "port",
	Aliases:     []string{"p"},
	Value:       23123,
	Usage:       "The port on which the game runs.",
	Destination: &port,
}
var allowCors bool
var allowCorsFlag = &cli.BoolFlag{
	Name:  "allowCors",
	Value: false, Usage: "If true, the backend sets CORS headers. " +
		"Should only be used for development and debugging.",
	Destination: &allowCors,
}
var playerKeyLength int
var playerKeyLengthFlag = &cli.IntFlag{
	Name:        "playerKeyLength",
	Value:       10,
	Usage:       "Number of chars used to identify players. " + keyLengthMessage,
	Action:      rangeValidation(2, 255, "playerKeyLength"),
	Destination: &playerKeyLength,
}
var roomKeyLength int
var roomKeyLengthFlag = &cli.IntFlag{
	Name:        "roomKeyLength",
	Value:       21,
	Usage:       "Number of chars used to identify rooms." + keyLengthMessage,
	Action:      rangeValidation(2, 255, "roomKeyLength"),
	Destination: &roomKeyLength,
}
var nominatimServer string
var nominatimServerFlag = &cli.StringFlag{
	Name:        "nominatimServer",
	Value:       "https://nominatim.openstreetmap.org",
	Usage:       "Base URL to the Nominatim backend API",
	Destination: &nominatimServer,
}
var tileServer string
var tileServerFlag = &cli.StringFlag{
	Name:        "tileServer",
	Value:       "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
	Usage:       "Base URL to the Tile backend API. Use {z}, {x}, {y} as placeholders (no \"$\").",
	Destination: &tileServer,
}

func main() {
	app := cli.App{
		Name:        "City Knowledge Context",
		Description: "A competitive quiz about your favorite city",
		Flags: []cli.Flag{
			portFlag,
			allowCorsFlag,
			playerKeyLengthFlag,
			roomKeyLengthFlag,
			nominatimServerFlag,
			tileServerFlag,
		},
		HideHelpCommand: true,
		Action: func(context *cli.Context) error {
			handler := webapi.New(webapi.Options{AllowCors: allowCors, TileServer: tileServer})
			keygen.SetPlayerKeyLength(playerKeyLength)
			keygen.SetRoomKeyLength(roomKeyLength)
			geodata.NominatimServer = nominatimServer
			log.Printf("Starting server on port %d", port)
			log.Printf("CORS mode enabled: %v", allowCors)
			log.Printf("Room key length set to %d", roomKeyLength)
			log.Printf("Player key length set to %d", playerKeyLength)
			log.Printf("Using Nominatim API backend at \"%s\"", geodata.NominatimServer)
			log.Printf("Using Tile API backend at \"%s\"", tileServer)
			go webapi.ClearTileCache()
			return http.ListenAndServe(":"+strconv.Itoa(port), handler)
		},
	}
	log.Fatal(app.Run(os.Args))
}

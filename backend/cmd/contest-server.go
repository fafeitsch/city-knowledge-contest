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
var useTileCache bool
var useTileCacheFlag = &cli.BoolFlag{
	Name:        "useTileCache",
	Value:       false,
	Usage:       "If true, OSM tiles will be cached in the ./file directory for 180 days.",
	Destination: &useTileCache,
}
var sslKey string
var sslKeyFlag = &cli.StringFlag{
	Name:        "sslKey",
	Value:       "",
	Usage:       "Path to the SSL key file.",
	Destination: &sslKey,
}
var sslCert string
var sslCertFlag = &cli.StringFlag{
	Name:        "sslCert",
	Value:       "",
	Usage:       "Path to the SSL certificate file.",
	Destination: &sslCert,
}
var imprintFile string
var imprintFileFlag = &cli.StringFlag{
	Name:        "imprintFile",
	Value:       "",
	Usage:       "Path to the imprint file.",
	Destination: &imprintFile,
}
var dataProtectionFile string
var dataProtectionFileFlag = &cli.StringFlag{
	Name:        "dataProtectionFile",
	Value:       "",
	Usage:       "Path to the data protection file.",
	Destination: &dataProtectionFile,
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
			useTileCacheFlag,
			sslCertFlag,
			sslKeyFlag,
			dataProtectionFileFlag,
			imprintFileFlag,
		},
		HideHelpCommand: true,
		Action: func(context *cli.Context) error {
			geodata.NominatimServer = nominatimServer
			handler := webapi.New(
				webapi.Options{
					AllowCors:          allowCors,
					TileServer:         tileServer,
					UseTileCache:       useTileCache,
					DataProtectionFile: dataProtectionFile,
					ImprintFile:        imprintFile,
				},
			)
			keygen.SetPlayerKeyLength(playerKeyLength)
			keygen.SetRoomKeyLength(roomKeyLength)
			log.Printf("Starting server on port %d", port)
			log.Printf("CORS mode enabled: %v", allowCors)
			log.Printf("Room key length set to %d", roomKeyLength)
			log.Printf("Player key length set to %d", playerKeyLength)
			log.Printf("Using Nominatim API backend at \"%s\"", geodata.NominatimServer)
			log.Printf("Using Tile API backend at \"%s\"", tileServer)
			log.Printf("Using Tile cache enabled: %v", useTileCache)
			log.Printf("Using data protection file at \"%s\"", dataProtectionFile)
			log.Printf("Using imprint file at \"%s\"", imprintFile)
			if sslKey != "" && sslCert != "" {
				log.Printf("SSL key file: %s", sslKey)
				log.Printf("SSL certificate file: %s", sslCert)
				return http.ListenAndServeTLS(":"+strconv.Itoa(port), sslCert, sslKey, handler)
			}
			log.Printf("Running in unencrypted mode, should only be used for development and debugging.")
			return http.ListenAndServe(":"+strconv.Itoa(port), handler)
		},
	}
	log.Fatal(app.Run(os.Args))
}

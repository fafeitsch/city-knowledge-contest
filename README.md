# City Knowledge Contest 
**Who knows finds the streets first?** *A PvP online game to get to know your city better.*

## The Game

Set up a room with your friends and find the required street first in the map of
your city! Answer fast – get more points!

The game is currently in alpha version.

## Requirements

In order to play, you need meet three non-technical requirements:

1. One or more list of streets. They must be provided in the following JSON format:
    ```json
   {
        "country": "Germany",
        "city": "Würzburg",
        "name": "Würzburg",
        "center": {
            "Lat": 9.933333,
            "Lng": 49.8
        },
        "streets":[
            {
                "Name": "Ingolstadter Hof",
                "coord": {
                    "Lat": 49.7944035,
                    "Lng": 9.9343492
                }
            },
            …
        ]
   }
   ```
   The `name` property donates the name of the list as shown in the game setup.
   The `streets` array contains the list of streets that might be asked in game. The coordinate
   contains the solution (i.e. where the street is). Instead of providing a coordinate object,
   you can provide `null`. Then, the app will query [Nominatim](https://nominatim.org/) for the solution.
   It will use `country` and `city` to refine the query. You can disable the Nominatim query. Then, or if
   the street is not found, no solution is shown after each round. Playing the game is 
   possible nonetheless because of the next step.
2. Access to an OSRM server. We recommend [hosting your own OSRM server](https://hub.docker.com/r/osrm/osrm-backend/).
   The OSRM server is used to check the players' answers for correctness.
3. Access to an Openstreet Map Tile server. Per default, the app uses the publicly available OSM tile servers,
   but you might run into quota limits here. We strongly recommend using your own OSM tile server or
   buy an API quota on one of the available public servers.

## Starting the server

The easiest way to start the server quickly is to use the `Dockerfile` to build
an image:

`docker build --rm --no-cache -t  city-contest:alpha .`
and then start the server:

`docker run -p 8081:23123 -it -v /path/to/streetlists:/streetlists city-contest:alpha`

As shown above, you have to specify the `streetlists` directory. You can pass the flag `--help` to the
docker container, which will print the list of available options (e.g. how to set your own OSRM and OSM tile servers).

The server is then available in the browser at `http://localhost:8081`

# License

MIT License

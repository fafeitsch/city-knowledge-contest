# City Knowledge Contest 
**Who knows finds the streets first?** *A PvP online game to get to know your city better.*

## The Game

Set up a room with your friends and find the required street first in the map of
your city! Answer fast – get more points!

The game is currently in beta version.

## Requirements

In order to play, you need meet three non-technical requirements:

1. One or more list of streets. They must be provided in the following JSON format:
    ```json
   {
        "country": "Germany",
        "city": "Würzburg",
        "name": "Würzburg",
        "center": {
            "lng": 9.933333,
            "lat": 49.8
        },
        "streets":["Ingolstadter Hof","Bahnhofstraße",…]
   }
   ```
   The `name` property donates the name of the list as shown in the game setup.
   The `streets` array contains the list of streets that might be asked in game: Opon every round, the app will randomly choose a street and 
   and query [Nominatim](https://nominatim.org/) for the solution.
   It will use `country` and `city` to refine the query. 
2. Access to a Nominatim server. We recommend [hosting your own Nominatim server](https://github.com/mediagis/nominatim-docker).
   The Nominatim server is used to get the solution for the randomly chosen street, as well as for checking the players' answers.
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
docker container, which will print the list of available options (e.g. how to set your own Nominatim and OSM tile servers).

The server is then available in the browser at `http://localhost:8081`

## Contributing

We are happy for contributions. You can either send a pull request for an idea of yours, or you can browse the open
issues for a suitable task. 

# License

MIT License

# spotauth
A CLI for generating and saving Spotify auth/refresh tokens to a file

# Usage
```
~ $ spotauth --help
Usage of spotauth:
  -out_file string
    	Path to file where Spotify auth token will be saved - default: spotify.json (default "spotify.json")
  -redirect_uri string
    	Redirect URI for Spotify app - default: http://localhost:8080/callback (default "http://localhost:8080/callback")
  -spotify_id string
    	Spotify app ID
  -spotify_secret string
    	Spotify app secret
```

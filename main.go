package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	outFilePathArg   string = "out_file"
	redirectURIArg   string = "redirect_uri"
	spotifyIDArg     string = "spotify_id"
	spotifySecretArg string = "spotify_secret"
)

var (
	outFilePath   string
	redirectURI   string
	spotifyID     string
	spotifySecret string
)

func init() {
	flag.StringVar(&outFilePath, outFilePathArg, "spotify.json", "Path to file where Spotify auth token will be saved - default: spotify.json")
	flag.StringVar(&redirectURI, redirectURIArg, "http://localhost:8080/callback", "Redirect URI for Spotify app - default: http://localhost:8080/callback")
	flag.StringVar(&spotifyID, spotifyIDArg, "", "Spotify app ID")
	flag.StringVar(&spotifySecret, spotifySecretArg, "", "Spotify app secret")
	flag.Parse()

	requiredArgs := map[string]string{outFilePathArg: outFilePath, redirectURIArg: redirectURI, spotifyIDArg: spotifyID, spotifySecretArg: spotifySecret}

	for argName, argVal := range requiredArgs {
		if argVal == "" {
			log.Fatalf("Arg '%s' is missing value", argName)
		}
	}
}

func main() {

	auth := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)
	auth.SetAuthInfo(spotifyID, spotifySecret)

	ch := make(chan *spotify.Client)

	http.HandleFunc("/callback", makeCompleteAuthHandler(auth, ch))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL("")
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-ch

	_, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Login successful! Token info written to:", outFilePath)
}

func makeCompleteAuthHandler(auth spotify.Authenticator, ch chan *spotify.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token("", r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}

		writeAuthTokenToFile(outFilePath, token)

		client := auth.NewClient(token)
		fmt.Fprintf(w, "Login Completed!")
		ch <- &client
	}
}

func writeAuthTokenToFile(filePath string, token *oauth2.Token) {
	j, _ := json.Marshal(token)
	err := ioutil.WriteFile(filePath, j, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quantinium03/lucy/internal/config"
	"github.com/quantinium03/lucy/internal/database"
	"github.com/quantinium03/lucy/internal/database/model"
	"gorm.io/gorm"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Artist struct {
	Text string `json:"#text"`
}

type Track struct {
	Artist Artist `json:"artist"`
	Name   string `json:"name"`
}

type RecentTracks struct {
	Track []Track `json:"track"`
}

const TRACK_SEARCH_BASE_URI = "https://api.spotify.com/v1/search?q="
const SPOTIFY_OEMBED_URI = "https://open.spotify.com/oembed"
const SPOTIFY_AUTH_BASE_URI = "https://accounts.spotify.com/api/token"
var API_KEY = config.Config("LAST_FM_API_KEY")
var LASTFM_URL = "https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=Quantinium_X&api_key="+ API_KEY +"&limit=1&format=json"

func FetchSpotifyData() {
	ticker := time.NewTicker(4 * time.Minute)
	defer ticker.Stop()
	db := database.DB.DB
	for {
		select {
		case <-ticker.C:
			if err := getCurrentlyPlayingTrack(db); err != nil {
				log.Printf("Couldn't fetch the current playing track: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}
		}
	}
}

func getCurrentlyPlayingTrack(db *gorm.DB) error {
	var spotify model.Spotify
	if err := db.First(&spotify, "username = ?", []byte("quantinium")).Error; err != nil {
		return fmt.Errorf("Failed to fetch spotify data: %v", err)
	}

	var accessToken = spotify.SpotifyAccessToken

	trackHtml, trackEmbedUri, err := getCurrentEmbed(accessToken)
	if err != nil {
		log.Println("Couldnt get the track html and uri")
		return err
	}

	var spotifyData model.Spotify
	if err := db.First(&spotifyData, "username = ?", "quantinium").Error; err != nil {
		return err
	}

	spotifyData.SpotifyTrackEmbedHtml = trackHtml
	spotifyData.SpotifyTrackEmbedURI = trackEmbedUri

	if err := db.Save(&spotifyData).Error; err != nil {
		return err
	}

	return nil
}

func getCurrentEmbed(accessToken string) (string, string, error) {
	uri, err := getCurrentTrackUri(accessToken)
	if err != nil {
		return "", "", err
	}

	encodedUri := url.QueryEscape(uri)
	agent := fiber.Get(fmt.Sprintf("%s?url=%s", SPOTIFY_OEMBED_URI, encodedUri))
	code, body, errs := agent.Bytes()
	if code > 300 {
		log.Println("Failed to get the spotify embed")
		return "", "", fmt.Errorf("Error in send in post request for accessToken %v with code %v", errs, code)
	}
	if len(errs) > 0 {
		log.Println("Couldn't get spotify embed body")
		return "", "", fmt.Errorf("Error in send in post request for accessToken %v with code %v", errs, code)
	}

	var data fiber.Map
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Couldnt marshal the spotify oembed body")
		return "", "", err
	}

	trackUri, ok := data["iframe_url"].(string)
	if !ok {
		log.Println("Error: URI field is missing or not a string")
		return "", "", fmt.Errorf("iframe_url field is missing or not a string")
	}

	trackHtml, ok := data["html"].(string)
	if !ok {
		log.Println("Error: html field missing or not a string")
		return "", "", fmt.Errorf("html field is missing or not a string")
	}

	return trackHtml, trackUri, nil
}

func getSpotifyTrackUriFromLastFM() (string, error) {
	resp, err := http.Get(LASTFM_URL)
	if err != nil {
		return "", fmt.Errorf("Error making request to LastFM: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("Non-successful response from LastFM: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("Error decoding JSON: %v", err)
	}

	recentTracks, ok := data["recenttracks"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Error: 'recenttracks' field is missing or invalid")
	}

	tracks, ok := recentTracks["track"].([]interface{})
	if !ok {
		return "", fmt.Errorf("Error: 'track' field is missing or invalid")
	}

	if len(tracks) == 0 {
		return "", fmt.Errorf("Error: No tracks found in the response")
	}

	track, ok := tracks[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Error: Track data is invalid")
	}

	name, ok := track["name"].(string)
	if !ok {
		return "", fmt.Errorf("Error: 'name' field is missing or invalid")
	}

	artistData, ok := track["artist"].(map[string]interface{})
	if !ok {
		log.Fatalf("Error: 'artist' field is missing or invalid")
	}

	artistName, ok := artistData["#text"].(string)
	if !ok {
		log.Fatalf("Error: 'artist.#text' field is missing or invalid")
	}

	uri := fmt.Sprintf("track:%s artist:%s", strings.ToLower(name), strings.ToLower(artistName))
	result := (TRACK_SEARCH_BASE_URI + url.QueryEscape(uri) + "&type=track")
	return result, nil
}

func getSpotifySongUri(spotifyTrackUri string, accessToken string) (string, error) {
	req, err := http.NewRequest("GET", spotifyTrackUri, nil)
	if err != nil {
		return "", fmt.Errorf("Couldnt create new request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error: failed to get Spotify embed URI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("non-successful response from Spotify: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("error decoding JSON: %v", err)
	}

	tracks, ok := data["tracks"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("tracks field is missing or invalid")
	}

	items, ok := tracks["items"].([]interface{})
	if !ok {
		return "", fmt.Errorf("items array is invalid or missing")
	}

	if len(items) == 0 {
		return "", fmt.Errorf("no items found in the response")
	}

	firstItem, ok := items[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("first item in the items array is not a valid object")
	}

	externalUrls, ok := firstItem["external_urls"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("external_urls field is missing or invalid")
	}

	spotifyUri, ok := externalUrls["spotify"].(string)
	if !ok {
		return "", fmt.Errorf("could not find the Spotify URI in external_urls")
	}

	return spotifyUri, nil
}

func getCurrentTrackUri(accessToken string) (string, error) {
	lastFmTrackUriToSearchSpotifyTrack, err := getSpotifyTrackUriFromLastFM()
	println("Spotify to search: last fm shit ", lastFmTrackUriToSearchSpotifyTrack)
	if err != nil {
		return "", err
	}
	spotifyTrackUri, err := getSpotifySongUri(lastFmTrackUriToSearchSpotifyTrack, accessToken)
	println("Original track url:", spotifyTrackUri)
	if err != nil {
		return "", err
	}
	return spotifyTrackUri, nil
}

func GetAccessToken() {
	db := database.DB.DB
	ticker := time.NewTicker(3600 * time.Second)
	defer ticker.Stop()

	refreshToken := func() error {
		auth := base64.StdEncoding.EncodeToString([]byte(config.Config("CLIENT_ID") + ":" + config.Config("CLIENT_SECRET")))
		data := "grant_type=client_credentials"
		reqBody := bytes.NewBufferString(data)

		req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", reqBody)
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Basic "+auth)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %w", err)
		}

		var token TokenResponse
		if err := json.Unmarshal(body, &token); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		var spotify model.Spotify
		if err := db.First(&spotify, "username = ?", "quantinium").Error; err != nil {
			return fmt.Errorf("failed to fetch spotify record: %w", err)
		}

		spotify.SpotifyAccessToken = token.AccessToken
		if err := db.Save(&spotify).Error; err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}

		log.Println("Successfully refreshed Spotify access token")
		return nil
	}

	if err := refreshToken(); err != nil {
		log.Printf("Initial token refresh failed: %v", err)
	}

	for {
		<-ticker.C
		if err := refreshToken(); err != nil {
			log.Printf("Failed to refresh token: %v", err)
		}
	}
}

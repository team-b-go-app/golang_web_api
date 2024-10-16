package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// 映画情報を格納する構造体
type Film struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	OriginalTitle string   `json:"original_title"`
	Description   string   `json:"description"`
	Director      string   `json:"director"`
	Producer      string   `json:"producer"`
	ReleaseDate   string   `json:"release_date"`
	RunningTime   string   `json:"running_time"`
	RTScore       string   `json:"rt_score"`
	People        []string `json:"people"`
	Species       []string `json:"species"`
	Locations     []string `json:"locations"`
	Vehicles      []string `json:"vehicles"`
	URL           string   `json:"url"`
}

type Character struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Gender     string   `json:"gender"`
	Age        string   `json:"age"`
	Eye_color  string   `json:"eye_color"`
	Hair_color string   `json:"hair_color"`
	Films      []string `json:"films"`
}

// グローバル変数として映画のリストを保持
var films []Film

var characters []Character

// Studio Ghibli APIから映画データを取得
func fetchFilms() {
	resp, err := http.Get("https://ghibliapi.vercel.app/films")
	if err != nil {
		log.Fatalf("データの取得に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("レスポンスの読み取りに失敗しました: %v", err)
	}

	err = json.Unmarshal(body, &films)
	if err != nil {
		log.Fatalf("JSONのパースに失敗しました: %v", err)
	}
}

func fetchCharacters() {
	resp, err := http.Get("https://ghibliapi.vercel.app/people")
	if err != nil {
		log.Fatalf("キャラクターデータの取得に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("レスポンスの読み取りに失敗しました: %v", err)
	}

	err = json.Unmarshal(body, &characters)
	if err != nil {
		log.Fatalf("JSONのパースに失敗しました: %v", err)
	}
}

// CORSを許可するためのミドルウェア
func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // 全てのオリジンを許可
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// 全ての映画を取得するハンドラ
func getAllFilms(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r) // CORSを許可
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(films)
}

// original_titleで映画を検索するハンドラ
func searchFilms(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r) // CORSを許可
	query := r.URL.Query().Get("title")
	var result []Film

	for _, film := range films {
		if strings.Contains(strings.ToLower(film.OriginalTitle), strings.ToLower(query)) {
			result = append(result, film)
			fmt.Println(film.People)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// 映画のIDでキャラクターを検索するハンドラ
func searchCharactersByFilm(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r) // CORSを許可
	query := r.URL.Query().Get("title")
	fmt.Println(query)
	var filmID string

	// 映画のIDを取得
	for _, film := range films {
		if strings.Contains(strings.ToLower(film.OriginalTitle), strings.ToLower(query)) {
			filmID = film.ID
			break
		}
	}

	if filmID == "" {
		http.Error(w, "映画が見つかりません", http.StatusNotFound)
		return
	}

	// 映画IDに基づいてキャラクターをフィルタリング
	var result []Character
	for _, character := range characters {
		for _, film := range character.Films {
			// filmからIDを抽出
			filmParts := strings.Split(film, "/")
			if len(filmParts) > 0 {
				extractedID := filmParts[len(filmParts)-1]
				if extractedID == filmID {
					result = append(result, character)
					break
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ルーティングの設定
func handleRequests() {
	http.HandleFunc("/films", getAllFilms)
	http.HandleFunc("/search", searchFilms)
	http.HandleFunc("/search/characters", searchCharactersByFilm) // 新しいエンドポイント
	log.Println("サーバーをポート8080で開始します...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	// サーバー起動時に映画データを取得
	fetchFilms()
	fetchCharacters()
	handleRequests()
}

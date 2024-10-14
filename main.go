package main

import (
    "encoding/json"
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

// グローバル変数として映画のリストを保持
var films []Film

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
        }
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

// ルーティングの設定
func handleRequests() {
    http.HandleFunc("/films", getAllFilms)
    http.HandleFunc("/search", searchFilms)
    log.Println("サーバーをポート8080で開始します...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
    fetchFilms()
    handleRequests()
}

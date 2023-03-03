package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strings"
)

var (
	helpInfo string = ("Для работы с плейлистом доступны следующие функции\n" +
		"playlist 		- показать список песен в плейлисте (и информацию о них) \n" +
		"play     		- воспроизвести песню, изначально запуск начинается с первой песни \n" +
		"pause    		- остановить воспроизведение песни\n" +
		"next     		- включить следующую песню\n" +
		"prev     		- включить предыдущую песню\n" +
		"add?song_id=0		- добавить песню с указанным id \n" +
		"all_songs 		- посмотреть список всех имеющихся песен в базе данных с их id , именем и длительностью\n" +
		"del?song_id=0 		- удалить песню с указанным id из плейлиста \n" +
		"help 			- вызов справки")

	DBPath = "data/db/DB.db"
)

func main() {

	var playList Playlist
	InitPlaylist(1, &playList)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleRequest(w, r, &playList)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}
}

func HandleRequest(w http.ResponseWriter, r *http.Request, playlist *Playlist) {
	switch r.URL.Path {
	case "/play":
		fmt.Fprintf(w, playlist.Play())

	case "/pause":
		fmt.Fprintf(w, playlist.Pause())

	case "/next":
		fmt.Fprintf(w, playlist.Next())

	case "/prev":
		fmt.Fprintf(w, playlist.Prev())

	case "/add":
		songId := r.URL.Query().Get("song_id")
		fmt.Fprintf(w, playlist.AddSong(songId))

	case "/all_songs":
		fmt.Fprintf(w, PrintAllSong())

	case "/playlist":
		fmt.Fprintf(w, playlist.PlaylistList())

	case "/del":
		songId := r.URL.Query().Get("song_id")
		fmt.Fprintf(w, playlist.DelSong(songId))
	case "/help":
		fmt.Fprintf(w, helpInfo)
	default:
		fmt.Fprintf(w, "Unknow command")
	}
}

//Initialisation playlist struct
func InitPlaylist(userId int, playlist *Playlist) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//получаем список песен в плейлисте у данного юзера

	var songsList sql.NullString
	err = db.QueryRow("SELECT songs_list FROM users WHERE user_id = 1").Scan(&songsList)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка на NULL
	if !songsList.Valid {
		log.Println("Значение поля songsList равно NULL")
		return
	}

	songs := strings.Split(songsList.String, " ")
	for _, val := range songs {
		playlist.AddSongToNode(NewSong(val, db))
	}

}

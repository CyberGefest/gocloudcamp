package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

//Import songs from database to chain Node
func ImportSongFromDB(idSong string, db *sql.DB, p *Playlist) (string, error) {
	//берем данные песни из бд из таблицы

	var name string
	var duration int
	err := db.QueryRow("SELECT songName, duration FROM songs WHERE song_id=?", idSong).Scan(&name, &duration)
	if err != nil {
		return "", err
	}
	//добавляем в ноду
	p.AddSongToNode(&Song{
		songId:   idSong,
		duration: duration,
		songName: name,
	})
	return name, nil
}

//Check playlist for empty in DB
func PlaylistIsEmpty() bool {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	var songsList sql.NullString
	err = db.QueryRow("SELECT songs_list FROM users WHERE user_id = 1").Scan(&songsList)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка на NULL
	if !songsList.Valid {
		return true
	}
	return false
}

//Create new song from database info
func NewSong(songId string, db *sql.DB) *Song {
	rows, err := db.Query("SELECT songName, duration FROM songs WHERE song_id = $1", songId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var songName string
	var duration int
	for rows.Next() {
		err = rows.Scan(&songName, &duration)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &Song{
		songId:   songId,
		duration: duration,
		songName: songName,
	}
}

//check for repeat song in playlist
func SongReply(idSong string) bool {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	var songs string
	err = db.QueryRow("SELECT songs_list FROM users WHERE user_id=1").Scan(&songs)
	if err != nil {
		log.Fatal(err)
	}

	sliceSongs := strings.Split(songs, " ")

	for _, val := range sliceSongs {
		if val == idSong {
			return true
		}
	}
	return false
}

//Print list all songs in database
func PrintAllSong() string {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM songs")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var data string
	for rows.Next() {
		var id int
		var songName string
		var duaration int

		if err := rows.Scan(&id, &songName, &duaration); err != nil {
			panic(err)
		}
		data += fmt.Sprintf("id:%d, songName: %s, duration: %d sec.\n", id, songName, duaration)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return data
}

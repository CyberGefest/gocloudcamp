package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

type PlaylistI interface {
	Play() string
	Pause() string
	AddSong(idSong string) string
	Next() string
	Prev() string
}

type Playlist struct {
	head     *Node
	tail     *Node
	cur      *Node
	playback bool
	pausing  bool
	songStop bool
	mu       sync.Mutex
}

//Play song or start work with song
func (p *Playlist) Play() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.head == nil {
		return "Плейлист пустой. Добавте песню в плейлист."
	}

	if !p.playback {
		p.playback = true
		p.songStop = false
		go p.Playback()
		return "Воспроивзедение началось"
	}
	if p.pausing {
		p.pausing = false
		return "Воспроивзедение продолжено"
	}

	return ""
}

//Activation the song mode
func (p *Playlist) Playback() {
	curTimeSong := 0
	for curTimeSong < p.cur.Value.duration {
		switch p.songStop {
		case true:
			p.playback = false
			go p.Play()
			return
		default:
			if !p.pausing {
				fmt.Println(p.cur.Value.songName, " :", curTimeSong)
				time.Sleep(time.Second)
				curTimeSong++
			}
		}

	}
	if p.cur == p.tail {
		p.playback = false
	}
	if p.cur.Next != nil {
		p.cur = p.cur.Next
		p.playback = false
		p.Play()
	}
}

// Pause song
func (p *Playlist) Pause() string {

	if !p.pausing {
		p.pausing = true
		return "Воспроизведение остановлено"
	}
	return ""
}

//Play next song
func (p *Playlist) Next() string {

	if p.cur.Next != nil {
		p.cur = p.cur.Next
		p.songStop = true
		p.pausing = false
		go p.Play()
		return "Включаем следующую песню"
	} else {
		return "Следующей песни в плейлисте не существует"
	}
}

//Play prev song
func (p *Playlist) Prev() string {
	if p.cur.Prev != nil {
		p.cur = p.cur.Prev
		p.songStop = true
		p.pausing = false
		go p.Play()
		return "Включаем предыдущую песню песню"
	} else {
		return "Предыдущей песни в плейлисте не существует"
	}
}

//Add song to Node
func (p *Playlist) AddSongToNode(s *Song) {
	node := &Node{
		Value: s,
	}
	if p.head == nil {
		p.head = node
		p.tail = node
		p.cur = node
	} else {
		node.Prev = p.tail
		p.tail.Next = node
		p.tail = node
	}
}

//Add song to playlist and database
func (p *Playlist) AddSong(idSong string) string {
	//открываем базу данных
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	//проверка на пустой плейлист
	if PlaylistIsEmpty() {
		// ищем песню в бд
		name, err := ImportSongFromDB(idSong, db, p)
		if err != nil {
			return "Ошибка при обращении к базе данных " + err.Error()
		}

		//добавляем в базу данных песню
		_, err = db.Exec("UPDATE users SET songs_list = $1 WHERE user_id = 1", idSong)
		if err != nil {
			return "Ошибка при обращении к базе данных " + err.Error()
		}

		fmt.Println("Песня: " + name + "- успешно добавлена")
		return "Песня: " + name + "- успешно добавлена"
	}

	//проверяем на наличие повторений песен в БД
	if SongReply(idSong) {
		fmt.Println("Песня уже есть в плейлисте")
		return "Песня уже есть в плейлисте"
	}

	//импортируем песни из базы данных в Node
	name, err := ImportSongFromDB(idSong, db, p)
	if err != nil {
		return "Ошибка при обращении к базе данных " + err.Error()
	}

	//добавляем в БД
	_, err = db.Exec("UPDATE users SET songs_list = songs_list || ' '|| $1 WHERE user_id = 1", idSong)
	if err != nil {
		return "Ошибка при обращении к базе данных " + err.Error()
	}

	fmt.Println("Песня: ", name, "- успешно добавлена")
	return "Песня: " + name + "- успешно добавлена"
}

//Print all song in playlist
func (p *Playlist) PlaylistList() string {
	// пройтись с первого по последнюю ноду и записать все имена.
	var playlistList string
	for node := p.head; node != nil; node = node.Next {
		playlistList += "id:" + node.Value.songId + "  " + node.Value.songName + "\n"
	}
	return playlistList
}

//Delete song from playlist and DB
func (p *Playlist) DelSong(idSong string) string {

	//проверить играет ли песня в данный момент времени
	if p.cur.Value.songId == idSong && p.playback && !p.pausing {
		return "Нельзя удалить из плейлиста активный трек"
	}

	//проверить есть ли в плейлисте песня с текущим id, вернуть ноду с данной песней
	nodeExist := isNodeExists(p.head, idSong)
	if nodeExist == nil {
		return "Песни с таким id в вашем плейлисте нет"
	}

	//удалить ноду из плейлиста
	DelNodeFromPlaylist(p, nodeExist)

	//удалить песню из базы данных

	//получить список всех песен в плейлисте
	songsIdList := AllSongsId(p.head)

	//Добавить в бд новое значение
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("UPDATE users SET songs_list = $1 WHERE user_id = 1", songsIdList)
	if err != nil {
		return "Ошибка при обращении к базе данных " + err.Error()
	}

	return "Удалено"
}

package main

import (
	"fmt"
	"strings"
)

type Song struct {
	songId   string
	duration int
	songName string
}

type Node struct {
	Value *Song
	Prev  *Node
	Next  *Node
}

//Get Last Node
func lastNode(head *Node) *Node {
	// Проходим по списку, начиная с головы
	curr := head
	for curr != nil && curr.Next != nil {
		curr = curr.Next
	}
	return curr
}

//return "nil" if song not exist. Return pointer to Node if song exist
func isNodeExists(head *Node, val string) *Node {
	// Проходим по списку, начиная с головы
	curr := head
	for curr != nil {
		// Если значение текущего узла совпадает с искомым значением
		if curr.Value.songId == val {
			return curr
		}
		curr = curr.Next
	}
	return nil
}

//Removes a node from a doubly linked list
func deleteNode(head *Node, delNode *Node) *Node {

	//если удаляемый узел является единственным в списке
	if delNode == head && delNode.Next == nil {
		fmt.Println("last")
		return nil
	}
	// Если удаляемый узел является головой списка
	if delNode == head {
		head = delNode.Next
	}

	// Если удаляемый узел является концом списка
	if delNode.Next == nil {
		delNode.Prev.Next = nil
	} else {
		delNode.Next.Prev = delNode.Prev
	}

	// Обновляем ссылки на предыдущий и следующий узлы
	if delNode.Prev != nil {
		delNode.Prev.Next = delNode.Next
	}

	return head
}

//Get all id contain in playlist
func AllSongsId(head *Node) string {
	idList := make([]string, 0)
	curr := head
	for curr != nil {
		idList = append(idList, curr.Value.songId)
		curr = curr.Next
	}
	res := strings.Join(idList, " ")
	return res
}

//Delete node from playlist
func DelNodeFromPlaylist(p *Playlist, node *Node) {
	//если мы стоим на песне которую необходимо удалить,
	if p.cur == node {
		//если песня последняя, то текущей песней станет прошлая в списке, иначе та которая идет после удаляемой
		if p.tail == node && node != p.head {
			if p.tail.Prev != nil {
				p.cur = p.tail.Prev
				p.songStop = true
			}
		} else if p.cur.Next != nil {
			p.cur = p.cur.Next
			p.songStop = true
		}
	}

	p.head = deleteNode(p.head, node)
	p.tail = lastNode(p.head)
}

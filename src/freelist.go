package BTreeKVDB

import (
	"encoding/binary"
	"fmt"
)

type FreelistNode struct {
	StartId uint32
	EndId   uint32
	Next    *FreelistNode
	Pre     *FreelistNode
}
type Freelist struct {
	Size int //node num
	Head *FreelistNode
}

func (f *Freelist) Init() {
	f.Head = &FreelistNode{}
	f.Head.Pre = nil
	f.Size = 0
	f.Head.Next = nil
}
func (f *Freelist) Insert(stId int, edId int) {
	tmp := f.Head.Next
	StartId := (uint32)(stId)
	EndId := (uint32)(edId)
	for tmp != nil {
		if tmp.StartId == EndId+1 {
			tmp.StartId = StartId
			return
		}
		if tmp.EndId == StartId-1 {
			tmp.EndId = EndId
			return
		}
		tmp = tmp.Next
	}
	tmp = &FreelistNode{}
	tmp.Next = f.Head.Next
	tmp.Pre = f.Head
	f.Size += 1
	tmp.StartId = StartId
	tmp.EndId = EndId
	if f.Head.Next != nil {
		f.Head.Next.Pre = tmp
	}
	f.Head.Next = tmp
	f.Head.Pre = nil
}
func (f *Freelist) Get(N int) (int, int) {

	n := (uint32)(N)
	tmp := f.Head.Next
	for tmp != nil {
		c := (uint32)(tmp.EndId - tmp.StartId + 1)
		if c >= n {
			//fmt.Println("TTTT ", c, n)
			st := tmp.StartId
			ed := tmp.StartId + n - 1
			if c <= (uint32)(n) {
				//fmt.Println("XXXX", tmp.StartId, tmp.EndId)
				tmp.Pre.Next = tmp.Next
				//fmt.Println("CCCCC", tmp.Pre.StartId, tmp.Pre.EndId)
				if tmp.Next != nil {
					tmp.Next.Pre = tmp.Pre
				}
				f.Size -= 1
				//sdelete(tmp)
			}
			tmp.StartId += n
			return (int)(st), (int)(ed)
		}
		tmp = tmp.Next
	}
	return -1, -1
}
func (f *Freelist) show() {
	tmp := f.Head.Next
	fmt.Print("Show")
	for tmp != nil {
		if tmp.Pre == nil {
			fmt.Println(tmp.StartId, tmp.EndId, "PPPPPP")
		}
		fmt.Print(tmp.StartId, tmp.EndId, " ")
		tmp = tmp.Next
	}
	fmt.Println("")
}
func (f *Freelist) Serealize() []byte {
	data := make([]byte, 4096*4)
	tmp := 0
	var now *FreelistNode
	now = f.Head.Next
	binary.BigEndian.PutUint32(data[0:4], (uint32)(f.Size))
	tmp += 4
	b := make([]byte, 4)
	for now != nil {
		binary.BigEndian.PutUint32(b, now.StartId)
		copy(data[tmp:tmp+4], b)
		tmp += 4
		binary.BigEndian.PutUint32(b, now.EndId)
		copy(data[tmp:tmp+4], b)
		tmp += 4
		now = now.Next
	}
	return data
}
func (f *Freelist) Deserealize(buffer []byte) {
	tmp := 0
	now := f.Head
	f.Size = (int)(binary.BigEndian.Uint32(buffer[0:4]))
	tmp += 4
	for i := 0; i < f.Size; i++ {
		st := binary.BigEndian.Uint32(buffer[tmp : tmp+4])
		tmp += 4
		ed := binary.BigEndian.Uint32(buffer[tmp : tmp+4])
		tmp += 4
		var nNode *FreelistNode
		nNode = &FreelistNode{}
		nNode.StartId = st
		nNode.EndId = ed

		now.Next = nNode
		nNode.Pre = now
		nNode.Next = nil

		now = nNode
	}
}


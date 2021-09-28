package BTreeKVDB

import "testing"
import "fmt"

func TestWriting(t *testing.T){
	db := &DB{}
	db.Init("./t.db")

	wt := &Transaction{}
	wt.Type = 2

	wt.db = db
	fmt.Println("L0")
	Key := ([]byte)("a")
	Value := ([]byte)("b")

	wt.Begin()
	fmt.Println("L1")
	wt.Write(Key,Value)
	//fmt.Println("L2")
	v := wt.ReadForWrting(Key)
	fmt.Println((string)(v))
	wt.Commit()
}
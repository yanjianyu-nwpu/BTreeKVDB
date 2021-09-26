package BTreeKVDB

import "testing"
import "os"
func TestDBInit(t *testing.T){
	db := &DB{}
	db.Init("./t.db")
}
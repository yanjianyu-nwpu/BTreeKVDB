package BTreeKVDB

import "testing"

func TestByteLess(t *testing.T){
	
	s1 := []byte("a")

	s2 := []byte("b")

	if ans:=ByteLess(s1,s2);ans!=true{
		t.Errorf("bytes[]10 , bytes[]20 Byteless is true")
	}

}
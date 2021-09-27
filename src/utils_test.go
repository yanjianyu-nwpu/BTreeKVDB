package BTreeKVDB

import "testing"
import "fmt"
func TestByteLess(t *testing.T){
	
	s1 := []byte("a")

	s2 := []byte("b")

	if ans:=ByteLess(s1,s2);ans!=true{
		t.Errorf("bytes[]10 , bytes[]20 Byteless is true")
	}

}
func TestLowerBound(t *testing.T){
	var keys [][]byte
	keys = append(keys,[]byte("b"))
	keys = append(keys,[]byte("c"))

	if ans:=LowerBound([]byte("a"),keys);ans!=1{
		fmt.Println(ans)
		t.Errorf("ecept 1,but got %d",ans)
	}
}
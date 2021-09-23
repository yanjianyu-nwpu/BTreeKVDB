package BTreeKVDB

import "testing"
import "fmt"
func TestInsert(t *testing.T){
	f := &Freelist{}
	f.Init()

	f.Insert(0,5)
	f.Insert(9,13)

	if st,ed:=f.Get(6);st!=0 || ed != 5 || f.Size!=1{
		t.Errorf("execpt 0,6 got %d %d",st,ed)
	}
	if st,ed:=f.Get(3);st!=9 || ed !=11  || f.Size!=1{
		t.Errorf("execpt 0,6 got %d %d",st,ed)
	}
}
func TestSerealize(t *testing.T){
	f := &Freelist{}
	f.Init()

	f.Insert(0,5)
	f.Insert(9,13)

	tmp := f.Serealize()
	nf := &Freelist{}
	nf.Init()
	nf.Deserealize(tmp)
	if nf.Size!=2{
		t.Errorf("loss fnode ")
	}
	if st,ed:=nf.Get(6);st!=0 || ed != 5 || nf.Size!=1{
		fmt.Print("PPPPPP")
		nf.show()
		t.Errorf("execpt 0,6 got %d %d",st,ed)
	}
	if st,ed:=nf.Get(3);st!=9 || ed !=11  || nf.Size!=1{
		t.Errorf("execpt 0,6 got %d %d",st,ed)
	}

}
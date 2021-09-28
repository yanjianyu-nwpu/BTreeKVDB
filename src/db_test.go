package BTreeKVDB
//import "fmt"
import "testing"
func TestDBInit(t *testing.T){
	db := &DB{}
	db.Init("./t.db")


	key := ([]byte)("yanjianyu")

	p1 := &Page{}
	p1.Version = 0
	p1.Id = 4
	p1.Put(key,([]byte)("handsome"))

	db.DirtyData[4] = []*Page{p1}

	p2 := &Page{}
	p2.Version = 1
	p2.Id = 4
	p2.Put(key,([]byte)("hcx"))
	db.Cache.Put(p2)
	
	if ans:=db.GetLastestPage(4);!ByteCompare(ans.Get(key),([]byte)("hcx")){
		t.Errorf("excpt got p2")
	}
	if ans:=db.GetPageWithVersion(4,0);ans!=p1{
		t.Errorf("except got p1")
	}
	if ans:=db.GetPageWithVersion(5,0);ans==nil{
		t.Errorf("except got nil")
	}
}
package BTreeKVDB
import "testing"
import "fmt"
func TestSearlize(t *testing.T){
	p := Page{}
	p.Id = 1
	p.Type = 1
	
	k:=([]byte)("yjy")
	v:=([]byte)("handsome")
	p.Put(k,v)

	k = ([]byte)("hcx")
	v = ([]byte)("ugly")
	p.Put(k,v)

	var tmp []byte
	tmp = p.Searlize()

	var newPage Page
	newPage.Desearlize(tmp)
	

	if newPage.Id != 1{
		t.Errorf("Desearlize page Id is wrong, is 1")
	}
	if newPage.Type != 1{
		t.Errorf("Desearlize page type is wrong, is 1")
	}
	if newPage.KVSize!=2{
		t.Errorf("Desearlize KVSize is wrong, is 2")
	}
	if (string)(newPage.kvs[0].Key) != "hcx"{
		t.Errorf("Desearlize page kvs[0].Key is wrong")
	}
	if (string)(newPage.kvs[0].Value) != "ugly"{
		t.Errorf("Desearlize page kvs[0].Value is wrong")
	}
	if (string)(newPage.kvs[1].Key) != "yjy"{
		t.Errorf("Desearlize page kvs[1].Key is wrong")
	}
	if (string)(newPage.kvs[1].Value)!="handsome"{
		t.Errorf("Desearlize page kvs[1s].value is wrong")
	}
}

func TestPutInsert(t *testing.T){
	p := &Page{}
	key := ([]byte)("yanjianyu")
	p.Put(key,([]byte)("handsome"))

	p.Put(key,([]byte)("hcx"))
	if ans:=p.Get(key); !ByteCompare(ans,([]byte)("hcx")){
		t.Errorf("except get hcx but got %s",(string)(ans))
	}

}
func TestPutInsertB(t *testing.T){
	p := &Page{}

	dic := make(map[string]string)
	tmpL := 0
	bs := "a"
	var r []*Page
	for i:=0;i<=61;i++{
		dic[bs] = bs
		
		tmpL += 2*(i+1)
		bs += "a"
		r = p.Put(([]byte)(bs),([]byte)(bs))
		if len(r)>=2{
			break
		}
	}
	fmt.Println(len(r))
	fmt.Println(r[0].CurrentLength,r[0].KVSize,(string)(r[0].kvs[r[0].KVSize-1].Key))
	for i:=0;i<int(r[0].KVSize);i++{
		fmt.Print(len(r[0].kvs[i].Key)," ")
	}
	fmt.Println(" ")
	for i:=0;i<int(r[1].KVSize);i++{
		fmt.Print(len(r[1].kvs[i].Key)," ")
	}
	fmt.Println(r[1].CurrentLength,r[1].KVSize,(string)(r[1].kvs[r[1].KVSize-1].Key))
	fmt.Println(tmpL)
}

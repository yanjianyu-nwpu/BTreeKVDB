package BTreeKVDB
import "testing"
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


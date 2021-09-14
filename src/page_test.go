package BTreeKVDB
import "testing"

func TestKV_Header(t *testing.T){
	
	kv_h := KV_Header{}
	kv_h.Offset_ = 0
	kv_h.KeySize = 10
	kv_h.ValueSize= 10
	if kv_h.KeySize != 10{
		t.Errorf("Kv_h Kesize got %d" ,kv_h.KeySize)
	}
}


package BTreeKVDB

func ByteLess(k1 []byte,k2 []byte) bool{
	l1 := len(k1)
	l2 := len(k2)

	minLen := l1
	if minLen > l2{
		minLen = l2
	}
	
	for i:=0;i<minLen;i++{
		if k2[i] > k1[i]{
			return true
		}
		if k1[i] > k2[i] {
			return false
		}
	}
	if l1 > l2{
		return false
	}else{
		return true
	}
}
func ByteBigger(k1 []byte,k2 []byte) bool{
	l1 := len(k1)
	l2 := len(k2)

	minLen := l1
	if minLen > l2{
		minLen = l2
	}
	
	for i:=0;i<minLen;i++{
		if k2[i] > k1[i]{
			return false
		}
		if k1[i] > k2[i] {
			return true
		}
	}
	if l1 > l2{
		return true
	}else{
		return false
	}
}

func ByteCompare(k1 []byte, k2 []byte)bool{
	l1 := len(k1)
	l2 := len(k2)

	if l1 != l2{
		return false
	}
	for i:=0;i<l1;i++{
		if k1[i] != k2[i]{
			return false
		}
	}
	return true
}
func LowerBound(key []byte, keys [][]byte) int {
	st := 0
	ed := len(keys)-1
	if (ed == -1){
		return 0
	}
	if (ByteBigger(key,keys[ed])){
		return ed+1
	}
	for st<ed{
		mid := (st+ed)/2

		if ByteCompare(keys[mid],key){
			return mid
		}
		if (ByteLess(keys[mid],key)){
			st = mid+1 
		}else{
			ed = mid-1
		}
	}
	return ed
}
func FindKV(key []byte, kvs []KV) int{
	l := len(kvs)
	if l<=0{
		return -1
	}
	ind := LowerBoundKV(key,kvs)
	if (ind >= l){
		return -1
	}
	if ByteCompare(key,kvs[ind].Key){
		return ind
	}
	return -1
}
func LowerBoundKV(key []byte, kvs []KV) int {
	st := 0
	ed := len(kvs)-1
	if (ed == -1){
		return 0
	}
	if (ByteBigger(key,kvs[ed].Key)){
		return ed+1
	}
	for st<ed{
		mid := (st+ed)/2

		if ByteCompare(kvs[mid].Key,key){
			return mid
		}
		if (ByteLess(kvs[mid].Key,key)){
			st = mid+1 
		}else{
			ed = mid-1
		}
	}
	return ed
}
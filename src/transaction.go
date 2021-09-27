package BTreeKVDB

import "sync/atomic"
import "encoding/binary"
const(
	ReadingTransaction int = 1
	WritingTransaction int = 2
)
type Transaction struct{
	
	Type int //01 for reading transaction, 02 for writing transaction
	db *DB
	
	DBVersion uint64

	Readpages map[PgId]*Page

	WritingPages map[PgId]*Page

	AllocatePages []PgId  // for writing transaction rollback
}

func (t *Transaction)Begin(){

	if t.Type == 1{
		t.DBVersion = atomic.LoadUint64(&(t.db.CurrentVersion))
		t.db.KeepVersionMu.Lock()
		t.db.KeepVersion = append(t.db.KeepVersion,t.DBVersion)
		t.db.KeepVersionMu.Unlock()
		return 
	}
	t.db.WritingMu.Lock()	
	t.DBVersion = atomic.AddUint64(&(t.db.CurrentVersion),1)
}

func (t *Transaction)Write(Key []byte,Value []byte){
	tmp,ok := t.WritingPages[9]
	if !ok{
		tmp = t.db.GetLastestPage(9)
	}
	Path := []*Page{}
	PathInd := []int{}
	if (tmp.KVSize == 0){
		newId,_ := t.db.Freelist.Get(1)
		b:=make([]byte,8)
		binary.BigEndian.PutUint64(b,(uint64)(newId))
		tmp.Put(Key,b)
		
		Path = append(Path,tmp)
		tmp = &Page{}
		tmp.Id = (PgId)(newId)
		tmp.Version = t.DBVersion
		tmp.Type = 2
		tmp.Put(Key,Value)
		tmp.SetPosition()

		Path = append(Path,tmp)
	}else{
		for tmp.Type!=3{
			Path = append(Path,tmp)
			ind := (uint32)(LowerBoundKV(Key,tmp.kvs))
			if (ind >= tmp.KVSize){
				ind = tmp.KVSize - 1
			}
			if ByteLess(Key,tmp.kvs[(int)(ind)].Key){
				ind -= 1
			}
			if (ind <= 0){
				ind = 0
			}
			PathInd = append(PathInd,(int)(ind))
			real_index := binary.BigEndian.Uint64(tmp.kvs[ind].Value)
			tmp = t.db.GetLastestPage((PgId)(real_index))
		}
		newP := tmp.Put(Key,Value)
		
		for i:=0;i<len(newP);i++{
			if newP[i].Id == 0{
				ns,_ := t.db.Freelist.Get(1)
				if ns == -1{
					t.db.CapacityExpansion()
					ns,_ = t.db.Freelist.Get(1)
				}
				newP[i].Type = 3
				newP[i].Id = (PgId)(ns)
				newP[i].Version = t.DBVersion
			}
		}


	}
}
/*
func (t *Transaction)Rollback(){
	//only for writing transaction 

}	
func (t *Transaction)Commit(){
	if t.Type == 1{
		Over := false
		t.db.KeepVersionMu.Lock()
		for i := 0;i<len(t.db.KeepVersion);i++{
			if t.db.KeepVersion[i] == t.DBVersion{
				t.db.KeepVersion = append(t.db.KeepVersion[:i],t.db.KeepVersion[i+1:]...)
				if (len(t.db.KeepVersion) == 0) || (len(t.db.KeepVersion)>0 && t.db.KeepVersion[0] > t.DBVersion){
					Over = true
				}
				break
			}
		}
		t.db.KeepVersionMu.Unlock()
		if Over{
			t.db.DirtyCacheMu.Lock()
			
			t.db.DirtyCacheMu.Unlock()
		}
	}
}*/
package BTreeKVDB
/*
import "sync/atomic"
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
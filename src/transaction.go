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
func (t *Transaction)GetPageForWrting(p PgId){
	//找到一个页面，逻辑是这样，先从writingPage里面找
	//如果没有，用db.GetLastedPage找最新版本的副本
	//不要将db中最新版本更新，因为这个再commit时候
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
		//这里的逻辑是这样
		//先找到叶子节点，招叶子节点的时候要注意小于最小的时候，和大于最大的时候，已经lowrboundkv的值，记录路径
		//插入叶子节点看看是不是有溢出
		//递归到根节点，如果一层有溢出，要插入新值，并且要申请新的PgId,都要加入writingPage这个map
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
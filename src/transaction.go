package BTreeKVDB

import (
	"encoding/binary"
	"sync/atomic"
	"fmt"
)

const (
	ReadingTransaction int = 1
	WritingTransaction int = 2
)

type Transaction struct {
	Type int //01 for reading transaction, 02 for writing transaction
	db   *DB

	DBVersion uint64

	Readpages map[PgId]*Page

	WritingPages map[PgId]*Page

	AllocatePages []PgId // for writing transaction rollback
}

func (t *Transaction) GetPageForWrting(p PgId) *Page {
	//找到一个页面，逻辑是这样，先从writingPage里面找
	//如果没有，用db.GetLastedPage找最新版本的副本
	//不要将db中最新版本更新，因为这个再commit时候

	r, ok := t.WritingPages[p]
	fmt.Println(r==nil,"TTTT")
	if !ok {
		r = t.db.GetLastestPage(p)
		t.WritingPages[p] = r
	}
	return r
}
func (t *Transaction) GetSingleNewPage() PgId {
	//allocate a new page else

	ns, _ := t.db.Freelist.Get(1)
	if ns == -1 {
		t.db.CapacityExpansion()
		ns, _ = t.db.Freelist.Get(1)
	}
	return (PgId)(ns)
}
func (t *Transaction) WrtingOver(path []*Page) {
	for i := 0; i < len(path); i++ {
		ns := path[i].Id
		t.WritingPages[ns] = path[i]
	}
}
func (t *Transaction) Begin() {

	if t.Type == 1 {
		t.DBVersion = atomic.LoadUint64(&(t.db.CurrentVersion))
		t.db.KeepVersionMu.Lock()
		t.db.KeepVersion = append(t.db.KeepVersion, t.DBVersion)
		t.db.KeepVersionMu.Unlock()
		return
	}
	t.db.WritingMu.Lock()
	t.DBVersion = atomic.AddUint64(&(t.db.CurrentVersion), 1)
}
func (t *Transaction) ReadForWrting(Key []byte) []byte {
	tmp := t.GetPageForWrting(9)
	for tmp.Type != 3 {
		ind := LowerBoundKV(Key, tmp.kvs)
		if ind >= int(tmp.KVSize) {
			ind = int(tmp.KVSize - 1)
		}
		if ByteLess(Key, tmp.kvs[ind].Key) {
			ind -= 1
		}
		if ind <= 0 {
			ind = 0
		}
		read_ind := (PgId)(binary.BigEndian.Uint64(tmp.kvs[ind].Value))
		tmp = t.GetPageForWrting(read_ind)
	}
	res := tmp.Get(Key)
	return res
}
func (t *Transaction) Write(Key []byte, Value []byte) {
	//这里的逻辑是这样
	//先找到叶子节点，招叶子节点的时候要注意小于最小的时候，和大于最大的时候，已经lowrboundkv的值，记录路径
	//插入叶子节点看看是不是有溢出
	//递归到根节点，如果一层有溢出，要插入新值，并且要申请新的PgId,都要加入writingPage这个map
	tmp := t.GetPageForWrting(9)
	fmt.Println("L2")
	//新打开的数据库
	Path := []*Page{tmp}
	PathInd := []PgId{}
	//if is a new database
	if tmp.KVSize == 0 {
		fmt.Println("L1")
		ns := t.GetSingleNewPage()
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, (uint64)(ns))
		tmp.Put(Key, b)
		tmp = &Page{
			Id:   ns,
			Type: 3,
		}
		tmp.Put(Key, Value)
		tmp.SetPosition()
		Path = append(Path, tmp)
		t.WrtingOver(Path)
		return
	}
	//find the leaf node and store the path
	for tmp.Type != 3 {
		ind := LowerBoundKV(Key, tmp.kvs)
		if ind >= int(tmp.KVSize) {
			ind = int(tmp.KVSize - 1)
		}
		if ByteLess(Key, tmp.kvs[ind].Key) {
			ind -= 1
		}
		if ind <= 0 {
			ind = 0
		}
		read_ind := (PgId)(binary.BigEndian.Uint64(tmp.kvs[ind].Value))
		PathInd = append(PathInd, read_ind)
		Path = append(Path, tmp)
		tmp = t.GetPageForWrting(read_ind)
	}
	Path = append(Path, tmp)
	final := (int)(len(PathInd)) - 1
	tmp = Path[final]
	newPage := tmp.Put(Key, Value)
	for tmp.Type != 3 && final >= 0 {
		tmp = Path[final]
		//the max size of newPage is 3
		//allocate new PgId for new page
		for i := 1; i < len(newPage); i++ {
			ns := t.GetSingleNewPage()
			newPage[i].Id = ns
			Path = append(Path, newPage[i])
		}
		tmp.kvs[PathInd[final]].Key = newPage[0].kvs[0].Key
		//iT := &[]*Page{tmp}
		var a1 *Page
		var a2 *Page
		a1 = nil
		a2 = nil
		for i := 1; i < len(newPage); i++ {
			if a1 == nil || (a1 != nil && ByteLess(Key, a1.kvs[0].Key)) {
				b := make([]byte, 8)
				binary.BigEndian.PutUint64(b, (uint64)(newPage[i].Id))
				nnn := tmp.Put(newPage[i].kvs[0].Key, b)
				if len(nnn) > 1 {
					a1 = nnn[1]
				}
				if len(nnn) > 2 {
					a2 = nnn[2]
				}
				continue
			}
			if a1 != nil && ((a2 == nil && !ByteLess(Key, a1.kvs[0].Key)) || (a2 != nil && ByteLess(Key, a2.kvs[0].Key))) {
				b := make([]byte, 8)
				binary.BigEndian.PutUint64(b, (uint64)(newPage[i].Id))
				nnn := tmp.Put(newPage[i].kvs[0].Key, b)
				if len(nnn) > 1 {
					a2 = nnn[1]
				}
			} else {
				b := make([]byte, 8)
				binary.BigEndian.PutUint64(b, (uint64)(newPage[i].Id))
				tmp.Put(newPage[i].kvs[0].Key, b)
			}
		}
		newPage = []*Page{tmp}
		if a1 != nil {
			newPage = append(newPage, a1)
		}
		if a2 != nil {
			newPage = append(newPage, a2)
		}
		final -= 1
	}
	if len(newPage) > 1 {
		ns := t.GetSingleNewPage()
		t.WritingPages[9].Id = ns
		t.WritingPages[9].Type = 2

		new_root := &Page{
			Id:   9,
			Type: 3,
		}
		for i := 0; i <= len(newPage); i++ {
			b := make([]byte, 8)
			binary.BigEndian.PutUint64(b, (uint64)(newPage[i].Id))
			new_root.Put(newPage[i].kvs[0].Key, b)
		}
		t.WritingPages[9] = new_root
		Path = append(Path, new_root)
	}
	for i := 0; i < len(Path); i++ {
		Path[i].SetPosition()
	}
	t.WrtingOver(Path)
}

/*
func (t *Transaction)Rollback(){
	//only for writing transaction

}
*/
func (t *Transaction)Commit(){
	if t.Type ==2 {
		t.db.LruCacheMu.Lock()
		t.db.DirtyCacheMu.Lock()
		for v,k := range t.WritingPages{
			s:=t.db.GetLastestPage(v)
			_,ok := t.db.DirtyData[v]
			if !ok{
				t.db.DirtyData[v] = []*Page{s}
			}else{
				t.db.DirtyData[v] = append(t.db.DirtyData[v],s)
			}
			t.db.Cache.Put(k)
		}
		
		t.db.DirtyCacheMu.Unlock()
		t.db.LruCacheMu.Unlock()
	}
	t.db.WritingMu.Unlock()
}

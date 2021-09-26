package BTreeKVDB
import "os"
import "sync"
import "sync/atomic"

type DB struct{
	Path string
	F *os.File
	DBSize int64

	LruCacheMu sync.Mutex 
	Cache *LruPageCache
	
	CurrentVersion uint64 //a writing transaction will add 1 

	KeepVersionMu sync.Mutex
	KeepVersion []uint64 //the version of current reading transaction event

	DirtyCacheMu sync.Mutex

	DirtyData map[PgId][]*Page  //the dirty page for reading transactio event
	
	WritingMu sync.Mutex //mutex for writing transaction event , only one writing transaction at one time

	Freelist *Freelist	
}

func (db *DB)Init(Path string){
	db.Path = Path
	db.F,_ = os.OpenFile(db.Path,os.O_CREATE | os.O_RDWR,0755)
	info,_ := db.F.Stat()
	db.DBSize = info.Size()

	db.CurrentVersion = 0

	db.Freelist = new(Freelist)
	db.Freelist.Init()

	db.Cache = &LruPageCache{}
	db.Cache.Init(20000)

	db.DirtyData = make(map[PgId][]*Page)
	if db.DBSize < 100{
		db.F.Truncate(4096*1024)
		db.DBSize = 1024

		db.Freelist.Insert(0,1024)

		rootPage := &Page{}
		rootPage.Id = 8
		rootPage.Type = 3
		rootPage.Version = db.CurrentVersion
		
		rootBuffer := rootPage.Searlize()
		db.Cache.Put(rootPage)
		db.F.WriteAt(rootBuffer,4096*8)
	}else{
		db.DBSize = db.DBSize/4096

		FreelistBuffer := make([]byte,4096*8)
		db.Freelist.Deserealize(FreelistBuffer)
		
		
		rootBuffer := make([]byte,4096)
		db.F.ReadAt(rootBuffer,4096*8)
		rootPage := &Page{}
		rootPage.Desearlize(rootBuffer)
		rootPage.Type = 3
		rootPage.Version = db.CurrentVersion
		db.Cache.Put(rootPage)
	}
}


//DB 扩容
func (db *DB)CapacityExpansion(){
	newDBSize := db.DBSize * 2

	db.F.Truncate(newDBSize)
	db.Freelist.Insert((int)(db.DBSize),(int)(db.DBSize*2-1))
	
	freelistBuffer := db.Freelist.Serealize()
	db.F.WriteAt(freelistBuffer,0)
}
func (db *DB) GetLastestPage(Id PgId) *Page{
	db.LruCacheMu.Lock()
	defer db.LruCacheMu.Unlock()
	tmp := db.Cache.Get(Id)

	if tmp != nil{
		return tmp
	}
	
	tmp = &Page{}
	readbuffer := make([]byte,4096)
	db.F.ReadAt(readbuffer,(int64)(Id*4096))
	tmp.Desearlize(readbuffer)
	tmp.Id = Id
	tmp.Version = atomic.LoadUint64(&(db.CurrentVersion))
	db.Cache.Put(tmp)
	return tmp
}
func (db *DB) GetPageWithVersion(Id PgId, Version uint64) *Page{
	var PagePtr *Page
	PagePtr = nil
	db.DirtyCacheMu.Lock()
	//defer db.DirtyCacheMu.Unlock()
	dirty,ok := db.DirtyData[Id]
	if (ok){
		for i:=0;i<len(dirty);i++{
			//if get the version return directly
			if dirty[i].Version == Version {
				db.DirtyCacheMu.Unlock()
				return dirty[i]
			}
			//if the dirty[i].version is more closer to version
			if dirty[i].Version < Version && (PagePtr == nil || PagePtr.Version < dirty[i].Version){
				PagePtr = dirty[i]
			}
		}
	}
	db.DirtyCacheMu.Unlock()
	db.LruCacheMu.Lock()
	defer db.LruCacheMu.Unlock()
	tmp := db.Cache.Get(Id)

	if tmp != nil && tmp.Version <= Version{
		return tmp
	}
	
	tmp = &Page{}
	readbuffer := make([]byte,4096)
	db.F.ReadAt(readbuffer,(int64)(Id*4096))
	tmp.Desearlize(readbuffer)
	tmp.Id = Id
	tmp.Version = Version
	db.Cache.Put(tmp)
	return tmp
}
package BTreeKVDB

type PgId uint32 
type TxId uint32
type KVs []KV
type KV_Headers []KV_Header

type KV struct{
	Key				[]byte
	Value			[]byte
	Next			*KV
}
type KV_Header struct{
	Offset_			int //the position  in this page
	KeySize			int //size of key
	ValueSize		int //size of value
	Next			*KV_Headers
}
type Page struct{
	Id 				PgId
	Current_Length 	int    	//the size of current page	
	Current_Txid 	TxId    //the latest id of writted Txid 

	KV_size 		int

	First_kv 		*KV
	First_kv_header *KV_Header		
}

func (p *Page)Put(key []byte,value []byte) bool {

	return true	
}
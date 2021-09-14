package BTreeKVDB


type PgId uint32 


type KV struct{
	Key				[]byte
	Value			[]byte
}
type KV_Header struct{
	Offset_			int //the position  in this page
	KeySize			int //size of key
	ValueSize		int //size of value
}
type Page struct{
	Id 				PgId
	Current_Length 	int    	//the size of current page	
}
package BTreeKVDB

import "unsafe"
import "encoding/binary"
//import "fmt"
type PgId uint64
type TxId uint64
type KVs []KV
type KVHeaders []KVHeader

const PageHeaderSize = (uint32)(unsafe.Sizeof(PgId(0)) + unsafe.Sizeof(int(0))*3 + unsafe.Sizeof(TxId(0)))
const KVHeaderSize = (uint32)(unsafe.Sizeof(KVHeader{}))
const HeaderSize = (uint32)(20)

type KV struct{
	Key				[]byte
	Value			[]byte
}
type KVHeader struct{
	Offset_			uint32 //the position  in this page
	KeySize			uint32 //size of key
	ValueSize		uint32 //size of value
}
type Page struct{
	Id 				PgId
	CurrentLength 	uint32   	//the size of current page	
	CurrentTxid 	TxId    //the latest id of writted Txid 
	Type 		uint32 //01 branch page, 02 leafpage, 03 rootpage
	KVSize 		uint32
	Version uint64
	kvHeaders			KVHeaders	
	kvs			KVs	

}
func (p *Page)ClearUp(){
	p.CurrentLength = 20
	p.Type = 0
	p.KVSize = 0
	p.kvHeaders = KVHeaders{}
	p.kvs= KVs{}
}
func (p *Page)SetPosition(){

	p.KVSize = (uint32)(len(p.kvs))
	tmpLength := PageHeaderSize + uint32(p.KVSize)*uint32(KVHeaderSize)
	
	for i:=0;i<int(p.KVSize);i++{
		p.kvHeaders[i].Offset_ = tmpLength
		p.kvHeaders[i].KeySize = (uint32)(len(p.kvs[i].Key))
		p.kvHeaders[i].ValueSize = (uint32)(len(p.kvs[i].Value))
		tmpLength += (uint32)(len(p.kvs[i].Key) + len(p.kvs[i].Value))
	}
	p.CurrentLength = tmpLength
}
func (p *Page)GetKVs() KVs{
	return p.kvs
}
//插入KV,return []*page
func (p *Page)Put(key []byte,value []byte) []*Page {
	
	//find the index of new kv
	l := len(p.kvs)

	index := LowerBoundKV(key,p.kvs)
	insertKV := KV{Key : key,Value : value}
	tmp := make([]KV,l+1)
	if l == 0{
		p.kvs = []KV{insertKV}
	}else{
		for i:=0;i<index;i++{
			tmp[i] = p.kvs[i]
 		}
		tmp[index] = insertKV
		for i:=index;i<(int)(p.KVSize);i++{
			tmp[i+1] = p.kvs[i]
		}
		p.kvs = tmp
	}
	insertKVHeader := KVHeader{}
	insertKVHeader.KeySize = (uint32)(len(key))
	insertKVHeader.ValueSize = (uint32)(len(value))

	tmpHeader := append(p.kvHeaders[:index],insertKVHeader)
	tmpHeader = append(tmpHeader,p.kvHeaders[index:]...)
	p.kvHeaders = tmpHeader	
	
	p.KVSize = uint32(len(p.kvs))
	var tmpLength uint32
	tmpLength = PageHeaderSize + uint32(p.KVSize)*uint32(KVHeaderSize)
	
	for i:=0;i<int(p.KVSize);i++{
		p.kvHeaders[i].Offset_ = tmpLength
		p.kvHeaders[i].KeySize = (uint32)(len(p.kvs[i].Key))
		p.kvHeaders[i].ValueSize = (uint32)(len(p.kvs[i].Value))
		tmpLength += (uint32)(len(p.kvs[i].Key) + len(p.kvs[i].Value))
	}
	p.CurrentLength = tmpLength
	
	
	res := []*Page{p}
	if p.CurrentLength > 4096{

		okvs := p.kvs
		okvHeaders := p.kvHeaders
		
		p.ClearUp()

	
		tmp_p := p

		var tmpLength uint32
		tmpLength = 0
		for cur_ind:=0;cur_ind<len(okvs);cur_ind++{
			header := KVHeader{}
			header.KeySize = (uint32)(len(okvs[cur_ind].Key))
			header.ValueSize = (uint32)(len(okvs[cur_ind].Value))
			if tmpLength >= 3000 || tmpLength + header.KeySize + header.ValueSize > 4096{
				tmp_p = &Page{}
				tmp_p.ClearUp()
				res = append(res,tmp_p)
				tmpLength = PageHeaderSize
			}
			tmp_p.kvs = append(tmp_p.kvs,okvs[cur_ind])
			tmp_p.kvHeaders = append(tmp_p.kvHeaders,okvHeaders[cur_ind])
			tmpLength += (header.KeySize + header.ValueSize)
		}
		for i:=0;i<len(res);i++{
			res[i].SetPosition()
		}
	}
	return res	
}
func (p *Page)Get(key []byte) []byte{
	ind:=FindKV(key,p.kvs)
	//fmt.Println(ind)
	if(ind == -1){
		return nil 
	}
	return p.kvs[ind].Value
}
func (p *Page)Delete(key []byte)bool{
	ind:=FindKV(key,p.kvs)
	if ind == -1{
		return false
	}

	p.kvHeaders = append(p.kvHeaders[:ind],p.kvHeaders[ind+1:]...)
	p.kvs = append(p.kvs[:ind],p.kvs[ind+1:]...)
	p.KVSize = (uint32)(len(p.kvs))

	tmpLength := PageHeaderSize + p.KVSize*KVHeaderSize
	
	for i:=0;i<(int)(p.KVSize);i++{
		p.kvHeaders[i].Offset_ = tmpLength
		tmpLength += (uint32)(len(p.kvs[i].Key) + len(p.kvs[i].Value))
	}
	p.CurrentLength = tmpLength

	return true	
	
}
func (p *Page)Searlize() []byte{
	buffer := make([]byte,4096)

	tmp := 0
	b := make([]byte,8)
	binary.BigEndian.PutUint64(b,(uint64)(p.Id))

	copy(buffer[tmp:tmp+8],b)
	tmp += 8

	binary.BigEndian.PutUint32(b,p.CurrentLength)
	copy(buffer[tmp:tmp+4],b[:4])
	tmp += 4

	binary.BigEndian.PutUint32(b,p.Type)
	copy(buffer[tmp:tmp+4],b[:4])
	tmp += 4

	binary.BigEndian.PutUint32(b,p.KVSize)
	copy(buffer[tmp:tmp+4],b[:4])
	tmp += 4
	
	for i:=0;i<(int)(p.KVSize);i++{
		binary.BigEndian.PutUint32(b,p.kvHeaders[i].Offset_)
		copy(buffer[tmp:tmp+4],b[:4])
		tmp += 4

		binary.BigEndian.PutUint32(b,p.kvHeaders[i].KeySize)
		copy(buffer[tmp:tmp+4],b[:4])
		tmp += 4

		binary.BigEndian.PutUint32(b,p.kvHeaders[i].ValueSize)
		copy(buffer[tmp:tmp+4],b[:4])
		tmp += 4
	}

	for i:=0;i<(int)(p.KVSize);i++{
		ksL:=(int)(p.kvHeaders[i].KeySize)
		kvL:=(int)(p.kvHeaders[i].ValueSize)
		copy(buffer[tmp:tmp+ksL],p.kvs[i].Key)
		tmp += ksL
		copy(buffer[tmp:tmp+kvL],p.kvs[i].Value)
		tmp += kvL
	}
	return buffer
}
func (p *Page)Desearlize(data []byte) {
	
	tmp := 0
	p.Id = (PgId)(binary.BigEndian.Uint64(data[tmp:tmp+8]))
	tmp +=8

	p.CurrentLength = binary.BigEndian.Uint32(data[tmp:tmp+4])
	tmp+=4

	p.Type = binary.BigEndian.Uint32(data[tmp:tmp+4])
	tmp+=4

	p.KVSize = binary.BigEndian.Uint32(data[tmp:tmp+4])
	tmp+=4

	p.kvHeaders = make([]KVHeader,(int)(p.KVSize))
	for i:=0;i<(int)(p.KVSize);i++{
		Header := KVHeader{}
		Header.Offset_ = (uint32)(binary.BigEndian.Uint32(data[tmp:tmp+4]))
		tmp += 4

		Header.KeySize = (uint32)(binary.BigEndian.Uint32(data[tmp:tmp+4]))
		tmp += 4

		Header.ValueSize = (uint32)(binary.BigEndian.Uint32(data[tmp:tmp+4]))
		tmp += 4
		p.kvHeaders[i] = Header
		//p.kvHeaders = append(p.kvHeaders,Header)
	}

	p.kvs = make([]KV,(int)(p.KVSize))
	for i:=0;i<(int)(p.KVSize);i++{
		kL := (int)(p.kvHeaders[i].KeySize)
		vL := (int)(p.kvHeaders[i].ValueSize)
		p.kvs[i].Key = data[tmp:tmp+kL]
		tmp += kL
		p.kvs[i].Value = data[tmp:tmp+vL]
		tmp += vL
	}

}
package BTreeKVDB

//LRU cache for reading and writing,

type LruPageCache struct{
	Capactiy 	int
	Size 		int
	ListHead 	*PageNode
	ListTail 	*PageNode
	Map 		map[PgId]*PageNode
}  

type PageNode struct{
	Id 			PgId
	PagePtr 	*Page
	Next 		*PageNode
	Pre 		*PageNode
}
func (c *LruPageCache)Init(cap int){
	c.Capactiy = cap
	c.Size = 0
	c.ListHead = new(PageNode)
	c.ListTail = new(PageNode)

	c.ListHead.Next = c.ListTail
	c.ListHead.Pre = nil
	c.ListTail.Next = nil
	c.ListTail.Pre = c.ListHead

	c.Map = make(map[PgId]*PageNode)
}

func (c *LruPageCache)Get(Id PgId)*Page{
	pageptr,ok :=  c.Map[Id]

	//key not exist
	if !ok{
		return nil
	}
	//newPage := new(PageNode)
	pageptr.Pre.Next = pageptr.Next
	pageptr.Next.Pre = pageptr.Pre

	//hot page push front of list
	c.ListHead.Next.Pre = pageptr

	pageptr.Next = c.ListHead.Next
	pageptr.Pre = c.ListHead
	c.ListHead.Next = pageptr

	return pageptr.PagePtr
}

func (c *LruPageCache)Put(p *Page){

	Id := p.Id
	ptr,ok :=  c.Map[Id]

	//key not exist
	if !ok{
		ptr = nil
	}

	if ptr != nil{
		ptr.Next.Pre = ptr.Pre
		ptr.Pre.Next = ptr.Next

		c.ListHead.Next.Pre = ptr
		ptr.Next = c.ListHead.Next
		ptr.Pre = c.ListHead
		c.ListHead.Next = ptr
		return
	}
	ptr = &PageNode{}
	ptr.PagePtr = p
	ptr.Next = c.ListHead.Next
	ptr.Pre = c.ListHead
	c.ListHead.Next = ptr
	c.Map[Id] = ptr
	if (c.Size < c.Capactiy){
		c.Size += 1
		return 
	}else{
		reId := c.ListTail.Pre.Id
		delete(c.Map,reId)
		c.ListTail.Pre.Pre.Next = c.ListTail
		c.ListTail.Pre = c.ListTail.Pre.Pre
		return
	}
	
}
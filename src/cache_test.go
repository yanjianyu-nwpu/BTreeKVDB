package BTreeKVDB

import "testing"
func TestLruCachePut(t *testing.T){
	c := LruPageCache{}
	c.Init(3)

	p1 := Page{}
	p1.Id = 1

	p2 := Page{}
	p2.Id = 2

	p3 := Page{}
	p3.Id = 3

	p4 := Page{}
	p4.Id = 4

	c.Put(&p1)
	c.Put(&p2)
	c.Put(&p3)
	c.Put(&p4)
	ans:=c.Get(1)
	
	if ans!=nil{
		t.Errorf("the capaciy of Lru is wrong")
	}
	if p:=c.Get(2);p.Id!=2{
		t.Errorf("the lru Get is wrong")
	}
}
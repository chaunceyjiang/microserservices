package data

import "testing"

func TestValidateProduct(t *testing.T)  {

	p:=&Product{Name: "s",Price: 1.2,SKU: "aa-aa-bb"}
	if err:=p.Validate();err!=nil{
		t.Fatal(err)
	}

}
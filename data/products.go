package data

import (
	"encoding/json"
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"io"
	"regexp"
	"time"
)

// 定义一个产品的结构体
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreateOn    string  `json:"-"`
	UpdateOn    string  `json:"-"`
	DeleteOn    string  `json:"-"`
}

// FromJSON 从数据流中解码
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Product) Validate() error {
	validate:=validator.New()
	validate.RegisterValidation("sku",validateSKU)
	return validate.Struct(p)
}
// 自定一个验证器
func validateSKU(fl validator.FieldLevel) bool {
	// 比如 SKU abc-def-erf
	reg:=regexp.MustCompile("[a-z]+-[a-z]+-[a-z]+")

	matches:=reg.FindAllString(fl.Field().String(),-1)
	if len(matches)!=1{
		return false
	}
	return true
}


var productList = Products{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "一种苦coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreateOn:    time.Now().UTC().String(),
		UpdateOn:    time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "milk",
		Description: "一种牛奶",
		Price:       1.99,
		SKU:         "fjd1234",
		CreateOn:    time.Now().UTC().String(),
		UpdateOn:    time.Now().UTC().String(),
	},
}

// 产品集合
type Products []*Product

// ToJSON 直接使用json 编码器输出到w
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) error {
	productList = append(productList, p)
	return nil
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[pos] = p
	return nil
}

var ErrProductNotFind = errors.New("product not found")

func findProduct(id int) (*Product, int, error) {
	for pos, p := range productList {
		if p.ID == id {
			return p, pos, nil
		}
	}
	return nil, 0, ErrProductNotFind
}

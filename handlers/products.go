package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"microserservices/data"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

// swagger:response Resp
// Resp 返回统一封装
type Resp struct {
	// Required: true
	// Code 状态码
	Code int `json:"code"`
	// Required: true
	// Message 状态码
	Message string `json:"message"`
	// Required: true
	// Data 状态码
	Data interface{} `json:"data"`
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l: l}
}

//
//func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		p.getProducts(w, r)
//		return
//	}
//	if r.Method == http.MethodPost {
//		p.addProducts(w, r)
//		return
//	}
//	if r.Method == http.MethodPut {
//		// except id in URL
//		reg := regexp.MustCompile(`/([0-9]+)`)
//		// 使用方法
//		// https://pkg.go.dev/regexp#Regexp.FindAllStringSubmatch
//		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
//		if len(g) != 1 {
//			http.Error(w, "Invalid URL", http.StatusBadRequest)
//			return
//		}
//		if len(g[0]) != 2 {
//			http.Error(w, "Invalid URL", http.StatusBadRequest)
//			return
//		}
//		idString := g[0][1]
//		id, err := strconv.Atoi(idString)
//		if err != nil {
//			p.l.Println(" id converted error")
//			http.Error(w, "Invalid URL", http.StatusBadRequest)
//			return
//		}
//		p.updateProduct(id, w, r)
//		return
//	}
//	//catch all
//	w.WriteHeader(http.StatusMethodNotAllowed)
//}

func (p *Products) AddProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Print("Handle POST Products")
	//prod := &data.Product{}
	//// r.Body 实现io.Reader
	//err := prod.FromJSON(r.Body)
	//if err != nil {
	//	http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	//	return
	//}
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}
// swagger:response ResponseProductsWrapper
// ResponseProducts
type ResponseProductsWrapper struct {
	// in: body
	Body []*Products
}

// swagger:route GET /products products listProducts
//  返回商品详情
// Responses:
// 200: ResponseProductsWrapper

// GetProducts 返回 商品详情
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	d := data.GetProducts()
	if err := d.ToJSON(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
// ErrIdConvertFail id 转化失败
var ErrIdConvertFail = errors.New("id convert error")

// swagger:route PUT /products/{id} products updateProducts
//  返回商品详情
// Responses:
// 200: Resp

// UpdateProduct 更新商品详情
func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// 加载变量
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, ErrIdConvertFail.Error(), http.StatusBadRequest)
	}
	p.l.Print("Handle Put Products")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	if err = data.UpdateProduct(id, prod); errors.Is(err, data.ErrProductNotFind) {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}
	err = data.UpdateProduct(id, prod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

}

type KeyProduct struct {
}

// MiddlewareProductValidation 中间件
func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		// r.Body 实现io.Reader
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}
		err = prod.Validate()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error Validate %s", err.Error()), http.StatusBadRequest)
			return
		}
		// 重新构造新的请求,并且通过ctx 传递变量
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		// 调用下一个handle,或者下一个中间件，或者终止
		next.ServeHTTP(w, req)
	})
}

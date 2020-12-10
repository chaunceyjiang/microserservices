package handlers

import (
	"errors"
	"log"
	"microserservices/data"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l: l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProducts(w, r)
		return
	}
	if r.Method == http.MethodPut {
		// except id in URL
		reg := regexp.MustCompile(`/([0-9]+)`)
		// 使用方法
		// https://pkg.go.dev/regexp#Regexp.FindAllStringSubmatch
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println(" id converted error")
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		p.updateProduct(id, w, r)
		return
	}
	//catch all
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) addProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Print("Handle POST Products")
	prod := &data.Product{}
	// r.Body 实现io.Reader
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}
	data.AddProduct(prod)
}

// getProducts 返回
func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	d := data.GetProducts()
	if err := d.ToJSON(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Print("Handle Put Products")
	prod := &data.Product{}
	// r.Body 实现io.Reader
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}
	if err = data.UpdateProduct(id, prod); errors.Is(err, data.ErrProductNotFind) {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}
	err = data.UpdateProduct(id, prod)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

}

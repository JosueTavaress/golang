package handler

import (
	"net/http"

	"app/internal"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
)

// NewProductsDefault returns a new ProductsDefault
func NewProductsDefault(sv internal.ServiceProduct) *ProductsDefault {
	return &ProductsDefault{sv: sv}
}

// ProductsDefault is a struct that returns the product handlers
type ProductsDefault struct {
	// sv is the product's service
	sv internal.ServiceProduct
}

// ProductJSON is a struct that represents a product in JSON format
type ProductJSON struct {
	Id          int     `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// GetAll returns all products
func (h *ProductsDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		// process
		p, err := h.sv.FindAll()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error getting products")
			return
		}

		// response
		// - serialize
		pJSON := make([]ProductJSON, len(p))
		for ix, v := range p {
			pJSON[ix] = ProductJSON{
				Id:          v.Id,
				Description: v.Description,
				Price:       v.Price,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "products found",
			"data":    pJSON,
		})
	}
}

type ProductAmountSoldResponseDto struct {
	Description string  `json:"description"`
	Total       float64 `json:"total"`
}

func (h *ProductsDefault) GetTopProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productAmount, err := h.sv.FindTopProductsByAmount(5)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error get top products")
			return
		}

		products := make([]ProductAmountSoldResponseDto, len(productAmount))
		for index, v := range productAmount {
			products[index] = ProductAmountSoldResponseDto{
				Description: v.Description,
				Total:       v.Total,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "products found",
			"data":    products,
		})
	}
}

// RequestBodyProduct is a struct that represents the request body for a product
type RequestBodyProduct struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// Create creates a new product
func (h *ProductsDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - body
		var reqBody RequestBodyProduct
		err := request.JSON(r, &reqBody)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "error parsing request body")
			return
		}

		// process
		// - deserialize
		p := internal.Product{
			ProductAttributes: internal.ProductAttributes{
				Description: reqBody.Description,
				Price:       reqBody.Price,
			},
		}
		// - save
		err = h.sv.Save(&p)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error creating product")
			return
		}

		// response
		// - serialize
		pr := ProductJSON{
			Id:          p.Id,
			Description: p.Description,
			Price:       p.Price,
		}
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "product created",
			"data":    pr,
		})
	}
}

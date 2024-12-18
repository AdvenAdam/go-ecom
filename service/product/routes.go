package product

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/AdvenAdam/go-ecom/types"
	"github.com/AdvenAdam/go-ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)

}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, ps)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload

	// TODO - get request body for image file
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		log.Println(err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	filename, uploadDestination, err := utils.GetFileNameDestination(fileHeader)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// NOTE - This prevents the file type in r.body from being parsed as a string
	// - get request body
	payload.Name = r.FormValue("name")
	payload.Description = r.FormValue("description")
	payload.Price, _ = strconv.ParseFloat(r.FormValue("price"), 64)
	payload.Quantity, _ = strconv.Atoi(r.FormValue("quantity"))
	payload.Image = uploadDestination

	// - validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to validate payload: %v", errors))
		return
	}
	// - create product
	err = h.store.CreateProduct(&types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	uploadedFilePath, err := utils.UploadFile(file, filename, uploadDestination)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("upload failed: %v", err))
		// Delete the uploaded file if there was an error
		os.Remove(uploadedFilePath)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}

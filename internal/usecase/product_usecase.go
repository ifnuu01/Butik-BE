package usecase

import (
	"butik/internal/domain"
	"butik/internal/domain/dto"
	"butik/internal/repository"
	"errors"
	"os"
)

type ProductUsecase interface {
	CreateProduct(req domain.CreateProductRequest, imageURL string) (*domain.CreateProductResponse, error)
	GetAllProducts(offset, limit int) ([]*domain.ProductResponse, int, error)
	GetProductByID(id uint) (*domain.ProductResponse, error)
	UpdateProduct(id uint, req domain.UpdateProductRequest, imageURL string) (*domain.UpdateProductResponse, error)
	DeleteProduct(id uint) (*domain.DeleteProductResponse, error)
	ReduceStock(productID uint, qty int) error
}

type productUsecase struct {
	productRepo  repository.ProductRepo
	categoryRepo repository.CategoryRepo
}

func NewProductUsecase(productRepo repository.ProductRepo, categoryRepo repository.CategoryRepo) ProductUsecase {
	return &productUsecase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (u *productUsecase) CreateProduct(req domain.CreateProductRequest, imageURL string) (*domain.CreateProductResponse, error) {
	// Validasi category
	category, err := u.categoryRepo.GetCategoryByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	product := domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		Category:    *category,
		ImageURL:    imageURL,
	}

	createdProduct, err := u.productRepo.CreateProduct(product)
	if err != nil {
		return nil, err
	}

	// Load category untuk response
	createdProduct.Category = *category

	return &domain.CreateProductResponse{
		Message: "Product created successfully",
		Product: *dto.ToProductResponse(createdProduct),
	}, nil
}

func (u *productUsecase) GetAllProducts(offset, limit int) ([]*domain.ProductResponse, int, error) {
	products, total, err := u.productRepo.GetAllProducts(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return dto.ToProductResponses(products), total, nil
}

func (u *productUsecase) GetProductByID(id uint) (*domain.ProductResponse, error) {
	product, err := u.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return dto.ToProductResponse(product), nil
}

func (u *productUsecase) UpdateProduct(id uint, req domain.UpdateProductRequest, imageURL string) (*domain.UpdateProductResponse, error) {
	// Cek product ada
	existingProduct, err := u.productRepo.GetProductByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	category, err := u.categoryRepo.GetCategoryByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Jika ada image baru, hapus image lama
	if imageURL != "" && existingProduct.ImageURL != "" {
		deleteFile(existingProduct.ImageURL)
	}

	// Jika tidak ada image baru, pakai image lama
	if imageURL == "" {
		imageURL = existingProduct.ImageURL
	}

	product := domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		Category:    *category,
		ImageURL:    imageURL,
	}

	updatedProduct, err := u.productRepo.UpdateProduct(id, product)
	if err != nil {
		return nil, err
	}

	updatedProduct.Category = *category

	return &domain.UpdateProductResponse{
		Message: "Product updated successfully",
		Product: *dto.ToProductResponse(updatedProduct),
	}, nil
}

func (u *productUsecase) DeleteProduct(id uint) (*domain.DeleteProductResponse, error) {
	// Get product untuk ambil image URL
	existingProduct, err := u.productRepo.GetProductByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Delete product dari database
	err = u.productRepo.DeleteProduct(id)
	if err != nil {
		return nil, err
	}

	// Hapus image dari disk
	if existingProduct.ImageURL != "" {
		deleteFile(existingProduct.ImageURL)
	}

	return &domain.DeleteProductResponse{
		Message: "Product deleted successfully",
	}, nil
}

func (u *productUsecase) ReduceStock(productID uint, qty int) error {
	product, err := u.productRepo.GetProductByID(productID)
	if err != nil {
		return err
	}
	if product.Stock < qty {
		return errors.New("stock not enough")
	}
	product.Stock -= qty
	_, err = u.productRepo.UpdateProduct(productID, *product)
	return err
}

// Helper untuk hapus file
func deleteFile(filePath string) {
	if filePath != "" {
		os.Remove(filePath)
	}
}

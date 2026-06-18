package handlers

import (
	"time"

	"github.com/dev-joaovitor/despensa-digital/models"
)

// auth
type LoginDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type SendRecoveryCodeDTO struct {
	Email string `json:"email"`
}

type VerifyRecoveryCodeDTO struct {
	Email string `json:"email"`
	Code string `json:"code"`
}

type ChangePasswordDTO struct {
	NewPassword string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

// users
type CreateUserDTO struct {
	FullName string `json:"full_name"`
	Email string `json:"email"`
	Password string `json:"password"`

	InvitationCode *string `json:"invitation_code"`
	HouseholdName *string `json:"household_name"`
}

type UpdateUserDTO struct {
	FullName string `json:"full_name"`
	Email string `json:"email"`

	Code string `json:"code"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

// households
type UpdateHouseholdDTO struct {
	Name string `json:"name"`
}

// establishments
type CreateEstablishmentDTO struct {
	Name string `json:"name"`
}

type UpdateEstablishmentDTO struct {
	Name string `json:"name"`
}

// brands
type CreateBrandDTO struct {
	Name string `json:"name"`
}

type UpdateBrandDTO struct {
	Name string `json:"name"`
}

// categories
type CreateCategoryDTO struct {
	Name string `json:"name"`
}

type UpdateCategoryDTO struct {
	Name string `json:"name"`
}

// price observations
type CreatePriceObservationDTO struct {
	ProductID int64 `json:"product_id"`
	EstablishmentID int64 `json:"establishment_id"`
	Price float64 `json:"price"`
}

type HistoryPriceObservationsDTO struct {
	ID int64 `json:"id"`
	ProductID int64 `json:"product_id"`
	Establishment models.Establishment `json:"establishment"`
	ObservedPrice float64 `json:"observed_price"`
	ObservedAt *time.Time `json:"observed_at"`
}

type ListPriceObservationsDTO struct {
	Product struct {
		ID int64 `json:"id"`
		Name string `json:"name"`
		Brand struct {
			Name string `json:"name"`
		} `json:"brand"`
		Measurement struct {
			Size int64 `json:"size"`
			Acronym string `json:"acronym"`
		} `json:"measurement"`
	} `json:"product"`
	Current struct {
		ObservedPrice float64 `json:"observed_price"`
		ObservedAt time.Time `json:"observed_at"`
		Establishment struct {
			Name string `json:"name"`
		} `json:"establishment"`
	} `json:"current"`
	AverageObservedPrice float64 `json:"average_observed_price"`
	Lowest struct {
		ObservedPrice float64 `json:"observed_price"`
		ObservedAt time.Time `json:"observed_at"`
		Establishment struct {
			Name string `json:"name"`
		} `json:"establishment"`
	} `json:"lowest"`
}

// products
type CreateProductDTO struct {
	Name string `json:"name"`
	BrandID int64 `json:"brand_id"`
	MeasurementID int64 `json:"measurement_id"`
	CategoryID int64 `json:"category_id"`
	UnitSize int64 `json:"unit_size"`
}

type ListProductsDTO struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Brand struct {
		Name string `json:"name"`
	} `json:"brand"`
	Category struct {
		Name string `json:"name"`
	} `json:"category"`
	Measurement struct {
		Size int64 `json:"size"`
		Acronym string `json:"acronym"`
	} `json:"measurement"`
}

// shopping list
type CreateShoppingItemDTO struct {
	ProductID int64 `json:"product_id"`
	Quantity int `json:"quantity"`
}

type UpdateShoppingItemDTO struct {
	Quantity int `json:"quantity"`
}

type ListShoppingItemsDTO struct {
	ID int64 `json:"id"`
	Product struct {
		ID int64 `json:"id"`
		Name string `json:"name"`
		Brand struct {
			Name string `json:"name"`
		} `json:"brand"`
		Measurement struct {
			Size int64 `json:"size"`
			Acronym string `json:"acronym"`
		} `json:"measurement"`
	} `json:"product"`
	Quantity int `json:"quantity"`
	IsChecked bool `json:"is_checked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubmitShoppingListDTO struct {
	Items []struct{
		ProductID int64 `json:"product_id"`
		EstablishmentID int64 `json:"establishment_id"`
		ExpirationDate string `json:"expiration_date"`
		Price float64 `json:"price"`
		Quantity int `json:"quantity"`
	} `json:"items"`
}

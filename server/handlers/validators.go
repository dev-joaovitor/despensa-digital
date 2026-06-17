package handlers

import (
	"errors"
	"fmt"
	"strings"
)

// users
func CreateUserValidator(user *CreateUserDTO) error {
	validationErrors := []string{}

	fullname := strings.TrimSpace(user.FullName)
	if fullname == "" {
		validationErrors = append(validationErrors, "Nome completo é obrigatório.")
	}

	if len(fullname) < 4 || len(fullname) > 100 {
		validationErrors = append(validationErrors, "Nome completo deve ter entre 4 a 100 caracteres.")
	}

	password := user.Password
	if password == "" {
		validationErrors = append(validationErrors, "Senha é obrigatória.")
	}

	if len(password) < 6 {
		validationErrors = append(validationErrors, "Senha deve ter pelo menos 6 caracteres.")
	}

	email := strings.TrimSpace(user.Email)
	if email == "" {
		validationErrors = append(validationErrors, "Email é obrigatório.")
	}

	if len(email) > 254 {
		validationErrors = append(validationErrors, "Email é muito grande.")
	}

	if user.HouseholdName != nil && user.InvitationCode == nil {
		if strings.TrimSpace(*user.HouseholdName) == "" {
			validationErrors = append(validationErrors, "É necessário criar uma residência quando não se tem um convite.")
		}
	}

	if user.InvitationCode != nil && user.HouseholdName == nil {
		if strings.TrimSpace(*user.InvitationCode) == "" {
			validationErrors = append(validationErrors, "É necessário ter um convite se não for criar uma residência.")
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

func UpdateUserValidator(user *UpdateUserDTO) error {
	validationErrors := []string{}

	fullname := strings.TrimSpace(user.FullName)
	if fullname != "" {
		if len(fullname) < 4 || len(fullname) > 100 {
			validationErrors = append(validationErrors, "Nome completo deve ter entre 4 a 100 caracteres.")
		}
	}

	email := strings.TrimSpace(user.Email)
	if email != "" {
		if len(email) > 254 {
			validationErrors = append(validationErrors, "Email é muito grande.")
		}
	}

	newPassword := user.NewPassword
	if newPassword != "" {
		if user.OldPassword == "" && user.Code == "" {
			validationErrors = append(
				validationErrors,
				"Você deve verificar sua identidade digitando a senha antiga ou confirmando pelo email.",
			)
		}

		if len(newPassword) < 6 {
			validationErrors = append(validationErrors, "Senha deve ter pelo menos 6 caracteres.")
		}

		if newPassword != user.NewPasswordConfirmation {
			validationErrors = append(validationErrors, "Nova senha deve ser igual a confirmação.")
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

// households
func UpdateHouseholdValidator(household *UpdateHouseholdDTO) error {
	validationErrors := []string{}

	name := strings.TrimSpace(household.Name)
	if name != "" {
		if len(name) < 4 || len(name) > 100 {
			validationErrors = append(validationErrors, "Nome deve ter entre 4 a 100 caracteres.")
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

// establishments
func CreateEstablishmentValidator(establishment *CreateEstablishmentDTO) error {
	validationErrors := []string{}

	name := strings.TrimSpace(establishment.Name)
	if name == "" {
		validationErrors = append(validationErrors, "Nome é obrigatório.")
	}

	if len(name) < 4 || len(name) > 100 {
		validationErrors = append(validationErrors, "Nome deve ter entre 4 a 100 caracteres.")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

func UpdateEstablishmentValidator(establishment *UpdateEstablishmentDTO) error {
	validationErrors := []string{}

	name := strings.TrimSpace(establishment.Name)
	if name != "" {
		if len(name) < 4 || len(name) > 100 {
			validationErrors = append(validationErrors, "Nome deve ter entre 4 a 100 caracteres.")
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

// brands
func CreateBrandValidator(brand *CreateBrandDTO) error {
	validationErrors := []string{}

	name := strings.TrimSpace(brand.Name)
	if name == "" {
		validationErrors = append(validationErrors, "Nome é obrigatório.")
	}

	if len(name) < 4 || len(name) > 100 {
		validationErrors = append(validationErrors, "Nome deve ter entre 4 a 100 caracteres.")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

func UpdateBrandValidator(brand *UpdateBrandDTO) error {
	validationErrors := []string{}

	name := strings.TrimSpace(brand.Name)
	if name != "" {
		if len(name) < 4 || len(name) > 100 {
			validationErrors = append(validationErrors, "Nome deve ter entre 4 a 100 caracteres.")
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

// categories
func CreateCategoryValidator(category *CreateCategoryDTO) error {
	validationErrors := []string{}

	name := strings.TrimSpace(category.Name)
	if name == "" {
		validationErrors = append(validationErrors, "Nome é obrigatório.")
	}

	if len(name) < 4 || len(name) > 100 {
		validationErrors = append(validationErrors, "Nome deve ter entre 4 a 100 caracteres.")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

func UpdateCategoryValidator(category *UpdateCategoryDTO) error {
	validationErrors := []string{}

	name := strings.TrimSpace(category.Name)
	if name != "" {
		if len(name) < 4 || len(name) > 100 {
			validationErrors = append(validationErrors, "Nome deve ter entre 4 a 100 caracteres.")
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

// price observations
func CreatePriceObservationValidator(observation *CreatePriceObservationDTO) error {
	validationErrors := []string{}

	productId := observation.ProductID
	if productId == 0 {
		validationErrors = append(validationErrors, "Produto é obrigatório.")
	}

	establishmentId := observation.EstablishmentID
	if establishmentId == 0 {
		validationErrors = append(validationErrors, "Estabelecimento é obrigatório.")
	}

	price := observation.Price
	if price == 0 {
		validationErrors = append(validationErrors, "Preço é obrigatório.")
	}

	if price < 0 {
		validationErrors = append(validationErrors, "Preço deve ser positivo.")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

// products
func CreateProductValidator(product *CreateProductDTO) error {
	validationErrors := []string{}

	name := strings.TrimSpace(product.Name)
	if name == "" {
		validationErrors = append(validationErrors, "Nome é obrigatório.")
	}

	if len(name) < 4 || len(name) > 100 {
		validationErrors = append(validationErrors, "Nome deve ter entre 4 a 100 caracteres.")
	}

	if product.BrandID <= 0 {
		validationErrors = append(validationErrors, "Marca é obrigatória")
	}

	if product.MeasurementID <= 0 {
		validationErrors = append(validationErrors, "Unidade de medida é obrigatória")
	}

	if product.CategoryID <= 0 {
		validationErrors = append(validationErrors, "Categoria é obrigatória")
	}

	if product.UnitSize <= 0 {
		validationErrors = append(validationErrors, "Tamanho de unidade é obrigatório")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

// shopping list
func CreateShoppingItemValidator(item *CreateShoppingItemDTO) error {
	validationErrors := []string{}

	if item.ProductID <= 0 {
		validationErrors = append(validationErrors, "Produto é obrigatório")
	}

	if item.Quantity <= 0 {
		validationErrors = append(validationErrors, "Quantidade é obrigatório. Quantidade deve ser maior que zero.")
	}

	if item.Quantity > 9999 {
		validationErrors = append(validationErrors, "Quantidade não deve exceder 9999.")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

func UpdateShoppingItemValidator(item *UpdateShoppingItemDTO) error {
	validationErrors := []string{}

	if item.Quantity <= 0 {
		validationErrors = append(validationErrors, "Quantidade é obrigatório. Quantidade deve ser maior que zero.")
	}

	if item.Quantity > 9999 {
		validationErrors = append(validationErrors, "Quantidade não deve exceder 9999.")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

func SubmitShoppingListValidator(list *SubmitShoppingListDTO) error {
	validationErrors := []string{}

	if len(list.Items) == 0 {
		validationErrors = append(validationErrors, "Deve haver pelo menos 1 item na lista.")
	}

	for index, item := range list.Items {
		if item.ItemID <= 0 {
			validationErrors = append(
				validationErrors,
				fmt.Sprintf("Id do item de índice %d é inválido", index),
			)
		}

		if item.EstablishmentID <= 0 {
			validationErrors = append(
				validationErrors,
				fmt.Sprintf("Id do estabelecimento de índice %d é inválido", index),
			)
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

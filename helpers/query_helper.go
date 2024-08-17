package helpers

import (
	"aswadwk/chatai/models"

	"gorm.io/gorm"
)

func Paginate(query models.QueryModel) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := query.Page
		if page <= 0 {
			page = 1
		}

		pageSize := query.PerPage
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func QueryPaginate(db *gorm.DB, model interface{}, result interface{}, page, perPage int, scopes ...func(*gorm.DB) *gorm.DB) (models.QueryModelResponse, error) {
	var response models.QueryModelResponse
	var total int64

	query := db.Model(model)
	for _, scope := range scopes {
		query = query.Scopes(scope)
	}

	// Count the total number of records with scopes applied
	query.Count(&total)

	if total == 0 {
		response = models.QueryModelResponse{
			Total:    0,
			PerPage:  0,
			CurPage:  0,
			LastPage: 0,
			From:     0,
			To:       0,
			Length:   0,
			Data:     result,
		}
		return response, nil
	}

	if perPage == 0 {
		perPage = 10
	}

	if page == 0 {
		page = 1
	}
	// Calculate the total number of pages
	lastPage := (int(total) + perPage - 1) / perPage

	// Fetch the records for the current page with scopes applied
	offset := (page - 1) * perPage
	query = query.Limit(perPage).Offset(offset)
	query.Find(&result)

	response = models.QueryModelResponse{
		Total:    int(total),
		PerPage:  perPage,
		CurPage:  page,
		LastPage: lastPage,
		From:     offset + 1,
		To:       offset + int(query.RowsAffected),
		Length:   int(query.RowsAffected),
		Data:     result,
	}

	return response, nil
}

func AgeGreaterThan(age int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("age > ?", age)
	}
}

// GroupBy is a function to group the records by a specific field.
func GroupBy(field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Group(field)
	}
}

// Select is a function to select specific fields from the records.
func Select(fields []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fields)
	}
}

func Count(model interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(model).Select("count(*) as total")
	}
}

func SearchByName(query models.QueryModel) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if query.Search != "" {
			return db.Where("name LIKE ?", "%"+query.Search+"%")
		}
		return db
	}
}

func SearchBy(field, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" LIKE ?", "%"+value+"%")
	}
}

func OrderBy(field, order string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if order == "desc" {
			return db.Order(field + " DESC")
		}
		return db.Order(field + " ASC")
	}
}

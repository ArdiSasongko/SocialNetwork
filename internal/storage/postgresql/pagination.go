package postgresql

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

type Pagination struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=10"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (p Pagination) Validate() error {
	return Validate.Struct(p)
}

func (p Pagination) Parse(r *http.Request) (Pagination, error) {
	queryString := r.URL.Query()

	limit := queryString.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return p, err
		}
		p.Limit = l
	}

	offset := queryString.Get("offset")
	if offset != "" {
		of, err := strconv.Atoi(offset)
		if err != nil {
			return p, err
		}
		p.Offset = of
	}

	sort := queryString.Get("sort")
	if sort != "" {
		p.Sort = sort
	}

	return p, nil
}

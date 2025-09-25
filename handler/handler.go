package handler

import (
	"fmt"
	"net/http"

	"github.com/MatheusJusto/gopportunities/config"
	"github.com/MatheusJusto/gopportunities/schemas"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func errParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type CreateOpeningRequest struct {
	Role     string  `json:"role"`
	Company  string  `json:"company"`
	Location string  `json:"location"`
	Remote   *bool   `json:"remote"`
	Link     string  `json:"link"`
	Salary   float64 `json:"salary"`
}

func (r *CreateOpeningRequest) Validate() error {
	if r.Role == "" {
		return errParamIsRequired("Role", "string")
	}
	if r.Company == "" {
		return errParamIsRequired("Company", "string")
	}
	if r.Location == "" {
		return errParamIsRequired("Location", "string")
	}
	if r.Link == "" {
		return errParamIsRequired("Link", "string")
	}
	if r.Remote == nil {
		return errParamIsRequired("Remote", "boolean")
	}
	if r.Salary <= 0 {
		return errParamIsRequired("Salary", "Number")
	}
	return nil
}

var (
	logger *config.Logger
	db     *gorm.DB
)

func InitializeHandler() {
	logger = config.GetLooger("handler")
	db = config.GetSQLite()
}

func CreateOpeningHandler(ctx *gin.Context) {
	request := CreateOpeningRequest{}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := request.Validate(); err != nil {
		logger.Errorf("validation error: %v", err.Error())
		//TRATAIVA DE ERRO COM O GIN
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//tratativa de erro com um metodo criado manualmente em request.go
		sendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	opening := schemas.Opening{
		Role:         request.Role,
		Company:      request.Company,
		Localization: request.Company,
		Remote:       *request.Remote,
		Link:         request.Link,
		Salary:       int64(request.Salary),
	}

	if err := db.Create(&opening).Error; err != nil {
		logger.Errorf("error creating opening: %v", err.Error())
		sendError(ctx, http.StatusInternalServerError, "error creating opening on database")
		return
	}

	sendSuccess(ctx, "create-openinng", opening)
}

func ShowOpeningHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "GET oppenings",
	})
}

func ShowOpeningsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "GET oppenings",
	})
}

func DeleteOpeningHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "DELETE oppening",
	})
}

func UpdateOpeningHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "PUT oppening",
	})
}

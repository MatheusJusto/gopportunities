package handler

import (
	"fmt"
	"net/http"

	"github.com/MatheusJusto/gopportunities/config"
	"github.com/MatheusJusto/gopportunities/schemas"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ======================= HELPERS =======================

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

// ======================= SWAGGER STRUCT =======================

// swagger:response OpeningResponse
type OpeningSwagger struct {
	ID           uint   `json:"id"`
	Role         string `json:"role"`
	Company      string `json:"company"`
	Localization string `json:"localization"`
	Remote       bool   `json:"remote"`
	Link         string `json:"link"`
	Salary       int64  `json:"salary"`
}

// ======================= HANDLER GLOBALS =======================

var (
	logger *config.Logger
	db     *gorm.DB
)

func InitializeHandler() {
	logger = config.GetLooger("handler")
	db = config.GetSQLite()
}

// ======================= CREATE =======================

// CreateOpeningHandler cria uma nova vaga
// @BasePath /api/v1
// @Summary Create Opening
// @Description Cria uma nova vaga de trabalho
// @Tags Openings
// @Accept json
// @Produce json
// @Param opening body handler.CreateOpeningRequest true "Opening Data"
// @Success 201 {object} handler.OpeningSwagger "Opening criada com sucesso"
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro interno no servidor"
// @Router /openings [post]
func CreateOpeningHandler(ctx *gin.Context) {
	request := CreateOpeningRequest{}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := request.Validate(); err != nil {
		logger.Errorf("validation error: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		sendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	opening := schemas.Opening{
		Role:         request.Role,
		Company:      request.Company,
		Localization: request.Location,
		Remote:       *request.Remote,
		Link:         request.Link,
		Salary:       int64(request.Salary),
	}

	if err := db.Create(&opening).Error; err != nil {
		logger.Errorf("error creating opening: %v", err.Error())
		sendError(ctx, http.StatusInternalServerError, "error creating opening on database")
		return
	}

	sendSuccessCreated(ctx, "create-opening", opening)
}

// ======================= READ =======================

// ShowOpeningHandler retorna uma vaga pelo ID
// @Summary Get Opening by ID
// @Description Busca uma vaga de trabalho pelo ID
// @Tags Openings
// @Accept json
// @Produce json
// @Param id query string true "Opening ID"
// @Success 200 {object} handler.OpeningSwagger "Opening encontrada"
// @Failure 400 {object} map[string]string "ID não informado"
// @Failure 404 {object} map[string]string "Opening não encontrada"
// @Router /openings/opening [get]
func ShowOpeningHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		sendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
		return
	}
	opening := schemas.Opening{}
	if err := db.First(&opening, id).Error; err != nil {
		sendError(ctx, http.StatusNotFound, fmt.Sprintf("opening with id: %s not found", id))
		return
	}
	sendSuccess(ctx, "opening-found", opening)
}

// ShowOpeningsHandler retorna todas as vagas
// @Summary List Openings
// @Description Lista todas as vagas de trabalho
// @Tags Openings
// @Accept json
// @Produce json
// @Success 200 {array} handler.OpeningSwagger "Lista de vagas"
// @Failure 404 {object} map[string]string "Nenhuma vaga encontrada"
// @Router /openings [get]
func ShowOpeningsHandler(ctx *gin.Context) {
	openings := []schemas.Opening{}

	if err := db.Find(&openings).Error; err != nil {
		sendError(ctx, http.StatusNotFound, "Nenhuma vaga encontrada")
		return
	}
	sendSuccess(ctx, "openings-found", openings)
}

// ======================= DELETE =======================

// DeleteOpeningHandler deleta uma vaga pelo ID
// @Summary Delete Opening
// @Description Deleta uma vaga de trabalho pelo ID
// @Tags Openings
// @Accept json
// @Produce json
// @Param id query string true "Opening ID"
// @Success 204 {string} string "Vaga deletada com sucesso"
// @Failure 400 {object} map[string]string "ID não informado"
// @Failure 404 {object} map[string]string "Opening não encontrada"
// @Failure 500 {object} map[string]string "Erro interno no servidor"
// @Router /openings [delete]
func DeleteOpeningHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		sendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
		return
	}

	opening := schemas.Opening{}
	if err := db.First(&opening, id).Error; err != nil {
		sendError(ctx, http.StatusNotFound, fmt.Sprintf("opening with id: %s not found", id))
		return
	}

	if err := db.Delete(&opening).Error; err != nil {
		sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error deleting opening with id: %s", id))
		return
	}
	sendSuccessDeleted(ctx, "deleted-opening", opening)
}

// ======================= UPDATE =======================

// UpdateOpeningHandler atualiza uma vaga pelo ID
// @Summary Update Opening
// @Description Atualiza os dados de uma vaga de trabalho pelo ID
// @Tags Openings
// @Accept json
// @Produce json
// @Param id query string true "Opening ID"
// @Param opening body handler.UpdateOpeningRequest true "Dados da vaga para atualizar"
// @Success 200 {object} handler.OpeningSwagger "Vaga atualizada com sucesso"
// @Failure 400 {object} map[string]string "Erro de validação ou ID não informado"
// @Failure 404 {object} map[string]string "Opening não encontrada"
// @Failure 500 {object} map[string]string "Erro interno no servidor"
// @Router /openings [put]
func UpdateOpeningHandler(ctx *gin.Context) {
	request := UpdateOpeningRequest{}

	if err := ctx.BindJSON(&request); err != nil {
		sendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := request.Validate(); err != nil {
		logger.Errorf("validation error %v", err.Error())
		sendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id := ctx.Query("id")
	if id == "" {
		sendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
		return
	}

	opening := schemas.Opening{}
	if err := db.First(&opening, id).Error; err != nil {
		sendError(ctx, http.StatusNotFound, "opening not found")
		return
	}

	if request.Role != "" {
		opening.Role = request.Role
	}
	if request.Company != "" {
		opening.Company = request.Company
	}
	if request.Location != "" {
		opening.Localization = request.Location
	}
	if request.Remote != nil {
		opening.Remote = *request.Remote
	}
	if request.Link != "" {
		opening.Link = request.Link
	}
	if request.Salary > 0 {
		opening.Salary = int64(request.Salary)
	}

	if err := db.Save(&opening).Error; err != nil {
		logger.Errorf("error updating opening: %v", err.Error())
		sendError(ctx, http.StatusInternalServerError, "error updating opening")
		return
	}

	sendSuccess(ctx, "update-opening", opening)
}

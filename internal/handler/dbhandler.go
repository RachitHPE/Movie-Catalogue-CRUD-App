package handler

import (
	"catalogue-app/internal/controller"
	"catalogue-app/internal/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	dbController controller.ControllerIntfc
}

func NewMovieHandler(dbController controller.ControllerIntfc) *MovieHandler {
	return &MovieHandler{dbController: dbController}
}

func (handler MovieHandler) GetMovies(ginCtx *gin.Context) {
	result, err := handler.dbController.GetMovies(ginCtx.Request.Context())
	if err != nil {
		ginCtx.AbortWithError(http.StatusNotFound, err)
		return
	}

	ginCtx.JSON(http.StatusOK, &result)
}

func (handler MovieHandler) GetMovieByID(ginCtx *gin.Context) {
	id := ginCtx.Param("id")

	result, err := handler.dbController.GetMovieByID(ginCtx.Request.Context(), id)
	if err != nil {
		ginCtx.AbortWithError(http.StatusNotFound, err)
		return
	}

	ginCtx.JSON(http.StatusOK, &result)
}

func (handler MovieHandler) CreateMovie(ginCtx *gin.Context) {
	var movieInfo model.MovieInfo

	if err := ginCtx.BindJSON(&movieInfo); err != nil {
		ginCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := handler.dbController.CreateMovie(ginCtx.Request.Context(), movieInfo)
	if err != nil {
		ginCtx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ginCtx.JSON(http.StatusCreated, &result)
}

func (handler MovieHandler) UpdateMovie(ginCtx *gin.Context) {
	id := ginCtx.Param("id")

	var movieInfo model.MovieInfo

	if err := ginCtx.BindJSON(&movieInfo); err != nil {
		ginCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := handler.dbController.UpdateMovie(ginCtx.Request.Context(), movieInfo, id)
	if err != nil {
		ginCtx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ginCtx.JSON(http.StatusOK, &result)
}

func (handler MovieHandler) DeleteMovie(ginCtx *gin.Context) {
	id := ginCtx.Param("id")

	err := handler.dbController.DeleteMovie(ginCtx.Request.Context(), id)
	if err != nil {
		ginCtx.AbortWithError(http.StatusNotFound, err)
		return
	}

	ginCtx.JSON(http.StatusNoContent, gin.H{"status": "deleted"})
}

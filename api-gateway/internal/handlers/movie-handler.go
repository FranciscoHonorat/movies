package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/FranciscoHonorat/movies/proto"
	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	client proto.MovieServiceClient
}

type CreateMovieStruct struct {
	Title string `json:"title" binding:"required"`
	Year  string `json:"year" binding:"required"`
}

type ListMovieStruct struct {
	Data  []*proto.Movie `json:"data"`
	Page  int32          `json:"page"`
	Limit int32          `json:"limit"`
	Total int32          `json:"total"`
}

func NewMovieHandler(client proto.MovieServiceClient) *MovieHandler {
	return &MovieHandler{client: client}
}

func (m *MovieHandler) GetMovie(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	movie, err := m.client.GetMovie(c.Request.Context(), &proto.GetMovieRequest{Id: int32(idInt)})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (m *MovieHandler) ListMovie(c *gin.Context) {
	title := c.Query("title")
	year := c.Query("year")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	sort := c.DefaultQuery("title", "year")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page (must be >= 1)"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit (must be between 1 and 100)"})
		return
	}
	validSorts := map[string]bool{"title": true, "year": true}
	if !validSorts[sort] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort (must be title, year)"})
		return
	}

	listReq := proto.ListMovieRequest{
		Title:  title,
		Year:   year,
		Page:   int32(pageInt),
		Limit:  int32(limitInt),
		SortBy: sort,
	}

	resp, err := m.client.ListMovie(c.Request.Context(), &listReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := ListMovieStruct{
		Data:  resp.Movie,
		Page:  int32(pageInt),
		Limit: int32(limitInt),
	}

	c.JSON(http.StatusOK, response)
}

func (m *MovieHandler) CreateMovie(c *gin.Context) {
	var req CreateMovieStruct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or missing required fields"})
		return
	}

	movieInput := proto.CreateMovieRequest{
		Title: req.Title,
		Year:  req.Year,
	}

	createMovie, err := m.client.CreateMovie(c.Request.Context(), &movieInput)

	if err != nil {
		slog.Error("CreateMovie error", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, createMovie)

}

func (m *MovieHandler) DeleteMovie(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = m.client.DeleteMovie(c.Request.Context(), &proto.DeleteMovieRequest{Id: int32(idInt)})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

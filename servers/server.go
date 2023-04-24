package servers

import (
	"net/http"
	"strconv"

	"key-value-ts/domains/entities"
	"key-value-ts/domains/storage"
	"key-value-ts/requests"

	"github.com/gin-gonic/gin"
)

// Server enclosure for the server object
type Server struct {
	storer storage.Storer
	router *gin.Engine
}

// New returns a new instance of server with the give storer
func New(storer storage.Storer) *Server {
	router := gin.Default()
	s := &Server{
		storer: storer,
	}
	s.router = router
	router.PUT("/", s.putSequence)
	router.GET("/", s.getSequence)
	return s
}

// Start trigger the server to start on the defined port
func (s *Server) Start(port string) error {
	return s.router.Run(port)
}

func (s *Server) putSequence(c *gin.Context) {
	var req requests.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	entity := entities.Sequence{
		Key:       req.Key,
		Timestamp: req.Timestamp,
		Value:     req.Value,
	}

	if err := s.storer.Save(entity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "sequence saved successfully",
	})
}

func (s *Server) getSequence(c *gin.Context) {
	key := c.Query("key")
	tsStr := c.Query("timestamp")

	ts, err := strconv.Atoi(tsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := s.storer.Get(key, int64(ts))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"value": res,
	})
}

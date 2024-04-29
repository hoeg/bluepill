package morpheus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoeg/bluepill/internal/values"
	"k8s.io/api/admission/v1beta1"
)

type API struct {
	whitlistedIP values.Whitelist
	enforce      bool
}

func NewAPI(config EnforcementConfig) *API {
	return &API{
		whitlistedIP: config.Whitelist,
		enforce:      config.Enforce,
	}
}

func (a *API) PostPill(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer c.Request.Body.Close()

	admReview := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal Server Error - unmarshaling request failed with %s", err)})
		return
	}

	admReview, err = a.Mutate(admReview)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal Server Error - mutation failed with %s", err)})
		return
	}
	c.JSON(http.StatusOK, admReview)
}

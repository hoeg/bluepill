package morpheus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type API struct {
	whitlistedIP []string
	enforce      bool
}

func NewAPI() *API {
	return &API{}
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

	if admReview.Request == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request - missing request field"})
		return
	}

	resp := &v1beta1.AdmissionResponse{}
	resp.Allowed = true
	resp.UID = admReview.Request.UID

	ingress := ""
	//check if label is set
	if err := json.Unmarshal(admReview.Request.Object.Raw, &ingress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to unmarshal object"})
		return
	}

	resp.AuditAnnotations = map[string]string{
		"bluepill": "internal enforced",
	}
	pt := v1beta1.PatchTypeJSONPatch
	resp.PatchType = &pt

	p := []map[string]string{
		{
			"op":    "add",
			"path":  "/metadata/annotations/nginx.ingress.kubernetes.io~1whitelist-source-range",
			"value": "127.0.0.1",
		},
	}
	resp.Patch, err = json.Marshal(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal Server Error - marshaling response failed with %s", err)})
		return
	}
	resp.Result = &metav1.Status{
		Status: "Success",
	}
	admReview.Response = resp
	c.JSON(http.StatusOK, admReview)
}

package morpheus

import (
	"encoding/json"
	"fmt"

	"k8s.io/api/admission/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *API) Mutate(admReview v1beta1.AdmissionReview) (v1beta1.AdmissionReview, error) {
	if admReview.Request == nil {
		return v1beta1.AdmissionReview{}, fmt.Errorf("missing request field")
	}

	resp := &v1beta1.AdmissionResponse{}
	resp.Allowed = true
	resp.UID = admReview.Request.UID

	ingress := networkingv1.Ingress{}
	if err := json.Unmarshal(admReview.Request.Object.Raw, &ingress); err != nil {
		return v1beta1.AdmissionReview{}, fmt.Errorf("unable to unmarshal object")
	}

	external := ingress.ObjectMeta.Labels["external"] != ""
	if a.enforce && !external {
		resp.AuditAnnotations = map[string]string{
			"bluepill": "internal enforced",
		}
		pt := v1beta1.PatchTypeJSONPatch
		resp.PatchType = &pt

		p := []map[string]string{
			{
				"op":    "add",
				"path":  "/metadata/annotations/nginx.ingress.kubernetes.io~1whitelist-source-range",
				"value": a.whitlistedIP.Value(),
			},
		}
		var err error
		resp.Patch, err = json.Marshal(p)
		if err != nil {
			return v1beta1.AdmissionReview{}, fmt.Errorf("Internal Server Error - marshaling response failed with %s", err)
		}
	}
	resp.Result = &metav1.Status{
		Status: "Success",
	}
	admReview.Response = resp
	return admReview, nil
}

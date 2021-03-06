package helm

import (
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/helm"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

type ReleaseStatuser interface {
	ReleaseStatus(rlsName string, opts ...helm.StatusOption) (*rls.GetReleaseStatusResponse, error)
}

type ReleaseResources struct {
	client ReleaseStatuser
}

func New(client ReleaseStatuser) *ReleaseResources {
	return &ReleaseResources{client: client}
}

func (r *ReleaseResources) Get(name string) (string, error) {
	resp, err := r.client.ReleaseStatus(name)
	if err != nil {
		return "", errors.Wrap(err, "helm client error")
	}

	return resp.GetInfo().GetStatus().GetResources(), nil
}

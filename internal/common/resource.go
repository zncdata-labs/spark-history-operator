package common

import (
	"context"

	sparkv1alpha1 "github.com/zncdatadev/spark-k8s-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var log = ctrl.Log.WithName("resourceFetcher")

type ResourceClient struct {
	Ctx       context.Context
	Client    client.Client
	Namespace string
}

func (r *ResourceClient) Get(obj client.Object) error {
	name := obj.GetName()
	kind := obj.GetObjectKind()
	if err := r.Client.Get(r.Ctx, client.ObjectKey{Namespace: r.Namespace, Name: name}, obj); err != nil {
		opt := []any{"ns", r.Namespace, "name", name, "kind", kind}
		if apierrors.IsNotFound(err) {
			log.Error(err, "Fetch resource NotFound", opt...)
		} else {
			log.Error(err, "Fetch resource occur some unknown err", opt...)
		}
		return err
	}
	return nil
}

type InstanceAttributes interface {
	GetClusterConfig() *sparkv1alpha1.ClusterConfigSpec
	GetClusterOperation() *sparkv1alpha1.ClusterOperationSpec
}

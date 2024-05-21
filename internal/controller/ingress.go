package controller

import (
	"context"

	sparkv1alpha1 "github.com/zncdatadev/spark-k8s-operator/api/v1alpha1"
	"github.com/zncdatadev/spark-k8s-operator/internal/common"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IngressReconciler struct {
	common.GeneralResourceStyleReconciler[*sparkv1alpha1.SparkHistoryServer, *sparkv1alpha1.RoleGroupSpec]
}

func NewIngress(
	scheme *runtime.Scheme,
	instance *sparkv1alpha1.SparkHistoryServer,
	client client.Client,
	groupName string,
	mergedLabels map[string]string,
	mergedCfg *sparkv1alpha1.RoleGroupSpec,
) *IngressReconciler {
	return &IngressReconciler{
		GeneralResourceStyleReconciler: *common.NewGeneraResourceStyleReconciler[*sparkv1alpha1.SparkHistoryServer,
			*sparkv1alpha1.RoleGroupSpec](
			scheme,
			instance,
			client,
			groupName,
			mergedLabels,
			mergedCfg),
	}
}

// Build implements the ResourceBuilder interface
func (i *IngressReconciler) Build(_ context.Context) (client.Object, error) {
	ingressSpec := i.getIngressSpec()
	pathTypePrefix := v1.PathTypePrefix
	ing := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createIngName(i.Instance.Name, i.GroupName),
			Namespace: i.Instance.Namespace,
			Labels:    i.MergedLabels,
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: ingressSpec.Host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathTypePrefix,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: createServiceName(i.Instance.Name, i.GroupName),
											Port: v1.ServiceBackendPort{
												Name: SparkHistoryHTTPPortName,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return ing, nil
}

// get ingress spec
func (i *IngressReconciler) getIngressSpec() *sparkv1alpha1.IngressSpec {
	spec := i.Instance.Spec.ClusterConfig.Ingress
	if spec == nil {
		spec = &sparkv1alpha1.IngressSpec{
			Host:    "spark-history-server.example.com",
			Enabled: true,
		}
	}
	return spec
}

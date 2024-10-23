package e2e

import (
	myapiv1 "github.com/KokoiRuby/module_07_application/api/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

// client & global var
var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func(done Done) {
	// logger
	ctrl.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(GinkoWriter)))

	// init envtest
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	// start envtest
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	// reg crd
	err = myapiv1.AddToScheme(scheme.Scheme)
	Expect(err).ToNot(HaveOccurred())

	// k8s client
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	close(done)
})

var _ = AfterSuite(func(done Done) {
	// close envtest
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

var _ = Describe("CRD Test", func() {
	ctx := context.Background()
	It("should create CRD successfully", func() {
		//
		myResource := &myapiv1.Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-application",
				Namespace: "default",
			},
			Spec: myapiv1.ApplicationSpec{
				Deployment: myapiv1.ApplicationDeployment{
					Replicas: 1,
					Image:    "nginx",
					Port:     80,
				},
				Service: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port:       80,
							TargetPort: intstr.FromInt32(80),
						},
					},
				},
				Ingress: networkingv1.IngressSpec{
					IngressClassName: "nginx",
					Rules: []networkingv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networkingv1.IngressRuleValue{
								HTTP: &networkingv1.HTTPIngressRuleValue{
									Paths: []networkingv1.HTTPIngressPath{
										{
											Path:     "/",
											PathType: networkingv1.PathTypePrefix,
											Backend: networkingv1.IngressBackend{
												Service: &networkingv1.IngressServiceBackend{
													Name: "test-application",
													Port: networkingv1.ServiceBackendPort{
														Number: 80,
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
			},
		}
		// create res
		Expect(k8sClient.Create(ctx, myResource)).Should(Succeed())

		// get res
		fetchedResource := &myapiv1.Application{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: "test-application"}, fetchedResource)).Should(Succeed())

	})
})

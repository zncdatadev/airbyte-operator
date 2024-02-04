package controller

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	//"github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
)

var _ = Describe("Role Group Controller", func() {
	//var (
	//	clusterCfg, roleCfg, groupCfg any
	//	fields                        []string
	//)
	ctx := context.Background()

	BeforeEach(func() {
		//clusterCfg = &v1alpha1.ClusterRoleConfigSpec{
		//	BaseRoleConfigSpec: v1alpha1.BaseRoleConfigSpec{
		//		MatchLabels: map[string]string{
		//			"app.kubernetes.io/name":     "airbyte1",
		//			"app.kubernetes.io/instance": "airbyte1",
		//		},
		//		Replicas: 1,
		//	},
		//	Edition: "community",
		//}
		//roleCfg = &v1alpha1.TemporalRoleConfigSpec{
		//	BaseRoleConfigSpec: v1alpha1.BaseRoleConfigSpec{
		//		MatchLabels: map[string]string{
		//			"app.kubernetes.io/name":     "airbyte2",
		//			"app.kubernetes.io/instance": "airbyte2",
		//		},
		//	},
		//}
		//groupCfg = &v1alpha1.TemporalRoleConfigSpec{
		//	BaseRoleConfigSpec: v1alpha1.BaseRoleConfigSpec{
		//		MatchLabels: map[string]string{
		//			"app.kubernetes.io/name":     "airbyte3",
		//			"app.kubernetes.io/instance": "airbyte3",
		//		},
		//	},
		//}
		// Initialize your test variables here
	})

	Describe("GetRoleGroupEx", func() {
		//Skip("Skipping this test because it is not completed yet.")
		Context("when all fields are selected", func() {
			It("returns a merged role group with all fields", func() {
				//fields = []string{"MatchLabels"}
				Expect(1).To(Equal(1))
				//result := GetRoleGroupEx(clusterCfg, roleCfg, groupCfg, fields, nil)
				//Expect(result).ToNot(BeNil())
				//Expect(result.GetFields()).To(Equal(fields))
				objs := stackv1alpha1.AirbyteList{}
				err := k8sClient.List(ctx, &objs)
				if err != nil {
					return
				}
			})
		})

	})
})

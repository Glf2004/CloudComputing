package utils

import (
	"context"
	"net"
	"os"
	"strings"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/utils"
)

var (
	CartHostPorts     = []string{"cart.default.svc.cluster.local:8883"}
	CheckoutHostPorts = []string{"checkout.default.svc.cluster.local:8884"}
	EmailHostPorts    = []string{"email.default.svc.cluster.local:8888"}
	OrderHostPorts    = []string{"order.default.svc.cluster.local:8885"}
	PaymentHostPorts  = []string{"payment.default.svc.cluster.local:8886"}
	ProductHostPorts  = []string{"product.default.svc.cluster.local:8881"}
	UserHostPorts     = []string{"user.default.svc.cluster.local:8882"}
)

func MyWithHostPorts(hostports ...string) client.Option {
	return client.Option{F: func(o *client.Options, di *utils.Slice) {
		o.Targets = strings.Join(hostports, ",")
		o.Resolver = &discovery.SynthesizedResolver{
			ResolveFunc: func(ctx context.Context, key string) (discovery.Result, error) {
				var ins []discovery.Instance
				for _, hp := range hostports {
					// Try Kubernetes DNS first
					if _, err := net.ResolveTCPAddr("tcp", hp); err == nil {
						ins = append(ins, discovery.NewInstance("tcp", hp, discovery.DefaultWeight, nil))
						continue
					}
					// Fallback to localhost for development
					if os.Getenv("GO_ENV") == "dev" {
						localHP := strings.Replace(hp, "-service.default.svc.cluster.local", "", 1)
						if _, err := net.ResolveTCPAddr("tcp", localHP); err == nil {
							ins = append(ins, discovery.NewInstance("tcp", localHP, discovery.DefaultWeight, nil))
							continue
						}
					}
					// Unix socket fallback
					if _, err := net.ResolveUnixAddr("unix", hp); err == nil {
						ins = append(ins, discovery.NewInstance("unix", hp, discovery.DefaultWeight, nil))
						continue
					}
				}
				return discovery.Result{
					Cacheable: true,
					CacheKey:  "fixed",
					Instances: ins,
				}, nil
			},
			NameFunc: func() string { return o.Targets },
			TargetFunc: func(ctx context.Context, target rpcinfo.EndpointInfo) string {
				return o.Targets
			},
		}
	}}
}

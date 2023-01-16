package gateway

import (
	"context"
	"net/http"
	"log"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"

	"github.com/walkergriggs/enoki/proto/manifests/v1"
)

func Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Println("Dialing")
	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	log.Println("Handling")
	gw, err := newGateway(ctx, conn, []runtime.ServeMuxOption{})
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzServer(conn))
	mux.Handle("/", gw)

	s := &http.Server{
		Addr:    "localhost:8001",
		Handler: allowCORS(mux),
	}

	go func() {
		<-ctx.Done()
		s.Shutdown(context.Background())
	}()

	log.Println("Serving")
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}



func RegisterManifestServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterManifestServiceHandlerClient(ctx, mux, v1.NewManifestServiceClient(conn))
}

func RegisterManifestServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client v1.ManifestServiceClient) error {
	mux.Handle("GET", pattern_ManifestService_Get_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}

		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)

		var protoReq v1.ManifestRequest
		var metadata runtime.ServerMetadata

		if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_ManifestService_Get_0); err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		
		resp, err := client.GetManifest(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		runtime.ForwardResponseMessage(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})

	return nil
}

var (
	pattern_ManifestService_Get_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "manifests"}, ""))
	filter_ManifestService_Get_0 =  &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

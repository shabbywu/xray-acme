package api

import (
	"context"
	"fmt"
	creflect "github.com/xtls/xray-core/common/reflect"
	"github.com/xtls/xray-core/main/commands/base"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"os"
	"reflect"
	"time"
)

func dialAPIServer() (conn *grpc.ClientConn, ctx context.Context, close func()) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(apiTimeout)*time.Second)
	conn, err := grpc.DialContext(ctx, apiServerAddrPtr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		base.Fatalf("failed to dial %s", apiServerAddrPtr)
	}
	close = func() {
		cancel()
		conn.Close()
	}
	return
}

func showJSONResponse(m proto.Message) {
	if isNil(m) {
		return
	}
	if j, ok := creflect.MarshalToJson(m, true); ok {
		fmt.Println(j)
	} else {
		fmt.Fprintf(os.Stdout, "%v\n", m)
		base.Fatalf("error encode json")
	}
}

func isNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return i == nil
}

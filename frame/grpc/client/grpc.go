package client

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/qq754174349/ht/ht-frame/common/result"
	"github.com/qq754174349/ht/ht-frame/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RpcClient[C any] struct {
	serviceName string
	timeout     time.Duration
	newClient   func(conn grpc.ClientConnInterface) C
}

type rpcResult[T any] struct {
	resp T
	err  error
}

func New[T any](serviceName string, pd func(cc grpc.ClientConnInterface) T) *RpcClient[T] {
	return &RpcClient[T]{
		serviceName: serviceName,
		newClient:   pd,
		timeout:     5 * time.Second,
	}
}

// Invoke 调用远程方法，使用默认超时
func Invoke[C any, T any](
	r *RpcClient[C],
	ctx context.Context,
	handler func(client C, ctx context.Context) (T, error),
) (T, error) {
	return InvokeWithTimeout(r, ctx, r.timeout, handler)
}

// InvokeWithTimeout 调用远程方法，自定义超时
func InvokeWithTimeout[C any, T any](
	r *RpcClient[C],
	ctx context.Context,
	timeout time.Duration,
	handler func(client C, ctx context.Context) (T, error),
) (T, error) {
	var zero T
	if timeout < 1 {
		timeout = r.timeout
	}

	conn, err := getConn(r.serviceName)
	if err != nil {
		return zero, err
	}
	defer conn.Close()

	childCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	copyKeys := []any{"traceID"}
	for _, key := range copyKeys {
		if val := ctx.Value(key); val != nil {
			childCtx = context.WithValue(childCtx, key, val)
		}
	}
	resultCh := make(chan rpcResult[T], 1)

	go func() {
		resp, err := handler(r.newClient(conn), childCtx)
		resultCh <- rpcResult[T]{resp: resp, err: err}
	}()

	select {
	case <-childCtx.Done():
		return zero, childCtx.Err()
	case resultsSt := <-resultCh:
		return resultsSt.resp, resultsSt.err
	}
}

// InvokeWithResult 调用远程方法，返回基础result,使用默认超时
func InvokeWithResult[C any](
	r *RpcClient[C],
	ctx context.Context,
	handler func(client C, ctx context.Context) (*result.Result, error),
) (*result.Result, error) {

	return InvokeWithTimeout(r, ctx, r.timeout, handler)
}

// InvokeWithTimeoutResult 调用远程方法，返回基础result,使用自定义超时
func InvokeWithTimeoutResult[C any](
	r *RpcClient[C],
	ctx context.Context,
	timeout time.Duration,
	handler func(client C, ctx context.Context) (*result.Result, error),
) *result.Result {
	if timeout <= 0 {
		timeout = r.timeout
	}

	// 获取连接
	conn, err := getConn(r.serviceName)
	if err != nil {
		log.Printf("连接失败: %v", err)
		return result.New(ctx, 500, err.Error(), nil)
	}
	defer conn.Close()

	childCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, key := range []any{"traceID"} {
		if val := ctx.Value(key); val != nil {
			childCtx = context.WithValue(childCtx, key, val)
		}
	}

	respCh := make(chan *result.Result, 1)
	go func() {
		resp, err := handler(r.newClient(conn), childCtx)
		if err != nil {
			respCh <- result.New(ctx, 500, err.Error(), nil)
		} else {
			respCh <- resp
		}
	}()

	select {
	case <-ctx.Done():
		return result.New(ctx, 500, "服务器响应超时", nil)
	case r := <-respCh:
		return r
	}
}

func getConn(serviceName string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s/%s", consul.GetConfig().Consul.Addr, serviceName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	)
	if err != nil {
		return nil, fmt.Errorf("获取grpc服务连接失败: %w", err)
	}
	return conn, nil
}

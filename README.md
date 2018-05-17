#### micService 

* client端

```go
 ./client
```

* server端

```go
./server
``` 
 
例如在kit里面的stringsvc1中最简单的例子:

* 业务逻辑逻辑

服务（Service）是从业务逻辑开始的，在 Go kit 中，我们将服务以interface 作为模型

```go
// StringService provides operations on strings.
type StringService interface {
    Uppercase(string) (string, error)
    Count(string) int
}
```

这个 interface 需要有一个“实现”

```go
//接口实现
type stringService struct{}


//处理Uppercase业务
func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}
//处理Count业务
func (stringService) Count(s string) int {
	return len(s)
}

// ErrEmpty is returned when input string is empty
var ErrEmpty = errors.New("Empty string")

```

* 请求和响应

在 Go kit 中，主要的消息模式是 RPC。因此，接口（ interface ）的每一个方法都会被模型化为远程过程调用（RPC）。
对于每一个方法，我们都定义了请求和响应的结构体，捕获输入、输出各自的所有参数。

```go

//定义Uppercase的输入参数的结构
type uppercaseRequest struct {
	S string `json:"s"`
}

//定义Uppercase的输出接口
type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}


//定义Count的输入参数结构
type countRequest struct {
	S string `json:"s"`
}

//定义Count的输入结构
type countResponse struct {
	V int `json:"v"`
}

```

* 端点 （endpoint）

Go kit 通过 endpoint 提供了非常丰富的功能。

一个端点代表一个RPC，也就是我们服务接口中的一个函数。我们将编写简单的适配器，将我们的服务的每一个方法转换成端点。

```go
//go-kit中，如果使用go-kit/kit/transport/http，那么还需要把StringService封装为endpoint来供调用。
//抽象 uppercase的RPC调用
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

//RPC调用封装成了更加通用的接口，输入参数和输出参数都为interface
//抽象 len的RPC调用
func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

```

* 传输（Transports）

现在我们需要将服务暴露给外界，这样它们才能被调用。对于服务如何与外界交互，你的组织可能已经有了定论。可能你会使用 Thrift、基于 HTTP 的自定义 JSON。
Go kit支持多种开箱即用的 传输 方式。(Adding support for new ones is easy—just 对新方式的支持是非常简单的。

针对我们现在的这个微型的服务例子，我们使用基于 HTTP 的 JSON。Go kit 中提供了一个辅助结构体，在 transport/http 中。

```go
//服务启动
func main() {
	svc := stringService{}

	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//从Request解码输入参数，编码输出到ResponseWriter
//第二步就是调用的是上面生成的endpoint，第一步需要我们传入解码器，用于将Request解码为输入参数，第三部需要我们传入编码器，输出到ResponseWriter。
//Uppercase输入解码器
func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

//Count输入解码器
func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

//由于Uppercase和Count对输出的处理一样，所以可以用一个通用的编码器，将结果写入到ResponseWriter
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

```

```go
$ curl -XPOST -d'{"s":"hello, world"}' localhost:8080/uppercase
{"v":"HELLO, WORLD","err":null}
$ curl -XPOST -d'{"s":"hello, world"}' localhost:8080/count
{"v":12}
```

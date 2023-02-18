package tcputils

import (
	"context"
	"github.com/Kirov7/fayUtils/loadBanlance"
	"github.com/Kirov7/fayUtils/net/tcp"
	"io"
	"log"
	"net"
	"time"
)

func NewTcpLoadBalanceReverseProxy(c *tcp.TcpSliceRouterContext, lb loadBanlance.LoadBalance) *TcpReverseProxy {
	return func() *TcpReverseProxy {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr fail")
		}
		return &TcpReverseProxy{
			ctx:             c.Ctx,
			Addr:            nextAddr,
			KeepAlivePeriod: time.Second,
			DialTimeout:     time.Second,
		}
	}()
}

//TCP反向代理
type TcpReverseProxy struct {
	ctx                  context.Context //单次请求单独设置
	Addr                 string
	KeepAlivePeriod      time.Duration //设置
	DialTimeout          time.Duration //设置超时时间
	DialContext          func(ctx context.Context, network, address string) (net.Conn, error)
	OnDialError          func(src net.Conn, dstDialErr error)
	ProxyProtocolVersion int
}

func (trp *TcpReverseProxy) dialTimeout() time.Duration {
	if trp.DialTimeout > 0 {
		return trp.DialTimeout
	}
	return 10 * time.Second
}

var defaultDialer = new(net.Dialer)

func (trp *TcpReverseProxy) dialContext() func(ctx context.Context, network, address string) (net.Conn, error) {
	if trp.DialContext != nil {
		return trp.DialContext
	}
	return (&net.Dialer{
		Timeout:   trp.DialTimeout,     //连接超时
		KeepAlive: trp.KeepAlivePeriod, //设置连接的检测时长
	}).DialContext
}

func (trp *TcpReverseProxy) keepAlivePeriod() time.Duration {
	if trp.KeepAlivePeriod != 0 {
		return trp.KeepAlivePeriod
	}
	return time.Minute
}

//传入上游 conn，在这里完成下游连接与数据交换
func (trp *TcpReverseProxy) ServeTCP(ctx context.Context, src net.Conn) {
	//设置连接超时
	var cancel context.CancelFunc
	if trp.DialTimeout >= 0 {
		ctx, cancel = context.WithTimeout(ctx, trp.dialTimeout())
	}
	dst, err := trp.dialContext()(ctx, "tcp", trp.Addr)
	if cancel != nil {
		cancel()
	}
	if err != nil {
		trp.onDialError()(src, err)
		return
	}

	defer func() { go dst.Close() }() //记得退出下游连接

	//设置dst的 keepAlive 参数,在数据请求之前
	if ka := trp.keepAlivePeriod(); ka > 0 {
		if c, ok := dst.(*net.TCPConn); ok {
			c.SetKeepAlive(true)
			c.SetKeepAlivePeriod(ka)
		}
	}
	errc := make(chan error, 1)
	go trp.proxyCopy(errc, src, dst)
	go trp.proxyCopy(errc, dst, src)
	<-errc
}

func (trp *TcpReverseProxy) onDialError() func(src net.Conn, dstDialErr error) {
	if trp.OnDialError != nil {
		return trp.OnDialError
	}
	return func(src net.Conn, dstDialErr error) {
		log.Printf("tcpproxy: for incoming conn %v, error dialing %q: %v", src.RemoteAddr().String(), trp.Addr, dstDialErr)
		src.Close()
	}
}

func (trp *TcpReverseProxy) proxyCopy(errc chan<- error, dst, src net.Conn) {
	_, err := io.Copy(dst, src)
	errc <- err
}

## 在p2p环境中进行HTTP通信
假设你正在构建一个分布式的web应用，该应用将运行在一个P2P网络中，而不是传统的客户端-服务器模型。在这种情况下，你可能需要一种方法来在网络的节点之间进行通信。

通常，你可能会使用一些专门为P2P环境设计的协议，如BitTorrent或Gnutella。但是，如果你更喜欢使用熟悉的HTTP语义，那么这个代码就可以派上用场了。

具体来说，你可以使用这个代码中的 RoundTripper 结构体来创建一个自定义的HTTP传输层，这个传输层使用libp2p的主机来路由请求。然后，你可以使用标准的 http.Client 结构体来发起HTTP请求，只不过这些请求会被路由到P2P网络中的其他节点，而不是传统的服务器。

例如，你可能有一个节点需要获取另一个节点上的某个资源。你可以简单地使用 http.Client 来发起一个GET请求，就像你请求一个普通的web服务器那样。但是，在后台，这个请求实际上会被路由到P2P网络中的正确节点。

总的来说，这个代码的目的是让你能够在libp2p环境中使用熟悉的HTTP语义进行通信，而无需关心底层的P2P路由细节。

libp2p网络堆栈上层可以实现HTTP协议的传输
### 配置
enableP2P 


### 调试
go run .\cmd\httpcat.go --port 9001 --p2pport 9002 --upload /home/web/website/download/ --download /home/web/website/download/ -C F:\open_code\httpcat\server\conf\svr.yml

go run cmd/httpcat.go  --port 9001 --p2pport 9002   --static=/home/web/website/upload/  --upload=/home/web/website/upload/ --download=/home/web/website/upload/  -C server/conf/svr.yml
go run cmd/httpcat.go  --port 9003 --p2pport 9004   --static=/home/web/website/upload/  --upload=/home/web/website/upload/ --download=/home/web/website/upload/  -C server/conf/svr.yml


### 节点发现
与常规的“host:port”寻址不同，“p2phttp”使用对等ID,并让LibP2P负责路由，通过单个连接利用上多路由、NAT和流复用等功能。

### web接口
http://{{ip}}:{{port}}/api/v1/p2p/send_message
POST
{
"topic": "httpcat",
"message": "ceshi cccccccccccc"
}
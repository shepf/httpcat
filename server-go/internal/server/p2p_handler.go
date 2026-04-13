package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
	"httpcat/internal/p2p"

	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

func runP2PServer(ctx context.Context, router *gin.Engine) {
	ip := common.P2pListenIP
	port := common.P2pListenPort
	listenAddr := fmt.Sprintf("/ip4/%s/tcp/%d", ip, port)

	fmt.Printf("[*] Listening on: %s with port: %d\n", ip, port)
	fmt.Println("p2p listenAddr:", listenAddr)

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(
			listenAddr,
		),
	)
	if err != nil {
		panic(err)
	}
	defer h.Close()

	fmt.Printf("\033[32mHello World, my p2p hosts ID is %s\033[0m\n", h.ID())

	// 节点发现
	go discoverPeers(ctx, h)

	// PubSub
	ps, _ = pubsub.NewGossipSub(ctx, h)
	ylog.Infof("runP2PServer", "join topic: %v", common.TopicName)
	topic, _ := ps.Join(common.TopicName)
	subscription, _ := topic.Subscribe()
	ylog.Infof("runP2PServer", "subscribed topic: %v", common.TopicName)

	// 将主题添加到已订阅的主题列表
	subscribedTopics[common.TopicName] = topic

	go func() {
		for {
			msg, err := subscription.Next(ctx)
			if err != nil {
				break
			}
			fmt.Printf("Received message from %s: %s\n", msg.GetFrom(), string(msg.GetData()))
		}
	}()

	// 等待上下文取消信号
	<-ctx.Done()
	fmt.Println("P2P server stopped")
}

func publishMessage(c *gin.Context, topicName string, message string) {
	topic, exists := subscribedTopics[topicName]
	if !exists {
		var err error
		topic, err = ps.Join(topicName)
		if err != nil {
			ylog.Errorf("publishMessage", "Failed to join the topic, err:%v", err)
			common.CreateResponse(c, common.ErrorCode, "Failed to join the topic")
			return
		}
		subscribedTopics[topicName] = topic
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := topic.Publish(ctx, []byte(message))
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "Failed to publish the message")
		return
	}

	common.CreateResponse(c, common.SuccessCode, "Message sent to topic: "+topicName)
}

func getSubscribedTopics(c *gin.Context) {
	common.CreateResponse(c, common.SuccessCode, subscribedTopics)
}

func discoverPeers(ctx context.Context, h host.Host) {
	if common.EnableMdns {
		fmt.Printf("Host ID is %s. Enabling MDNS for discovering nodes!\n", h.ID())

		peerChan := p2p.InitMDNS(h, common.RendezvousString)

		connectedPeers := map[peer.ID]bool{}

		for {
			peer := <-peerChan
			if peer.ID == h.ID() {
				continue
			}

			fmt.Println("Found peer:")
			fmt.Println("ID:", peer.ID)

			if _, ok := connectedPeers[peer.ID]; !ok {
				if err := h.Connect(ctx, peer); err != nil {
					fmt.Println("Connection failed:", err)
					continue
				}

				connectedPeers[peer.ID] = true

				fmt.Println("Connected to:", peer.ID)
				fmt.Println("Connected peers:")
				for connectedPeer := range connectedPeers {
					fmt.Println("- ", connectedPeer)
				}
			}
		}
	}
}

// SendP2pMessage P2P 发送消息 handler
func SendP2pMessage(c *gin.Context) {
	type MessageData struct {
		Topic   string `json:"topic"`
		Message string `json:"message"`
	}

	var data MessageData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}

	publishMessage(c, data.Topic, data.Message)
}

// GetSubscribedTopics 获取已订阅的主题列表 handler
func GetSubscribedTopics(c *gin.Context) {
	getSubscribedTopics(c)
}

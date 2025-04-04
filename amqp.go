// Package amqp contains AMQP API for a remote server.
package amqp

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/Azure/go-amqp"
	"go.k6.io/k6/js/modules"
)

type AMQP struct {
	session *amqp.Session
	conn    *amqp.Conn
}

// Options defines configuration options for an AMQP session.
type Options struct {
	ConnectionURL string
	Username      string
	Password      string
}

type PublishOptions struct {
	Session      *amqp.Session
	ConnectionID int
	QueueName    string
	Body         string
}

type ListenOptions struct {
	Session      *amqp.Session
	ConnectionID int
	QueueName    string
}

func (a *AMQP) Start(options Options) *amqp.Session {
	ctx := context.Background()
	connOptions := &amqp.ConnOptions{
		SASLType: amqp.SASLTypePlain(options.Username, options.Password),
	}
	// Dial the connection
	conn, err := amqp.Dial(ctx, options.ConnectionURL, connOptions)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	fmt.Println("Connected to AMQP server")

	// Open a session
	session, err := conn.NewSession(ctx, nil)
	if err != nil {
		conn.Close()
		log.Fatalf("Failed to create session: %v", err)
	}
	a.session = session
	a.conn = conn
	return session
}

// Close cleans up the AMQP connection and session.
func (a *AMQP) Close() {
	if a.session != nil {
		a.session.Close(context.Background())
	}
	if a.conn != nil {
		a.conn.Close()
	}
}
func (a *AMQP) Publish(options PublishOptions) error {
	// Create a sender
	//session, err := amqp.conn.NewSession(options.ConnectionURL)
	ctx := context.Background()
	sender, err := options.Session.NewSender(ctx, options.QueueName, nil)
	if err != nil {
		log.Fatalf("Failed to create sender: %v", err)
		return err
	}
	defer sender.Close(ctx)

	// Send a message
	message := amqp.NewMessage([]byte(options.Body))
	log.Printf("Sending message: %s\n", message.GetData())
	if err := sender.Send(ctx, message, nil); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	fmt.Println("Message sent successfully")
	return nil
}

func (a *AMQP) Listen(options ListenOptions) string {
	// Accept the message
	ctx := context.Background()
	receiver, err := options.Session.NewReceiver(
		ctx, options.QueueName, nil)
	if err != nil {
		log.Fatalf("Failed to create receiver: %v", err)
	}
	defer receiver.Close(ctx)

	// Receive the message
	msg, err := receiver.Receive(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to receive message: %v", err)
	}
	fmt.Printf("Received message : %v\n", string(msg.GetData()))
	message := msg.GetData()
	error := receiver.AcceptMessage(ctx, msg)
	if error != nil {
		log.Fatalf("Failed to accept message: %v", err)
	}
	fmt.Println("Message accepted successfully")
	return string(message)
}

// Register the AMQP module for k6
func init() {
	modules.Register("k6/x/amqp", &AMQP{})
}

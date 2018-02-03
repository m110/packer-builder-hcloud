package hcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type checker func(ctx context.Context) (bool, error)

type Waiter struct {
	client  *hcloud.Client
	timeout time.Duration
}

func NewWaiter(client *hcloud.Client, timeout time.Duration) Waiter {
	return Waiter{
		client:  client,
		timeout: timeout,
	}
}

func (w *Waiter) WaitForServer(ctx context.Context, serverID int, status hcloud.ServerStatus) error {
	check := func(ctx context.Context) (bool, error) {
		log.Printf("Checking server %d status", serverID)

		server, _, err := w.client.Server.GetByID(ctx, serverID)
		if err != nil {
			return false, err
		}

		if server.Status == status {
			return true, nil
		}

		return false, nil
	}

	log.Printf("Waiting for server %d to become %s", serverID, status)
	return w.wait(ctx, check)
}

func (w *Waiter) WaitForImage(ctx context.Context, imageID int, status hcloud.ImageStatus) error {
	check := func(ctx context.Context) (bool, error) {
		log.Printf("Checking image %d status", imageID)

		image, _, err := w.client.Image.GetByID(ctx, imageID)
		if err != nil {
			return false, err
		}

		if image.Status == status {
			return true, nil
		}

		return false, nil
	}

	log.Printf("Waiting for image %d to become %s", imageID, status)
	return w.wait(ctx, check)
}

func (w *Waiter) wait(ctx context.Context, check checker) error {
	done := make(chan struct{})
	defer close(done)

	result := make(chan error, 1)
	go func() {
		for {
			ready, err := check(ctx)
			if err != nil {
				result <- err
				return
			}

			if ready {
				result <- nil
				return
			}

			time.Sleep(2 * time.Second)

			select {
			case <-done:
				return
			default:
			}
		}
	}()

	select {
	case err := <-result:
		return err
	case <-time.After(w.timeout):
		return fmt.Errorf("waiter timed out")
	}
}

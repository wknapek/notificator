package notificator

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

func sendMessage(ctx context.Context, url string, message string) error {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(message))
	if err != nil {
		return err
	}
	req.WithContext(ctx)
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}

package notificator

import (
	"context"
	"errors"
	"fmt"
	"io"
	_ "log"
	"net/http"
	"strings"
	"time"
	_ "time"
)

const NoNotificators = 10

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
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}

func SendMessages(ctx context.Context, url *string, messages chan string) {
	for idx := 0; idx < NoNotificators; idx++ {
		go func() {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			err := sendMessage(ctx, *url, <-messages)
			if err != nil {
				fmt.Print(err.Error())
			}
		}()
	}
}

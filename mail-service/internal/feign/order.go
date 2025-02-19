package feign

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/constant"
	. "github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/dto"
	. "github.com/jordanmarcelino/learn-go-microservices/pkg/dto"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/jwtutils"
)

type OrderClient interface {
	Get(ctx context.Context, req *GetOrderRequest) (*OrderResponse, error)
}

type orderClientImpl struct {
	URL string
}

func NewOrderClient(url string) OrderClient {
	return &orderClientImpl{
		URL: url,
	}
}

func (c *orderClientImpl) Get(ctx context.Context, req *GetOrderRequest) (*OrderResponse, error) {
	url := fmt.Sprintf("%s%s", c.URL, fmt.Sprintf(constant.OrderDetail, req.OrderID))

	token, err := jwtUtil.Sign(&jwtutils.JWTPayload{UserID: req.UserID, Email: req.Email})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	resp, err := feignClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := new(WebResponse[OrderResponse])
	if err := sonic.Unmarshal(bytes, res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

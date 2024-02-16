package persistance

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-redis/redis/v8"
	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
)

type (
	//セッションカート
	SessionCart struct {
		SessionCartProducts []SessionCartProduct `json:"sessionCartProducts"`
	}

	//セッションカート商品
	SessionCartProduct struct {
		ProductID string `json:"productID"`
		Count     int    `json:"count"`
	}

	sessionCartRepository struct {
		redisClient *redis.Client
	}
)

func NewSessionCartRepository(redisClient *redis.Client) sessionCartRepository {
	return sessionCartRepository{
		redisClient: redisClient,
	}
}

func (scr sessionCartRepository) FindBySessionID(ctx context.Context, sessionID string) (entity.SessionCart, bool, error) {
	// Redisからセッションカート情報を取得
	data, err := scr.redisClient.Get(ctx, sessionID).Bytes()
	if err != nil {
		// セッションIDが見つからない場合
		if errors.Is(err, redis.Nil) {
			return entity.SessionCart{}, false, nil
		}

		// その他のエラーの場合
		return entity.SessionCart{}, false, errors.WithStack(err)
	}

	// セッションIDが見つからない場合
	if data == nil {
		return entity.SessionCart{}, false, nil
	}

	// JSONデコードしてセッションカート構造体に変換
	var sessionCart SessionCart
	err = json.Unmarshal(data, &sessionCart)
	if err != nil {
		return entity.SessionCart{}, false, err
	}

	return scr.toEntity(sessionCart, sessionID), true, nil
}

func (src sessionCartRepository) Delete(ctx context.Context, sessionCart entity.SessionCart) error {
	err := src.redisClient.Del(ctx, sessionCart.SessionID).Err()
	return errors.WithStack(err)
}

// セションカートの有効期限を更新する
func (scr sessionCartRepository) UpdateExpiration(ctx context.Context, sessionCart entity.SessionCart, expiration time.Duration) error {
	err := scr.redisClient.Expire(ctx, sessionCart.SessionID, expiration).Err()
	return errors.WithStack(err)
}

func (scr sessionCartRepository) toEntity(sessionCart SessionCart, sessionID string) entity.SessionCart {
	sessionCartProducts := make([]entity.SessionCartProduct, 0, len(sessionCart.SessionCartProducts))
	for _, p := range sessionCart.SessionCartProducts {
		sessionCartProducts = append(sessionCartProducts, entity.SessionCartProduct{
			ProductID: p.ProductID,
			Count:     p.Count,
		})
	}

	return entity.SessionCart{
		SessionID:           sessionID,
		SessionCartProducts: sessionCartProducts,
	}
}

// func (scr sessionCartRepository) toModel(sessionCart entity.SessionCart) SessionCart {
// 	sessionCartProducts := make([]SessionCartProduct, 0, len(sessionCart.SessionCartProducts))
// 	for _, p := range sessionCart.SessionCartProducts {
// 		sessionCartProducts = append(sessionCartProducts, SessionCartProduct{
// 			ProductID: p.ProductID,
// 			Count:     p.Count,
// 		})
// 	}

// 	return SessionCart{
// 		SessionCartProducts: sessionCartProducts,
// 	}
// }

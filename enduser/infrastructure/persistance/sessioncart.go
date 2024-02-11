package persistance

import (
	"context"
	"encoding/json"

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

func (scr sessionCartRepository) FindBySessionID(ctx context.Context, sessionID string) (entity.SessionCart, error) {
	// Redisからセッションカート情報を取得
	data, err := scr.redisClient.Get(ctx, sessionID).Bytes()
	if err != nil {
		// セッションIDが見つからない場合
		if errors.Is(err, redis.Nil) {
			return entity.SessionCart{}, errors.WithStack(ErrSessionNotFound)
		}

		// その他のエラーの場合
		return entity.SessionCart{}, err
	}

	// 有効期限が切れている場合
	if data == nil {
		return entity.SessionCart{}, errors.WithStack(ErrSessionExpired)
	}

	// JSONデコードしてセッションカート構造体に変換
	var sessionCart SessionCart
	err = json.Unmarshal(data, &sessionCart)
	if err != nil {
		return entity.SessionCart{}, err
	}

	return scr.toEntity(sessionCart, sessionID), nil
}

func (src sessionCartRepository) Delete(ctx context.Context, sessionCart entity.SessionCart) error {
	err := src.redisClient.Del(ctx, sessionCart.SessionID).Err()
	return errors.WithStack(err)
}

func (scr sessionCartRepository) toEntity(sessionCart SessionCart, sessionID string) entity.SessionCart {
	sessionCartProducts := make([]entity.SessionCartProduct, len(sessionCart.SessionCartProducts))
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
// 	sessionCartProducts := make([]SessionCartProduct, len(sessionCart.SessionCartProducts))
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

package persistance

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-redis/redis/v8"
	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/share"
)

type (
	sessionAccountRepository struct {
		redisClient *redis.Client
	}
)

func NewSessionAccountRepository(redisClient *redis.Client) sessionAccountRepository {
	return sessionAccountRepository{
		redisClient: redisClient,
	}
}

func (sar sessionAccountRepository) Insert(ctx context.Context, sessionAccount *entity.SessionAccount, expiration time.Duration, eventPublisher share.DomainEventPublisher) error {
	err := sar.redisClient.Set(ctx, sessionAccount.SessionID, sessionAccount.AccountID, expiration).Err()
	if err != nil {
		return errors.WithStack(err)
	}

	if eventPublisher == nil {
		return nil
	}

	events := sessionAccount.ClearEvents()
	err = eventPublisher.Publish(events)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (sar sessionAccountRepository) FindBySessionID(ctx context.Context, sessionID string) (entity.SessionAccount, bool, error) {
	accountID, err := sar.redisClient.Get(ctx, sessionID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// セッションIDが見つからない場合
			return entity.SessionAccount{}, false, nil
		}

		// その他のエラーの場合
		return entity.SessionAccount{}, false, errors.WithStack(err)
	}

	// セッションIDが見つからない場合
	if accountID == "" {
		return entity.SessionAccount{}, false, nil
	}

	return entity.SessionAccount{
		SessionID: sessionID,
		AccountID: accountID,
	}, true, nil
}

// セションアカウントの有効期限を更新する
func (sar sessionAccountRepository) UpdateExpiration(ctx context.Context, sessionAccount entity.SessionAccount, expiration time.Duration) error {
	err := sar.redisClient.Expire(ctx, sessionAccount.SessionID, expiration).Err()
	return errors.WithStack(err)
}

// セッションアカウントを削除する
func (sar sessionAccountRepository) Delete(ctx context.Context, sessionAccount entity.SessionAccount) error {
	err := sar.redisClient.Del(ctx, sessionAccount.SessionID).Err()
	return errors.WithStack(err)
}

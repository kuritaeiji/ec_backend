package persistance_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kuritaeiji/ec_backend/config"
	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/enduser/infrastructure/persistance"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/share/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type sessionAccountRepositoryTestSuite struct {
	suite.Suite
	sessionAccountRepository repository.SessionAccountRepository
	redisClient              *redis.Client
}

func TestSessionAccountRepository(t *testing.T) {
	err := config.SetupEnv()
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("環境変数設定時にエラーが発生しました。\n+%+v", err))
	}
	redisClient := config.NewRedisClient()
	suite.Run(t, &sessionAccountRepositoryTestSuite{
		sessionAccountRepository: persistance.NewSessionAccountRepository(redisClient),
		redisClient:              redisClient,
	})
}

func (suite *sessionAccountRepositoryTestSuite) tearDown() {
	err := suite.redisClient.FlushAll(context.Background()).Err()
	if err != nil {
		suite.FailNow("Redisのデータ全削除時にエラー発生\n+%v", err)
	}
}

func (suite *sessionAccountRepositoryTestSuite) TestInsert() {
	defer suite.tearDown()

	// given（前提条件）
	event := entity.SessionAccountCreatedEvent{}
	sessionAccount := entity.SessionAccount{
		AccountID: "accountID",
		SessionID: "sessionID",
		Events:    []share.DomainEvent{event},
	}
	expiration := 1 * time.Hour
	eventPublisherMock := mocks.NewDomainEventPublisher(suite.T())
	eventPublisherMock.On("Publish", []share.DomainEvent{event}).Return(nil)

	// when（操作）
	err := suite.sessionAccountRepository.Insert(context.Background(), sessionAccount, expiration, eventPublisherMock)

	// then（期待する結果）
	suite.Nil(err)

	sessionAccountResult, ok, err := suite.sessionAccountRepository.FindBySessionID(context.Background(), sessionAccount.SessionID)
	if err != nil {
		suite.FailNow("アカウントID取得時にエラー発生\n+%+v", err)
	}
	suite.True(ok)
	suite.Equal(sessionAccount.AccountID, sessionAccountResult.AccountID)

	ttl, err := suite.redisClient.TTL(context.Background(), sessionAccount.SessionID).Result()
	if err != nil {
		suite.FailNow("セッションIDの有効期限取得時にエラー発生\n+%+v", err)
	}
	suite.Equal(expiration, ttl)
}

func (suite *sessionAccountRepositoryTestSuite) TestFindBySessionID() {
	sessionID := "sessionID"
	accountID := "accountID"

	type expected struct {
		SessionAccount       entity.SessionAccount
		ExistsSessionAccount bool
		Error                error
	}

	tests := []struct {
		Name     string
		Setup    func(t *testing.T)
		Expected expected
	}{
		{
			Name: "セッションアカウントが存在する場合",
			Setup: func(t *testing.T) {
				err := suite.redisClient.Set(context.Background(), sessionID, accountID, 1*time.Hour).Err()
				if err != nil {
					assert.FailNow(t, fmt.Sprintf("セッションアカウントを登録できませんでした\n%+v", err))
				}
			},
			Expected: expected{SessionAccount: entity.SessionAccount{SessionID: sessionID, AccountID: accountID}, ExistsSessionAccount: true, Error: nil},
		},
		{
			Name: "セッションIDが存在しない場合",
			Setup: func(t *testing.T) {
				err := suite.redisClient.Set(context.Background(), "invalidID", accountID, 1*time.Hour).Err()
				if err != nil {
					assert.FailNow(t, fmt.Sprintf("セッションアカウントを登録できませんでした\n%+v", err))
				}
			},
			Expected: expected{SessionAccount: entity.SessionAccount{}, ExistsSessionAccount: false, Error: nil},
		},
		{
			Name: "有効期限切れの場合",
			Setup: func(t *testing.T) {
				err := suite.redisClient.Set(context.Background(), sessionID, accountID, 1*time.Nanosecond).Err()
				if err != nil {
					assert.FailNow(t, fmt.Sprintf("セッションアカウントを登録できませんでした\n%+v", err))
				}
				time.Sleep(1 * time.Second)
			},
			Expected: expected{SessionAccount: entity.SessionAccount{}, ExistsSessionAccount: false, Error: nil},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.Name, func(t *testing.T) {
			defer suite.tearDown()

			tt.Setup(t)

			sessionAccount, ok, err := suite.sessionAccountRepository.FindBySessionID(context.Background(), sessionID)
			suite.Equal(tt.Expected.SessionAccount, sessionAccount)
			suite.Equal(tt.Expected.ExistsSessionAccount, ok)
			suite.Equal(tt.Expected.Error, err)
		})
	}
}

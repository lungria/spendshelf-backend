package sync_mono

import (
	"context"
	"log"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/lungria/spendshelf-backend/src/webhooks"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lungria/spendshelf-backend/src/models"

	shalmono "github.com/shal/mono"

	"github.com/lungria/spendshelf-backend/src/transactions"
)

type MonoSync struct {
	sync.RWMutex
	txnRepo      transactions.Repository
	monoClient   *shalmono.Personal
	accountUAH   *shalmono.Account
	logger       *zap.SugaredLogger
	transactions chan []shalmono.Transaction
	errChan      chan error
}

func NewMonoSync(token string, txnRepo transactions.Repository, logger *zap.SugaredLogger) (*MonoSync, error) {
	s := MonoSync{
		monoClient:   shalmono.NewPersonal(token),
		transactions: make(chan []shalmono.Transaction),
		errChan:      make(chan error, 1),
		txnRepo:      txnRepo,
		logger:       logger,
	}

	accUAH, err := getAccount(*s.monoClient)
	if err != nil {
		return nil, err
	}
	s.accountUAH = accUAH

	return &s, nil
}

func (s *MonoSync) Transactions(createdAtAccount time.Time) {
	ctx := context.Background()
	defer ctx.Done()

	from := createdAtAccount
	for from.Before(time.Now().UTC()) {
		to := from.Add(time.Hour * 744)
		log.Println("Start", from.String())
		log.Println("End", to.String())
		txns, err := s.monoClient.Transactions(ctx, s.accountUAH.ID, from, to)
		if err != nil {
			s.logger.Errorw("Unable to fetch transactions from mono bank", "Error", err.Error())
			s.errChan <- err
		}
		go func() {
			s.transactions <- txns
		}()
		from = to

		time.Sleep(time.Second * 61)
	}
}

func getAccount(monoPersonal shalmono.Personal) (*shalmono.Account, error) {
	ctx := context.Background()
	defer ctx.Done()
	user, err := monoPersonal.User(ctx)
	if err != nil {
		return nil, err
	}

	var account shalmono.Account

	for _, acc := range user.Accounts {
		ccy, _ := shalmono.CurrencyFromISO4217(acc.CurrencyCode)
		if ccy.Code == "UAH" {
			account = acc
			break
		}
	}
	return &account, nil
}

func (s *MonoSync) trimDuplicate(syncTxns []shalmono.Transaction) []models.Transaction {
	s.RLock()
	defer s.RUnlock()
	unique := []models.Transaction{}

	currentTxns, err := s.txnRepo.FindAll()
	if err != nil {
		s.logger.Error("Unable to find transactions from transactions collection", "Error", err.Error())
		s.errChan <- err
	}
	curr := make(map[string]models.Transaction, len(currentTxns))

	for i := 0; i < len(currentTxns); i++ {
		currID := currentTxns[i].BankTransaction.StatementItem.ID
		curr[currID] = currentTxns[i]
	}

	for i := 0; i < len(syncTxns); i++ {
		_, found := curr[syncTxns[i].ID]
		if !found {
			unique = append(unique, s.txnFromSyncTxn(syncTxns[i]))
		}
	}
	return unique
}

func (s *MonoSync) txnFromSyncTxn(syncTxn shalmono.Transaction) models.Transaction {
	var txn models.Transaction

	txn.ID = primitive.NewObjectID()
	txn.BankTransaction.AccountID = s.accountUAH.ID
	txn.BankTransaction.StatementItem = syncTxn
	txn.Bank = webhooks.MonoBankName
	txn.Time = time.Unix(int64(syncTxn.Time), 0)
	txn.Description = syncTxn.Description
	txn.Amount = syncTxn.Amount
	txn.Balance = syncTxn.Balance

	return txn
}

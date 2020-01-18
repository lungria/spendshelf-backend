package sync_mono

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/lungria/spendshelf-backend/src/webhooks"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lungria/spendshelf-backend/src/models"

	shalmono "github.com/shal/mono"

	"github.com/lungria/spendshelf-backend/src/transactions"
)

type Sync struct {
	sync.RWMutex
	txnRepo      transactions.Repository
	client       *shalmono.Personal
	accountUAH   *shalmono.Account
	transactions chan []shalmono.Transaction
	errChan      chan error
}

func NewSync(token string, txnRepo transactions.Repository) (*Sync, error) {
	s := Sync{
		client:       shalmono.NewPersonal(token),
		transactions: make(chan []shalmono.Transaction),
		errChan:      make(chan error, 1),
		txnRepo:      txnRepo,
	}

	accUAH, err := getAccount(*s.client)
	if err != nil {
		return nil, err
	}
	s.accountUAH = accUAH

	go s.run()

	return &s, nil
}

func (s *Sync) Transactions(createdAtAccount time.Time) {
	ctx := context.Background()
	defer ctx.Done()

	from := createdAtAccount
	for from.Before(time.Now().UTC()) {
		to := from.Add(time.Hour * 744)
		log.Println("Start", from.String())
		log.Println("End", to.String())
		txns, err := s.client.Transactions(ctx, s.accountUAH.ID, from, to)
		if err != nil {
			log.Println(err)
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

func (s *Sync) run() {
	for {
		select {
		case err := <-s.errChan:
			log.Println(err)
			return
		case txns := <-s.transactions:
			toInsert := s.trimDuplicate(txns)
			s.Lock()
			err := s.txnRepo.InsertManyTransactions(toInsert)
			s.Unlock()
			s.errChan <- err
		}
	}
}

func (s *Sync) trimDuplicate(syncTxns []shalmono.Transaction) []models.Transaction {
	s.RLock()
	defer s.RUnlock()
	unique := []models.Transaction{}

	currentTxns, err := s.txnRepo.FindAll()
	if err != nil {
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

func (s *Sync) txnFromSyncTxn(syncTxn shalmono.Transaction) models.Transaction {
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

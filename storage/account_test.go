package storage_test

//
//func TestAccountStorageSave_WithLocalDb_NoErrorReturned(t *testing.T) {
//	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
//		os.Exit(1)
//	}
//	defer dbpool.Close()
//
//	db := storage.NewAccountsStorage(dbpool)
//
//	// test insert
//	err = db.Save(context.Background(), storage.Account{
//		ID:       "acc1",
//		Balance:  10000,
//		Currency: "UAH",
//	})
//
//	assert.NoError(t, err)
//
//	// test on conflict update
//	err = db.Save(context.Background(), storage.Account{
//		ID:       "acc1",
//		Balance:  20000,
//		Currency: "UAH",
//	})
//
//	assert.NoError(t, err)
//}

package mono_test

//
//func TestClient_GetTransactions_NoErrorsReturned(t *testing.T) {
//	key := os.Getenv("MONO_API_KEY")
//	if key == "" {
//		t.Fatal("MONO_API_KEY is required for test")
//		return
//	}
//
//	client := mono.NewClient("https://api.monobank.ua", key)
//
//	accID := os.Getenv("MONO_ACCOUNT_ID")
//	if accID == "" {
//		t.Fatal("MONO_ACCOUNT_ID is required for test")
//		return
//	}
//
//	_, err := client.GetTransactions(context.Background(), mono.GetTransactionsQuery{
//		Account: accID,
//		From:    time.Now().Add(-48 * time.Hour),
//		To:      time.Now(),
//	})
//
//	assert.NoError(t, err)
//}
//
//func TestJsonUnmarshal_TimeParsedCorrectly(t *testing.T) {
//	jsonData := []byte(`[
//   {
//      "id":"tr-id",
//      "time":1601022957,
//      "description":"Era-in-ear",
//      "mcc":5732,
//      "amount":-402100,
//      "operationAmount":-402100,
//      "currencyCode":980,
//      "commissionRate":0,
//      "cashbackAmount":0,
//      "balance":38817,
//      "hold":true
//   }]`)
//	tr := make([]mono.Transaction, 0)
//
//	err := json.Unmarshal(jsonData, &tr)
//
//	assert.NoError(t, err)
//}

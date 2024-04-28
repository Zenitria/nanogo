package nanogo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"
)

// Client is a client for the Nano RPC protocol,
// Url: the url of the RPC server,
// AuthHeader: the authentication header of the RPC server (optional).
// AuthToken: the authorization token of the RPC server (optional),
type Client struct {
	Url        string
	AuthHeader string // optional
	AuthToken  string // optional
}

// AccountInfo is the account info of a wallet,
// Frontier: the frontier of the wallet,
// ConfirmedFrontier: the confirmed frontier of the wallet,
// OpenBlock: the open block of the wallet,
// RepresentativeBlock: the representative block of the wallet,
// Balance: the balance of the wallet in raw,
// ConfirmedBalance: the confirmed balance of the wallet in raw,
// Representative: the representative of the wallet,
// ConfirmedRepresentative: the confirmed representative of the wallet,
// ModifiedTimestamp: the modified timestamp of the wallet,
// BlockCount: the block count of the wallet,
// ConfirmedHeight: the confirmed height of the wallet,
// AccountVersion: the account version of the wallet,
// ConfirmationHeight: the confirmation height of the wallet,
// ConfirmationHeightFrontier: the confirmation height frontier of the wallet,
// Error: the error of the request.
type AccountInfo struct {
	Frontier                   string `json:"frontier"`
	ConfirmedFrontier          string `json:"confirmed_frontier"`
	OpenBlock                  string `json:"open_block"`
	RepresentativeBlock        string `json:"representative_block"`
	Balance                    string `json:"balance"`
	ConfirmedBalance           string `json:"confirmed_balance"`
	Representative             string `json:"representative"`
	ConfirmedRepresentative    string `json:"confirmed_representative"`
	ModifiedTimestamp          string `json:"modified_timestamp"`
	BlockCount                 string `json:"block_count"`
	ConfirmedHeight            string `json:"confirmed_height"`
	AccountVersion             string `json:"account_version"`
	ConfirmationHeight         string `json:"confirmation_height"`
	ConfirmationHeightFrontier string `json:"confirmation_height_frontier"`

	Error any `json:"error"`
}

// AccountBalance is the balance of a wallet,
// Balance: the balance of the wallet in raw,
// Pending: the pending balance of the wallet in raw,
// Receivable: the receivable balance of the wallet in raw,
// Error: the error of the request.
type AccountBalance struct {
	Balance    string `json:"balance"`
	Pending    string `json:"pending"`
	Receivable string `json:"receivable"`

	Error any `json:"error"`
}

// AccountHistory is the history of a wallet,
// Account: the wallet address,
// History: the history of the wallet,
// Previous: the previous block of the wallet,
// Error: the error of the request.
type AccountHistory struct {
	Account string `json:"account"`
	History []struct {
		Type           string `json:"type"`
		Account        string `json:"account"`
		Amount         string `json:"amount"`
		LocalTimestamp string `json:"local_timestamp"`
		Hash           string `json:"hash"`
		Confirmed      string `json:"confirmed"`
	}
	Previous string `json:"previous"`

	Error any `json:"error"`
}

// Receivable is the receivable blocks of a wallet,
// Blocks: the receivable blocks of the wallet,
// Error: the error of the request.
type Receivable struct {
	Blocks map[string]struct {
		Amount string `json:"amount"`
		Source string `json:"source"`
	} `json:"blocks"`

	Error any `json:"error"`
}

// Representatives is the online representatives of the network,
// Representatives: list of the online representatives of the network,
// Error: the error of the request.
type Representatives struct {
	Representatives []string `json:"representatives"`

	Error any `json:"error"`
}

// RPC sends a JSON-RPC request to the RPC server,
// body: the body of the request,
// returns the response or an error.
func (c *Client) RPC(data map[string]any) ([]byte, error) {
	dataJson, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.Url, bytes.NewBuffer(dataJson))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.AuthHeader != "" && c.AuthToken != "" {
		req.Header.Set(c.AuthHeader, c.AuthToken)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

// GetAccountBalance gets the balance of a wallet,
// address: the wallet address to get the balance of,
// returns the balance in raw or an error.
func (c *Client) GetAccountBalance(address string) (AccountBalance, error) {
	data := map[string]any{
		"action":  "account_balance",
		"account": address,
	}

	res, err := c.RPC(data)

	if err != nil {
		return AccountBalance{}, err
	}

	var body AccountBalance
	json.Unmarshal(res, &body)

	if body.Error != nil {
		return AccountBalance{}, fmt.Errorf("%v", body.Error)
	}

	return body, nil
}

// GetAccountInfo gets the account info of a wallet,
// address: the wallet address to get the info of,
// returns the account info or an error.
func (c *Client) GetAccountInfo(address string) (AccountInfo, error) {
	data := map[string]any{
		"action":            "account_info",
		"account":           address,
		"include_confirmed": "true",
		"representative":    "true",
	}

	res, err := c.RPC(data)

	if err != nil {
		return AccountInfo{}, err
	}

	var info AccountInfo
	json.Unmarshal(res, &info)

	if info.Error == "Account not found" {
		return AccountInfo{}, ErrAccountNotFound
	}

	if info.Error != nil {
		return AccountInfo{}, fmt.Errorf("%v", info.Error)
	}

	return info, nil
}

// GetAccountHistory gets the history of a wallet,
// address: the wallet address to get the history of,
// count: the count of the history to get (-1 for all),
// returns the account history or an error.
func (c *Client) GetAccountHistory(address string, count int) (AccountHistory, error) {
	data := map[string]any{
		"action":  "account_history",
		"account": address,
		"count":   count,
	}

	res, err := c.RPC(data)

	if err != nil {
		return AccountHistory{}, err
	}

	var history AccountHistory
	json.Unmarshal(res, &history)

	if history.Error != nil {
		return AccountHistory{}, fmt.Errorf("%v", history.Error)
	}

	return history, nil
}

// GetReceivable gets the receivable blocks of a wallet,
// address: the wallet address to get the receivable blocks of,
// returns the receivable blocks or an error.
func (c *Client) GetReceivable(address string) (Receivable, error) {
	data := map[string]any{
		"action":  "receivable",
		"account": address,
		"source":  "true",
	}

	res, err := c.RPC(data)

	if err != nil {
		return Receivable{}, nil
	}

	var receivable Receivable
	json.Unmarshal(res, &receivable)

	if receivable.Error != nil {
		return Receivable{}, fmt.Errorf("%v", receivable.Error)
	}

	return receivable, nil
}

// GetRepresentatives gets the online representatives of the network,
// returns the representatives or an error.
func (c *Client) GetRepresentatives() (Representatives, error) {
	data := map[string]any{
		"action": "representatives_online",
	}

	res, err := c.RPC(data)

	if err != nil {
		return Representatives{}, err
	}

	var reps Representatives
	json.Unmarshal(res, &reps)

	if reps.Error != nil {
		return Representatives{}, fmt.Errorf("%v", reps.Error)
	}

	return reps, nil
}

// Process processes a block,
// subtype: the subtype of the block,
// block: the block to process,
// returns the block hash or an error.
func (c *Client) Process(subtype string, block Block) (string, error) {
	data := map[string]any{
		"action":     "process",
		"subtype":    subtype,
		"json_block": "true",
		"block":      block,
	}

	res, err := c.RPC(data)

	if err != nil {
		return "", err
	}

	var body struct {
		Hash string `json:"hash"`

		Error any `json:"error"`
	}
	json.Unmarshal(res, &body)

	if body.Error != nil {
		return "", fmt.Errorf("%v", body.Error)
	}

	return body.Hash, nil
}

// GenerateWork generates work for a block using RPC,
// hash: the block hash to generate work for,
// returns the work or an error.
func (c *Client) GenerateWork(block Block) (string, error) {
	var hash string

	if block.Previous == "0" || block.Previous == "0000000000000000000000000000000000000000000000000000000000000000" {
		pubKey, err := AddressToPublicKey(block.Account)

		if err != nil {
			return "", err
		}

		hash = fmt.Sprintf("%064X", pubKey)
	} else {
		hash = block.Previous
	}

	data := map[string]any{
		"action": "work_generate",
		"hash":   hash,
	}

	res, err := c.RPC(data)

	if err != nil {
		return "", err
	}

	var body struct {
		Work string `json:"work"`

		Error any `json:"error"`
	}
	json.Unmarshal(res, &body)

	if body.Error != nil {
		return "", fmt.Errorf("%v", body.Error)
	}

	return body.Work, nil
}

// Send sends a raw amount of Nano to a wallet,
// toAddr: the destination wallet address,
// raw: the amount to send in raw,
// seed: the seed of the sending wallet,
// index: the index of the sending wallet (usually 0),
// returns the block hash or an error.
func (c *Client) Send(toAddress, raw, seed string, index int) (string, error) {
	privKey, err := SeedToPrivateKey(seed, index)

	if err != nil {
		return "", err
	}

	pubKey, err := PrivateKeyToPublicKey(privKey)

	if err != nil {
		return "", err
	}

	addr, err := PublicKeyToAddress(pubKey)

	if err != nil {
		return "", err
	}

	info, err := c.GetAccountInfo(addr)

	if err != nil {
		return "", err
	}

	bal, ok := new(big.Int).SetString(info.ConfirmedBalance, 10)

	if !ok {
		return "", fmt.Errorf("could not convert string to big int")
	}

	rawBigInt, ok := new(big.Int).SetString(raw, 10)

	if !ok {
		return "", fmt.Errorf("could not convert string to big int")
	}

	if bal.Cmp(rawBigInt) < 0 {
		return "", fmt.Errorf("raw is bigger than wallet balance")
	}

	balAfter := new(big.Int).Sub(bal, rawBigInt)
	rcptPubKey, err := AddressToPublicKey(toAddress)

	if err != nil {
		return "", err
	}

	block := Block{
		Type:           "state",
		Account:        addr,
		Previous:       info.Frontier,
		Representative: info.Representative,
		Balance:        balAfter.String(),
		Link:           fmt.Sprintf("%064X", rcptPubKey),
		LinkAsAccount:  toAddress,
	}

	err = block.Sign(privKey)

	if err != nil {
		return "", err
	}

	work, err := c.GenerateWork(block)

	if err != nil {
		return "", err
	}

	block.AddWork(work)

	return c.Process("send", block)
}

// ChangeRepresentative changes the representative of a wallet,
// representative: the new representative wallet address,
// seed: the seed of the wallet,
// index: the index of the wallet (usually 0),
// returns the block hash or an error.
func (c *Client) ChangeRepresentative(representative, seed string, index int) (string, error) {
	privKey, err := SeedToPrivateKey(seed, index)

	if err != nil {
		return "", err
	}

	pubKey, err := PrivateKeyToPublicKey(privKey)

	if err != nil {
		return "", err
	}

	addr, err := PublicKeyToAddress(pubKey)

	if err != nil {
		return "", err
	}

	info, err := c.GetAccountInfo(addr)

	if err != nil {
		return "", err
	}

	block := Block{
		Type:           "state",
		Account:        addr,
		Previous:       info.Frontier,
		Representative: representative,
		Balance:        info.ConfirmedBalance,
		Link:           "0000000000000000000000000000000000000000000000000000000000000000",
		LinkAsAccount:  "nano_1111111111111111111111111111111111111111111111111111hifc8npp",
	}

	err = block.Sign(privKey)

	if err != nil {
		return "", err
	}

	work, err := c.GenerateWork(block)

	if err != nil {
		return "", err
	}

	block.AddWork(work)

	return c.Process("change", block)
}

// Receive receives a block,
// hash: the block hash to receive,
// sourceAddress: the source wallet address,
// raw: the amount to receive in raw,
// representative: the representative wallet address,
// seed: the seed of the receiving wallet,
// index: the index of the receiving wallet (usually 0),
func (c *Client) Receive(hash, sourceAddress, raw, seed string, index int) (string, error) {
	privKey, err := SeedToPrivateKey(seed, index)

	if err != nil {
		return "", err
	}

	pubKey, err := PrivateKeyToPublicKey(privKey)

	if err != nil {
		return "", err
	}

	addr, err := PublicKeyToAddress(pubKey)

	if err != nil {
		return "", err
	}

	info, err := c.GetAccountInfo(addr)

	if errors.Is(err, ErrAccountNotFound) {
		reps, err := c.GetRepresentatives()

		if err != nil {
			return "", err
		}

		info.ConfirmedBalance = "0"
		info.Frontier = "0000000000000000000000000000000000000000000000000000000000000000"
		info.Representative = reps.Representatives[0]

	} else if err != nil {
		return "", err
	}

	bal, ok := new(big.Int).SetString(info.ConfirmedBalance, 10)

	if !ok {
		return "", fmt.Errorf("could not convert string to big int")
	}

	rawBigInt, ok := new(big.Int).SetString(raw, 10)

	if !ok {
		return "", fmt.Errorf("could not convert string to big int")
	}

	balAfter := new(big.Int).Add(bal, rawBigInt)

	block := Block{
		Type:           "state",
		Account:        addr,
		Previous:       info.Frontier,
		Representative: info.Representative,
		Balance:        balAfter.String(),
		Link:           hash,
		LinkAsAccount:  sourceAddress,
	}

	err = block.Sign(privKey)

	if err != nil {
		return "", err
	}

	work, err := c.GenerateWork(block)

	if err != nil {
		return "", err
	}

	block.AddWork(work)

	return c.Process("receive", block)
}

// ReceiveAll receives all receivable blocks of a wallet,
// representative: the representative wallet address,
// seed: the seed of the receiving wallet,
// index: the index of the receiving wallet (usually 0),
// returns the block hashes or an error.
func (c *Client) ReceiveAll(seed string, index int) ([]string, error) {
	privKey, err := SeedToPrivateKey(seed, index)

	if err != nil {
		return []string{}, err
	}

	pubKey, err := PrivateKeyToPublicKey(privKey)

	if err != nil {
		return []string{}, err
	}

	addr, err := PublicKeyToAddress(pubKey)

	if err != nil {
		return []string{}, err
	}

	receivable, err := c.GetReceivable(addr)

	if err != nil {
		return []string{}, err
	}

	var hashes []string

	for h, b := range receivable.Blocks {
		hash, err := c.Receive(h, b.Source, b.Amount, seed, index)

		if err != nil {
			return []string{}, err
		}

		hashes = append(hashes, hash)

		time.Sleep(2 * time.Second)
	}

	return hashes, nil
}

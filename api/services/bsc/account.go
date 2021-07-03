package bsc

import (
	"fmt"
	"net/url"
	"strings"
)

type BalanceResponse string

type MultiBalanceResult struct {
	Account string `json:"account"`
	Balance string `json:"balance"`
}

type TransactionResponse struct {
	BlockNumber        string `json:"blockNumber"`
	TimeStamp          string `json:"timeStamp"`
	Hash               string `json:"hash"`
	Nonce              string `json:"nonce"`
	BlockHash          string `json:"blockHash"`
	TransactionIndex   string `json:"transactionIndex"`
	From               string `json:"from"`
	To                 string `json:"to"`
	Value              string `json:"value"`
	Gas                string `json:"gas"`
	GasPrice           string `json:"gasPrice"`
	IsError            string `json:"isError"`
	TransactionReceipt string `json:"txreceipt_status"`
	Input              string `json:"input"`
	ContractAddress    string `json:"contractAddress"`
	CumulativeGasUsed  string `json:"cumulativeGasUsed"`
	GasUsed            string `json:"gasUsed"`
	Confirmations      string `json:"confirmations"`
}

type TransactionInternalResponse struct {
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	ContractAddress string `json:"contractAddress"`
	Input           string `json:"input"`
	Type            string `json:"type"`
	Gas             string `json:"gas"`
	GasUsed         string `json:"gasUsed"`
	TraceId         string `json:"traceId"`
	IsError         string `json:"isError"`
	ErrorCode       string `json:"errCode"`
}

type TokenTransactionResponse struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"`
	To                string `json:"to"`
	Value             string `json:"value"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      string `json:"tokenDecimal"`
	TransactionIndex  string `json:"transactionIndex"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Input             string `json:"input"`
	Confirmations     string `json:"confirmations"`
}

type BalanceOptions struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type TransactionsOptions struct {
	StartBlock int32  `json:"startBlock"`
	EndBlock   int32  `json:"endBlock"`
	Sort       string `json:"sort"`
}

func (c *Client) GetBalance(address string) (BalanceResponse, error) {
	queryValues := url.Values{}
	queryValues.Add("module", "account")
	queryValues.Add("action", "balance")
	queryValues.Add("address", address)

	var res BalanceResponse
	if err := c.sendRequest(queryValues, &res); err != nil {
		return "", err
	}

	return res, nil
}

func (c *Client) GetMultiBalance(addresses []string) (*[]MultiBalanceResult, error) {
	queryValues := url.Values{}
	queryValues.Add("module", "account")
	queryValues.Add("action", "balancemulti")
	queryValues.Add("address", strings.Join(addresses, ","))

	var res []MultiBalanceResult
	if err := c.sendRequest(queryValues, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetTransactionList(address string, options *TransactionsOptions) (*[]TransactionResponse, error) {
	startblock := "1"
	endblock := "99999999"
	sort := "asc"
	if options != nil {
		startblock = fmt.Sprintf("%d", options.StartBlock)
		endblock = fmt.Sprintf("%d", options.EndBlock)
		sort = options.Sort
	}

	queryValues := url.Values{}
	queryValues.Add("module", "account")
	queryValues.Add("action", "txlist")
	queryValues.Add("address", address)
	queryValues.Add("startblock", startblock)
	queryValues.Add("endblock", endblock)
	queryValues.Add("sort", sort)

	var res []TransactionResponse
	if err := c.sendRequest(queryValues, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetInternalTransactionList(address string, options *TransactionsOptions) (*[]TransactionInternalResponse, error) {
	startblock := "1"
	endblock := "99999999"
	sort := "asc"
	if options != nil {
		startblock = fmt.Sprintf("%d", options.StartBlock)
		endblock = fmt.Sprintf("%d", options.EndBlock)
		sort = options.Sort
	}

	queryValues := url.Values{}
	queryValues.Add("module", "account")
	queryValues.Add("action", "txlistinternal")
	queryValues.Add("address", address)
	queryValues.Add("startblock", startblock)
	queryValues.Add("endblock", endblock)
	queryValues.Add("sort", sort)

	var res []TransactionInternalResponse
	if err := c.sendRequest(queryValues, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetInternalTransactionsTxList(txHash string) (*[]TransactionInternalResponse, error) {
	queryValues := url.Values{}
	queryValues.Add("module", "account")
	queryValues.Add("action", "txlistinternal")
	queryValues.Add("txhash", txHash)

	var res []TransactionInternalResponse
	if err := c.sendRequest(queryValues, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetTokenTransfers(address string, options *TransactionsOptions) (*[]TokenTransactionResponse, error) {
	startblock := "1"
	endblock := "99999999"
	sort := "asc"
	if options != nil {
		startblock = fmt.Sprintf("%d", options.StartBlock)
		endblock = fmt.Sprintf("%d", options.EndBlock)
		sort = options.Sort
	}

	queryValues := url.Values{}
	queryValues.Add("module", "account")
	queryValues.Add("action", "tokentx")
	queryValues.Add("address", address)
	queryValues.Add("startblock", startblock)
	queryValues.Add("endblock", endblock)
	queryValues.Add("sort", sort)

	var res []TokenTransactionResponse
	if err := c.sendRequest(queryValues, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

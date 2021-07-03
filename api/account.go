package bsc

import (
	"fmt"
	"github.com/antonioplacerda/go_bsc/api/utils"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Balance struct {
	Account string
	Balance float64
}

type TransactionsOptions struct {
	StartBlock int32  `json:"startBlock"`
	EndBlock   int32  `json:"endBlock"`
	Sort       string `json:"sort"`
}

type Transaction struct {
	TimeStamp       time.Time `json:"timeStamp"`
	Hash            string    `json:"hash"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	Value           float64   `json:"value"`
	IsError         bool      `json:"isError"`
	ContractAddress string    `json:"contractAddress"`
	Fee             float64   `json:"fee"`
}

type TokenTransfer struct {
	TimeStamp       time.Time `json:"timeStamp"`
	Hash            string    `json:"hash"`
	From            string    `json:"from"`
	ContractAddress string    `json:"contractAddress"`
	To              string    `json:"to"`
	Value           float64   `json:"value"`
	TokenName       string    `json:"tokenName"`
	TokenSymbol     string    `json:"tokenSymbol"`
	TokenDecimal    string    `json:"tokenDecimal"`
	Fee             float64   `json:"fee"`
}

func (c *Client) GetBalance(address string) (float64, error) {
	queryValues := url.Values{}
	queryValues.Add("module", "account")
	queryValues.Add("action", "balance")
	queryValues.Add("address", address)

	var res string
	if err := c.sendRequest(queryValues, &res); err != nil {
		return 0, err
	}

	balance, err := utils.ConvertStringToBNB(res)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (c *Client) GetMultiBalance(addresses []string) (*[]Balance, error) {
	queryValues := url.Values{}
	queryValues.Add("module", "account")
	queryValues.Add("action", "balancemulti")
	queryValues.Add("address", strings.Join(addresses, ","))

	type multiBalanceResult struct {
		Account string `json:"account"`
		Balance string `json:"balance"`
	}

	var res []multiBalanceResult
	if err := c.sendRequest(queryValues, &res); err != nil {
		return nil, err
	}

	var balances []Balance
	for i := 0; i < len(res); i++ {
		balance, err := utils.ConvertStringToBNB(res[i].Balance)
		if err != nil {
			return nil, err
		}
		balances = append(balances, Balance{
			res[i].Account,
			balance,
		})
	}

	return &balances, nil
}

func (c *Client) GetTransactionList(address string, options *TransactionsOptions) (*[]Transaction, error) {
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

	type transactionResult struct {
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

	var res []transactionResult
	if err := c.sendRequest(queryValues, &res); err != nil {
		return nil, err
	}

	var transactions []Transaction
	for i := 0; i < len(res); i++ {
		timeInt, err := strconv.ParseInt(res[i].TimeStamp, 10, 64)
		if err != nil {
			panic(err)
		}
		timestamp := time.Unix(timeInt, 0)
		amount, err := utils.ConvertStringToBNB(res[i].Value)
		if err != nil {
			return nil, err
		}
		fee, err := utils.ComputeFee(res[i].GasUsed, res[i].GasPrice)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, Transaction{
			timestamp,
			res[i].Hash,
			res[i].From,
			res[i].To,
			amount,
			res[i].IsError != "0",
			res[i].ContractAddress,
			fee,
		})
	}

	return &transactions, nil
}

func (c *Client) getInternalTransactions(queryValues url.Values) (*[]Transaction, error) {
	type transactionResult struct {
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

	var res []transactionResult
	if err := c.sendRequest(queryValues, &res); err != nil {
		return nil, err
	}

	var transactions []Transaction
	for i := 0; i < len(res); i++ {
		timeInt, err := strconv.ParseInt(res[i].TimeStamp, 10, 64)
		if err != nil {
			panic(err)
		}
		timestamp := time.Unix(timeInt, 0)
		amount, err := utils.ConvertStringToBNB(res[i].Value)
		if err != nil {
			return nil, err
		}
		fee, err := utils.ConvertStringToBNB(res[i].GasUsed)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, Transaction{
			timestamp,
			res[i].Hash,
			res[i].From,
			res[i].To,
			amount,
			res[i].IsError != "0",
			res[i].ContractAddress,
			fee,
		})
	}

	return &transactions, nil
}

func (c *Client) GetInternalTransactionList(address string, options *TransactionsOptions) (*[]Transaction, error) {
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

	return c.getInternalTransactions(queryValues)
}

func (c *Client) GetInternalTransactionsTxList(txHash string) (*[]Transaction, error) {
	queryValues := url.Values{}
	queryValues.Add("module", "account")
	queryValues.Add("action", "txlistinternal")
	queryValues.Add("txhash", txHash)

	return c.getInternalTransactions(queryValues)
}

func (c *Client) GetTokenTransfers(address string, options *TransactionsOptions) (*[]TokenTransfer, error) {
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

	type tokenTransactionResult struct {
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

	var res []tokenTransactionResult
	if err := c.sendRequest(queryValues, &res); err != nil {
		return nil, err
	}

	var transactions []TokenTransfer
	for i := 0; i < len(res); i++ {
		timeInt, err := strconv.ParseInt(res[i].TimeStamp, 10, 64)
		if err != nil {
			panic(err)
		}
		timestamp := time.Unix(timeInt, 0)
		amount, err := utils.ConvertStringToBNB(res[i].Value)
		if err != nil {
			return nil, err
		}
		gas, err := strconv.ParseFloat(res[i].GasUsed, 64)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, TokenTransfer{
			timestamp,
			res[i].Hash,
			res[i].From,
			res[i].ContractAddress,
			res[i].To,
			amount,
			res[i].TokenName,
			res[i].TokenSymbol,
			res[i].TokenDecimal,
			gas,
		})
	}

	return &transactions, nil
}

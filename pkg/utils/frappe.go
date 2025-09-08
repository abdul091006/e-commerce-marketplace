package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type BalanceTypeResponse struct {
	Data []struct {
		Name     string `json:"name"`
		TypeName string `json:"type_name"`
	} `json:"data"`
}

// ValidateBalanceTypeFromFrappe checks balance type against Frappe
func ValidateBalanceTypeFromFrappe(balanceType string) error {
	baseURL := os.Getenv("FRAPPE_URL")       
	apiKey := os.Getenv("FRAPPE_API_KEY")
	apiSecret := os.Getenv("FRAPPE_API_SECRET")

	url := fmt.Sprintf("%s/api/resource/Balance%%20Type?fields=[\"name\",\"type_name\"]", baseURL)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+apiKey+":"+apiSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return NewWalletError(CodeInternalError, "failed to connect to frappe", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewWalletError(CodeInternalError, "failed to fetch balance types from frappe", fmt.Sprintf("status code: %d", resp.StatusCode))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("FRAPPE RESPONSE:", string(body)) // DEBUG

	var result BalanceTypeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return NewWalletError(CodeInternalError, "invalid response from frappe", string(body))
	}

	
	for _, b := range result.Data {
		if strings.TrimSpace(b.TypeName) == strings.TrimSpace(balanceType) {
			return nil
		}
	}

	return NewWalletError(CodeInvalidBalanceType, "invalid balance type", "type not found in frappe")
}

// GetAllBalanceTypesFromFrappe fetch list of all balance types
func GetAllBalanceTypesFromFrappe() ([]string, error) {
	baseURL := os.Getenv("FRAPPE_URL")
	apiKey := os.Getenv("FRAPPE_API_KEY")
	apiSecret := os.Getenv("FRAPPE_API_SECRET")

	url := fmt.Sprintf("%s/api/resource/Balance%%20Type?fields=[\"name\",\"type_name\"]", baseURL)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+apiKey+":"+apiSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []map[string]string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	types := []string{}
	for _, item := range result.Data {
		t := strings.TrimSpace(item["type_name"])
		if t != "" {
			types = append(types, t)
		}
	}
	return types, nil
}
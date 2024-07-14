package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"sync"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type KeyVaultResponse []struct {
	Attributes struct {
		Created         time.Time   `json:"created"`
		Enabled         bool        `json:"enabled"`
		Expires         interface{} `json:"expires"`
		NotBefore       interface{} `json:"notBefore"`
		RecoverableDays int         `json:"recoverableDays"`
		RecoveryLevel   string      `json:"recoveryLevel"`
		Updated         time.Time   `json:"updated"`
	} `json:"attributes"`
	ContentType interface{} `json:"contentType"`
	ID          string      `json:"id"`
	Managed     interface{} `json:"managed"`
	Name        string      `json:"name"`
	Tags        struct {
	} `json:"tags"`
}

func getValuesFromKeys(kvname string, keys []string) map[string]string {
    query := `az keyvault secret show --name %s --vault-name %s --query value -o tsv`
    values := make(map[string]string)
    mu := &sync.Mutex{}

    const numWorkers = 10
    keysChan := make(chan string, numWorkers)

    var wg sync.WaitGroup
    wg.Add(numWorkers)

    for i := 0; i < numWorkers; i++ {
        go func() {
            defer wg.Done()
            for k := range keysChan {
                cmd := exec.Command("bash", "-c", fmt.Sprintf(query, k, kvname))
                out, err := cmd.Output()
                if err != nil {
                    fmt.Printf("Error fetching key %s: %v\n", k, err)
                    continue
                }
                mu.Lock()
                values[k] = string(out)
                mu.Unlock()
            }
        }()
    }

    for _, key := range keys {
        keysChan <- key
    }

    close(keysChan)
    wg.Wait()

    return values
}

func ListKvKeys(name string) []string {
	var response KeyVaultResponse
	const query = `az keyvault secret list --vault-name %s`
	cmd := exec.Command("bash", "-c", fmt.Sprintf(query, name))
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(out, &response)
	if err != nil {
		panic(err)
	}
	var keys []string
	for _, key := range response {
		keys = append(keys, key.Name)
	}
	return keys
}

func displayUsage() {
	fmt.Println("Usage: secrets show <vault-name>")
	fmt.Println("Usage: secrets diff <vault-name> <vault-name>")
	fmt.Println("Usage: secrets list <vault-name>")
	fmt.Println("Usage: secrets get <vault-name> <key-name>")
}

func getKvKeyValPairsAsString(name string) string {
	keys := ListKvKeys(name)
	values := getValuesFromKeys(name, keys)

	// Extract the keys and sort them
	var sortedKeys []string
	for key := range values {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// Iterate over the sorted keys to access the values
	var kvPairs string
	for _, key := range sortedKeys {
		kvPairs += fmt.Sprintf("%s: %s", key, values[key])
	}

	return kvPairs
}

func main() {
	args := os.Args
	if len(args) < 3 {
		displayUsage()
		return
	}

	switch args[1] {
	case "get":
		vaultName := args[2]
		keyName := args[3]
		value := getValuesFromKeys(vaultName, []string{keyName})[keyName]
		fmt.Print(value)
	case "list":
		vaultName := args[2]
		keys := ListKvKeys(vaultName)
		for _, key := range keys {
			fmt.Println(key)
		}
	case "show":
		keys := ListKvKeys(args[2])
		fmt.Printf("Keys in %s:\n", args[2])
		fmt.Printf("Found %d keys\n", len(keys))
		values := getValuesFromKeys(args[2], keys)
		for key, value := range values {
			fmt.Printf("%s: %s", key, value)
		}

	case "diff":
		fmt.Println("Comparing keys in", args[2], "and", args[3])
		keys1 := getKvKeyValPairsAsString(args[2])
		keys2 := getKvKeyValPairsAsString(args[3])
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(keys1, keys2, false)
		fmt.Println(dmp.DiffPrettyText(diffs))

	default:
		fmt.Println("Invalid command")
    displayUsage()
	}
}

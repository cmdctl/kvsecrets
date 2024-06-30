# Azure Key Vault Secrets Utility

This utility provides a command-line interface to interact with Azure Key Vault secrets. It allows you to list, show, get, and compare secrets stored in different Azure Key Vaults.

## Features

- **List Secrets**: List all secrets in a specified Key Vault.
- **Show Secrets**: Display all secrets and their values in a specified Key Vault.
- **Get Secret**: Retrieve the value of a specific secret from a Key Vault.
- **Diff Secrets**: Compare secrets between two Key Vaults and display the differences.

## Prerequisites

- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) installed and configured.
- Valid Azure subscription with access to the Key Vaults.

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/cmdctl/kvsecrets.git
    cd kvsecrets
    ```

2. Build the project:

    ```sh
    go build -o secrets
    ```

## Usage

### List Secrets

List all secrets in a specified Key Vault.

```sh
./secrets list <vault-name>
```

### Show Secrets

Display all secrets and their values in a specified Key Vault.

```sh
./secrets show <vault-name>
```

### Get Secret

Retrieve the value of a specific secret from a Key Vault.

```sh
./secrets get <vault-name> <key-name>
```

### Diff Secrets

Compare secrets between two Key Vaults and display the differences.

```sh
./secrets diff <vault-name1> <vault-name2>
```

## Examples

### List Secrets

```sh
./secrets list my-vault
```

### Show Secrets

```sh
./secrets show my-vault
```

### Get Secret

```sh
./secrets get my-vault my-secret
```

### Diff Secrets

```sh
./secrets diff my-vault1 my-vault2
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
```


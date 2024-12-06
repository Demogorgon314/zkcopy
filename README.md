# Zookeeper Node Copier

This project provides a tool to copy nodes from one Zookeeper instance to another. It ensures that the destination path exists and copies the data and children nodes recursively.

## Prerequisites

- Go 1.23 or later
- Zookeeper instances

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/Demogorgon314/zkcopy.git
    cd zkcopy
    ```

2. Build the project:
    ```sh
    go build -o zkcopy
    ```

## Usage

```sh
./zkcopy -source-zk <source_zk_connection_string> -destination-zk <destination_zk_connection_string> -source-path <source_path> -destination-path <destination_path> [-delete-before-copy]
```

### Arguments

- `-source-zk`: Source Zookeeper connection string (e.g., `localhost:2181`)
- `-destination-zk`: Destination Zookeeper connection string (e.g., `localhost:2182`)
- `-source-path`: Source Zookeeper path (must start with `/`)
- `-destination-path`: Destination Zookeeper path (must start with `/`)
- `-delete-before-copy`: Delete destination path before copying (optional, default is `false`)

### Example

```sh
./zkcopy -source-zk localhost:2181 -destination-zk localhost:2182 -source-path /sourceNode -destination-path /destNode -delete-before-copy
```

## Development

### Running Tests

To run the tests, use the following command:

```sh
go test ./...
```

### GitHub Actions

This repository uses GitHub Actions for CI/CD. The workflow is defined in `.github/workflows/build-on-tag.yml` and triggers on tags starting with `v`.

## Contributing

1. Fork the repository
2. Create a new branch (`git checkout -b feature-branch`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature-branch`)
5. Create a new Pull Request
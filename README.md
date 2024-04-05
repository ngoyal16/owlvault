# OwlVault

We are excited to announce the initial release of OwlVault, your trusted key vault service designed to securely store and manage encryption keys. With OwlVault, you can safeguard your sensitive data and ensure its confidentiality with ease.

## Key Features

1. **Flexible Storage Options:** OwlVault supports multiple storage backends, including MySQL, DynamoDB, and more. Choose the storage solution that best fits your needs and seamlessly integrate OwlVault into your existing infrastructure.

2. **Robust Encryption:** Protect your data with strong encryption using customizable encryption algorithms. OwlVault provides support for various encryption methods, allowing you to tailor the encryption to your specific security requirements.

3. **Dynamic Key Providers:** OwlVault offers support for dynamic key providers, enabling you to retrieve encryption keys from various sources such as local files, AWS Key Management Service (KMS), and more. Easily manage and rotate encryption keys to enhance the security of your data.

4. **REST API-based Interaction:** Interact with OwlVault programmatically using its REST API. With API versioning support, you can ensure compatibility and seamless upgrades as the service evolves.

5. **Secure Communication:** OwlVault ensures secure communication between clients and the service using industry-standard encryption protocols, keeping your data protected during transit.

## How to Get Started

To start using OwlVault, follow these simple steps:

1. **Download the OwlVault binary:** Visit our GitHub repository and download the latest release of OwlVault for your platform.

2. **Configure OwlVault:** Modify the configuration file (`config.yaml`) to specify your desired storage backend, encryption settings, and key provider preferences.

3. **Run OwlVault:** Launch the OwlVault service using the provided binary and start storing and retrieving your encryption keys securely.

```shell
export OWLVAULT_CONFIG_PATH=/home/nitin/owlvault-go/config.yaml
```


# API Documentation

## Introduction
This document provides detailed information about the API endpoints and their functionalities for the Key Storage service (KS2).

## Base URL
`BASE_URL`: The base URL for accessing the KS2 API.

## Version
`v1`: Version 1 of the KS2 API.

## Endpoints

### 1. StoreKey
Store a key along with its associated data in the KS2.

#### Endpoint
`BASE_URL/v1/ks2?Action=StoreKey`

#### Method
POST

#### Input
- `keyPath` (string, required): The path to the key to be stored.
- `data` (object, required): The data associated with the key.

#### Output
- `requestId` (string): Unique identifier for the request.
- `data` (object):
    - `keyPath` (string): The path to the stored key.
    - `version` (integer): Version number of the stored key.

#### Response Codes
- `200 OK`: Successfully stored the key.
- `400 Bad Request`: Invalid input data.
- `500 Internal Server Error`: Server encountered an error while processing the request.

### 2. RetrieveKey
Retrieve the data associated with a stored key from the KS2.

#### Endpoint
`BASE_URL/v1/ks2?Action=RetrieveKey`

#### Method
POST

#### Input
- `keyPath` (string, required): The path to the key to be retrieved.
- `version` (integer, optional): The version number of the key to be retrieved. If not provided, the latest version will be retrieved.

#### Output
- `requestId` (string): Unique identifier for the request.
- `data` (object):
    - `keyPath` (string): The path to the retrieved key.
    - `data` (object): The data associated with the retrieved key.

#### Response Codes
- `200 OK`: Successfully retrieved the key and its associated data.
- `400 Bad Request`: Invalid input data.
- `404 Not Found`: Key not found.
- `500 Internal Server Error`: Server encountered an error while processing the request.

### 3. StoreKeys
Store multiple keys along with their associated data in the KS2.

#### Endpoint
`BASE_URL/v1/ks2?Action=StoreKeys`

#### Method
POST

#### Input
- `keysToStore` (array, required): An array of objects representing keys and their associated data to be stored.
    - `keyPath` (string, required): The path to the key to be stored.
    - `data` (object, required): The data associated with the key.

#### Output
- `requestId` (string): Unique identifier for the request.
- `data` (array): An array of objects representing the stored keys and their versions.
    - `keyPath` (string): The path to the stored key.
    - `version` (integer): Version number of the stored key.

#### Response Codes
- `200 OK`: Successfully stored all keys.
- `400 Bad Request`: Invalid input data.
- `500 Internal Server Error`: Server encountered an error while processing the request.

## Error Responses
In case of error, the response will include an error message along with the corresponding HTTP status code.

Example Error Response:
```json
{
    "error": "Invalid input data",
    "status": 400
}
```

#### Notes
1. Ensure that the keyPath for each key is unique.
2. All input and output data is in JSON format.
3. Handle errors gracefully by checking the HTTP status code and error message in the response.
4. Make sure to include proper authentication and authorization mechanisms for accessing the API endpoints.

## Feedback and Support

We value your feedback and are committed to continuously improving OwlVault to meet your needs. If you encounter any issues or have suggestions for enhancements, please don't hesitate to reach out to us through our GitHub repository or contact our support team.

Thank you for choosing OwlVault. We look forward to helping you safeguard your data and enhance your security posture.

Stay secure,
The OwlVault Team

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

#### Sample Input
```json
{
  "keyPath": "kv1",
  "data": {
    "key1": "val1",
    "key2": "val2"
  }
}
```

#### Output
- `requestId` (string): Unique identifier for the request.
- `data` (object):
    - `keyPath` (string): The path to the stored key.
    - `version` (integer): Version number of the stored key.

#### Sample Output
```json
{
  "requestId": "abcdabcd-abcd-abcd-abcd-abcdabcdabcd",
  "data": {
    "keyPath": "kv1",
    "version": 1
  }
}
```

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

#### Sample Input
For latest version
```json
{
  "keyPath": "kv1"
}
```

For specific version

```json
{
  "keyPath": "kv1",
  "version": 1
}
```

#### Output
- `requestId` (string): Unique identifier for the request.
- `data` (object):
    - `keyPath` (string): The path to the retrieved key.
    - `data` (object): The data associated with the retrieved key.

#### Sample Output
```json
{
  "requestId": "abcdabcd-abcd-abcd-abcd-abcdabcdabcd",
  "data": {
    "keyPath": "kv1",
    "data": {
      "key1": "val1",
      "key2": "val2"
    }
  }
}
```


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

#### Sample Input
```json
{
  "keysToStore": [
    {
      "keyPath": "kv1",
      "data": {
        "key1": "val1",
        "key2": "val2"
      }
    },
    {
      "keyPath": "kv2",
      "data": {
        "key1": "val1",
        "key2": "val2"
      }
    }
  ]
}
```

#### Output
- `requestId` (string): Unique identifier for the request.
- `data` (array): An array of objects representing the stored keys and their versions.
    - `keyPath` (string): The path to the stored key.
    - `version` (integer): Version number of the stored key.

#### Sample Output
```json
{
  "requestId": "abcdabcd-abcd-abcd-abcd-abcdabcdabcd",
  "data": [
    {
      "keyPath": "kv1",
      "version": 1
    },
    {
      "keyPath": "kv2",
      "version": 1
    }
  ]
}
```

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

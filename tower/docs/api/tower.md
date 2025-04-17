# Tower API Documentation

This document outlines the REST API endpoints provided by the Tower API Middleware service. Tower enables seamless connection and data transformation between different API systems.

## Base URL

All API endpoints are prefixed with `/api/v1`.

## Authentication

Authentication is not explicitly covered in this documentation. Implement appropriate authentication mechanisms for your production environment.

## Response Format

All API endpoints return JSON responses with a standard format:

```json
{
  "success": true,
  "data": { ... }  // Response data
}
```

For errors:

```json
{
  "success": false,
  "error": "Error message"
}
```

## API Endpoints

### Health Check

```
GET /health
```

Returns a simple status check for the API service.

**Response**: `200 OK` with body `"OK"` if the service is running.

---

### Connectors

#### List Connectors

```
GET /api/v1/connectors
```

Returns a list of all available connectors.

**Response**:

```json
{
  "success": true,
  "data": {
    "connectors": [
      {
        "id": "woocommerce",
        "name": "WooCommerce",
        "type": "woocommerce",
        "description": "Connector for product"
      },
      {
        "id": "mirakl",
        "name": "Mirakl",
        "type": "mirakl",
        "description": "Connector for product"
      }
    ]
  }
}
```

#### Get Connector Schema

```
GET /api/v1/connectors/{connectorID}/schema
```

Returns the data schema for a specific connector.

**Parameters**:
- `connectorID` - The ID of the connector

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "woocommerce",
    "entityName": "product",
    "fields": {
      "id": {
        "type": "number",
        "required": false,
        "path": "id"
      },
      "name": {
        "type": "string",
        "required": true,
        "path": "name"
      },
      "price": {
        "type": "number",
        "required": true,
        "path": "price"
      }
    }
  }
}
```

#### Test Connector

```
POST /api/v1/connectors/test
```

Tests the connection to a specific connector.

**Request Body**:

```json
{
  "id": "woocommerce",
  "config": {
    "api_url": "https://example.com/wp-json/wc/v3",
    "consumer_key": "your_consumer_key",
    "consumer_secret": "your_consumer_secret"
  }
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Connection successful"
  }
}
```

#### Fetch Data from Connector

```
POST /api/v1/connectors/fetch
```

Fetches data from a connector based on query parameters.

**Request Body**:

```json
{
  "id": "woocommerce",
  "query": {
    "entity_type": "product",
    "per_page": 10,
    "page": 1
  }
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "success": true,
    "count": 2,
    "data": [
      {
        "id": 1,
        "name": "Product 1",
        "price": "19.99"
      },
      {
        "id": 2,
        "name": "Product 2",
        "price": "29.99"
      }
    ]
  }
}
```

---

### Transformers

#### List Transformers

```
GET /api/v1/transformers
```

Returns a list of all transformers.

**Response**:

```json
{
  "success": true,
  "data": {
    "transformers": [
      {
        "id": "123456789",
        "name": "WooCommerce to Mirakl",
        "description": "Maps WooCommerce products to Mirakl format",
        "mappings": [
          {
            "sourceField": "name",
            "targetField": "product_title"
          },
          {
            "sourceField": "price",
            "targetField": "price"
          }
        ],
        "functions": [],
        "createdAt": "2025-04-17T10:30:00Z",
        "updatedAt": "2025-04-17T10:30:00Z"
      }
    ]
  }
}
```

#### Get Transformer

```
GET /api/v1/transformers/{transformerID}
```

Returns a specific transformer.

**Parameters**:
- `transformerID` - The ID of the transformer

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "123456789",
    "name": "WooCommerce to Mirakl",
    "description": "Maps WooCommerce products to Mirakl format",
    "mappings": [
      {
        "sourceField": "name",
        "targetField": "product_title"
      },
      {
        "sourceField": "price",
        "targetField": "price"
      }
    ],
    "functions": [],
    "createdAt": "2025-04-17T10:30:00Z",
    "updatedAt": "2025-04-17T10:30:00Z"
  }
}
```

#### Create Transformer

```
POST /api/v1/transformers
```

Creates a new transformer.

**Request Body**:

```json
{
  "name": "WooCommerce to Mirakl",
  "description": "Maps WooCommerce products to Mirakl format",
  "mappings": [
    {
      "sourceField": "name",
      "targetField": "product_title"
    },
    {
      "sourceField": "price",
      "targetField": "price"
    }
  ],
  "functions": []
}
```

**Response**: Same as Get Transformer with status `201 Created`.

#### Update Transformer

```
PUT /api/v1/transformers/{transformerID}
```

Updates an existing transformer.

**Parameters**:
- `transformerID` - The ID of the transformer

**Request Body**: Same as Create Transformer

**Response**: Same as Get Transformer

#### Delete Transformer

```
DELETE /api/v1/transformers/{transformerID}
```

Deletes a transformer.

**Parameters**:
- `transformerID` - The ID of the transformer

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "Transformer deleted successfully"
  }
}
```

#### Generate Transformer

```
POST /api/v1/transformers/generate
```

Automatically generates a transformer between two connectors.

**Request Body**:

```json
{
  "sourceId": "woocommerce",
  "targetId": "mirakl",
  "name": "WooCommerce to Mirakl",
  "description": "Auto-generated transformer from WooCommerce to Mirakl",
  "threshold": 0.7
}
```

**Response**: Same as Get Transformer with status `201 Created`.

#### Test Transformation

```
POST /api/v1/transformers/transform
```

Tests a transformation with sample data.

**Request Body**:

```json
{
  "transformerId": "123456789",
  "data": {
    "name": "Sample Product",
    "price": "19.99",
    "sku": "SAMPLE-001"
  }
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "result": {
      "product_title": "Sample Product",
      "price": "19.99"
    }
  }
}
```

---

### Connections

#### List Connections

```
GET /api/v1/connections
```

Returns a list of all connections.

**Response**:

```json
{
  "success": true,
  "data": {
    "connections": [
      {
        "id": "987654321",
        "name": "WooCommerce to Mirakl Sync",
        "description": "Sync products from WooCommerce to Mirakl",
        "sourceId": "woocommerce",
        "targetId": "mirakl",
        "transformerId": "123456789",
        "query": {
          "entity_type": "product",
          "per_page": 100
        },
        "schedule": "1h",
        "active": true,
        "lastRun": "2025-04-17T10:45:00Z",
        "status": "success",
        "createdAt": "2025-04-17T10:30:00Z",
        "updatedAt": "2025-04-17T10:45:00Z"
      }
    ]
  }
}
```

#### Get Connection

```
GET /api/v1/connections/{connectionID}
```

Returns a specific connection.

**Parameters**:
- `connectionID` - The ID of the connection

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "987654321",
    "name": "WooCommerce to Mirakl Sync",
    "description": "Sync products from WooCommerce to Mirakl",
    "sourceId": "woocommerce",
    "targetId": "mirakl",
    "transformerId": "123456789",
    "query": {
      "entity_type": "product",
      "per_page": 100
    },
    "schedule": "1h",
    "active": true,
    "lastRun": "2025-04-17T10:45:00Z",
    "status": "success",
    "createdAt": "2025-04-17T10:30:00Z",
    "updatedAt": "2025-04-17T10:45:00Z"
  }
}
```

#### Create Connection

```
POST /api/v1/connections
```

Creates a new connection.

**Request Body**:

```json
{
  "name": "WooCommerce to Mirakl Sync",
  "description": "Sync products from WooCommerce to Mirakl",
  "sourceId": "woocommerce",
  "targetId": "mirakl",
  "transformerId": "123456789",
  "query": {
    "entity_type": "product",
    "per_page": 100
  },
  "schedule": "1h",
  "active": true
}
```

**Response**: Same as Get Connection with status `201 Created`.

#### Update Connection

```
PUT /api/v1/connections/{connectionID}
```

Updates an existing connection.

**Parameters**:
- `connectionID` - The ID of the connection

**Request Body**: Same as Create Connection

**Response**: Same as Get Connection

#### Delete Connection

```
DELETE /api/v1/connections/{connectionID}
```

Deletes a connection.

**Parameters**:
- `connectionID` - The ID of the connection

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "Connection deleted successfully"
  }
}
```

#### Execute Connection

```
POST /api/v1/connections/{connectionID}/execute
```

Executes a connection once.

**Parameters**:
- `connectionID` - The ID of the connection

**Response**:

```json
{
  "success": true,
  "data": {
    "executionId": "567890123",
    "message": "Connection execution started"
  }
}
```

#### Set Connection Active Status

```
PATCH /api/v1/connections/{connectionID}/active
```

Enables or disables a connection.

**Parameters**:
- `connectionID` - The ID of the connection

**Request Body**:

```json
{
  "active": true
}
```

**Response**: Same as Get Connection

#### Get Connection Executions

```
GET /api/v1/connections/{connectionID}/executions
```

Returns the execution history for a connection.

**Parameters**:
- `connectionID` - The ID of the connection
- `page` (query) - Page number (default: 1)
- `size` (query) - Page size (default: 20)

**Response**:

```json
{
  "success": true,
  "data": {
    "executions": [
      {
        "id": "567890123",
        "connectionId": "987654321",
        "status": "success",
        "startTime": "2025-04-17T10:45:00Z",
        "endTime": "2025-04-17T10:45:15Z",
        "sourceData": [...],
        "targetData": [...],
        "createdAt": "2025-04-17T10:45:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  },
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "lastPage": 1
  }
}
```

## Error Codes

| Status Code | Meaning               |
|-------------|------------------------|
| 200         | OK                     |
| 201         | Created                |
| 400         | Bad Request            |
| 404         | Not Found              |
| 500         | Internal Server Error  |

## Rate Limiting

Rate limiting is not explicitly implemented in this version. Consider adding appropriate rate limiting for production deployments.

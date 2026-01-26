# WebSocket API Documentation

Dokumentasi untuk WebSocket API yang digunakan untuk fitur real-time pada Wallet Exhibition Application.

## Endpoint

```
ws://localhost:3000/ws?token=<jwt_token>
```

## Autentikasi

WebSocket menggunakan autentikasi via query parameter `token`. Token yang digunakan adalah JWT token yang didapat dari endpoint `/users/login`.

### Contoh Koneksi

```javascript
// JavaScript/Browser
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...";
const ws = new WebSocket(`ws://localhost:3000/ws?token=${token}`);

ws.onopen = function() {
    console.log("WebSocket connected");
};

ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log("Received:", message);
};

ws.onclose = function() {
    console.log("WebSocket disconnected");
};

ws.onerror = function(error) {
    console.error("WebSocket error:", error);
};
```

```go
// Go Client
import "github.com/gorilla/websocket"

token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
url := fmt.Sprintf("ws://localhost:3000/ws?token=%s", token)

conn, _, err := websocket.DefaultDialer.Dial(url, nil)
if err != nil {
    log.Fatal("dial:", err)
}
defer conn.Close()

for {
    _, message, err := conn.ReadMessage()
    if err != nil {
        log.Println("read:", err)
        return
    }
    log.Printf("received: %s", message)
}
```

## Format Message

Semua message yang dikirim melalui WebSocket menggunakan format JSON dengan struktur berikut:

```json
{
    "type": "<message_type>",
    "payload": { ... }
}
```

## Tipe Message

### 1. Transaction Notification

Dikirim ketika terjadi transaksi yang melibatkan pengguna (top-up atau transfer).

**Type:** `transaction`

**Payload:**

| Field | Type | Description |
|-------|------|-------------|
| transaction_id | integer | ID transaksi |
| transaction_type | string | Tipe transaksi (`top_up` atau `transfer`) |
| amount | string | Jumlah transaksi |
| from_user_id | integer (nullable) | ID pengirim (null untuk top-up) |
| to_user_id | integer | ID penerima |
| performed_by_user_id | integer | ID yang melakukan transaksi |
| description | string (nullable) | Deskripsi transaksi |
| created_at | string | Waktu transaksi (RFC3339 format) |

**Contoh:**

```json
{
    "type": "transaction",
    "payload": {
        "transaction_id": 1,
        "transaction_type": "transfer",
        "amount": "50000",
        "from_user_id": 1,
        "to_user_id": 2,
        "performed_by_user_id": 1,
        "description": "Bayar makan siang",
        "created_at": "2026-01-26T12:00:00+07:00"
    }
}
```

### 2. Wallet Update Notification

Dikirim ketika saldo wallet pengguna berubah.

**Type:** `wallet_update`

**Payload:**

| Field | Type | Description |
|-------|------|-------------|
| wallet_id | integer | ID wallet |
| new_balance | string | Saldo baru setelah mutasi |
| mutation_type | string | Tipe mutasi (`credit` atau `debit`) |
| mutation_id | integer | ID mutasi |
| transaction_id | integer | ID transaksi terkait |
| amount | string | Jumlah mutasi |
| updated_at | string | Waktu update (RFC3339 format) |

**Contoh:**

```json
{
    "type": "wallet_update",
    "payload": {
        "wallet_id": 1,
        "new_balance": "150000",
        "mutation_type": "credit",
        "mutation_id": 1,
        "transaction_id": 1,
        "amount": "50000",
        "updated_at": "2026-01-26T12:00:00+07:00"
    }
}
```

## Use Cases

### 1. Menerima Notifikasi Top-Up

Ketika super admin melakukan top-up ke wallet pengguna:

1. Pengguna terkoneksi ke WebSocket
2. Super admin melakukan POST ke `/transactions/topup`
3. Pengguna menerima message `transaction` dan `wallet_update`

### 2. Menerima Notifikasi Transfer

Ketika pengguna lain melakukan transfer:

1. Pengguna A dan B terkoneksi ke WebSocket
2. Pengguna A melakukan POST ke `/transactions/transfer` ke B
3. Pengguna A menerima konfirmasi (message `transaction` dan `wallet_update` dengan `mutation_type: debit`)
4. Pengguna B menerima notifikasi (message `transaction` dan `wallet_update` dengan `mutation_type: credit`)

## Error Handling

### Connection Errors

| Error | Description | Solution |
|-------|-------------|----------|
| 401 Unauthorized | Token tidak valid atau tidak diberikan | Pastikan token JWT valid dan belum expired |
| 426 Upgrade Required | Request bukan WebSocket upgrade | Gunakan WebSocket client yang benar |

### Connection Close Codes

| Code | Description |
|------|-------------|
| 1000 | Normal closure |
| 1001 | Going away (server shutdown) |
| 1006 | Abnormal closure (connection lost) |

## Best Practices

1. **Reconnection**: Implementasikan reconnection logic jika koneksi terputus
2. **Heartbeat**: Kirim ping secara berkala untuk menjaga koneksi tetap aktif
3. **Token Refresh**: Reconnect dengan token baru sebelum token expired
4. **Error Handling**: Handle semua error dan close events dengan baik

### Contoh Reconnection Logic

```javascript
class WebSocketClient {
    constructor(url) {
        this.url = url;
        this.reconnectInterval = 5000;
        this.connect();
    }

    connect() {
        this.ws = new WebSocket(this.url);
        
        this.ws.onopen = () => {
            console.log("Connected");
        };
        
        this.ws.onmessage = (event) => {
            this.handleMessage(JSON.parse(event.data));
        };
        
        this.ws.onclose = () => {
            console.log("Disconnected, reconnecting...");
            setTimeout(() => this.connect(), this.reconnectInterval);
        };
    }

    handleMessage(message) {
        switch(message.type) {
            case "transaction":
                this.onTransaction(message.payload);
                break;
            case "wallet_update":
                this.onWalletUpdate(message.payload);
                break;
        }
    }

    onTransaction(data) {
        console.log("Transaction received:", data);
    }

    onWalletUpdate(data) {
        console.log("Wallet updated:", data);
    }
}

// Usage
const token = "your_jwt_token";
const client = new WebSocketClient(`ws://localhost:3000/ws?token=${token}`);
```

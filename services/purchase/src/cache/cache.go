package serviceCache

import (
	"encoding/json"
	"log"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/dgraph-io/ristretto/v2"
)

const (
	MaxCacheSize = 256 << 20 //  256 MB
	DefaultTtl   = 5 * time.Minute

	CacheAuthEmailToToken      = "auth:%s"
	CacheUserIdToProfile       = "user:%s"
	CacheInvalidatedUserIds    = "inv_usr" // Value is comma-separated, e.g., 1,3,5
	CacheEmployeesWithParams   = "employees:v%d:%s"
	CacheDepartmentsWithParams = "departments:v%d:%s"

	CacheProductById = "product:%s"
	CacheSellerById  = "seller:%s"
)

var (
	EmployeeNamespaceVersion   atomic.Int64
	DepartmentNamespaceVersion atomic.Int64

	ProductNamespaceVersion atomic.Int64
	SellerNamespaceVersion  atomic.Int64
)

var Cache *ristretto.Cache[string, string]

// Initialize menginisialisasi cache dengan konfigurasi default dan mengatur versi namespace.
// Fungsi ini membuat cache Ristretto dengan pelacakan frekuensi, batas biaya maksimum,
// dan buffer untuk operasi Get. Jika inisialisasi gagal, program akan dihentikan.
func Initialize() {
	var err error

	Cache, err = ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e6,          // 1 juta penghitung untuk pelacakan frekuensi
		MaxCost:     MaxCacheSize, // Ukuran maksimum cache
		BufferItems: 64,           // Buffer 64 kunci per operasi Get
	})
	if err != nil {
		log.Fatalf("failed to initialize cache: %v", err)
	}

	// Mengatur versi namespace awal
	EmployeeNamespaceVersion.Store(1)
	DepartmentNamespaceVersion.Store(1)
	ProductNamespaceVersion.Store(1)
	SellerNamespaceVersion.Store(1)
}

// Set menyimpan pasangan kunci-nilai dalam cache dengan biaya yang dihitung otomatis.
//
// Parameters:
//   - key: Kunci string untuk menyimpan nilai di cache.
//   - value: Nilai string yang akan disimpan di cache.
func Set(key string, value string) {
	cost := int64(len(key) + len(value))
	SetWithCost(key, value, cost)
}

// SetAsMap menyimpan map sebagai nilai dalam cache.
// Data map akan diserialisasi menjadi string JSON, dan biaya dihitung
// berdasarkan panjang data yang telah diserialisasi dan kunci.
//
// Parameters:
//   - key: Kunci string untuk menyimpan nilai di cache.
//   - value: Map (key-value) yang akan disimpan sebagai nilai.
func SetAsMap(key string, value map[string]string) {
	data, err := sonic.Marshal(value)
	if err != nil {
		panic(err)
	}

	cost := int64(len(key) + len(data))
	SetWithCost(key, string(data), cost)
}

// SetAsMapArrayWithTtlAndCostMultiplier menyimpan array map ke dalam cache dengan TTL
// (Time-to-Live) dan pengganda biaya yang ditentukan. Data array akan diserialisasi ke JSON.
//
// Parameters:
//   - key: Kunci string untuk menyimpan nilai di cache.
//   - value: Array map yang akan disimpan sebagai nilai.
//   - costMultiplier: Faktor pengali untuk menghitung total biaya penyimpanan.
//   - ttl: Durasi waktu cache tetap aktif sebelum kedaluwarsa.
func SetAsMapArrayWithTtlAndCostMultiplier(
	key string,
	value []map[string]string,
	costMultiplier int,
	ttl time.Duration,
) {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	totalCost := int64(len(key))
	for _, v := range value {
		for k, str := range v {
			stringCost := int64(len(k)) + int64(len(str))
			totalCost += stringCost
		}
	}

	Cache.SetWithTTL(key, string(data), totalCost*int64(costMultiplier), ttl)
}

// SetWithCost menyimpan pasangan kunci-nilai dalam cache dengan biaya tertentu
// dan TTL default.
//
// Parameters:
//   - key: Kunci string untuk menyimpan nilai di cache.
//   - value: Nilai string yang akan disimpan di cache.
//   - cost: Biaya penyimpanan berdasarkan ukuran kunci dan nilai.
func SetWithCost(key string, value string, cost int64) {
	Cache.SetWithTTL(key, value, cost, DefaultTtl)
}

// Get mengambil nilai string dari cache berdasarkan kunci.
//
// Parameters:
//   - key: Kunci string untuk mengambil nilai dari cache.
//
// Returns:
//   - string: Nilai string yang ditemukan di cache.
//   - bool: Boolean untuk menunjukkan apakah kunci ditemukan.
func Get(key string) (string, bool) {
	return Cache.Get(key)
}

// GetAsMap mengambil nilai map dari cache berdasarkan kunci.
// Data yang diambil akan di-deserialize dari string JSON.
//
// Parameters:
//   - key: Kunci string untuk mengambil nilai dari cache.
//
// Returns:
//   - map[string]string: Map string yang ditemukan di cache.
//   - bool: Boolean untuk menunjukkan apakah kunci ditemukan.
func GetAsMap(key string) (map[string]string, bool) {
	val, found := Cache.Get(key)

	if !found {
		return nil, false
	}

	var result map[string]string
	err := sonic.Unmarshal([]byte(val), &result)
	if err != nil {
		log.Printf("failed to unmarshal cache value: %v", val)
		panic(err)
	}

	return result, true
}

// GetAsMapArray mengambil nilai array map dari cache berdasarkan kunci.
// Data yang diambil akan di-deserialize dari string JSON.
//
// Parameters:
//   - key: Kunci string untuk mengambil nilai dari cache.
//
// Returns:
//   - []map[string]string: Array map string yang ditemukan di cache.
//   - bool: Boolean untuk menunjukkan apakah kunci ditemukan.
func GetAsMapArray(key string) ([]map[string]string, bool) {
	val, found := Cache.Get(key)
	if !found {
		return nil, false
	}

	var result []map[string]string
	err := json.Unmarshal([]byte(val), &result)
	if err != nil {
		panic(err)
	}

	return result, true
}

// Delete menghapus nilai dari cache berdasarkan kunci.
//
// Parameters:
//   - key: Kunci string untuk menghapus nilai dari cache.
func Delete(key string) {
	Cache.Del(key)
}

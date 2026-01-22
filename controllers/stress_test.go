package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pclubiitk/puppylove2.0_backend/db"
	"github.com/pclubiitk/puppylove2.0_backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =============================================================================
// CONFIGURATION - Adjust these for different stress levels
// =============================================================================

const (
	// Large scale concurrent tests
	STRESS_HEARTS_CONCURRENT = 2000 // Number of hearts for concurrent claims
	STRESS_VERIFY_CONCURRENT = 2000 // Number of hearts for concurrent verification

	// Contention tests (many users -> same resource)
	STRESS_SINGLE_HEART_USERS  = 500 // Users competing for 1 heart
	STRESS_SINGLE_VERIFY_USERS = 500 // Users competing to verify 1 heart

	// Batch size for DB inserts (performance optimization)
	DB_BATCH_SIZE = 100
)

func initTestDB(t *testing.T) *gorm.DB {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	database.AutoMigrate(
		&models.User{},
		&models.SendHeart{},
		&models.HeartClaims{},
		&models.ReturnHearts{},
		&models.MatchTable{},
	)

	Db = db.PuppyDb{DB: database}
	return database
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// batchInsertHearts inserts hearts in batches for better performance
func batchInsertHearts(database *gorm.DB, count int) {
	hearts := make([]models.SendHeart, 0, DB_BATCH_SIZE)
	for i := 0; i < count; i++ {
		hearts = append(hearts, models.SendHeart{
			SHA: fmt.Sprintf("sha_%d", i),
			ENC: fmt.Sprintf("enc_%d", i),
		})
		if len(hearts) >= DB_BATCH_SIZE {
			database.Create(&hearts)
			hearts = hearts[:0]
		}
	}
	if len(hearts) > 0 {
		database.Create(&hearts)
	}
}

// batchInsertUsers inserts users in batches for better performance
func batchInsertUsers(database *gorm.DB, prefix string, count int) {
	users := make([]models.User, 0, DB_BATCH_SIZE)
	for i := 0; i < count; i++ {
		users = append(users, models.User{Id: fmt.Sprintf("%s_%d", prefix, i)})
		if len(users) >= DB_BATCH_SIZE {
			database.Create(&users)
			users = users[:0]
		}
	}
	if len(users) > 0 {
		database.Create(&users)
	}
}

// =============================================================================
// STRESS TEST: HeartClaim - 2000 different hearts claimed concurrently
// =============================================================================

func TestStress_HeartClaim_Concurrent(t *testing.T) {
	database := initTestDB(t)
	router := setupRouter()
	router.POST("/heartclaim", func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		HeartClaim(c)
	})

	numHearts := STRESS_HEARTS_CONCURRENT

	t.Logf("Setting up %d hearts and users...", numHearts)
	start := time.Now()

	// Batch insert for performance
	batchInsertHearts(database, numHearts)
	batchInsertUsers(database, "user", numHearts)

	t.Logf("Setup complete in %v", time.Since(start))

	var wg sync.WaitGroup
	start = time.Now()

	for i := 0; i < numHearts; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			body, _ := json.Marshal(map[string]string{
				"sha":            fmt.Sprintf("sha_%d", idx),
				"enc":            fmt.Sprintf("enc_%d", idx),
				"genderOfSender": "M",
			})

			req := httptest.NewRequest("POST", "/heartclaim", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", fmt.Sprintf("user_%d", idx))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	var claimsInDB int64
	var unclaimed int64
	database.Model(&models.HeartClaims{}).Count(&claimsInDB)
	database.Model(&models.SendHeart{}).Count(&unclaimed)

	t.Logf("=== HeartClaim Concurrent Results ===")
	t.Logf("Total Hearts:    %d", numHearts)
	t.Logf("Claims in DB:    %d (%.1f%%)", claimsInDB, float64(claimsInDB)/float64(numHearts)*100)
	t.Logf("Unclaimed:       %d", unclaimed)
	t.Logf("Time elapsed:    %v", elapsed)
	t.Logf("Throughput:      %.1f claims/sec", float64(claimsInDB)/elapsed.Seconds())

	assert.Equal(t, int64(numHearts), claimsInDB, "All hearts should be claimed")
	assert.Equal(t, int64(0), unclaimed, "No hearts should remain")
}

// =============================================================================
// STRESS TEST: HeartClaim - 500 users trying to claim the SAME heart
// =============================================================================

func TestStress_HeartClaim_SingleHeart(t *testing.T) {
	database := initTestDB(t)
	router := setupRouter()
	router.POST("/heartclaim", func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		HeartClaim(c)
	})

	numUsers := STRESS_SINGLE_HEART_USERS

	// Create ONE heart and many users competing for it
	database.Create(&models.SendHeart{SHA: "shared_sha", ENC: "shared_enc"})
	batchInsertUsers(database, "user", numUsers)

	t.Logf("Testing %d users competing for 1 heart...", numUsers)

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			body, _ := json.Marshal(map[string]string{
				"sha":            "shared_sha",
				"enc":            "shared_enc",
				"genderOfSender": "M",
			})

			req := httptest.NewRequest("POST", "/heartclaim", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", fmt.Sprintf("user_%d", idx))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	var claimsInDB int64
	database.Model(&models.HeartClaims{}).Count(&claimsInDB)

	t.Logf("=== HeartClaim Single Heart Contention ===")
	t.Logf("Competing users: %d", numUsers)
	t.Logf("Claims in DB:    %d (expected: 1)", claimsInDB)
	t.Logf("Time elapsed:    %v", elapsed)

	assert.Equal(t, int64(1), claimsInDB, "Only 1 claim should exist - no duplicates!")
}

// =============================================================================
// STRESS TEST: VerifyReturnHeart - 2000 different hearts verified concurrently
// =============================================================================

func TestStress_VerifyReturnHeart_Concurrent(t *testing.T) {
	database := initTestDB(t)
	router := setupRouter()
	router.POST("/verify", func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		VerifyReturnHeart(c)
	})

	numHearts := STRESS_VERIFY_CONCURRENT

	t.Logf("Setting up %d return hearts and claims...", numHearts)
	start := time.Now()

	// Batch insert return hearts and claims
	returnHearts := make([]models.ReturnHearts, 0, DB_BATCH_SIZE)
	heartClaims := make([]models.HeartClaims, 0, DB_BATCH_SIZE)

	for i := 0; i < numHearts; i++ {
		secret := fmt.Sprintf("secret_%d", i)
		h := sha256.New()
		h.Write([]byte(secret))
		hash := fmt.Sprintf("%x", h.Sum(nil))
		enc := fmt.Sprintf("enc_%d", i)

		returnHearts = append(returnHearts, models.ReturnHearts{SHA: hash, ENC: enc})
		heartClaims = append(heartClaims, models.HeartClaims{
			Id:   fmt.Sprintf("claim_id_%d", i),
			SHA:  hash,
			Roll: fmt.Sprintf("other_%d", i),
		})

		if len(returnHearts) >= DB_BATCH_SIZE {
			database.Create(&returnHearts)
			database.Create(&heartClaims)
			returnHearts = returnHearts[:0]
			heartClaims = heartClaims[:0]
		}
	}
	if len(returnHearts) > 0 {
		database.Create(&returnHearts)
		database.Create(&heartClaims)
	}

	t.Logf("Setup complete in %v", time.Since(start))

	var wg sync.WaitGroup
	start = time.Now()

	for i := 0; i < numHearts; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			body, _ := json.Marshal(map[string]string{
				"secret": fmt.Sprintf("secret_%d", idx),
				"enc":    fmt.Sprintf("enc_%d", idx),
			})

			req := httptest.NewRequest("POST", "/verify", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", fmt.Sprintf("user_%d", idx))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	var matchesInDB int64
	var unverified int64
	database.Model(&models.MatchTable{}).Count(&matchesInDB)
	database.Model(&models.ReturnHearts{}).Count(&unverified)

	t.Logf("=== VerifyReturnHeart Concurrent Results ===")
	t.Logf("Total Hearts:    %d", numHearts)
	t.Logf("Matches in DB:   %d (%.1f%%)", matchesInDB, float64(matchesInDB)/float64(numHearts)*100)
	t.Logf("Unverified:      %d", unverified)
	t.Logf("Time elapsed:    %v", elapsed)
	t.Logf("Throughput:      %.1f matches/sec", float64(matchesInDB)/elapsed.Seconds())

	assert.Equal(t, int64(numHearts), matchesInDB, "All should create matches")
	assert.Equal(t, int64(0), unverified, "No return hearts should remain")
}

// =============================================================================
// STRESS TEST: VerifyReturnHeart - 500 users verifying the SAME heart
// =============================================================================

func TestStress_VerifyReturnHeart_SingleHeart(t *testing.T) {
	database := initTestDB(t)
	router := setupRouter()
	router.POST("/verify", func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		VerifyReturnHeart(c)
	})

	numUsers := STRESS_SINGLE_VERIFY_USERS

	secret := "shared_secret"
	h := sha256.New()
	h.Write([]byte(secret))
	hash := fmt.Sprintf("%x", h.Sum(nil))

	database.Create(&models.ReturnHearts{SHA: hash, ENC: "shared_enc"})
	database.Create(&models.HeartClaims{SHA: hash, Roll: "other_user"})

	t.Logf("Testing %d users competing to verify 1 heart...", numUsers)

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			body, _ := json.Marshal(map[string]string{
				"secret": secret,
				"enc":    "shared_enc",
			})

			req := httptest.NewRequest("POST", "/verify", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", fmt.Sprintf("user_%d", idx))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	var matchesInDB int64
	database.Model(&models.MatchTable{}).Count(&matchesInDB)

	t.Logf("=== VerifyReturnHeart Single Heart Contention ===")
	t.Logf("Competing users: %d", numUsers)
	t.Logf("Matches in DB:   %d (expected: 1)", matchesInDB)
	t.Logf("Time elapsed:    %v", elapsed)

	assert.Equal(t, int64(1), matchesInDB, "Only 1 match should exist - no duplicates!")
}

// =============================================================================
// STRESS TEST: Mixed Load - Simultaneous HeartClaim and VerifyReturnHeart
// =============================================================================

func TestStress_MixedLoad(t *testing.T) {
	database := initTestDB(t)
	router := setupRouter()

	router.POST("/heartclaim", func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		HeartClaim(c)
	})
	router.POST("/verify", func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		VerifyReturnHeart(c)
	})

	numHeartClaims := 1000
	numVerifies := 1000

	t.Logf("Setting up mixed load: %d claims + %d verifies", numHeartClaims, numVerifies)
	start := time.Now()

	// Setup HeartClaim data
	batchInsertHearts(database, numHeartClaims)
	batchInsertUsers(database, "claim_user", numHeartClaims)

	// Setup VerifyReturnHeart data
	returnHearts := make([]models.ReturnHearts, 0, DB_BATCH_SIZE)
	heartClaims := make([]models.HeartClaims, 0, DB_BATCH_SIZE)

	for i := 0; i < numVerifies; i++ {
		secret := fmt.Sprintf("verify_secret_%d", i)
		h := sha256.New()
		h.Write([]byte(secret))
		hash := fmt.Sprintf("%x", h.Sum(nil))
		enc := fmt.Sprintf("verify_enc_%d", i)

		returnHearts = append(returnHearts, models.ReturnHearts{SHA: hash, ENC: enc})
		heartClaims = append(heartClaims, models.HeartClaims{
			Id:   fmt.Sprintf("verify_claim_%d", i),
			SHA:  hash,
			Roll: fmt.Sprintf("verify_other_%d", i),
		})

		if len(returnHearts) >= DB_BATCH_SIZE {
			database.Create(&returnHearts)
			database.Create(&heartClaims)
			returnHearts = returnHearts[:0]
			heartClaims = heartClaims[:0]
		}
	}
	if len(returnHearts) > 0 {
		database.Create(&returnHearts)
		database.Create(&heartClaims)
	}

	t.Logf("Setup complete in %v", time.Since(start))

	var wg sync.WaitGroup
	start = time.Now()

	// Launch HeartClaim goroutines
	for i := 0; i < numHeartClaims; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			body, _ := json.Marshal(map[string]string{
				"sha":            fmt.Sprintf("sha_%d", idx),
				"enc":            fmt.Sprintf("enc_%d", idx),
				"genderOfSender": "M",
			})
			req := httptest.NewRequest("POST", "/heartclaim", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", fmt.Sprintf("claim_user_%d", idx))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}(i)
	}

	// Launch VerifyReturnHeart goroutines
	for i := 0; i < numVerifies; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			body, _ := json.Marshal(map[string]string{
				"secret": fmt.Sprintf("verify_secret_%d", idx),
				"enc":    fmt.Sprintf("verify_enc_%d", idx),
			})
			req := httptest.NewRequest("POST", "/verify", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", fmt.Sprintf("verify_user_%d", idx))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	var claimsInDB int64
	var matchesInDB int64
	var unclaimedHearts int64
	var unverifiedHearts int64

	// Count new claims (those with Roll starting with "claim_user_")
	database.Model(&models.HeartClaims{}).Where("roll LIKE 'claim_user_%'").Count(&claimsInDB)
	database.Model(&models.MatchTable{}).Count(&matchesInDB)
	database.Model(&models.SendHeart{}).Count(&unclaimedHearts)
	database.Model(&models.ReturnHearts{}).Count(&unverifiedHearts)

	t.Logf("=== Mixed Load Results ===")
	t.Logf("HeartClaims:     %d/%d (%.1f%%)", claimsInDB, numHeartClaims, float64(claimsInDB)/float64(numHeartClaims)*100)
	t.Logf("Matches:         %d/%d (%.1f%%)", matchesInDB, numVerifies, float64(matchesInDB)/float64(numVerifies)*100)
	t.Logf("Unclaimed:       %d", unclaimedHearts)
	t.Logf("Unverified:      %d", unverifiedHearts)
	t.Logf("Time elapsed:    %v", elapsed)
	t.Logf("Total ops:       %d", numHeartClaims+numVerifies)
	t.Logf("Throughput:      %.1f ops/sec", float64(numHeartClaims+numVerifies)/elapsed.Seconds())

	// Mixed load test: verify no data loss - source data should be consumed
	// Note: HeartClaims count may vary due to VerifyReturnHeart deleting entries
	assert.Equal(t, int64(0), unclaimedHearts, "All SendHearts should be processed")
	assert.Equal(t, int64(0), unverifiedHearts, "All ReturnHearts should be processed")
	assert.Equal(t, int64(numHeartClaims), claimsInDB, "All hearts should be claimed")
	assert.Equal(t, int64(numVerifies), matchesInDB, "All verifies should create matches")
}

// =============================================================================
// STRESS TEST: Burst Load - Rapid successive waves
// =============================================================================

func TestStress_BurstLoad(t *testing.T) {
	database := initTestDB(t)
	router := setupRouter()
	router.POST("/heartclaim", func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-User-ID"))
		HeartClaim(c)
	})

	wavesCount := 5
	heartsPerWave := 500
	totalHearts := wavesCount * heartsPerWave

	t.Logf("Setting up %d waves x %d hearts = %d total", wavesCount, heartsPerWave, totalHearts)

	// Setup all data upfront
	batchInsertHearts(database, totalHearts)
	batchInsertUsers(database, "user", totalHearts)

	var totalClaims int64
	start := time.Now()

	for wave := 0; wave < wavesCount; wave++ {
		var wg sync.WaitGroup
		waveStart := wave * heartsPerWave

		for i := 0; i < heartsPerWave; i++ {
			idx := waveStart + i
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				body, _ := json.Marshal(map[string]string{
					"sha":            fmt.Sprintf("sha_%d", idx),
					"enc":            fmt.Sprintf("enc_%d", idx),
					"genderOfSender": "M",
				})
				req := httptest.NewRequest("POST", "/heartclaim", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-ID", fmt.Sprintf("user_%d", idx))
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
			}(idx)
		}
		wg.Wait()

		var waveClaims int64
		database.Model(&models.HeartClaims{}).Count(&waveClaims)
		t.Logf("Wave %d complete: %d/%d claims so far", wave+1, waveClaims, (wave+1)*heartsPerWave)
	}

	elapsed := time.Since(start)
	database.Model(&models.HeartClaims{}).Count(&totalClaims)

	var unclaimed int64
	database.Model(&models.SendHeart{}).Count(&unclaimed)

	t.Logf("=== Burst Load Results ===")
	t.Logf("Total waves:     %d", wavesCount)
	t.Logf("Hearts/wave:     %d", heartsPerWave)
	t.Logf("Total claims:    %d/%d (%.1f%%)", totalClaims, totalHearts, float64(totalClaims)/float64(totalHearts)*100)
	t.Logf("Unclaimed:       %d", unclaimed)
	t.Logf("Time elapsed:    %v", elapsed)
	t.Logf("Throughput:      %.1f claims/sec", float64(totalClaims)/elapsed.Seconds())

	assert.Equal(t, int64(totalHearts), totalClaims, "All hearts should be claimed across waves")
}

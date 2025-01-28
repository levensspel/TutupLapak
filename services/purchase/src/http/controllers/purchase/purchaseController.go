package appController

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	serviceCache "github.com/TimDebug/FitByte/src/cache"
	helper "github.com/TimDebug/FitByte/src/helper/validator"
	functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"
	loggerZap "github.com/TimDebug/FitByte/src/logger/zap"
	"github.com/TimDebug/FitByte/src/model/dtos/request"
	"github.com/TimDebug/FitByte/src/model/dtos/response"
	"github.com/TimDebug/FitByte/src/services/proto/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
)

type PurchaseController struct {
	logger             loggerZap.LoggerInterface
	validator          helper.XValidator
	userServiceAddress string
}

func NewPurchaseController(logger loggerZap.LoggerInterface) IPurchaseController {
	xValidator := helper.XValidator{Validator: validator.New()}
	xValidator.Validator.RegisterValidation("sender_email_or_phone", func(fl validator.FieldLevel) bool {
		contactType := fl.Parent().FieldByName("SenderContactType").String()
		value := fl.Field().String()

		switch contactType {
		case "email":
			emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
			return regexp.MustCompile(emailRegex).MatchString(value)
		case "phone":
			// phoneRegex := `^\+?[1-9]\d{1,14}$` // E.164 format
			phoneRegex := `^[+0]\d{1,14}$`
			return regexp.MustCompile(phoneRegex).MatchString(value)
		default:
			return false
		}
	})

	return &PurchaseController{logger: logger, validator: xValidator, userServiceAddress: "host.docker.internal:50051"}
}

func NewPurchaseControllerInject(i do.Injector) (IPurchaseController, error) {
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)
	return NewPurchaseController(_logger), nil
}

// Purchase godoc
// @Summary Add items to the cart
// @Description Pembeli akan memasukkan detail produk dan jumlah yang akan dibeli, kemudian mengembalikan daftar detail produk beserta dengan daftar detail bank dari masing-masing penjual
// @Tags Purchase
// @Accept json
// @Produce json
// @Param request body request.CartDto true "Cart Data"
// @Success 201 {object} response.PurchaseResponseDTO "success response"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /v1/purchase [post]
func (pc *PurchaseController) Cart(c *fiber.Ctx) error {
	//// todo; Parse Body
	requestBody := new(request.CartDto)
	if err := c.BodyParser(requestBody); err != nil {
		pc.logger.Error(err.Error(), functionCallerInfo.PurhcaseControllerPutCart)
		return err
	}

	// Validation
	if errs := pc.validator.Validate(requestBody); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}
		pc.logger.Error(strings.Join(errMsgs, " || "), functionCallerInfo.PurhcaseControllerPutCart)
		return fiber.NewError(fiber.StatusBadRequest, strings.Join(errMsgs, " || "))
	}

	//// todo; buat cart
	// detail dari cart nanti ditambahkan setelah produk dan seller semua terkompilasi
	var cart response.PurchaseResponseDTO

	// sambil jalan mengambil produk dan seller, juga akumulasi totalPrice
	cart.TotalPrice = 0

	//// todo; check cache untuk ProdukID

	var cachedProducts []response.ProductItemDTO
	var cachedSeller []response.SellerBankDetailDTO

	// teknik indexing ketimbang for loop dan "find" dari go-lodash, gimik cuy
	//
	// parameter
	// 	- string: sellerId
	// 	- int: index
	var mapSellerId map[string]int
	mapSellerId = make(map[string]int)
	// berperan sebagai pointer index untuk mapSellerId
	var pointerIndexSellerId int = 0
	// berisi total harga yang akan dibawar oleh pembeli, index berdasarkan sellerId
	var sellerIdTotalPrices []float64

	//// todo; buat satu kontainer apabila produk dan seller tidak ada di cache
	var toGetProductsById []string
	var toGetSellersById []string

	// todo; task (go-routine)
	// forloop ambil cache by productIDs
	for _, item := range requestBody.PurchasedItems {
		if cachedProductValue, found := serviceCache.GetAsMap(fmt.Sprintf(serviceCache.CacheProductById, item.ProductId)); found {
			/*
				type ProductItemDTO struct {
					ProductId        string    `json:"productId"`
					Name             string    `json:"name"`
					Category         string    `json:"category"`
					Qty              int       `json:"qty"`
					Price            float64   `json:"price"`
					SKU              string    `json:"sku"`
					FileID           string    `json:"fileId"`
					FileURI          string    `json:"fileUri"`
					FileThumbnailURI string    `json:"fileThumbnailUri"`
					CreatedAt        time.Time `json:"createdAt"`
					UpdatedAt        time.Time `json:"updatedAt"`
				}
			*/
			_qty, err := strconv.Atoi(cachedProductValue["Qty"]) // Atoi = ASCII to Integer
			if err != nil {
				// Tangani error jika string tidak bisa dikonversi ke angka
				fmt.Println("Error:", err)
				pc.logger.Error(err.Error(), functionCallerInfo.CachePurhcaseControllerPutCart, "Parse Qty")
				return err
			}
			_price, err := strconv.ParseFloat(cachedProductValue["Price"], 64) // Atoi = ASCII to Integer
			if err != nil {
				// Tangani error jika string tidak bisa dikonversi ke floating number
				fmt.Println("Error:", err)
				pc.logger.Error(err.Error(), functionCallerInfo.CachePurhcaseControllerPutCart, "Parse Price")
				return err
			}
			_createdAt, err := time.Parse(time.RFC3339, cachedProductValue["CreatedAt"])
			if err != nil {
				// Tangani error jika string tidak bisa dikonversi ke Time sesuai dengan ISO
				fmt.Println("Error:", err)
				pc.logger.Error(err.Error(), functionCallerInfo.CachePurhcaseControllerPutCart, "Parse CreatedAt")
				return err
			}
			_modifiedAt, err := time.Parse(time.RFC3339, cachedProductValue["UpdatedAt"])
			if err != nil {
				// Tangani error jika string tidak bisa dikonversi ke Time sesuai dengan ISO
				fmt.Println("Error:", err)
				pc.logger.Error(err.Error(), functionCallerInfo.CachePurhcaseControllerPutCart, "Parse ModifiedAt")
				return err
			}

			_sellerId := cachedProductValue["SellerId"]

			cachedProducts = append(cachedProducts, response.ProductItemDTO{
				ProductId:        item.ProductId,
				Name:             cachedProductValue["Name"],
				Category:         cachedProductValue["Category"],
				Qty:              _qty,
				Price:            _price,
				SKU:              cachedProductValue["SKU"],
				FileID:           cachedProductValue["FileID"],
				FileURI:          cachedProductValue["FileURI"],
				FileThumbnailURI: cachedProductValue["FileThumbnailURI"],
				CreatedAt:        _createdAt,
				UpdatedAt:        _modifiedAt,
				SellerId:         _sellerId,
			})

			// todo; distinct seller_id
			// yang jelas, pertama kali ini pasti enggak ada catatannya
			indexSellerId, exists := mapSellerId[_sellerId]
			if !exists {
				// key tidak ada dalam map, inisiasi total harga awal

				// tambahkan terlebih dulu informasi seller
				mapSellerId[_sellerId] = pointerIndexSellerId

				// append baru ke larik total harga, index dari urutan seller
				sellerIdTotalPrices = append(sellerIdTotalPrices, _price)

				// todo; check cache untuk sellerID
				/*
					type SellerBankDetailDTO struct {
						SellerId          string
						BankAccountName   string  `json:"bankAccountName"`
						BankAccountHolder string  `json:"bankAccountHolder"`
						BankAccountNumber string  `json:"bankAccountNumber"`
						TotalPrice        float64 `json:"totalPrice"`
					}
				*/
				if cachedSellerValue, found := serviceCache.GetAsMap(fmt.Sprintf(serviceCache.CacheSellerById, _sellerId)); found {
					cachedSeller = append(cachedSeller, response.SellerBankDetailDTO{
						SellerId:          cachedSellerValue["SellerId"],
						BankAccountName:   cachedSellerValue["BankAccountName"],
						BankAccountHolder: cachedSellerValue["BankAccountHolder"],
						BankAccountNumber: cachedSellerValue["BankAccountNumber"],
						TotalPrice:        sellerIdTotalPrices[pointerIndexSellerId],
					})
				} else {
					// todo; kompilasi sellerId yang akan diambil
					toGetSellersById = append(toGetSellersById, _sellerId)
				}

				// increment pointerIndexSellerId
				pointerIndexSellerId = pointerIndexSellerId + 1
			} else {
				// key ada,

				// update akumulasi total harga
				sellerIdTotalPrices[indexSellerId] = sellerIdTotalPrices[indexSellerId] + _price

				// update total harga dari cache seller detail, karena cachedSeller berupa larik kita ambil elemen dengan index yang sudah diambil dari mapSellerId
				cachedSeller[indexSellerId].TotalPrice = sellerIdTotalPrices[indexSellerId]
			}

			// todo; insert cart
			cart.TotalPrice = cart.TotalPrice + _price

		} else {
			// todo; kompilasi produkId yang akan diambil
			toGetProductsById = append(toGetProductsById, item.ProductId)

		}
	}

	// todo; get produkId di produk service
	if len(toGetProductsById) > 0 {
		// kirim batch

		// dapat response
		// masukan ke respons.dto
		// untuk sementara pakai toGetProductsById terlebih dulu sampai produk service siap
		for _, item := range toGetProductsById {
			_sellerId := "d29bdf02-25dc-4742-b425-728374370351"
			_price := 0.0
			cart.PurchasedItems = append(cart.PurchasedItems, response.ProductItemDTO{
				ProductId: item,
				Price:     _price,
				SellerId:  _sellerId,
			})

			// todo; distinct seller_id
			// yang jelas, pertama kali ini pasti enggak ada catatannya
			indexSellerId, exists := mapSellerId[_sellerId]
			if !exists {
				// key tidak ada dalam map, inisiasi total harga awal

				// tambahkan terlebih dulu informasi seller
				mapSellerId[_sellerId] = pointerIndexSellerId

				// append baru ke larik total harga, index dari urutan seller
				sellerIdTotalPrices = append(sellerIdTotalPrices, _price)

				// todo; check cache untuk sellerID
				/*
					type SellerBankDetailDTO struct {
						SellerId          string
						BankAccountName   string  `json:"bankAccountName"`
						BankAccountHolder string  `json:"bankAccountHolder"`
						BankAccountNumber string  `json:"bankAccountNumber"`
						TotalPrice        float64 `json:"totalPrice"`
					}
				*/
				if cachedSellerValue, found := serviceCache.GetAsMap(fmt.Sprintf(serviceCache.CacheSellerById, _sellerId)); found {
					cachedSeller = append(cachedSeller, response.SellerBankDetailDTO{
						SellerId:          cachedSellerValue["SellerId"],
						BankAccountName:   cachedSellerValue["BankAccountName"],
						BankAccountHolder: cachedSellerValue["BankAccountHolder"],
						BankAccountNumber: cachedSellerValue["BankAccountNumber"],
						TotalPrice:        sellerIdTotalPrices[pointerIndexSellerId],
					})
				} else {
					// todo; kompilasi sellerId yang akan diambil
					toGetSellersById = append(toGetSellersById, _sellerId)
				}

				// increment pointerIndexSellerId
				pointerIndexSellerId = pointerIndexSellerId + 1
			} else {
				// key ada,

				// update akumulasi total harga
				sellerIdTotalPrices[indexSellerId] = sellerIdTotalPrices[indexSellerId] + _price

				// update total harga dari cache seller detail, karena cachedSeller berupa larik kita ambil elemen dengan index yang sudah diambil dari mapSellerId
				cachedSeller[indexSellerId].TotalPrice = sellerIdTotalPrices[indexSellerId]
			}

			// todo; insert cart
			cart.TotalPrice = cart.TotalPrice + _price
		}
	}
	// todo; get SellerId berdasarkan produkId
	if len(toGetSellersById) > 0 {
		// Mulai pengukuran waktu
		start := time.Now()

		// kirim batch
		connection, err := grpc.Dial(pc.userServiceAddress, grpc.WithInsecure())
		if err != nil {
			pc.logger.Error(err.Error(), functionCallerInfo.PurhcaseControllerPutCart, "Grpc Connection")
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to connect to gRPC server: %s\n", err.Error()))
		}
		defer connection.Close()

		// Create new userService client
		userServiceClient := user.NewUserServiceClient(connection)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call gRPC method
		// dapat response
		grpcResponse, err := userServiceClient.GetUserDetailsWithId(ctx, &user.UserRequest{UserIds: toGetSellersById})
		if err != nil {
			pc.logger.Error(err.Error(), functionCallerInfo.PurhcaseControllerPutCart, "Grpc Call")
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("error during gRPC call: %s\n", err.Error()))
		}
		// masukan ke respons.dto
		for _, item := range grpcResponse.Users {
			cart.PaymentDetails = append(cart.PaymentDetails, response.SellerBankDetailDTO{
				SellerId:          item.UserId,
				BankAccountName:   item.BankAccountName,
				BankAccountHolder: item.BankAccountHolder,
				BankAccountNumber: item.BankAccountNumber,
			})
		}
		connection.Close()

		// Catat waktu selesai
		elapsed := time.Since(start)
		fmt.Printf("GRPC CALL>> Batch processing took %s", elapsed)
	}
	cleanStart := time.Now()
	// todo; compile respond

	// todo; free temporaries
	mapSellerId = nil
	pointerIndexSellerId = 0
	sellerIdTotalPrices = nil
	cachedProducts = nil
	toGetSellersById = nil
	runtime.GC()

	// Catat waktu selesai
	cleanElapsed := time.Since(cleanStart)
	fmt.Printf("CLEAN FREE>> Batch processing took %s", cleanElapsed)
	//
	return c.Status(fiber.StatusOK).JSON(cart)
}

package userRepository

import (
	"context"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/repository"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/response"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

type UserRepositoryInterface interface {
	CreateUserByEmail(ctx context.Context, pool *pgxpool.Pool, email, passwordHash string) (userId string, err error)
	CreateUserByPhone(ctx context.Context, pool *pgxpool.Pool, phone, passwordHash string) (userId string, err error)
	GetAuthByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (auth repository.AuthByEmail, err error)
	GetAuthByPhone(ctx context.Context, pool *pgxpool.Pool, phone string) (auth repository.AuthByPhone, err error)
	UpdateEmail(ctx context.Context, pool *pgxpool.Pool, email, userId string) (user *repository.User, err error)
	UpdatePhone(ctx context.Context, pool *pgxpool.Pool, phone, userId string) (user *repository.User, err error)
	GetUserProfile(ctx context.Context, pool *pgxpool.Pool, userId string) (user *repository.User, err error)
	UpdateUserProfile(ctx context.Context, pool *pgxpool.Pool, input repository.UpdateUser, userId string) (*repository.User, error)

	GetUserProfiles(ctx context.Context, pool *pgxpool.Pool, userId []string) (user []response.UserResponse, err error)
	GetUserProfilesWithId(ctx context.Context, pool *pgxpool.Pool, userIds []string) (user []response.UserWithIdResponse, err error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

func NewUserRepositoryInject(i do.Injector) (UserRepositoryInterface, error) {
	return NewUserRepository(
		do.MustInvoke[*pgxpool.Pool](i),
	), nil
}

func (ur *UserRepository) CreateUserByEmail(ctx context.Context, pool *pgxpool.Pool, email, passwordHash string) (userId string, err error) {
	query := `INSERT INTO users(email, password_hash) VALUES($1, $2) RETURNING id`

	row := pool.QueryRow(ctx, query, email, passwordHash)
	err = row.Scan(&userId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (ur *UserRepository) CreateUserByPhone(ctx context.Context, pool *pgxpool.Pool, phone, passwordHash string) (userId string, err error) {
	query := `INSERT INTO users(phone, password_hash) VALUES($1, $2) RETURNING id`

	row := pool.QueryRow(ctx, query, phone, passwordHash)
	err = row.Scan(&userId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (ur *UserRepository) GetAuthByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (auth repository.AuthByEmail, err error) {
	query := `SELECT id, password_hash, phone FROM users WHERE email = $1 ;`

	row := pool.QueryRow(ctx, query, email)

	var phone *string
	err = row.Scan(&auth.UserId, &auth.HashPassword, &phone)
	if err != nil {
		return repository.AuthByEmail{}, err
	}

	if phone != nil {
		auth.Phone = *phone
	}

	return auth, nil
}

func (ur *UserRepository) GetAuthByPhone(ctx context.Context, pool *pgxpool.Pool, phone string) (auth repository.AuthByPhone, err error) {
	query := `SELECT id, password_hash, email FROM users WHERE phone = $1 ;`

	row := pool.QueryRow(ctx, query, phone)

	var email *string
	err = row.Scan(&auth.UserId, &auth.HashPassword, &email)
	if err != nil {
		return repository.AuthByPhone{}, err
	}

	if email != nil {
		auth.Email = *email
	}

	return auth, nil
}

func (ur *UserRepository) UpdateEmail(ctx context.Context, pool *pgxpool.Pool, email, userId string) (*repository.User, error) {
	query := `
		UPDATE users 
		SET email = $1 
		WHERE id = $2 
		RETURNING
			phone,
			fileId,
			fileUri,
			fileThumbnailUri,
			bankAccountName,
			bankAccountHolder,
			bankAccountNumber;`

	var user repository.User
	err := pool.QueryRow(ctx, query, email, userId).Scan(
		&user.Phone,
		&user.FileId,
		&user.FileUri,
		&user.FileThumbnailUri,
		&user.BankAccountName,
		&user.BankAccountHolder,
		&user.BankAccountNumber,
	)
	if err != nil {
		return &repository.User{}, err
	}

	return &user, nil
}

func (ur *UserRepository) UpdatePhone(ctx context.Context, pool *pgxpool.Pool, phone, userId string) (*repository.User, error) {
	query := `
		UPDATE users 
		SET phone = $1 
		WHERE id = $2 
		RETURNING
			email,
			fileId,
			fileUri,
			fileThumbnailUri,
			bankAccountName,
			bankAccountHolder,
			bankAccountNumber;`

	var user repository.User
	err := pool.QueryRow(ctx, query, phone, userId).Scan(
		&user.Email,
		&user.FileId,
		&user.FileUri,
		&user.FileThumbnailUri,
		&user.BankAccountName,
		&user.BankAccountHolder,
		&user.BankAccountNumber,
	)
	if err != nil {
		return &repository.User{}, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserProfile(ctx context.Context, pool *pgxpool.Pool, userId string) (*repository.User, error) {
	query := `
		SELECT 
			email,
			phone,
			fileId,
			fileUri,
			fileThumbnailUri,
			bankAccountName,
			bankAccountHolder,
			bankAccountNumber
		FROM users 
		WHERE id = $1;`

	var user repository.User
	err := pool.QueryRow(ctx, query, userId).Scan(
		&user.Email,
		&user.Phone,
		&user.FileId,
		&user.FileUri,
		&user.FileThumbnailUri,
		&user.BankAccountName,
		&user.BankAccountHolder,
		&user.BankAccountNumber,
	)
	if err != nil {
		return &repository.User{}, err
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUserProfile(ctx context.Context, pool *pgxpool.Pool, input repository.UpdateUser, userId string) (*repository.User, error) {
	// Target query:
	// `UPDATE users
	// 	SET
	// 		bankAccountName = $1,
	// 		bankAccountHolder = $2,
	// 		bankAccountNumber = $3,
	// 		fileId = $5,
	// 		fileUri = $6,
	// 		fileThumbnailUri = $7,
	// 	WHERE id = $4
	// 	RETURNING
	// 		email,
	// 		phone,
	//		fileId,
	// 		fileUri,
	// 		fileThumbnailUri;`

	query := `
		UPDATE users 
		SET
			bankAccountName = $1,
			bankAccountHolder = $2,
			bankAccountNumber = $3`

	args := make([]interface{}, 4)
	args[0] = input.BankAccountName
	args[1] = input.BankAccountHolder
	args[2] = input.BankAccountNumber
	args[3] = userId

	if input.FileId != nil {
		query += `,
			fileId = $5,
			fileUri = $6,
			fileThumbnailUri = $7`
		args = append(args, input.FileId, input.FileUri, input.FileThumbnailUri)
	}

	query += `
		WHERE id = $4
		RETURNING
			email,
			phone,
			fileId,
			fileUri,
			fileThumbnailUri;`

	user := &repository.User{
		BankAccountName:   input.BankAccountName,
		BankAccountHolder: input.BankAccountHolder,
		BankAccountNumber: input.BankAccountNumber,
	}
	err := pool.QueryRow(
		ctx,
		query,
		args...,
	).Scan(&user.Email, &user.Phone, &user.FileId, &user.FileUri, &user.FileThumbnailUri)
	if err != nil {
		return &repository.User{}, err
	}

	return user, nil
}

// Ambil semua profile user dengan array of userId
//
// Returns:
//   - List<response.UserResponse>
//   - error
func (ur *UserRepository) GetUserProfiles(ctx context.Context, pool *pgxpool.Pool, userIds []string) ([]response.UserResponse, error) {
	var result []response.UserResponse

	// Menyusun query SQL untuk mengambil profil user berdasarkan userId
	query := `
		SELECT 
			id,
			email,
			phone,
			fileId,
			fileUri,
			fileThumbnailUri,
			bankAccountName,
			bankAccountHolder,
			bankAccountNumber
		FROM users 
		WHERE id = ANY($1::text[]);`

	// Menjalankan query untuk mendapatkan hasil berdasarkan userIds
	rows, err := pool.Query(ctx, query, userIds)
	if err != nil {
		// Mengembalikan error jika query gagal dijalankan
		return nil, err
	}
	// Pastikan rows ditutup setelah selesai digunakan
	defer rows.Close()

	// Loop untuk membaca setiap baris hasil query
	for rows.Next() {
		var user response.UserWithIdResponseSqlNulString
		// Memindahkan hasil query ke dalam struct user
		err := rows.Scan(
			&user.UserId,
			&user.Email,
			&user.Phone,
			&user.FileId,
			&user.FileUri,
			&user.FileThumbnailUri,
			&user.BankAccountName,
			&user.BankAccountHolder,
			&user.BankAccountNumber,
		)
		if err != nil {
			// Mengembalikan error jika terjadi kesalahan saat memindahkan data
			return nil, err
		}
		// Menambahkan user ke dalam hasil
		_res := response.UserResponse{
			Email:             user.Email.String,
			Phone:             user.Phone.String,
			FileId:            user.FileId.String,
			FileUri:           user.FileUri.String,
			FileThumbnailUri:  user.FileThumbnailUri.String,
			BankAccountName:   user.BankAccountName.String,
			BankAccountHolder: user.BankAccountHolder.String,
			BankAccountNumber: user.BankAccountNumber.String,
		}
		result = append(result, _res)
	}

	// Memeriksa apakah ada error selama iterasi
	if err = rows.Err(); err != nil {
		// Mengembalikan error jika ada masalah saat membaca hasil query
		return nil, err
	}
	// Mengembalikan hasil yang sudah lengkap
	return result, nil
}

// Ambil semua profile user dengan array of userId dengan userId yang ngikut
//
// Returns:
//   - List<response.UserWithIdResponse>
//   - error
func (ur *UserRepository) GetUserProfilesWithId(ctx context.Context, pool *pgxpool.Pool, userIds []string) ([]response.UserWithIdResponse, error) {
	var result []response.UserWithIdResponse

	// Menyusun query SQL untuk mengambil profil user berdasarkan userId
	query := `
		SELECT 
			id,
			email,
			phone,
			fileId,
			fileUri,
			fileThumbnailUri,
			bankAccountName,
			bankAccountHolder,
			bankAccountNumber
		FROM users 
		WHERE id = ANY($1::text[]);`

	// Menjalankan query untuk mendapatkan hasil berdasarkan userIds
	rows, err := pool.Query(ctx, query, userIds)
	if err != nil {
		// Mengembalikan error jika query gagal dijalankan
		return nil, err
	}
	// Pastikan rows ditutup setelah selesai digunakan
	defer rows.Close()

	// Loop untuk membaca setiap baris hasil query
	for rows.Next() {
		var user response.UserWithIdResponseSqlNulString
		// Memindahkan hasil query ke dalam struct user
		err := rows.Scan(
			&user.UserId,
			&user.Email,
			&user.Phone,
			&user.FileId,
			&user.FileUri,
			&user.FileThumbnailUri,
			&user.BankAccountName,
			&user.BankAccountHolder,
			&user.BankAccountNumber,
		)
		if err != nil {
			// Mengembalikan error jika terjadi kesalahan saat memindahkan data
			return nil, err
		}
		// Menambahkan user ke dalam hasil
		_res := response.UserWithIdResponse{
			UserId:            user.UserId.String,
			Email:             user.Email.String,
			Phone:             user.Phone.String,
			FileId:            user.FileId.String,
			FileUri:           user.FileUri.String,
			FileThumbnailUri:  user.FileThumbnailUri.String,
			BankAccountName:   user.BankAccountName.String,
			BankAccountHolder: user.BankAccountHolder.String,
			BankAccountNumber: user.BankAccountNumber.String,
		}
		result = append(result, _res)
	}

	// Memeriksa apakah ada error selama iterasi
	if err = rows.Err(); err != nil {
		// Mengembalikan error jika ada masalah saat membaca hasil query
		return nil, err
	}
	// Mengembalikan hasil yang sudah lengkap
	return result, nil
}

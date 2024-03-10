package user

import (
	"21-api/helper"
	"21-api/middlewares"
	"21-api/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	Model model.UserModel
}

func (us *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input RegisterRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(input)

		if err != nil {
			// var message = []string{}
			// for _, val := range err.(validator.ValidationErrors) {
			// 	if val.Tag() == "required" {
			// 		message = append(message, fmt.Sprint(val.Field(), " wajib diisi"))
			// 	} else if val.Tag() == "min" {
			// 		message = append(message, fmt.Sprint(val.Field(), " minimal 10 digit"))
			// 	} else {
			// 		message = append(message, fmt.Sprint(val.Field(), " ", val.Tag()))
			// 	}
			// }
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirim kurang sesuai", nil))
		}

		var processInput model.User
		processInput.Hp = input.Hp
		processInput.Nama = input.Nama
		processInput.Password = input.Password

		//Untuk memeriksa apakah nomoor sudah terdaftar
		if us.Model.CekUser(processInput.Hp) {
			return c.JSON(http.StatusConflict, helper.ResponseFormat(http.StatusConflict, "Nomor sudah terdaftar", nil))
		}

		err = us.Model.AddUser(processInput) // ini adalah fungsi yang kita buat sendiri
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "selamat data sudah terdaftar", nil))
	}
}

func (us *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(input)

		if err != nil {
			for _, val := range err.(validator.ValidationErrors) {
				fmt.Println(val.Error())
			}
		}

		result, err := us.Model.Login(input.Hp, input.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		token, err := middlewares.GenerateJWT(result.Hp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem, gagal memproses data", nil))
		}

		var responseData LoginResponse
		responseData.Hp = result.Hp
		responseData.Nama = result.Nama
		responseData.Token = token

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "selamat anda berhasil login", responseData))

	}
}

func (us *UserController) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var hp = c.Param("hp")
		var input model.User
		err := c.Bind(&input)
		if err != nil {
			log.Println("masalah baca input:", err.Error())
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		isFound := us.Model.CekUser(hp)

		if !isFound {
			return c.JSON(http.StatusNotFound,
				helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
		}

		err = us.Model.Update(hp, input)

		if err != nil {
			log.Println("masalah database :", err.Error())
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan saat update data", nil))
		}

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "data berhasil di update", nil))
	}
}

func (us *UserController) ListUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		listUser, err := us.Model.GetAllUser()
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", listUser))
	}
}

func (us *UserController) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var hp = c.Param("hp")
		result, err := us.Model.GetProfile(hp)

		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound,
					helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
			}
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
	}
}

// Fungsi Controller untuk menambah AddActivity
func (us *UserController) AddActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mengambil token JWT dari konteks
		token := c.Get("user").(*jwt.Token)

		// Mengambil claims dari token JWT
		claims := token.Claims.(jwt.MapClaims)

		// Mengambil nomor HP dari claims
		hp := claims["hp"].(string)

		var newActivity model.Activity
		if err := c.Bind(&newActivity); err != nil {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Gagal memproses permintaan", nil))
		}

		if err := us.Model.AddActivity(hp, newActivity); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Gagal menambahkan", nil))
		}

		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "Berhasil Ditambahkan", nil))
	}
}

// Fungsi Controller untuk mengubah kegiatan milik pengguna yang sedang login
func (us *UserController) UpdateActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token JWT dari konteks
		token := c.Get("user").(*jwt.Token)

		// Mendapatkan klaim dari token JWT
		claims := token.Claims.(jwt.MapClaims)

		// Mendapatkan nomor HP dari klaim
		hp := claims["hp"].(string)

		// Mendapatkan ID kegiatan yang ingin diubah dari parameter URL
		activityID := c.Param("id")

		// Konversi activityID menjadi tipe data uint
		activityIDUint, err := strconv.ParseUint(activityID, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "ID kegiatan tidak valid", nil))
		}

		// Mendapatkan kegiatan dari database
		var activity model.Activity
		if err := us.Model.Connection.Where("id = ? AND hp = ?", activityIDUint, hp).First(&activity).Error; err != nil {
			return c.JSON(http.StatusNotFound, helper.ResponseFormat(http.StatusNotFound, "Kegiatan tidak ditemukan", nil))
		}

		// Binding data kegiatan baru yang ingin diubah
		var updatedActivity model.Activity
		if err := c.Bind(&updatedActivity); err != nil {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Gagal memproses permintaan", nil))
		}

		// Memastikan kegiatan yang ingin diubah dimiliki oleh pengguna yang sedang login
		if activity.UserID != activity.UserID {
			return c.JSON(http.StatusForbidden, helper.ResponseFormat(http.StatusForbidden, "Anda tidak memiliki izin untuk mengubah kegiatan ini", nil))
		}

		// Mengupdate kegiatan dalam database
		if err := us.Model.Connection.Model(&activity).Updates(&updatedActivity).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Gagal mengupdate kegiatan", nil))
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Kegiatan berhasil diubah", nil))
	}
}

// GetAllActivities mengembalikan daftar kegiatan berdasarkan pengguna yang terautentikasi
func (us *UserController) GetAllActivities() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token JWT dari konteks
		token := c.Get("user").(*jwt.Token)

		// Mendapatkan klaim dari token JWT
		claims := token.Claims.(jwt.MapClaims)

		// Mendapatkan nomor HP pengguna dari klaim
		hp := claims["hp"].(string)

		// Mengambil daftar kegiatan berdasarkan nomor HP pengguna
		activities, err := us.Model.GetActivitiesByHp(hp)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil daftar kegiatan")
		}

		// Membuat respons
		response := helper.ResponseFormat(http.StatusOK, "Daftar kegiatan berhasil diambil", activities)

		// Mengembalikan respons
		return c.JSON(http.StatusOK, response)
	}
}

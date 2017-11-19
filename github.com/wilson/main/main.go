package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"html"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
var db *sql.DB
var err error
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var router = mux.NewRouter()
var tpl *template.Template

type email struct {
	Name                 string
	Email                string
	NewPasswordGenerated string
}

func addAdmin(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)

	if req.Method != "POST" {
		if adminEmail != "" {
			var adminName string
			var adminRole string
			var newId int

			rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&adminName, &adminRole)
				if err != nil {
					log.Fatal(err)
				}
			}

			rowsNewId, err := db.Query("SELECT max(id_admin)+1 as new_id FROM admin")
			if err != nil {
				log.Fatal(err)
			}
			defer rowsNewId.Close()
			for rowsNewId.Next() {
				err := rowsNewId.Scan(&newId)
				if err != nil {
					log.Fatal(err)
				}
			}

			varmap := map[string]interface{}{
				"AdminEmail": adminEmail,
				"AdminName":  adminName,
				"AdminRole":  adminRole,
				"NewId":      newId,
			}
			tpl.ExecuteTemplate(res, "addAdmin.gohtml", varmap)
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	}
}

func addUser(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)

	if req.Method != "POST" {
		if adminEmail != "" {
			var adminName string
			var adminRole string
			var groupStatus string
			var newId int
			var groupStatusArr []string

			rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&adminName, &adminRole)
				if err != nil {
					log.Fatal(err)
				}
			}

			rowsNewId, err := db.Query("SELECT max(id_pengguna)+1 as new_id FROM pengguna")
			if err != nil {
				log.Fatal(err)
			}
			defer rowsNewId.Close()
			for rowsNewId.Next() {
				err := rowsNewId.Scan(&newId)
				if err != nil {
					log.Fatal(err)
				}
			}

			rowsGroups, err := db.Query("SELECT nama_grup_pengguna FROM grup_pengguna")
			if err != nil {
				log.Fatal(err)
			}
			defer rowsGroups.Close()
			for rowsGroups.Next() {
				err := rowsGroups.Scan(&groupStatus)
				groupStatusArr = append(groupStatusArr, groupStatus)
				if err != nil {
					log.Fatal(err)
				}
			}

			varmap := map[string]interface{}{
				"AdminEmail":     adminEmail,
				"AdminName":      adminName,
				"AdminRole":      adminRole,
				"GroupStatusArr": groupStatusArr,
				"NewId":          newId,
			}
			tpl.ExecuteTemplate(res, "addUser.gohtml", varmap)
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	}
}

func adminLists(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)

	if req.Method != "POST" {
		if adminEmail != "" {
			var adminName string
			var adminRole string

			rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&adminName, &adminRole)
				if err != nil {
					log.Fatal(err)
				}
			}

			varmap := map[string]interface{}{
				"AdminEmail": adminEmail,
				"AdminName":  adminName,
				"AdminRole":  adminRole,
			}
			tpl.ExecuteTemplate(res, "adminLists.gohtml", varmap)
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	} else if req.Method == "POST" {
		var adminName string
		var adminRole string
		var duplicateEmail int
		var newId int

		rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&adminName, &adminRole)
			if err != nil {
				log.Fatal(err)
			}
		}

		rowsNewId, err := db.Query("SELECT max(id_admin)+1 as new_id FROM admin")
		if err != nil {
			log.Fatal(err)
		}
		defer rowsNewId.Close()
		for rowsNewId.Next() {
			err := rowsNewId.Scan(&newId)
			if err != nil {
				log.Fatal(err)
			}
		}

		id_admin := html.EscapeString(req.FormValue("admin_id"))
		admin_name := html.EscapeString(req.FormValue("admin_name"))
		admin_email := html.EscapeString(req.FormValue("admin_email"))
		peran := html.EscapeString(req.FormValue("peran"))
		admin_password := html.EscapeString(req.FormValue("admin_password"))
		admin_password_confirm := html.EscapeString(req.FormValue("admin_password_confirm"))

		rowsDuplicateEmail, err := db.Query("SELECT COUNT(email_admin) FROM admin WHERE email_admin=?", admin_email)
		if err != nil {
			log.Fatal(err)
		}
		defer rowsDuplicateEmail.Close()
		for rowsDuplicateEmail.Next() {
			err := rowsDuplicateEmail.Scan(&duplicateEmail)
			if err != nil {
				log.Fatal(err)
			}
		}

		if len(admin_password) < 6 {
			varmap := map[string]interface{}{
				"AdminEmail":   adminEmail,
				"AdminName":    adminName,
				"AdminRole":    adminRole,
				"ErrorMessage": "Panjang kata sandi minimal 6 karakter",
				"NewId":        newId,
			}
			tpl.ExecuteTemplate(res, "addAdmin.gohtml", varmap)
		} else if len(admin_password_confirm) < 6 {
			varmap := map[string]interface{}{
				"AdminEmail":   adminEmail,
				"AdminName":    adminName,
				"AdminRole":    adminRole,
				"ErrorMessage": "Panjang konfirmasi kata sandi minimal 6 karakter",
				"NewId":        newId,
			}
			tpl.ExecuteTemplate(res, "addAdmin.gohtml", varmap)
		} else if admin_password != admin_password_confirm {
			varmap := map[string]interface{}{
				"AdminEmail":   adminEmail,
				"AdminName":    adminName,
				"AdminRole":    adminRole,
				"ErrorMessage": "Konfirmasi kata sandi tidak cocok",
				"NewId":        newId,
			}
			tpl.ExecuteTemplate(res, "addAdmin.gohtml", varmap)
		} else if duplicateEmail > 0 {
			varmap := map[string]interface{}{
				"AdminEmail":   adminEmail,
				"AdminName":    adminName,
				"AdminRole":    adminRole,
				"ErrorMessage": "Email telah terdaftar",
				"NewId":        newId,
			}
			tpl.ExecuteTemplate(res, "addAdmin.gohtml", varmap)
		} else {
			t := time.Now().Local().Add(time.Hour * time.Duration(7))
			hashed, _ := HashPassword(admin_password)
			rows := db.QueryRow("INSERT INTO admin VALUES(?,?,?,?,?,?)", id_admin, admin_name, admin_email, peran, hashed, t)

			if rows != nil {
				varmap := map[string]interface{}{
					"AdminEmail": adminEmail,
					"AdminName":  adminName,
					"AdminRole":  adminRole,
				}
				tpl.ExecuteTemplate(res, "adminLists.gohtml", varmap)
			}
		}
	}
}

func adminPassword(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)

	if req.Method != "POST" {
		if adminEmail != "" {
			var adminName string
			var adminRole string

			rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&adminName, &adminRole)
				if err != nil {
					log.Fatal(err)
				}
			}

			varmap := map[string]interface{}{
				"AdminEmail": adminEmail,
				"AdminName":  adminName,
				"AdminRole":  adminRole,
			}
			tpl.ExecuteTemplate(res, "adminPassword.gohtml", varmap)
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	} else if req.Method == "POST" {
		if adminEmail != "" {
			var adminName string
			var adminRole string

			rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&adminName, &adminRole)
				if err != nil {
					log.Fatal(err)
				}
			}

			// Grab the adminname & email from the submitted post form
			oldPassword := html.EscapeString(req.FormValue("old-password"))
			newPassword := html.EscapeString(req.FormValue("new-password"))
			newPasswordConfirm := html.EscapeString(req.FormValue("new-password-confirm"))

			var id int
			var databasePassword string

			rowsAdmin, err := db.Query("SELECT id_admin, kata_sandi_admin FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rowsAdmin.Close()
			for rowsAdmin.Next() {
				err := rowsAdmin.Scan(&id, &databasePassword)
				if err != nil {
					log.Fatal(err)
				}
			}

			err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(oldPassword))
			if err != nil {
				// Old password incorrect
				varmap := map[string]interface{}{
					"AdminEmail":   adminEmail,
					"AdminName":    adminName,
					"AdminRole":    adminRole,
					"ErrorMessage": "Kata sandi lama salah",
				}
				tpl.ExecuteTemplate(res, "adminPassword.gohtml", varmap)
			} else if len(oldPassword) < 6 {
				// The password length less than 6 characters
				varmap := map[string]interface{}{
					"AdminEmail":   adminEmail,
					"AdminName":    adminName,
					"AdminRole":    adminRole,
					"ErrorMessage": "Panjang kata sandi lama minimal 6 karakter",
				}
				tpl.ExecuteTemplate(res, "adminPassword.gohtml", varmap)
			} else if len(newPassword) < 6 {
				// The password length less than 6 characters
				varmap := map[string]interface{}{
					"AdminEmail":   adminEmail,
					"AdminName":    adminName,
					"AdminRole":    adminRole,
					"ErrorMessage": "Panjang kata sandi baru minimal 6 karakter",
				}
				tpl.ExecuteTemplate(res, "adminPassword.gohtml", varmap)
			} else if len(newPasswordConfirm) < 6 {
				// The password length less than 6 characters
				varmap := map[string]interface{}{
					"AdminEmail":   adminEmail,
					"AdminName":    adminName,
					"AdminRole":    adminRole,
					"ErrorMessage": "Panjang konfirmasi kata sandi baru minimal 6 karakter",
				}
				tpl.ExecuteTemplate(res, "adminPassword.gohtml", varmap)
			} else if newPassword != newPasswordConfirm {
				// New password and new password confirmation doesn't match
				varmap := map[string]interface{}{
					"AdminEmail":   adminEmail,
					"AdminName":    adminName,
					"AdminRole":    adminRole,
					"ErrorMessage": "Kata sandi baru dan konfirmasi kata sandi baru tidak cocok",
				}
				tpl.ExecuteTemplate(res, "adminPassword.gohtml", varmap)
			} else {
				t := time.Now().Local().Add(time.Hour * time.Duration(7))
				hashed, _ := HashPassword(newPassword)
				resUpdateProfile := db.QueryRow("UPDATE admin SET kata_sandi_admin=?, tanggal_diperbaharui=? WHERE id_admin=?", hashed, t, id)

				if resUpdateProfile != nil {
					varmap := map[string]interface{}{
						"AdminEmail":     adminEmail,
						"AdminName":      adminName,
						"AdminRole":      adminRole,
						"SuccessMessage": "Kata sandi berhasil diperbaharui",
					}
					tpl.ExecuteTemplate(res, "adminPassword.gohtml", varmap)
				}
			}
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	}
}

func adminProfile(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)

	if req.Method != "POST" {
		if adminEmail != "" {
			var adminName string
			var adminRole string

			rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&adminName, &adminRole)
				if err != nil {
					log.Fatal(err)
				}
			}

			varmap := map[string]interface{}{
				"AdminEmail": adminEmail,
				"AdminName":  adminName,
				"AdminRole":  adminRole,
			}
			tpl.ExecuteTemplate(res, "adminProfile.gohtml", varmap)
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	} else if req.Method == "POST" {
		if adminEmail != "" {
			// Grab the adminname & email from the submitted post form
			name := html.EscapeString(req.FormValue("name"))
			status := html.EscapeString(req.FormValue("status"))

			var id int

			rows, err := db.Query("SELECT id_admin FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&id)
				if err != nil {
					log.Fatal(err)
				}
			}

			t := time.Now().Local().Add(time.Hour * time.Duration(7))
			resUpdateProfile := db.QueryRow("UPDATE admin SET nama_admin=?, peran=?, tanggal_diperbaharui=? WHERE id_admin=?", name, status, t, id)
			if resUpdateProfile != nil {
				varmap := map[string]interface{}{
					"AdminEmail":     adminEmail,
					"AdminName":      name,
					"AdminRole":      status,
					"SuccessMessage": "Profile berhasil diperbaharui",
				}
				tpl.ExecuteTemplate(res, "adminProfile.gohtml", varmap)

				return
			}
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	}
}

func clearSession(res http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(res, cookie)
}

func deleteUser(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)
	user_id := html.EscapeString(req.FormValue("id_pengguna"))

	var adminName string
	var adminRole string

	rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&adminName, &adminRole)
		if err != nil {
			log.Fatal(err)
		}
	}

	db.QueryRow("DELETE FROM pengguna WHERE id_pengguna=?", user_id)
	varmap := map[string]interface{}{
		"AdminEmail": adminEmail,
		"AdminName":  adminName,
		"AdminRole":  adminRole,
	}
	tpl.ExecuteTemplate(res, "userLists.gohtml", varmap)
}

func editUserLists(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)
	vars := mux.Vars(req)

	var adminName string
	var adminRole string
	var duplicateEmail int
	var groupAkses string
	var groupAksesArr []string
	var groupStatus string
	var groupStatusArr []string
	var userDesc string
	var userEmail string
	var userCreatedAt string
	var userName string
	var userStatus string
	var userUpdatedAt string

	user_id := vars["id"]

	// Get the logged in name and role
	rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&adminName, &adminRole)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Get user details
	rowsUser, err := db.Query("SELECT nama_pengguna, email_pengguna, deskripsi, p.tanggal_dibuat, p.tanggal_diperbaharui, gp.nama_grup_pengguna as 'status' FROM pengguna p JOIN grup_pengguna gp ON(gp.id_grup_pengguna=p.id_grup_pengguna) WHERE id_pengguna=?", user_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rowsUser.Close()
	for rowsUser.Next() {
		err := rowsUser.Scan(&userName, &userEmail, &userDesc, &userCreatedAt, &userUpdatedAt, &userStatus)
		if err != nil {
			log.Fatal(err)
		}
	}

	rowsGroups, err := db.Query("SELECT nama_grup_pengguna FROM grup_pengguna")
	if err != nil {
		log.Fatal(err)
	}
	defer rowsGroups.Close()
	for rowsGroups.Next() {
		err := rowsGroups.Scan(&groupStatus)
		groupStatusArr = append(groupStatusArr, groupStatus)
		if err != nil {
			log.Fatal(err)
		}
	}

	rowsGroupAkses, err := db.Query("SELECT DISTINCT(nama_akses) FROM grup g JOIN grup_pengguna gp ON(gp.id_grup_pengguna=g.id_grup_pengguna) JOIN grup_akses ga ON(ga.id_grup_akses=g.id_grup_akses) JOIN akses a ON(a.id_grup_akses=g.id_grup_akses) WHERE nama_grup_pengguna=? ORDER BY nama_grup_pengguna", userStatus)
	if err != nil {
		log.Fatal(err)
	}
	defer rowsGroupAkses.Close()
	for rowsGroupAkses.Next() {
		err := rowsGroupAkses.Scan(&groupAkses)
		groupAksesArr = append(groupAksesArr, groupAkses)
		if err != nil {
			log.Fatal(err)
		}
	}

	if req.Method != "POST" {
		if adminEmail != "" {
			varmap := map[string]interface{}{
				"AdminEmail":     adminEmail,
				"AdminName":      adminName,
				"AdminRole":      adminRole,
				"GroupAksesArr":  groupAksesArr,
				"GroupStatusArr": groupStatusArr,
				"UserId":         user_id,
				"UserName":       userName,
				"UserEmail":      userEmail,
				"UserDesc":       userDesc,
				"UserCreatedAt":  userCreatedAt,
				"UserUpdatedAt":  userUpdatedAt,
				"UserStatus":     userStatus,
			}
			tpl.ExecuteTemplate(res, "editUser.gohtml", varmap)
			return
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	} else if req.Method == "POST" {
		if adminEmail != "" {
			deskripsi_pengguna := html.EscapeString(req.FormValue("deskripsi_pengguna"))
			email_pengguna := html.EscapeString(req.FormValue("email_pengguna"))
			nama_pengguna := html.EscapeString(req.FormValue("nama_pengguna"))
			user_status := html.EscapeString(req.FormValue("user_status"))

			rowsDuplicateEmail, err := db.Query("SELECT COUNT(email_pengguna) FROM pengguna WHERE email_pengguna=? AND id_pengguna!=?", email_pengguna, user_id)
			if err != nil {
				log.Fatal(err)
			}
			defer rowsDuplicateEmail.Close()
			for rowsDuplicateEmail.Next() {
				err := rowsDuplicateEmail.Scan(&duplicateEmail)
				if err != nil {
					log.Fatal(err)
				}
			}

			if user_status == "" {
				varmap := map[string]interface{}{
					"AdminEmail":     adminEmail,
					"AdminName":      adminName,
					"AdminRole":      adminRole,
					"ErrorMessage":   "Status belum diisi",
					"GroupStatusArr": groupStatusArr,
					"UserId":         user_id,
					"UserName":       userName,
					"UserEmail":      userEmail,
					"UserDesc":       userDesc,
					"UserCreatedAt":  userCreatedAt,
					"UserUpdatedAt":  userUpdatedAt,
					"UserStatus":     userStatus,
				}
				tpl.ExecuteTemplate(res, "editUser.gohtml", varmap)
			} else if duplicateEmail > 0 {
				varmap := map[string]interface{}{
					"AdminEmail":     adminEmail,
					"AdminName":      adminName,
					"AdminRole":      adminRole,
					"ErrorMessage":   "Email telah terdaftar",
					"GroupStatusArr": groupStatusArr,
					"UserId":         user_id,
					"UserName":       userName,
					"UserEmail":      userEmail,
					"UserDesc":       userDesc,
					"UserCreatedAt":  userCreatedAt,
					"UserUpdatedAt":  userUpdatedAt,
					"UserStatus":     userStatus,
				}
				tpl.ExecuteTemplate(res, "editUser.gohtml", varmap)
			} else {
				var userStatusId int

				rowsUserStatusId, err := db.Query("SELECT id_grup_pengguna FROM grup_pengguna WHERE nama_grup_pengguna=?", user_status)
				if err != nil {
					log.Fatal(err)
				}
				defer rowsUserStatusId.Close()
				for rowsUserStatusId.Next() {
					err := rowsUserStatusId.Scan(&userStatusId)
					if err != nil {
						log.Fatal(err)
					}
				}

				t := time.Now().Local().Add(time.Hour * time.Duration(7))

				resUpdateUser := db.QueryRow("UPDATE pengguna SET nama_pengguna=?, email_pengguna=?, deskripsi=?, tanggal_diperbaharui=?, id_grup_pengguna=? WHERE id_pengguna=?", nama_pengguna, email_pengguna, deskripsi_pengguna, t, userStatusId, user_id)
				if resUpdateUser != nil {
					var userStatusName string

					rowsUserStatus, err := db.Query("SELECT nama_grup_pengguna FROM grup_pengguna gp JOIN pengguna p ON(p.id_grup_pengguna=gp.id_grup_pengguna) WHERE id_pengguna=?", user_id)
					if err != nil {
						log.Fatal(err)
					}
					defer rowsUserStatus.Close()
					for rowsUserStatus.Next() {
						err := rowsUserStatus.Scan(&userStatusName)
						if err != nil {
							log.Fatal(err)
						}
					}
					varmap := map[string]interface{}{
						"AdminEmail":     adminEmail,
						"AdminName":      adminName,
						"AdminRole":      adminRole,
						"GroupStatusArr": groupStatusArr,
						"SuccessMessage": "Pengguna berhasil diperbaharui",
						"UserId":         user_id,
						"UserName":       nama_pengguna,
						"UserEmail":      email_pengguna,
						"UserDesc":       deskripsi_pengguna,
						"UserCreatedAt":  userCreatedAt,
						"UserUpdatedAt":  t,
						"UserStatus":     userStatusName,
					}
					tpl.ExecuteTemplate(res, "editUser.gohtml", varmap)

					return
				}
			}
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	}
}

func forgotPassword(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		tpl.ExecuteTemplate(res, "forgotPassword.gohtml", nil)
		return
	} else if req.Method == "POST" {
		// Grab the email from the submitted post form
		email := html.EscapeString(req.FormValue("email"))

		var validEmail int

		rows, err := db.Query("SELECT COUNT(email_admin) FROM admin WHERE email_admin=?", email)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&validEmail)
			if err != nil {
				log.Fatal(err)
			}
		}

		if validEmail == 0 {
			varmap := map[string]interface{}{
				"ErrorMessage": "Email belum terdaftar",
			}

			tpl.ExecuteTemplate(res, "forgotPassword.gohtml", varmap)
		} else {
			var adminName string

			rowsAdminName, err := db.Query("SELECT nama_admin FROM admin WHERE email_admin=?", email)
			if err != nil {
				log.Fatal(err)
			}
			defer rowsAdminName.Close()
			for rowsAdminName.Next() {
				err := rowsAdminName.Scan(&adminName)
				if err != nil {
					log.Fatal(err)
				}
			}

			newPasswordGenerated := RandStringRunes(52)
			send(adminName, email, newPasswordGenerated, email)

			t := time.Now().Local().Add(time.Hour * time.Duration(7))
			hashed, _ := HashPassword(newPasswordGenerated)
			resUpdatePassword := db.QueryRow("UPDATE admin SET kata_sandi_admin=?, tanggal_diperbaharui=? WHERE email_admin=?", hashed, t, email)

			if resUpdatePassword != nil {
				varmap := map[string]interface{}{
					"SuccessMessage": "Instruksi pengaturan ulang kata sandi telah dikirim ke email Anda",
				}
				tpl.ExecuteTemplate(res, "forgotPassword.gohtml", varmap)
			}
		}
	}
}

func getAdminEmail(req *http.Request) (adminEmail string) {
	if cookie, err := req.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			adminEmail = cookieValue["email"]
		}
	}
	return adminEmail
}

func getAdminsJson(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(getJSON("SELECT id_admin, nama_admin, email_admin, peran, tanggal_diperbaharui FROM admin")))
}

func getJSON(sqlString string) string {
	rows, err := db.Query(sqlString)
	if err != nil {
		return ""
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return ""
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return ""
	}
	// fmt.Println(string(jsonData))

	return string(jsonData)
}

func getGroupsJson(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(getJSON("SELECT * FROM grup_pengguna")))
}

func getUsersJson(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(getJSON("SELECT id_pengguna, nama_pengguna, email_pengguna, nama_admin, nama_grup_pengguna, pengguna.tanggal_diperbaharui as 'tanggal_diperbaharui' FROM pengguna JOIN admin ON (admin.id_admin=pengguna.id_admin) JOIN grup_pengguna ON (grup_pengguna.id_grup_pengguna=pengguna.id_grup_pengguna)")))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func home(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)

	if req.Method != "POST" {
		if adminEmail != "" {
			var adminName string
			var adminRole string

			// Get the logged in name and role
			rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&adminName, &adminRole)
				if err != nil {
					log.Fatal(err)
				}
			}

			varmap := map[string]interface{}{
				"AdminEmail": adminEmail,
				"AdminName":  adminName,
				"AdminRole":  adminRole,
			}
			tpl.ExecuteTemplate(res, "home.gohtml", varmap)
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		tpl.ExecuteTemplate(res, "index.gohtml", nil)
		return
	}

	// Grab the email & password from the submitted post form
	email := html.EscapeString(req.FormValue("email"))
	password := html.EscapeString(req.FormValue("password"))

	// Grab from the database
	var databaseEmail string
	var databasePassword string

	err := db.QueryRow("SELECT email_admin, kata_sandi_admin FROM admin WHERE email_admin=?", email).Scan(&databaseEmail, &databasePassword)

	// If not then redirect to the login page
	if err != nil {
		varmap := map[string]interface{}{
			"ErrorMessage": "Akun tidak ditemukan",
		}
		tpl.ExecuteTemplate(res, "index.gohtml", varmap)
		return
	}

	// Validate the password
	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	// If wrong password, send error message and redirect to the login page
	if err != nil {
		varmap := map[string]interface{}{
			"ErrorMessage": "Kata sandi salah",
		}
		tpl.ExecuteTemplate(res, "index.gohtml", varmap)
		return
	} else {
		// If the login succeeded
		setSession(email, res)
		http.Redirect(res, req, "/home", 302)
		return
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func logout(res http.ResponseWriter, req *http.Request) {
	clearSession(res)
	http.Redirect(res, req, "/", 302)
}

func main() {
	// Create an sql.DB and check for errors
	db, err = sql.Open("mysql", "root:root@/kudo_admin")
	if err != nil {
		panic(err.Error())
	}
	// sql.DB should be long lived "defer" closes it once this function ends
	defer db.Close()

	// Test the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/addAdmin", addAdmin)
	r.HandleFunc("/addUser", addUser)
	r.HandleFunc("/adminLists", adminLists)
	r.HandleFunc("/deleteUser", deleteUser)
	r.HandleFunc("/forgotPassword", forgotPassword)
	r.HandleFunc("/getAdminsJson", getAdminsJson)
	r.HandleFunc("/getGroupsJson", getGroupsJson)
	r.HandleFunc("/getUsersJson", getUsersJson)
	r.HandleFunc("/home", home)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/adminPassword", adminPassword)
	r.HandleFunc("/adminProfile", adminProfile)
	r.HandleFunc("/userGroup", userGroup)
	r.HandleFunc("/userLists", userLists)
	r.HandleFunc("/userLists/{id}", editUserLists)
	r.HandleFunc("/userPermission", userPermission)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.ListenAndServe(":8080", r)
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func send(adminname string, adminemail string, passwordgenerated string, to string) {
	// Create the variables for the template
	myEmail := email{adminname, adminemail, passwordgenerated}

	// Create a template using template.html
	tmpl, err := template.New("email").ParseFiles("templates/email.gohtml")
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	// Stores the parsed template
	var buff bytes.Buffer

	// Send the parsed template to buff
	err = tmpl.Execute(&buff, myEmail)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	from := "practicalkudo@gmail.com"
	password := "cfuibgcablwwevfh"

	// replace body with buff.String()
	msg := "From: " + "Kudo Admin" + "\r\n" +
		"To: " + to + "\r\n" +
		"MIME-Version: 1.0" + "\r\n" +
		"Content-type: text/html" + "\r\n" +
		"Subject: Pengaturan ulang kata sandi" + "\r\n\r\n" +
		buff.String() + "\r\n"

	errMail := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if errMail != nil {
		log.Printf("Error: %s", errMail)
		return
	}

	// log.Print("Message sent!")
}

func setSession(adminEmail string, res http.ResponseWriter) {
	value := map[string]string{
		"email": adminEmail,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(res, cookie)
	}
}

func userGroup(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)

	if req.Method != "POST" {
		if adminEmail != "" {
			var adminName string
			var adminRole string

			rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&adminName, &adminRole)
				if err != nil {
					log.Fatal(err)
				}
			}

			varmap := map[string]interface{}{
				"AdminEmail": adminEmail,
				"AdminName":  adminName,
				"AdminRole":  adminRole,
			}
			tpl.ExecuteTemplate(res, "userGroup.gohtml", varmap)
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	}
}

func userLists(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)

	if req.Method != "POST" {
		if adminEmail != "" {
			var adminName string
			var adminRole string

			rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&adminName, &adminRole)
				if err != nil {
					log.Fatal(err)
				}
			}

			varmap := map[string]interface{}{
				"AdminEmail": adminEmail,
				"AdminName":  adminName,
				"AdminRole":  adminRole,
			}
			tpl.ExecuteTemplate(res, "userLists.gohtml", varmap)
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	} else if req.Method == "POST" {
		var userStatusId string
		var admin_id string
		var adminName string
		var adminRole string
		var duplicateEmail int
		var newId int
		var groupStatus string
		var groupStatusArr []string

		rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&adminName, &adminRole)
			if err != nil {
				log.Fatal(err)
			}
		}

		rowsGroups, err := db.Query("SELECT nama_grup_pengguna FROM grup_pengguna")
		if err != nil {
			log.Fatal(err)
		}
		defer rowsGroups.Close()
		for rowsGroups.Next() {
			err := rowsGroups.Scan(&groupStatus)
			groupStatusArr = append(groupStatusArr, groupStatus)
			if err != nil {
				log.Fatal(err)
			}
		}

		rowsNewId, err := db.Query("SELECT max(id_pengguna)+1 as new_id FROM pengguna")
		if err != nil {
			log.Fatal(err)
		}
		defer rowsNewId.Close()
		for rowsNewId.Next() {
			err := rowsNewId.Scan(&newId)
			if err != nil {
				log.Fatal(err)
			}
		}

		user_id := html.EscapeString(req.FormValue("user_id"))
		user_name := html.EscapeString(req.FormValue("user_name"))
		user_email := html.EscapeString(req.FormValue("user_email"))
		user_description := html.EscapeString(req.FormValue("user_description"))
		user_status := html.EscapeString(req.FormValue("user_status"))

		rowsDuplicateEmail, err := db.Query("SELECT COUNT(email_pengguna) FROM pengguna WHERE email_pengguna=?", user_email)
		if err != nil {
			log.Fatal(err)
		}
		defer rowsDuplicateEmail.Close()
		for rowsDuplicateEmail.Next() {
			err := rowsDuplicateEmail.Scan(&duplicateEmail)
			if err != nil {
				log.Fatal(err)
			}
		}

		if user_status == "" {
			varmap := map[string]interface{}{
				"AdminEmail":     adminEmail,
				"AdminName":      adminName,
				"AdminRole":      adminRole,
				"ErrorMessage":   "Status belum diisi",
				"GroupStatusArr": groupStatusArr,
				"NewId":          newId,
			}
			tpl.ExecuteTemplate(res, "addUser.gohtml", varmap)
		} else if duplicateEmail > 0 {
			varmap := map[string]interface{}{
				"AdminEmail":     adminEmail,
				"AdminName":      adminName,
				"AdminRole":      adminRole,
				"ErrorMessage":   "Email telah terdaftar",
				"GroupStatusArr": groupStatusArr,
				"NewId":          newId,
			}
			tpl.ExecuteTemplate(res, "addUser.gohtml", varmap)
		} else {
			rowsUserStatusId, err := db.Query("SELECT id_grup_pengguna FROM grup_pengguna WHERE nama_grup_pengguna=?", user_status)
			if err != nil {
				log.Fatal(err)
			}
			defer rowsUserStatusId.Close()
			for rowsUserStatusId.Next() {
				err := rowsUserStatusId.Scan(&userStatusId)
				if err != nil {
					log.Fatal(err)
				}
			}

			rowsAdminId, err := db.Query("SELECT id_admin FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rowsAdminId.Close()
			for rowsAdminId.Next() {
				err := rowsAdminId.Scan(&admin_id)
				if err != nil {
					log.Fatal(err)
				}
			}

			t := time.Now().Local().Add(time.Hour * time.Duration(7))
			rows := db.QueryRow("INSERT INTO pengguna VALUES(?,?,?,?,?,?,?,?)", user_id, user_name, user_email, user_description, t, t, admin_id, userStatusId)

			if rows != nil {
				varmap := map[string]interface{}{
					"AdminEmail": adminEmail,
					"AdminName":  adminName,
					"AdminRole":  adminRole,
				}
				tpl.ExecuteTemplate(res, "userLists.gohtml", varmap)
			}
		}
	}
}

func userPermission(res http.ResponseWriter, req *http.Request) {
	adminEmail := getAdminEmail(req)

	if req.Method != "POST" {
		if adminEmail != "" {
			var adminName string
			var adminRole string

			rows, err := db.Query("SELECT nama_admin, peran FROM admin WHERE email_admin=?", adminEmail)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&adminName, &adminRole)
				if err != nil {
					log.Fatal(err)
				}
			}

			varmap := map[string]interface{}{
				"AdminEmail": adminEmail,
				"AdminName":  adminName,
				"AdminRole":  adminRole,
			}
			tpl.ExecuteTemplate(res, "userPermission.gohtml", varmap)
		} else {
			// Redirect to login page if admin not authenticated
			http.Redirect(res, req, "/", 301)
		}
	}
}

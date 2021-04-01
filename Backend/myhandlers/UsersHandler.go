package myhandlers

import (
	"Backend/config"
	"Backend/mymodels"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// CreateUser endpoint creates in db an user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser mymodels.User
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "(USER) %v", err.Error())
		return
	}

	Db, err := config.MYSQLConnection()
	defer Db.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
	json.Unmarshal(reqBody, &newUser)
	switch {
	case (newUser.FirstLastName == nil) || (len(*newUser.FirstLastName) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "FirstLastName is empty or not valid")
		return
	case (newUser.SecondLastName == nil) || (len(*newUser.SecondLastName) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "SecondLastName is empty or not valid")
		return
	case (newUser.FirstName == nil) || (len(*newUser.FirstName) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "FirstName is empty or not valid")
		return
	case (newUser.OtherNames == nil):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "OtherNames is empty or not valid")
		return
	case (newUser.CountryID == nil) || (*newUser.CountryID*1 <= 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "CountryID is empty or not valid")
		return
	case (newUser.DocumentTypeID == nil) || (*newUser.DocumentTypeID*1 <= 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "DocumentTypeID is empty or not valid")
		return
	case (newUser.Document == nil) || (len(*newUser.Document) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Document is empty or not valid")
		return
	case (newUser.StartDate == nil) || (len(*newUser.StartDate) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "StartDate is empty or not valid")
		return
	case (newUser.AreaID == nil) || (*newUser.AreaID*1 <= 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "AreaID is empty or not valid")
		return
	case (newUser.Status == nil) || (len(*newUser.Status) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Status is empty or not valid")
		return
	default:
		rows, err := Db.Query("CALL CreateUser(?,?,?,?,?,?,?,?,?,?)", newUser.FirstLastName, newUser.SecondLastName, newUser.FirstName, newUser.OtherNames, newUser.CountryID, newUser.DocumentTypeID, newUser.Document, newUser.StartDate, newUser.AreaID, newUser.Status)
		defer rows.Close()
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			fmt.Println("(SQL) ", err.Error())
			return
		}
		for rows.Next() {
			var userID int
			var result string
			if err := rows.Scan(&userID, &result); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("(SQL) ", err.Error())
				return
			}
			newUser.ID = &userID
			newUser.Email = &result
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newUser)

		}
		return
	}
}

// GetAllUsers endpoint returns all users in db
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users mymodels.AllUsers
	var Db, err = config.MYSQLConnection()
	defer Db.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
	rows, err := Db.Query("CALL GetAllUsers()")
	defer rows.Close()
	if err != nil {
		fmt.Println("-> ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "(SQL) %v", err.Error())
		return
	}
	for rows.Next() {
		var userID, countryID, documentTypeID, areaID int
		var firstLastName, secondLastName, firstName, otherNames, country, documentType, document, email, startDate, status, area, registrationDate string
		if err := rows.Scan(&userID, &firstLastName, &secondLastName, &firstName, &otherNames, &countryID, &country, &documentTypeID, &documentType, &document, &email, &startDate, &status, &areaID, &area, &registrationDate); err != nil {
			fmt.Fprintf(w, "(SQL) %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Next line is used in production because it has
		// RegistrationDate attribute
		//var user = mymodels.User{ID: &userID, FirstLastName: &firstLastName, SecondLastName: &secondLastName, FirstName: &firstName, OtherNames: &otherNames, CountryID: &countryID, Country: &country, IdTypes: &idTypes, IdTypesName: &idTypesName, IdentificationDocument: &identificationDocument, Email: &email, StartDate: &startDate, AreaID: &areaID, AreaName: &areaName, Status: &status, RegistrationDate: &registrationDate}

		// Next line is used in development because RegistrationDate
		// attribute is always changing.
		var user = mymodels.User{ID: &userID, FirstLastName: &firstLastName, SecondLastName: &secondLastName, FirstName: &firstName, OtherNames: &otherNames, CountryID: &countryID, Country: &country, DocumentTypeID: &documentTypeID, DocumentType: &documentType, Document: &document, Email: &email, StartDate: &startDate, AreaID: &areaID, Area: &area, Status: &status}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
	return
}

// GetUser returns one user filtered by UserID
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID is empty or not valid")
		return
	}
	userID, err := strconv.Atoi(keys[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID is empty or not valid")
		return
	}
	Db, err := config.MYSQLConnection()
	defer Db.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	row, err := Db.Query("CALL GetUser(?)", userID)
	defer row.Close()
	if err != nil {
		fmt.Println("-> ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "(SQL) %v", err.Error())
		return
	}
	for row.Next() {
		var userID, countryID, documentTypeID, areaID int
		var firstLastName, secondLastName, firstName, otherNames, country, documentType, document, email, startDate, status, area, registrationDate string
		if err := row.Scan(&userID, &firstLastName, &secondLastName, &firstName, &otherNames, &countryID, &country, &documentTypeID, &documentType, &document, &email, &startDate, &status, &areaID, &area, &registrationDate); err != nil {
			fmt.Fprintf(w, "(SQL) %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Next line is used in production because it has
		// RegistrationDate attribute
		//var user = mymodels.User{ID: &userID, FirstLastName: &firstLastName, SecondLastName: &secondLastName, FirstName: &firstName, OtherNames: &otherNames, CountryID: &countryID, Country: &country, IdTypes: &idTypes, IdTypesName: &idTypesName, IdentificationDocument: &identificationDocument, Email: &email, StartDate: &startDate, AreaID: &areaID, AreaName: &areaName, Status: &status, RegistrationDate: &registrationDate}

		// Next line is used in development because RegistrationDate
		// attribute is always changing.
		var user = mymodels.User{ID: &userID, FirstLastName: &firstLastName, SecondLastName: &secondLastName, FirstName: &firstName, OtherNames: &otherNames, CountryID: &countryID, Country: &country, DocumentTypeID: &documentTypeID, DocumentType: &documentType, Document: &document, Email: &email, StartDate: &startDate, AreaID: &areaID, Area: &area, Status: &status}
		json.NewEncoder(w).Encode(user)
	}
	return
}

// UpdateUser updates all user's data except RegisterDate
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Something went wrong")
		fmt.Fprintf(w, "(USER) %v", err.Error())
		return
	}
	var updatedUser mymodels.User
	json.Unmarshal(reqBody, &updatedUser)
	switch {
	case (updatedUser.ID == nil) || (*updatedUser.ID*1 <= 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID is empty or not valid")
		return
	case (updatedUser.FirstLastName == nil) || (len(*updatedUser.FirstLastName) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "FirstLastName is empty or not valid")
		return
	case (updatedUser.SecondLastName == nil) || (len(*updatedUser.SecondLastName) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "SecondLastName is empty or not valid")
		return
	case (updatedUser.FirstName == nil) || (len(*updatedUser.FirstName) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "FirstName is empty or not valid")
		return
	case (updatedUser.OtherNames == nil):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "OtherNames is empty or not valid")
		return
	case (updatedUser.CountryID == nil) || (*updatedUser.CountryID*1 <= 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "CountryID is empty or not valid")
		return
	case (updatedUser.DocumentTypeID == nil) || (*updatedUser.DocumentTypeID*1 <= 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "DocumentTypeID is empty or not valid")
		return
	case (updatedUser.Document == nil) || (len(*updatedUser.Document) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Document is empty or not valid")
		return
	case (updatedUser.StartDate == nil) || (len(*updatedUser.StartDate) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "StartDate is empty or not valid")
		return
	case (updatedUser.AreaID == nil) || (*updatedUser.AreaID*1 <= 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "AreaID is empty or not valid")
		return
	case (updatedUser.Status == nil) || (len(*updatedUser.Status) == 0):
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Status is empty or not valid")
		return
	default:
		Db, err := config.MYSQLConnection()
		defer Db.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
		}
		row, err := Db.Query("SELECT UpdateUser(?,?,?,?,?,?,?,?,?,?,?)", updatedUser.ID, updatedUser.FirstLastName, updatedUser.SecondLastName, updatedUser.FirstName, updatedUser.OtherNames, updatedUser.CountryID, updatedUser.DocumentTypeID, updatedUser.Document, updatedUser.StartDate, updatedUser.AreaID, updatedUser.Status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "(SQL) %v", err.Error())
			return
		}

		for row.Next() {
			var result string
			if err := row.Scan(&result); err != nil {
				fmt.Println("(SQL) ", err.Error())
				return
			}
			if result == "Not updated" {
				w.WriteHeader(http.StatusInternalServerError)
			}
			updatedUser.Email = &result
			json.NewEncoder(w).Encode(updatedUser)
		}
		return
	}
}

// DeleteUser deletes one user by UserID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "(USER) %v", err.Error())
		return
	}
	var deletedUser mymodels.User
	json.Unmarshal(reqBody, &deletedUser)

	if (deletedUser.ID) == nil || (*deletedUser.ID*1 == 0) || (*deletedUser.ID*1 < 0) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID is empty or not valid")
		return
	}

	Db, err := config.MYSQLConnection()
	defer Db.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
	row, err := Db.Exec("DELETE FROM Users WHERE ID=?", deletedUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "(SQL) %v", err.Error())
		return
	}

	count, err := row.RowsAffected()
	if err != nil {
		fmt.Fprintf(w, "(SQL) %v", err.Error())
		return
	}
	if count == 1 {
		fmt.Fprintf(w, "One row deleted")
	} else {
		fmt.Fprintf(w, "No rows deleted")
	}
	return
}

package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/shivendr0102/Assignment_02_OrderServies/model"
)

type SortOrder int

const (
	SortIDAsc SortOrder = iota
	SortIDDesc
	SortSTATUSAsc
	SortSTATUSDesc
	SortTOTALAsc
	SortTOTALDesc
	SortCURRUNITAsc
	SortCURRUNITDesc
)

func (db Database) Add_Orders(orders model.Order) (string, error) {

	// SQL QUERIES Used in this Function
	const (
		getallCountqueryStatement = `SELECT COUNT(*) FROM [Order] WHERE ID = ?;`
		InsertqueryStatement      = `INSERT INTO [Order] (ID, Order_Status, Items, Total, CurrencyUnit ) VALUES (?, ?, ?, ?, ?);`
		UpdatequeryStatement      = `UPDATE INTO [Order] SET Order_Status = ? WHERE ID = ?;`
	)

	// Checking if ORDER_ID is already present in Database or not
	var cnt int
	cntStmt, err := db.SqlDb.Prepare(getallCountqueryStatement)
	if err != nil {
		return "", nil
	}

	errr := cntStmt.QueryRow(orders.ID).Scan(&cnt)
	if errr != nil {
		return "Unsuccessfull", err
	}

	// String Generation of all the present Items ID in HTTPrequest to be clubbed together
	Items_id_String := ""
	for _, ids := range orders.Items {
		Items_id_String = Items_id_String + ids.Id + ","
	}

	// If ORDER_ID not present then Insert Everything into DataBase
	if cnt == 0 {
		stmt, err := db.SqlDb.Prepare(InsertqueryStatement)
		if err != nil {
			return "UnSuccesfull", err
		}
		defer stmt.Close()

		_, err = stmt.Exec(orders.ID, orders.Status, Items_id_String, orders.Total, orders.CurrencyUnit)
		if err != nil {
			return "UnSuccesfull", err
		}
		return "Successfull", nil
	}

	// Else , Update the STATUS in ORDER Table for the matching ORDER_ID
	stmt, err := db.SqlDb.Prepare(UpdatequeryStatement)
	if err != nil {
		return "UnSuccesfull", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(orders.Status, orders.ID)
	if err != nil {
		return "UnSuccesfull", err
	}

	return "Successfull", nil

}

func (db Database) Get_Orders(Sorter SortOrder, search *string) ([]model.Order, error) {

	// SQL QUERIES Used in this Function
	const (
		getOrderIdAscQuery        = `SELECT * FROM [Order] OVER (ORDER BY [ID] ASC) WHERE ( Order_Status LIKE ? OR CurrencyUnit LIKE ?)`
		getOrderIdDescQuery       = `SELECT * FROM [Order] OVER (ORDER BY [ID] DESC) WHERE ( Order_Status LIKE ? OR CurrencyUnit LIKE ?)`
		getOrderStatusAscQuery    = `SELECT * FROM [Order] OVER (ORDER BY [Status] ASC) WHERE ( Order_Status LIKE ? OR CurrencyUnit LIKE ?)`
		getOrderStatusDescQuery   = `SELECT * FROM [Order] OVER (ORDER BY [Status] DESC) WHERE ( Order_Status LIKE ? OR CurrencyUnit LIKE ?)`
		getOrderTotalAscQuery     = `SELECT * FROM [Order] OVER (ORDER BY [Total] ASC) WHERE ( Order_Status LIKE ? OR CurrencyUnit LIKE ?)`
		getOrderTotalDescQuery    = `SELECT * FROM [Order] OVER (ORDER BY [Total] DESC) WHERE ( Order_Status LIKE ? OR CurrencyUnit LIKE ?)`
		getOrderCurrUnitAscQuery  = `SELECT * FROM [Order] OVER (ORDER BY [CurrencyUnit] ASC) WHERE ( Order_Status LIKE ? OR CurrencyUnit LIKE ?)`
		getOrderCurrUnitDescQuery = `SELECT * FROM [Order] OVER (ORDER BY [CurrencyUnit] DESC) WHERE ( Order_Status LIKE ? OR CurrencyUnit LIKE ?)`
		GetItemsByIdQuery         = `SELECT * FROM [Items] WHERE Id = ?`
	)

	vals := make([]interface{}, 0)

	searchFilterVal := "%"
	if search != nil {
		searchFilterVal = fmt.Sprintf("%%%s%%", *search)
	}

	var stmt *sql.Stmt
	switch Sorter {
	case SortIDAsc:
		stmt_prepare, err := db.SqlDb.Prepare(getOrderIdAscQuery)
		if err != nil {
			return nil, err
		}
		stmt = stmt_prepare
	case SortIDDesc:
		stmt_prepare, err := db.SqlDb.Prepare(getOrderIdDescQuery)
		if err != nil {
			return nil, err
		}
		stmt = stmt_prepare
	case SortSTATUSAsc:
		stmt_prepare, err := db.SqlDb.Prepare(getOrderStatusAscQuery)
		if err != nil {
			return nil, err
		}
		stmt = stmt_prepare
	case SortSTATUSDesc:
		stmt_prepare, err := db.SqlDb.Prepare(getOrderStatusDescQuery)
		if err != nil {
			return nil, err
		}
		stmt = stmt_prepare
	case SortTOTALAsc:
		stmt_prepare, err := db.SqlDb.Prepare(getOrderTotalAscQuery)
		if err != nil {
			return nil, err
		}
		stmt = stmt_prepare
	case SortTOTALDesc:
		stmt_prepare, err := db.SqlDb.Prepare(getOrderTotalDescQuery)
		if err != nil {
			return nil, err
		}
		stmt = stmt_prepare
	case SortCURRUNITAsc:
		stmt_prepare, err := db.SqlDb.Prepare(getOrderCurrUnitAscQuery)
		if err != nil {
			return nil, err
		}
		stmt = stmt_prepare
	case SortCURRUNITDesc:
		stmt_prepare, err := db.SqlDb.Prepare(getOrderCurrUnitDescQuery)
		if err != nil {
			return nil, err
		}
		stmt = stmt_prepare
	default:
		stmt_prepare, err := db.SqlDb.Prepare(getOrderIdAscQuery)
		if err != nil {
			return nil, err
		}
		stmt = stmt_prepare
	}

	// Two times adding searcfilter ( 1-> Status , 2-> CurrencyUnit)
	vals = append(vals, searchFilterVal)
	vals = append(vals, searchFilterVal)
	rows, err := stmt.Query(vals...)
	if err != nil {
		return nil, err
	}

	Orders := []model.Order{}

	for rows.Next() {
		var id, status, items, total, currencyUnit *string
		err := rows.Scan(&id, &status, &items, &total, &currencyUnit)
		if err != nil {
			return nil, err
		}

		// Fetching all the Items IDs by splitting the string on ","
		ItemsIDs := strings.Split(*items, ",")

		ItemsList := []model.Items{}

		// Fetching Item Details for each particular Item ID
		for _, ids := range ItemsIDs {
			stmt, err := db.SqlDb.Prepare(GetItemsByIdQuery)
			if err != nil {
				return nil, err
			}
			Itemsdetails := stmt.QueryRow(strings.TrimSpace(ids))

			var item_id, description, price, quantity *string
			item_err := Itemsdetails.Scan(&item_id, &description, &price, &quantity)
			if item_err != nil {
				return nil, err
			}
			ItemsList = append(ItemsList, model.Items{
				Id:          *item_id,
				Description: *description,
				Price:       *price,
				Quantity:    *quantity,
			})
		}

		Orders = append(Orders, model.Order{
			ID:           *id,
			Status:       *status,
			Items:        ItemsList,
			Total:        *total,
			CurrencyUnit: *currencyUnit,
		})
	}
	return Orders, nil
}

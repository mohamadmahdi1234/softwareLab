package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"simpleAPI/config"
	"simpleAPI/db/entity"
	"sync"
	"time"
)

const (
	ADDITEM               = 1
	CHECKHALLISVALID      = 2
	CHECKWITHHALLDATETIME = 3
	CHECKWITHRESERVETABLE = 4
	INSERTRESERVE         = 5
)

type MySQL struct {
	sync.RWMutex
	Master     *sql.DB
	Statements map[uint8]*sql.Stmt
}

func NewMySQL(cfg *config.Config) (*MySQL, error) {
	conn, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&collation=utf8_unicode_ci",
			cfg.DBUsername,
			cfg.DBPassword,
			cfg.DBHost,
			"3306",
			"conference_hall",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("open master db conn: %w", err)
	}
	log.Println("ping master db...")
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping master db conn: %w", err)
	}
	return &MySQL{
		Master:     conn,
		Statements: make(map[uint8]*sql.Stmt),
	}, nil

}

func (db *MySQL) ReserveHall(ctx context.Context, req entity.ReserveReq, date time.Time) error {
	var hallIsValid bool
	hallISValidStmt := db.getStmt(CHECKHALLISVALID)
	if hallISValidStmt == nil {
		ps, err := db.Master.PrepareContext(ctx, `SELECT EXISTS(SELECT 1 FROM halls WHERE id = ?)`)
		if err != nil {
			return err
		}
		db.setStmt(CHECKHALLISVALID, hallISValidStmt)
		hallISValidStmt = ps
	}
	err := hallISValidStmt.QueryRowContext(ctx, req.HallID).Scan(&hallIsValid)
	if err != nil {
		return fmt.Errorf("err in database %w", err)
	}
	if !hallIsValid {
		return fmt.Errorf("hall doesnot exist")
	}

	dateTimeHallIsValidStmt := db.getStmt(CHECKWITHHALLDATETIME)
	if dateTimeHallIsValidStmt == nil {
		ps, err := db.Master.PrepareContext(ctx, `SELECT COUNT(*) FROM hall_date_time WHERE hall_id = ? AND date = ? AND FIND_IN_SET(?, times)`)
		if err != nil {
			return err
		}
		db.setStmt(CHECKWITHHALLDATETIME, dateTimeHallIsValidStmt)
		dateTimeHallIsValidStmt = ps
	}
	var count int
	err = dateTimeHallIsValidStmt.QueryRowContext(ctx, req.HallID, date, req.Time).Scan(&count)
	if err != nil {
		return fmt.Errorf("herr happen %w", err)
	}
	if count == 0 {
		return fmt.Errorf("hall has not this time in its valid times for reserve")
	}

	var conflict int
	reserveIsValidStmt := db.getStmt(CHECKWITHRESERVETABLE)
	if reserveIsValidStmt == nil {
		ps, err := db.Master.PrepareContext(ctx, `SELECT COUNT(*) FROM reserve WHERE  hall_id = ? AND date = ? AND time = ?`)
		if err != nil {
			return err
		}
		db.setStmt(CHECKWITHRESERVETABLE, reserveIsValidStmt)
		reserveIsValidStmt = ps
	}
	err = reserveIsValidStmt.QueryRowContext(ctx, req.HallID, date, req.Time).Scan(&conflict)
	if err != nil {
		return fmt.Errorf("conflict in reserve %w", err)
	}

	if conflict > 0 {
		return fmt.Errorf("conflict in reserve")
	}

	_, err = db.Master.ExecContext(
		ctx,
		`insert into reserve (user_id, hall_id, date, time, status) VALUES (?, ?, ?, ?, ?)`,
		req.UserID,
		req.HallID,
		date,
		req.Time,
		0,
	)
	if err != nil {
		return fmt.Errorf("problem in insert reserve %w", err)
	}
	return nil

}

func (db *MySQL) GetHalls(ctx context.Context) ([]*entity.Hall, error) {
	sqlstmt := `select * from halls`
	rows, err := db.Master.QueryContext(ctx, sqlstmt)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	var todos []*entity.Hall
	for rows.Next() {
		todo := new(entity.Hall)
		if err = db.scanHall(rows, todo); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil

}

func (db *MySQL) scanHall(row *sql.Rows, hall *entity.Hall) error {
	err := row.Scan(
		&hall.ID,
		&hall.Name,
		&hall.Description,
		&hall.Tools,
		&hall.ImageLink,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *MySQL) Close() {
	for _, stmt := range db.Statements {
		_ = stmt.Close()
	}
	_ = db.Master.Close()
}

func (db *MySQL) getStmt(id uint8) *sql.Stmt {
	db.RLock()
	defer db.RUnlock()
	return db.Statements[id]
}

func (db *MySQL) setStmt(id uint8, stmt *sql.Stmt) {
	db.Lock()
	defer db.Unlock()
	db.Statements[id] = stmt
}

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

var db *sql.DB
var storeStmt *sql.Stmt
var noteName string

type notebook struct {
	Id      int64
	PrevId  int64
	content *TextBlock
}

func InitDB() (err error) {
	db, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/brycenote?timeout=30s")
	if err != nil {
		log.Print(err)
		return err
	}
	db.SetConnMaxLifetime(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Minute * 15)
	db.SetMaxOpenConns(10)
	return db.Ping()
}

func storeWithDB(tb *TextBlock) (err error) {
	err = InitDB()
	if err != nil {
		return err
	}
	noteName = tb.Tittle
	_, err = deleteFromDB(noteName)
	if err != nil {
		log.Panicln(err)
	}
	return storeAll(tb, -1)
}

func deleteFromDB(tittle string) (int64, error) {
	result, err := db.Exec("delete from notebook where note_name = ?", tittle)
	if err != nil {
		log.Panicln(err)
	}
	return result.RowsAffected()
}

func LoadFromDB(tittle string) (tb *TextBlock, err error) {
	err = InitDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("select id,prev_id, content from notebook where note_name = ? and prev_id = -1 limit 1", tittle)
	b := getNoteBook(row)
	rows, err := db.Query("select id, prev_id, content from notebook where note_name = ? and prev_id  > -1 ", tittle)
	if err != nil {
		return nil, err
	}
	var notebooks []*notebook
	for rows.Next() {
		notebooks = append(notebooks, getNoteBook(rows))
	}
	setSubTexts(b.content, b.Id, notebooks)
	return b.content, nil
}

func getNoteBook(row any) *notebook {
	type scanner interface {
		Scan(dest ...any) error
	}
	r, ok := row.(scanner)
	if !ok {
		log.Panicln("row must have. [Scan(dest ...any) error] method")
	}
	var content string
	var id, pid int64
	err := r.Scan(&id, &pid, &content)
	if err != nil {
		log.Panicln(err)
	}
	b := new(TextBlock)
	err = json.Unmarshal([]byte(content), b)
	if err != nil {
		log.Panicln(err)
	}
	t := &notebook{Id: id, PrevId: pid, content: b}
	return t
}

func setSubTexts(t *TextBlock, id int64, bs []*notebook) {
	var idm = map[*TextBlock]int64{}
	for _, b := range bs {
		if b.PrevId == id {
			b.content.prev = t
			idm[b.content] = b.Id
			t.subBlocks = append(t.subBlocks, b.content)
		}
	}
	for _, block := range t.subBlocks {
		setSubTexts(block, idm[block], bs)
	}
}

func storeOne(tb *TextBlock, preId int64) (id int64, err error) {
	if storeStmt == nil {
		storeStmt, err = db.Prepare("insert into notebook(prev_id,note_name,content) values (?,?,?)")
	}
	bs, err := json.Marshal(tb)
	if err != nil {
		return -1, err
	}
	bf := &bytes.Buffer{}
	err = json.Compact(bf, bs)
	if err != nil {
		return -1, err
	}
	result, err := storeStmt.Exec(preId, noteName, bf.String())
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func storeAll(tb *TextBlock, preId int64) (err error) {
	id, err := storeOne(tb, preId)
	if err != nil {
		return err
	}
	for _, block := range tb.subBlocks {
		err = storeAll(block, id)
		if err != nil {
			return err
		}
	}
	return nil
}

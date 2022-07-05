package main

import (
	"context"
	"encoding/binary"
	"io"
	"net/http"

	"github.com/jackc/pgx/v4"
)

func load(db *pgx.Conn) (string, [][]byte, error) {
	var source string
	var name string
	err := db.QueryRow(context.Background(), "SELECT name, source FROM sounds").Scan(&name, &source)
	if err != nil {
		return "", nil, err
	}

	res, err := http.Get(source)
	if err != nil {
		return "", nil, err
	}

	var opuslen int16
	var buffer = make([][]byte, 0)

	for {
		err = binary.Read(res.Body, binary.LittleEndian, &opuslen)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return name, buffer, nil
		}

		if err != nil {
			return "", nil, err
		}

		in := make([]byte, opuslen)
		err = binary.Read(res.Body, binary.LittleEndian, &in)

		if err != nil {
			return "", nil, err
		}

		buffer = append(buffer, in)
	}
}

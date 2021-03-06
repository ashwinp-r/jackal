/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package pgsql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/ortuman/jackal/xmpp"
	"github.com/ortuman/jackal/xmpp/jid"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
)

func TestInsertOfflineMessages(t *testing.T) {
	j, _ := jid.NewWithString("ortuman@jackal.im/balcony", false)
	message := xmpp.NewElementName("message")
	message.SetID(uuid.New())
	message.AppendElement(xmpp.NewElementName("body"))
	m, _ := xmpp.NewMessageFromElement(message, j, j)
	messageXML := m.String()

	s, mock := NewMock()
	mock.ExpectExec("INSERT INTO offline_messages (.+)").
		WithArgs("ortuman", messageXML).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.InsertOfflineMessage(m, "ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.Nil(t, err)

	s, mock = NewMock()
	mock.ExpectExec("INSERT INTO offline_messages (.+)").
		WithArgs("ortuman", messageXML).
		WillReturnError(errGeneric)

	err = s.InsertOfflineMessage(m, "ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.NotNil(t, err)
}

func TestCountOfflineMessages(t *testing.T) {
	countColums := []string{"count"}

	s, mock := NewMock()
	mock.ExpectQuery("SELECT COUNT(.+) FROM offline_messages (.+)").
		WithArgs("ortuman").
		WillReturnRows(sqlmock.NewRows(countColums).AddRow(1))

	cnt, _ := s.CountOfflineMessages("ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.Equal(t, 1, cnt)

	s, mock = NewMock()
	mock.ExpectQuery("SELECT COUNT(.+) FROM offline_messages (.+)").
		WithArgs("ortuman").
		WillReturnRows(sqlmock.NewRows(countColums))

	cnt, _ = s.CountOfflineMessages("ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.Equal(t, 0, cnt)

	s, mock = NewMock()
	mock.ExpectQuery("SELECT COUNT(.+) FROM offline_messages (.+)").
		WithArgs("ortuman").
		WillReturnError(errGeneric)

	_, err := s.CountOfflineMessages("ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.Equal(t, errGeneric, err)
}

func TestFetchOfflineMessages(t *testing.T) {
	var offlineMessagesColumns = []string{"data"}

	s, mock := NewMock()
	mock.ExpectQuery("SELECT (.+) FROM offline_messages (.+)").
		WithArgs("ortuman").
		WillReturnRows(sqlmock.NewRows(offlineMessagesColumns).AddRow("<message id='abc'><body>Hi!</body></message>"))

	msgs, _ := s.FetchOfflineMessages("ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.Equal(t, 1, len(msgs))

	s, mock = NewMock()
	mock.ExpectQuery("SELECT (.+) FROM offline_messages (.+)").
		WithArgs("ortuman").
		WillReturnRows(sqlmock.NewRows(offlineMessagesColumns))

	msgs, _ = s.FetchOfflineMessages("ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.Equal(t, 0, len(msgs))

	s, mock = NewMock()
	mock.ExpectQuery("SELECT (.+) FROM offline_messages (.+)").
		WithArgs("ortuman").
		WillReturnRows(sqlmock.NewRows(offlineMessagesColumns).AddRow("<message id='abc'><body>Hi!"))

	_, err := s.FetchOfflineMessages("ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.NotNil(t, err)

	s, mock = NewMock()
	mock.ExpectQuery("SELECT (.+) FROM offline_messages (.+)").
		WithArgs("ortuman").
		WillReturnError(errGeneric)

	_, err = s.FetchOfflineMessages("ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.Equal(t, errGeneric, err)
}

func TestDeleteOfflineMessages(t *testing.T) {
	s, mock := NewMock()
	mock.ExpectExec("DELETE FROM offline_messages (.+)").
		WithArgs("ortuman").WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.DeleteOfflineMessages("ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.Nil(t, err)

	s, mock = NewMock()
	mock.ExpectExec("DELETE FROM offline_messages (.+)").
		WithArgs("ortuman").WillReturnError(errGeneric)

	err = s.DeleteOfflineMessages("ortuman")
	require.Nil(t, mock.ExpectationsWereMet())
	require.Equal(t, errGeneric, err)
}

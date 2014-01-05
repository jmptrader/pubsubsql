/* Copyright (C) 2013 CompleteDB LLC.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with PubSubSQL.  If not, see <http://www.gnu.org/licenses/>.
 */

package pubsubsql

import (
	"encoding/json"
	"net"
)

type Client interface {

	// Connect connects to the pubsubsql server.
	// Address has the form host:port.
	Connect(address string) bool

	// Disconnect disconnects from the pubsubsql server.
	Disconnect()

	// Ok determines if last operation succeeded. 
	Ok() bool

	// Failed determines if last operation failed.
	Failed() bool

	// Error returns error string when last operation fails.
	Error() string

	// Execute executes a command.
	// Returns true on success.
	Execute(command string) bool

	// JSON() returns JSON response string returned from the last operation.
	JSON() string

	// Action returns action for last operation.
	Action() string

	// I returns unique record identifier generated by the database table.
	// Valid for actions: insert, select (when id is in selected columns),
	// pubsub insert, pubsub add, pubsub delete, pubsub update.
	Id() string

	// PubSubId returns unique pubsub identifier generated by the database.  
	// Valid for pubsub related actions: subscribe. insert, add, delete, update.
	// PubSubId is generated by pubsubsql server when client subscribes to a table and
	// used to uniqely identify particual subscription for the connected client.
	PubSubId() string

	// RecordCount returns number of records in the returned result set.
	RecordCount() int

	// NextRecord move cusrsor to the next data record of the returned result set.    
	// Returns false when all records are read.
	// Must be called initially to position cursor to the first record. 
	NextRecord() bool

	// JSONRecord returns current record in JSON format.
	JSONRecord() string

	// Value returns column value by column name.
	// If column does not exist in current result set it returns empty string.	
	Value(column string) string

	// ValueByOrdinal returns column value by column ordinal.
	// If column ordinal does not exist in current result set it returns empty string.	
	ValueByColumnOrdinal(ordinal int) string

	// Columns returns array of valid column names returned by last operation. 		
	Columns() []string

	// ColumnCount returns number of valid columns
	ColumnCount() int

	// WaitForPubSub waits until publish message is retreived or
	// timeout expired.
	// Returns false on timeout.
	WaitForPubSub(timeout int) bool
}

func NewClient() Client {
	var c client
	return &c
}

var CLIENT_DEFAULT_BUFFER_SIZE int = 2048

// respnoseData holds unmarshaled result from pubsubsql JSON response
type responseData struct {
	Status   string
	Msg      string
	Action   string
	Id       string
	PubSubId string
	Rows     int
	Fromrow  int
	Torow    int
	Data     []map[string]string
}

func (this *responseData) reset() {
	this.Status = ""
	this.Msg = ""
	this.Action = ""
	this.Id = ""
	this.Rows = 0
	this.Fromrow = 0
	this.Torow = 0
	this.Data = nil
}

type client struct {
	Client
	rw        NetMessageReaderWriter
	requestId uint32
	err       string
	rawjson   []byte
	//
	response responseData
	record   int
}

func (this *client) Connect(address string) bool {
	this.Disconnect()
	conn, err := net.Dial("tcp", address)
	if err != nil {
		this.setError(err)
		return false
	}
	this.rw.Set(conn, CLIENT_DEFAULT_BUFFER_SIZE)
	return true
}

func (this *client) Disconnect() {
	this.write("close")
	// write may generate error so we reset after instead
	this.reset()
	this.rw.Close()
}

func (this *client) Ok() bool {
	return this.err == ""
}

func (this *client) Failed() bool {
	return !this.Ok()
}

func (this *client) Error() string {
	return this.err
}

func (this *client) Execute(command string) bool {
	this.reset()
	ok := this.write(command)
	var bytes []byte
	var header *NetworkHeader
	for ok {
		header, bytes, ok = this.read()
		if !ok {
			break
		}
		if header.RequestId == this.requestId {
			// response we are waiting for
			return this.unmarshalJSON(bytes)
		} else if header.RequestId == 0 {
			// pubsub action, save it and skip it for now
		} else if header.RequestId < this.requestId {
			// we did not read full result set from previous command ignore it or flag and error?
			// TODO DECIDE
			this.setErrorString("previous result was not fully read")
			ok = false
		} else {
			// this should never happen
			this.setErrorString("protocol error invalid requestId")
			ok = false
		}
	}
	return ok
}

func (this *client) JSON() string {
	return string(this.rawjson)
}

func (this *client) Action() string {
	return this.response.Action
}

func (this *client) Id() string {
	return this.response.Id
}

func (this *client) PubSubId() string {
	return this.response.PubSubId
}

func (this *client) RecordCount() int {
	return this.response.Rows
}

func (this *client) NextRecord() bool {
	for this.Ok() {
		// no result set
		if this.response.Rows == 0 {
			return false
		}
		if this.response.Fromrow == 0 || this.response.Torow == 0 {
			this.setErrorString("protocol error invalid fromrow, torow values")
		}
		// the current record is valid 
		this.record++
		if this.record <= (this.response.Torow - this.response.Fromrow) {
			return true
		}
		// we reached the end of result set
		if this.response.Rows == this.response.Torow {
			return false
		}
		// get another batch 
		this.reset()
		header, bytes, ok := this.read()
		if !ok {
			return false
		}
		// should not happen but check anyway
		if header.RequestId != this.requestId {
			this.setErrorString("protocol error")
			return false
		}
		// we got another batch unmarshall the data	
		this.unmarshalJSON(bytes)
	}
	return false
}

func (this *client) unmarshalJSON(bytes []byte) bool {
	this.rawjson = bytes
	err := json.Unmarshal(bytes, &this.response)
	if err != nil {
		this.setError(err)
		return false
	}
	if this.response.Status != "ok" {
		this.setErrorString(this.response.Msg)
		return false
	}
	return true
}

func (this *client) reset() {
	this.resetError()
	this.response.reset()
	this.rawjson = nil
	this.record = -1
}

func (this *client) resetError() {
	this.err = ""
}

func (this *client) setErrorString(err string) {
	this.reset()
	this.err = err
}

func (this *client) setError(err error) {
	this.setErrorString(err.Error())
}

func (this *client) write(message string) bool {
	this.requestId++
	if this.rw.Valid() {
		err := this.rw.WriteHeaderAndMessage(this.requestId, []byte(message))
		if err == nil {
			return true
		}
		this.setError(err)
		return false
	}
	this.setErrorString("Not connected")
	return false
}

func (this *client) read() (*NetworkHeader, []byte, bool) {
	if this.rw.Valid() {
		header, bytes, err := this.rw.ReadMessage()
		if err == nil {
			return header, bytes, true
		}
		this.setError(err)
		return nil, nil, false
	}
	this.setErrorString("Not connected")
	return nil, nil, false
}

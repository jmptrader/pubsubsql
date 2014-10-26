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

package server

func (this *parser) parseConnectionAddress(connectionAddress *string) request {
	tok := this.tokens.Produce()
	if tok.typ != tokenTypeSqlValue {
		return this.parseError("expected connection address, but got: " + tok.typ.String())
	}
	*connectionAddress = tok.val
	return nil
}

// mysql connect connectionAddress
func (this *parser) parseMysqlConnect() request {
	req := new(mysqlConnectRequest)
	// connectionAddress
	if errReq := this.parseConnectionAddress(&(req.address)); errReq != nil {
		return errReq
	}
	return this.parseEOF(req)
}

// mysql disconnect
func (this *parser) parseMysqlDisconnect() request {
	req := new(mysqlDisconnectRequest)
	return this.parseEOF(req)
}

// mysql status
func (this *parser) parseMysqlStatus() request {
	req := new(mysqlStatusRequest)
	return this.parseEOF(req)
}

// mysql subscribe
func (this *parser) parseMysqlSubscribe() request {
	req := new(mysqlSubscribeRequest)
	r := this.parseSqlSubscribe()
	req.sqlSubscribeReq = r.(*sqlSubscribeRequest)
	return req
}

// mysql unsubscribe
func (this *parser) parseMysqlUnsubscribe() request {
	req := new(mysqlUnsubscribeRequest)
	r := this.parseSqlUnsubscribe()
	req.sqlUnsubscribeReq = r.(*sqlUnsubscribeRequest)
	return req
}

// mysql tables
func (this *parser) parseMysqlTables() request {
	req := new(mysqlTablesRequest)
	return this.parseEOF(req)
}

// mysql
func (this *parser) parseSqlMysql() request {
	tok := this.tokens.Produce()
	switch tok.typ {
	case tokenTypeSqlConnect:
		return this.parseMysqlConnect()
	case tokenTypeSqlDisconnect:
		return this.parseMysqlDisconnect()
	case tokenTypeCmdStatus:
		return this.parseMysqlStatus()
	case tokenTypeSqlSubscribe:
		return this.parseMysqlSubscribe()
	case tokenTypeSqlUnsubscribe:
		return this.parseMysqlUnsubscribe()
	case tokenTypeSqlTables:
		return this.parseMysqlTables()
	}
	return this.parseError("invalid mysql request")
}

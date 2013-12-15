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

import "log"
import "os"
import "fmt"

var debugLogger = log.New(os.Stderr, "debug: ", log.LstdFlags)
var infoLogger = log.New(os.Stderr, "", log.LstdFlags)
var warnLogger = log.New(os.Stderr, "warning: ", log.LstdFlags)
var errLogger = log.New(os.Stderr, "error: ", log.LstdFlags)

func debug(v ...interface{}) {
	if LOG_DEBUG {
		debugLogger.Output(2, fmt.Sprintln(v...))
	}
}

func loginfo(v ...interface{}) {
	if LOG_INFO {
		infoLogger.Output(2, fmt.Sprintln(v...))
	}
}

func logwarn(v ...interface{}) {
	if LOG_WARN {
		warnLogger.Output(2, fmt.Sprintln(v...))
	}
}

func logerror(v ...interface{}) {
	if LOG_INFO {
		infoLogger.Output(2, fmt.Sprintln(v...))
	}
}
/* Copyright (C) 2014 CompleteDB LLC.
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

package pubsubsql; 

import java.util.ArrayList;

public class ResponseData {

	String status;
	String msg;
	String action;
	String pubsubid;
	int rows;
	int fromrow;
	int torow;
	ArrayList<String> columns = new ArrayList<String>();
	ArrayList<ArrayList<String>> data = new ArrayList<ArrayList<String>>();

}

/*
 * Radon
 *
 * Copyright 2018 The Radon Authors.
 * Code is licensed under the GPLv3.
 *
 */

package proxy

import (
	"testing"

	"fakedb"

	"github.com/stretchr/testify/assert"
	"github.com/xelabs/go-mysqlstack/driver"
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/sqltypes"
	"github.com/xelabs/go-mysqlstack/xlog"
)

func TestProxyInsert(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := MockProxy(log)
	defer cleanup()
	address := proxy.Address()

	// fakedbs.
	{
		fakedbs.AddQueryPattern("create .*", &sqltypes.Result{})
		fakedbs.AddQueryPattern("insert .*", &sqltypes.Result{})
	}

	// create database.
	{
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		query := "create database test"
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}

	// create test table.
	{
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		query := "create table test.t1(id int, b int) partition by hash(id)"
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}

	// Delete.
	{
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		query := "insert into test.t1 (id, b) values(1,2),(3,4)"
		fakedbs.AddQuery(query, fakedb.Result3)
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}
}

func TestProxyInsertQuerys(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := MockProxy(log)
	defer cleanup()
	address := proxy.Address()

	// fakedbs.
	{
		fakedbs.AddQueryPattern("create .*", &sqltypes.Result{})
		fakedbs.AddQueryPattern("insert .*", &sqltypes.Result{})
	}

	// create database.
	{
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		query := "create database test"
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}

	tables := []string{
		"create table test.t1(id int, b int) partition by hash(id)",
		"create table test.t2(id datetime, b int) partition by hash(id)",
		"create table test.t3(id varchar(200), b int) partition by hash(id)",
		"create table test.t4(id decimal, b int) partition by hash(id)",
		"create table test.t5(id float, b int) partition by hash(id)",
	}

	querys := []string{
		"insert into test.t1(id, b) values(1, 1)",
		"insert into test.t2(id, b) values(20111218131717, 1)",
		"insert into test.t3(id, b) values('xx', 1)",
		"insert into test.t4(id, b) values(1.11, 1)",
		"insert into test.t5(id, b) values(0.3333, 1)",
	}

	for _, table := range tables {
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		_, err = client.FetchAll(table, -1)
		assert.Nil(t, err)
	}

	for _, query := range querys {
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}
}

func TestProxyLongTimeQuerys(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := MockProxy(log)
	defer cleanup()
	address := proxy.Address()
	proxy.SetLongQueryTime(0)

	// fakedbs.
	{
		fakedbs.AddQueryPattern("create .*", &sqltypes.Result{})
		fakedbs.AddQueryPattern("insert .*", &sqltypes.Result{})
	}

	// create database.
	{
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		query := "create database test"
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}

	tables := []string{
		"create table test.t1(id int, b int) partition by hash(id)",
		"create table test.t2(id datetime, b int) partition by hash(id)",
		"create table test.t3(id varchar(200), b int) partition by hash(id)",
		"create table test.t4(id decimal, b int) partition by hash(id)",
		"create table test.t5(id float, b int) partition by hash(id)",
	}

	for _, table := range tables {
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		_, err = client.FetchAll(table, -1)
		assert.Nil(t, err)
	}

	querysError := []string{
		"insert into test.t6(id, b) values(1, 1)",
		"insert into test.t7(id, b) values(1, 1)",
	}
	for _, query := range querysError {
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		_, err = client.FetchAll(query, -1)
		assert.NotNil(t, err)
	}
}

func TestProxyInsertAutoIncrement(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := MockProxy(log)
	defer cleanup()
	address := proxy.Address()

	// fakedbs.
	{
		fakedbs.AddQueryPattern("create .*", &sqltypes.Result{})
		fakedbs.AddQueryPattern("insert .*", &sqltypes.Result{})
		fakedbs.AddQueryPattern("replace .*", &sqltypes.Result{})
	}

	// create database.
	{
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		query := "create database test"
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}

	tables := []string{
		"create table test.t1(`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT, b int) partition by hash(id)",
	}

	querys := []string{
		"insert into test.t1(b) values(1)",
		"explain insert into test.t1(b) values(1)",
		"insert into test.t1(id, b) values(1, 1)",
		"explain insert into test.t1(id, b) values(1, 1)",
		"replace into test.t1(b) values(1)",
		"explain replace into test.t1(b) values(1)",
		"replace into test.t1(id, b) values(1, 1)",
		"explain replace into test.t1(id, b) values(1, 1)",
	}

	for _, table := range tables {
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		_, err = client.FetchAll(table, -1)
		assert.Nil(t, err)
	}

	for _, query := range querys {
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}
}

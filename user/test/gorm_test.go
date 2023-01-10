package test

import (
	"context"
	"fmt"
	"jam3.com/user/pgk/dao"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	str := strconv.Itoa(rand.Intn(10000000))
	user := &dao.User{
		Username:  "user" + str,
		Fullname:  "user1fullname" + str,
		Birthdata: time.Now().Unix(),
		Bio:       "BioTest" + str,
		Token:     "token" + str,
		Create_at: time.Now(),
		Update_at: time.Now(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	dao.Save(ctx, user)
	println(user)
}
func TestUpdateUser(t *testing.T) {
	str := strconv.Itoa(rand.Intn(10000000))
	user := dao.User{
		Username:  "testupdate" + str,
		Fullname:  "testupdatefullname" + str,
		Birthdata: time.Now().Unix(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	dao.Update(ctx, 2, user)

}
func TestDeleteUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	dao.Delete(ctx, 1)

}

func TestSelectUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	u := dao.GetInfo(ctx, 2)
	fmt.Println(u)
	u = dao.GetInfo(ctx, 3)
	fmt.Println(u)
}

func TestPageUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	u := dao.GetList(ctx, 1, 2)
	fmt.Println(u)
	u = dao.GetList(ctx, 5, 2)
	fmt.Println(u)
}

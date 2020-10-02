package main

import (
	"database/sql"
	"log"
	"time"

	_ "./models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

//Coupon is a coupon
type Coupon struct {
	CouponID           int       `json:"coupon_id"`
	CouponSerialNumber string    `json:"coupon_serial_number"`
	CouponStatus       string    `json:"coupon_status"`
	CouponUpdateTime   time.Time `json:"coupon_update_time"`
}

//Award is a award
type Award struct {
	AwardID           int       `json:"award_id"`
	AwardName         string    `json:"award_name"`
	AwardSerialNumber string    `json:"award_serial_number"`
	AwardStatus       string    `json:"award_status"`
	AwardUpdateTime   time.Time `json:"award_update_time"`
}

var db *sql.DB
var (
	coupon Coupon
	award  Award
)

func main() {

	err := initDB()
	if err != nil {
		log.Println("conncet db error", err)
		return
	}
	log.Println("conncet db success ! ")

	r := gin.Default()

	r.POST("/lottery", func(c *gin.Context) {
		r.Use(cors.Default())
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		serialNumber := c.PostForm("serialNumber")
		log.Println("some user send serialNumber = ", serialNumber)

		couponRows, err := db.Query("select * from coupon where coupon_status = 0")
		if err != nil {
			log.Println("select db error", err)
			return
		}

		isValidCoupon := false
		for couponRows.Next() {
			err = couponRows.Scan(&coupon.CouponID, &coupon.CouponSerialNumber, &coupon.CouponStatus, &coupon.CouponUpdateTime)
			if coupon.CouponSerialNumber == serialNumber {

				stmt, _ := db.Prepare("update coupon set coupon_status= ? ,coupon_update_time = now() where coupon_id= ?")

				stmt.Exec("1", coupon.CouponID)
				log.Println("serialNumber : ", serialNumber, " used update status to 1")
				isValidCoupon = true
				break
			}
			if err != nil {
				log.Print(err.Error())
			}
		}

		if isValidCoupon == true {
			awardNumber := "0"
			awardName := ""

			err := db.QueryRow("SELECT * FROM award where award_status='0' ORDER BY RAND() LIMIT 1").Scan(&award.AwardID, &award.AwardName, &award.AwardSerialNumber, &award.AwardStatus, &award.AwardUpdateTime)
			if err != nil {
				log.Print(err.Error())
			}
			stmt, _ := db.Prepare("update award set award_status= ? ,award_update_time = now() where award_id= ?")
			stmt.Exec("1", award.AwardID)

			awardName = award.AwardName
			awardNumber = award.AwardSerialNumber

			c.JSON(200, gin.H{
				"awardNumber": awardNumber,
				"awardName":   awardName,
			})
		} else {
			c.JSON(200, gin.H{
				"awardNumber": 0,
				"awardName":   "您輸入的序號無效或是已被使用",
			})
			return
		}
	})

	r.Run(":8082")
}

func initDB() (err error) {
	dsn := "frank:123456@tcp(localhost)/lottery?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

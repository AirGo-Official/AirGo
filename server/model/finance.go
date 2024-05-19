package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// 用户余额明细
type BalanceStatement struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"  gorm:"index"`
	ID        int64      `json:"id" gorm:"primaryKey"`

	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`        //充值=Recharge，支出=Expenditure，提现=Withdraw
	Type        string `json:"type"`         //增加=Plus，减少=Reduce
	Amount      string `json:"amount"`       //变化值
	FinalAmount string `json:"final_amount"` //最终金额
}

// 邀请佣金明细
type CommissionStatement struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id"      gorm:"primaryKey"`

	UserID int64 `json:"user_id"`

	OrderUserID    int64   `json:"order_user_id"`
	OrderUserName  string  `json:"order_user_name"`
	OrderID        int64   `json:"order_id"`
	Subject        string  `json:"subject"`
	TotalAmount    string  `json:"total_amount"`
	CommissionRate float64 `json:"commission_rate"` //佣金率
	Commission     string  `json:"commission"`      //佣金
	IsWithdrew     bool    `json:"is_withdrew"`     //是否已提现
}

type FinanceSummary struct {
	TotalInvitation         string `json:"total_invitation"`          //总邀请人数
	TotalCommissionAmount   string `json:"total_commission_amount"`   //总佣金
	PendingWithdrawalAmount string `json:"pending_withdrawal_amount"` //待提现佣金
	TotalConsumptionAmount  string `json:"total_consumption_amount"`  //总消费
}

// 自定义gorm数据类型，注意：值接收器 和 指针接收器
type Jackpot []JackpotItem

type JackpotItem struct {
	Balance float64 `json:"balance"`
	Weight  int     `json:"weight"`
}

func (s *Jackpot) Scan(value any) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, s)
}
func (s Jackpot) Value() (driver.Value, error) {
	return json.Marshal(s)
}

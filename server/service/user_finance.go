package service

import (
	"errors"
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
	"strconv"
)

type Finance struct {
}

var FinanceSvc *Finance

func (f *Finance) NewBalanceStatement(params *model.BalanceStatement) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(params).Error
	})
}

func (f *Finance) NewCommissionStatement(params *model.CommissionStatement) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(params).Error
	})
}

func (f *Finance) SetWithdrew(uID int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Exec("UPDATE commission_statement SET is_withdrew = ? WHERE user_id = ?", 1, uID).Error
	})
}

func (f *Finance) GetWithdrawalAmount(uID int64) (totalCommission, currentCommission float64, err error) {
	var list []model.CommissionStatement
	err = global.DB.
		Model(model.CommissionStatement{}).
		Where(&model.CommissionStatement{UserID: uID}).
		Find(&list).Error
	if err != nil {
		return 0, 0, err
	}

	for k, _ := range list {
		com, err := strconv.ParseFloat(list[k].Commission, 64)
		if err != nil {
			continue
		}
		totalCommission += com
		if !list[k].IsWithdrew {
			currentCommission += com
		}
	}
	return totalCommission, currentCommission, nil

}

func (f *Finance) WithdrawToBalance(uID int64) error {
	_, amount, err := f.GetWithdrawalAmount(uID)
	if err != nil {
		return err
	}
	if amount < global.Server.Finance.WithdrawThreshold {
		return errors.New(constant.ERROR_COMMISSION_IS_NOT_ENOUGH)
	}
	user, err := UserSvc.FirstUser(&model.User{ID: uID})
	if err != nil {
		return err
	}
	user.Balance += amount
	endAmount := fmt.Sprintf("%.2f", user.Balance)
	user.Balance, err = strconv.ParseFloat(endAmount, 64)
	if err != nil {
		return err
	}
	//将明细标记为已提现
	err = f.SetWithdrew(uID)
	if err != nil {
		return err
	}
	//保存用户信息
	err = UserSvc.SaveUser(user)
	if err != nil {
		return err
	}
	//处理余额流水
	return f.NewBalanceStatement(&model.BalanceStatement{
		UserID:      uID,
		Title:       constant.BALANCE_STATEMENT_TITLE_WITHDRAW,
		Type:        constant.BALANCE_STATEMENT_TYPE_PLUS,
		Amount:      fmt.Sprintf("%.2f", amount),
		FinalAmount: endAmount,
	})
}

func (f *Finance) GetCommissionSummary(uID int64) (*model.FinanceSummary, error) {
	//总消费
	totalConsumption, err := OrderSvc.GetUserTotalConsumptionAmount(uID)
	if err != nil {
		return nil, err
	}
	//总佣金, 待提现佣金
	totalCommission, currentCommission, err := f.GetWithdrawalAmount(uID)
	if err != nil {
		return nil, err
	}
	//邀请人数
	var TotalInvitation int64
	err = global.DB.
		Model(&model.User{}).
		Where(&model.User{ReferrerUserID: uID}).
		Count(&TotalInvitation).Error
	if err != nil {
		return nil, err
	}
	return &model.FinanceSummary{
		TotalInvitation:         fmt.Sprintf("%d", TotalInvitation),
		TotalCommissionAmount:   fmt.Sprintf("%.2f", totalCommission),
		PendingWithdrawalAmount: fmt.Sprintf("%.2f", currentCommission),
		TotalConsumptionAmount:  fmt.Sprintf("%.2f", totalConsumption),
	}, nil

}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Orders struct {
	Id             int            `gorm:"column:id;primaryKey;autoIncrement;" uri:"id" json:"id" form:"id" comment:"自增主键"`                   //自增主键
	OrderNo        string         `gorm:"column:order_no" json:"order_no" form:"order_no" comment:"主订单编号（唯一）"`                               //主订单编号（唯一）
	UserId         string         `gorm:"column:user_id" json:"user_id" form:"user_id" comment:"用户唯一标识"`                                     //用户唯一标识
	CardId         int            `gorm:"column:card_id" json:"card_id" form:"card_id" comment:"权益卡标识"`                                      //权益卡标识
	CardName       string         `gorm:"column:card_name" json:"card_name" form:"card_name" comment:"权益卡名称"`                                //权益卡名称
	TotalPeriods   int            `gorm:"column:total_periods" json:"total_periods" form:"total_periods" comment:"总期数"`                      //总期数
	TotalAmount    int            `gorm:"column:total_amount" json:"total_amount" form:"total_amount" comment:"订单总金额（单位分）"`                  //订单总金额（单位分）
	PeriodAmount   int            `gorm:"column:period_amount" json:"period_amount" form:"period_amount" comment:"每期金额（单位分）"`                //每期金额（单位分）
	ContractStatus int            `gorm:"column:contract_status" json:"contract_status" form:"contract_status" comment:"签约状态：1=签约成功;2=解约成功"` //签约状态：1=签约成功;2=解约成功
	SignTime       time.Time      `gorm:"column:sign_time" json:"sign_time" form:"sign_time" comment:"签约时间"`                                 //签约时间
	CancelTime     time.Time      `gorm:"column:cancel_time" json:"cancel_time" form:"cancel_time" comment:"解约时间"`                           //解约时间
	StartTime      time.Time      `gorm:"column:start_time" json:"start_time" form:"start_time" comment:"整个权益生效开始时间（第1期开始）"`                 //整个权益生效开始时间（第1期开始）
	EndTime        time.Time      `gorm:"column:end_time" json:"end_time" form:"end_time" comment:"整个权益生效结束时间（最后一期结束）"`                      //整个权益生效结束时间（最后一期结束）
	Status         int            `gorm:"column:status" json:"status" form:"status" comment:"主订单状态：0=进行中,1=已完成,2=已取消,3=已过期"`                 //主订单状态：0=进行中,1=已完成,2=已取消,3=已过期
	FailureReason  string         `gorm:"column:failure_reason" json:"failure_reason" form:"failure_reason" comment:"整体失败原因（如签约失败）"`         //整体失败原因（如签约失败）
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at" form:"created_at" comment:"创建时间"`                              //创建时间
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at" form:"updated_at" comment:"更新时间"`                              //更新时间
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-" form:"deleted_at" comment:"软删除标记"`                                      //软删除标记
}

func (Orders) TableName() string {
	return "benefit.orders"
}

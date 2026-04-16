package models

import (
	"time"
)

type Orders struct {
	Id                   int        `gorm:"column:id;primaryKey;autoIncrement;" uri:"id" json:"id" form:"id" comment:"主键ID"`                                  //主键ID
	AppCode              int        `gorm:"column:app_code" json:"app_code" form:"app_code" comment:"APP标识"`                                                  //APP标识
	ChannelCode          string     `gorm:"column:channel_code" json:"channel_code" form:"channel_code" comment:"渠道标识"`                                       //渠道标识
	Code                 string     `gorm:"column:code" json:"code" form:"code" comment:"编码"`                                                                 //编码
	Phone                string     `gorm:"column:phone"  json:"phone" form:"phone" comment:"手机号"`                                                            //手机号
	OrderNo              string     `gorm:"column:order_no" json:"order_no" form:"order_no" comment:"业务订单编号（平台侧唯一）"`                                          //业务订单编号（平台侧唯一）
	PayOrderNo           string     `gorm:"column:pay_order_no" json:"pay_order_no" form:"pay_order_no" comment:"接入方支付订单编号"`                                  //接入方支付订单编号
	BillNo               string     `gorm:"column:bill_no" json:"bill_no" form:"bill_no" comment:"账单编号"`                                                      //账单编号
	BizOrderNo           string     `gorm:"column:biz_order_no" json:"biz_order_no" form:"biz_order_no" comment:"交易流水号"`                                      //交易流水号
	UserId               int        `gorm:"column:user_id" json:"user_id" form:"user_id" comment:"用户ID"`                                                      //用户ID
	OpenUserId           string     `gorm:"column:open_user_id" json:"open_user_id" form:"open_user_id" comment:"第三方用户ID"`                                    //第三方用户ID
	LoanUserId           int        `gorm:"column:loan_user_id" json:"loan_user_id" form:"loan_user_id" comment:"贷款信息标识"`                                     //贷款信息标识
	RewardProductId      int        `gorm:"column:reward_product_id" json:"reward_product_id" form:"reward_product_id" comment:"权益产品ID"`                      //权益产品ID
	ProductCode          string     `gorm:"column:product_code" json:"product_code" form:"product_code" comment:"卡标识或者对应权益包编码"`                               //卡标识或者对应权益包编码
	ProductName          string     `gorm:"column:product_name" json:"product_name" form:"product_name" comment:"权益产品名称"`                                     //权益产品名称
	PlanPeriods          int        `gorm:"column:plan_periods" json:"plan_periods" form:"plan_periods" comment:"权益计划期数"`                                     //权益计划期数
	Periods              int        `gorm:"column:periods" json:"periods" form:"periods" comment:"权益实际支付期数"`                                                  //权益实际支付期数
	Amount               int        `gorm:"column:amount" json:"amount" form:"amount" comment:"金额（单位：分）"`                                                     //金额（单位：分）
	PayAmount            int        `gorm:"column:pay_amount" json:"pay_amount" form:"pay_amount" comment:"支付金额（单位：分）"`                                       //支付金额（单位：分）
	PayChannel           string     `gorm:"column:pay_channel" json:"pay_channel" form:"pay_channel" comment:"支付通道标识"`                                        //支付通道标识
	PayStatus            int        `gorm:"column:pay_status" json:"pay_status" form:"pay_status" comment:"支付状态:0=进行中;1=成功;2=失败;3=取消中;4=取消成功"`                //支付状态:0=进行中;1=成功;2=失败;3=取消中;4=取消成功
	PayTime              *time.Time `gorm:"column:pay_time" json:"pay_time" form:"pay_time" comment:"支付时间"`                                                   //支付时间
	RefundOrderNo        string     `gorm:"column:refund_order_no" json:"refund_order_no" form:"refund_order_no" comment:"退款订单号"`                             //退款订单号
	RefundAmount         int        `gorm:"column:refund_amount" json:"refund_amount" form:"refund_amount" comment:"退款金额（单位：分）"`                              //退款金额（单位：分）
	RefundStatus         int        `gorm:"column:refund_status" json:"refund_status" form:"refund_status" comment:"退款状态:1=全额退款;2=失败;3=部分退款"`                 //退款状态:1=全额退款;2=失败;3=部分退款
	RefundTime           *time.Time `gorm:"column:refund_time" json:"refund_time" form:"refund_time" comment:"退款时间"`                                          //退款时间
	FailReason           string     `gorm:"column:fail_reason" json:"fail_reason" form:"fail_reason" comment:"失败原因（支付或退款）"`                                   //失败原因（支付或退款）
	RefundReason         string     `gorm:"column:refund_reason" json:"refund_reason" form:"refund_reason" comment:"退款原因"`                                    //退款原因
	ContractStatus       int        `gorm:"column:contract_status" json:"contract_status" form:"contract_status" comment:"签约状态：1=签约成功;2=解约成功"`                //签约状态：1=签约成功;2=解约成功
	CancelTime           *time.Time `gorm:"column:cancel_time" json:"cancel_time" form:"cancel_time" comment:"解约时间"`                                          //解约时间
	StartTime            *time.Time `gorm:"column:start_time" json:"start_time" form:"start_time" comment:"权益开始时间"`                                           //权益开始时间
	EndTime              *time.Time `gorm:"column:end_time" json:"end_time" form:"end_time" comment:"权益结束时间"`                                                 //权益结束时间
	PlatformRefundAmount int        `gorm:"column:platform_refund_amount" json:"platform_refund_amount" form:"platform_refund_amount" comment:"平台退款金额（单位：分）"` //平台退款金额（单位：分）
	PlatformShareAmount  int        `gorm:"column:platform_share_amount" json:"platform_share_amount" form:"platform_share_amount" comment:"平台分账金额（单位：分）"`    //平台分账金额（单位：分）
	CreatedAt            time.Time  `gorm:"column:created_at" json:"created_at" form:"created_at" comment:"创建时间"`                                             //创建时间
	UpdatedAt            time.Time  `gorm:"column:updated_at" json:"updated_at" form:"updated_at" comment:"修改时间"`                                             //修改时间
}

func (Orders) TableName() string {
	return "admin.orders"
}

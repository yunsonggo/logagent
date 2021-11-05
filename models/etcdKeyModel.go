package models

type EtcdKey struct {
	Id int `orm:"auto;pk"`
	KeyName string `orm:"key_name" form:"key_name"`
}

func (ek *EtcdKey) TableName() string {
	return "etcd_key"
}

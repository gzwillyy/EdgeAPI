package models

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/types"
)

const (
	NodeIPAddressStateEnabled  = 1 // 已启用
	NodeIPAddressStateDisabled = 0 // 已禁用
)

type NodeIPAddressDAO dbs.DAO

func NewNodeIPAddressDAO() *NodeIPAddressDAO {
	return dbs.NewDAO(&NodeIPAddressDAO{
		DAOObject: dbs.DAOObject{
			DB:     Tea.Env,
			Table:  "edgeNodeIPAddresses",
			Model:  new(NodeIPAddress),
			PkName: "id",
		},
	}).(*NodeIPAddressDAO)
}

var SharedNodeIPAddressDAO = NewNodeIPAddressDAO()

// 启用条目
func (this *NodeIPAddressDAO) EnableAddress(id int64) (err error) {
	_, err = this.Query().
		Pk(id).
		Set("state", NodeIPAddressStateEnabled).
		Update()
	return err
}

// 禁用IP地址
func (this *NodeIPAddressDAO) DisableAddress(id int64) (err error) {
	_, err = this.Query().
		Pk(id).
		Set("state", NodeIPAddressStateDisabled).
		Update()
	return err
}

// 禁用节点的所有的IP地址
func (this *NodeIPAddressDAO) DisableAllAddressesWithNodeId(nodeId int64) error {
	if nodeId <= 0 {
		return errors.New("invalid nodeId")
	}
	_, err := this.Query().
		Attr("nodeId", nodeId).
		Set("state", NodeIPAddressStateDisabled).
		Update()
	return err
}

// 查找启用中的IP地址
func (this *NodeIPAddressDAO) FindEnabledAddress(id int64) (*NodeIPAddress, error) {
	result, err := this.Query().
		Pk(id).
		Attr("state", NodeIPAddressStateEnabled).
		Find()
	if result == nil {
		return nil, err
	}
	return result.(*NodeIPAddress), err
}

// 根据主键查找名称
func (this *NodeIPAddressDAO) FindAddressName(id int64) (string, error) {
	return this.Query().
		Pk(id).
		Result("name").
		FindStringCol("")
}

// 创建IP地址
func (this *NodeIPAddressDAO) CreateAddress(nodeId int64, name string, ip string) (addressId int64, err error) {
	op := NewNodeIPAddressOperator()
	op.NodeId = nodeId
	op.Name = name
	op.IP = ip
	op.State = NodeIPAddressStateEnabled
	_, err = this.Save(op)
	if err != nil {
		return 0, err
	}
	return types.Int64(op.Id), nil
}

// 修改IP地址
func (this *NodeIPAddressDAO) UpdateAddress(addressId int64, name string, ip string) (err error) {
	if addressId <= 0 {
		return errors.New("invalid addressId")
	}

	op := NewNodeIPAddressOperator()
	op.Id = addressId
	op.Name = name
	op.IP = ip
	_, err = this.Save(op)
	return err
}

// 修改IP地址所属节点
func (this *NodeIPAddressDAO) UpdateAddressNodeId(addressId int64, nodeId int64) error {
	_, err := this.Query().
		Pk(addressId).
		Set("nodeId", nodeId).
		Set("state", NodeIPAddressStateEnabled). // 恢复状态
		Update()
	return err
}

// 查找某个节点所有的IP地址
func (this *NodeIPAddressDAO) FindAllEnabledAddressesWithNode(nodeId int64) (result []*NodeIPAddress, err error) {
	_, err = this.Query().
		Attr("nodeId", nodeId).
		State(NodeIPAddressStateEnabled).
		Desc("order").
		AscPk().
		Slice(&result).
		FindAll()
	return
}
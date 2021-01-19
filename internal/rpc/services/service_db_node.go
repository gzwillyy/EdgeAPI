package services

import (
	"context"
	"github.com/TeaOSLab/EdgeAPI/internal/db/models"
	rpcutils "github.com/TeaOSLab/EdgeAPI/internal/rpc/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/types"
)

// 数据库节点相关服务
type DBNodeService struct {
	BaseService
}

// 创建数据库节点
func (this *DBNodeService) CreateDBNode(ctx context.Context, req *pb.CreateDBNodeRequest) (*pb.CreateDBNodeResponse, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	nodeId, err := models.SharedDBNodeDAO.CreateDBNode(tx, req.IsOn, req.Name, req.Description, req.Host, req.Port, req.Database, req.Username, req.Password, req.Charset)
	if err != nil {
		return nil, err
	}
	return &pb.CreateDBNodeResponse{NodeId: nodeId}, nil
}

// 修改数据库节点
func (this *DBNodeService) UpdateDBNode(ctx context.Context, req *pb.UpdateDBNodeRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	err = models.SharedDBNodeDAO.UpdateNode(tx, req.NodeId, req.IsOn, req.Name, req.Description, req.Host, req.Port, req.Database, req.Username, req.Password, req.Charset)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// 删除节点
func (this *DBNodeService) DeleteDBNode(ctx context.Context, req *pb.DeleteDBNodeRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	err = models.SharedDBNodeDAO.DisableDBNode(tx, req.NodeId)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// 计算可用的数据库节点数量
func (this *DBNodeService) CountAllEnabledDBNodes(ctx context.Context, req *pb.CountAllEnabledDBNodesRequest) (*pb.RPCCountResponse, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	count, err := models.SharedDBNodeDAO.CountAllEnabledNodes(tx)
	if err != nil {
		return nil, err
	}
	return this.SuccessCount(count)
}

// 列出单页的数据库节点
func (this *DBNodeService) ListEnabledDBNodes(ctx context.Context, req *pb.ListEnabledDBNodesRequest) (*pb.ListEnabledDBNodesResponse, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	nodes, err := models.SharedDBNodeDAO.ListEnabledNodes(tx, req.Offset, req.Size)
	if err != nil {
		return nil, err
	}

	result := []*pb.DBNode{}
	for _, node := range nodes {
		result = append(result, &pb.DBNode{
			Id:          int64(node.Id),
			Name:        node.Name,
			Description: node.Description,
			IsOn:        node.IsOn == 1,
			Host:        node.Host,
			Port:        types.Int32(node.Port),
			Database:    node.Database,
			Username:    node.Username,
			Password:    node.Password,
			Charset:     node.Charset,
		})
	}
	return &pb.ListEnabledDBNodesResponse{Nodes: result}, nil
}

// 根据ID查找可用的数据库节点
func (this *DBNodeService) FindEnabledDBNode(ctx context.Context, req *pb.FindEnabledDBNodeRequest) (*pb.FindEnabledDBNodeResponse, error) {
	// 校验请求
	_, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin)
	if err != nil {
		return nil, err
	}

	tx := this.NullTx()

	node, err := models.SharedDBNodeDAO.FindEnabledDBNode(tx, req.NodeId)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return &pb.FindEnabledDBNodeResponse{Node: nil}, nil
	}
	return &pb.FindEnabledDBNodeResponse{Node: &pb.DBNode{
		Id:          int64(node.Id),
		Name:        node.Name,
		Description: node.Description,
		IsOn:        node.IsOn == 1,
		Host:        node.Host,
		Port:        types.Int32(node.Port),
		Database:    node.Database,
		Username:    node.Username,
		Password:    node.Password,
		Charset:     node.Charset,
	}}, nil
}
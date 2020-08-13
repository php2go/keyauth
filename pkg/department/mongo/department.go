package mongo

import (
	"context"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/types/ftime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/keyauth/common"
	"github.com/infraboard/keyauth/pkg/department"
)

func (s *service) QueryDepartment(req *department.QueryDepartmentRequest) (
	*department.Set, error) {
	r := newQueryDepartmentRequest(req)
	resp, err := s.col.Find(context.TODO(), r.FindFilter(), r.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find department error, error is %s", err)
	}

	set := department.NewDepartmentSet(req.PageRequest)
	// 循环
	for resp.Next(context.TODO()) {
		ins := department.NewDefaultDepartment()
		if err := resp.Decode(ins); err != nil {
			return nil, exception.NewInternalServerError("decode department error, error is %s", err)
		}

		set.Add(ins)
	}

	// count
	count, err := s.col.CountDocuments(context.TODO(), r.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get department count error, error is %s", err)
	}
	set.Total = count

	return set, nil
}

func (s *service) DescribeDepartment(req *department.DescribeDeparmentRequest) (
	*department.Department, error) {
	r, err := newDescribeDepartment(req)
	if err != nil {
		return nil, err
	}

	ins := department.NewDefaultDepartment()
	if err := s.col.FindOne(context.TODO(), r.FindFilter()).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("department %s not found", req)
		}

		return nil, exception.NewInternalServerError("find department %s error, %s", req.ID, err)
	}

	return ins, nil
}

func (s *service) CreateDepartment(req *department.CreateDepartmentRequest) (
	*department.Department, error) {
	ins, err := department.NewDepartment(req, s)
	if err != nil {
		return nil, err
	}

	count, err := s.counter.GetNextSequenceValue(ins.CounterKey())
	if err != nil {
		return nil, err
	}
	ins.Number = count.Value

	if _, err := s.col.InsertOne(context.TODO(), ins); err != nil {
		return nil, exception.NewInternalServerError("inserted department(%s) document error, %s",
			ins.Name, err)
	}

	return ins, nil
}

func (s *service) DeleteDepartment(id string) error {
	result, err := s.col.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return exception.NewInternalServerError("delete department(%s) error, %s", id, err)
	}

	if result.DeletedCount == 0 {
		return exception.NewNotFound("department %s not found", id)
	}

	return nil
}

func (s *service) UpdateDepartment(req *department.UpdateDepartmentRequest) (*department.Department, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate update department error, %s", err)
	}

	dp, err := s.DescribeDepartment(department.NewDescriptDepartmentRequestWithID(req.ID))
	if err != nil {
		return nil, err
	}
	switch req.UpdateMode {
	case common.PutUpdateMode:
		*dp.CreateDepartmentRequest = *req.CreateDepartmentRequest
	case common.PatchUpdateMode:
		dp.CreateDepartmentRequest.Patch(req.CreateDepartmentRequest)
	default:
		return nil, exception.NewBadRequest("unknown update mode: %s", req.UpdateMode)
	}

	dp.UpdateAt = ftime.Now()
	_, err = s.col.UpdateOne(context.TODO(), bson.M{"_id": dp.ID}, bson.M{"$set": dp})
	if err != nil {
		return nil, exception.NewInternalServerError("update domain(%s) error, %s", dp.Name, err)
	}

	return dp, nil
}

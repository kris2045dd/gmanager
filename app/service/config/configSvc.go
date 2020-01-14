package config

import (
	"errors"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"gmanager/app/model/config"
	"gmanager/utils/base"
)

// 请求参数
type Request struct {
	config.Entity
}

func GetById(id int64) (*config.Entity, error) {
	if id <= 0 {
		glog.Error(" get id error")
		return new(config.Entity), errors.New("参数不合法")
	}

	return config.Model.FindOne(" id = ?", id)
}

func GetOne(form *base.BaseForm) (*config.Entity, error) {
	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["id"] != "" {
		where += " and id = ? "
		params = append(params, gconv.Int(form.Params["id"]))
	}

	return config.Model.FindOne(where, params)
}

func List(form *base.BaseForm) ([]*config.Entity, error) {
	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["name"] != "" {
		where += " and name like ? "
		params = append(params, "%"+form.Params["name"]+"%")
	}
	if form.Params != nil && form.Params["parentId"] != "" {
		where += " and parent_id = ? "
		params = append(params, gconv.Int(form.Params["parentId"]))
	}

	return config.Model.Order(form.OrderBy).FindAll(where, params)
}

func Delete(id int64) (int64, error) {
	if id <= 0 {
		glog.Error("delete id error")
		return 0, errors.New("参数不合法")
	}

	r, err := config.Model.Delete(" id = ?", id)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

func Update(request *Request) (int64, error) {
	entity := (*config.Entity)(nil)
	err := gconv.StructDeep(request.Entity, &entity)
	if err != nil {
		return 0, errors.New("数据错误")
	}

	if entity.Id <= 0 {
		glog.Error("update id error")
		return 0, errors.New("参数不合法")
	}

	r, err := config.Model.OmitEmpty().Where(" id = ?", entity.Id).Update(entity)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

func Insert(request *Request) (int64, error) {
	entity := (*config.Entity)(nil)
	err := gconv.StructDeep(request.Entity, &entity)
	if err != nil {
		return 0, errors.New("数据错误")
	}

	if entity.Id > 0 {
		glog.Error("insert id error")
		return 0, errors.New("参数不合法")
	}

	r, err := config.Model.Insert(entity)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

// 分页查询
func Page(form *base.BaseForm) ([]config.Entity, error) {
	if form.Page <= 0 || form.Rows <= 0 {
		glog.Error("page param error", form.Page, form.Rows)
		return []config.Entity{}, nil
	}

	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil {
		if form.Params["name"] != "" {
			where += " and t.name like ? "
			params = append(params, "%"+form.Params["name"]+"%")
		}
		if form.Params["key"] != "" {
			where += " and t.key like ? "
			params = append(params, "%"+form.Params["key"]+"%")
		}
		if gconv.Int(form.Params["parentId"]) > 0 {
			where += " and t.parent_id = ? "
			params = append(params, gconv.Int(form.Params["parentId"]))
		}
	}

	num, err := config.Model.As("t").FindCount(where, params)
	form.TotalSize = num
	form.TotalPage = num / form.Rows

	if err != nil {
		glog.Error("page count error", err)
		return []config.Entity{}, err
	}

	// 没有数据直接返回
	if num == 0 {
		form.TotalPage = 0
		form.TotalSize = 0
		return []config.Entity{}, err
	}

	var resData []config.Entity
	dbModel := config.Model.As("t").Fields(config.Model.Columns() + ",su1.real_name as updateName,su2.real_name as createName")
	dbModel = dbModel.LeftJoin("sys_user su1", " t.update_id = su1.id ")
	dbModel = dbModel.LeftJoin("sys_user su2", " t.update_id = su2.id ")
	err = dbModel.Where(where, params).Order(form.OrderBy).Page(form.Page, form.Rows).M.Structs(&resData)
	if err != nil {
		glog.Error("page list error", err)
		return []config.Entity{}, err
	}

	return resData, nil
}

package funcall

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

//FunCall 函数映射
type FunCall struct {
	funcMap map[FuncCallEnum]reflect.Value
}

var (
	funcall *FunCall
)

func init() {
	funcall = NewFuncCall()
}

//NewFuncCall 获取funcall
func NewFuncCall() *FunCall {
	return &FunCall{
		funcMap: make(map[FuncCallEnum]reflect.Value),
	}
}

//BindFunc bind call
func BindFunc(tp FuncCallEnum, fc interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(strconv.Itoa(int(tp)) + " bind func error .")
		}
	}()
	if _, ok := funcall.funcMap[tp]; ok {
		err = errors.New(strconv.Itoa(int(tp)) + " had binded yet .")
		return
	}
	va := reflect.ValueOf(fc)
	va.Type().NumIn()
	funcall.funcMap[tp] = va
	return
}

//Call bind call
func Call(tp FuncCallEnum, params ...interface{}) (result []reflect.Value, err error) {
	defer func() {
		if e := recover(); e != nil {
			//err = errors.New(strconv.Itoa(int(tp)) + " is not callable.")
			fmt.Println("funccall err ", e)
		}
	}()
	if _, ok := funcall.funcMap[tp]; !ok {
		err = errors.New(strconv.Itoa(int(tp)) + " does not exist.")
		return
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = funcall.funcMap[tp].Call(in)
	return
}

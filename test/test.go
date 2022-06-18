package main

import (
	"reflect"
	"strings"
	"strconv"
	"fmt"
)

// Reflect String to Base Type : int
func RefStrToInt(dst, src interface{}) {
	// Judge data Kind
	if reflect.TypeOf(dst).Kind() == reflect.Ptr && reflect.TypeOf(src).Kind() == reflect.Ptr {
		dstType, dstVal := reflect.TypeOf(dst).Elem(), reflect.ValueOf(dst).Elem()
		srcType, srcVal := reflect.TypeOf(src).Elem(), reflect.ValueOf(src).Elem()
		// Traverse
		for i := 0; i < srcType.NumField(); i++ {
			// Get source StructField
			srcSF := srcType.Field(i)
			for j := 0; j < dstType.NumField(); j++ {
				// Obtain destination StructField
				dstSF := dstType.Field(j)
				// Acquire destination Kind and source Kind
				dstKind, srcKind := dstSF.Type.Kind(), srcSF.Type.Kind()
				// Judge FieldName and Kind
				if strings.EqualFold(srcSF.Name, dstSF.Name) && dstKind == srcKind {
					// Src Assignment to Dst
					dstVal.Field(j).Set(srcVal.Field(i))
					break
				} else if strings.EqualFold(srcSF.Name, dstSF.Name) && srcKind == reflect.String {
					// destination kind is Int,Int32,Int64
					if dstKind == reflect.Int || dstKind == reflect.Int64 || dstKind == reflect.Int32 {
						// String convent to Int64
						atoi, _ := strconv.Atoi(srcVal.Field(i).Interface().(string))
						// src Assignment to dst
						dstVal.Field(j).SetInt(int64(atoi))
						break
					}
					// destination kind is Uint,Uint32,Uint64
					if dstKind == reflect.Uint || dstKind == reflect.Uint32 || dstKind == reflect.Uint64 {
						// String convent to Uint64
						atoi, _ := strconv.Atoi(srcVal.Field(i).Interface().(string))
						// src Assignment to dst
						dstVal.Field(j).SetUint(uint64(atoi))
						break
					}
				}
			}
		}
	}
}

type DST struct {
	Name  string
	Age   int
	Level int32
	Grade int64
	Score uint
	Month uint32
	Year  uint64
}

type SRC struct {
	Level string
	Grade string
	Month string
	Year  string
	Name  string
	Age   string
	Score string
}

func main() {
	dst, src := DST{}, SRC{
		Name:  "zs",
		Age:   "18",
		Score: "68",
		Level: "33",
		Grade: "66",
		Month: "8",
		Year:  "2018",
	}
	RefStrToInt(&dst, &src)
	fmt.Println(dst)
}
